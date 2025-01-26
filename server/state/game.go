package state

import (
	"bytes"
	"fmt"
	"math/rand"
	"server/consts"
	"server/database"
	"server/model"
	"server/util"
	"strings"
	"time"
)

type Game struct{}

var (
	stateRob       = 1
	statePlay      = 2
	stateReset     = 3
	stateWaiting   = 4
	stateFirstCard = 5
	stateTakeCard  = 6
)

func InitGame(room *database.Room) (*database.Game, error) {
	// 得到牌的分发
	distributes := util.Distribute(room.Players)
	players := make([]int64, 0)
	roomPlayers := database.RoomPlayers(room.ID)
	for playerId := range roomPlayers {
		players = append(players, playerId)
	}

	states := map[int64]chan int{}
	groups := map[int64]int{}
	pokers := map[int64]model.Pokers{}
	playTimes := map[int64]int{}
	playTimeout := map[int64]time.Duration{}
	for i := range players {
		states[players[i]] = make(chan int, 1) // 打牌的人的状态
		groups[players[i]] = 0
		pokers[players[i]] = distributes[i] // 打牌的人持有的牌
		playTimes[players[i]] = 1
		playTimeout[players[i]] = consts.PlayTimeout
	}
	states[players[rand.Intn(len(states))]] <- stateRob
	return &database.Game{
		Room:        room,
		States:      states,
		Players:     players,
		Groups:      groups,
		Pokers:      pokers,
		Additional:  distributes[len(distributes)-1], // 附加的牌
		Multiple:    1,
		PlayTimes:   playTimes,
		PlayTimeOut: playTimeout,
		Discards:    model.Pokers{},
	}, nil
}

func (g *Game) Next(player *database.Player) (consts.StateID, error) {
	room := database.GetRoom(player.RoomID)
	if room == nil {
		return 0, player.WriteError(consts.ErrorsExist)
	}
	game := room.Game.(*database.Game)
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
			/*for i, id := range game.Players {
				game.Groups[id] = i
			}
			game.States[player.ID] <- statePlay*/
			handleRob(player, game)
		case statePlay:
			err := handlePlay(player, game)
			if err != nil {
				return 0, err
			}
		case stateWaiting:
			return consts.StateWaiting, nil
		default:
			return 0, consts.ErrorsChanClosed
		}
	}
}

func (*Game) Exit(player *database.Player) consts.StateID {
	return consts.StateHome
}

func handlePlay(player *database.Player, game *database.Game) error {
	master := player.ID == game.LastPlayer || game.LastPlayer == 0
	database.Broadcast(player.RoomID, fmt.Sprintf("%s turn to play\n", player.Name))
	return playing(player, game, master, game.PlayTimes[player.ID])
}

func playing(player *database.Player, game *database.Game, master bool, playTimes int) error {
	timeout := game.PlayTimeOut[player.ID]
	for {
		buf := bytes.Buffer{}
		buf.WriteString("\n")
		if !master && len(game.LastPokers) > 0 {
			buf.WriteString(fmt.Sprintf("Last player: %s (%s), played: %s\n", database.GetPlayer(game.LastPlayer).Name, game.Team(game.LastPlayer), game.LastPokers.String()))
		}
		buf.WriteString(fmt.Sprintf("Timeout: %ds, pokers: %s\n", int(timeout.Seconds()), game.Pokers[player.ID].String()))
		_ = player.WriteString(buf.String())
		before := time.Now().Unix()
		pokers := game.Pokers[player.ID]
		ans, err := player.AskForString(timeout)
		if err != nil {
			if master {
				ans = util.GetAlias(pokers[0].Key)
			} else {
				ans = "p"
			}
		} else {
			timeout -= time.Second * time.Duration(time.Now().Unix()-before)
		}
		ans = strings.ToLower(ans)
		if ans == "" {
			_ = player.WriteString(fmt.Sprintf("%s\n", consts.ErrorsPokersFacesInvalid.Error()))
			continue
		} else if ans == "p" || ans == "pass" {
			if master {
				_ = player.WriteError(consts.ErrorsHaveToPlay)
				continue
			} else {
				nextPlayer := database.GetPlayer(game.NextPlayer(player.ID))
				database.Broadcast(player.RoomID, fmt.Sprintf("%s passed, next %s\n", player.Name, nextPlayer.Name))
				game.States[nextPlayer.ID] <- statePlay
				return nil
			}
		}
		normalPokers := map[int]model.Pokers{}
		for _, v := range pokers {
			normalPokers[v.Key] = append(normalPokers[v.Key], v)
		}
		sells := make(model.Pokers, 0) // 要出的
		for _, alias := range ans {
			key := util.GetKey(string(alias))
			sells = append(sells, normalPokers[key][len(normalPokers[key])-1]) // 出的牌
			normalPokers[key] = normalPokers[key][:len(normalPokers[key])-1]   // 剩的牌
		}
		pokers = make(model.Pokers, 0)
		for _, curr := range normalPokers {
			pokers = append(pokers, curr...)
		}
		game.Pokers[player.ID] = pokers
		game.LastPlayer = player.ID
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
		if master {
			playTimes--
			if playTimes > 0 {
				database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s\n", player.Name, sells.String()))
				return playing(player, game, master, playTimes)
			}
		}

		game.Pokers[player.ID].SortByValue()
		err = database.GetPlayer(player.ID).WriteString(fmt.Sprintf(">> poker left:%s\n", game.Pokers[player.ID].String()))
		if err != nil {
			return err
		}
		nextPlayer := database.GetPlayer(game.NextPlayer(player.ID))
		database.Broadcast(player.RoomID, fmt.Sprintf("%s played %s, next %s\n", player.Name, sells.String(), nextPlayer.Name))
		game.States[nextPlayer.ID] <- statePlay
		return nil
	}
}

func handleRob(player *database.Player, game *database.Game) {
	// 随机取一个，作为地主
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
}
