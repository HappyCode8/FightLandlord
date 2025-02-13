package model

import (
	"bytes"
	"server/consts"
	"sort"
	"strconv"
)

type Faces struct {
	Keys   []int            `json:"keys"`
	Values []int            `json:"values"`
	Score  int64            `json:"score"`
	Type   consts.FacesType `json:"type"`
	Main   int              `json:"main"`
	Extra  int              `json:"extra"`
}

func (f *Faces) SetValues(values []int) *Faces {
	f.Values = values
	return f
}

func (f *Faces) SetKeys(keys []int) *Faces {
	f.Keys = keys
	return f
}

func (f *Faces) SetMain(main int) *Faces {
	f.Main = main
	return f
}

func (f *Faces) SetExtra(extra int) *Faces {
	f.Extra = extra
	return f
}

func (f *Faces) SetScore(score int64) *Faces {
	f.Score = score
	return f
}

func (f *Faces) SetType(t consts.FacesType) *Faces {
	f.Type = t
	return f
}

func (f *Faces) String() string {
	buf := bytes.Buffer{}
	for _, k := range f.Keys {
		buf.WriteString(strconv.Itoa(k))
	}
	return buf.String()
}

func (f *Faces) Compare(lastFaces Faces) bool {
	if f.Type == consts.Bomb {
		return f.Score > lastFaces.Score
	}
	if f.Type != lastFaces.Type {
		return false
	}
	return f.Score > lastFaces.Score && f.Main == lastFaces.Main && f.Extra == lastFaces.Extra
}

