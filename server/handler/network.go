package handler

import (
	"log"
	"server/consts"
	"server/database"
	"server/model"
	"server/network"
	"server/protocol"
	"server/state"
	"server/util"
	"time"
)

// Network is interface of all kinds of network.
type Network interface {
	Serve() error
}

func handle(rwc protocol.ReadWriteCloser) error {
	// 给新进入的用户分配资源，一个id对应一个conn
	c := network.Wrapper(rwc)
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("new player connected! ")
	authInfo, err := loginAuth(c)
	if err != nil || authInfo.ID == 0 {
		_ = c.Write(protocol.ErrorPacket(err))
		return err
	}
	player := database.Connected(c, authInfo)
	log.Printf("player auth accessed, ip %s, %d:%s\n\n", player.IP, authInfo.ID, authInfo.Name)
	go state.Run(player)
	defer player.Offline()
	return player.Listening()
}

// 登陆验签
func loginAuth(c *network.Conn) (*model.AuthInfo, error) {
	authChan := make(chan *model.AuthInfo)
	defer close(authChan)
	util.Async(func() {
		packet, err := c.Read()
		if err != nil {
			log.Println(err)
			return
		}
		authInfo := &model.AuthInfo{}
		err = packet.Unmarshal(authInfo)
		if err != nil {
			log.Println(err)
			return
		}
		authChan <- authInfo
	})
	select {
	case authInfo := <-authChan:
		return authInfo, nil
	// 最多等待3s
	case <-time.After(3 * time.Second):
		return nil, consts.ErrorsAuthFail
	}
}
