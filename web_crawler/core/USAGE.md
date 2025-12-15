# Contain 使用说明

## 全局池 vs 非全局池

`Contain` 支持两种浏览器管理方式：

### 1. 全局池模式（默认）

使用全局浏览器池，多个 `Contain` 实例共享浏览器资源。

**优点**：
- 资源复用，减少浏览器创建开销
- 适合高并发场景
- 自动管理浏览器生命周期

**缺点**：
- 所有实例共享同一个池，可能相互影响
- 无法独立控制每个实例的浏览器

**使用方式**：

```go
// 使用全局池（默认）
chrome := core.NewCotain("path/to/chrome", true)
defer chrome.Close()
```

### 2. 非全局池模式

每个 `Contain` 实例独立管理自己的浏览器，不使用全局池。

**优点**：
- 完全独立，互不影响
- 可以精确控制每个实例的浏览器生命周期
- 适合需要隔离的场景

**缺点**：
- 每次创建新浏览器，资源开销较大
- 需要手动管理浏览器生命周期

**使用方式**：

```go
// 不使用全局池，独立管理
chrome := core.NewContainWithoutPool("path/to/chrome", true)
defer chrome.Close()
```

## 浏览器和页面管理

### 浏览器管理

#### 全局池模式

```go
chrome := core.NewCotain("path/to/chrome", true)

// 使用浏览器
page := chrome.PageWithStealth("https://example.com")
// ... 使用页面
page.Close()

// 关闭浏览器
chrome.Close()
// 行为：
// 1. 关闭所有页面
// 2. Kill 当前浏览器
// 3. 创建新浏览器放入全局池（如果池未满）
```

#### 非全局池模式

```go
chrome := core.NewContainWithoutPool("path/to/chrome", true)

// 使用浏览器
page := chrome.PageWithStealth("https://example.com")
// ... 使用页面
page.Close()

// 关闭浏览器
chrome.Close()
// 行为：
// 1. 关闭所有页面
// 2. Kill 当前浏览器
// 3. 不创建新浏览器（完全关闭）
```

### 页面管理

无论使用哪种模式，页面管理都是相同的：

```go
// 创建新页面
page := chrome.PageWithStealth("https://example.com")

// 使用页面
// ... 操作页面

// 关闭页面（重要：必须关闭，否则会占用浏览器资源）
defer page.Close()
```

## 完整示例

### 示例 1：全局池模式

```go
package main

import (
    "github.com/karosown/katool-go/web_crawler/core"
)

func main() {
    // 创建使用全局池的浏览器实例
    chrome := core.NewCotain("C:\\path\\to\\chrome.exe", true)
    defer chrome.Close()
    
    // 创建页面
    page := chrome.PageWithStealth("https://example.com")
    defer page.Close()
    
    // 使用页面
    html := page.MustHTML()
    println(html)
}
```

### 示例 2：非全局池模式

```go
package main

import (
    "github.com/karosown/katool-go/web_crawler/core"
)

func main() {
    // 创建不使用全局池的浏览器实例
    chrome := core.NewContainWithoutPool("C:\\path\\to\\chrome.exe", true)
    defer chrome.Close()
    
    // 创建页面
    page := chrome.PageWithStealth("https://example.com")
    defer page.Close()
    
    // 使用页面
    html := page.MustHTML()
    println(html)
}
```

### 示例 3：多个独立实例

```go
package main

import (
    "github.com/karosown/katool-go/web_crawler/core"
)

func main() {
    // 创建多个独立的浏览器实例（不使用全局池）
    chrome1 := core.NewContainWithoutPool("C:\\path\\to\\chrome.exe", true)
    chrome2 := core.NewContainWithoutPool("C:\\path\\to\\chrome.exe", true)
    defer chrome1.Close()
    defer chrome2.Close()
    
    // 每个实例完全独立，互不影响
    page1 := chrome1.PageWithStealth("https://example1.com")
    page2 := chrome2.PageWithStealth("https://example2.com")
    defer page1.Close()
    defer page2.Close()
    
    // 使用页面...
}
```

## API 参考

### 构造函数

- `NewCotain(path string, headless bool) *Contain`
  - 创建使用全局池的浏览器实例

- `NewContainWithoutPool(path string, headless bool) *Contain`
  - 创建不使用全局池的浏览器实例（独立管理）

- `NewContainWithPool(path string, headless bool, useGlobalPool bool) *Contain`
  - 创建浏览器实例，可选择是否使用全局池

### 方法

- `PageWithStealth(url string) *rod.Page`
  - 创建新页面（带 Stealth 模式）

- `Close()`
  - 关闭浏览器实例
  - 全局池模式：kill 后创建新浏览器放入池中
  - 非全局池模式：直接 kill，不创建新的

- `ReStart()`
  - 重启浏览器（保持相同的池配置）

## 注意事项

1. **页面必须关闭**：使用 `defer page.Close()` 确保页面被正确关闭
2. **浏览器生命周期**：
   - 全局池模式：浏览器会被复用，`Close()` 后可能被其他实例使用
   - 非全局池模式：浏览器完全独立，`Close()` 后完全销毁
3. **并发安全**：所有操作都是线程安全的
4. **资源管理**：确保在不再使用时调用 `Close()` 释放资源
