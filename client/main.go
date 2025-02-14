package main

import (
	"client/shell"
	"flag"
	"fmt"
	"log"
)

var (
	host string
	port int
	name string
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.IntVar(&port, "p", 9998, "port")
	flag.StringVar(&name, "n", "", "name")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Fatal(shell.New(addr, name).Start())
}
