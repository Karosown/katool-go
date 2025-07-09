package container_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/random"
	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/helper/jsonhp"
	"github.com/karosown/katool-go/sys"
)

type user struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Sex   int    `json:"sex"`
	Money int    `json:"money"`
	Class string `json:"class"`
	Id    int    `json:"id"`
}
type userVo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  int    `json:"sex"`
	Id   int    `json:"id"`
}

var userList []user

func TestOfStream(t *testing.T) {

	arr := []int{1, 3, 2, 3, 3, 3, 3, 1, 3, 2, 3, 3, 3, 3, 1, 3, 2, 3, 3, 3, 3, 1, 3, 2, 3, 3, 3, 3, 1, 3, 2, 3, 3, 3, 3, 1, 3, 2, 3, 3, 3, 3}

	distinct := stream.ToStream(&arr).
		Parallel().
		Filter(func(i int) bool {
			return i > 1
		}).Map(func(item int) any {
		return strconv.Itoa(item) + "w "
	}).DistinctBy(algorithm.HASH_WITH_JSON_SUM)

	fmt.Println(distinct.Reduce("", func(cntValue any, nxt any) any {
		return cntValue.(string) + nxt.(string)
	}, func(sum1, sum2 any) any {
		return sum1.(string) + sum2.(string)
	}))
	list := distinct.ToOptionList()
	list.ForEach(func(s any) {
		fmt.Println(s)
	})

	toMap := stream.ToStream(&arr).
		//Parallel().
		ToMap(func(index int, item int) any {
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
	userStream := stream.ToStream(&ul).Parallel()
	println(userStream.Count())
	// 排序
	stream.ToStream(&ul).
		Parallel().
		Sort(func(a user, b user) bool { return a.Id < b.Id }).ForEach(func(item user) { println(convert.ToString(item.Id) + " " + item.Name) })
	// 求和
	totalMoney := userStream.Reduce(int64(0), func(cntValue any, nxt user) any { return cntValue.(int64) + int64(nxt.Money) }, func(sum1, sum2 any) any {
		return sum1.(int64) + sum2.(int64)
	})
	println(totalMoney.(int64))
	// 过滤
	userStream.Filter(func(item user) bool { return item.Sex != 0 }).DistinctBy(algorithm.HASH_WITH_JSON_MD5).ToOptionList().ForEach(func(item user) { println(item.Name) })
	// 转换
	s := userStream.Map(func(item user) any {
		properties, err := convert.CopyProperties(&item, &userVo{})
		if err != nil {
			sys.Panic(err)
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
		return user.Class
	})
	println(by)
}

func init() {
	classes := []string{"一班", "二班", "三班", "四班", "五班"}
	userList = make([]user, 0)
	for i := 0; i < 100; i++ {
		userList = append(userList, user{
			Name:  random.RandString(10),
			Class: classes[rand.Int()%len(classes)],
			Age:   rand.Intn(100),
			Sex:   rand.Intn(2),
		})
		time.Sleep(1)
	}
}
func Test(t *testing.T) {
	by := stream.ToStream(&userList).Parallel().GroupBy(func(user user) any { return user.Class })
	maputil.ForEach(by, func(key any, value []user) {
		println(key.(string))
		toStream := stream.ToStream(&value).Parallel()
		toStream.ForEach(func(item user) {
			fmt.Println(item)
		})
		reduce := toStream.Reduce(0, func(cntValue any, nxt user) any {
			return cntValue.(int) + nxt.Sex
		}, func(cntValue any, nxt any) any {
			return cntValue.(int) + nxt.(int)
		})
		println("男生总数：" + convert.ToString(reduce))
		reduce = toStream.Reduce(0, func(cntValue any, nxt user) any {
			return cntValue.(int) + (nxt.Sex ^ 1)
		}, func(cntValue any, nxt any) any {
			return cntValue.(int) + nxt.(int)
		})
		println("女生总数：" + convert.ToString(reduce))
	})
	toStream := stream.ToStream(&userList).Parallel()
	count := toStream.Filter(func(item user) bool { return item.Age >= 60 }).Count()
	println("年龄大于等于60岁的共" + convert.ToString(count) + "人")
	count = toStream.Filter(func(item user) bool { return item.Age < 60 }).Count()
	println("年龄小于60岁的共" + convert.ToString(count) + "人")
	toStream.FlatMap(func(i user) *stream.Stream[any, []any] {
		array := convert.ToAnySlice([]rune(i.Name))
		return stream.ToStream(&array)
	}).ForEach(func(item any) {
		print(string(item.(rune)) + " ")
	})
}
func Test_Sorted(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 1, 2, 8, 9, 10}
	for i := 0; i < 100; i++ {
		arr = append(arr, 1, 2, 3, 4, 5, 6, 7, 1, 2, 8, 9, 10)
	}
	stream.ToStream(&arr).Parallel().Sort(func(a, b int) bool { return a < b }).ForEach(func(item int) {
		fmt.Println(item)
	})
}
func Test_Distinct(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 10000; i++ {
		arr = append(arr, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	}
	stream.ToStream(&arr).Parallel().Distinct().ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_FromAnySlice(t *testing.T) {
	slice := convert.ToAnySlice(userList)
	// 使用FromAnySlice方法
	stream.FromAnySlice[user, []user](slice).ForEach(func(item user) {
		fmt.Println(item)
	})
	// 使用Cast方法，对于Map可以配套使用
	anyStream := stream.ToStream(&slice)
	stream.Cast[user](anyStream).ForEach(func(item user) {
		fmt.Println(item)
	})
}

func Test_Merge(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	arr1 := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	stream.ToStream(&arr).Parallel().Merge(arr1).ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_Intersect(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	arr1 := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 5, 6, 7}
	stream.ToStream(&arr).Parallel().Intersect(arr1).ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_Difference(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	arr1 := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 5, 6, 7}
	stream.ToStream(&arr).Parallel().Difference(arr1).ForEach(func(item int) {
		fmt.Println(item)
	})
}

// 其他玩儿法
func Test_Intersect_filter(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream.ToStream(&arr).Parallel().Intersect(arr1, func(a, b int) bool {
		return a < b
	}).ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_Sub(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream.ToStream(&arr).Sub(-2, -1).ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_Skip(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream.ToStream(&arr).Skip(3).ForEach(func(item int) {
		fmt.Println(item)
	})
}

func Test_JSONLINE(t *testing.T) {
	println(jsonhp.ToJsonLine[user](userList))

	marshal, err := json.Marshal(userList)
	if err != nil {
		t.Error(err)
	}
	println(jsonhp.ToJsonLine[user](marshal))

	println(jsonhp.ToJsonLine[user](stream.Cast[[]byte](stream.ToStream(&userList).Map(func(i user) any {
		bytes, err2 := json.Marshal(i)
		if err2 != nil {
			return []byte{}
		}
		return bytes
	})).ToList()))
}
func Test_ToJsonLine(t *testing.T) {
	println(jsonhp.ToJsonLine[user](userList))
}
func Test_ToJSON(t *testing.T) {
	println(jsonhp.ToJSON(jsonhp.ToJsonLine[user](userList)))
}

func Test_FixJson(t *testing.T) {
	strs := []string{
		"{a:1,C:'1',b:[2,3,5,",
		`{"name""SINpUGSroN""age"37"sex"0"money"0"class""一班""id"0}
{"name""RqPaNmmEcJ""age"31"sex"1"money"0"class""四班""id"0}
{"name""jzHCCSepro""age"17"sex":1,"money":0,"class":"五班","id":0}
{"name":"uHZBCFNrGE","age":6,"sex":0"money":0,"class":"三班","id":0}
{"name":"SGwKKLCfPL","age":19,"sex":0,"money":0,"class":"一班","id":0}
{"name":"zENFnzjImU","age":98,"sex":1,"money":0,"class":"二班","id":0}`,
	}

	stream.ToStream(&strs).ForEach(func(item string) {
		fmt.Println(jsonhp.FixJson(item))
	})
}
