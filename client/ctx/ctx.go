package ctx

import (
	"client/const"
	"client/model"
	"client/network"
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

	conn *network.Conn
}

type netConnector func(addr string) (*network.Conn, error)

var netConnectors = map[string]netConnector{
	"tcp": tcpConnect,
	/*"ws":  websocketConnect,*/
}

func New(user model.LoginRespData) *Context {
	return &Context{
		id:   user.ID,
		name: user.Name,
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
	return c.conn.Write(protocol.ObjectPacket(model.AuthInfo{
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
	return c.conn.Accept(func(packet protocol.Packet, conn *network.Conn) {
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

// 传入一个地址，返回一个连接
func tcpConnect(addr string) (*network.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("tcp server error: %v", err))
	}
	return network.Wrapper(protocol.NewTcpReadWriteCloser(conn)), nil
}

/*func websocketConnect(addr string) (*network.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ws server error: %v", err))
	}
	return network.Wrapper(protocol.NewWebsocketReadWriteCloser(conn)), nil
}
*/
