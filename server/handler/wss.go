package handler

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/protocol"
)

type Websocket struct {
	addr string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketServer(addr string) Websocket {
	return Websocket{addr: addr}
}

func (w Websocket) Serve() error {
	http.HandleFunc("/ws", serveWs)
	log.Printf("Websocket server listener on %s\n", w.addr)
	return http.ListenAndServe(w.addr, nil)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	// 将http请求升级为websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = handle(protocol.NewWebsocketReadWriteCloser(conn))
	if err != nil {
		log.Println(err)
	}
}
