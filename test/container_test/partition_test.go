package container_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	lynx "github.com/Tangerg/lynx/pkg/sync"
	"github.com/karosown/katool/collect/lists"
	"github.com/karosown/katool/container/stream"
	"github.com/karosown/katool/convert"
)

func Test_Partition(t *testing.T) {
	sum := lists.Partition(userList[:], 15).Stream().Map(func(i []user) any {
		return stream.ToStream(&i).Map(func(user user) any {
			properties, _ := convert.CopyProperties(user, &userVo{})
			return *properties
		}).ToList()
	}).Reduce("", func(cntValue any, nxt any) any {
		anies := nxt.([]any)
		return stream.ToStream(&anies).Reduce(cntValue, func(sumValue any, nxt any) any {
			return sumValue.(string) + nxt.(userVo).Name
		})
	})
	println(sum.(string))
}

func Test_ForEach(t *testing.T) {
	type user struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Sex   int    `json:"sex"`
		Money int    `json:"money"`
		Id    int    `json:"id"`
	}
	type userVo struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Sex  int    `json:"sex"`
		Id   int    `json:"id"`
	}
	userList := [...]user{
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
	i := atomic.Uint32{}

	lists.Partition(userList[:], 15).ForEach(func(automicDatas []user) error {
		i.Add(1)
		fmt.Println("分批处理" + convert.ConvertToString(i.Load()))
		stream.ToStream(&automicDatas).ForEach(func(data user) {
			fmt.Println(data)
		})
		return nil
	}, true, lynx.NewLimiter(3))
	//println(sum.(string))
}
