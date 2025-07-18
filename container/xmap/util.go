package xmap

import (
	"encoding/json"
)

func ToMap[T any](obj T) (Map[string, any], error) {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	newMap := NewMap[string, any]()
	json.Unmarshal(marshal, newMap)
	return newMap, nil
}
