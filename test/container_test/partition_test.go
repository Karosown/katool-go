package container

import (
	"testing"

	"github.com/karosown/katool/container/lists"
	"github.com/karosown/katool/container/stream"
)

func Test_Partition(t *testing.T) {
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
		}, {
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
		}, {
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
		}, {
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
		}, {
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
		}, {
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
		}, {
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
		}, {
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
	sum := lists.Partition(userList, 3).Stream().Reduce(0, func(cntValue any, nxt []User) any {
		return stream.ToStream(&nxt).Reduce(cntValue, func(sumValue any, nxt User) any {
			return sumValue.(int) + nxt.Money
		})
	})
	println(sum.(int))
}
