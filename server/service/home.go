package service

import (
	"server/consts"
	"server/database"
	"server/errdef"
)

type home struct{}

func (*home) Next(player *database.Player) (consts.StateID, error) {
	chooseStr := "1.Join\n2.New\n"
	err := player.WriteString(chooseStr)
	if err != nil {
		return 0, player.WriteError(err)
	}
	selected, err := player.AskForInt()
	if err != nil {
		return 0, player.WriteError(err)
	}
	switch selected {
	case 1:
		return consts.StateJoin, nil
	case 2:
		return consts.StateCreate, nil
	}
	return 0, player.WriteError(errdef.ErrorsInputInvalid)
}

func (*home) Exit(player *database.Player) consts.StateID {
	return 0
}
