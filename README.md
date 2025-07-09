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

<p>
  <a href="README.md">🇨🇳 中文</a> |
  <a href="README_EN.md">🇺🇸 English</a>
</p>

</div>

<hr>

## 📋 目录

- [📝 简介](#简介)
- [✨ 特性](#特性)
- [📦 安装](#安装)
- [🚀 快速开始](#快速开始)
- [🔧 核心模块](#核心模块)
  - [📚 容器与集合](#容器与集合)
    - [Optional 可选值容器](#optional-可选值容器)
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

<details>
<summary><b>📚 Optional 容器 - 安全处理空值</b></summary>

```go
package main

import (
	"fmt"
	"strings"
	"github.com/karosown/katool-go/container/optional"
)

func main() {
	// 1. 基础用法：安全处理可能为空的值
	fmt.Println("=== Optional 基础用法 ===")
	
	// 创建包含值的Optional
	nameOpt := optional.Of("张三")
	nameOpt.IfPresent(func(name string) {
		fmt.Printf("用户名: %s\n", name)
	})
	
	// 处理空值情况
	emptyOpt := optional.Empty[string]()
	username := emptyOpt.OrElse("匿名用户")
	fmt.Printf("用户名（带默认值）: %s\n", username)
	
	// 2. 函数式链式操作
	fmt.Println("\n=== 链式操作 ===")
	
	// 用户输入处理链
	userInput := "  HELLO WORLD  "
	processedInput := optional.MapTyped(optional.Of(userInput), strings.TrimSpace).
		Filter(func(s string) bool { return len(s) > 0 }).         // 过滤空字符串
		Map(func(s any) any { return strings.ToLower(s.(string)) }). // 转小写
		OrElse("无效输入")
	
	fmt.Printf("处理后的输入: %s\n", processedInput)
	
	// 3. 字符串专用处理
	fmt.Println("\n=== 字符串专用处理 ===")
	
	// StringOptional 链式处理
	result := optional.NewStringOptional("  hello world  ").
		TrimSpace().                    // 去除空格
		FilterNonEmpty().              // 过滤空字符串
		OrElse("空字符串")
	
	fmt.Printf("字符串处理结果: %s\n", result)
	
	// 4. 配置值处理
	fmt.Println("\n=== 配置值处理 ===")
	
	// 模拟从环境变量获取配置
	getConfig := func(key string) optional.Optional[string] {
		configs := map[string]string{
			"database_url": "postgres://localhost:5432/mydb",
			"redis_url":    "",  // 空值
		}
		return optional.OfNullable(configs[key])
	}
	
	// 获取数据库配置，带默认值
	dbUrl := getConfig("database_url").OrElse("sqlite://memory")
	fmt.Printf("数据库URL: %s\n", dbUrl)
	
	// 获取Redis配置，空值处理
	redisUrl := getConfig("redis_url").OrElse("redis://localhost:6379")
	fmt.Printf("Redis URL: %s\n", redisUrl)
	
	// 5. 用户验证链
	fmt.Println("\n=== 用户验证链 ===")
	
	type User struct {
		Name  string
		Age   int
		Email string
	}
	
	validateUser := func(user User) optional.Optional[User] {
		return optional.Of(user).
			Filter(func(u User) bool { return u.Name != "" }).        // 验证姓名
			Filter(func(u User) bool { return u.Age >= 18 }).         // 验证年龄
			Filter(func(u User) bool { return strings.Contains(u.Email, "@") }) // 验证邮箱
	}
	
	// 测试有效用户
	validUser := User{Name: "张三", Age: 25, Email: "zhangsan@example.com"}
	result1 := validateUser(validUser)
	result1.IfPresentOrElse(
		func(u User) { fmt.Printf("验证通过: %+v\n", u) },
		func() { fmt.Println("验证失败") },
	)
	
	// 测试无效用户
	invalidUser := User{Name: "", Age: 16, Email: "invalid-email"}
	result2 := validateUser(invalidUser)
	result2.IfPresentOrElse(
		func(u User) { fmt.Printf("验证通过: %+v\n", u) },
		func() { fmt.Println("验证失败") },
	)
	
	// 6. 条件工具函数
	fmt.Println("\n=== 条件工具函数 ===")
	
	isVIP := true
	userType := optional.IsTrue(isVIP, "VIP用户", "普通用户")
	fmt.Printf("用户类型: %s\n", userType)
	
	// 根据条件执行不同函数
	message := optional.IsTrueByFunc(isVIP,
		func() string { return "欢迎VIP用户，享受专属服务！" },
		func() string { return "欢迎使用我们的服务！" },
	)
	fmt.Printf("欢迎消息: %s\n", message)
}
```
</details>

<details>
<summary><b>⚡ 规则引擎 - 灵活的业务逻辑</b></summary>

```go
package main

import (
	"fmt"
	"time"
	"github.com/karosown/katool-go/ruleengine"
)

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Email    string    `json:"email"`
	VIPLevel int       `json:"vip_level"`
	Balance  float64   `json:"balance"`
	IDCard   string    `json:"id_card"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// 1. 创建规则引擎
	fmt.Println("=== 规则引擎基础用法 ===")
	
	engine := ruleengine.NewRuleEngine[User]()
	
	// 2. 注册验证规则
	engine.RegisterRule("validate_basic_info",
		func(user User, _ any) bool { return true },
		func(user User, _ any) (User, any, error) {
			if user.Name == "" {
				return user, "用户名不能为空", ruleengine.EOF
			}
			if len(user.Name) < 2 {
				return user, "用户名太短", ruleengine.EOF
			}
			return user, "基础信息验证通过", nil
		},
	)
	
	// 3. 年龄检查规则（含流程控制）
	engine.RegisterRule("check_age",
		func(user User, _ any) bool { return true },
		func(user User, _ any) (User, any, error) {
			if user.Age < 13 {
				return user, "用户年龄过小", ruleengine.EOF // 立即终止
			} else if user.Age < 18 {
				return user, "未成年用户", ruleengine.FALLTHROUGH // 跳过成年用户逻辑
			}
			return user, "成年用户", nil
		},
	)
	
	// 4. 成年用户身份验证（未成年用户会跳过）
	engine.RegisterRule("adult_identity_check",
		func(user User, _ any) bool { return user.Age >= 18 },
		func(user User, _ any) (User, any, error) {
			if user.IDCard == "" {
				return user, "成年用户需要身份证", ruleengine.EOF
			}
			return user, "身份验证完成", nil
		},
	)
	
	// 5. VIP特权检查
	engine.RegisterRule("vip_privilege_check",
		func(user User, _ any) bool { return user.VIPLevel > 0 },
		func(user User, _ any) (User, any, error) {
			if user.VIPLevel >= 3 {
				user.Balance += 100.0  // VIP3以上赠送余额
				return user, "VIP特权已激活", nil
			} else if user.VIPLevel >= 1 {
				user.Balance += 50.0   // VIP1-2赠送部分余额
				return user, "VIP福利已发放", nil
			}
			return user, "普通用户", nil
		},
	)
	
	// 6. 最终注册
	engine.RegisterRule("complete_registration",
		func(user User, _ any) bool { return true },
		func(user User, _ any) (User, any, error) {
			if user.ID == 0 {
				user.ID = int(time.Now().Unix()) // 生成ID
			}
			user.CreatedAt = time.Now()
			return user, "注册完成", nil
		},
	)
	
	// 7. 添加日志中间件
	engine.AddMiddleware(func(data User, next func(User) (User, any, error)) (User, any, error) {
		fmt.Printf("  → 处理用户: %s (年龄: %d)\n", data.Name, data.Age)
		result, info, err := next(data)
		if err == ruleengine.EOF {
			fmt.Printf("  ✖ 流程终止: %v\n", info)
		} else if err == ruleengine.FALLTHROUGH {
			fmt.Printf("  ⚡ 规则跳过: %v\n", info)
		} else if err == nil {
			fmt.Printf("  ✓ 执行成功: %v\n", info)
		} else {
			fmt.Printf("  ✗ 执行失败: %v\n", err)
		}
		return result, info, err
	})
	
	// 8. 构建注册流程链
	_, err := engine.NewBuilder("user_registration").
		AddRule("validate_basic_info").
		AddRule("check_age").
		AddRule("adult_identity_check").
		AddRule("vip_privilege_check").
		AddRule("complete_registration").
		Build()
	
	if err != nil {
		fmt.Printf("构建规则链失败: %v\n", err)
		return
	}
	
	// 9. 测试不同场景
	fmt.Println("\n=== 测试场景 1: 正常成年VIP用户 ===")
	adultVIP := User{
		Name:     "张三",
		Age:      25,
		Email:    "zhangsan@example.com",
		VIPLevel: 3,
		IDCard:   "123456789012345678",
		Balance:  0,
	}
	result1 := engine.Execute("user_registration", adultVIP)
	fmt.Printf("最终结果: ID=%d, 余额=%.2f\n", result1.Data.ID, result1.Data.Balance)
	
	fmt.Println("\n=== 测试场景 2: 未成年用户（跳过身份验证）===")
	minor := User{
		Name:     "李四",
		Age:      16,
		Email:    "lisi@example.com",
		VIPLevel: 1,
		Balance:  0,
	}
	result2 := engine.Execute("user_registration", minor)
	fmt.Printf("最终结果: ID=%d, 余额=%.2f\n", result2.Data.ID, result2.Data.Balance)
	
	fmt.Println("\n=== 测试场景 3: 年龄过小（立即终止）===")
	child := User{
		Name:     "王五",
		Age:      10,
		Email:    "wangwu@example.com",
		VIPLevel: 0,
		Balance:  0,
	}
	result3 := engine.Execute("user_registration", child)
	if result3.Error != nil {
		fmt.Printf("注册失败: %v\n", result3.Error)
	}
	
	fmt.Println("\n=== 测试场景 4: 批量处理多个用户 ===")
	users := []User{
		{Name: "用户A", Age: 25, VIPLevel: 2, IDCard: "111111111111111111"},
		{Name: "用户B", Age: 17, VIPLevel: 1},
		{Name: "", Age: 30, VIPLevel: 0},  // 无效用户名
	}
	
	for i, user := range users {
		fmt.Printf("\n--- 处理用户 %d ---\n", i+1)
		result := engine.Execute("user_registration", user)
		if result.Error != nil && result.Error != ruleengine.EOF && result.Error != ruleengine.FALLTHROUGH {
			fmt.Printf("处理失败: %v\n", result.Error)
		} else {
			fmt.Printf("处理完成: ID=%d\n", result.Data.ID)
		}
	}
}
```
</details>

<hr>

## 🔧 核心模块

### 📚 容器与集合

Katool-Go 提供了丰富的容器和集合类型，全部支持泛型，提供类型安全的操作。

#### Optional 可选值容器

Optional 是一个用于安全处理可能为空值的容器类型，灵感来自 Java 的 Optional 类，提供类型安全的空值处理机制。

##### 🚀 基础用法

```go
import "github.com/karosown/katool-go/container/optional"

// 创建包含值的Optional
opt := optional.Of("Hello World")

// 创建空的Optional
emptyOpt := optional.Empty[string]()

// 根据值是否为零值创建Optional
nullableOpt := optional.OfNullable("")  // 空字符串会创建空Optional
```

##### 🔍 安全检查和获取

```go
// 安全检查和获取
if opt.IsPresent() {
    value := opt.Get()
    fmt.Println("值存在:", value)
}

// 检查是否为空
if emptyOpt.IsEmpty() {
    fmt.Println("Optional为空")
}

// 提供默认值的几种方式
defaultValue := emptyOpt.OrElse("默认值")
lazyDefault := emptyOpt.OrElseGet(func() string {
    return "延迟计算的默认值"
})
safeValue := opt.OrElsePanic("Optional不能为空!")
```

##### ⚡ 函数式操作

```go
// 条件执行 - 有值时执行
opt.IfPresent(func(v string) {
    fmt.Println("处理值:", v)
})

// 双分支执行 - 有值执行第一个函数，无值执行第二个
opt.IfPresentOrElse(
    func(v string) { fmt.Println("有值:", v) },
    func() { fmt.Println("无值") },
)

// 过滤操作
filtered := opt.Filter(func(s string) bool {
    return len(s) > 5
})

// 类型安全的映射（推荐）
result := optional.MapTyped(optional.Of("  hello  "), strings.TrimSpace).
    Filter(func(s string) bool { return len(s) > 0 }).
    OrElse("空字符串")
```

##### 🔤 字符串处理专用

为了更好地支持字符串处理，提供了专用的 StringOptional：

```go
// 专用的StringOptional进行链式字符串处理
result := optional.NewStringOptional("  hello  ").
    TrimSpace().                    // 去除空格
    FilterNonEmpty().              // 过滤空字符串
    OrElse("空字符串")             // 提供默认值

fmt.Println("处理结果:", result) // 输出: 处理结果: hello
```

##### 🛠️ 实用工具函数

```go
// 根据条件返回不同的值
enabled := optional.IsTrue(condition, "启用", "禁用")

// 根据条件调用不同的函数
result := optional.IsTrueByFunc(condition, 
    func() string { return "功能已启用" },
    func() string { return "功能已禁用" },
)

// 根据函数条件调用不同的函数
result := optional.FuncIsTrueByFunc(
    func() bool { return someComplexCondition() },
    enabledFunc,
    disabledFunc,
)
```

##### 📝 实用示例

**用户输入处理**
```go
func processUserInput(input string) string {
    return optional.MapTyped(optional.Of(input), strings.TrimSpace).
        Filter(func(s string) bool { return len(s) > 0 }).
        Map(func(s any) any { return strings.ToLower(s.(string)) }).
        OrElse("无效输入").(string)
}
```

**配置值处理**
```go
func getConfig(key string) optional.Optional[string] {
    if value := os.Getenv(key); value != "" {
        return optional.Of(value)
    }
    return optional.Empty[string]()
}

// 使用
dbUrl := getConfig("DATABASE_URL").OrElse("sqlite://default.db")
```

**用户验证链式处理**
```go
func validateUser(user User) optional.Optional[User] {
    return optional.Of(user).
        Filter(func(u User) bool { return u.Name != "" }).
        Filter(func(u User) bool { return u.Age >= 18 }).
        Filter(func(u User) bool { return u.Email != "" })
}

// 使用
validUser := validateUser(user).OrElsePanic("用户验证失败")
```

##### 📋 API 参考

**核心方法：**
- `Of[T](value T)` - 创建包含值的Optional
- `Empty[T]()` - 创建空Optional
- `OfNullable[T](value T)` - 根据零值创建Optional

**检查方法：**
- `IsPresent()` - 检查是否有值
- `IsEmpty()` - 检查是否为空

**获取方法：**
- `Get()` - 获取值（空时panic）
- `OrElse(T)` - 提供默认值
- `OrElseGet(func() T)` - 延迟计算默认值
- `OrElsePanic(string)` - 空时panic并显示消息

**函数式方法：**
- `IfPresent(func(T))` - 条件执行
- `IfPresentOrElse(func(T), func())` - 双分支执行
- `Filter(func(T) bool)` - 过滤
- `Map(func(T) any)` - 映射（实例方法）
- `MapTyped[T,R](Optional[T], func(T) R)` - 类型安全映射

##### ⚠️ 注意事项

1. **类型安全**: 使用 `MapTyped` 进行类型安全的映射操作
2. **链式调用**: 实例方法支持链式调用，但要注意类型转换
3. **性能**: Optional 会带来轻微的性能开销，在性能敏感的场景中谨慎使用
4. **空指针**: Optional 本身不会为 nil，但内部值可能是零值

### 🌊 流式处理

提供类似 Java 8 Stream API 的强大流式处理能力，支持并行计算和链式操作。

```go
import "github.com/karosown/katool-go/container/stream"

// 并行流处理
results := stream.ToStream(&data).
    Parallel().                               // 启用并行处理
    Filter(func(item Item) bool { return item.IsValid() }).
    Map(func(item Item) ProcessedItem { return item.Process() }).
    Sort(func(a, b ProcessedItem) bool { return a.Priority > b.Priority }).
    ToList()
```

### 🔄 数据转换

强大的数据转换和结构体处理能力。

```go
import "github.com/karosown/katool-go/convert"

// 结构体复制
var dest DestStruct
convert.CopyStruct(&dest, &source)

// 数据导出
convert.ExportToCSV(data, "output.csv")
convert.ExportToJSON(data, "output.json")
```

### 💉 依赖注入

轻量级IOC容器，简化依赖管理。

```go
import "github.com/karosown/katool-go/container/ioc"

// 注册服务
container := ioc.NewContainer()
container.Register("userService", &UserService{})

// 获取服务
userSvc := container.Get("userService").(*UserService)
```

### 🔒 并发控制

提供类似Java的并发控制工具。

```go
import "github.com/karosown/katool-go/lock"

// LockSupport 类似Java的park/unpark
lock.LockSupport.Park()        // 阻塞当前协程
lock.LockSupport.Unpark(goroutineId) // 唤醒指定协程
```

### 🕸️ Web爬虫

智能内容提取和网页爬取工具。

```go
import "github.com/karosown/katool-go/web_crawler"

// 内容提取
extractor := web_crawler.NewContentExtractor()
content := extractor.ExtractFromURL("https://example.com")

// Chrome渲染支持
renderer := web_crawler.NewChromeRenderer()
html := renderer.RenderPage("https://spa-app.com")
```

### 📁 文件操作

完整的文件系统操作工具。

```go
import "github.com/karosown/katool-go/file"

// 文件下载
downloader := file.NewDownloader()
downloader.Download("https://example.com/file.zip", "./downloads/")

// 序列化操作
file.SerializeToFile(data, "data.json")
data := file.DeserializeFromFile[MyStruct]("data.json")
```

### 💾 数据库支持

MongoDB等数据库操作增强。

```go
import "github.com/karosown/katool-go/db"

// MongoDB分页查询
paginator := db.NewMongoPaginator(collection)
result := paginator.Page(1).Size(20).Find(filter)
```

### 🌐 网络通信

现代化HTTP客户端和网络工具。

```go
import "github.com/karosown/katool-go/net/http/remote"

// 链式HTTP请求构建
var result APIResponse
resp, err := remote.NewRemoteRequest("https://api.example.com").
    Headers(map[string]string{"Authorization": "Bearer " + token}).
    QueryParam(map[string]string{"page": "1"}).
    Method("GET").
    Url("/api/data").
    Build(&result)
```

### 📝 日志系统

结构化日志和链式构建。

```go
import "github.com/karosown/katool-go/xlog"

// 结构化日志
logger := xlog.NewLogger().
    WithField("service", "user-api").
    WithField("version", "1.0.0")

logger.Info("用户登录成功").
    WithField("userId", userId).
    WithField("ip", clientIP).
    Log()
```

### ⚙️ 算法工具

实用算法和数据结构。

```go
import "github.com/karosown/katool-go/algorithm"

// 有序数组合并
merged := algorithm.MergeSortedArrays(arr1, arr2)

// 哈希计算
hash := algorithm.ComputeHash(data)
```

### 🔤 文本处理

中文分词和文本分析。

```go
import "github.com/karosown/katool-go/words"

// 中文分词
segmenter := words.NewJiebaSegmenter()
tokens := segmenter.Cut("这是一个中文分词测试", true)

// 词频统计
counter := words.NewWordCounter()
frequencies := counter.Count(tokens)
```

### 🧰 辅助工具

实用的开发辅助工具。

```go
import "github.com/karosown/katool-go/util"

// 日期处理
date := util.ParseDate("2023-12-25")
formatted := util.FormatDate(date, "YYYY-MM-DD")

// 随机数生成
randomStr := util.RandomString(10)
randomInt := util.RandomInt(1, 100)
```

### ⚡ 规则引擎

灵活强大的业务规则处理引擎，支持规则链、规则树和中间件机制。支持泛型、并发安全，提供EOF和FALLTHROUGH流程控制。

#### 🚀 快速开始

```go
import "github.com/karosown/katool-go/ruleengine"

// 1. 创建规则引擎
engine := ruleengine.NewRuleEngine[User]()

// 2. 注册规则
engine.RegisterRule("validate_age",
    func(user User, _ any) bool { return user.Age > 0 },  // 验证函数
    func(user User, _ any) (User, any, error) {           // 执行函数
        if user.Age < 18 {
            return user, "未成年用户", nil
        }
        return user, "成年用户", nil
    },
)

// 3. 构建规则链
engine.NewBuilder("user_processing").
    AddRule("validate_age").
    Build()

// 4. 执行规则
user := User{Name: "张三", Age: 25}
result := engine.Execute("user_processing", user)
fmt.Printf("处理结果: %v\n", result.Result)
```

#### 🛠️ 高级功能

##### 中间件支持

```go
// 添加日志中间件
engine.AddMiddleware(func(data User, next func(User) (User, any, error)) (User, any, error) {
    fmt.Printf("执行前: %+v\n", data)
    result, info, err := next(data)
    fmt.Printf("执行后: %+v\n", result)
    return result, info, err
})

// 添加性能监控中间件
engine.AddMiddleware(func(data User, next func(User) (User, any, error)) (User, any, error) {
    start := time.Now()
    result, info, err := next(data)
    fmt.Printf("执行耗时: %v\n", time.Since(start))
    return result, info, err
})
```

##### 错误控制机制

**EOF - 立即终止执行**
```go
// 当遇到严重问题时，立即终止整个规则链
return user, "用户被禁用", ruleengine.EOF
```

**FALLTHROUGH - 跳过当前规则继续执行**
```go
// 跳过当前规则，但继续执行后续规则
return user, "跳过此步骤", ruleengine.FALLTHROUGH
```

#### 🌳 规则树结构

除了线性的规则链，还支持树形结构的规则组织：

##### 基础用法

```go
type TestData struct {
    Value int
}

// 创建规则节点
leafNode := ruleengine.NewRuleNode[TestData](
    func(data TestData, _ any) bool { return data.Value > 5 },
    func(data TestData, _ any) (TestData, any, error) {
        return TestData{Value: data.Value + 10}, "处理完成", nil
    },
)

// 创建规则树
tree := ruleengine.NewRuleTree[TestData](leafNode)

// 执行规则树
result, info, err := tree.Run(TestData{Value: 3})
```

##### 复杂树形结构

```go
// 构建复杂的规则树
rootNode := ruleengine.NewRuleNode[User](
    func(user User, _ any) bool { return user.ID > 0 },
    func(user User, _ any) (User, any, error) {
        return user, "用户ID验证通过", nil
    },
)

// 添加子节点
ageNode := ruleengine.NewRuleNode[User](
    func(user User, _ any) bool { return user.Age > 0 },
    func(user User, _ any) (User, any, error) {
        return user, "年龄验证通过", nil
    },
)

emailNode := ruleengine.NewRuleNode[User](
    func(user User, _ any) bool { return user.Email != "" },
    func(user User, _ any) (User, any, error) {
        return user, "邮箱验证通过", nil
    },
)

// 构建树形结构
rootNode.AddChild(ageNode)
rootNode.AddChild(emailNode)

tree := ruleengine.NewRuleTree[User](rootNode)
```

#### 📝 复杂业务场景示例

##### 用户注册验证流程

```go
func setupUserRegistrationEngine() *ruleengine.RuleEngine[User] {
    engine := ruleengine.NewRuleEngine[User]()
    
    // 基础信息验证
    engine.RegisterRule("validate_basic_info",
        func(user User, _ any) bool { return true },
        func(user User, _ any) (User, any, error) {
            if user.Name == "" {
                return user, "用户名不能为空", ruleengine.EOF
            }
            if len(user.Name) < 2 {
                return user, "用户名太短", ruleengine.EOF
            }
            return user, "基础信息验证通过", nil
        },
    )
    
    // 年龄检查
    engine.RegisterRule("check_age",
        func(user User, _ any) bool { return true },
        func(user User, _ any) (User, any, error) {
            if user.Age < 13 {
                return user, "用户年龄过小", ruleengine.EOF
            } else if user.Age < 18 {
                return user, "未成年用户", ruleengine.FALLTHROUGH
            }
            return user, "成年用户", nil
        },
    )
    
    // 成年用户身份验证（未成年用户会跳过）
    engine.RegisterRule("adult_identity_check",
        func(user User, _ any) bool { return user.Age >= 18 },
        func(user User, _ any) (User, any, error) {
            if user.IDCard == "" {
                return user, "成年用户需要身份证", ruleengine.EOF
            }
            return user, "身份验证完成", nil
        },
    )
    
    // 邮箱验证
    engine.RegisterRule("validate_email",
        func(user User, _ any) bool { return user.Email != "" },
        func(user User, _ any) (User, any, error) {
            if !isValidEmail(user.Email) {
                return user, "邮箱格式错误", ruleengine.EOF
            }
            return user, "邮箱验证通过", nil
        },
    )
    
    // 最终注册
    engine.RegisterRule("complete_registration",
        func(user User, _ any) bool { return true },
        func(user User, _ any) (User, any, error) {
            user.ID = generateUserID()
            user.CreatedAt = time.Now()
            return user, "注册完成", nil
        },
    )
    
    // 构建注册流程链
    engine.NewBuilder("user_registration").
        AddRule("validate_basic_info").
        AddRule("check_age").
        AddRule("adult_identity_check").
        AddRule("validate_email").
        AddRule("complete_registration").
        Build()
    
    return engine
}

// 使用示例
func registerUser(userData User) {
    engine := setupUserRegistrationEngine()
    
    result := engine.Execute("user_registration", userData)
    if result.Error != nil {
        fmt.Printf("注册失败: %v\n", result.Error)
        return
    }
    
    fmt.Printf("注册成功: %+v\n", result.Data)
    fmt.Printf("处理信息: %v\n", result.Result)
}
```

##### 复杂执行场景分析

```go
// 执行结果分析：
用户年龄 12: validate_basic_info(✅) → check_age(EOF 🛑) → 后续规则全部跳过
用户年龄 16: validate_basic_info(✅) → check_age(FALLTHROUGH ⚡) → adult_identity_check(跳过) → validate_email(✅) → complete_registration(✅)
用户年龄 25: validate_basic_info(✅) → check_age(✅) → adult_identity_check(✅) → validate_email(✅) → complete_registration(✅)
```

#### 🔄 批量执行

```go
// 批量执行多个规则链
users := []User{
    {Name: "张三", Age: 25, Email: "zhang@example.com"},
    {Name: "李四", Age: 17, Email: "li@example.com"},
    {Name: "王五", Age: 30, Email: "wang@example.com"},
}

chains := []string{"user_registration", "user_validation"}

for _, user := range users {
    results := engine.BatchExecute(chains, user)
    for i, result := range results {
        fmt.Printf("用户 %s 执行链 %s: ", user.Name, chains[i])
        if result.Error != nil {
            fmt.Printf("失败 - %v\n", result.Error)
        } else {
            fmt.Printf("成功 - %v\n", result.Result)
        }
    }
}
```

#### 📚 API 参考

##### 核心引擎方法

**创建与配置：**
- `NewRuleEngine[T]()` - 创建新的规则引擎
- `RegisterRule(name, validFunc, execFunc)` - 注册规则
- `AddMiddleware(middleware)` - 添加中间件

**规则链构建：**
- `NewBuilder(chainName)` - 创建规则链构建器
- `AddRule(ruleName)` - 添加已注册的规则
- `AddCustomRule(validFunc, execFunc)` - 添加临时规则
- `Build()` - 构建规则链

**执行方法：**
- `Execute(chainName, data)` - 执行指定规则链
- `BatchExecute(chainNames, data)` - 批量执行多个规则链

##### 规则树方法

**树结构构建：**
- `NewRuleNode[T](validFunc, execFunc)` - 创建规则节点
- `AddChild(childNode)` - 添加子节点
- `AddChildren(childNodes...)` - 添加多个子节点

**树执行：**
- `NewRuleTree[T](rootNode)` - 创建规则树
- `Run(data)` - 执行规则树
- `ToQueue()` - 转换为队列形式

##### 错误控制常量

- `ruleengine.EOF` - 立即终止执行
- `ruleengine.FALLTHROUGH` - 跳过当前规则继续执行

##### 执行结果结构

```go
type ExecuteResult[T any] struct {
    Data   T           // 处理后的数据
    Result any         // 执行结果信息
    Error  error       // 错误信息
}
```

#### ⚠️ 注意事项

1. **系统要求**: 需要 Go 1.18+ (泛型支持)
2. **线程安全**: 引擎实例支持并发访问
3. **规则命名**: 建议使用 `动词_名词` 格式，如 `validate_email`
4. **错误控制**: 
   - 使用 `EOF` 处理严重错误，立即终止
   - 使用 `FALLTHROUGH` 跳过可选逻辑
5. **性能优化**: 
   - 合理设计规则粒度，避免单个规则过于复杂
   - 善用中间件处理横切关注点
   - 规则链顺序影响性能，将高频失败的规则前置

#### 🎯 最佳实践

1. **单一职责**: 每个规则只处理一种业务逻辑
2. **合理分层**: 基础验证 → 业务逻辑 → 数据处理 → 最终确认
3. **错误处理**: 区分业务错误（FALLTHROUGH）和系统错误（EOF）
4. **中间件使用**: 用于日志、监控、缓存等横切关注点
5. **测试覆盖**: 为每个规则和规则链编写单元测试

```go
// 规则测试示例
func TestValidateAgeRule(t *testing.T) {
    engine := ruleengine.NewRuleEngine[User]()
    
    engine.RegisterRule("validate_age",
        func(user User, _ any) bool { return user.Age > 0 },
        func(user User, _ any) (User, any, error) {
            if user.Age < 18 {
                return user, "未成年", ruleengine.FALLTHROUGH
            }
            return user, "成年", nil
        },
    )
    
    engine.NewBuilder("test_chain").AddRule("validate_age").Build()
    
    // 测试未成年用户
    minorResult := engine.Execute("test_chain", User{Age: 16})
    assert.Equal(t, ruleengine.FALLTHROUGH, minorResult.Error)
    assert.Equal(t, "未成年", minorResult.Result)
    
    // 测试成年用户
    adultResult := engine.Execute("test_chain", User{Age: 25})
    assert.Nil(t, adultResult.Error)
    assert.Equal(t, "成年", adultResult.Result)
}
```

#### 📊 可视化流程图

##### EOF 机制 - 立即终止执行

```
正常执行流程：
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   规则 A    │───▶│   规则 B    │───▶│   规则 C    │───▶│   规则 D    │
│  (验证通过)  │    │  (验证通过)  │    │  (验证通过)  │    │  (验证通过)  │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       ✅               ✅               ✅               ✅

EOF 终止流程：
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   规则 A    │───▶│   规则 B    │ ╳  │   规则 C    │    │   规则 D    │
│  (验证通过)  │    │ (返回 EOF)  │    │  (未执行)   │    │  (未执行)   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       ✅               🛑               ⭕               ⭕
                   立即终止，后续规则不执行

规则树中的 EOF：
                    根节点
                       │
                   ┌───▼───┐
                   │ 规则A │ ✅
                   └───┬───┘
                       │
            ┌──────────┼──────────┐
            ▼          ▼          ▼
        ┌──────┐   ┌──────┐   ┌──────┐
        │规则B1│   │规则B2│   │规则B3│
        │ (EOF)│   │(未执行)│ │(未执行)│
        └──────┘   └──────┘   └──────┘
            🛑         ⭕         ⭕
        
        当B1返回EOF时，整个树立即终止
        B2、B3 以及所有后续节点都不会执行
```

##### FALLTHROUGH 机制 - 跳过继续执行

```
正常执行流程：
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   规则 A    │───▶│   规则 B    │───▶│   规则 C    │───▶│   规则 D    │
│  (验证通过)  │    │  (验证通过)  │    │  (验证通过)  │    │  (验证通过)  │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       ✅               ✅               ✅               ✅

FALLTHROUGH 跳过流程：
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   规则 A    │───▶│   规则 B    │~~~▶│   规则 C    │───▶│   规则 D    │
│  (验证通过)  │    │(FALLTHROUGH)│    │  (验证通过)  │    │  (验证通过)  │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       ✅               ⚡               ✅               ✅
                   跳过但继续执行后续规则

规则树中的 FALLTHROUGH：
                      根节点
                         │
                     ┌───▼───┐
                     │ 规则A │ ✅
                     └───┬───┘
                         │
              ┌──────────┼──────────┐
              ▼          ▼          ▼
          ┌──────┐   ┌──────┐   ┌──────┐
          │规则B1│   │规则B2│   │规则B3│
          │(FALL)│   │ ✅   │   │ ✅   │
          └──┬───┘   └──┬───┘   └──┬───┘
             │⚡       │        │
             ▼          ▼          ▼
          ┌──────┐   ┌──────┐   ┌──────┐
          │规则C1│   │规则C2│   │规则C3│
          │(跳过) │   │ ✅   │   │ ✅   │
          └──────┘   └──────┘   └──────┘
             ⭕
        
        当B1返回FALLTHROUGH时：
        - B1的子节点C1被跳过
        - B2、B3 继续正常执行
        - C2、C3 继续正常执行
```

<hr>

## 💡 最佳实践

### 🚀 性能优化建议

<details>
<summary><b>🌊 流式处理性能优化</b></summary>

- **合理使用并行流**：大数据集(>1000元素)时启用`Parallel()`
- **避免频繁装箱**：使用具体类型而非interface{}
- **链式操作排序**：先Filter再Map，减少处理元素数量

```go
// ✅ 推荐：先过滤再处理
stream.ToStream(&data).
    Filter(func(item Item) bool { return item.IsValid() }).  // 先减少数据量
    Map(func(item Item) ProcessedItem { return item.Process() }).
	ToList()

// ❌ 避免：先处理再过滤
stream.ToStream(&data).
    Map(func(item Item) ProcessedItem { return item.Process() }).  // 处理所有数据
    Filter(func(item ProcessedItem) bool { return item.IsValid() }). // 再过滤
	ToList()
```
</details>

<details>
<summary><b>📚 Optional 容器最佳实践</b></summary>

- **避免嵌套Optional**：不要创建`Optional[Optional[T]]`
- **使用类型安全的MapTyped**：避免类型断言错误
- **合理使用OrElsePanic**：仅在确定不会为空时使用

```go
// ✅ 推荐：使用MapTyped进行类型安全转换
result := optional.MapTyped(optional.Of("  hello  "), strings.TrimSpace).
    Filter(func(s string) bool { return len(s) > 0 }).
    OrElse("默认值")

// ❌ 避免：使用Map需要类型断言
result := optional.Of("  hello  ").
    Map(func(s any) any { return strings.TrimSpace(s.(string)) }). // 需要断言
    OrElse("默认值")
```
</details>

<details>
<summary><b>⚡ 规则引擎最佳实践</b></summary>

- **规则粒度控制**：单个规则只处理一种业务逻辑
- **合理使用中间件**：用于日志、监控，避免业务逻辑
- **错误控制策略**：EOF用于严重错误，FALLTHROUGH用于跳过逻辑

```go
// ✅ 推荐：单一职责的规则
engine.RegisterRule("validate_email",
    func(user User, _ any) bool { return user.Email != "" },
    func(user User, _ any) (User, any, error) {
        if !isValidEmail(user.Email) {
            return user, "邮箱格式错误", ruleengine.EOF
        }
        return user, "邮箱验证通过", nil
    },
)

// ❌ 避免：复杂的多职责规则
engine.RegisterRule("validate_user",  // 太宽泛
    func(user User, _ any) bool { return true },
    func(user User, _ any) (User, any, error) {
        // 验证邮箱、手机、身份证等多种逻辑混合
        // 违反单一职责原则
    },
)
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