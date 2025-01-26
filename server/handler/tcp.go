package handler

import (
	"log"
	"net"
	"server/protocol"
	"server/util"
)

type Tcp struct {
	addr string
}

func NewTcpServer(addr string) Tcp {
	return Tcp{addr: addr}
}

func (t Tcp) Serve() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	log.Println("tcp server listening on", t.addr)
	for {
		// 监听连接
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			log.Println("listener.Accept err", err)
			continue
		}
		// 每有一个连接，就处理
		util.Async(func() {
			handleErr := handle(protocol.NewTcpReadWriteCloser(conn))
			if handleErr != nil {
				log.Println("handle err", handleErr)
			}
		})
	}
}
