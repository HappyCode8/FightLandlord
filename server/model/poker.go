package model

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Poker struct {
	Key  int       `json:"key"`  // 牌的数字 1-13，14小王，15大王
	Val  int       `json:"val"`  // 牌的大小权重
	Desc string    `json:"desc"` // 牌的描述，A，2，3，4，5，6，7，8，9，10，J，Q，K，S，X
	Suit PokerSuit `json:"suit"` // 牌的花色
}

type Pokers []Poker

// Shuffle 洗牌算法，依次从i位置取出一张与前i张随机取出的一张做交换，第一张的可能性有54，第二张有53，总共的洗法位54!,只有54!的洗牌算法是正确的洗牌算法
func (pokers Pokers) Shuffle(n int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := n - 1; i > 0; i -= 1 {
		j := int(r.Int31n(int32(i + 1)))
		pokers.Swap(i, j)
	}
}

func (pokers Pokers) Swap(i, j int) {
	pokers[i], pokers[j] = pokers[j], pokers[i]
}

// SortByKey 按牌值排序
func (pokers Pokers) SortByKey() {
	sort.Slice(pokers, func(i, j int) bool {
		return pokers[i].Key < pokers[j].Key
	})
}

// SortByValue 按牌的权重排序
func (pokers Pokers) SortByValue() {
	sort.Slice(pokers, func(i, j int) bool {
		return pokers[i].Val < pokers[j].Val
	})
}

func (pokers Pokers) String() string {
	buf := bytes.Buffer{}
	for i := len(pokers) - 1; i >= 0; i-- {
		poker := pokers[i]
		buf.WriteString(fmt.Sprintf("%s", poker.Desc))
		if i != 0 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}
