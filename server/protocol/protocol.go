package protocol

import (
	"encoding/json"
	"strconv"
)

type Packet struct {
	Body []byte `json:"data"`
}

func (p Packet) Int() (int, error) {
	v, err := strconv.ParseInt(p.String(), 10, 64)
	return int(v), err
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

func ErrorPacket(err error) Packet {
	return Packet{
		Body: []byte(err.Error()),
	}
}
