# RuleEngine 规则引擎

一个强大的Go规则引擎，支持泛型、并发安全、中间件和树形规则组织。

## 📦 安装

```bash
go get github.com/karosown/katool-go/ruleengine
```

## 🚀 快速开始

### 1. 创建规则引擎

```go
import "github.com/karosown/katool-go/ruleengine"

// 创建规则引擎
engine := ruleengine.NewRuleEngine[int]()

// 用户数据类型
type User struct {
    ID       int
    Name     string
    Age      int
    VipLevel int
    Balance  float64
}

userEngine := ruleengine.NewRuleEngine[User]()
```

### 2. 注册规则

```go
// 注册简单规则
engine.RegisterRule("add_ten",
    func(data int, _ any) bool { return data > 0 },  // 验证函数
    func(data int, _ any) (int, any, error) {        // 执行函数
        return data + 10, "添加了10", nil
    },
)

// 注册用户规则
userEngine.RegisterRule("validate_age",
    func(user User, _ any) bool { return user.Age > 0 },
    func(user User, _ any) (User, any, error) {
        if user.Age < 18 {
            return user, "未成年用户", nil
        }
        return user, "成年用户", nil
    },
)
```

### 3. 构建规则链

```go
// 构建规则链
_, err := engine.NewBuilder("simple_chain").
    AddRule("add_ten").
    Build()

// 构建复杂规则链
_, err := userEngine.NewBuilder("user_processing").
    AddRule("validate_age").
    AddCustomRule(
        func(user User, _ any) bool { return user.Age >= 18 },
        func(user User, _ any) (User, any, error) {
            return user, "成年验证通过", nil
        },
    ).
    Build()
```

### 4. 执行规则

```go
// 执行规则链
result := engine.Execute("simple_chain", 5)
if result.Error != nil {
    fmt.Printf("执行失败: %v\n", result.Error)
} else {
    fmt.Printf("结果: %d, 信息: %v\n", result.Data, result.Result)
}

// 批量执行
results := engine.BatchExecute([]string{"chain1", "chain2"}, data)
```

## 🛠 高级功能

### 中间件

```go
// 添加日志中间件
engine.AddMiddleware(func(data int, next func(int) (int, any, error)) (int, any, error) {
    fmt.Printf("执行前: %d\n", data)
    result, info, err := next(data)
    fmt.Printf("执行后: %d\n", result)
    return result, info, err
})
```

### 错误控制

```go
// 使用EOF终止执行
return data, "提前结束", ruleengine.EOF

// 使用FALLTHROUGH跳过继续
return data, "跳过继续", ruleengine.FALLTHROUGH
```

#### EOF 机制可视化 - 立即终止执行

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

#### FALLTHROUGH 机制可视化 - 跳过继续执行

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

#### 复杂场景示例

```go
// 示例：用户验证规则链
engine.RegisterRule("check_age",
    func(user User, _ any) bool { return true },
    func(user User, _ any) (User, any, error) {
        if user.Age < 13 {
            return user, "用户过小", ruleengine.EOF // 立即终止，不再处理
        } else if user.Age < 18 {
            return user, "未成年用户", ruleengine.FALLTHROUGH // 跳过成年用户逻辑
        }
        return user, "成年用户", nil
    },
)

engine.RegisterRule("adult_verification",
    func(user User, _ any) bool { return user.Age >= 18 },
    func(user User, _ any) (User, any, error) {
        // 这个规则在 FALLTHROUGH 时会被跳过
        return user, "成年验证完成", nil
    },
)

engine.RegisterRule("final_process",
    func(user User, _ any) bool { return true },
    func(user User, _ any) (User, any, error) {
        return user, "最终处理", nil
    },
)

// 执行结果分析：
用户年龄 12: check_age(EOF) → 立即终止，后续规则都不执行
用户年龄 16: check_age(FALLTHROUGH) → adult_verification(跳过) → final_process(执行)
用户年龄 25: check_age(正常) → adult_verification(执行) → final_process(执行)
```

#### 树形结构中的混合场景

```
复杂规则树示例：
                        根节点
                           │
                    ┌──────▼──────┐
                    │  年龄检查    │
                    │ (可能返回EOF) │
                    └──────┬──────┘
                           │
               ┌───────────┼───────────┐
               ▼           ▼           ▼
           ┌──────┐    ┌──────┐    ┌──────┐
           │身份验证│    │邮箱验证│    │手机验证│
           │(正常) │    │(FALL) │    │(正常) │
           └──┬───┘    └──┬───┘    └──┬───┘
              │           │⚡         │
              ▼           ▼           ▼
           ┌──────┐    ┌──────┐    ┌──────┐
           │权限检查│    │邮箱配置│    │短信配置│
           │ ✅   │    │(跳过) │    │ ✅   │
           └──────┘    └──────┘    └──────┘

执行流程说明：
1. 根节点年龄检查：
   - 如果返回 EOF → 整个树立即终止
   - 如果正常 → 继续执行三个分支

2. 三个并行分支：
   - 身份验证 → 权限检查 (正常执行)
   - 邮箱验证 → 返回FALLTHROUGH → 邮箱配置被跳过
   - 手机验证 → 短信配置 (正常执行)

3. 最终结果：除了邮箱配置，其他都正常执行
```

## 🌳 规则树使用

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

## 📝 完整示例

```go
func main() {
    engine := ruleengine.NewRuleEngine[User]()

    // 注册规则
    engine.RegisterRule("validate_age",
        func(user User, _ any) bool { return user.Age > 0 },
        func(user User, _ any) (User, any, error) {
            if user.Age < 18 {
                return user, "未成年", nil
            }
            return user, "成年", nil
        },
    )

    // 构建规则链
    engine.NewBuilder("user_chain").AddRule("validate_age").Build()

    // 执行
    user := User{Name: "张三", Age: 25}
    result := engine.Execute("user_chain", user)
    fmt.Printf("结果: %v\n", result.Result)
}
```

## 🔧 注意事项

- 需要 Go 1.18+ (泛型支持)
- 线程安全，支持并发执行
- `EOF`: 终止规则链
- `FALLTHROUGH`: 跳过当前规则
- 规则命名建议: `动词_名词` 格式

## 📚 API 参考

### 主要方法

- `NewRuleEngine[T]()` - 创建引擎
- `RegisterRule(name, valid, exec)` - 注册规则
- `NewBuilder(name)` - 创建构建器
- `Execute(chain, data)` - 执行规则链
- `BatchExecute(chains, data)` - 批量执行
- `AddMiddleware(middleware)` - 添加中间件