package shell

import (
	"client/ctx"
	"client/model"
	"client/util"
	"log"
	"strings"
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
	if strings.HasSuffix(s.addr, "9998") {
		net = "ws"
	}
	err := s.ctx.Connect(net, s.addr)
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.ctx.Auth()
	if err != nil {
		log.Println(err)
		return err
	}
	return s.ctx.Listener()
}
