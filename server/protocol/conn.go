package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"server/consts"
	"strconv"
	"sync/atomic"
)

var connId int64

type Conn struct {
	id    int64
	state int
	conn  ReadWriteCloser
}

func Wrapper(conn ReadWriteCloser) *Conn {
	return &Conn{
		id:   atomic.AddInt64(&connId, 1),
		conn: conn,
	}
}

func (c *Conn) ID() int64 {
	return c.id
}

func (c *Conn) IP() string {
	return c.conn.IP()
}

func (c *Conn) Close() error {
	c.state = 1
	return c.conn.Close()
}

func (c *Conn) State() int {
	return c.state
}

func (c *Conn) Write(packet Packet) error {
	return c.conn.Write(packet)
}

func (c *Conn) Read() (*Packet, error) {
	return c.conn.Read()
}

type ReadWriteCloser interface {
	Read() (*Packet, error)
	Write(msg Packet) error
	Close() error
	IP() string
}

func readUint32(reader io.Reader) (uint32, error) {
	data := make([]byte, 4)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func encode(msg Packet) []byte {
	var lenSize = 4
	lenBytes := make([]byte, lenSize)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msg.Body)))
	data := make([]byte, 0)
	data = append(data, lenBytes...)
	return append(data, msg.Body...)
}

func decode(r io.Reader) (*Packet, error) {
	l, err := readUint32(r)
	if err != nil {
		return nil, err
	}
	if l > consts.MaxPacketSize {
		return nil, errors.New("Overflow max packet size " + strconv.Itoa(consts.MaxPacketSize))
	}
	dataBytes := make([]byte, l)
	_, err = io.ReadFull(r, dataBytes)
	if err != nil {
		return nil, err
	}
	return &Packet{
		Body: dataBytes,
	}, nil
}
