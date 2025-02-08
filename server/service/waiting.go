package service

import (
	"bytes"
	"errors"
	"fmt"
	"server/consts"
	"server/database"
	"server/errdef"
	"server/util"
	"strings"
	"time"
)

type waiting struct{}

func (s *waiting) Next(player *database.Player) (consts.StateID, error) {
	room := database.GetRoom(player.RoomID)
	if room == nil {
		return 0, errdef.ErrorsExist
	}
	access, err := waitingForStart(player, room)
	if err != nil {
		return 0, err
	}
	if access {
		return consts.StateGame, nil
	}
	return s.Exit(player), nil
}

func (*waiting) Exit(player *database.Player) consts.StateID {
	room := database.GetRoom(player.RoomID)
	if room != nil {
		isOwner := room.Creator == player.ID
		database.LeaveRoom(room.ID, player.ID)
		database.Broadcast(room.ID, fmt.Sprintf("%s exited room! room current has %d players\n", player.Name, room.Players))
		if isOwner {
			newOwner := database.GetPlayer(room.Creator)
			database.Broadcast(room.ID, fmt.Sprintf("%s become new owner\n", newOwner.Name))
		}
	}
	return consts.StateHome
}

func waitingForStart(player *database.Player, room *database.Room) (bool, error) {
	access := false
	//对局类别
	player.StartTransaction()
	defer player.StopTransaction()
	for {
		signal, err := player.AskForStringWithoutTransaction(time.Second)
		if err != nil && !errors.Is(err, errdef.ErrorsTimeout) {
			return access, err
		}
		if room.State == consts.RoomStateRunning {
			access = true
			break
		}
		signal = strings.ToLower(signal)
		// ls指令时列出当前的人员,start开始游戏，其余指令视为聊天广播
		if signal == "ls" {
			viewRoomPlayers(room, player)
		} else if signal == "start" && room.Creator == player.ID && room.Players == 3 {
			err = startGame(player, room)
			if err != nil {
				return access, err
			}
			access = true
			break
		} else if len(signal) > 0 {
			database.BroadcastChat(player, fmt.Sprintf("%s say: %s\n", player.Name, signal))
		}
	}
	return access, nil
}

func startGame(player *database.Player, room *database.Room) (err error) {
	room.Lock()
	defer room.Unlock()
	room.Game, err = InitGame(room)
	if err != nil {
		_ = player.WriteError(err)
		return err
	}
	room.State = consts.RoomStateRunning
	return nil
}

func viewRoomPlayers(room *database.Room, currPlayer *database.Player) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Room ID: %d\n", room.ID))
	buf.WriteString(fmt.Sprintf("%-20s%-10s\n", "Name", "Title"))
	for playerId := range database.RoomPlayers(room.ID) {
		player := database.GetPlayer(playerId)
		buf.WriteString(fmt.Sprintf("%-20s%-10s\n", player.Name, util.ChooseIf(playerId == room.Creator, "owner", "player")))
	}
	_ = currPlayer.WriteString(buf.String())
}
