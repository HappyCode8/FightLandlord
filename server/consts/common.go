package consts

import "time"

const (
	MaxPacketSize = 65536
	MaxPlayers    = 3

	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"

	RoomStateWaiting = 1
	RoomStateRunning = 2
)

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

var (
	RoomStates = map[int]string{
		RoomStateWaiting: "Waiting",
		RoomStateRunning: "Running",
	}
)
