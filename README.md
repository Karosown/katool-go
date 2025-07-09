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

<b><i>一个功能丰富的 Go 工具库，借鉴 Java 生态优秀设计，为 Go 开发提供全方位支持</i></b>

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
  - [⚡ 规则引擎](#规则引擎)
  - [🧰 辅助工具](#辅助工具)
- [💡 最佳实践](#最佳实践)
- [👥 贡献指南](#贡献指南)
- [📄 许可证](#许可证)

<hr>

## 📝 简介

**Katool-Go** 是一个现代化的 Go 语言综合工具库，专为提高开发效率和代码质量而设计。它借鉴了 Java 生态系统中的成熟设计模式，同时充分利用 Go 语言的现代特性，如泛型、协程等，为开发者提供了一套完整的工具解决方案。

本库采用**模块化、类型安全、高性能**的设计理念，适用于各种规模的 Go 项目，从微服务到大型企业应用，都能提供强有力的支持。

### 🎯 设计目标

- **类型安全**：充分利用 Go 1.18+ 泛型特性，提供类型安全的 API
- **性能优异**：内置并发优化，充分发挥 Go 语言性能优势
- **易于使用**：提供类似 Java Stream API 的链式操作，降低学习成本
- **生产就绪**：完整的错误处理、日志系统和测试覆盖

<hr>

## ✨ 特性

Katool-Go 提供以下核心特性：

<table>
  <tr>
    <td><b>🌊 流式处理</b></td>
    <td>类似 Java 8 Stream API 的链式操作，支持并行处理、map/filter/reduce/collect 等完整操作集</td>
  </tr>
  <tr>
    <td><b>📚 容器与集合</b></td>
    <td>增强的集合类型：Map、SafeMap、SortedMap、HashBasedMap、Optional 等，全部支持泛型</td>
  </tr>
  <tr>
    <td><b>💉 依赖注入</b></td>
    <td>轻量级 IOC 容器，支持组件注册、获取和生命周期管理，简化依赖管理</td>
  </tr>
  <tr>
    <td><b>🔒 并发控制</b></td>
    <td>LockSupport（类似Java的park/unpark）、同步锁封装等协程控制工具</td>
  </tr>
  <tr>
    <td><b>🔄 数据转换</b></td>
    <td>结构体属性复制、类型转换、文件导出（CSV/JSON）、序列化等全方位数据处理</td>
  </tr>
  <tr>
    <td><b>🕸️ Web爬虫</b></td>
    <td>智能内容提取、Chrome渲染支持、RSS订阅解析等完整爬虫解决方案</td>
  </tr>
  <tr>
    <td><b>📁 文件操作</b></td>
    <td>文件下载器、序列化工具、路径处理等文件系统操作</td>
  </tr>
  <tr>
    <td><b>💾 数据库支持</b></td>
    <td>MongoDB 增强工具、分页查询器等数据库操作简化</td>
  </tr>
  <tr>
    <td><b>🌐 网络通信</b></td>
    <td>现代化 HTTP 客户端、OAuth2 支持、SSE 实时通信、RESTful API 封装</td>
  </tr>
  <tr>
    <td><b>📝 日志系统</b></td>
    <td>结构化日志、链式构建、多级别输出、自定义格式化等完整日志方案</td>
  </tr>
  <tr>
    <td><b>⚙️ 算法工具</b></td>
    <td>有序数组合并、多种哈希计算、数据结构算法等实用算法集</td>
  </tr>
  <tr>
    <td><b>🔤 文本处理</b></td>
    <td>中文分词（jieba）、词频统计、文本分析、语言检测等NLP工具</td>
  </tr>
  <tr>
    <td><b>⚡ 规则引擎</b></td>
    <td>灵活的业务规则处理、规则链构建、中间件支持等企业级规则管理</td>
  </tr>
  <tr>
    <td><b>🧰 辅助工具</b></td>
    <td>日期处理、随机数生成、调试工具、系统工具等开发辅助功能</td>
  </tr>
</table>

<hr>

## 📦 安装

使用 `go get` 安装最新版本：

```bash
go get -u github.com/karosown/katool-go
```

> ⚠️ **系统要求**
> - Go 版本 >= 1.23.1
> - 支持泛型特性
> - 推荐使用最新版本以获得最佳性能

<hr>

## 🚀 快速开始

下面是几个核心功能的快速示例，展示 Katool-Go 的强大能力：

<details open>
<summary><b>🌊 流式处理 - Java风格的链式操作</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/algorithm"
)

// 定义用户结构
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Sex   int    `json:"sex"`    // 0-女性，1-男性
	Money int    `json:"money"`
	Class string `json:"class"`
	Id    int    `json:"id"`
}

func main() {
	users := []User{
		{Name: "Alice", Age: 25, Sex: 1, Money: 1000, Class: "A", Id: 1},
		{Name: "Bob", Age: 30, Sex: 0, Money: 1500, Class: "B", Id: 2},
		{Name: "Charlie", Age: 35, Sex: 0, Money: 2000, Class: "A", Id: 3},
		{Name: "David", Age: 40, Sex: 1, Money: 2500, Class: "B", Id: 4},
	}
	
	// 创建并行流
	userStream := stream.ToStream(&users).Parallel()
	
	// 链式操作：过滤 -> 排序 -> 统计
	adultUsers := userStream.
		Filter(func(u User) bool { 
			return u.Age >= 30 
		}).
		Sort(func(a, b User) bool { 
			return a.Money > b.Money  // 按收入降序
		}).
		ToList()
	
	fmt.Printf("30岁以上用户（按收入排序）: %+v\n", adultUsers)
	
	// 聚合计算：总收入
	totalMoney := userStream.Reduce(int64(0), 
		func(sum any, u User) any { 
			return sum.(int64) + int64(u.Money) 
		}, 
		func(sum1, sum2 any) any {
			return sum1.(int64) + sum2.(int64)
		})
	fmt.Printf("总收入: %d\n", totalMoney)
	
	// 分组统计：按班级分组
	groups := stream.ToStream(&users).GroupBy(func(u User) any {
		return u.Class
	})
	
	for class, members := range groups {
		fmt.Printf("班级 %s: %d人\n", class, len(members))
	}
	
	// 去重操作（基于JSON序列化）
	uniqueUsers := userStream.DistinctBy(algorithm.HASH_WITH_JSON_MD5).ToList()
	fmt.Printf("去重后用户数: %d\n", len(uniqueUsers))
	
	// 转换为映射
	userMap := stream.ToStream(&users).ToMap(
		func(index int, u User) any { return u.Id },
		func(index int, u User) any { return u.Name },
	)
	fmt.Printf("用户ID->姓名映射: %+v\n", userMap)
}
```
</details>

<details>
<summary><b>📚 增强集合 - 类型安全的容器</b></summary>

```go
package main

import (
	"fmt"
	"encoding/json"
	"github.com/karosown/katool-go/container/xmap"
	"github.com/karosown/katool-go/container/optional"
)

func main() {
	// 1. 基础Map - 泛型支持
	userMap := xmap.NewMap[string, User]()
	userMap.Set("alice", User{Name: "Alice", Age: 25})
	userMap.Set("bob", User{Name: "Bob", Age: 30})
	
	if user, exists := userMap.Get("alice"); exists {
		fmt.Printf("找到用户: %+v\n", user)
	}
	
	// 2. 线程安全Map - 并发场景
	safeMap := xmap.NewSafeMap[string, int]()
	
	// 原子操作：获取或存储
	value, loaded := safeMap.LoadOrStore("counter", 1)
	fmt.Printf("计数器值: %d, 是否已存在: %v\n", value, loaded)
	
	// 原子操作：获取并删除
	value, exists := safeMap.LoadAndDelete("counter")
	fmt.Printf("删除的值: %d, 是否存在: %v\n", value, exists)
	
	// 3. 有序Map - 按键排序，支持JSON序列化
	sortedMap := xmap.NewSortedMap[string, string]()
	sortedMap.Set("3", "three")
	sortedMap.Set("1", "one")
	sortedMap.Set("2", "two")
	
	jsonBytes, _ := json.Marshal(sortedMap)
	fmt.Printf("有序JSON: %s\n", string(jsonBytes))  // 按键排序输出
	
	// 4. 双层键映射
	dbMap := xmap.NewHashBasedMap[string, int, User]()
	dbMap.Set("users", 1, User{Name: "Alice", Age: 25})
	dbMap.Set("users", 2, User{Name: "Bob", Age: 30})
	dbMap.Set("admins", 1, User{Name: "Admin", Age: 40})
	
	if user, exists := dbMap.Get("users", 1); exists {
		fmt.Printf("用户表中ID=1的用户: %+v\n", user)
	}
	
	// 5. Optional - 避免空指针
	opt := optional.Of("Hello World")
	opt.IfPresent(func(s string) {
		fmt.Printf("Optional值: %s\n", s)
	})
	
	emptyOpt := optional.Empty[string]()
	defaultValue := emptyOpt.OrElse("默认值")
	fmt.Printf("空Optional的默认值: %s\n", defaultValue)
}
```
</details>

<details>
<summary><b>🔒 并发控制 - 协程同步</b></summary>

```go
package main

import (
	"fmt"
	"time"
	"sync"
	"github.com/karosown/katool-go/lock"
	"github.com/karosown/katool-go/container/stream"
)

func main() {
	// 1. LockSupport - 类似Java的park/unpark
	fmt.Println("=== LockSupport 示例 ===")
	ls := lock.NewLockSupport()
	
	go func() {
		fmt.Println("子协程：准备阻塞等待...")
		ls.Park()  // 阻塞直到被唤醒
		fmt.Println("子协程：被成功唤醒！")
	}()
	
	time.Sleep(time.Second)
	fmt.Println("主协程：发送唤醒信号")
	ls.Unpark()  // 唤醒阻塞的协程
	
	time.Sleep(100 * time.Millisecond)  // 等待输出
	
	// 2. 批量协程管理
	fmt.Println("\n=== 批量协程管理 ===")
	supports := make([]*lock.LockSupport, 5)
	for i := 0; i < 5; i++ {
		supports[i] = lock.NewLockSupport()
		idx := i
		go func() {
			fmt.Printf("协程 %d: 等待唤醒\n", idx)
			supports[idx].Park()
			fmt.Printf("协程 %d: 被唤醒\n", idx)
		}()
	}
	
	time.Sleep(500 * time.Millisecond)
	
	// 使用流式API批量唤醒
	stream.ToStream(&supports).ForEach(func(ls *lock.LockSupport) {
		ls.Unpark()
		time.Sleep(100 * time.Millisecond)  // 依次唤醒
	})
	
	// 3. 同步代码块
	fmt.Println("\n=== 同步代码块 ===")
	var counter int
	var mutex sync.Mutex
	
	// 传统方式 vs 封装方式
	lock.Synchronized(&mutex, func() {
		counter++
		fmt.Printf("同步块中的计数器: %d\n", counter)
	})
	
	// 带返回值的同步
	result := lock.SynchronizedT(&mutex, func() string {
		return fmt.Sprintf("最终计数: %d", counter)
	})
	fmt.Println(result)
	
	time.Sleep(100 * time.Millisecond)
}
```
</details>

<details>
<summary><b>🔄 数据转换 - 全方位数据处理</b></summary>

```go
package main

import (
	"fmt"
	"github.com/karosown/katool-go/convert"
)

// 源结构体
type UserEntity struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	CreateAt string `json:"create_at"`
}

// 目标DTO
type UserDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Status   string `json:"status"`  // 新增字段
}

func main() {
	users := []UserEntity{
		{ID: 1, Name: "Alice", Age: 25, Email: "alice@example.com", CreateAt: "2024-01-01"},
		{ID: 2, Name: "Bob", Age: 30, Email: "bob@example.com", CreateAt: "2024-01-02"},
		{ID: 3, Name: "Charlie", Age: 35, Email: "charlie@example.com", CreateAt: "2024-01-03"},
	}
	
	// 1. 属性复制（同名字段自动复制）
	fmt.Println("=== 属性复制 ===")
	sourceUser := users[0]
	targetDTO := &UserDTO{Status: "Active"}  // 预设新字段
	
	result, err := convert.CopyProperties(sourceUser, targetDTO)
	if err == nil {
		fmt.Printf("复制结果: %+v\n", result)
	}
	
	// 2. 批量转换
	fmt.Println("\n=== 批量转换 ===")
	dtos := convert.Convert(users, func(user UserEntity) UserDTO {
		return UserDTO{
			ID:     user.ID,
			Name:   user.Name,
			Age:    user.Age,
			Email:  user.Email,
			Status: "Active",
		}
	})
	fmt.Printf("转换后的DTO列表: %+v\n", dtos)
	
	// 3. 类型转换
	fmt.Println("\n=== 类型转换 ===")
	fmt.Printf("整数转字符串: %s\n", convert.ToString(123))
	fmt.Printf("布尔转字符串: %s\n", convert.ToString(true))
	fmt.Printf("切片转字符串: %s\n", convert.ToString([]int{1, 2, 3}))
	
	// 4. 类型擦除和恢复
	fmt.Println("\n=== 类型擦除和恢复 ===")
	anySlice := convert.ToAnySlice(users)
	fmt.Printf("类型擦除后长度: %d\n", len(anySlice))
	
	recoveredUsers := convert.FromAnySlice[UserEntity](anySlice)
	fmt.Printf("恢复类型后第一个用户: %+v\n", recoveredUsers[0])
	
	// 5. 文件导出
	fmt.Println("\n=== 文件导出 ===")
	// 导出为JSON文件
	err = convert.StructToJsonFile(users, "users.json")
	if err == nil {
		fmt.Println("成功导出JSON文件: users.json")
	}
	
	// 导出为CSV文件（需要csv标签）
	type UserCSV struct {
		ID   int    `csv:"用户ID"`
		Name string `csv:"姓名"`
		Age  int    `csv:"年龄"`
	}
	
	csvUsers := convert.Convert(users, func(u UserEntity) UserCSV {
		return UserCSV{ID: u.ID, Name: u.Name, Age: u.Age}
	})
	
	err = convert.StructToCSV(csvUsers, "users.csv")
	if err == nil {
		fmt.Println("成功导出CSV文件: users.csv")
	}
}
```
</details>

<hr>

## 🔧 核心模块

### 📚 容器与集合

<details>
<summary><b>🗂️ XMap - 增强的映射类型</b></summary>

XMap 提供了比标准 map 更丰富的功能和类型安全保证：

```go
import "github.com/karosown/katool-go/container/xmap"

// 1. 基础Map - 泛型支持
regularMap := xmap.NewMap[string, int]()
regularMap.Set("one", 1)
regularMap.Set("two", 2)

// 2. 线程安全Map - 并发安全
safeMap := xmap.NewSafeMap[string, int]()
safeMap.Set("counter", 1)

// 原子操作
value, loaded := safeMap.LoadOrStore("new_key", 100)  // 不存在则存储
value, exists := safeMap.LoadAndDelete("counter")     // 获取并删除

// 3. 有序Map - 按键排序
sortedMap := xmap.NewSortedMap[string, string]()
sortedMap.Set("c", "third")
sortedMap.Set("a", "first")
sortedMap.Set("b", "second")

// JSON序列化自动按键排序
jsonBytes, _ := json.Marshal(sortedMap)  // {"a":"first","b":"second","c":"third"}

// 4. 双层键映射
hashMap := xmap.NewHashBasedMap[string, int, User]()
hashMap.Set("users", 1, User{Name: "Alice"})
hashMap.Set("users", 2, User{Name: "Bob"})
hashMap.Set("admins", 1, User{Name: "Admin"})

user, exists := hashMap.Get("users", 1)  // 通过两个键定位
```
</details>

<details>
<summary><b>📦 Optional - 空值安全处理</b></summary>

Optional 提供了处理可能为空值的安全方式，避免空指针异常：

```go
import "github.com/karosown/katool-go/container/optional"

// 创建Optional
opt := optional.Of("Hello World")
emptyOpt := optional.Empty[string]()

// 安全检查和获取
if opt.IsPresent() {
    value := opt.Get()
    fmt.Println("值存在:", value)
}

// 提供默认值
defaultValue := emptyOpt.OrElse("默认值")

// 条件执行
opt.IfPresent(func(v string) {
    fmt.Println("处理值:", v)
})

// 链式操作
result := optional.Of("  hello  ").
    Map(strings.TrimSpace).
    Filter(func(s string) bool { return len(s) > 0 }).
    OrElse("空字符串")

// 工具函数
enabled := optional.IsTrue(condition, "启用", "禁用")
result := optional.IsTrueByFunc(condition, enabledFunc, disabledFunc)
```
</details>

### 🌊 流式处理

<details>
<summary><b>🔄 Stream API - Java风格的链式操作</b></summary>

```go
import "github.com/karosown/katool-go/container/stream"

users := []User{
    {Name: "Alice", Age: 25, Salary: 5000},
    {Name: "Bob", Age: 30, Salary: 6000},
    {Name: "Charlie", Age: 35, Salary: 7000},
}

// 1. 基本操作链
result := stream.ToStream(&users).
    Parallel().                           // 启用并行处理
    Filter(func(u User) bool {           // 过滤
        return u.Age >= 30
    }).
    Map(func(u User) any {               // 转换
        return u.Name
    }).
    Sort(func(a, b any) bool {           // 排序
        return a.(string) < b.(string)
    }).
    ToList()

// 2. 聚合操作
totalSalary := stream.ToStream(&users).
    Reduce(0, func(sum any, u User) any {
        return sum.(int) + u.Salary
    }, func(sum1, sum2 any) any {
        return sum1.(int) + sum2.(int)
    })

// 3. 分组操作
ageGroups := stream.ToStream(&users).GroupBy(func(u User) any {
    if u.Age < 30 {
        return "young"
    }
    return "senior"
})

// 4. 去重操作
uniqueUsers := stream.ToStream(&users).
    DistinctBy(algorithm.HASH_WITH_JSON_MD5).
    ToList()

// 5. 扁平化处理
departments := []Department{
    {Name: "IT", Users: []User{{Name: "Alice"}, {Name: "Bob"}}},
    {Name: "HR", Users: []User{{Name: "Charlie"}}},
}

allUsers := stream.ToStream(&departments).
    FlatMap(func(dept Department) *stream.Stream[any, []any] {
        userAnySlice := convert.ToAnySlice(dept.Users)
        return stream.ToStream(&userAnySlice)
    }).
    ToList()

// 6. 转换为Map
userMap := stream.ToStream(&users).ToMap(
    func(index int, u User) any { return u.ID },
    func(index int, u User) any { return u.Name },
)

// 7. 统计操作
count := stream.ToStream(&users).Count()
seniorCount := stream.ToStream(&users).
    Filter(func(u User) bool { return u.Age >= 35 }).
    Count()

// 8. 集合操作
newUsers := []User{{Name: "David", Age: 28}}
mergedStream := stream.ToStream(&users).Merge(newUsers)

intersection := stream.ToStream(&users).
    Intersect(newUsers, func(a, b User) bool {
        return a.Name == b.Name
    })
```
</details>

### 💉 依赖注入

<details>
<summary><b>🏭 IOC容器 - 轻量级依赖管理</b></summary>

```go
import "github.com/karosown/katool-go/container/ioc"

// 定义接口和实现
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

type DatabaseUserRepository struct {
    connectionString string
}

func (r *DatabaseUserRepository) FindByID(id int) (*User, error) {
    // 数据库查询逻辑
    return &User{ID: id, Name: "User" + strconv.Itoa(id)}, nil
}

func (r *DatabaseUserRepository) Save(user *User) error {
    // 保存逻辑
    return nil
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.FindByID(id)
}

func main() {
    // 1. 注册值对象
    ioc.RegisterValue("dbConnection", "localhost:5432/mydb")
    
    // 2. 注册工厂函数
    ioc.Register("userRepo", func() any {
        connStr := ioc.Get("dbConnection").(string)
        return &DatabaseUserRepository{connectionString: connStr}
    })
    
    // 3. 注册依赖其他组件的服务
    ioc.Register("userService", func() any {
        repo := ioc.Get("userRepo").(UserRepository)
        return &UserService{repo: repo}
    })
    
    // 4. 获取服务使用
    service := ioc.Get("userService").(*UserService)
    user, err := service.GetUser(1)
    if err == nil {
        fmt.Printf("获取到用户: %+v\n", user)
    }
    
    // 5. 获取带默认值的组件
    cache := ioc.GetDef("cache", &MemoryCache{})
    
    // 6. 强制注册（覆盖已存在的）
    ioc.ForceRegister("userRepo", func() UserRepository {
        return &MockUserRepository{}  // 测试时替换为Mock
    })
    
    // 7. 延迟注册（通过函数）
    config := ioc.GetDefFunc("config", func() *Config {
        return &Config{Debug: true, Port: 8080}
    })
}
```
</details>

### 🕸️ Web爬虫

<details>
<summary><b>📄 智能内容提取</b></summary>

```go
import "github.com/karosown/katool-go/web_crawler"

// 1. 基础内容提取
article := web_crawler.GetArticleWithURL("https://example.com/article")
if !article.IsErr() {
    fmt.Println("标题:", article.Title)
    fmt.Println("内容:", article.Content)
    fmt.Println("摘要:", article.Excerpt)
    fmt.Println("作者:", article.Byline)
    fmt.Println("发布时间:", article.PublishedTime)
}

// 2. 自定义请求头
article = web_crawler.GetArticleWithURL("https://example.com/article",
    func(r *http.Request) {
        r.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyBot/1.0)")
        r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
    })

