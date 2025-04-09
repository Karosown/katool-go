package test

import (
	"sync"
	"testing"

	"github.com/karosown/katool/container/xmap"
	"github.com/stretchr/testify/assert"
)

// 测试 Map 的基本功能
func TestMap(t *testing.T) {
	// 创建新的 Map
	m := xmap.NewMap[string, int]()

	// 测试初始状态
	assert.Equal(t, 0, m.Len(), "新创建的 Map 应该为空")

	// 测试 Set 和 Get
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	val, exists := m.Get("one")
	assert.True(t, exists, "键 'one' 应该存在")
	assert.Equal(t, 1, val, "键 'one' 的值应为 1")

	// 测试 Len
	assert.Equal(t, 3, m.Len(), "Map 应该有 3 个元素")

	// 测试 Has
	assert.True(t, m.Has("two"), "键 'two' 应该存在")
	assert.False(t, m.Has("four"), "键 'four' 不应该存在")

	// 测试 Delete
	m.Delete("two")
	assert.False(t, m.Has("two"), "删除后键 'two' 不应该存在")
	assert.Equal(t, 2, m.Len(), "删除后 Map 应该有 2 个元素")

	// 测试 Keys 和 Values
	keys := m.Keys()
	assert.Equal(t, 2, len(keys), "应该有 2 个键")
	assert.Contains(t, keys, "one", "键应包含 'one'")
	assert.Contains(t, keys, "three", "键应包含 'three'")

	values := m.Values()
	assert.Equal(t, 2, len(values), "应该有 2 个值")
	assert.Contains(t, values, 1, "值应包含 1")
	assert.Contains(t, values, 3, "值应包含 3")

	// 测试 Clear
	m.Clear()
	assert.Equal(t, 0, m.Len(), "清空后 Map 应为空")

	// 测试 Reset
	m.Set("test", 100)
	assert.Equal(t, 1, m.Len(), "添加元素后长度应为 1")

	m = m.Reset()
	assert.Equal(t, 0, m.Len(), "重置后 Map 应为空")

	// 测试 ForEach
	m.Set("a", 1)
	m.Set("b", 2)

	count := 0
	sum := 0
	m.ForEach(func(k string, v int) {
		count++
		sum += v
	})

	assert.Equal(t, 2, count, "应遍历 2 个元素")
	assert.Equal(t, 3, sum, "值的总和应为 3")

	// 测试 ToStream (简单验证，不测试 stream 包的功能)
	stream := m.ToStream()
	assert.NotNil(t, stream, "应返回非空的 stream")
}

// 测试 SafeMap 的基本功能
func TestSafeMap(t *testing.T) {
	// 创建新的 SafeMap
	sm := xmap.NewSafeMap[string, int]()

	// 测试初始状态
	assert.Equal(t, 0, sm.Len(), "新创建的 SafeMap 应该为空")

	// 测试 Set 和 Get
	sm.Set("one", 1)
	sm.Set("two", 2)
	sm.Set("three", 3)

	val, exists := sm.Get("one")
	assert.True(t, exists, "键 'one' 应该存在")
	assert.Equal(t, 1, val, "键 'one' 的值应为 1")

	// 测试 Len
	assert.Equal(t, 3, sm.Len(), "SafeMap 应该有 3 个元素")

	// 测试 Has
	assert.True(t, sm.Has("two"), "键 'two' 应该存在")
	assert.False(t, sm.Has("four"), "键 'four' 不应该存在")

	// 测试 Delete
	sm.Delete("two")
	assert.False(t, sm.Has("two"), "删除后键 'two' 不应该存在")
	assert.Equal(t, 2, sm.Len(), "删除后 SafeMap 应该有 2 个元素")

	// 测试 LoadOrStore
	actual, loaded := sm.LoadOrStore("one", 100)
	assert.True(t, loaded, "键 'one' 应该已存在")
	assert.Equal(t, 1, actual, "应返回已存在的值 1 而非 100")

	actual, loaded = sm.LoadOrStore("four", 4)
	assert.False(t, loaded, "键 'four' 不应该已存在")
	assert.Equal(t, 4, actual, "应返回新存储的值 4")

	// 测试 LoadAndDelete
	val, exists = sm.LoadAndDelete("four")
	assert.True(t, exists, "键 'four' 应该存在")
	assert.Equal(t, 4, val, "应返回键 'four' 的值 4")
	assert.False(t, sm.Has("four"), "LoadAndDelete 后键 'four' 不应该存在")

	// 测试 Clear
	sm.Clear()
	assert.Equal(t, 0, sm.Len(), "清空后 SafeMap 应为空")

	// 测试 Reset
	sm.Set("test", 100)
	assert.Equal(t, 1, sm.Len(), "添加元素后长度应为 1")

	sm.Reset()
	assert.Equal(t, 0, sm.Len(), "重置后 SafeMap 应为空")

	// 测试 ToMap
	sm.Set("a", 1)
	sm.Set("b", 2)

	regularMap := sm.ToMap()
	assert.Equal(t, 2, len(regularMap), "转换后的普通 Map 应有 2 个元素")
	assert.Equal(t, 1, regularMap["a"], "键 'a' 的值应为 1")
	assert.Equal(t, 2, regularMap["b"], "键 'b' 的值应为 2")

	// 测试 ToStream (简单验证)
	stream := sm.ToStream()
	assert.NotNil(t, stream, "应返回非空的 stream")
}

// 测试 SafeMap 的并发安全性
func TestSafeMapConcurrency(t *testing.T) {
	sm := xmap.NewSafeMap[int, int]()
	const goroutines = 100
	const countsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 并发写入
	for i := 0; i < goroutines; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < countsPerGoroutine; j++ {
				key := base*countsPerGoroutine + j
				sm.Set(key, key*2)
			}
		}(i)
	}

	wg.Wait()

	// 验证结果
	assert.Equal(t, goroutines*countsPerGoroutine, sm.Len(),
		"SafeMap 应该包含所有写入的元素")

	// 并发读取和验证
	var readWg sync.WaitGroup
	readWg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(base int) {
			defer readWg.Done()
			for j := 0; j < countsPerGoroutine; j++ {
				key := base*countsPerGoroutine + j
				val, exists := sm.Get(key)
				assert.True(t, exists, "键应该存在")
				assert.Equal(t, key*2, val, "值应该正确")
			}
		}(i)
	}

	readWg.Wait()
}
