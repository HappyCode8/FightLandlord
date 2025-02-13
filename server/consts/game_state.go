package consts

import "time"

type StateID int

const (
	_ StateID = iota
	StateWelcome
	StateHome
	StateJoin
	StateCreate
	StateWaiting
	StateGame

	PlayTimeout = 40 * time.Second
)
