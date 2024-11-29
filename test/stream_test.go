package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"katool/algorithm"
	"katool/convert"
	"katool/stream"
)

func TestOfStream(t *testing.T) {

	arr := []int{1, 3, 2, 3, 3, 3, 3}

	distinct := stream.ToStream(&arr).Filter(func(i int) bool {
		return i > 1
	}).Map(func(item int) any {
		return strconv.Itoa(item) + "w "
	}).Distinct(algorithm.HASH_WITH_JSON)

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
	// 计数
	userStream := stream.ToStream(&userList)
	println(userStream.Count())
	// 排序
	stream.ToStream(&userList).Sort(func(a User, b User) bool { return a.Id > b.Id }).ToOptionList().ForEach(func(item User) {
		println(convert.ConvertToString(item.Id) + "" + item.Name)
	})
	// 求和
	totalMoney := userStream.Reduce(int64(0), func(cntValue any, nxt User) any {
		return cntValue.(int64) + int64(nxt.Money)
	})
	println(totalMoney.(int64))
	// 过滤
	userStream.Filter(func(item User) bool { return item.Sex != 0 }).ToOptionList().ForEach(func(item User) {
		println(item.Name)
	})
	// 转换
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
