package model

import (
	"client/const"
	"client/protocol"
	"client/util"
	"errors"
	"fmt"
	"log"
	"net"
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

func NewContext(userId int64, userName string) *Context {
	return &Context{
		id:   userId,
		name: userName,
	}
}

func (c *Context) Connect(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	c.conn = protocol.Wrapper(conn)
	return nil
}

func (c *Context) Auth() error {
	return c.conn.Write(protocol.ObjectPacket(AuthInfo{
		ID:   c.id,
		Name: c.name,
	}))
}

func (c *Context) Listener() error {
	is := false
	// 这个异步任务的作用是接收输入然后往conn里面写，传给服务器
	util.Async(func() {
		for {
			// 等着控制台输入
			line, err := util.Readline()
			if err != nil {
				log.Panic(err)
			}
			if !is {
				//log.Println("debug发送:进入了!is")
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
			//log.Println("debug接收: 接收到了服务器的开始信息")
			if !is {
				c.print(fmt.Sprintf(cleanLine+"[%s@ratel %s]# ", strings.TrimSpace(strings.ToLower(c.name)), "~"))
			}
			is = true
			return
		} else if data == consts.IsStop {
			//log.Println("debug接收: 接收到了服务器的结束信息")
			if is {
				c.print(cleanLine)
			}
			is = false
			return
		}
		if is {
			//log.Println("debug接收:进入了is")
			c.print(cleanLine + data + fmt.Sprintf(cleanLine+"[%s@ratel %s]# ", strings.TrimSpace(strings.ToLower(c.name)), "~"))
		} else {
			//log.Println("debug接收:进入了is else")
			// 还没开始的时候服务器返回的信息，直接打印
			c.print(data)
		}
	})
}

func (c *Context) print(str string) {
	c.Lock()
	defer c.Unlock()
	fmt.Print(str)
}
