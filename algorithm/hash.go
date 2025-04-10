package algorithm

import (
	"crypto/md5"
	"encoding/json"
	"strconv"

	"github.com/karosown/katool-go/xlog"
)

type HashType string
type IDType int64
type HashComputeFunction func(any2 any) HashType
type IDComputeFunction func(any2 any) IDType

func HASH_WITH_JSON(cnt any) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		xlog.KaToolLoggerWrapper.Warn().ApplicationDesc(err.Error()).Panic()
	}
	return HashType(marshal)
}

func HASH_WITH_JSON_MD5(cnt any) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		xlog.KaToolLoggerWrapper.Warn().ApplicationDesc(err.Error()).Panic()
	}
	sum := md5.Sum(marshal)
	return HashType(sum[:])
}

func HASH_WITH_JSON_SUM(cnt any) HashType {
	jsonStr := HASH_WITH_JSON(cnt)
	sum := 0
	for _, v := range jsonStr {
		sum += (sum<<2 + sum) + int(v)
	}
	return HashType(strconv.Itoa(sum))
}
