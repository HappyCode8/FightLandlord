package database

import (
	"fmt"
	"log"
	"server/consts"
	"server/network"
	"server/protocol"
	"server/util"
	"strings"
	"sync"
	"time"
)

type Player struct {
	ID     int64  `json:"id"`
	IP     string `json:"ip"`
	Name   string `json:"name"`
	RoomID int64  `json:"roomId"`

	conn   *network.Conn
	data   chan *protocol.Packet
	read   bool
	state  consts.StateID
	online bool
}

func (p *Player) Conn(conn *network.Conn) {
	p.conn = conn
	p.data = make(chan *protocol.Packet, 8)
	p.online = true
}

func (p *Player) State(s consts.StateID) {
	p.state = s
}

func (p *Player) GetState() consts.StateID {
	return p.state
}

func (p *Player) WriteString(data string) error {
	time.Sleep(30 * time.Millisecond)
	return p.conn.Write(protocol.Packet{
		Body: []byte(data),
	})
}

func (p *Player) WriteObject(data interface{}) error {
	return p.conn.Write(protocol.Packet{
		Body: util.Marshal(data),
	})
}

func (p *Player) WriteError(err error) error {
	if err == consts.ErrorsExist {
		return err
	}
	return p.conn.Write(protocol.Packet{
		Body: []byte(err.Error() + "\n"),
	})
}

func (p *Player) AskForString(timeout ...time.Duration) (string, error) {
	packet, err := p.AskForPacket(timeout...)
	if err != nil {
		return "", err
	}
	return packet.String(), nil
}

func (p *Player) AskForInt(timeout ...time.Duration) (int, error) {
	packet, err := p.AskForPacket(timeout...)
	if err != nil {
		return 0, err
	}
	return packet.Int()
}

func (p *Player) AskForPacket(timeout ...time.Duration) (*protocol.Packet, error) {
	p.StartTransaction()
	defer p.StopTransaction()
	return p.askForPacket(timeout...)
}

func (p *Player) AskForStringWithoutTransaction(timeout ...time.Duration) (string, error) {
	packet, err := p.askForPacket(timeout...)
	if err != nil {
		return "", err
	}
	return packet.String(), nil
}

func (p *Player) StartTransaction() {
	p.read = true
	_ = p.WriteString(consts.IsStart)
}

func (p *Player) StopTransaction() {
	p.read = false
	_ = p.WriteString(consts.IsStop)
}

func (p *Player) askForPacket(timeout ...time.Duration) (*protocol.Packet, error) {
	var packet *protocol.Packet
	if len(timeout) > 0 {
		select {
		case packet = <-p.data:
		case <-time.After(timeout[0]):
			return nil, consts.ErrorsTimeout
		}
	} else {
		packet = <-p.data
	}
	if packet == nil {
		return nil, consts.ErrorsChanClosed
	}
	single := strings.ToLower(packet.String())
	if single == "exit" || single == "e" {
		return nil, consts.ErrorsExist
	}
	return packet, nil
}

type Room struct {
	sync.Mutex

	ID             int64     `json:"id"` // 房间id
	Game           RoomGame  `json:"gameId"`
	State          int       `json:"state"`      // 房间状态
	Players        int       `json:"players"`    // 房间人数
	Creator        int64     `json:"creator"`    // 房间创建者
	ActiveTime     time.Time `json:"activeTime"` // 房间活跃时间
	MaxPlayers     int       `json:"maxPlayers"` // 房间最大人数
	EnableChat     bool      `json:"enableChat"` // 房间是否开启聊天
	EnableLandlord bool      `json:"enableLandlord"`
}

type RoomGame interface {
	Clean()
}

func (p *Player) Offline() {
	p.online = false
	_ = p.conn.Close()
	close(p.data)
	room := getRoom(p.RoomID)
	if room != nil {
		room.Lock()
		defer room.Unlock()
		broadcast(room, fmt.Sprintf("%s lost connection! \n", p.Name))
		if room.State == consts.RoomStateWaiting {
			leaveRoom(room, p)
		}
		roomCancel(room)
	}
}

func roomCancel(room *Room) {
	if room.ActiveTime.Add(24 * time.Hour).Before(time.Now()) {
		log.Println("room %d is timeout 24 hours, removed.\n", room.ID)
		deleteRoom(room)
		return
	}
	living := false
	playerIds := getRoomPlayers(room.ID)
	for id := range playerIds {
		if getPlayer(id).online {
			living = true
			break
		}
	}
	if !living {
		log.Println("room %d is not living, removed.\n", room.ID)
		deleteRoom(room)
	}
}

func (p *Player) Listening() error {
	for {
		pack, err := p.conn.Read()
		if err != nil {
			log.Println(err)
			return err
		}
		if p.read {
			p.data <- pack
		}
	}
}