// 3. Chrome渲染支持（处理JavaScript渲染的页面）
article = web_crawler.GetArticleWithChrome(
    "https://spa-example.com/article",
    func(page *rod.Page) {
        // 等待页面加载完成
        page.WaitLoad()
        // 等待特定元素出现
        page.MustElement(".article-content").WaitVisible()
        // 模拟用户交互
        page.MustElement("#load-more").Click()
        time.Sleep(2 * time.Second)
    },
    func(article web_crawler.Article) bool {
        // 重试条件：内容为空时重启Chrome
        return len(article.Content) == 0
    },
)

// 4. 源码提取
sourceCode := web_crawler.ReadSourceCode(
    "https://example.com",
    "",  // CSS选择器，空表示全页面
    func(page *rod.Page) {
        page.WaitLoad()
    },
)

if !sourceCode.IsErr() {
    fmt.Println("页面源码长度:", len(sourceCode.String()))
}

// 5. 路径解析工具
absoluteURL := web_crawler.ParsePath("https://example.com/page", "./image.jpg")
// 结果: "https://example.com/page/image.jpg"

relativeURL := web_crawler.ParsePath("https://example.com/page", "/api/data")
// 结果: "https://example.com/api/data"
```
</details>

<details>
<summary><b>📰 RSS订阅解析</b></summary>

```go
import "github.com/karosown/katool-go/web_crawler/rss"

