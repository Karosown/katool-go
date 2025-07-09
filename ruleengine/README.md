# 规则引擎 (Rule Engine)

一个基于Go泛型的灵活、高性能规则引擎，支持复杂的业务规则管理和执行。

## 特性

- 🚀 **基于泛型**：类型安全，支持任意数据类型
- 🏗️ **构建器模式**：链式API，易于使用
- ⚡ **高性能**：基于队列的异步执行机制
- 🔗 **规则链**：支持复杂的规则组合和串联
- 🔧 **中间件**：支持执行前后的数据处理
- 🔄 **并发安全**：内置读写锁，支持并发访问
- 📊 **批量执行**：支持并发批量处理
- 🎯 **条件分支**：支持复杂的条件逻辑
- 📈 **动态管理**：运行时添加/删除规则

## 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "yourproject/ruleengine"
)

func main() {
    // 创建规则引擎
    engine := ruleengine.NewRuleEngine[int]()

    // 注册规则
    engine.RegisterRule("add_ten",
        func(data int, _ any) bool { return data > 0 }, // 验证函数
        func(data int, _ any) (int, any, error) {       // 执行函数
            return data + 10, "添加了10", nil
        },
    )

    // 构建规则链
    _, err := engine.NewBuilder("simple_chain").
        AddRule("add_ten").
        Build()
    if err != nil {
        panic(err)
    }

    // 执行规则
    result := engine.Execute("simple_chain", 5)
    fmt.Printf("结果: %d, 信息: %v\n", result.Data, result.Result)
    // 输出: 结果: 15, 信息: 添加了10
}
```

### 复杂示例

```go
// 用户数据处理示例
type User struct {
    ID       int
    Name     string
    Age      int
    VipLevel int
    Balance  float64
}

func main() {
    engine := ruleengine.NewRuleEngine[User]()

    // 注册年龄验证规则
    engine.RegisterRule("validate_age",
        func(user User, _ any) bool { return user.Age > 0 },
        func(user User, _ any) (User, any, error) {
            if user.Age >= 18 {
                return user, "成年用户", nil
            }
            return user, "未成年用户", nil
        },
    )

    // 注册VIP折扣计算规则
    engine.RegisterRule("calculate_discount",
        func(user User, _ any) bool { return user.VipLevel > 0 },
        func(user User, _ any) (User, any, error) {
            discount := float64(user.VipLevel) * 0.1
            return user, fmt.Sprintf("VIP折扣: %.0f%%", discount*100), nil
        },
    )

    // 构建用户处理链
    engine.NewBuilder("user_processing").
        AddRule("validate_age").
        AddCustomRule(
            func(user User, _ any) bool { return user.Age >= 18 },
            func(user User, _ any) (User, any, error) {
                return user, "成年验证通过", nil
            },
        ).
        AddRule("calculate_discount").
        Build()

    // 执行处理
    user := User{ID: 1, Name: "张三", Age: 25, VipLevel: 2, Balance: 1000}
    result := engine.Execute("user_processing", user)
    
    fmt.Printf("处理结果: %+v\n", result.Data)
}
```

## 核心概念

### 1. 规则节点 (RuleNode)

规则节点是规则引擎的基础单元，包含：
- **验证函数** (`Valid`): 决定规则是否应该执行
- **执行函数** (`Exec`): 包含具体的业务逻辑

### 2. 规则链 (RuleChain)

规则链是多个规则节点的有序组合，支持：
- 顺序执行
- 条件分支
- 错误处理

### 3. 规则引擎 (RuleEngine)

规则引擎是管理中心，提供：
- 规则注册和查找
- 规则链构建和执行
- 中间件管理
- 并发控制

## API 参考

### 规则引擎管理

```go
// 创建新引擎
engine := NewRuleEngine[T]()

// 注册规则
engine.RegisterRule(name, validFunc, execFunc)

// 获取规则
rule, exists := engine.GetRule(name)

// 移除规则
removed := engine.RemoveRule(name)

// 列出所有规则
rules := engine.ListRules()

// 获取统计信息
stats := engine.Stats()
```

### 规则链构建

```go
// 创建构建器
builder := engine.NewBuilder("chain_name")

// 添加已注册的规则
builder.AddRule("rule_name")

// 添加自定义规则
builder.AddCustomRule(validFunc, execFunc)

// 添加条件分支
builder.AddConditionalChain(condition, trueChain, falseChain)

// 构建规则链
tree, err := builder.Build()
```

### 规则执行

```go
// 执行单个规则链
result := engine.Execute("chain_name", data)

// 批量执行多个规则链（并发）
results := engine.BatchExecute(chainNames, data)

// 执行所有规则链
allResults := engine.ExecuteAll(data)
```

### 中间件

```go
// 添加中间件
engine.AddMiddleware(func(data T, next func(T) (T, any, error)) (T, any, error) {
    // 前置处理
    fmt.Println("开始处理")
    
    // 调用下一个处理器
    result, info, err := next(data)
    
    // 后置处理
    fmt.Println("处理完成")
    
    return result, info, err
})
```

## 高级功能

### 条件分支

```go
// 构建条件分支规则
trueChain := []*RuleNode[T]{
    // 条件为真时执行的规则
}

falseChain := []*RuleNode[T]{
    // 条件为假时执行的规则
}

builder.AddConditionalChain(
    func(data T, _ any) bool {
        // 条件判断逻辑
        return someCondition
    },
    trueChain,
    falseChain,
)
```

### 错误处理

规则引擎提供两种特殊错误：

```go
// 结束执行
return data, result, ruleengine.EOF

// 继续执行下一个规则
return data, result, ruleengine.FALLTHROUGH
```

### 并发执行

```go
// 并发执行多个规则链
chainNames := []string{"chain1", "chain2", "chain3"}
results := engine.BatchExecute(chainNames, data)

// 处理结果
for chainName, result := range results {
    if result.Error != nil {
        fmt.Printf("规则链 %s 执行失败: %v\n", chainName, result.Error)
    } else {
        fmt.Printf("规则链 %s 执行成功: %v\n", chainName, result.Result)
    }
}
```

## 测试

```bash
# 运行所有测试
go test ./ruleengine

# 运行性能测试
go test -bench=. ./ruleengine

# 查看测试覆盖率
go test -cover ./ruleengine

# 运行示例
go test -run TestExample ./ruleengine
```

## 性能特性

- **内存高效**: 基于队列的执行机制，避免递归栈溢出
- **并发安全**: 内置读写锁，支持高并发访问
- **异步执行**: 支持非阻塞的规则执行
- **批量处理**: 并发执行多个规则链，提高吞吐量

## 使用场景

### 1. 业务规则引擎
- 订单处理流程
- 用户权限验证
- 价格计算规则

### 2. 数据处理管道
- 数据验证和转换
- ETL流程
- 数据清洗

### 3. 工作流引擎
- 审批流程
- 状态机实现
- 业务流程自动化

### 4. 配置驱动的业务逻辑
- 动态规则配置
- A/B测试规则
- 特性开关

## 最佳实践

### 1. 规则设计
- 保持规则单一职责
- 合理设计验证函数
- 避免规则间的强耦合

### 2. 性能优化
- 将常用规则放在链的前端
- 合理使用条件验证减少不必要的执行
- 对于大量数据使用批量执行

### 3. 错误处理
- 合理使用 EOF 和 FALLTHROUGH
- 在规则中进行适当的错误检查
- 使用中间件进行统一的错误处理

### 4. 并发安全
- 避免在规则中修改共享状态
- 使用引擎的内置并发控制
- 合理设计数据结构避免竞态条件

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！ 