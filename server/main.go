package main

import (
	"flag"
	"log"
	"server/handler"
	"strconv"
)

var (
	TCPPort int
)

func main() {
	flag.IntVar(&TCPPort, "t", 9999, "TcpServer Port")
	flag.Parse()

	server := handler.NewTcpServer(":" + strconv.Itoa(TCPPort))
	log.Panic(server.Serve())
}
