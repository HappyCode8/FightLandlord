package state

import (
	"bytes"
	"fmt"
	"server/consts"
	"server/database"
)

type welcome struct{}

func (*welcome) Next(player *database.Player) (consts.StateID, error) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Hi %s, Welcome to ratel online!\n", player.Name))
	err := player.WriteString(buf.String())
	if err != nil {
		return 0, player.WriteError(err)
	}
	return consts.StateHome, nil
}

func (*welcome) Exit(player *database.Player) consts.StateID {
	return 0
}
