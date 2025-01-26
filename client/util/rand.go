package util

import (
	"fmt"
	"math/rand"
	"time"
)

var lastName = []string{
	"Zhao", "Qian", "Sun", "Li", "Zhou", "Wu", "Zheng", "Wang",
}
var firstName = []string{
	"Tom", "Mike", "Jack", "Mary", "Jim", "Tim", "Tomas", "Json", "Black",
}

func RandomName() string {
	rand.Seed(time.Now().UnixNano()) //设置随机数种子
	return fmt.Sprintf("%s%s", fmt.Sprint(lastName[rand.Intn(len(lastName)-1)]), fmt.Sprint(firstName[rand.Intn(len(firstName)-1)]))
}
