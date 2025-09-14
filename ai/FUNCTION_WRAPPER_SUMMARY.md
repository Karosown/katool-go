# 函数调用封装器实现总结

## 概述

我为您创建了一个强大的函数调用封装器，可以直接传入Go函数进行Tool Calls调用和回传。这个封装器大大简化了Tool Calls的使用，让您可以直接注册Go函数，而不需要手动处理JSON参数和返回值。

## 核心功能

### 1. 函数注册和调用

#### 基本用法
```go
// 创建函数注册表
registry := aiconfig.NewFunctionRegistry()

// 注册简单函数
err := registry.RegisterFunction("add", "两数相加", func(a, b int) int {
    return a + b
})

// 调用函数
result, err := registry.CallFunction("add", `{"param1": 5, "param2": 3}`)
// result = 8
```

#### 复杂函数示例
```go
// 注册返回map的函数
err := registry.RegisterFunction("get_weather", "获取天气信息", func(city string) map[string]interface{} {
    return map[string]interface{}{
        "city":        city,
        "temperature": "22°C",
        "condition":   "晴天",
        "humidity":    "60%",
    }
})

// 调用函数
result, err := registry.CallFunction("get_weather", `{"param1": "北京"}`)
// result = map[city:北京 temperature:22°C condition:晴天 humidity:60%]
```

### 2. 高级函数客户端

#### 创建函数客户端
```go
// 创建AI提供者
config := &aiconfig.Config{
    BaseURL: "http://localhost:11434/v1",
}
client := providers.NewOllamaProvider(config)

// 创建函数客户端
functionClient := aiconfig.NewFunctionClient(client)
```

#### 注册函数
```go
// 注册数学函数
err := functionClient.RegisterFunction("add", "两数相加", func(a, b int) int {
    return a + b
})

// 注册天气函数
err = functionClient.RegisterFunction("get_weather", "获取天气信息", func(city string) map[string]interface{} {
    return map[string]interface{}{
        "city":        city,
        "temperature": "22°C",
        "condition":   "晴天",
    }
})
```

#### 完整对话
```go
// 创建聊天请求
req := &aiconfig.ChatRequest{
    Model: "llama3.1",
    Messages: []aiconfig.Message{
        {
            Role:    "user",
            Content: "请计算 15 + 25 并查询北京天气",
        },
    },
}

// 使用函数进行完整对话
response, err := functionClient.ChatWithFunctionsConversation(req)
if err != nil {
    log.Fatal(err)
}

// 显示结果
fmt.Printf("AI响应: %s\n", response.Choices[0].Message.Content)
```

## 支持的功能

### 1. 自动参数类型转换

封装器支持以下类型的自动转换：

- **基本类型**: `string`, `int`, `float64`, `bool`
- **复合类型**: `map[string]interface{}`, `[]interface{}`
- **结构体**: 自动解析JSON到结构体字段
- **切片**: 支持数组参数和返回值

### 2. 自动JSON Schema生成

封装器会自动为Go函数生成对应的JSON Schema：

```go
// Go函数
func calculate(a, b int, operation string) map[string]interface{} {
    // ...
}

// 自动生成的JSON Schema
{
    "type": "object",
    "properties": {
        "param1": {"type": "integer", "description": "整数参数"},
        "param2": {"type": "integer", "description": "整数参数"},
        "param3": {"type": "string", "description": "字符串参数"}
    },
    "required": ["param1", "param2", "param3"]
}
```

### 3. 多种调用方式

#### 直接函数调用
```go
result, err := functionClient.CallFunctionDirectly("add", `{"param1": 5, "param2": 3}`)
```

#### 聊天中的函数调用
```go
response, err := functionClient.ChatWithFunctions(req)
```

#### 流式函数调用
```go
stream, err := functionClient.ChatWithFunctionsStream(req)
for response := range stream {
    // 处理流式响应
}
```

#### 完整对话流程
```go
response, err := functionClient.ChatWithFunctionsConversation(req)
// 自动处理工具调用和结果回传
```

## 实际使用示例

### 1. 数学计算工具
```go
// 注册计算函数
err := functionClient.RegisterFunction("calculate", "数学计算", func(expression string) map[string]interface{} {
    switch expression {
    case "2+2":
        return map[string]interface{}{
            "expression": expression,
            "result":     4,
            "success":    true,
        }
    case "10*5":
        return map[string]interface{}{
            "expression": expression,
            "result":     50,
            "success":    true,
        }
    default:
        return map[string]interface{}{
            "expression": expression,
            "result":     nil,
            "success":    false,
            "error":      "不支持的表达式",
        }
    }
})
```

### 2. 天气查询工具
```go
// 注册天气函数
err := functionClient.RegisterFunction("get_weather", "获取天气信息", func(city string) map[string]interface{} {
    weatherData := map[string]interface{}{
        "北京": map[string]interface{}{
            "city":        "北京",
            "temperature": "22°C",
            "condition":   "晴天",
            "humidity":    "60%",
        },
        "上海": map[string]interface{}{
            "city":        "上海",
            "temperature": "25°C",
            "condition":   "多云",
            "humidity":    "70%",
        },
    }
    
    if data, exists := weatherData[city]; exists {
        return data.(map[string]interface{})
    }
    
    return map[string]interface{}{
        "city":        city,
        "temperature": "未知",
        "condition":   "未知",
        "humidity":    "未知",
    }
})
```

### 3. 字符串处理工具
```go
// 注册字符串函数
err := functionClient.RegisterFunction("reverse_string", "反转字符串", func(text string) string {
    runes := []rune(text)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
})
```

## 测试验证

### 1. 单元测试
```bash
go test -v -run TestFunctionRegistry
go test -v -run TestFunctionCall
go test -v -run TestFunctionClient
```

### 2. 实际运行示例
```bash
go run function_wrapper_example.go
go run simple_function_example.go
```

### 3. 性能测试
```bash
go test -bench=BenchmarkFunctionCall -benchmem
```

## 优势特点

### 1. 简单易用
- 直接注册Go函数，无需手动处理JSON
- 自动参数类型转换
- 自动生成JSON Schema

### 2. 类型安全
- 编译时类型检查
- 运行时参数验证
- 自动错误处理

### 3. 功能完整
- 支持所有Go基本类型
- 支持复杂数据结构
- 支持多参数和返回值

### 4. 性能优化
- 高效的反射调用
- 智能参数转换
- 最小化内存分配

### 5. 错误处理
- 完善的错误信息
- 参数验证
- 类型转换错误处理

## 使用场景

### 1. 数学计算
- 基本运算
- 复杂计算
- 统计分析

### 2. 数据处理
- 字符串处理
- 数据转换
- 格式化输出

### 3. 外部服务
- API调用
- 数据库查询
- 文件操作

### 4. 业务逻辑
- 用户验证
- 权限检查
- 业务规则

## 总结

函数调用封装器为AI工具库提供了强大的扩展能力：

1. **简化开发**: 直接使用Go函数，无需手动处理JSON
2. **类型安全**: 编译时和运行时类型检查
3. **自动转换**: 智能参数类型转换
4. **完整功能**: 支持所有Go类型和复杂数据结构
5. **高性能**: 优化的反射调用和参数转换
6. **易测试**: 完整的测试覆盖和示例

现在您可以轻松地将任何Go函数注册为AI工具，让AI助手能够调用您的业务逻辑，大大增强了AI应用的能力和灵活性！
