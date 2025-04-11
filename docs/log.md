# 日志工具

Katool 提供了轻量级的日志工具，支持多种日志级别和适配器，满足不同场景下的日志记录需求。

## 核心功能

- 支持多个日志级别：INFO、WARN、ERROR、DEBUG 等
- 支持多种日志输出目标
- 提供 logrus 适配器
- 简洁易用的 API

## 基本用法

### 基本日志记录

```go
// 使用默认日志记录器
log.Info("这是一条信息日志")
log.Warn("这是一条警告日志")
log.Error("这是一条错误日志")
log.Debug("这是一条调试日志")

// 使用格式化的日志消息
log.Infof("用户 %s 登录系统", username)
log.Warnf("CPU 使用率高: %.2f%%", cpuUsage)
log.Errorf("操作失败: %v", err)
```

### 带上下文的日志

```go
// 创建带有上下文字段的日志
log.WithFields(map[string]interface{}{
    "user":      "admin",
    "requestID": "req-123456",
    "ip":        "192.168.1.1",
}).Info("用户登录")

// 链式调用添加上下文
log.WithField("module", "auth").
    WithField("method", "POST").
    Info("处理请求")
```

## Logrus 适配器

Katool 提供了对 [Logrus](https://github.com/sirupsen/logrus) 的适配器，可以无缝集成 Logrus 的强大功能。

### 配置 Logrus 适配器

```go
import (
    "github.com/sirupsen/logrus"
    "github.com/karosown/katool-go/log"
)

func initLogger() {
    // 创建 logrus 实例
    logrusLogger := logrus.New()
    
    // 配置 logrus
    logrusLogger.SetFormatter(&logrus.JSONFormatter{})
    logrusLogger.SetLevel(logrus.InfoLevel)
    
    // 设置文件输出
    file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    logrusLogger.SetOutput(file)
    
    // 创建适配器
    adapter := log.NewLogrusAdapter(logrusLogger)
    
    // 设置为全局日志适配器
    log.SetGlobalLogger(adapter)
}
```

### 使用适配后的日志

```go
// 初始化后的使用方式与标准日志相同
log.Info("使用 Logrus 记录的日志")
log.WithField("component", "database").Error("连接失败")
```

## 自定义日志适配器

你可以通过实现 `Logger` 接口创建自己的日志适配器：

```go
// Logger 接口定义
type Logger interface {
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Warn(args ...interface{})
    Warnf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
}

// 实现自定义适配器
type MyCustomLogger struct {
    // 自定义字段
}

// 实现 Logger 接口的方法
func (l *MyCustomLogger) Info(args ...interface{}) {
    // 自定义实现
}

// ... 实现其他方法 ...

// 设置为全局日志适配器
log.SetGlobalLogger(&MyCustomLogger{})
```

## 日志包装器

`xlog` 包提供了更高级的日志包装功能，支持日志轮转、过滤、格式化等更多功能。

```go
import "github.com/karosown/katool-go/xlog"

// 创建日志包装器
logger := xlog.NewLogWrapper(
    xlog.WithLevel(xlog.InfoLevel),
    xlog.WithFormat(xlog.JSONFormat),
    xlog.WithOutput("app.log"),
    xlog.WithRotation(24 * time.Hour, 7), // 每天轮转，保留 7 个文件
)

// 使用包装器记录日志
logger.Info("应用启动")
logger.WithFields(xlog.Fields{
    "user": "admin",
    "role": "superuser",
}).Info("用户执行管理操作")

// 关闭日志
defer logger.Close()
```

## 最佳实践

1. **合理设置日志级别**：生产环境通常设置为 INFO 或 WARN 级别，测试环境可设置为 DEBUG 级别。

2. **添加上下文信息**：使用 WithField/WithFields 添加上下文信息，便于日志分析。

3. **日志轮转**：对于长期运行的应用，配置日志轮转防止日志文件过大。

4. **结构化日志**：在需要机器处理日志的场景，使用 JSON 格式的结构化日志。

5. **错误日志处理**：记录错误日志时包含完整的错误信息和堆栈跟踪。
   ```go
   if err != nil {
       log.WithField("stack", debug.Stack()).Errorf("操作失败: %v", err)
   }
   ``` 