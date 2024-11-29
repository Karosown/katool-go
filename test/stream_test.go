package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"go.uber.org/zap"
	"katool/convert"
	remote "katool/net/http"
	"katool/stream"
)

func TestOfStream(t *testing.T) {

	arr := []int{1, 2, 3, 3, 3, 3}

	distinct := stream.ToStream(&arr).Filter(func(i int) bool {
		return i > 1
	}).Map(func(item int) any {
		return strconv.Itoa(item) + "w "
	}).Distinct(func(cnt, nxt any) bool {
		return cnt.(string) == nxt.(string)
	})
	fmt.Println(distinct.Reduce("", func(cntValue any, nxt any) any {
		return cntValue.(string) + nxt.(string)
	}))
	list := distinct.ToOptionList()
	list.ForEach(func(s any) {
		fmt.Println(s)
	})

	toMap := stream.ToStream(&arr).ToMap(func(index int, item int) any {
		return index
	}, func(index int, item int) any {
		return item
	})

	maputil.ForEach(toMap, func(key any, value any) {
		fmt.Printf("key: %v, value: %v\n", key, value)
	})
	req := remote.OAuth2Req{}
	logger := zap.SugaredLogger{}
	req.SetLogger(&logger)
}

func Test_Map(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Sex   int    `json:"sex"`
		Money int    `json:"money"`
		Id    int    `json:"id"`
	}
	type UserVo struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Sex  int    `json:"sex"`
		Id   int    `json:"id"`
	}
	userList := []User{
		{
			Name:  "张三",
			Age:   18,
			Sex:   0,
			Money: 23456789,
		},
		{
			Name:  "李四",
			Age:   28,
			Sex:   1,
			Money: 23456789,
			Id:    1,
		},
		{
			Name:  "王五",
			Age:   38,
			Sex:   0,
			Money: 23456789,
			Id:    2,
		},
	}
	userStream := stream.ToStream(&userList)
	println(userStream.Count())
	totalMoney := userStream.Reduce(int64(0), func(cntValue any, nxt User) any {
		return cntValue.(int64) + int64(nxt.Money)
	})
	println(totalMoney.(int64))
	s := userStream.Map(func(item User) any {
		properties, err := convert.CopyProperties(&item, &UserVo{})
		if err != nil {
			panic(err)
		}
		return properties
	}).ToOptionList()
	s.ForEach(func(s any) {
		fmt.Println(s)
	})
}
