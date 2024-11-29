package algorithm

import (
	"encoding/json"
)

type HashType string

func HASH_WITH_JSON(cnt any) HashType {
	marshal, err := json.Marshal(cnt)
	if err != nil {
		panic(err)
	}
	return HashType(marshal)
}