// 解析RSS源
feed, err := rss.ParseURL("https://example.com/feed.xml")
if err == nil {
    fmt.Println("Feed标题:", feed.Title)
    fmt.Println("Feed描述:", feed.Description)
    fmt.Println("更新时间:", feed.LastBuildDate)
    
    // 遍历文章
    for _, item := range feed.Items {
        fmt.Printf("文章: %s\n", item.Title)
        fmt.Printf("链接: %s\n", item.Link)
        fmt.Printf("描述: %s\n", item.Description)
        fmt.Printf("发布时间: %s\n", item.PubDate)
        fmt.Printf("作者: %s\n", item.Author)
        fmt.Println("---")
    }
}

// 使用流式处理RSS数据
import "github.com/karosown/katool-go/container/stream"

recentArticles := stream.ToStream(&feed.Items).
    Filter(func(item rss.Item) bool {
        // 过滤最近一周的文章
        pubDate, _ := time.Parse(time.RFC1123, item.PubDate)
        return time.Since(pubDate) <= 7*24*time.Hour
    }).
    Sort(func(a, b rss.Item) bool {
        // 按发布时间降序排序
        dateA, _ := time.Parse(time.RFC1123, a.PubDate)
        dateB, _ := time.Parse(time.RFC1123, b.PubDate)
        return dateA.After(dateB)
    }).
    ToList()
