package model

import (
	"strconv"
	"sync"
)

var (
	base      = make(Pokers, 0)    // 存储所有的牌，1-13存花色，14，15存大小王
	desc      = map[int]string{}   // 存每一张牌对应的名字，1是A，11,12,13是JQK，14,15是SX，其余是数字
	keysAlias = map[int][]string{} // 每张牌对应的别名
	aliasKeys = map[string]int{}   // 每张别名对应的牌，出牌时使用这个结构找到真正的牌
	once      sync.Once
)

func InitPackPoker() {
	once.Do(createPackNewPokers)
}

func createPackNewPokers() {
	for k := 1; k <= 15; k++ {
		switch k {
		case 1:
			desc[k] = "A"
			keysAlias[k] = []string{"1", "a", "A"}
		case 10:
			desc[k] = "10"
			keysAlias[k] = []string{"0", "t", "T"}
		case 11:
			desc[k] = "J"
			keysAlias[k] = []string{"j", "J"}
		case 12:
			desc[k] = "Q"
			keysAlias[k] = []string{"q", "Q"}
		case 13:
			desc[k] = "K"
			keysAlias[k] = []string{"k", "K"}
		case 14:
			desc[k] = "S"
			keysAlias[k] = []string{"s", "S"}
		case 15:
			desc[k] = "X"
			keysAlias[k] = []string{"x", "X"}
		default:
			desc[k] = strconv.Itoa(k)
			keysAlias[k] = []string{strconv.Itoa(k)}
		}
	}
	//每一个别名都能对应到相应的数字上
	for k, aliases := range keysAlias {
		for _, alias := range aliases {
			aliasKeys[alias] = k
		}
	}
	// 加上花色
	for k := 1; k <= 13; k++ {
		for t := 0; t < 4; t++ {
			base = append(base, Poker{
				Key:  k,
				Desc: desc[k],
				Suit: PokerSuit(t),
			})
		}
	}
	// 加上大小王
	for k := 14; k <= 15; k++ {
		base = append(base, Poker{
			Key:  k,
			Desc: desc[k],
		})
	}
}

// Distribute 洗牌，洗成[[17张],[17张],[17张],[3张]]的结果
func Distribute(number int) []Pokers {
	pokers := make(Pokers, 0)
	//没洗过的牌，没有权重
	pokers = append(pokers, base...)
	// 给牌值加上权重
	for i := range pokers {
		pokers[i].Val = GetValueByKey(pokers[i].Key)
	}
	size := len(pokers)
	// 洗牌
	pokers.Shuffle(size)
	avgNum := 17
	pokersArr := make([]Pokers, 0)
	for i := 0; i < number; i++ {
		// 每个人的牌数
		pokerArr := make([]Poker, 0)
		pokersArr = append(pokersArr, append(pokerArr, pokers[i*avgNum:(i+1)*avgNum]...))
	}

	pokerArr := make([]Poker, 0)
	pokersArr = append(pokersArr, append(pokerArr, pokers[size-3:]...))

	for i := range pokersArr {
		pokersArr[i].SortByValue()
	}
	return pokersArr
}

// GetKey 根据出的牌找到真正的牌的key
func GetKey(alias string) int {
	return aliasKeys[alias]
}

// GetValueByKey Value对牌型大小计分，大小关系是345678910JQKA2SX，对应的计分分别是3 4 5 6 7 8 9 10 11 12 13 14 15
func GetValueByKey(key int) int {
	if key == 1 {
		return 12 // A
	} else if key == 2 {
		return 13 // 2
	} else if key > 13 {
		return key // 大小王
	}
	return key - 2 //其余牌
}
