package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"server/consts"
	"server/database"
	"server/errdef"
	"server/model"
	"strings"
	"time"
)

type Game struct{}

var (
	stateRob     = 1
	statePlay    = 2
	stateWaiting = 3
)

func InitGame(room *database.Room) (*database.Game, error) {
	// 得到牌的分发
	distributes := model.Distribute(room.Players)
	players := make([]int64, 0)
	roomPlayers := database.RoomPlayers(room.ID)
	for playerId := range roomPlayers {
		players = append(players, playerId)
	}

	states := map[int64]chan int{}
	pokers := map[int64]model.Pokers{}
	playTimeout := map[int64]time.Duration{}
	for i := range players {
		states[players[i]] = make(chan int, 1) // 打牌的人的状态
		pokers[players[i]] = distributes[i]    // 打牌的人持有的牌
		playTimeout[players[i]] = consts.PlayTimeout
	}
	states[players[rand.Intn(len(states))]] <- stateRob // 随机选一个人作为抢状态
	return &database.Game{
		Room:        room,
		States:      states,
		Players:     players,
		Pokers:      pokers,
		Additional:  distributes[len(distributes)-1], // 附加的牌
		PlayTimeOut: playTimeout,
	}, nil
}

func (g *Game) Next(player *database.Player) (consts.StateID, error) {
	room := database.GetRoom(player.RoomID)
	if room == nil {
		return 0, player.WriteError(errdef.ErrorsExist)
	}
	game := room.Game
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Game starting!\n"))
	buf.WriteString(fmt.Sprintf("Your pokers: %s\n", game.Pokers[player.ID].String()))
	_ = player.WriteString(buf.String())
	for {
		if room.State == consts.RoomStateWaiting {
			return consts.StateWaiting, nil
		}
		state := <-game.States[player.ID]
		switch state {
		case stateRob:
			// 没有完全实现，随机取一个作为地主
			handleRob(player, game)
		case statePlay:
			// 没有完全实现，只实现了把一个人出的牌全局广播
			err := handlePlay(player, game)
			if err != nil {
				return 0, err
			}
		case stateWaiting:
			return consts.StateWaiting, nil
		default:
			return 0, errdef.ErrorsChanClosed
		}
	}
}

func (*Game) Exit(player *database.Player) consts.StateID {
	return consts.StateHome
}

func handlePlay(player *database.Player, game *database.Game) error {
	database.Broadcast(player.RoomID, fmt.Sprintf("%s turn to play\n", player.Name))
	return playing(player, game)
}

func playing(player *database.Player, game *database.Game) error {
	timeout := game.PlayTimeOut[player.ID]
	for {
		buf := bytes.Buffer{}
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("Timeout: %ds, pokers: %s\n", int(timeout.Seconds()), game.Pokers[player.ID].String()))
		_ = player.WriteString(buf.String())
		before := time.Now().Unix()
		pokers := game.Pokers[player.ID]
		ans, err := player.AskForString(timeout)
		if err != nil {
			ans = "pass"
		} else {
			timeout -= time.Second * time.Duration(time.Now().Unix()-before)
		}
		ans = strings.ToLower(ans)
		if ans == "" {
			_ = player.WriteString(fmt.Sprintf("%s\n", errdef.ErrorsPokersFacesInvalid.Error()))
			continue
		} else if ans == "pass" {
			nextPlayer := database.GetPlayer(game.NextPlayer(player.ID))
			database.Broadcast(player.RoomID, fmt.Sprintf("%s passed, next %s\n", player.Name, nextPlayer.Name))
			game.States[nextPlayer.ID] <- statePlay
			return nil
		}
		normalPokers := map[int]model.Pokers{} //<3,333 5,555>样式
		for _, v := range pokers {
			normalPokers[v.Key] = append(normalPokers[v.Key], v)
		}
		sells := make(model.Pokers, 0) // 要出的
		for _, alias := range ans {
			key := model.GetKey(string(alias))
			sells = append(sells, normalPokers[key][len(normalPokers[key])-1]) // 出的牌，从normalPokers取最后一张
			normalPokers[key] = normalPokers[key][:len(normalPokers[key])-1]   // 剩的牌，从normalPokers取前几张
		}
		// 把剩下的牌重新分给用户
		pokers = make(model.Pokers, 0)
		for _, curr := range normalPokers {
			pokers = append(pokers, curr...)
		}
		game.Pokers[player.ID] = pokers
		game.LastPlayer = player.ID
		// 出完牌以后，重新初始化房间的信息
		if len(pokers) == 0 {
			database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s, won the game! \n", player.Name, sells.String()))
			room := database.GetRoom(player.RoomID)
			if room != nil {
				room.Game = nil
				room.State = consts.RoomStateWaiting
			}
			for _, playerId := range game.Players {
				game.States[playerId] <- stateWaiting
			}
			return nil
		}
		// 给用户发消息
		game.Pokers[player.ID].SortByValue()
		err = database.GetPlayer(player.ID).WriteString(fmt.Sprintf(">> poker left:%s\n", game.Pokers[player.ID].String()))
		if err != nil {
			return err
		}
		// 给全局进行广播
		nextPlayer := database.GetPlayer(game.NextPlayer(player.ID))
		database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s, next %s\n", player.Name, sells.String(), nextPlayer.Name))
		game.States[nextPlayer.ID] <- statePlay
		return nil
	}
}

// 模拟抢，随机取一个，作为地主
func handleRob(player *database.Player, game *database.Game) {
	players := game.Players
	randLordIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(3)
	landLordId := players[randLordIndex]
	game.Pokers[landLordId] = append(game.Pokers[landLordId], game.Additional...)
	game.Pokers[landLordId].SortByValue()
	game.States[landLordId] <- statePlay
	game.FirstPlayer = landLordId
	buf := bytes.Buffer{}
	landLord := database.GetPlayer(landLordId)
	buf.WriteString(fmt.Sprintf("%s became landlord, got pokers: %s\n", landLord.Name, game.Additional.String()))
	database.Broadcast(player.RoomID, buf.String())
}