```
</details>

### 🌐 网络通信

<details>
<summary><b>🔗 现代化HTTP客户端</b></summary>

```go
import "github.com/karosown/katool-go/net/http"

// 1. 基础HTTP请求
client := remote.NewRemoteRequest("https://api.example.com")

// GET请求
var users []User
resp, err := client.
    QueryParam(map[string]string{
        "page":     "1",
        "pageSize": "10",
    }).
    Headers(map[string]string{
        "Authorization": "Bearer your-token",
        "Content-Type":  "application/json",
    }).
    Method("GET").
    Url("/users").
    Build(&users)

// POST请求
newUser := User{Name: "Alice", Age: 25}
var createdUser User
resp, err = client.
    Data(newUser).
    Method("POST").
    Url("/users").
    Build(&createdUser)

// 2. 链式构建复杂请求
response, err := client.
    Url("/api/data").
    QueryParam(map[string]string{"filter": "active"}).
    Headers(map[string]string{"X-API-Version": "v2"}).
    FormData(map[string]string{
        "name":  "test",
        "value": "123",
    }).
    Files(map[string]string{
        "upload": "/path/to/file.txt",
    }).
    Method("POST").
    DecodeHandler(format.Json).  // 自定义解码器
    Build(&result)

// 3. 自定义HTTP客户端
customClient := resty.New().
    SetTimeout(30 * time.Second).
    SetRetryCount(3)

