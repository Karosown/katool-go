# Katool - Go

## Stream

支持像Java一样的Stream流，但是由于Go不支持方法泛型，在使用过程中需要自行处理

包含map、reduce、filter、groupBy、distinct、sort、flatMap、orderBy等操作, 支持foreach遍历，也支持collect进行自定义函数逻辑采集，同时支持异步流调用（目前异步流支持map、reduce、fillter、flatMap、ToMap、ForEach、OrderBy、OrderByID、Sort等方法

具体使用查看test

```go
stream.ToStream(&data).Filter(func(i any) bool { return i.(*response.AppleCrashReport).Date == "2024-11-04";}).ToList()
	
})
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
		fmt.Println("分批处理 第" + convert.ToString(pos) + "批")
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

