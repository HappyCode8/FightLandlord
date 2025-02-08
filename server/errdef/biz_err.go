package errdef

type Error struct {
	Code int
	Msg  string
	Exit bool
}

func (e Error) Error() string {
	return e.Msg
}

func NewErr(code int, exit bool, msg string) Error {
	return Error{Code: code, Exit: exit, Msg: msg}
}

var (
	ErrorsExist        = NewErr(1, true, "Exist. ")
	ErrorsChanClosed   = NewErr(1, true, "Chan closed. ")
	ErrorsTimeout      = NewErr(1, false, "Timeout. ")
	ErrorsInputInvalid = NewErr(1, false, "Input invalid. ")
	ErrorsAuthFail     = NewErr(1, true, "Auth fail. ")
	ErrorsRoomInvalid  = NewErr(1, true, "Room invalid. ")

	ErrorsRoomPlayersIsFull = NewErr(1, false, "Room players is fill. ")

	ErrorsJoinFailForRoomRunning = NewErr(1, false, "Join fail, room is running. ")

	ErrorsPokersFacesInvalid = NewErr(1, false, "Pokers faces invalid. ")
	ErrorsHaveToPlay         = NewErr(1, false, "Have to play. ")
)
