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

	i := atomic.Uint32{}

	lists.Partition(userList[:], 15).ForEach(func(pos int, automicDatas []user) error {
		i.Add(1)
		fmt.Println("分批处理" + convert.ConvertToString(i.Load()))
		stream.ToStream(&automicDatas).ForEach(func(data user) {
			fmt.Println(data)
		})
		return nil
	}, true, lynx.NewLimiter(3))
	//println(sum.(string))
}
