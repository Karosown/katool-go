# 快速开始

本指南将帮助你快速上手 Katool-Go 库，了解其核心功能并开始使用。

## 安装

使用 Go 模块安装 Katool-Go：

```bash
go get github.com/karosown/katool-go
```

## 导入

根据需要导入相应的包：

```go
import (
    // 流式处理
    "github.com/karosown/katool-go/container/stream"
    
    // 集合操作
    "github.com/karosown/katool-go/collect/lists"
    
    // IOC 容器
    "github.com/karosown/katool-go/container/ioc"
    
    // 锁支持
    "github.com/karosown/katool-go/lock"
    
    // 转换工具
    "github.com/karosown/katool-go/convert"
    
    // 其他可选包
    "github.com/karosown/katool-go/algorithm"
    "github.com/karosown/katool-go/web_crawler"
    "github.com/karosown/katool-go/log"
)
```

## 基本使用

### 使用 Stream 流处理数据

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/container/stream"
)

func main() {
    // 示例数据
    data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 使用流处理：过滤、映射、收集
    result := stream.ToStream(&data).
        Filter(func(n int) bool {
            return n%2 == 0 // 只保留偶数
        }).
        Map(func(n int) any {
            return n * n // 计算平方
        }).
        ToList() // 收集结果
    
    fmt.Println(result) // 输出: [4 16 36 64 100]
}
```

### 使用集合分区并行处理

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/collect/lists"
    "github.com/karosown/katool-go/container/stream"
    "github.com/Tangerg/lynx/pkg/sync"
)

func main() {
    // 示例数据
    data := make([]int, 100)
    for i := 0; i < 100; i++ {
        data[i] = i
    }
    
    // 将数据分成批次并并行处理
    lists.Partition(data, 10).ForEach(func(pos int, batch []int) error {
        fmt.Printf("处理第 %d 批数据，大小: %d\n", pos, len(batch))
        
        // 每个批次内部使用流处理
        sum := 0
        stream.ToStream(&batch).ForEach(func(n int) {
            sum += n
        })
        
        fmt.Printf("批次 %d 的和: %d\n", pos, sum)
        return nil
    }, true, sync.NewLimiter(4)) // 并行处理，最多4个并发
}
```

### 使用 IOC 容器管理依赖

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/container/ioc"
)

// 定义服务接口
type UserService interface {
    GetUserName(id int) string
}

// 实现服务
type UserServiceImpl struct{}

func (s *UserServiceImpl) GetUserName(id int) string {
    return fmt.Sprintf("User-%d", id)
}

func main() {
    // 注册服务
    ioc.RegisterValue("userService", &UserServiceImpl{})
    
    // 获取服务
    service := ioc.Get("userService").(UserService)
    
    // 使用服务
    fmt.Println(service.GetUserName(123)) // 输出: User-123
}
```

## 下一步

- 阅读各功能模块的详细文档
- 查看示例代码深入理解用法
- 探索更高级的功能和配置选项 