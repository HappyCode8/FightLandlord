package state

import (
	"log"
	"server/consts"
	"server/database"
	"server/util"
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
		state := states[player.GetState()]
		// 以欢迎页为例，输出一段欢迎信息，然后输出homeid
		stateId, err := state.Next(player)
		if err != nil {
			if err1, ok := err.(consts.Error); ok {
				if err1.Exit {
					stateId = state.Exit(player)
				}
			} else {
				log.Println(err)
				state.Exit(player)
				break
			}
		}
		if stateId > 0 {
			player.State(stateId)
		}
	}
}

func isExit(signal string) bool {
	signal = strings.ToLower(signal)
	return isX(signal, "exit", "e")
}

func isLs(signal string) bool {
	return isX(signal, "ls")
}

func isX(signal string, x ...string) bool {
	signal = strings.ToLower(signal)
	for _, v := range x {
		if v == signal {
			return true
		}
	}
	return false
}
