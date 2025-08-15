package algorithm

import (
	"crypto/md5"
	"encoding/json"
	"strconv"

	"github.com/karosown/katool-go/sys"
)

// HashType 哈希值类型
// HashType represents a hash value type
type HashType string

// IDType ID类型
// IDType represents an ID type
type IDType int64

// HashComputeFunction 哈希计算函数类型
// HashComputeFunction represents a function type for computing hash values
type HashComputeFunction[T any] func(any2 T) HashType

// IDComputeFunction ID计算函数类型
// IDComputeFunction represents a function type for computing ID values
type IDComputeFunction func(any2 any) IDType

// HASH_WITH_JSON 使用JSON序列化计算哈希值
// HASH_WITH_JSON computes hash value using JSON serialization
func HASH_WITH_JSON[T any](cnt T) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		sys.Warn(err.Error())
	}
	return HashType(marshal)
}

// HASH_WITH_JSON_MD5 使用JSON序列化后MD5计算哈希值
// HASH_WITH_JSON_MD5 computes hash value using MD5 after JSON serialization
func HASH_WITH_JSON_MD5[T any](cnt T) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		sys.Warn(err.Error())
	}
	sum := md5.Sum(marshal)
	return HashType(sum[:])
}

// HASH_WITH_JSON_SUM 使用JSON序列化后累加计算哈希值
// HASH_WITH_JSON_SUM computes hash value using sum calculation after JSON serialization
func HASH_WITH_JSON_SUM(cnt any) HashType {
	jsonStr := HASH_WITH_JSON(cnt)
	sum := 0
	for _, v := range jsonStr {
		sum += (sum<<2 + sum) + int(v)
	}
	return HashType(strconv.Itoa(sum))
}
