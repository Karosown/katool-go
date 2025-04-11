# Lock 锁支持

Lock 模块提供了类似 Java 中的锁机制，用于协程同步和并发控制。

## LockSupport

`LockSupport` 是一个类似 Java 中 `LockSupport` 的实现，可以用于阻塞和恢复协程。

### 基本用法

```go
// 创建一个 LockSupport 实例
lockSupport := lock.NewLockSupport()

// 在某个协程中
go func() {
    // 阻塞当前协程，等待被唤醒
    lockSupport.Park()
    
    // 被唤醒后继续执行
    fmt.Println("协程被唤醒")
}()

// 在另一个协程中唤醒被阻塞的协程
go func() {
    time.Sleep(1 * time.Second) // 模拟延迟
    
    // 恢复被阻塞的协程
    err := lockSupport.Unpark()
    if err != nil {
        fmt.Println("唤醒协程失败:", err)
    }
}()
```

### 注意事项

- 当前版本不支持自动恢复协程，需要显式调用 `Unpark` 方法
- `Park` 方法返回一个布尔值，指示是否成功阻塞
- `Unpark` 方法返回一个错误，如果唤醒失败，会返回相应的错误

## Synchronized

`Synchronized` 方法提供了类似 Java 中 synchronized 关键字的功能，用于确保代码块在锁的保护下执行。

### 基本用法

```go
// 创建一个互斥锁
mutex := &sync.Mutex{}

// 在锁的保护下执行代码
lock.Synchronized(mutex, func() {
    // 这里的代码会在锁的保护下执行
    // 只有一个协程能同时执行这个代码块
    fmt.Println("在锁的保护下执行")
})
```

## LockMap

`LockMap` 是一个字符串到锁的映射，可以用于管理多个命名锁。

### 基本用法

```go
// 创建一个命名锁映射
lockMap := lock.LockMap{}

// 获取或创建特定名称的锁
lockName := "resourceA"
var locker sync.Locker
var ok bool

if locker, ok = lockMap[lockName]; !ok {
    locker = &sync.Mutex{}
    lockMap[lockName] = locker
}

// 使用这个命名锁
lock.Synchronized(locker, func() {
    // 对资源A的访问受到保护
})
```

## 应用场景

1. **协程同步**：使用 `LockSupport` 在协程之间进行同步，一个协程等待另一个协程的信号。

2. **资源保护**：使用 `Synchronized` 保护共享资源，确保在任何时刻只有一个协程能访问资源。

3. **细粒度锁**：使用 `LockMap` 为不同的资源创建独立的锁，减少锁竞争，提高并发性能。

## 示例：实现生产者-消费者模式

```go
func Example_ProducerConsumer() {
    queue := make([]string, 0, 10)
    queueMutex := &sync.Mutex{}
    
    // 创建信号用的 LockSupport
    signal := lock.NewLockSupport()
    
    // 消费者协程
    go func() {
        for {
            // 等待生产者的信号
            signal.Park()
            
            var item string
            // 从队列中取出一个项目
            lock.Synchronized(queueMutex, func() {
                if len(queue) > 0 {
                    item = queue[0]
                    queue = queue[1:]
                }
            })
            
            if item != "" {
                fmt.Println("消费:", item)
            }
        }
    }()
    
    // 生产者
    for i := 1; i <= 5; i++ {
        item := fmt.Sprintf("项目-%d", i)
        
        // 向队列添加一个项目
        lock.Synchronized(queueMutex, func() {
            queue = append(queue, item)
            fmt.Println("生产:", item)
        })
        
        // 通知消费者
        signal.Unpark()
        
        // 生产速度控制
        time.Sleep(500 * time.Millisecond)
    }
    
    // 等待消费完成
    time.Sleep(1 * time.Second)
}
``` 