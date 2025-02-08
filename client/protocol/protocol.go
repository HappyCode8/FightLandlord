package protocol

import (
	"encoding/json"
	"strconv"
)

var (
	lenSize = 4
)

type Packet struct {
	Body []byte `json:"data"`
}

func (p Packet) Int() (int, error) {
	v, err := strconv.ParseInt(p.String(), 10, 64)
	return int(v), err
}

func (p Packet) Int64() (int64, error) {
	v, _ := strconv.ParseInt(p.String(), 10, 64)
	return v, nil
}

func (p Packet) String() string {
	return string(p.Body)
}

func (p Packet) Unmarshal(v interface{}) error {
	err := json.Unmarshal(p.Body, v)
	if err != nil {
		return err
	}
	return nil
}

func ObjectPacket(obj interface{}) Packet {
	marshal, _ := json.Marshal(obj)
	return Packet{
		Body: marshal,
	}
}
