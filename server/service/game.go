package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"server/consts"
	"server/database"
	"server/errdef"
	"server/model"
	"slices"
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
			handleRob(player, game)
		case statePlay:
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
		var (
			before            = time.Now().Unix()      // 进入时间打点，便于计算超时
			reservePokers     = game.Pokers[player.ID] // 持有的牌
			reservePokersMap  = map[int]model.Pokers{} // 手中的牌map <3,33 5,555>样式
			reservePokersKeys = make([]int, 0)         // 手中的牌keys
			sellPokers        = make(model.Pokers, 0)  // 出的牌
			inValid           bool                     // 出的牌是否合法
			remainPokers      = make(model.Pokers, 0)  // 出完以后手中剩下的牌
			sellFaces         *model.Faces             // 出的牌牌型解析
			lastFaces         *model.Faces             // 上家出的牌牌型解析
		)
		_ = player.WriteString(fmt.Sprintf("Timeout: %ds, reservePokers: %s\n", int(timeout.Seconds()), reservePokers.String()))
		// 请求用户出牌信息
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
			// 过的话轮到下一个
			nextPlayer := database.GetPlayer(game.NextPlayer(player.ID))
			database.Broadcast(player.RoomID, fmt.Sprintf("%s passed, next %s\n", player.Name, nextPlayer.Name))
			game.States[nextPlayer.ID] <- statePlay
			return nil
		}
		for _, v := range reservePokers {
			reservePokersMap[v.Key] = append(reservePokersMap[v.Key], v)
			reservePokersKeys = append(reservePokersKeys, v.Key)
		}

		for _, alias := range ans {
			key := model.GetKey(string(alias))
			if !slices.Contains(reservePokersKeys, key) {
				inValid = true
				break
			}
			sellPokers = append(sellPokers, reservePokersMap[key][len(reservePokersMap[key])-1]) // 出的牌，从normalPokers取最后一张
			reservePokersMap[key] = reservePokersMap[key][:len(reservePokersMap[key])-1]         // 剩的牌，从normalPokers取前几张
		}
		sellFaces = model.ParseFaces(sellPokers)
		lastFaces = game.LastFaces
		// 出的牌型不合法或者出的牌不在手里
		if sellFaces.Type == consts.Invalid || inValid {
			_ = player.WriteString(fmt.Sprintf("%s\n", errdef.ErrorsPokersFacesInvalid.Error()))
			continue
		}
		// 出的牌型跟上家不一样或者没有上家的大，但是上次出牌的是自己的除外
		if game.LastPlayer != player.ID && lastFaces != nil && (!sellFaces.Valid(lastFaces) || !sellFaces.MaxThan(lastFaces)) {
			_ = player.WriteString(fmt.Sprintf("%s\n", errdef.ErrorsPokersFacesInvalid.Error()))
			continue
		}

		for _, curr := range reservePokersMap {
			remainPokers = append(remainPokers, curr...)
		}
		game.Pokers[player.ID] = remainPokers
		game.LastPlayer = player.ID
		game.LastFaces = sellFaces
		// 出完牌以后，刷新房间信息
		if len(remainPokers) == 0 {
			database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s, won the game! \n", player.Name, sellPokers.String()))
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
		database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s, next %s\n", player.Name, sellPokers.String(), nextPlayer.Name))
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
