package model

import (
	"client/const"
	"client/protocol"
	"client/util"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

const cleanLine = "\r\r                                                                                              \r\r"

type Context struct {
	sync.Mutex
	id   int64
	name string

	conn *protocol.Conn
}

type netConnector func(addr string) (*protocol.Conn, error)

var netConnectors = map[string]netConnector{
	"tcp": tcpConnect,
	"ws":  websocketConnect,
}

func NewContext(userId int64, userName string) *Context {
	return &Context{
		id:   userId,
		name: userName,
	}
}

func (c *Context) Connect(net string, addr string) error {
	if connector, ok := netConnectors[net]; ok {
		conn, err := connector(addr)
		if err != nil {
			return err
		}
		c.conn = conn
		return nil
	}
	return errors.New(fmt.Sprintf("unsupported net type: %s", net))
}

func (c *Context) Auth() error {
	return c.conn.Write(protocol.ObjectPacket(AuthInfo{
		ID:   c.id,
		Name: c.name,
	}))
}

func (c *Context) Listener() error {
	inputIsValid := false // 用来控制是向下游求包还是展示信息
	// 这个异步任务的作用是接收输入然后往conn里面写，传给服务器
	util.Async(func() {
		for {
			// 等着控制台输入
			line, err := util.Readline()
			if err != nil {
				log.Panic(err)
			}
			if !inputIsValid {
				continue
			}
			c.print(fmt.Sprintf(cleanLine+"[%s@ratel %s]# ", strings.TrimSpace(strings.ToLower(c.name)), "~"))
			err = c.conn.Write(protocol.Packet{
				Body: line,
			})
			if err != nil {
				continue
			}
		}
	})
	// 这个的作用是接收服务器的信息，
	return c.conn.Accept(func(packet protocol.Packet, conn *protocol.Conn) {
		data := string(packet.Body)
		if data == consts.IsStart {
			if !inputIsValid {
				c.print(fmt.Sprintf(cleanLine+"[%s@ratel %s]# ", strings.TrimSpace(strings.ToLower(c.name)), "~"))
			}
			inputIsValid = true
			return
		} else if data == consts.IsStop {
			if inputIsValid {
				c.print(cleanLine)
			}
			inputIsValid = false
			return
		}
		// 求包的时候展示输入信息，否则展示服务器信息
		if inputIsValid {
			c.print(cleanLine + data + fmt.Sprintf(cleanLine+"[%s@ratel %s]# ", strings.TrimSpace(strings.ToLower(c.name)), "~"))
		} else {
			c.print(data)
		}
	})
}

func (c *Context) print(str string) {
	c.Lock()
	defer c.Unlock()
	fmt.Print(str)
}

func tcpConnect(addr string) (*protocol.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	return protocol.Wrapper(protocol.NewTcpReadWriteCloser(conn)), nil
}

func websocketConnect(addr string) (*protocol.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ws server error: %v", err))
	}
	return protocol.Wrapper(protocol.NewWebsocketReadWriteCloser(conn)), nil
}