resp, err = client.
    HttpClient(customClient).
    Method("GET").
    Url("/api/retry-endpoint").
    Build(&result)
```
</details>

<details>
<summary><b>🔐 OAuth2 支持</b></summary>

```go
// OAuth2认证客户端
oauth := remote.NewOAuth2Request(
    "https://api.example.com",     // API基础URL
    "your-client-id",              // 客户端ID
    "your-client-secret",          // 客户端密钥
    "https://auth.example.com/token", // Token端点
)

// 自动处理Token获取和刷新
var protectedData ApiResponse
resp, err := oauth.
    Headers(map[string]string{"X-API-Version": "v1"}).
    Method("GET").
    Url("/protected-resource").
    Build(&protectedData)

// Token会自动管理，无需手动处理
```
</details>

<details>
<summary><b>📡 SSE 实时通信</b></summary>

```go
// SSE客户端
sseClient := remote.NewSSERequest("https://api.example.com")

// 连接SSE流
err := sseClient.
    Headers(map[string]string{"Authorization": "Bearer token"}).
    Connect("/events", func(event remote.SSEEvent) {
        fmt.Printf("收到事件: %s\n", event.Type)
        fmt.Printf("数据: %s\n", event.Data)
        fmt.Printf("ID: %s\n", event.ID)
    })

// 处理连接错误
if err != nil {
    fmt.Printf("SSE连接失败: %v\n", err)
}
```
</details>

### 📝 日志系统

<details>
<summary><b>📊 结构化日志</b></summary>

```go
import "github.com/karosown/katool-go/xlog"

// 1. 基础日志使用
xlog.Info("应用启动成功")
xlog.Errorf("处理失败: %v", err)
xlog.Debug("调试信息: 变量值为 %d", value)

// 2. 链式日志构建
logger := xlog.NewLogWrapper().
    Header("MyApplication").              // 应用标识
    FunctionByFunc(func(layer int) string {  // 自动获取函数名
        pc, _, _, _ := runtime.Caller(layer)
        return runtime.FuncForPC(pc).Name()
    }).
    ApplicationDesc("用户服务模块")         // 模块描述

