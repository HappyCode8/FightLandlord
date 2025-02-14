package handler

import (
	"log"
	"server/database"
	"server/errdef"
	"server/model"
	"server/protocol"
	"server/service"
	"server/util"
	"time"
)

// func handle(rwc protocol.ReadWriteCloser) error {
func handle(conn protocol.ReadWriteCloser) error {
	// 给新进入的用户分配资源，一个id对应一个conn
	c := protocol.Wrapper(conn)
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
	// 开一个线程处理
	go service.Run(player)
	defer player.Offline()
	// 开始监听连接信息,这个连接把包写入player的data里，比如在home里要取一个选择创建房间还是加入房间的askpacket
	return player.Listening()
}

// 登陆验签
func loginAuth(c *protocol.Conn) (*model.AuthInfo, error) {
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
		return nil, errdef.ErrorsAuthFail
	}
}
