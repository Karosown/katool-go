# Katool-Go

<div align="center">

<img src="logo.png" alt="Katool-Go Logo" width="400">

<h1>🛠️ Katool-Go</h1>

<p>
  <a href="https://pkg.go.dev/github.com/karosown/katool-go"><img src="https://pkg.go.dev/badge/github.com/karosown/katool-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/karosown/katool-go"><img src="https://goreportcard.com/badge/github.com/karosown/katool-go" alt="Go Report Card"></a>
  <a href="https://github.com/karosown/katool-go/releases"><img src="https://img.shields.io/github/v/release/karosown/katool-go" alt="GitHub release"></a>
  <a href="https://github.com/karosown/katool-go/blob/main/LICENSE"><img src="https://img.shields.io/github/license/karosown/katool-go" alt="License"></a>
  <a href="https://golang.org/dl/"><img src="https://img.shields.io/github/go-mod/go-version/karosown/katool-go" alt="Go Version"></a>
</p>

<b><i>一个功能丰富的 Go 工具库，借鉴 Java 生态优秀设计，为 Go 开发提供全方位支持 （以下内容为介绍，具体使用建议看test文件）</i></b>

</div>

<hr>

## 📋 目录

- [📝 简介](#简介)
- [✨ 特性](#特性)
- [📦 安装](#安装)
- [🚀 快速开始](#快速开始)
- [🔧 核心模块](#核心模块)
  - [📚 容器与集合](#容器与集合)
  - [🌊 流式处理](#流式处理)
  - [🔄 数据转换](#数据转换)
  - [💉 依赖注入](#依赖注入)
  - [🔒 并发控制](#并发控制)
  - [🕸️ Web爬虫](#web爬虫)
  - [📁 文件操作](#文件操作)
  - [💾 数据库支持](#数据库支持)
  - [🌐 网络通信](#网络通信)
  - [📝 日志系统](#日志系统)
  - [⚙️ 算法工具](#算法工具)
  - [🔤 文本处理](#文本处理)
  - [🧰 辅助工具](#辅助工具)
- [💡 最佳实践](#最佳实践)
- [👥 贡献指南](#贡献指南)
- [📄 许可证](#许可证)

<hr>

## 📝 简介

**Katool-Go** 是一个综合性的 Go 语言工具库，旨在提供丰富的功能组件和实用工具，帮助开发者提高开发效率。它借鉴了 Java 生态中的成熟设计模式和经验，同时充分利用 Go 语言的特性，如并发、泛型等，提供了一系列易用且高效的工具。

本库的设计理念是：**模块化、可组合、高性能**，适用于各种规模的 Go 项目。无论是构建微服务、Web应用，还是数据处理系统，Katool-Go 都能提供有力支持。

<hr>

## ✨ 特性

Katool-Go 提供以下核心特性：

<table>
  <tr>
    <td><b>🌊 流式处理</b></td>
    <td>提供类似 Java 8 Stream API 的链式操作，支持 map/filter/reduce/collect 等操作</td>
  </tr>
  <tr>
    <td><b>📚 容器与集合</b></td>
    <td>增强的集合类型，如 XMap、HashBasedMap、Optional 等</td>
  </tr>
  <tr>
    <td><b>💉 依赖注入</b></td>
    <td>轻量级 IOC 容器，支持组件注册、获取和生命周期管理</td>
  </tr>
  <tr>
    <td><b>🔒 并发控制</b></td>
    <td>协程控制工具，如 LockSupport、同步锁封装等</td>
  </tr>
  <tr>
    <td><b>🔄 数据转换</b></td>
    <td>对象属性复制、类型转换、序列化等工具</td>
  </tr>
  <tr>
    <td><b>🕸️ Web爬虫</b></td>
    <td>网页内容抓取、解析、RSS订阅支持等</td>
  </tr>
  <tr>
    <td><b>📁 文件操作</b></td>
    <td>文件下载、序列化、路径工具等</td>
  </tr>
  <tr>
    <td><b>💾 数据库支持</b></td>
    <td>MongoDB工具、分页查询等</td>
  </tr>
  <tr>
    <td><b>🌐 网络通信</b></td>
    <td>HTTP客户端、RESTful支持、SSE、OAuth2等</td>
  </tr>
  <tr>
    <td><b>📝 日志系统</b></td>
    <td>多级别日志、适配器、自定义格式等</td>
  </tr>
  <tr>
    <td><b>⚙️ 算法工具</b></td>
    <td>数组操作、哈希计算、二进制处理等</td>
  </tr>
  <tr>
    <td><b>🔤 文本处理</b></td>
    <td>中文分词、文本分析等</td>
  </tr>
  <tr>
    <td><b>🧰 辅助工具</b></td>
    <td>日期、随机数、调试等实用工具</td>
  </tr>
</table>

<hr>

## 📦 安装

使用 `go get` 安装最新版本：

```bash
go get -u github.com/karosown/katool-go
```

> ⚠️ 要求 Go 版本 >= 1.23.1

<hr>

## 🚀 快速开始

下面是几个简单示例，展示 Katool-Go 的基本用法：

<details open>
<summary><b>🌊 流式处理</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/convert"
	"strconv"
)

// 定义用户结构
type user struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Sex   int    `json:"sex"`
	Money int    `json:"money"`
	Class string `json:"class"`
	Id    int    `json:"id"`
}

func main() {
	users := []user{
		{Name: "Alice", Age: 25, Sex: 1, Money: 1000, Class: "A", Id: 1},
		{Name: "Bob", Age: 30, Sex: 0, Money: 1500, Class: "B", Id: 2},
		{Name: "Charlie", Age: 35, Sex: 0, Money: 2000, Class: "A", Id: 3},
		{Name: "David", Age: 40, Sex: 1, Money: 2500, Class: "B", Id: 4},
	}
	
	// 并行流处理
	userStream := stream.ToStream(&users).Parallel()
	
	// 计算总人数
	fmt.Println("总人数:", userStream.Count())
	
	// 按ID排序
	stream.ToStream(&users).Parallel().
		Sort(func(a, b user) bool { 
			return a.Id < b.Id 
		}).ForEach(func(item user) { 
			fmt.Println(item.Id, item.Name) 
		})
	
	// 计算总金额
	totalMoney := userStream.Reduce(int64(0), 
		func(sum any, u user) any { 
			return sum.(int64) + int64(u.Money) 
		}, 
		func(sum1, sum2 any) any {
			return sum1.(int64) + sum2.(int64)
		})
	fmt.Println("总金额:", totalMoney)
	
	// 按班级分组
	groups := stream.ToStream(&users).GroupBy(func(u user) any {
		return u.Class
	})
	
	// 输出各班级成员
	for class, members := range groups {
		fmt.Printf("班级: %s, 人数: %d\n", class, len(members))
	}

	// 过滤操作
	filtered := userStream.Filter(func(u user) bool {
		return u.Age > 30
	}).ToList()

	// 映射操作
	userStream.Map(func(u user) any {
		return u.Name
	}).ForEach(func(name any) {
		fmt.Println(name)
	})

	// 排序操作
	sorted := stream.ToStream(&users).Sort(func(a, b user) bool {
		return a.Age < b.Age
	}).ToList()

	// 元素去重
	numbers := []int{1, 2, 3, 1, 2, 3, 4, 5}
	stream.ToStream(&numbers).Distinct().ToList()

	// 自定义去重
	userStream.DistinctBy(algorithm.HASH_WITH_JSON_MD5).ToList()

	// 字符串操作
	arr := []int{1, 2, 3}
	result := stream.ToStream(&arr).Map(func(i int) any {
		return strconv.Itoa(i) + "w"
	}).Reduce("", func(sum any, item any) any {
		return sum.(string) + item.(string)
	}, func(sum1, sum2 any) any {
		return sum1.(string) + sum2.(string)
	})
	fmt.Println(result)
}
```
</details>

<details>
<summary><b>📚 容器操作</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/container/xmap"
	"encoding/json"
)

func main() {
	// 创建普通Map
	m := xmap.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)
	
	// 获取和验证
	val, exists := m.Get("one")
	fmt.Printf("键'one'存在: %v, 值: %d\n", exists, val) // true, 1
	
	// 删除元素
	m.Delete("two")
	fmt.Printf("Map大小: %d\n", m.Len()) // 2
	
	// 遍历
	m.ForEach(func(k string, v int) {
		fmt.Printf("%s: %d\n", k, v)
	})
	
	// 创建线程安全Map
	sm := xmap.NewSafeMap[string, int]()
	sm.Set("a", 1)
	sm.Set("b", 2)
	
	// 安全地获取或存储
	val, loaded := sm.LoadOrStore("a", 100)
	fmt.Printf("键'a'已存在: %v, 值: %d\n", loaded, val) // true, 1
	
	val, loaded = sm.LoadOrStore("c", 3)
	fmt.Printf("键'c'已存在: %v, 值: %d\n", loaded, val) // false, 3
	
	// 创建有序Map
	sortedMap := xmap.NewSortedMap[string, string]()
	sortedMap.Set("3", "three")
	sortedMap.Set("1", "one")
	sortedMap.Set("2", "two")
	
	// 序列化为JSON (按键排序)
	jsonBytes, _ := json.Marshal(sortedMap) 
}
```
</details>

<details>
<summary><b>🔒 并发控制</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/lock"
	"github.com/karosown/katool-go/container/stream"
	"time"
)

func main() {
	// 单个锁支持
	support := lock.NewLockSupport()
	
	go func() {
		fmt.Println("即将进入阻塞，等待异步唤醒")
		support.Park() // 阻塞当前协程，直到有人调用Unpark
		fmt.Println("唤醒成功")
	}()
	
	time.Sleep(time.Second) // 等待协程启动
	fmt.Println("主协程准备唤醒子协程")
	support.Unpark() // 解除阻塞
	
	// 多个LockSupport的管理
	locks := make([]*lock.LockSupport, 5)
	for i := 0; i < 5; i++ {
		locks[i] = lock.NewLockSupport()
		idx := i
		go func() {
			fmt.Printf("协程 %d 等待唤醒\n", idx)
			locks[idx].Park()
			fmt.Printf("协程 %d 被唤醒\n", idx)
		}()
	}
	
	// 依次唤醒所有协程
	for i, ls := range locks {
		fmt.Printf("唤醒协程 %d\n", i)
		ls.Unpark()
		time.Sleep(100 * time.Millisecond) // 间隔唤醒
	}
}
```
</details>

<details>
<summary><b>🔤 文本处理</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/words/jieba"
)

func main() {
	// 创建分词客户端
	jb := jieba.New()
	defer jb.Free() // 使用完必须释放资源
	
	// 精确模式分词
	text := "我测试一下中文分词 Hello World"
	words := jb.Cut(text)
	fmt.Println(words) // ["我", "测试", "一下", "中文", "分词", "Hello", "World"]
	
	// 全模式分词
	text = "下面是一个简洁的Go语言SDK"
	allWords := jb.CutAll(text)
	fmt.Println(allWords) // 包含所有可能的分词结果
	
	// 搜索引擎模式分词 (更细粒度，适合搜索)
	searchWords := jb.CutForSearch("清华大学位于北京市")
	fmt.Println(searchWords) // ["清华", "华大", "大学", "位于", "北京", "北京市"]
	
	// 词频统计
	wordFreq := jb.CutAll("重复的词重复的词重复的词").Frequency()
	for word, count := range wordFreq {
		fmt.Printf("%s: %d次\n", word, count)
	}
}
```
</details>

<hr>

## 🔧 核心模块

### 📚 容器与集合

<details>
<summary><b>✨ XMap - 增强的映射</b></summary>

XMap 提供了比标准 map 更丰富的功能：

```go
import "github.com/karosown/katool-go/container/xmap"

// 创建不同类型的Map
regularMap := xmap.NewMap[string, int]()      // 普通Map
safeMap := xmap.NewSafeMap[string, int]()     // 线程安全Map
sortedMap := xmap.NewSortedMap[string, int]() // 有序Map

// 设置值
regularMap.Set("one", 1)
safeMap.Set("one", 1)
sortedMap.Set("one", 1)

// 安全地获取或存储 (仅SafeMap支持)
value, loaded := safeMap.LoadOrStore("two", 2) // 如果不存在则存储
if !loaded {
	fmt.Println("键'two'不存在，已存储值:", value) // 2
}

// 获取和删除 (仅SafeMap支持)
value, exists := safeMap.LoadAndDelete("one")
if exists {
	fmt.Println("获取并删除键'one'的值:", value) // 1
}

// 遍历
regularMap.ForEach(func(k string, v int) {
	fmt.Printf("%s: %d\n", k, v)
})

// JSON序列化 (SortedMap按键排序)
jsonBytes, _ := json.Marshal(sortedMap) 
```
</details>

<details>
<summary><b>🔑 HashBasedMap - 双层键映射</b></summary>

HashBasedMap 支持使用两个键来索引值：

```go
import "github.com/karosown/katool-go/container/hash_based_map"

// 创建双层映射
m := hash_based_map.NewHashBasedMap[string, int, User]()

// 设置值
m.Set("group1", 1, User{Name: "Alice"})
m.Set("group1", 2, User{Name: "Bob"})
m.Set("group2", 1, User{Name: "Charlie"})

// 获取值
user, exists := m.Get("group1", 1) // user={Name:"Alice"}, exists=true

// 获取所有键
firstKeys, secondKeys := m.Keys()
```
</details>

<details>
<summary><b>📦 Optional - 避免空指针</b></summary>

Optional 提供了处理可能为空值的安全方式：

```go
import "github.com/karosown/katool-go/container/optional"

// 创建Optional
opt := optional.Of("value")
emptyOpt := optional.Empty[string]()

// 检查是否存在值
if opt.IsPresent() {
    value := opt.Get()
    fmt.Println(value) // 输出: value
}

// 提供默认值
value := emptyOpt.OrElse("default") // value="default"

// 条件执行
opt.IfPresent(func(v string) {
    fmt.Println("Value present:", v)
})
```
</details>

### 🌊 流式处理

<details>
<summary><b>🔄 基本操作</b></summary>

```go
import (
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/algorithm"
	"strconv"
)

// 准备数据
users := []user{
	{Name: "Alice", Age: 25, Class: "A"},
	{Name: "Bob", Age: 30, Class: "B"},
	{Name: "Charlie", Age: 35, Class: "A"},
	{Name: "David", Age: 40, Class: "B"},
}

// 创建流 (可选并行处理)
s := stream.ToStream(&users).Parallel()

// 过滤操作
filtered := s.Filter(func(u user) bool {
	return u.Age > 30
}).ToList() // [{Name:Charlie Age:35...}, {Name:David Age:40...}]

// 映射操作
s.Map(func(u user) any {
	return u.Name
}).ForEach(func(name any) {
	fmt.Println(name) // 输出所有名字
})

// 排序操作
sorted := stream.ToStream(&users).Sort(func(a, b user) bool {
	return a.Age < b.Age // 按年龄升序
}).ToList()

// 元素去重
numbers := []int{1, 2, 3, 1, 2, 3, 4, 5}
stream.ToStream(&numbers).Distinct().ToList() // [1, 2, 3, 4, 5]

// 自定义去重 (使用自定义哈希函数)
s.DistinctBy(algorithm.HASH_WITH_JSON_MD5).ToList()

// 字符串操作
arr := []int{1, 2, 3}
result := stream.ToStream(&arr)
	.Map(func(i int) any {
		return strconv.Itoa(i) + "w" // 转为字符串并添加后缀
	})
	.Reduce("", func(sum any, item any) any {
		return sum.(string) + item.(string) // 拼接字符串
	}, func(sum1, sum2 any) any {
		return sum1.(string) + sum2.(string) // 拼接结果合并
	})
// result = "1w2w3w"
```
</details>

<details>
<summary><b>🚀 高级操作</b></summary>

```go
// 分组操作
classGroups := stream.ToStream(&users).GroupBy(func(u user) any {
	return u.Class
}) // map["A":[用户列表], "B":[用户列表]]

// 对每个分组进行统计
for class, members := range classGroups {
	fmt.Printf("班级: %s, 人数: %d\n", class, len(members))
	
	// 对每个分组创建流进行处理
	maleCount := stream.ToStream(&members).Reduce(0, 
		func(count any, u user) any {
			return count.(int) + u.Sex // 假设Sex=0为女性，Sex=1为男性
		}, 
		func(a, b any) any {
			return a.(int) + b.(int)
		})
	
	fmt.Printf("  男生人数: %d\n", maleCount)
	fmt.Printf("  女生人数: %d\n", len(members) - maleCount.(int))
}

// 扁平化操作 (将多个集合合并处理)
nameChars := stream.ToStream(&users).FlatMap(func(u user) *stream.Stream[any, []any] {
	// 将每个用户名拆分为字符
	chars := []rune(u.Name)
	array := convert.ToAnySlice(chars)
	return stream.ToStream(&array)
}).ToList()
// 结果为所有用户名中的字符列表

// 转换为Map
userMap := stream.ToStream(&users).ToMap(
	func(index int, u user) any {
		return u.Id // 键
	}, 
	func(index int, u user) any {
		return u.Name // 值
	}
) // map[1:"Alice" 2:"Bob" 3:"Charlie" 4:"David"]

// 类型安全转换
anySlice := convert.ToAnySlice(users)
typedUsers := stream.FromAnySlice[user, []user](anySlice).ToList()
```
</details>

<details>
<summary><b>📊 收集操作</b></summary>

```go
// 求和统计
sum := stream.ToStream(&users).Reduce(0, 
	func(acc any, u user) any {
		return acc.(int) + u.Age
	}, 
	func(a, b any) any {
		return a.(int) + b.(int)
	}
).(int) // sum=130

// 统计元素数量
count := stream.ToStream(&users).Count() // 4

// 条件统计
seniorCount := stream.ToStream(&users)
	.Filter(func(u user) bool { 
		return u.Age >= 60 
	})
	.Count() // 年龄大于等于60的人数

// 聚合统计
totalMoney := stream.ToStream(&users).Reduce(int64(0),
	func(sum any, u user) any {
		return sum.(int64) + int64(u.Money)
	},
	func(sum1, sum2 any) any {
		return sum1.(int64) + sum2.(int64)
	}
).(int64)
fmt.Printf("总金额: %d\n", totalMoney)
```
</details>

### 🔄 数据转换

<details>
<summary><b>📋 属性复制</b></summary>

```go
import "github.com/karosown/katool-go/convert"

// 源对象和目标对象
type Source struct {
	ID   int
	Name string
	Age  int
}

type Destination struct {
	ID   int
	Name string
	Age  int
	Extra string // 额外字段
}

// 复制属性
src := Source{ID: 1, Name: "Alice", Age: 30}
dest := &Destination{Extra: "Additional info"}

result, err := convert.CopyProperties(src, dest)
// result={ID:1 Name:"Alice" Age:30 Extra:"Additional info"}
```
</details>

<details>
<summary><b>🔄 类型转换</b></summary>

```go
// 转换为字符串
str := convert.ToString(123) // "123"
str = convert.ToString(true) // "true"
str = convert.ToString([]int{1, 2, 3}) // "[1,2,3]"

// 类型批量转换
type UserDTO struct {
	ID   string
	Name string
}

users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}

dtos := convert.Convert(users, func(u User) UserDTO {
	return UserDTO{
		ID:   convert.ToString(u.ID),
		Name: u.Name,
	}
})
// dtos=[{ID:"1" Name:"Alice"}, {ID:"2" Name:"Bob"}]

// 任意类型转换
anySlice := convert.ToAnySlice(users) // []any{User{...}, User{...}}
typedSlice := convert.FromAnySlice[User](anySlice) // []User{...}
```
</details>

### 💉 依赖注入

<details>
<summary><b>🏭 IOC容器</b></summary>

```go
import "github.com/karosown/katool-go/container/ioc"

// 定义接口和实现
type UserRepository interface {
	FindByID(id int) string
}

type UserRepositoryImpl struct{}

func (r *UserRepositoryImpl) FindByID(id int) string {
	return fmt.Sprintf("User %d", id)
}

type UserService interface {
	GetUser(id int) string
}

type UserServiceImpl struct {
	Repository UserRepository
}

func (s *UserServiceImpl) GetUser(id int) string {
	return s.Repository.FindByID(id)
}

// 注册组件
ioc.RegisterValue("userRepo", &UserRepositoryImpl{})

// 注册工厂方法
ioc.Register("userService", func() any {
	repo := ioc.Get("userRepo").(UserRepository)
	return &UserServiceImpl{Repository: repo}
})

// 获取服务
service := ioc.Get("userService").(UserService)
result := service.GetUser(1) // "User 1"

// 获取带默认值
repo := ioc.GetDef("missingRepo", &UserRepositoryImpl{})
```
</details>

### 🔒 并发控制

<details>
<summary><b>⏱️ LockSupport</b></summary>

LockSupport 提供了类似 Java 的 park/unpark 机制，用于协程间的精确控制：

```go
import (
	"github.com/karosown/katool-go/lock"
	"fmt"
	"time"
)

// 创建LockSupport
ls := lock.NewLockSupport()

// 在单个协程中使用
go func() {
	fmt.Println("协程开始")
	ls.Park() // 阻塞当前协程，直到有人调用Unpark
	fmt.Println("协程继续执行") // 只有在Unpark后才会执行
}()

time.Sleep(time.Second) // 等待协程启动
fmt.Println("主协程准备唤醒子协程")
ls.Unpark() // 解除阻塞

// 多个LockSupport的管理
locks := make([]*lock.LockSupport, 10)
for i := 0; i < 10; i++ {
	locks[i] = lock.NewLockSupport()
	
	go func(i int, support *lock.LockSupport) {
		fmt.Printf("协程 %d 等待唤醒\n", i)
		support.Park()
		fmt.Printf("协程 %d 被唤醒\n", i)
	}(i, locks[i])
}

// 使用流式API依次唤醒
stream.ToStream(&locks).ForEach(func(support *lock.LockSupport) {
	fmt.Println("准备唤醒")
	support.Unpark()
	time.Sleep(100 * time.Millisecond) // 间隔唤醒
})
```
</details>

<details>
<summary><b>🔐 同步工具</b></summary>

```go
import (
	"github.com/karosown/katool-go/lock"
	"sync"
)

// 同步代码块
mutex := &sync.Mutex{}
counter := 0

lock.Synchronized(mutex, func() {
	// 临界区代码
	counter++
})

// 锁映射
lockMap := lock.NewLockMap()
// 适用于需要对不同对象分别加锁的场景
```
</details>

### 🕸️ Web爬虫

<details>
<summary><b>📄 页面抓取</b></summary>

```go
import "github.com/karosown/katool-go/web_crawler"

// 获取文章内容
article, err := web_crawler.FromURL("https://example.com", 30*time.Second)
if err == nil {
	fmt.Println("标题:", article.Title)
	fmt.Println("内容:", article.Content)
	fmt.Println("长度:", article.Length)
	fmt.Println("摘要:", article.Excerpt)
}

// 使用自定义请求选项
article, err = web_crawler.FromURLWithOptions("https://example.com", 
	30*time.Second, 
	func(r *http.Request) {
		r.Header.Set("User-Agent", "Mozilla/5.0...")
	})

// 解析路径
absolutePath := web_crawler.ParsePath("https://example.com/page", "./image.jpg")
// absolutePath = "https://example.com/page/image.jpg"
```
</details>

<details>
<summary><b>📰 RSS订阅</b></summary>

```go
import "github.com/karosown/katool-go/web_crawler/rss"

// 解析RSS源
feed, err := rss.ParseURL("https://example.com/feed.xml")
if err == nil {
	fmt.Println("源标题:", feed.Title)
	
	// 遍历条目
	for _, item := range feed.Items {
		fmt.Println("- 文章:", item.Title)
		fmt.Println("  链接:", item.Link)
		fmt.Println("  发布时间:", item.PubDate)
	}
}
```
</details>

### 📁 文件操作

<details>
<summary><b>⬇️ 文件下载</b></summary>

```go
import "github.com/karosown/katool-go/file/file_downloader"

// 下载文件
downloader := file_downloader.NewDownloader(
	file_downloader.WithTimeout(30*time.Second),
	file_downloader.WithRetries(3),
)
err := downloader.Download("https://example.com/file.zip", "local.zip")
```
</details>

<details>
<summary><b>💾 序列化</b></summary>

```go
import "github.com/karosown/katool-go/file/file_serialize"

// 序列化数据
data := map[string]any{"name": "Alice", "age": 30}
err = file_serialize.SerializeToFile(data, "data.json")

// 反序列化数据
var result map[string]any
err = file_serialize.DeserializeFromFile("data.json", &result)
```
</details>

### 💾 数据库支持

<details>
<summary><b>🍃 MongoDB</b></summary>

```go
import (
	"github.com/karosown/katool-go/db/xmongo"
	"go.mongodb.org/mongo-driver/bson"
)

// 创建MongoDB客户端
client := xmongo.NewClient("mongodb://localhost:27017")
coll := client.Database("test").Collection("users")

// 插入文档
_, err := coll.InsertOne(context.Background(), bson.M{
	"name": "Alice",
	"age":  30,
})
```
</details>

<details>
<summary><b>📄 分页查询</b></summary>

```go
import "github.com/karosown/katool-go/db/pager"

// 使用分页器查询
p := pager.NewPager(1, 10) // 第1页，每页10条
query := bson.M{"age": bson.M{"$gt": 25}}

cursor, err := coll.Find(context.Background(), query).
	Skip(p.Skip()).
	Limit(p.Limit()).
	Sort(bson.M{"name": 1}).
	Cursor()
```
</details>

### 🌐 网络通信

<details>
<summary><b>🌍 HTTP请求</b></summary>

```go
import "github.com/karosown/katool-go/net/http"

// 发送HTTP请求
client := http.NewRemoteRequest("https://api.example.com")
resp, err := client.Get("/users")
if err == nil {
	var users []User
	resp.UnmarshalJson(&users)
}

// POST请求与JSON
resp, err = client.PostJson("/users", User{Name: "Alice", Age: 30})
```
</details>

<details>
<summary><b>🔑 OAuth2支持</b></summary>

```go
// OAuth2支持
oauth := http.NewOAuth2Request(
	"https://api.example.com",
	"client_id",
	"client_secret",
	"https://auth.example.com/token",
)
resp, err = oauth.Get("/protected-resource")
```
</details>

<details>
<summary><b>📊 格式化</b></summary>

```go
import "github.com/karosown/katool-go/net/format"

// 格式化
jsonData := `{"name":"Alice","age":30}`
user := &User{}
err = format.Json.Decode([]byte(jsonData), user)

// 格式化响应
resp, err = client.GetWithFormat("/data", format.Json)
```
</details>

### 📝 日志系统

<details>
<summary><b>📋 基本日志</b></summary>

```go
import "github.com/karosown/katool-go/xlog"

// 基本日志
xlog.Info("这是一条信息")
xlog.Errorf("错误: %v", err)
```
</details>

<details>
<summary><b>📊 自定义日志</b></summary>

```go
import "github.com/karosown/katool-go/xlog"

// 创建自定义logger
logger := xlog.NewLogger(
	xlog.WithLevel(xlog.InfoLevel),
	xlog.WithFormat(xlog.JSONFormat),
	xlog.WithOutput("app.log"),
)

logger.Info("应用启动")
logger.WithFields(xlog.Fields{
	"user": "admin",
	"action": "login",
}).Info("用户登录")
```
</details>

### ⚙️ 算法工具

<details>
<summary><b>🔢 数组操作</b></summary>

```go
import "github.com/karosown/katool-go/algorithm"

// 合并有序数组
arr1 := []int{1, 3, 5}
arr2 := []int{2, 4, 6}
merged := algorithm.MergeSortedArrayWithEntity[int](func(a, b int) bool {
	return a < b // 升序
})(arr1, arr2)
// merged = [1, 2, 3, 4, 5, 6]

// 更多合并函数
result := algorithm.MergeSortedArrayWithPrimaryData[MyType](false, hashFunc)(array1, array2)
result := algorithm.MergeSortedArrayWithPrimaryId[MyType](false, idFunc)(array1, array2)
```
</details>

<details>
<summary><b>🔐 哈希计算</b></summary>

```go
// 哈希计算
data := map[string]any{"id": 123, "name": "test"}
hash := algorithm.HASH_WITH_JSON(data) // 使用JSON序列化计算哈希
md5Hash := algorithm.HASH_WITH_JSON_MD5(data) // 使用MD5计算哈希
sumHash := algorithm.HASH_WITH_JSON_SUM(data) // 使用累加计算哈希
```
</details>

### 🔤 文本处理

<details>
<summary><b>📝 中文分词</b></summary>

```go
import "github.com/karosown/katool-go/words/cgojieba"

// 创建分词客户端
jb := jieba.New()
defer jb.Free() // 使用完必须释放资源

// 精确模式分词
text := "我测试一下中文分词 Hello World"
words := jb.Cut(text)
fmt.Println(words) // ["我", "测试", "一下", "中文", "分词", "Hello", "World"]

// 全模式分词
text = "下面是一个简洁的Go语言SDK"
allWords := jb.CutAll(text)
fmt.Println(allWords) // 包含所有可能的分词结果

// 搜索引擎模式分词 (更细粒度，适合搜索)
searchWords := jb.CutForSearch("清华大学位于北京市")
fmt.Println(searchWords) // ["清华", "华大", "大学", "位于", "北京", "北京市"]

// 词频统计
wordFreq := jb.CutAll("下面是一个简洁的Go语言SDK，封装了 gojieba 库以简化中文分词的调用").Frequency()
for word, count := range wordFreq {
	fmt.Printf("%s: %d次\n", word, count)
}
```
</details>

### 🧰 辅助工具

<details>
<summary><b>📅 日期工具</b></summary>

```go
import "github.com/karosown/katool-go/util/dateutil"

// 日期工具
now := dateutil.Now()
formatted := dateutil.Format(now, "yyyy-MM-dd")
tomorrow := dateutil.AddDays(now, 1)
```
</details>

<details>
<summary><b>🎲 随机数工具</b></summary>

```go
import "github.com/karosown/katool-go/util/randutil"

// 随机数工具
randomInt := randutil.Int(1, 100)
randomString := randutil.String(10)
uuid := randutil.UUID()
```
</details>

<details>
<summary><b>📁 路径工具</b></summary>

```go
import "github.com/karosown/katool-go/util/pathutil"

// 路径工具
abs := pathutil.Abs("config.json")
joined := pathutil.Join("dir", "file.txt")
exists := pathutil.Exists("data.json")
```
</details>

<details>
<summary><b>🔍 调试工具</b></summary>

```go
import (
	"github.com/karosown/katool-go/util/dumper"
	"github.com/karosown/katool-go/sys"
)

// 调试工具
dumper.Dump(complexObject) // 打印对象结构

// 系统工具
sys.Warn("警告信息")
sys.Panic("发生严重错误") // 会导致panic
```
</details>

<hr>

## 💡 最佳实践

<details>
<summary><b>🌊 流式处理</b></summary>

- 对于大数据集，使用 `Parallel()` 开启并行处理
- 使用 `Reduce` 时注意提供合适的初始值和合并函数
- 在链式操作中，尽量将过滤操作放在前面，减少后续处理的数据量
- 根据实际测试案例，流式处理在处理大量数据时比传统循环更具可读性

```go
// ✅ 推荐写法：先过滤再转换
result := stream.ToStream(&users).
	Parallel().  // 启用并行处理
	Filter(func(u user) bool { 
		return u.Sex != 0 // 先过滤，减少数据量
	}).
	Map(func(u user) any {
		// 对过滤后的数据进行转换
		return u.Name
	}).
	ToList()

// ❌ 不推荐写法：先转换再过滤
result := stream.ToStream(&users).
	Map(func(u user) any {
		// 转换所有数据，包括最终会被过滤掉的
		return u.Name
	}).
	Filter(func(name any) bool {
		// 过滤转换后的数据，浪费了转换资源
		return someCondition
	}).
	ToList()
```
</details>

<details>
<summary><b>🔒 并发控制</b></summary>

- 使用 `Synchronized` 替代直接操作锁，减少忘记解锁的风险
- 注意协程泄漏，确保每个 `Park()` 都有对应的 `Unpark()`
- 推荐使用 `defer` 语句确保资源被正确释放
- 对于多个 `LockSupport` 的管理，可结合流式处理进行批量操作

```go
// ✅ 推荐写法：使用流式API管理多个LockSupport
supports := make([]*lock.LockSupport, n)
for i := 0; i < n; i++ {
	supports[i] = lock.NewLockSupport()
	// 启动工作协程...
}

// 批量唤醒所有协程
stream.ToStream(&supports).ForEach(func(ls *lock.LockSupport) {
	ls.Unpark()
})

// ✅ 推荐写法：使用defer确保Unpark
func someFunction() {
	ls := lock.NewLockSupport()
	done := false
	
	go func() {
		defer func() { done = true }()
		// 执行某些操作...
		ls.Park() // 阻塞等待信号
		// 继续操作...
	}()
	
	// 等待条件满足
	for !done {
		// 检查条件...
		if conditionMet {
			ls.Unpark() // 发送信号
			break
		}
		time.Sleep(checkInterval)
	}
}
```
</details>

<details>
<summary><b>🔤 文本处理</b></summary>

- 使用 `jieba` 分词时，记得使用 `defer` 确保调用 `Free()` 释放资源
- 根据不同需求选择合适的分词模式：
  - `Cut`: 精确模式，适合文本分析和提取关键信息
  - `CutAll`: 全模式，会把句子中所有可能的词都扫描出来
  - `CutForSearch`: 搜索引擎模式，在精确模式基础上对长词再次切分
- 使用 `Frequency()` 方法可以快速获取文本中的词频统计

```go
// ✅ 推荐写法：资源管理
func processText(text string) map[string]int {
	client := jieba.New()
	defer client.Free() // 确保资源释放
	
	// 根据需求选择合适的分词模式
	words := client.Cut(text)      // 一般场景
	// 或
	words = client.CutForSearch(text) // 搜索场景
	
	// 词频统计
	return words.Frequency()
}
```
</details>

<hr>

## 👥 贡献指南

我们欢迎各种形式的贡献，包括但不限于：

- 📝 报告问题和提出建议
- ✨ 提交修复和新功能
- 📚 改进文档和示例
- 🔧 优化性能和代码质量

### 贡献流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 代码要求

请确保代码符合以下要求：

- ✅ 通过所有测试
- 📏 遵循 Go 代码规范
- 📝 添加必要的文档和注释
- 🧪 包含适当的测试用例

<hr>

## 📄 许可证

Katool-Go 采用 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。

<div align="center">
  <sub>Made with ❤️ by Karosown Team</sub>
</div>