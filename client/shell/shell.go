package shell

import (
	"client/ctx"
	"client/model"
	"client/util"
	"log"
	"time"
)

type shell struct {
	ctx  *ctx.Context
	addr string
	name string
}

func New(addr, name string) *shell {
	return &shell{
		addr: addr,
		name: name,
	}
}

func (s *shell) Start() error {
	name := util.RandomName()
	s.ctx = ctx.New(model.LoginRespData{
		ID:   time.Now().UnixNano(),
		Name: name,
	})
	net := "tcp"
	/*if strings.HasSuffix(s.addr, "9998") {
		net = "ws"
	}*/
	// 建立一个连接
	err := s.ctx.Connect(net, s.addr)
	if err != nil {
		log.Println(err)
		return err
	}
	// 往连接里边写一个包
	err = s.ctx.Auth()
	if err != nil {
		log.Println(err)
		return err
	}
	// 然后开始监听连接信息
	return s.ctx.Listener()
}
