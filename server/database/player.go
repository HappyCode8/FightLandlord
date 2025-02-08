package database

import (
	"fmt"
	"log"
	"server/consts"
	"server/network"
	"server/protocol"
	"strings"
	"time"
)

type Player struct {
	ID     int64  `json:"id"`
	IP     string `json:"ip"`
	Name   string `json:"name"`
	RoomID int64  `json:"roomId"` // 所属房间id

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
