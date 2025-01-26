package util

import (
	"server/model"
)

var (
	base      = make(model.Pokers, 0)
	desc      = map[int]string{}
	keysAlias = map[int][]string{}
	aliasKeys = map[string]int{}
)

func init() {
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
			desc[k] = String(k)
			keysAlias[k] = []string{String(k)}
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
			base = append(base, model.Poker{
				Key:  k,
				Val:  0,
				Desc: desc[k],
				Suit: model.PokerSuit(t),
			})
		}
	}
	for k := 14; k <= 15; k++ {
		base = append(base, model.Poker{
			Key:  k,
			Val:  0,
			Desc: desc[k],
		})
	}
}

func GetKey(alias string) int {
	return aliasKeys[alias]
}

func GetAlias(key int) string {
	if len(keysAlias[key]) > 0 {
		return keysAlias[key][0]
	}
	return ""
}

func GetDesc(key int) string {
	return desc[key]
}

// Distribute number is players number. n is shuffle times.
func Distribute(number int) []model.Pokers {
	pokers := make(model.Pokers, 0)
	//没洗过的牌，没有权重
	pokers = append(pokers, base...)
	// 给牌值加上权重
	for i := range pokers {
		pokers[i].Val = GetValueByKey(pokers[i].Key)
	}
	size := len(pokers)
	// 洗牌
	pokers.Shuffle(size, 1)
	avgNum := 17
	pokersArr := make([]model.Pokers, 0)
	for i := 0; i < number; i++ {
		// 每个人的牌数
		pokerArr := make([]model.Poker, 0)
		pokersArr = append(pokersArr, append(pokerArr, pokers[i*avgNum:(i+1)*avgNum]...))
	}

	pokerArr := make([]model.Poker, 0)
	pokersArr = append(pokersArr, append(pokerArr, pokers[size-3:]...))

	for i := range pokersArr {
		pokersArr[i].SortByValue()
	}
	for i := range pokersArr {
		pokersArr[i].SortByValue()
	}
	return pokersArr
}

// Value 对牌型大小计分，大小关系是 大小王21345678910JQK对应的计分分别是15 14 13 12 3 4 5 6 7 8 9 10 11 12 13
func GetValueByKey(key int) int {
	if key == 1 {
		return 12
	} else if key == 2 {
		return 13
	} else if key > 13 {
		return key
	}
	return key - 2
}
