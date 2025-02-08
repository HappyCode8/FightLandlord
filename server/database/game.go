package database

import (
	"server/model"
	"time"
)

type Game struct {
	Room        *Room                   `json:"room"`        // 持有的房间
	Players     []int64                 `json:"players"`     // 持有的人
	States      map[int64]chan int      `json:"states"`      // 每个人的状态
	Pokers      map[int64]model.Pokers  `json:"pokers"`      // 每个人的牌
	Additional  model.Pokers            `json:"pocket"`      // 附加的牌
	FirstPlayer int64                   `json:"firstPlayer"` // 第一个出牌的人
	LastPlayer  int64                   `json:"lastPlayer"`  // 最后一个出牌的人
	PlayTimeOut map[int64]time.Duration `json:"playTimeOut"` // 超时时间
}

func (game *Game) Clean() {
	if game != nil {
		for _, state := range game.States {
			close(state)
		}
	}
}

// NextPlayer 获取下一个出牌人，获取的方式是当前人的(id+1)%总人数
func (game *Game) NextPlayer(curr int64) int64 {
	var idx int
	for index, value := range game.Players {
		if value == curr {
			idx = index
		}
	}
	return game.Players[(idx+1)%len(game.Players)]
}
