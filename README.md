# Katool - Go

Katool 是一个功能丰富的 Go 工具库，借鉴了 Java 生态中的优秀设计，提供了诸多实用功能，包括流式处理、IOC 容器、锁支持、集合操作、类型转换等。

## 主要功能

- [Stream 流式处理](#stream)
- [集合操作](#lists)
- [IOC 容器](#ioc)
- [LockSupport](#locksupport)
- [数据转换](#convert)
- [Web 爬虫](#web_crawler)
- [文件操作](#file)
- [DB 操作工具](#db)
- [日志工具](#log)
- [算法工具](#algorithm)

## 安装

```bash
go get github.com/karosown/katool-go
```

## Stream

支持像Java一样的Stream流，但是由于Go不支持方法泛型，在使用过程中需要自行处理

包含map、reduce、filter、groupBy、distinct、sort、flatMap、orderBy等操作, 支持foreach遍历，也支持collect进行自定义函数逻辑采集，同时支持异步流调用（目前异步流支持map、reduce、filter、flatMap、ToMap、ForEach、OrderBy、OrderByID、Sort等方法

具体使用查看test

```go
stream.ToStream(&data).Filter(func(i any) bool { 
    return i.(*response.AppleCrashReport).Date == "2024-11-04"
}).ToList()
```

## Lists

### list.Partition
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

## IOC

IoC (控制反转) 容器实现，类似 Spring 的依赖注入功能，支持注册、获取和强制注册等操作。

```go
// 注册一个值
ioc.RegisterValue("key", value)

// 获取一个值，如果不存在则使用默认值
value := ioc.GetDef("key", defaultValue)

// 使用函数注册
ioc.Register("key", func() any {
    return newValue
})
```

## LockSupport

类似Java LockSupport 对协程进行控制，但是目前不支持自动恢复协程（预计下个版本加入）

```go
lockSupport := lock.NewLockSupport()
// 阻塞当前协程
lockSupport.Park()
// 恢复协程
lockSupport.Unpark()
```

还提供了 Synchronized 方法，类似 Java 的 synchronized 关键字：

```go
mutex := &sync.Mutex{}
lock.Synchronized(mutex, func() {
    // 在锁保护下执行代码
})
```

## Convert

数据转换工具，支持类型属性复制、任意类型转字符串、数组转换等功能。

```go
// 复制属性
dest, err := convert.CopyProperties(src, &dest{})

// 转换为字符串
str := convert.ToString(value)

// 类型转换
result := convert.Convert(srcSlice, func(item SrcType) DestType {
    // 转换逻辑
    return converted
})
```

## Web_Crawler

网页爬虫工具，支持从URL读取内容，解析文章等功能。

```go
// 示例代码，具体请查看文档
article, err := web_crawler.FromURL("https://example.com", timeout)
```

## 更多功能

更多功能和详细用法，请查看各模块文档：

- [Stream 流式处理](docs/stream.md)
- [集合操作](docs/lists.md)
- [IOC 容器](docs/ioc.md)
- [锁支持](docs/lock.md)
- [数据转换](docs/convert.md)
- [Web 爬虫](docs/web_crawler.md)
- [文件操作](docs/file.md)
- [日志工具](docs/log.md)
- [算法工具](docs/algorithm.md)

## 贡献

欢迎提交 PR 和 Issues!
