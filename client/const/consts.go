package consts

const (
	MaxPacketSize = 65536

	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"
)

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}
