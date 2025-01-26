package database

import (
	"log"
	"server/consts"
	"server/model"
	"server/network"
	"server/util"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

var roomIds int64 = 0
var players = make(map[int64]*Player) // 存储连接过服务器的全部用户
var connPlayers = make(map[int64]*Player)
var rooms = make(map[int64]*Room)
var roomPlayers = make(map[int64]map[int64]bool)

func Connected(conn *network.Conn, info *model.AuthInfo) *Player {
	player := &Player{
		ID:   info.ID,
		IP:   conn.IP(),
		Name: info.Name,
	}
	player.Conn(conn)               // 初始化play对象
	players[info.ID] = player       // 写入用户池
	connPlayers[conn.ID()] = player // 写入连接用户池
	return player
}

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
		ID:             atomic.AddInt64(&roomIds, 1),
		State:          consts.RoomStateWaiting,
		Creator:        creator,
		ActiveTime:     time.Now(),
		MaxPlayers:     consts.MaxPlayers,
		EnableLandlord: true,
		EnableChat:     true,
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

	// 房间状态检查
	if room.State == consts.RoomStateRunning {
		return consts.ErrorsJoinFailForRoomRunning
	}

	//房间人数检查
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

func getPlayer(playerId int64) *Player {
	if v, ok := players[playerId]; ok {
		return v
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

func BroadcastChat(player *Player, msg string, exclude ...int64) {
	log.Println("chat msg, player %s[%d] %s say: %s\n", player.Name, player.ID, player.IP, msg)
	Broadcast(player.RoomID, msg, exclude...)
}

func RoomPlayers(roomId int64) map[int64]bool {
	return getRoomPlayers(roomId)
}

type Game struct {
	Room        *Room                  `json:"room"`
	Players     []int64                `json:"players"`
	Groups      map[int64]int          `json:"groups"`
	States      map[int64]chan int     `json:"states"`
	Pokers      map[int64]model.Pokers `json:"pokers"`
	Universals  []int                  `json:"universals"`
	Decks       int                    `json:"decks"`
	Additional  model.Pokers           `json:"pocket"`
	Multiple    int                    `json:"multiple"`
	FirstPlayer int64                  `json:"firstPlayer"`
	LastPlayer  int64                  `json:"lastPlayer"`
	Robs        []int64                `json:"robs"`
	FirstRob    int64                  `json:"firstRob"`
	LastRob     int64                  `json:"lastRob"`
	FinalRob    bool                   `json:"finalRob"`
	LastPokers  model.Pokers           `json:"lastPokers"`
	Mnemonic    map[int]int            `json:"mnemonic"`
	//Skills      map[int64]int           `json:"skills"`
	PlayTimes   map[int64]int           `json:"playTimes"`
	PlayTimeOut map[int64]time.Duration `json:"playTimeOut"`
	//Rules       poker.Rules             `json:"rules"`
	Discards model.Pokers `json:"discards"` // 废弃的牌
}

func (game *Game) Clean() {
	if game != nil {
		for _, state := range game.States {
			close(state)
		}
	}
}

func (game *Game) Start() {

}

func (g Game) NextPlayer(curr int64) int64 {
	idx := util.IndexOf(g.Players, curr)
	return g.Players[(idx+1)%len(g.Players)]
}

func (g Game) PrevPlayer(curr int64) int64 {
	idx := util.IndexOf(g.Players, curr)
	return g.Players[(idx+len(g.Players))%len(g.Players)]
}

func (g Game) IsTeammate(player1, player2 int64) bool {
	return g.Groups[player1] == g.Groups[player2]
}

func (g Game) IsLandlord(playerId int64) bool {
	return g.Groups[playerId] == 1
}

func (g Game) Team(playerId int64) string {
	if !g.Room.EnableLandlord {
		return "team" + strconv.Itoa(g.Groups[playerId])
	} else {
		if !g.IsLandlord(playerId) {
			return "peasant"
		} else {
			return "landlord"
		}
	}
}
