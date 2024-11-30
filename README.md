# Katool - Go

## Stream

支持像Java一样的Stream流，但是由于Go不支持方法泛型，在使用过程中需要自行处理

包含map、reduce、filter、groupBy、distinct、sort、flatMap、orderBy等操作, 支持foreach遍历，也支持collect进行自定义函数逻辑采集，同时支持异步流调用（目前异步流支持map、reduce、fillter、flatMap、ToMap、ForEach等方法

具体使用查看test

```go
package container_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool/algorithm"
	"github.com/karosown/katool/container/stream"
	"github.com/karosown/katool/convert"
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
	// ...过多无效内容，自行填充
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
	by := stream.ToStream(&users).OrderBy(true, func(u any) algorithm.HashType {
		return algorithm.HashType(u.(user).Name)
	}).ToList()
	println(by)
}

func Test_FlatMap(t *testing.T) {
	users := userList[:]
	stream.ToStream(&users).Parallel().FlatMap(func(user user) *stream.Stream[any, []any] {
		split := strings.Split(user.Name, "")
		res := make([]any, len(split))
		for i, v := range split {
			res[i] = v
		}
		return stream.NewStream(&res)
	}).ForEach(func(item any) {
		println(item.(string))
	})

}
```



## list.Partition
进行分片操作，可以在分片后配合convert包转换为stream流，同时也支持对分片进行foreach遍历（内部采用协程、同时可以控制协程大小），同时支持判断批次（用于数据原始位置判断）

```go
func Test_Partition(t *testing.T) {
	sum := convert.PatitonToStreamp(lists.Partition(userList[:], 15)).
		Parallel().
		Map(func(i []user) any {
			return stream.ToStream(&i).Map(func(user user) any {
				properties, _ := convert.CopyProperties(user, &userVo{})
				return *properties
			}).ToList()
		}).
		Reduce("", func(cntValue any, nxt any) any {
			anies := nxt.([]any)
			return stream.ToStream(&anies).Reduce(cntValue, func(sumValue any, nxt any) any {
				return sumValue.(string) + nxt.(userVo).Name
			}, func(sum1, sum2 any) any {
				return sum1.(string) + sum2.(string)
			})
		}, func(sum1, sum2 any) any {
			return sum1.(string) + sum2.(string)
		})
	println(sum.(string))
}

func Test_ForEach(t *testing.T) {
	lists.Partition(userList[:], 15).ForEach(func(pos int, automicDatas []user) error {
		fmt.Println("分批处理 第" + convert.ConvertToString(pos) + "批")
		stream.ToStream(&automicDatas).ForEach(func(data user) {
			fmt.Println(data)
		})
		return nil
	}, true, lynx.NewLimiter(3))
	//println(sum.(string))
}
```



## lockSupport
类似Java LockSupport 对协程进行控制，但是目前不支持自动恢复协程（预计下个版本加入）

## remote

二次封装的resty请求库，适用于工作流处理，例如请求链处理，完成请求链路即可，对于google ads和bing ads等请求链路有用

## convert
一些数据转换的util