// 不同级别的日志
logger.Info().ApplicationDesc("用户登录成功").String()
logger.Warn().ApplicationDesc("内存使用率过高").String()
logger.Error().ApplicationDesc("数据库连接失败").Panic()  // 会触发panic

// 3. 自定义格式化
customLogger := xlog.NewLogWrapper().
    Header("CustomApp").
    Format(func(msg xlog.LogMessage) string {
        return fmt.Sprintf("[%s] %s: %v", 
            msg.Header, msg.Type, msg.ApplicationDesc)
    }).
    Info()

// 4. 内置工具日志器
xlog.KaToolLoggerWrapper.ApplicationDesc("工具库内部错误").Error()

// 5. 自定义Logger配置
logger := xlog.NewLogger(
    xlog.WithLevel(xlog.InfoLevel),
    xlog.WithFormat(xlog.JSONFormat),
    xlog.WithOutput("app.log"),
    xlog.WithRotation(xlog.DailyRotation),
)

logger.WithFields(xlog.Fields{
    "userID": 12345,
    "action": "login",
    "ip":     "192.168.1.1",
}).Info("用户操作记录")
```
</details>

### ⚙️ 算法工具

<details>
<summary><b>🔢 数组和哈希算法</b></summary>

```go
import "github.com/karosown/katool-go/algorithm"

// 1. 有序数组合并
arr1 := []int{1, 3, 5, 7}
arr2 := []int{2, 4, 6, 8}

// 自定义比较函数的合并
mergeFunc := algorithm.MergeSortedArrayWithEntity[int](func(a, b int) bool {
    return a < b  // 升序
})
merged := mergeFunc(convert.ToAnySlice(arr1), convert.ToAnySlice(arr2))
// 结果: [1, 2, 3, 4, 5, 6, 7, 8]

// 基于哈希值的合并（用于复杂对象）
users1 := []User{{ID: 1, Name: "Alice"}, {ID: 3, Name: "Charlie"}}
users2 := []User{{ID: 2, Name: "Bob"}, {ID: 4, Name: "David"}}

userMergeFunc := algorithm.MergeSortedArrayWithPrimaryData[User](
    false,  // 升序
    func(user any) algorithm.HashType {
        u := user.(User)
        return algorithm.HashType(fmt.Sprintf("%d", u.ID))
    },
)
mergedUsers := userMergeFunc(convert.ToAnySlice(users1), convert.ToAnySlice(users2))

// 基于ID的合并
idMergeFunc := algorithm.MergeSortedArrayWithPrimaryId[User](
    false,  // 升序
    func(user any) algorithm.IDType {
        return algorithm.IDType(user.(User).ID)
    },
)

// 2. 哈希计算
data := map[string]any{
    "id":   123,
    "name": "test",
    "tags": []string{"go", "tool"},
}

// 基于JSON序列化的哈希
jsonHash := algorithm.HASH_WITH_JSON(data)
fmt.Printf("JSON Hash: %s\n", jsonHash)

// MD5哈希
md5Hash := algorithm.HASH_WITH_JSON_MD5(data)
fmt.Printf("MD5 Hash: %s\n", md5Hash)

// 简单累加哈希（性能更好）
sumHash := algorithm.HASH_WITH_JSON_SUM(data)
fmt.Printf("Sum Hash: %s\n", sumHash)

// 3. 在流式处理中使用
uniqueData := stream.ToStream(&dataList).
    DistinctBy(algorithm.HASH_WITH_JSON_MD5).  // 使用MD5去重
    ToList()

// 4. 二进制工具
binary := algorithm.ToBinary(42)        // 转二进制
decimal := algorithm.FromBinary("101010") // 从二进制转回
```
</details>

### 🔤 文本处理

<details>
<summary><b>📝 中文分词和文本分析</b></summary>

```go
import "github.com/karosown/katool-go/words/split/jieba"

// 1. 基础分词
jb := jieba.New()
defer jb.Free()  // 必须释放资源

text := "我正在测试Katool-Go的中文分词功能，效果很不错！"

// 精确模式分词（推荐）
words := jb.Cut(text)
fmt.Printf("精确分词: %v\n", words)
// 输出: ["我", "正在", "测试", "Katool-Go", "的", "中文", "分词", "功能", "效果", "很", "不错"]

// 全模式分词（包含所有可能的词）
allWords := jb.CutAll(text)
fmt.Printf("全模式分词: %v\n", allWords)

// 搜索引擎模式（适合搜索索引）
searchWords := jb.CutForSearch("清华大学计算机科学与技术系")
fmt.Printf("搜索模式: %v\n", searchWords)
// 输出: ["清华", "华大", "大学", "清华大学", "计算", "计算机", "科学", "技术", "系"]

// 2. 词频统计
document := "机器学习是人工智能的一个重要分支。机器学习算法能够从数据中学习模式。"
words = jb.Cut(document)

// 获取词频统计
frequency := words.Frequency()
frequency.Range(func(word string, count int64) bool {
    fmt.Printf("词: %s, 频次: %d\n", word, count)
    return true
})

// 3. 流式处理分词结果
meaningfulWords := words.ToStream().
    Filter(func(word string) bool {
        // 过滤停用词和标点
        return len(word) > 1 && !words.IsStopWord(word)
    }).
    Distinct().  // 去重
    Sort(func(a, b string) bool {
        return len(a) > len(b)  // 按长度排序
    }).
    ToList()

// 4. 文本工具函数
import "github.com/karosown/katool-go/words"

// 字符串截取
content := words.SubString("Hello [World] End", "[", "]")  // "World"

