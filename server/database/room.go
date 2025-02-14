package database

import (
	"log"
	"sync"
	"time"
)

type Room struct {
	sync.Mutex // 锁

	ID         int64     `json:"id"`         // 房间id
	Game       *Game     `json:"gameId"`     // 房间持有的游戏
	State      int       `json:"service"`    // 房间状态
	Players    int       `json:"players"`    // 房间人数
	Creator    int64     `json:"creator"`    // 房间创建者
	ActiveTime time.Time `json:"activeTime"` // 房间活跃时间
	MaxPlayers int       `json:"maxPlayers"` // 房间最大人数
}

func roomCancel(room *Room) {
	if room.ActiveTime.Add(24 * time.Hour).Before(time.Now()) {
		log.Printf("room %d is timeout 24 hours, removed.\n", room.ID)
		deleteRoom(room)
		return
	}
	living := false
	playerIds := getRoomPlayers(room.ID)
	for id := range playerIds {
		if getPlayer(id).online {
			living = true
			break
		}
	}
	if !living {
		log.Printf("room %d is not living, removed.\n", room.ID)
		deleteRoom(room)
	}
}
