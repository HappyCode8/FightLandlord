package model

import (
	"github.com/stretchr/testify/assert"
	"server/consts"
	"testing"
)

func TestIsStraight(t *testing.T) {
	assert.Equal(t, true, isStraight([]int{3, 4, 5, 6, 7}, 1))
	assert.Equal(t, true, isStraight([]int{1, 2, 3, 4, 5, 6}, 1))
	assert.Equal(t, true, isStraight([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 1))
	assert.Equal(t, false, isStraight([]int{1, 2, 3, 4, 7}, 1))
	assert.Equal(t, false, isStraight([]int{1, 2, 3, 4}, 1))

	assert.Equal(t, true, isStraight([]int{1, 2, 3}, 2))
	assert.Equal(t, true, isStraight([]int{1, 2, 3, 4}, 2))
	assert.Equal(t, false, isStraight([]int{1, 2, 4}, 2))
	assert.Equal(t, false, isStraight([]int{1, 2}, 2))

	assert.Equal(t, true, isStraight([]int{1, 2}, 3))
	assert.Equal(t, true, isStraight([]int{1, 2, 3}, 3))
	assert.Equal(t, false, isStraight([]int{1, 3}, 3))
	assert.Equal(t, false, isStraight([]int{1}, 3))
}

func getPokers(keys ...int) Pokers {
	pokers := make(Pokers, 0)
	for _, k := range keys {
		pokers = append(pokers, Poker{
			Key:  k,
			Desc: desc[k],
		})
	}
	return pokers
}

func TestParseFaces(t *testing.T) {
	assert.Equal(t, consts.KingBomb, ParseFaces(getPokers(14, 15)).Type)
	assert.Equal(t, consts.Bomb, ParseFaces(getPokers(3, 3, 3, 3)).Type)
	assert.Equal(t, consts.Single, ParseFaces(getPokers(3)).Type)
	assert.Equal(t, consts.Single, ParseFaces(getPokers(15)).Type)
	assert.Equal(t, consts.Double, ParseFaces(getPokers(3, 3)).Type)
	assert.Equal(t, consts.Triple, ParseFaces(getPokers(3, 3, 3)).Type)
	assert.Equal(t, consts.TripleWithSingle, ParseFaces(getPokers(3, 3, 3, 4)).Type)
	assert.Equal(t, consts.TripleWithDouble, ParseFaces(getPokers(3, 3, 3, 4, 4)).Type)
	assert.Equal(t, consts.TripleWithDouble, ParseFaces(getPokers(3, 2, 2, 2, 3)).Type)
	assert.Equal(t, consts.QuarterWithTwoSingle, ParseFaces(getPokers(3, 3, 3, 3, 4, 5)).Type)
	assert.Equal(t, consts.QuarterWithTwoDouble, ParseFaces(getPokers(3, 3, 3, 3, 5, 5, 6, 6)).Type)
	assert.Equal(t, consts.SingleStraight, ParseFaces(getPokers(3, 4, 5, 6, 7)).Type)
	assert.Equal(t, consts.SingleStraight, ParseFaces(getPokers(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1)).Type)
	assert.Equal(t, consts.DoubleStraight, ParseFaces(getPokers(3, 3, 4, 4, 5, 5)).Type)
	assert.Equal(t, consts.TripleStraight, ParseFaces(getPokers(3, 3, 3, 4, 4, 4)).Type)
	assert.Equal(t, consts.TripleStraightSingle, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 6)).Type)
	assert.Equal(t, consts.TripleStraightSingle, ParseFaces(getPokers(9, 9, 9, 8, 8, 8, 7, 7, 7, 11, 1, 15)).Type)
	assert.Equal(t, consts.TripleStraightDouble, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 5, 6, 6)).Type)

	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(2, 3, 4, 5, 6, 7)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(10, 11, 12, 13, 1, 2)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 4)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 3, 4, 4, 4, 4, 4)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 4, 5, 6)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 4, 5, 6, 7, 8)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 3)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 3, 4, 5)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(14, 15, 15)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(14, 14, 15, 15)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(14, 14, 14, 15, 15, 16)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(4, 4, 4, 4, 6, 6, 6, 6)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(4, 4, 4, 16)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 5)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 5, 5, 7, 7, 7)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 7, 7)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 4, 3, 3, 4, 4, 4, 5, 5, 5, 5)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10, 10, 10)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 3, 4, 4, 4, 4, 4)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 8)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 8, 8)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 8, 8)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4, 4, 6, 6, 6, 7, 7, 7, 8, 8, 8, 9, 9, 9, 9, 9)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 3, 4, 4, 6, 6, 6, 7, 7, 7, 8, 8, 8, 9, 9, 9, 9, 9)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 4, 4, 5)).Type)
	assert.Equal(t, consts.Invalid, ParseFaces(getPokers(3, 3, 3, 3, 4)).Type)
}
