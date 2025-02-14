package protocol

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type WebsocketReadWriteCloser struct {
	conn *websocket.Conn
}

func NewWebsocketReadWriteCloser(conn *websocket.Conn) WebsocketReadWriteCloser {
	return WebsocketReadWriteCloser{conn: conn}
}

func (w WebsocketReadWriteCloser) Read() (*Packet, error) {
	_, b, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msg := &Packet{}
	_ = json.Unmarshal(b, msg)
	return msg, nil
}

func (w WebsocketReadWriteCloser) Write(msg Packet) error {
	marshal, _ := json.Marshal(msg)
	return w.conn.WriteMessage(websocket.BinaryMessage, marshal)
}

func (w WebsocketReadWriteCloser) Close() error {
	return w.conn.Close()
}

func (w WebsocketReadWriteCloser) IP() string {
	return w.conn.RemoteAddr().String()
}