func ParseFaces(pokers Pokers) *Faces {
	invalidFaces := &Faces{
		Type: consts.Invalid,
	}
	if len(pokers) == 0 { //33344
		return invalidFaces
	}
	sCount, xCount, score := 0, 0, int64(0)
	valueCountMap := map[int]int{}        // 记录牌与张数的关系，牌用的是value <3,3  4,2>,3张3，4张2
	valueCountGroupMap := map[int][]int{} // 记录几张的有哪些，<3,[3]  2,[4]>,3张的有3，2张的有4
	countNums := make([]int, 0)           // 记录有几种张数 3 2, 有3张的，有2张的
	values := make([]int, 0)              // 1 1 1 2 2, 原先的值-2
	for _, poker := range pokers {
		if poker.Key < 0 || poker.Key > 15 {
			return invalidFaces
		}
		poker.Val = GetValueByKey(poker.Key)
		score += int64(poker.Val)
		values = append(values, poker.Val)
		valueCountMap[poker.Val]++
		if poker.Key == 14 {
			sCount++
		} else if poker.Key == 15 {
			xCount++
		}
	}
	for v, c := range valueCountMap {
		valueCountGroupMap[c] = append(valueCountGroupMap[c], v)
	}
	for c := range valueCountGroupMap {
		countNums = append(countNums, c)
		sort.Ints(valueCountGroupMap[c]) // 对每种张数的排序
	}
	sort.Ints(countNums)
	for i := 0; i < len(countNums)/2; i++ {
		countNums[i], countNums[len(countNums)-i-1] = countNums[len(countNums)-i-1], countNums[i]
	}
	// 王炸
	if sCount+xCount == len(pokers) && sCount+xCount == 2 {
		return &Faces{
			Values: values,
			Score:  int64(sCount*14+xCount*15)*2 + int64(len(pokers)*2*1000), // 炸弹计4张，带的计4张
			Type:   consts.KingBomb,
		}
	}
	// 最多的是单牌
	if countNums[0] == 1 {
		if len(valueCountGroupMap[countNums[0]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.Single,
			}
		}
		// 顺子
		if isStraight(valueCountGroupMap[countNums[0]], countNums[0]) {
			return &Faces{
				Values: values,
				Score:  score,
				Main:   len(valueCountGroupMap[countNums[0]]),
				Type:   consts.SingleStraight,
			}
		}
	}
	// 最多的是两张
	if countNums[0] == 2 {
		// 对
		if len(countNums) == 1 && len(valueCountGroupMap[countNums[0]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.Double,
			}
		}
		// 连对
		if isStraight(valueCountGroupMap[countNums[0]], countNums[0]) {
			return &Faces{
				Values: values,
				Score:  score,
				Main:   len(valueCountGroupMap[countNums[0]]),
				Type:   consts.DoubleStraight,
			}
		}
	}
	// 最多的是3张
	if countNums[0] == 3 {
		// 三张
		if len(countNums) == 1 && len(valueCountGroupMap[countNums[0]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.Triple,
			}
		}
		// 三带一
		if len(countNums) == 2 && len(valueCountGroupMap[countNums[0]]) == 1 && countNums[1] == 1 && len(valueCountGroupMap[countNums[1]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.TripleWithSingle,
			}
		}
		// 三带二
		if len(countNums) == 2 && len(valueCountGroupMap[countNums[0]]) == 1 && countNums[1] == 2 && len(valueCountGroupMap[countNums[1]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.TripleWithDouble,
			}
		}
		// 飞机
		if len(countNums) == 1 && isStraight(valueCountGroupMap[countNums[0]], countNums[0]) {
			if countNums[0] == 3 {
				// 只有一种，而且这一种大于3张而且是顺子
				return &Faces{
					Values: values,
					Score:  score,
					Main:   len(valueCountGroupMap[countNums[0]]),
					Type:   consts.TripleStraight,
				}
			}
		}
		// 飞机带单
		if len(countNums) == 2 && isStraight(valueCountGroupMap[countNums[0]], countNums[0]) && countNums[1] == 1 && len(valueCountGroupMap[countNums[1]]) == len(valueCountGroupMap[countNums[0]]) {
			if countNums[0] == 3 {
				// 只有一种，而且这一种大于3张而且是顺子
				return &Faces{
					Values: values,
					Score:  score,
					Main:   len(valueCountGroupMap[countNums[0]]),
					Type:   consts.TripleStraightSingle,
				}
			}
		}
		// 飞机带对
		if len(countNums) == 2 && isStraight(valueCountGroupMap[countNums[0]], countNums[0]) && countNums[1] == 2 && len(valueCountGroupMap[countNums[1]]) == len(valueCountGroupMap[countNums[0]]) {
			if countNums[0] == 3 {
				// 只有一种，而且这一种大于3张而且是顺子
				return &Faces{
					Values: values,
					Score:  score,
					Main:   len(valueCountGroupMap[countNums[0]]),
					Type:   consts.TripleStraightDouble,
				}
			}
		}
	}

	// 最多的是4张
	if countNums[0] == 4 {
		// 炸弹
		if len(countNums) == 1 && len(valueCountGroupMap[countNums[0]]) == 1 {
			return &Faces{
				Values: values,
				Score:  score,
				Type:   consts.Bomb,
			}
		}
		// 四带两张
		if len(countNums) == 2 && len(valueCountGroupMap[countNums[0]]) == 1 && (countNums[1] == 1 && len(valueCountGroupMap[countNums[1]]) == 2) {
			return &Faces{
				Values: values,
				Score:  int64(valueCountGroupMap[countNums[0]][0] * countNums[0]),
				Main:   len(valueCountGroupMap[countNums[0]]),
				Extra:  countNums[1],
				Type:   consts.QuarterWithTwoSingle,
			}
		}
		// 四带两对
		if len(countNums) == 2 && len(valueCountGroupMap[countNums[0]]) == 1 && (countNums[1] == 2 && len(valueCountGroupMap[countNums[1]]) == 2) {
			return &Faces{
				Values: values,
				Score:  int64(valueCountGroupMap[countNums[0]][0] * countNums[0]),
				Main:   len(valueCountGroupMap[countNums[0]]),
				Extra:  countNums[1],
				Type:   consts.QuarterWithTwoDouble,
			}
		}
	}

	return invalidFaces
}

// isStraight 是不是顺子(包括单、双、三),faces存的是val值
func isStraight(faces []int, count int) bool {
	if faces[len(faces)-1]-faces[0] != len(faces)-1 {
		return false
	}
	// 不能带2,X,S
	if faces[len(faces)-1] > 12 {
		return false
	}
	if count == 1 {
		// 单张的必须5连
		return len(faces) >= 5
	}
	if count == 2 {
		// 对的必须3连
		return len(faces) >= 3
	}
	if count > 2 {
		// 三张的必须2连
		return len(faces) >= 2
	}
	return false
}
