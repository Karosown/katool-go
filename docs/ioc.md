# IOC 容器

IOC (Inversion of Control) 容器是一个轻量级依赖注入解决方案，类似于 Java Spring 框架中的依赖注入功能。它帮助管理对象实例、解决依赖关系，并使代码更易于测试和维护。

## 基本用法

### 获取实例

```go
// 获取已注册的实例，键名为"userService"
userService := ioc.Get("userService")

// 使用类型断言
userService := ioc.Get("userService").(*UserService)

// 获取带类型的实例
userService := ioc.GetDef("userService", &UserService{}).(*UserService)
```

### 注册值

直接注册一个值：

```go
// 注册一个值，如果键名已存在会panic
ioc.RegisterValue("config", myConfig)

// 强制注册值，会覆盖已存在的键
ioc.MustRegisterValue("config", myConfig)
```

### 使用工厂函数注册

注册一个工厂函数，只有在首次获取时才会执行：

```go
// 注册一个工厂函数，如果键名已存在会panic
ioc.Register("database", func() any {
    // 创建数据库连接
    db, _ := sql.Open("postgres", "connection-string")
    return db
})

// 强制注册，会覆盖已存在的键
ioc.MustRegister("database", func() *sql.DB {
    db, _ := sql.Open("postgres", "connection-string")
    return db
})

// 强制替换现有注册
ioc.ForceRegister("database", func() *sql.DB {
    db, _ := sql.Open("postgres", "connection-string")
    return db
})
```

### 带默认值的获取

获取一个值，如果不存在则使用默认值：

```go
// 如果"maxConnections"不存在，则注册并返回10
maxConn := ioc.GetDef("maxConnections", 10).(int)

// 如果"database"不存在，则执行函数创建并注册
db := ioc.GetDefFunc("database", func() *sql.DB {
    db, _ := sql.Open("postgres", "connection-string")
    return db
})
```

## 应用场景

1. **单例模式**: 确保全局只有一个实例
   ```go
   // 在程序启动时注册
   ioc.RegisterValue("logger", log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime))
   
   // 在需要的地方获取
   logger := ioc.Get("logger").(*log.Logger)
   ```

2. **依赖注入**: 集中管理组件依赖
   ```go
   // 注册服务
   ioc.Register("userRepository", func() any {
       return &UserRepository{DB: ioc.Get("database").(*sql.DB)}
   })
   
   ioc.Register("userService", func() any {
       return &UserService{Repo: ioc.Get("userRepository").(*UserRepository)}
   })
   
   // 获取服务
   userService := ioc.Get("userService").(*UserService)
   ```

3. **配置管理**: 统一管理应用配置
   ```go
   // 注册配置
   ioc.RegisterValue("appConfig", &AppConfig{
       Port:     8080,
       LogLevel: "info",
   })
   
   // 在需要的地方获取
   config := ioc.Get("appConfig").(*AppConfig)
   ```

4. **测试中模拟依赖**: 
   ```go
   // 在测试前替换真实实现
   ioc.ForceRegister("emailService", func() any {
       return &MockEmailService{}
   })
   ```

## 最佳实践

- 使用描述性的键名，如 "userService" 而不是 "service1"
- 对于复杂对象，使用工厂函数而不是直接注册值
- 在应用启动时集中注册所有依赖
- 避免循环依赖
- 在测试中使用 ForceRegister 替换实现 