// 语言检测
hasChinese := words.ContainsLanguage("Hello世界", unicode.Han)  // true
onlyChinese := words.OnlyLanguage("世界", unicode.Han)         // true

// 大小写转换
shifted := words.CaseShift("Hello")  // "hELLO"

// 5. 自定义分词器
customJieba := jieba.New("/path/to/custom/dict.txt")
defer customJieba.Free()

customWords := customJieba.Cut("自定义词典测试")
```
</details>

### 🧰 辅助工具

<details>
<summary><b>📅 日期工具</b></summary>

```go
import "github.com/karosown/katool-go/util/dateutil"

// 性能测试
duration := dateutil.BeginEndTimeComputed(func() {
    // 测试的代码
    time.Sleep(100 * time.Millisecond)
})
fmt.Printf("执行耗时: %d 纳秒\n", duration)

// 时间段分割
start := time.Now()
end := start.Add(24 * time.Hour)
periods := dateutil.GetPeriods(start, end, time.Hour)

fmt.Printf("24小时分成%d个小时段:\n", len(periods))
for i, period := range periods {
    fmt.Printf("段%d: %v - %v\n", i+1, period.Start.Format("15:04"), period.End.Format("15:04"))
}
```
</details>

<details>
<summary><b>🎲 随机数和路径工具</b></summary>

```go
import (
    "github.com/karosown/katool-go/util/randutil"
    "github.com/karosown/katool-go/util/pathutil"
)

// 随机数生成
randomInt := randutil.Int(1, 100)        // 1-99之间的随机整数
randomStr := randutil.String(16)         // 16位随机字符串
uuid := randutil.UUID()                  // UUID生成

// 路径工具
currentDir := pathutil.CurrentDir()
absolutePath := pathutil.Abs("config.json")
joined := pathutil.Join("data", "files", "image.jpg")
exists := pathutil.Exists("important.txt")

if !exists {
    pathutil.EnsureDir("data/backup")    // 确保目录存在
}
```
</details>

<details>
<summary><b>🔍 调试和系统工具</b></summary>

```go
import (
    "github.com/karosown/katool-go/util/dumper"
    "github.com/karosown/katool-go/sys"
)

// 调试输出
complexObject := map[string]any{
    "users": []User{{Name: "Alice", Age: 25}},
    "config": map[string]int{"timeout": 30},
}
dumper.Dump(complexObject)  // 美化输出对象结构

// 系统工具
funcName := sys.GetLocalFunctionName()  // 获取当前函数名
fmt.Printf("当前函数: %s\n", funcName)

// 错误处理
sys.Warn("这是一个警告消息")
// sys.Panic("严重错误，程序终止")  // 会触发panic
```
</details>

<hr>

## 💡 最佳实践

<details>
<summary><b>🌊 流式处理最佳实践</b></summary>

- **优先使用 `Parallel()`**：对于大数据集，启用并行处理可显著提升性能
- **合理安排操作顺序**：先过滤再转换，减少后续处理的数据量
- **正确使用 `Reduce`**：注意提供合适的初始值和合并函数
- **避免嵌套过深**：复杂逻辑可拆分为多个步骤

```go
// ✅ 推荐写法：先过滤再转换
result := stream.ToStream(&users).
    Parallel().
    Filter(func(u User) bool { 
        return u.Age > 25  // 先过滤，减少数据量
    }).
    Map(func(u User) any {
        return u.Name  // 对过滤后的数据进行转换
    }).
    ToList()

// ❌ 不推荐：先转换再过滤
result := stream.ToStream(&users).
    Map(func(u User) any {
        return u.Name  // 转换所有数据
    }).
    Filter(func(name any) bool {
        // 过滤转换后的数据，浪费了转换资源
        return len(name.(string)) > 3
    }).
    ToList()
```
</details>

<details>
<summary><b>🔒 并发控制最佳实践</b></summary>

- **使用 `defer` 确保资源释放**：避免协程泄漏
- **合理使用 `LockSupport`**：确保每个 `Park()` 都有对应的 `Unpark()`
- **批量操作用流式API**：简化多个 `LockSupport` 的管理
- **避免死锁**：合理设计锁的获取顺序

```go
// ✅ 推荐写法：使用流式API管理多个LockSupport
supports := make([]*lock.LockSupport, n)
for i := 0; i < n; i++ {
    supports[i] = lock.NewLockSupport()
    // 启动工作协程...
}

// 批量唤醒
stream.ToStream(&supports).ForEach(func(ls *lock.LockSupport) {
    ls.Unpark()
})

// ✅ 推荐写法：确保资源释放
func processWithTimeout() {
    ls := lock.NewLockSupport()
    done := make(chan bool, 1)
    
    go func() {
        defer func() { done <- true }()
        ls.Park()
        // 处理逻辑...
    }()
    
    select {
    case <-done:
        // 正常完成
    case <-time.After(5 * time.Second):
        ls.Unpark()  // 超时唤醒
    }
}
```
</details>

<details>
<summary><b>🔤 文本处理最佳实践</b></summary>

- **及时释放资源**：使用 `defer jb.Free()` 释放分词器资源
- **选择合适的分词模式**：根据场景选择精确、全模式或搜索模式
- **合理使用词频统计**：大文本处理时注意内存使用

```go
// ✅ 推荐写法：资源管理
func processText(text string) map[string]int64 {
    jb := jieba.New()
    defer jb.Free()  // 确保资源释放
    
    words := jb.Cut(text)
    return words.Frequency().ToMap()  // 转为普通map避免持有引用
}

