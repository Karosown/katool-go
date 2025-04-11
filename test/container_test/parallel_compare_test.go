package container_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/util"
)

func orderById(flag bool, ids []int) []int {
	toStream := stream.ToStream(&ids)
	if flag {
		toStream.Parallel()
	}
	return toStream.OrderById(false, func(u any) algorithm.IDType {
		return algorithm.IDType(u.(int))
	}).ToList()
}

func BenchmarkOrderByIdUnParallel(b *testing.B) {
	users := make([]int, 0)
	for i := 0; i < 6e6+5e5+3e3+2e2+1e1; i++ {
		users = append(users, rand.Int()%10000)
	}
	orderById(false, users)
}

func BenchmarkOrderByIdParallel(b *testing.B) {
	users := make([]int, 0)
	for i := 0; i < 6e6+5e5+3e3+2e2+1e1; i++ {
		users = append(users, rand.Int()%10000)
	}
	orderById(true, users)
}

func Test_OrderBy_ID(t *testing.T) {
	users := make([]int, 0)
	for i := 0; i < 6e6+5e5+3e3+2e2+1e1; i++ {
		users = append(users, rand.Int()%10000)
	}
	var unParallel []int
	computed := util.BeginEndTimeComputed(func() {
		unParallel = orderById(false, users)
	})
	println(computed)
	var parallel []int
	computed = util.BeginEndTimeComputed(func() {
		parallel = orderById(true, users)
	})
	println(computed)
	for i := 0; i < len(unParallel); i++ {
		//println(unParallel[i], parallel[i])
		if unParallel[i] != parallel[i] {
			sys.Panic("unparallel not equal parallel" + convert.ToString(i))
		}
	}
}

func Test_OrderBy(t *testing.T) {
	users := userList[:]
	for i := 0; i < 100000; i++ {
		users = append(users, userList[:]...)
	}
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
			sys.Panic("unparallel not equal parallel" + convert.ToString(i))
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
