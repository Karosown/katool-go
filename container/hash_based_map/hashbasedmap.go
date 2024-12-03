package hash_based_map

type HashBasedMap[K1, K2, V comparable] struct {
	m map[K1]map[K2]V
}

// New 创建一个新的 HashBasedMap
func NewHashBasedMap[K1, K2, V comparable]() *HashBasedMap[K1, K2, V] {
	return &HashBasedMap[K1, K2, V]{
		m: make(map[K1]map[K2]V),
	}
}

// Set 设置特定的值
func (hm *HashBasedMap[K1, K2, V]) Set(key1 K1, key2 K2, value V) {
	if _, exists := hm.m[key1]; !exists {
		hm.m[key1] = make(map[K2]V)
	}
	hm.m[key1][key2] = value
}

// Get 获取特定的值，并返回是否存在
func (hm *HashBasedMap[K1, K2, V]) Get(key1 K1, key2 K2) (V, bool) {
	if subMap, exists := hm.m[key1]; exists {
		value, found := subMap[key2]
		return value, found
	}
	var zero V
	return zero, false
}

// Delete 删除特定的值
func (hm *HashBasedMap[K1, K2, V]) Delete(key1 K1, key2 K2) {
	if _, exists := hm.m[key1]; exists {
		delete(hm.m[key1], key2)
		if len(hm.m[key1]) == 0 {
			delete(hm.m, key1)
		}
	}
}

// Keys 返回所有的一级和二级键
func (hm *HashBasedMap[K1, K2, V]) Keys() ([]K1, map[K1][]K2) {
	firstLevelKeys := make([]K1, 0, len(hm.m))
	secondLevelKeys := make(map[K1][]K2)

	for k1, subMap := range hm.m {
		firstLevelKeys = append(firstLevelKeys, k1)
		secondLevelKeys[k1] = make([]K2, 0, len(subMap))
		for k2 := range subMap {
			secondLevelKeys[k1] = append(secondLevelKeys[k1], k2)
		}
	}

	return firstLevelKeys, secondLevelKeys
}

// Size 返回总的元素数量
func (hm *HashBasedMap[K1, K2, V]) Size() int {
	total := 0
	for _, subMap := range hm.m {
		total += len(subMap)
	}
	return total
}

// Clear 清空整个映射
func (hm *HashBasedMap[K1, K2, V]) Clear() {
	hm.m = make(map[K1]map[K2]V)
}
