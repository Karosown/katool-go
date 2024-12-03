package container_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool/algorithm"
	"github.com/karosown/katool/container/stream"
	"github.com/karosown/katool/convert"
	"github.com/karosown/katool/util"
)

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

var userList = [...]user{
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

func TestOfStream(t *testing.T) {

	arr := []int{1, 3, 2, 3, 3, 3, 3}

	distinct := stream.ToStream(&arr).
		Parallel().
		Filter(func(i int) bool {
			return i > 1
		}).Map(func(item int) any {
		return strconv.Itoa(item) + "w "
	}).Distinct(algorithm.HASH_WITH_JSON_SUM)

	fmt.Println(distinct.Reduce("", func(cntValue any, nxt any) any {
		return cntValue.(string) + nxt.(string)
	}, func(sum1, sum2 any) any {
		return sum1.(string) + sum2.(string)
	}))
	list := distinct.ToOptionList()
	list.ForEach(func(s any) {
		fmt.Println(s)
	})

	toMap := stream.ToStream(&arr).Parallel().ToMap(func(index int, item int) any {
		return index
	}, func(index int, item int) any {
		return item
	})

	maputil.ForEach(toMap, func(key any, value any) {
		fmt.Printf("key: %v, value: %v\n", key, value)
	})
}

func Test_Map(t *testing.T) {
	ul := userList[:]
	// 计数
	userStream := stream.ToStream(&ul)
	println(userStream.Count())
	// 排序
	stream.ToStream(&ul).
		Parallel().
		Sort(func(a user, b user) bool { return a.Id > b.Id }).ForEach(func(item user) { println(convert.ConvertToString(item.Id) + " " + item.Name) })
	// 求和
	totalMoney := userStream.Reduce(int64(0), func(cntValue any, nxt user) any { return cntValue.(int64) + int64(nxt.Money) }, func(sum1, sum2 any) any {
		return sum1.(int64) + sum2.(int64)
	})
	println(totalMoney.(int64))
	// 过滤
	userStream.Filter(func(item user) bool { return item.Sex != 0 }).Distinct(algorithm.HASH_WITH_JSON_MD5).ToOptionList().ForEach(func(item user) { println(item.Name) })
	// 转换
	s := userStream.Map(func(item user) any {
		properties, err := convert.CopyProperties(&item, &userVo{})
		if err != nil {
			panic(err)
		}
		return properties
	}).ToOptionList()
	s.ForEach(func(s any) {
		fmt.Println(s)
	})
}

func Test_GroupBy(t *testing.T) {
	users := userList[:]
	by := stream.ToStream(&users).GroupBy(func(user user) any {
		return user.Name
	})
	println(by)
}

func Test_OrderBy(t *testing.T) {
	users := userList[:]
	//for i := 0; i < 100; i++ {
	//	users = append(users, userList[:]...)
	//}
	var unParallel []user
	computed := util.BeginEndTimeComputed(func() {
		unParallel = stream.ToStream(&users).OrderBy(false, func(u any) algorithm.HashType {
			return algorithm.HashType(u.(user).Name)
		}).ToList()
	})
	println(computed)
	var parallel []user
	computed = util.BeginEndTimeComputed(func() {
		parallel = stream.ToStream(&users).Parallel().OrderBy(false, func(u any) algorithm.HashType {
			return algorithm.HashType(u.(user).Name)
		}).ToList()
	})
	println(computed)
	for i := 0; i < len(unParallel); i++ {
		if unParallel[i].Name != parallel[i].Name {
			panic("unparallel not equal parallel" + convert.ConvertToString(i))
		}
	}
}

func Test_OrderBy_Round(t *testing.T) {
	users := make([]int, 0)
	for i := 0; i < 6e6+5e5+3e3+2e2+1e1; i++ {
		users = append(users, rand.Int()%10000)
	}
	var unParallel []int
	computed := util.BeginEndTimeComputed(func() {
		unParallel = stream.ToStream(&users).OrderById(false, func(u any) algorithm.IDType {
			return algorithm.IDType(u.(int))
		}).ToList()
	})
	println(computed)
	var parallel []int
	computed = util.BeginEndTimeComputed(func() {
		parallel = stream.ToStream(&users).Parallel().OrderById(false, func(u any) algorithm.IDType {
			return algorithm.IDType(u.(int))
		}).ToList()
	})
	println(computed)
	for i := 0; i < len(unParallel); i++ {
		//println(unParallel[i], parallel[i])
		if unParallel[i] != parallel[i] {
			panic("unparallel not equal parallel" + convert.ConvertToString(i))
		}
	}
}

func Test_FlatMap(t *testing.T) {
	users := userList[:]
	for i := 0; i < 10; i++ {
		users = append(users, userList[:]...)
	}
	//print("123")
	computed := util.BeginEndTimeComputed(func() {
		stream.ToStream(&users).FlatMap(func(user user) *stream.Stream[any, []any] {
			split := strings.Split(user.Name, "")
			res := make([]any, len(split))
			for i, v := range split {
				res[i] = v
			}
			return stream.NewStream(&res)
		}).ForEach(func(item any) {
			print(item.(string))
		})
	})
	println()
	println(computed)
	computed = util.BeginEndTimeComputed(func() {
		stream.ToParallelStream(&users).FlatMap(func(user user) *stream.Stream[any, []any] {
			split := strings.Split(user.Name, "")
			res := make([]any, len(split))
			for i, v := range split {
				res[i] = v
			}
			return stream.NewStream(&res)
		}).ForEach(func(item any) {
			print(item.(string))
		})
	})
	println()
	println(computed)

}
