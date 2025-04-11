# Lists 集合操作

Lists 模块提供了对集合的各种操作工具，如批处理、分区等功能。

## Batch 批处理

Batch 结构支持将大集合分割成多个小批次进行处理，特别适合需要并行处理大数据集的场景。

### Partition 分片

将集合分割成多个相等大小的子集合：

```go
// 将 userList 分成每组最多 15 个元素的批次
batches := lists.Partition(userList[:], 15)
```

### ForEach 批处理

对分割后的每个批次执行操作，支持并行处理：

```go
// 同步处理每个批次
lists.Partition(userList[:], 15).ForEach(func(pos int, automicDatas []user) error {
    // pos 表示当前是第几个批次
    fmt.Println("处理第", pos, "批数据")
    
    // 处理当前批次的数据
    for _, data := range automicDatas {
        // 处理单个数据
    }
    
    return nil
}, false, nil) // false表示同步处理，不并行

// 并行处理每个批次，限制最大并发为3
lists.Partition(userList[:], 15).ForEach(func(pos int, automicDatas []user) error {
    fmt.Println("分批处理 第" + convert.ToString(pos) + "批")
    
    // 处理当前批次的数据
    stream.ToStream(&automicDatas).ForEach(func(data user) {
        fmt.Println(data)
    })
    
    return nil
}, true, lynx.NewLimiter(3)) // true表示异步处理，最大并发为3
```

### 结合 Stream 处理

可以将分批处理与 Stream 流式处理结合使用，实现更复杂的数据处理逻辑：

```go
// 将批次转换为流，对每个批次应用Map操作，最后使用Reduce合并结果
sum := convert.PatitonToStreamp(lists.Partition(userList[:], 15)).
    Parallel(). // 对批次并行处理
    Map(func(i []user) any {
        // 对每个批次内的数据再使用流处理
        return stream.ToStream(&i).Map(func(user user) any {
            properties, _ := convert.CopyProperties(user, &userVo{})
            return *properties
        }).ToList()
    }).
    Reduce("", func(cntValue any, nxt any) any {
        // 对每个批次处理的结果进行归约
        anies := nxt.([]any)
        return stream.ToStream(&anies).Reduce(cntValue, func(sumValue any, nxt any) any {
            return sumValue.(string) + nxt.(userVo).Name
        }, func(sum1, sum2 any) any {
            return sum1.(string) + sum2.(string)
        })
    }, func(sum1, sum2 any) any {
        // 合并并行处理的结果
        return sum1.(string) + sum2.(string)
    })
```

## 最佳实践

- 针对大数据集，选择合适的批次大小（通常是处理器核心数的整数倍）
- 使用并行处理时，注意控制并发数，避免资源耗尽
- 对于计算密集型任务，并行处理通常能带来性能提升
- 对于IO密集型任务，适当增加并发数可以提高吞吐量 