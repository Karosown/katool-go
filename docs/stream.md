# Stream 流式处理

Stream 模块提供了类似 Java 中的 Stream API 的流式处理能力，可以链式调用各种操作方法进行数据处理。

## 基本用法

### 创建 Stream

```go
// 从切片创建流
stream := stream.ToStream(&mySlice)

// 或使用 Of 方法（与 ToStream 功能相同）
stream := stream.Of(&mySlice)

// 创建并行流
parallelStream := stream.ToParallelStream(&mySlice)
// 或将现有流转为并行流
parallelStream := stream.ToStream(&mySlice).Parallel()
```

### 常用操作

#### Map 转换

将元素转换为其他类型：

```go
result := stream.ToStream(&users).Map(func(user User) any {
    return user.Name
}).ToList()
```

#### Filter 过滤

筛选满足条件的元素：

```go
filtered := stream.ToStream(&numbers).Filter(func(n int) bool {
    return n > 10
}).ToList()
```

#### Reduce 归约

将所有元素组合成单一结果：

```go
sum := stream.ToStream(&numbers).Reduce(0, 
    // 元素组合函数
    func(acc any, n int) any {
        return acc.(int) + n
    }, 
    // 并行结果组合函数（并行模式需要）
    func(result1, result2 any) any {
        return result1.(int) + result2.(int)
    })
```

#### Distinct 去重

去除重复元素：

```go
// 使用默认去重方式（JSON序列化比较）
distinct := stream.ToStream(&numbers).Distinct().ToList()

// 使用自定义去重方式
distinct := stream.ToStream(&users).DistinctBy(func(user any) algorithm.HashType {
    return algorithm.HashType(user.(User).ID)
}).ToList()
```

#### FlatMap 扁平化映射

将每个元素映射为一个流并扁平化：

```go
result := stream.ToStream(&departments).FlatMap(func(dept Department) *stream.Stream[any, []any] {
    return stream.ToStream(&dept.Employees).Map(func(emp Employee) any {
        return emp.Name
    })
}).ToList()
```

#### GroupBy 分组

按条件将元素分组：

```go
grouped := stream.ToStream(&users).GroupBy(func(user User) any {
    return user.Department
})
// 返回 map[any][]User
```

#### 排序

```go
// 使用 OrderBy 排序
sorted := stream.ToStream(&users).OrderBy(false, func(user any) algorithm.HashType {
    return algorithm.HashType(user.(User).Age)
}).ToList()

// 使用 Sort 方法自定义排序
sorted := stream.ToStream(&users).Sort(func(a, b User) bool {
    return a.Age < b.Age
}).ToList()
```

#### ForEach 遍历

对每个元素执行操作：

```go
stream.ToStream(&users).ForEach(func(user User) {
    fmt.Println(user.Name)
})
```

#### Collect 收集

自定义收集结果：

```go
result := stream.ToStream(&users).Collect(func(data stream.Options[User], sourceData []User) any {
    // 自定义收集逻辑
    return customCollect(data)
})
```

#### ToMap 转为映射

```go
userMap := stream.ToStream(&users).ToMap(
    // 键提取函数
    func(index int, user User) any {
        return user.ID
    },
    // 值提取函数
    func(index int, user User) any {
        return user
    },
)
```

#### Count 计数

```go
count := stream.ToStream(&users).Count()
```

## 并行流

将流转换为并行流以利用多核能力：

```go
result := stream.ToStream(&largeDataSet).
    Parallel().       // 转换为并行流
    Filter(func(item int) bool {
        return item > 100
    }).
    Map(func(item int) any {
        return item * 2
    }).
    UnParallel().     // 可选：转回串行流
    ToList()
```

## 注意事项

1. 由于 Go 的泛型限制，使用时需要自行处理类型转换
2. 链式调用中的方法调用顺序会影响处理效率
3. 并行流适合计算密集型的大数据集处理，小数据集可能因为调度开销反而更慢 