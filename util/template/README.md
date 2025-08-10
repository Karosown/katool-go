# 模板引擎 (Template Engine)

这是一个通用的模板引擎，支持短信和邮件模板的变量替换和发送功能。

## 功能特性

- 支持自定义分隔符
- 类型安全的泛型设计
- 模板验证功能
- 链式调用API
- 支持批量映射添加

## 基本用法

### 1. 创建短信模板

```go
// 创建短信适配器
smsAdapter := SMSAdapter{
    ID:       "sms_001",
    Name:     "阿里云短信",
    Provider: "aliyun",
    Status:   "active",
    Config: map[string]interface{}{
        "accessKeyId":     "your_access_key",
        "accessKeySecret": "your_secret_key",
        "signName":        "公司名称",
    },
}

// 创建模板引擎
engine := NewEngine[SMSAdapter]("您好，您的验证码是：{{code}}，有效期{{expire}}分钟")

// 添加变量映射
engine.AddMapping("code", "123456")
engine.AddMapping("expire", "5")

// 加载模板
result := engine.Load()
// 输出: "您好，您的验证码是：123456，有效期5分钟"
```

### 2. 创建邮件模板

```go
// 创建邮件适配器
mailAdapter := MailAdapter{
    ID:       "mail_001",
    Name:     "SMTP邮件",
    Provider: "smtp",
    Status:   "active",
    Config: map[string]interface{}{
        "host":     "smtp.example.com",
        "port":     587,
        "username": "user@example.com",
        "password": "password",
    },
}

// 创建模板引擎
engine := NewEngine[MailAdapter]("欢迎 {{name}}！您的账户已激活，登录地址：{{loginUrl}}")

// 批量添加映射
mappings := map[string]string{
    "name":     "张三",
    "loginUrl": "https://example.com/login",
}
engine.AddMappings(mappings)

// 加载模板
result := engine.Load()
// 输出: "欢迎 张三！您的账户已激活，登录地址：https://example.com/login"
```

### 3. 自定义分隔符

```go
engine := NewEngine[SMSAdapter]("Hello ${name}, your code is ${code}")
engine.SetDelimiters("${", "}")

engine.AddMapping("name", "John")
engine.AddMapping("code", "ABC123")

result := engine.Load()
// 输出: "Hello John, your code is ABC123"
```

### 4. 模板验证

```go
engine := NewEngine[SMSAdapter]("Hello {{name}}, your code is {{code}}")
engine.AddMapping("name", "John")
// 故意不添加code的映射

if err := engine.Validate(); err != nil {
    fmt.Println("模板验证失败:", err)
    // 输出: "模板验证失败: template contains unresolved variables"
}
```

### 5. 链式调用

```go
result := NewEngine[SMSAdapter]("Hello {{name}}")
    .AddMapping("name", "World")
    .SetDelimiters("{{", "}}")
    .Load()
// 输出: "Hello World"
```

## 适配器结构

### SMSAdapter

短信适配器包含以下字段：

- `ID`: 适配器唯一标识
- `Name`: 适配器名称
- `Provider`: 短信服务提供商（如：aliyun, tencent等）
- `Status`: 状态（active, inactive, disabled）
- `Config`: 配置信息（API密钥、签名等）

### MailAdapter

邮件适配器包含以下字段：

- `ID`: 适配器唯一标识
- `Name`: 适配器名称
- `Provider`: 邮件服务提供商（如：smtp, sendgrid等）
- `Status`: 状态（active, inactive, disabled）
- `Config`: 配置信息（SMTP服务器、用户名密码等）

## 注意事项

1. 适配器结构体只包含配置信息，不包含具体的客户端实例
2. 模板引擎使用泛型确保类型安全
3. 默认分隔符为 `{{` 和 `}}`
4. 验证功能会检查是否有未解析的模板变量
5. 所有方法都支持链式调用

## 运行测试

```bash
go test ./util/template/test/
``` 