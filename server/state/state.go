package state

import (
	"errors"
	"log"
	"server/consts"
	"server/database"
	"server/util"
	"slices"
	"strings"
)

var states = map[consts.StateID]State{}

func init() {
	register(consts.StateWelcome, &welcome{})
	register(consts.StateHome, &home{})
	register(consts.StateCreate, &create{})
	register(consts.StateJoin, &join{})
	register(consts.StateWaiting, &waiting{})
	register(consts.StateGame, &Game{})
}

func register(id consts.StateID, state State) {
	states[id] = state
}

type State interface {
	Next(player *database.Player) (consts.StateID, error)
	Exit(player *database.Player) consts.StateID
}

func Run(player *database.Player) {
	player.State(consts.StateWelcome)
	defer func() {
		if err := recover(); err != nil {
			util.PrintStackTrace(err)
		}
		log.Println("player %s state machine break up.\n", player)
	}()
	for {
		// 获取状态对应的处理对象
		state := states[player.GetState()]
		stateId, err := state.Next(player)
		if err != nil {
			var err1 consts.Error
			if errors.As(err, &err1) {
				if err1.Exit {
					stateId = state.Exit(player)
				}
			}
		}
		if stateId > 0 {
			player.State(stateId)
		}
	}
}

func isExitSignal(signal string) bool {
	signal = strings.ToLower(signal)
	return isXSignal(signal, "exit", "e")
}

func isLsSignal(signal string) bool {
	return isXSignal(signal, "ls")
}

func isXSignal(signal string, x ...string) bool {
	signal = strings.ToLower(signal)
	return slices.Contains(x, signal)
}
