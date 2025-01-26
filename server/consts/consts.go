package consts

import "time"

const (
	MaxPacketSize = 65536

	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"

	MaxPlayers = 3

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

// todo: Exit的作用？
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
	RoomStates               = map[int]string{
		RoomStateWaiting: "Waiting",
		RoomStateRunning: "Running",
	}
)

type FacesType int

const (
	_                   FacesType = iota
	FacesBomb                     = 1 //炸弹
	FacesSingle                   = 2 //单牌
	FacesDouble                   = 3 //对子
	FacesTriple                   = 4 //三张
	FacesUnion3                   = 5 //三带一
	FacesUnion4                   = 6 //四带二
	FacesStraight                 = 7 //顺子 or 连对
	FacesUnion3Straight           = 8 //飞机
	FacesInvalid                  = 9
)
