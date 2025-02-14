package main

import (
	"flag"
	"log"
	"server/handler"
	"server/model"
	"server/util"
	"strconv"
)

var (
	TCPPort int
	WSPort  int
)

func init() {
	model.InitPackPoker()
}

func main() {
	flag.IntVar(&WSPort, "w", 9998, "WebsocketServer Port")
	flag.IntVar(&TCPPort, "t", 9999, "TcpServer Port")
	flag.Parse()

	// 这里必须异步，否则会阻塞
	util.Async(func() {
		wsServer := handler.NewWebsocketServer(":" + strconv.Itoa(WSPort))
		log.Panic(wsServer.Serve())
	})

	server := handler.NewTcpServer(":" + strconv.Itoa(TCPPort))
	log.Panic(server.Serve())
}
