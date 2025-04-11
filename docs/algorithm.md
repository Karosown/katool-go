# Algorithm 算法工具

Algorithm 模块提供了常用的算法工具，包括哈希计算函数、二进制操作等。这些工具可用于数据处理、去重、排序等场景。

## 哈希计算

哈希计算功能主要用于生成对象的哈希值，可用于比较、去重等场景。

### 基本类型和函数

```go
// 哈希值类型
type HashType string

// ID 类型（用于排序）
type IDType int64

// 哈希计算函数类型
type HashComputeFunction func(any2 any) HashType

// ID 计算函数类型
type IDComputeFunction func(any2 any) IDType
```

### HASH_WITH_JSON

将对象序列化为 JSON 格式并作为哈希值返回：

```go
// 基本使用
obj := SomeStruct{Field1: "value", Field2: 123}
hash := algorithm.HASH_WITH_JSON(obj)

// 在 Stream 中使用（用于去重）
uniqueItems := stream.ToStream(&items).DistinctBy(algorithm.HASH_WITH_JSON).ToList()
```

### HASH_WITH_JSON_MD5

将对象序列化为 JSON 后计算 MD5 哈希值：

```go
// 计算对象的 MD5 哈希
obj := SomeStruct{Field1: "value", Field2: 123}
md5Hash := algorithm.HASH_WITH_JSON_MD5(obj)

// 使用 MD5 哈希进行去重（更高效）
uniqueItems := stream.ToStream(&items).DistinctBy(algorithm.HASH_WITH_JSON_MD5).ToList()
```

### HASH_WITH_JSON_SUM

将对象序列化为 JSON 后计算哈希和：

```go
// 计算哈希和
obj := SomeStruct{Field1: "value", Field2: 123}
sumHash := algorithm.HASH_WITH_JSON_SUM(obj)

// 在需要更快计算的场景使用
uniqueItems := stream.ToStream(&items).DistinctBy(algorithm.HASH_WITH_JSON_SUM).ToList()
```

## 二进制操作

`bin.go` 模块提供了与二进制操作相关的工具。

```go
// 获取两个数之间的最大值
max := algorithm.Max(10, 20) // 返回 20

// 获取两个数之间的最小值
min := algorithm.Min(10, 20) // 返回 10

// 计算适合并行处理的线程数（2的幂次）
threads := algorithm.NumOfTwoMultiply(dataSize)
```

## 数组操作

`array.go` 模块提供了数组和切片相关的操作工具。

```go
// 示例数组操作
array := []int{1, 2, 3, 4, 5}

// 数组转换
converted := algorithm.ArrayConvert(array, func(i int) string {
    return fmt.Sprintf("Item-%d", i)
})
// 结果: []string{"Item-1", "Item-2", "Item-3", "Item-4", "Item-5"}

// 数组过滤
filtered := algorithm.ArrayFilter(array, func(i int) bool {
    return i > 3
})
// 结果: []int{4, 5}
```

## 在 Stream 中使用算法工具

算法工具与 Stream 模块紧密集成，特别是在排序和去重操作中：

```go
// 使用 HASH_WITH_JSON 进行去重
unique := stream.ToStream(&items).Distinct().ToList() // 默认使用 HASH_WITH_JSON

// 使用自定义哈希函数进行去重
unique := stream.ToStream(&items).DistinctBy(func(item any) algorithm.HashType {
    return algorithm.HashType(item.(User).ID)
}).ToList()

// 使用哈希函数进行排序
sorted := stream.ToStream(&users).OrderBy(false, func(user any) algorithm.HashType {
    return algorithm.HashType(fmt.Sprintf("%d", user.(User).Age))
}).ToList()

// 使用 ID 计算函数进行排序
sorted := stream.ToStream(&users).OrderById(false, func(user any) algorithm.IDType {
    return algorithm.IDType(user.(User).Age)
}).ToList()
```

## 自定义哈希函数

如果内置的哈希函数不满足需求，可以自定义哈希函数：

```go
// 自定义哈希函数
customHash := func(obj any) algorithm.HashType {
    user := obj.(User)
    // 使用用户名和邮箱的组合作为哈希值
    return algorithm.HashType(user.Username + ":" + user.Email)
}

// 在 Stream 中使用自定义哈希函数
uniqueUsers := stream.ToStream(&users).DistinctBy(customHash).ToList()
```

## 最佳实践

1. **选择合适的哈希函数**：
   - `HASH_WITH_JSON`: 适用于简单对象和调试，可读性好
   - `HASH_WITH_JSON_MD5`: 适用于需要固定长度哈希的场景，碰撞几率低
   - `HASH_WITH_JSON_SUM`: 计算速度快，适用于大数据量场景

2. **避免重复计算哈希值**：
   ```go
   // 不好的方式：每次比较都会重新计算哈希
   items = stream.ToStream(&items).
       Filter(func(item any) bool { return algorithm.HASH_WITH_JSON_SUM(item) != "0" }).
       DistinctBy(algorithm.HASH_WITH_JSON).
       ToList()
   
   // 更好的方式：预计算哈希
   itemWithHash := stream.ToStream(&items).Map(func(item any) any {
       return struct {
           Item any
           Hash algorithm.HashType
       }{
           Item: item,
           Hash: algorithm.HASH_WITH_JSON(item),
       }
   }).ToList()
   ```

3. **并行处理**：
   - 使用 `NumOfTwoMultiply` 计算合适的并行度
   - 结合 Stream 并行流提高性能 