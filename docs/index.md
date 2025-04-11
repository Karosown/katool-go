---
home: true
heroImage: /logo.png
heroText: Katool-Go
tagline: 功能丰富的 Go 工具库，借鉴 Java 生态的优秀设计
actionText: 快速上手 →
actionLink: /guide/getting-started/
features:
- title: 流式处理
  details: 提供类似 Java Stream API 的链式调用，支持映射、过滤、归约等操作，让数据处理更简洁优雅。
- title: 依赖注入
  details: 轻量级 IOC 容器，管理依赖关系，使代码更易于测试和维护。
- title: 强大工具集
  details: 包含集合操作、锁支持、类型转换、Web 爬虫等实用工具，满足多种开发场景需求。
- title: 并发支持
  details: 提供高效的并行流和协程控制工具，充分利用多核能力。
- title: 易于扩展
  details: 模块化设计，接口清晰，便于扩展和自定义功能。
- title: 类型安全
  details: 充分利用 Go 泛型，在保持灵活性的同时提供类型安全。
footer: MIT Licensed | Copyright © 2023-present Karosown
---

## 轻松上手

### 安装

```bash
go get github.com/karosown/katool-go
```

### 简单示例

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/container/stream"
)

func main() {
    // 示例数据
    data := []int{1, 2, 3, 4, 5}
    
    // 使用流处理
    result := stream.ToStream(&data).
        Filter(func(n int) bool {
            return n % 2 == 0
        }).
        Map(func(n int) any {
            return n * n
        }).
        ToList()
    
    fmt.Println(result) // [4, 16]
}
```

访问 [快速开始](/guide/getting-started/) 了解更多使用方式。 