package main

import (
	"flag"
	"log"
	"server/handler"
	"server/model"
	"strconv"
)

var (
	TCPPort int
)

func init() {
	model.InitPackPoker()
}

func main() {
	flag.IntVar(&TCPPort, "t", 9999, "TcpServer Port")
	flag.Parse()

	server := handler.NewTcpServer(":" + strconv.Itoa(TCPPort))
	log.Panic(server.Serve())
}
