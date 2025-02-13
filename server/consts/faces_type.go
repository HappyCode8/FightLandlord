package consts

type FacesType int

const (
	_ FacesType = iota

	KingBomb             //王炸
	Bomb                 //炸弹
	Single               //单牌
	Double               //对子
	Triple               //三张
	TripleWithSingle     //三带一
	TripleWithDouble     //三带二
	QuarterWithTwoSingle //四带二
	QuarterWithTwoDouble //四带两对
	SingleStraight       //单顺子
	DoubleStraight       //连对
	TripleStraight       //飞机
	TripleStraightSingle //飞机带单
	TripleStraightDouble //飞机带对
	Invalid
)
