package state

import (
	"fmt"
	"server/consts"
	"server/database"
)

type create struct{}

func (*create) Next(player *database.Player) (consts.StateID, error) {
	// 创建房间
	room := database.CreateRoom(player.ID)
	err := player.WriteString(fmt.Sprintf("Create room successful, id : %d\n", room.ID))
	if err != nil {
		return 0, player.WriteError(err)
	}
	// 创建完以后加入
	err = database.JoinRoom(room.ID, player.ID)
	if err != nil {
		return 0, player.WriteError(err)
	}
	return consts.StateWaiting, nil
}

func (*create) Exit(_ *database.Player) consts.StateID {
	return consts.StateHome
}