// ✅ 推荐写法：流式处理分词结果
meaningfulWords := jb.Cut(text).ToStream().
    Filter(func(word string) bool {
        return len(word) > 1  // 过滤单字
    }).
    Distinct().
    ToList()
```
</details>

<details>
<summary><b>🔄 数据转换最佳实践</b></summary>

- **注意类型安全**：使用泛型确保类型安全
- **合理使用属性复制**：确保源和目标结构体字段类型匹配
- **大批量转换使用并行流**：提升性能

```go
// ✅ 推荐写法：类型安全的批量转换
dtos := convert.Convert(users, func(u User) UserDTO {
    return UserDTO{
        ID:     u.ID,
        Name:   u.Name,
        Status: "Active",
    }
})

// ✅ 推荐写法：并行处理大批量数据
result := stream.ToStream(&largeDataSet).
    Parallel().
    Map(func(item DataItem) any {
        return convert.ToString(item.Value)
    }).
    ToList()
```
</details>

<details>
<summary><b>🌐 网络请求最佳实践</b></summary>

- **设置合理的超时**：避免请求hang住
- **使用链式构建**：提高代码可读性
- **正确处理错误**：检查响应状态和错误

```go
// ✅ 推荐写法：完整的错误处理
client := remote.NewRemoteRequest("https://api.example.com").
    HttpClient(resty.New().SetTimeout(30*time.Second))

var result ApiResponse
resp, err := client.
    Headers(map[string]string{"Authorization": "Bearer " + token}).
    QueryParam(map[string]string{"page": "1"}).
    Method("GET").
    Url("/api/data").
    Build(&result)

if err != nil {
    log.Printf("请求失败: %v", err)
    return
}

// 检查业务逻辑错误
if result.Code != 0 {
    log.Printf("业务错误: %s", result.Message)
    return
}
```
</details>

<hr>

## 👥 贡献指南

我们热烈欢迎社区贡献！无论是报告问题、提出建议，还是提交代码，都对项目的发展很有帮助。

### 🚀 如何参与

<table>
  <tr>
    <td><b>🐛 报告问题</b></td>
    <td>发现Bug或有改进建议？请在 <a href="https://github.com/karosown/katool-go/issues">Issues</a> 中提交</td>
  </tr>
  <tr>
    <td><b>✨ 贡献代码</b></td>
    <td>提交新功能或修复，遵循下面的开发流程</td>
  </tr>
  <tr>
    <td><b>📚 完善文档</b></td>
    <td>改进文档、添加示例或翻译</td>
  </tr>
  <tr>
    <td><b>🔧 性能优化</b></td>
    <td>提升代码性能和质量</td>
  </tr>
</table>

### 📝 开发流程

1. **Fork 本仓库**
   ```bash
   git clone https://github.com/your-username/katool-go.git
   cd katool-go
   ```

2. **创建特性分支**
   ```bash
   git checkout -b feature/amazing-feature
   # 或
   git checkout -b fix/bug-description
   ```

3. **开发和测试**
   ```bash
   # 运行测试确保不破坏现有功能
   go test ./...
   
   # 运行性能测试
   go test -bench=. ./...
   
   # 检查代码格式
   go fmt ./...
   go vet ./...
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新的流式处理功能"
   # 或
   git commit -m "fix: 修复并发访问问题"
   ```

5. **推送和创建PR**
   ```bash
   git push origin feature/amazing-feature
   ```
   然后在GitHub上创建 Pull Request

### ✅ 代码规范

请确保您的代码符合以下要求：

- **✅ 通过所有测试**：`go test ./...` 无错误
- **📏 遵循Go规范**：使用 `go fmt`、`go vet` 检查
- **📝 添加文档**：公开函数和结构体需要有注释
- **🧪 包含测试**：新功能需要有对应的测试用例
- **⚡ 性能考虑**：避免明显的性能问题

### 📋 提交信息规范

使用以下格式的提交信息：

```
type(scope): 简短描述

详细描述（可选）

Closes #issue_number
```

**类型说明：**
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建或工具相关

**示例：**
```
feat(stream): 添加并行流处理支持

- 新增 Parallel() 方法启用并行处理
- 优化大数据集的处理性能
- 添加相关测试用例

Closes #123
```

### 🔍 代码审查

我们会仔细审查每个PR，确保：

- 代码质量和性能
- 测试覆盖率
- 文档完整性
- 与现有架构的兼容性

### 🆘 获取帮助

如果您在贡献过程中遇到问题：

- 查看现有的 [Issues](https://github.com/karosown/katool-go/issues)
- 阅读项目文档和示例
- 在 Issue 中提问或讨论

<hr>

## 📄 许可证

Katool-Go 采用 **MIT 许可证**。详情请参见 [LICENSE](LICENSE) 文件。

### 📜 许可证摘要

- ✅ **商业使用**：可用于商业项目
- ✅ **修改**：可以修改源代码
- ✅ **分发**：可以分发原版或修改版
- ✅ **私用**：可用于私人项目
- ❗ **责任**：作者不承担任何责任
- ❗ **保证**：不提供任何保证

### 🤝 致谢

感谢所有为 Katool-Go 做出贡献的开发者和用户！

特别感谢以下开源项目：
- [Go 官方团队](https://golang.org/) - 提供优秀的编程语言
- [resty](https://github.com/go-resty/resty) - HTTP客户端库
- [rod](https://github.com/go-rod/rod) - Chrome控制库
- [jieba](https://github.com/yanyiwu/gojieba) - 中文分词库
- [logrus](https://github.com/sirupsen/logrus) - 日志库

---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/karosown">Karosown Team</a></sub>
  <br>
  <sub>⭐ 如果这个项目对您有帮助，请给我们一个Star！</sub>
</div>