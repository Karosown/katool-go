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
  - [🧰 辅助工具](#辅助工具)
- [💡 最佳实践](#最佳实践)
- [👥 贡献指南](#贡献指南)
- [📄 许可证](#许可证)

<hr>

## 📝 简介

> **Katool-Go** 是一个综合性的 Go 语言工具库，旨在提供丰富的功能组件和实用工具，帮助开发者提高开发效率。它借鉴了 Java 生态中的成熟设计模式和经验，同时充分利用 Go 语言的特性，如并发、泛型等，提供了一系列易用且高效的工具。

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

> ⚠️ 要求 Go 版本 >= 1.18 (支持泛型)

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
)

func main() {
	// 准备数据
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	// 使用流式操作：过滤偶数、乘以2、求和
	result := stream.ToStream(&numbers).
		Filter(func(n int) bool {
			return n%2 == 0 // 过滤偶数
		}).
		Map(func(n int) any {
			return n * 2 // 乘以2
		}).
		Reduce(0, func(sum any, n int) any {
			return sum.(int) + n.(int) // 求和
		}, func(sum1, sum2 any) any {
			return sum1.(int) + sum2.(int)
		})
	
	fmt.Println("结果:", result) // 输出: 结果: 60
}
```
</details>

<details>
<summary><b>💉 依赖注入</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/container/ioc"
)

type UserService interface {
	GetUsername() string
}

type UserServiceImpl struct {
	username string
}

func (u *UserServiceImpl) GetUsername() string {
	return u.username
}

func main() {
	// 注册服务
	ioc.RegisterValue("userService", &UserServiceImpl{username: "admin"})
	
	// 获取服务
	service := ioc.Get("userService").(UserService)
	
	fmt.Println("用户名:", service.GetUsername()) // 输出: 用户名: admin
}
```
</details>

<details>
<summary><b>🔄 数据转换</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/convert"
)

type User struct {
	ID   int
	Name string
	Age  int
}

type UserDTO struct {
	ID   int
	Name string
	Age  int
}

func main() {
	user := User{ID: 1, Name: "Alice", Age: 30}
	
	// 属性复制
	userDTO, _ := convert.CopyProperties(user, &UserDTO{})
	
	fmt.Printf("原始用户: %+v\n", user)     // 输出: 原始用户: {ID:1 Name:Alice Age:30}
	fmt.Printf("转换后DTO: %+v\n", userDTO) // 输出: 转换后DTO: &{ID:1 Name:Alice Age:30}
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

// 创建XMap
m := xmap.New[string, int]()

// 设置值
m.Put("one", 1)
m.Put("two", 2)

// 获取值
val, exists := m.Get("one") // val=1, exists=true

// 遍历
m.ForEach(func(k string, v int) {
    fmt.Printf("%s: %d\n", k, v)
})

// 转换为流
stream := m.Stream()
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
import "github.com/karosown/katool-go/container/stream"

// 准备数据
users := []User{
    {ID: 1, Name: "Alice", Age: 25},
    {ID: 2, Name: "Bob", Age: 30},
    {ID: 3, Name: "Charlie", Age: 35},
    {ID: 4, Name: "David", Age: 40},
}

// 创建流
s := stream.ToStream(&users)

// 过滤操作
filtered := s.Filter(func(u User) bool {
    return u.Age > 30
}).ToList() // [{ID:3 Name:Charlie Age:35}, {ID:4 Name:David Age:40}]

// 映射操作
names := stream.ToStream(&users).
    Map(func(u User) any {
        return u.Name
    }).ToList() // ["Alice", "Bob", "Charlie", "David"]

// 排序操作
sorted := stream.ToStream(&users).
    Sort(func(a, b User) bool {
        return a.Age < b.Age // 按年龄升序
    }).ToList()
```
</details>

<details>
<summary><b>🚀 高级操作</b></summary>

```go
// 分组操作
groups := stream.ToStream(&users).
    GroupBy(func(u User) any {
        if u.Age < 30 {
            return "young"
        }
        return "senior"
    }) // map[young:[{ID:1 Name:Alice Age:25}] senior:[{ID:2 Name:Bob Age:30}, ...]]

// 并行处理
result := stream.ToStream(&users).
    Parallel(). // 启用并行处理
    Filter(func(u User) bool {
        return u.Age > 25
    }).
    Map(func(u User) any {
        // 模拟耗时操作
        time.Sleep(100 * time.Millisecond)
        return u.Name
    }).ToList()

// 扁平化映射
departments := []Department{
    {Name: "Engineering", Members: []User{{Name: "Alice"}, {Name: "Bob"}}},
    {Name: "Marketing", Members: []User{{Name: "Charlie"}, {Name: "David"}}},
}

allMembers := stream.ToStream(&departments).
    FlatMap(func(d Department) *stream.Stream[any, []any] {
        members := d.Members
        return stream.ToStream(&members).Map(func(u User) any {
            return u.Name
        })
    }).ToList() // ["Alice", "Bob", "Charlie", "David"]
```
</details>

<details>
<summary><b>📊 收集操作</b></summary>

```go
// 转换为列表
list := stream.ToStream(&users).ToList()

// 转换为映射
userMap := stream.ToStream(&users).
    ToMap(func(i int, u User) any {
        return u.ID // 键
    }, func(i int, u User) any {
        return u.Name // 值
    }) // map[1:"Alice" 2:"Bob" 3:"Charlie" 4:"David"]

// 汇总统计
sum := stream.ToStream(&users).
    Reduce(0, func(acc any, u User) any {
        return acc.(int) + u.Age
    }, func(a, b any) any {
        return a.(int) + b.(int)
    }).(int) // sum=130

// 自定义收集
result := stream.ToStream(&users).
    Collect(func(data stream.Options[User], srcData []User) any {
        // 自定义收集逻辑
        total := 0
        for _, opt := range data {
            total += opt.opt.Age
        }
        return total / len(data)
    }).(int) // 平均年龄
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

```go
import "github.com/karosown/katool-go/lock"

// 创建LockSupport
ls := lock.NewLockSupport()

// 在协程中使用
go func() {
    fmt.Println("协程开始")
    ls.Park() // 阻塞协程
    fmt.Println("协程继续")
}()

// 等待一段时间
time.Sleep(time.Second)

// 恢复协程
ls.Unpark() // 解除阻塞
```
</details>

<details>
<summary><b>🔐 同步工具</b></summary>

```go
// 同步代码块
mutex := &sync.Mutex{}
counter := 0

lock.Synchronized(mutex, func() {
    // 临界区代码
    counter++
})

// 锁映射
lockMap := lock.LockMap{}
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
import "github.com/karosown/katool-go/log"

// 基本日志
log.Info("这是一条信息")
log.Errorf("错误: %v", err)
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
```
</details>

<details>
<summary><b>🔐 哈希计算</b></summary>

```go
// 哈希计算
data := map[string]any{"id": 123, "name": "test"}
hash := algorithm.HASH_WITH_JSON(data) // 使用JSON序列化计算哈希
md5Hash := algorithm.HASH_WITH_JSON_MD5(data) // 使用MD5计算哈希
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

```go
// ✅ 推荐写法
result := stream.ToStream(&largeDataset).
    Parallel().                // 启用并行处理
    Filter(func(d Data) bool { // 先过滤，减少数据量
        return d.IsValid
    }).
    Map(func(d Data) any {     // 然后转换
        return d.Value
    }).
    ToList()

// ❌ 不推荐写法
result := stream.ToStream(&largeDataset).
    Map(func(d Data) any {     // 先转换所有数据
        return d.Value
    }).
    Filter(func(d Data) bool { // 再过滤
        return d.IsValid
    }).
    ToList()
```
</details>

<details>
<summary><b>💉 依赖注入</b></summary>

- 优先注册接口而非具体实现
- 使用工厂方法注册有依赖关系的组件
- 注意避免循环依赖

```go
// ✅ 推荐写法
ioc.RegisterValue("userRepo", &UserRepositoryImpl{})
ioc.Register("userService", func() any {
    repo := ioc.Get("userRepo").(UserRepository)
    return &UserServiceImpl{Repository: repo}
})

// ❌ 不推荐写法: 硬编码依赖
ioc.RegisterValue("userService", &UserServiceImpl{
    Repository: &UserRepositoryImpl{},
})
```
</details>

<details>
<summary><b>🔒 并发控制</b></summary>

- 使用 `Synchronized` 替代直接操作锁，减少忘记解锁的风险
- 注意协程泄漏，确保每个 `Park()` 都有对应的 `Unpark()`
- 在高并发场景下，考虑使用 `LockMap` 减少锁冲突

```go
// ✅ 推荐写法
mutex := &sync.Mutex{}
lock.Synchronized(mutex, func() {
    // 临界区代码
})

// ❌ 不推荐写法: 容易忘记解锁
mutex.Lock()
// 临界区代码
mutex.Unlock()
```
</details>

<details>
<summary><b>🔄 数据转换</b></summary>

- 使用 `CopyProperties` 时注意字段类型和名称匹配
- 对于复杂对象，考虑实现自定义转换逻辑
- 使用泛型版本的 `Convert` 函数处理批量转换

```go
// ✅ 推荐写法: 使用泛型Convert
dtos := convert.Convert(users, func(u User) UserDTO {
    return UserDTO{
        ID:   convert.ToString(u.ID),
        Name: u.Name,
    }
})

// ❌ 不推荐写法: 手动循环转换
dtos := make([]UserDTO, len(users))
for i, u := range users {
    dtos[i] = UserDTO{
        ID:   convert.ToString(u.ID),
        Name: u.Name,
    }
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