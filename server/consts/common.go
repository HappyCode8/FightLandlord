package consts

const (
	MaxPacketSize = 65536
	MaxPlayers    = 3

	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"

	RoomStateWaiting = 1
	RoomStateRunning = 2
)

var (
	RoomStates = map[int]string{
		RoomStateWaiting: "Waiting",
		RoomStateRunning: "Running",
	}
)
