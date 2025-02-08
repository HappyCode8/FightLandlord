package database

import (
	"log"
	"server/consts"
	"server/model"
	"server/network"
	"sort"
	"sync/atomic"
	"time"
)

var roomIds int64 = 0                            //房间id
var players = make(map[int64]*Player)            // <userId,player>
var connPlayers = make(map[int64]*Player)        //<connId,player>
var rooms = make(map[int64]*Room)                // <roomId,room>
var roomPlayers = make(map[int64]map[int64]bool) //<roomId,<playerId, bool>>

func Connected(conn *network.Conn, info *model.AuthInfo) *Player {
	player := &Player{
		// 客户端的id是时间戳
		ID:   info.ID,
		Name: info.Name,
		IP:   conn.IP(),
	}
	player.Conn(conn)               // 初始化player对象
	players[info.ID] = player       // 写入用户池
	connPlayers[conn.ID()] = player // 写入连接用户池
	return player
}

// GetRooms 获取所有房间，按照id排序
func GetRooms() []*Room {
	list := make([]*Room, 0)
	for _, room := range rooms {
		list = append(list, room)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	return list
}

func CreateRoom(creator int64) *Room {
	room := &Room{
		ID:         atomic.AddInt64(&roomIds, 1),
		State:      consts.RoomStateWaiting,
		Creator:    creator,
		ActiveTime: time.Now(),
		MaxPlayers: consts.MaxPlayers,
	}
	rooms[room.ID] = room
	roomPlayers[room.ID] = map[int64]bool{}
	return room
}

func JoinRoom(roomId, playerId int64) error {
	// 资源检查
	player := getPlayer(playerId)
	if player == nil {
		return consts.ErrorsExist
	}
	room := getRoom(roomId)
	if room == nil {
		return consts.ErrorsRoomInvalid
	}

	// 加锁防止并发异常
	room.Lock()
	defer room.Unlock()

	room.ActiveTime = time.Now()

	// 房间状态检查，已经运行的不能加入
	if room.State == consts.RoomStateRunning {
		return consts.ErrorsJoinFailForRoomRunning
	}

	// 房间人数检查，已经超过2人的不能加入
	if room.Players >= room.MaxPlayers {
		return consts.ErrorsRoomPlayersIsFull
	}

	playersIds := getRoomPlayers(roomId)
	if playersIds != nil {
		playersIds[playerId] = true
		room.Players++
		player.RoomID = roomId
	} else {
		deleteRoom(room)
		return consts.ErrorsRoomInvalid
	}
	return nil
}

func GetRoom(roomId int64) *Room {
	return getRoom(roomId)
}

func getRoom(roomId int64) *Room {
	if v, ok := rooms[roomId]; ok {
		return v
	}
	return nil
}

func getRoomPlayers(roomId int64) map[int64]bool {
	if v, ok := roomPlayers[roomId]; ok {
		return v
	}
	return nil
}

func deleteRoom(room *Room) {
	if room != nil {
		rooms[room.ID] = nil
		roomPlayers[room.ID] = nil
		if room.Game != nil {
			room.Game.Clean()
		}
	}
}

func broadcast(room *Room, msg string, exclude ...int64) {
	room.ActiveTime = time.Now()
	excludeSet := map[int64]bool{}
	for _, exc := range exclude {
		excludeSet[exc] = true
	}
	for playerId := range getRoomPlayers(room.ID) {
		if player := getPlayer(playerId); player != nil && !excludeSet[playerId] {
			_ = player.WriteString(">> " + msg)
		}
	}
}

func Broadcast(roomId int64, msg string, exclude ...int64) {
	room := getRoom(roomId)
	if room == nil {
		return
	}
	broadcast(room, msg, exclude...)
}

func LeaveRoom(roomId, playerId int64) {
	room := getRoom(roomId)
	if room != nil {
		room.Lock()
		defer room.Unlock()
		leaveRoom(room, getPlayer(playerId))
	}
}

func leaveRoom(room *Room, player *Player) {
	if room == nil || player == nil {
		return
	}
	room.ActiveTime = time.Now()
	playersIds := getRoomPlayers(room.ID)
	if _, ok := playersIds[player.ID]; ok {
		room.Players--
		player.RoomID = 0
		delete(playersIds, player.ID)
		if len(playersIds) > 0 && room.Creator == player.ID {
			for k := range playersIds {
				room.Creator = k
				break
			}
		}
	}
	if len(playersIds) == 0 {
		deleteRoom(room)
	}
}

func GetPlayer(playerId int64) *Player {
	return getPlayer(playerId)
}

func getPlayer(playerId int64) *Player {
	if v, ok := players[playerId]; ok {
		return v
	}
	return nil
}

func BroadcastChat(player *Player, msg string, exclude ...int64) {
	log.Printf("chat msg, player %s[%d] %s say: %s\n\n", player.Name, player.ID, player.IP, msg)
	Broadcast(player.RoomID, msg, exclude...)
}

func RoomPlayers(roomId int64) map[int64]bool {
	return getRoomPlayers(roomId)
}
