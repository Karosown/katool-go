# Convert 转换工具

Convert 模块提供了各种数据转换功能，帮助在不同类型之间进行转换、复制属性、处理通道数据等。

## 基本类型转换

### ToString - 将任意类型转为字符串

支持将基本类型、结构体、切片等转换为字符串表示：

```go
// 基本类型转字符串
intStr := convert.ToString(123)         // "123"
floatStr := convert.ToString(123.45)    // "123.45"
boolStr := convert.ToString(true)       // "true"

// 时间类型转字符串 
timeStr := convert.ToString(time.Now()) // 格式化的时间字符串

// 结构体、map等复杂类型转为JSON字符串
type User struct {
    Name string
    Age  int
}
userStr := convert.ToString(User{Name: "张三", Age: 30}) // 转为JSON字符串
```

## 对象属性复制

### CopyProperties - 在两个对象之间复制同名属性

```go
// 源对象
type UserEntity struct {
    ID       int64
    Username string
    Password string
    Email    string
    Age      int
}

// 目标对象
type UserDTO struct {
    ID       int64
    Username string
    Email    string
    Age      int
}

// 复制属性（只会复制同名同类型的字段）
entity := UserEntity{1, "admin", "secret", "admin@example.com", 30}
dto := &UserDTO{}
result, err := convert.CopyProperties(entity, dto)

if err != nil {
    fmt.Println("属性复制失败:", err)
} else {
    fmt.Printf("复制后的对象: %+v\n", result)
    // 输出: 复制后的对象: &{ID:1 Username:admin Email:admin@example.com Age:30}
}
```

## 切片转换

### Convert - 对切片中的每个元素应用转换函数

```go
// 源切片
users := []UserEntity{
    {1, "user1", "pwd1", "user1@example.com", 25},
    {2, "user2", "pwd2", "user2@example.com", 30},
}

// 使用转换函数将每个元素转换为新类型
dtos := convert.Convert(users, func(user UserEntity) UserDTO {
    dto := UserDTO{}
    result, _ := convert.CopyProperties(user, &dto)
    return *result.(*UserDTO)
})

// dtos 现在包含转换后的 UserDTO 对象切片
```

### ToAnySlice - 将类型化切片转为 []any

```go
// 类型化切片
numbers := []int{1, 2, 3, 4, 5}

// 转换为 []any
anySlice := convert.ToAnySlice(numbers)

// anySlice 现在是 []any{1, 2, 3, 4, 5}
```

### FromAnySlice - 将 []any 转为类型化切片

```go
// any 切片
anySlice := []any{1, 2, 3, 4, 5}

// 转换为类型化切片
intSlice := convert.FromAnySlice[int](anySlice)

// intSlice 现在是 []int{1, 2, 3, 4, 5}
```

## 通道数据处理

### ChanToArray - 将通道中的数据转为切片

```go
// 创建一个通道并发送数据
ch := make(chan string, 3)
ch <- "one"
ch <- "two"
ch <- "three"
close(ch)

// 将通道数据转为切片
strArray := convert.ChanToArray(ch)

// strArray 现在是 []string{"one", "two", "three"}
```

### ChanToFlatArray - 将通道中的切片数据平铺为单一切片

```go
// 创建一个通道，其中包含多个切片
ch := make(chan []int, 3)
ch <- []int{1, 2}
ch <- []int{3, 4}
ch <- []int{5, 6}
close(ch)

// 将通道中的切片平铺为单一切片
flatArray := convert.ChanToFlatArray(ch)

// flatArray 现在是 []int{1, 2, 3, 4, 5, 6}
```

## 批次处理转换

结合 lists.Partition 和 Stream 进行批量处理：

```go
// 假设 userList 是一个大型用户列表
// 将其分割成批次，然后转换为流进行处理
result := convert.PatitonToStreamp(lists.Partition(userList[:], 100)).
    Parallel().
    Map(func(batch []User) any {
        // 处理每一批用户
        processedBatch := make([]ProcessedUser, 0, len(batch))
        for _, user := range batch {
            processed := ProcessUser(user)
            processedBatch = append(processedBatch, processed)
        }
        return processedBatch
    }).
    ToList()
```

## 最佳实践

- 使用 `CopyProperties` 替代手动赋值，特别是在 DTO 和实体对象之间转换时
- 对于大型数据集，考虑使用批处理和并行流组合提高性能
- 处理未知类型时，使用 `ToString` 进行安全的字符串转换 