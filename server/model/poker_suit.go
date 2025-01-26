package model

type PokerSuit int

const (
	Spade   PokerSuit = iota // 黑桃
	Heart                    // 红桃
	Club                     // 梅花
	Diamond                  // 方片
)

func (s PokerSuit) String() string {
	return map[PokerSuit]string{
		Spade:   "♠",
		Heart:   "♥",
		Club:    "♣",
		Diamond: "♦",
	}[s]
}
