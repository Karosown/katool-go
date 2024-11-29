package algorithm

import (
	"crypto/md5"
	"encoding/json"
)

type HashType string
type HashComputeFunction[T any] func(T) HashType

func HASH_WITH_JSON(cnt any) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		panic(err)
	}
	return HashType(marshal)
}

func HASH_WITH_JSON_MD5(cnt any) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		panic(err)
	}
	sum := md5.Sum(marshal)
	return HashType(sum[:])
}

func HASH_WITH_JSON_SUM(cnt any) HashType {
	json := HASH_WITH_JSON(cnt)
	sum := 0
	for _, v := range json {
		sum += (sum<<2 + sum) + int(v)
	}
	return HashType(sum)
}
