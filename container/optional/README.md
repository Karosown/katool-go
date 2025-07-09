# Optional 可选值容器

Optional 是一个用于安全处理可能为空值的容器类型，灵感来自 Java 的 Optional 类。

## 基础用法

### 创建 Optional

```go
import "github.com/karosown/katool-go/container/optional"

// 创建包含值的Optional
opt := optional.Of("Hello World")

// 创建空的Optional
emptyOpt := optional.Empty[string]()

// 根据值是否为零值创建Optional
nullableOpt := optional.OfNullable("")  // 空字符串会创建空Optional
```

### 检查和获取值

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
```

### 提供默认值

```go
// 直接提供默认值
defaultValue := emptyOpt.OrElse("默认值")
fmt.Println("结果:", defaultValue) // 输出: 结果: 默认值

// 通过函数提供默认值
lazyDefault := emptyOpt.OrElseGet(func() string {
    return "延迟计算的默认值"
})

// 为空时panic
safeValue := opt.OrElsePanic("Optional不能为空!")
```

## 函数式操作

### 条件执行

```go
// 如果有值则执行函数
opt.IfPresent(func(v string) {
    fmt.Println("处理值:", v)
})

// 有值执行第一个函数，无值执行第二个函数
opt.IfPresentOrElse(
    func(v string) { fmt.Println("有值:", v) },
    func() { fmt.Println("无值") },
)
```

### 过滤

```go
// 过滤满足条件的值
filtered := opt.Filter(func(s string) bool {
    return len(s) > 5
})
```

### 映射转换

```go
// 类型安全的映射（推荐用于链式操作）
result := optional.MapTyped(optional.Of("  hello  "), strings.TrimSpace).
    Filter(func(s string) bool { return len(s) > 0 }).
    OrElse("空字符串")

// 实例方法映射（返回any类型）
mapped := opt.Map(func(s any) any {
    if str, ok := s.(string); ok {
        return strings.ToUpper(str)
    }
    return s
})
```

## 字符串处理专用类型

为了更好地支持字符串处理，我们提供了专用的 StringOptional：

```go
import "strings"

// 使用StringOptional进行链式字符串处理
result := optional.NewStringOptional("  hello  ").
    TrimSpace().                    // 去除空格
    FilterNonEmpty().              // 过滤空字符串
    OrElse("空字符串")             // 提供默认值

fmt.Println("处理结果:", result) // 输出: 处理结果: hello
```

## 工具函数

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

## 实用示例

### 用户输入处理

```go
func processUserInput(input string) string {
    return optional.MapTyped(optional.Of(input), strings.TrimSpace).
        Filter(func(s string) bool { return len(s) > 0 }).
        Map(func(s any) any { return strings.ToLower(s.(string)) }).
        OrElse("无效输入").(string)
}
```

### 配置值处理

```go
func getConfig(key string) optional.Optional[string] {
    if value := os.Getenv(key); value != "" {
        return optional.Of(value)
    }
    return optional.Empty[string]()
}

// 使用
dbUrl := getConfig("DATABASE_URL").
    OrElse("sqlite://default.db")
```

### 链式验证

```go
func validateUser(user User) optional.Optional[User] {
    return optional.Of(user).
        Filter(func(u User) bool { return u.Name != "" }).
        Filter(func(u User) bool { return u.Age >= 18 }).
        Filter(func(u User) bool { return u.Email != "" })
}

// 使用
validUser := validateUser(user).
    OrElsePanic("用户验证失败")
```

## 注意事项

1. **类型安全**: 使用 `MapTyped` 进行类型安全的映射操作
2. **链式调用**: 实例方法支持链式调用，但要注意类型转换
3. **性能**: Optional 会带来轻微的性能开销，在性能敏感的场景中谨慎使用
4. **空指针**: Optional 本身不会为 nil，但内部值可能是零值

## API 列表

### 核心方法
- `Of[T](value T)` - 创建包含值的Optional
- `Empty[T]()` - 创建空Optional
- `OfNullable[T](value T)` - 根据零值创建Optional

### 检查方法
- `IsPresent()` - 检查是否有值
- `IsEmpty()` - 检查是否为空

### 获取方法
- `Get()` - 获取值（空时panic）
- `OrElse(T)` - 提供默认值
- `OrElseGet(func() T)` - 延迟计算默认值
- `OrElsePanic(string)` - 空时panic并显示消息

### 函数式方法
- `IfPresent(func(T))` - 条件执行
- `IfPresentOrElse(func(T), func())` - 双分支执行
- `Filter(func(T) bool)` - 过滤
- `Map(func(T) any)` - 映射（实例方法）
- `MapTyped[T,R](Optional[T], func(T) R)` - 类型安全映射

### 字符串专用
- `NewStringOptional(string)` - 创建字符串Optional
- `TrimSpace()` - 去除空格
- `FilterNonEmpty()` - 过滤空字符串 