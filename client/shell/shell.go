package shell

import (
	"client/model"
	"client/util"
	"log"
	"time"
)

type shell struct {
	ctx  *model.Context
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
	s.ctx = model.NewContext(
		time.Now().UnixNano(),
		name,
	)
	// 建立一个连接
	err := s.ctx.Connect(s.addr)
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
