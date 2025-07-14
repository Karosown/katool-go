package jsonhp

import (
	"encoding/json"

	"github.com/kaptinlin/jsonrepair"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/sys"
)

/*
*
对于残缺的json进行修复
*/
func FixJson(unfixJson string) string {
	fixedJson, err := jsonrepair.JSONRepair(unfixJson)
	if err != nil {
		return ""
	}
	return fixedJson
}
func ToJSON(entity any) string {
	switch entity.(type) {
	case string:
		return FixJson(entity.(string))
	default:
		marshal, err := json.Marshal(entity)
		if err != nil {
			return ""
		}
		return string(marshal)
	}
}
func ToJSONIndent(entity any, prefixAndIndent ...string) string {
	if cutil.IsEmpty(prefixAndIndent) {
		prefixAndIndent = []string{"", " "}
	}
	if cutil.IsBlank(prefixAndIndent[0]) {
		prefixAndIndent[0] = ""
	}
	if cutil.IsBlank(prefixAndIndent[1]) {
		prefixAndIndent[1] = " "
	}
	switch entity.(type) {
	case string:
		return FixJson(entity.(string))
	default:
		marshal, err := json.MarshalIndent(entity, prefixAndIndent[0], prefixAndIndent[1])
		if err != nil {
			return ""
		}
		return string(marshal)
	}
}
func JsonUnMarshal[T any](unfixJson string) *T {
	fixedJson := FixJson(unfixJson)
	back := new(T)
	err := json.Unmarshal([]byte(fixedJson), back)
	if err != nil {
		return nil
	}
	return back
}
func ToJsonLine[T any](entities any) string {
	switch entities.(type) {
	case []T:
		ts := entities.([]T)
		reduce := stream.ToStream(&ts).Reduce("", func(cntValue any, nxt T) any {
			marshal := ToJSON(nxt)
			if marshal == "" {
				return cntValue
			}
			return cntValue.(string) + string(marshal) + "\n"
		}, func(cntValue any, nxt any) any {
			return cntValue.(string) + nxt.(string) + "\n"
		})
		return reduce.(string)
	case string:
		ts := JsonUnMarshal[[]T](entities.(string))
		reduce := stream.ToStream(ts).Reduce("", func(cntValue any, nxt T) any {
			marshal := ToJSON(nxt)
			if marshal == "" {
				return cntValue
			}
			return cntValue.(string) + marshal + "\n"
		}, func(cntValue any, nxt any) any {
			return cntValue.(string) + nxt.(string) + "\n"
		})
		return reduce.(string)
	case []byte:
		return ToJsonLine[T](string(entities.([]byte)))
	case [][]byte:
		entries := entities.([][]byte)
		reduce := stream.ToStream(&entries).Reduce([]T{}, func(cntValue any, nxt []byte) any {
			marshal := JsonUnMarshal[T](string(nxt))
			return append(cntValue.([]T), *marshal)
		}, func(cntValue any, nxt any) any {
			return append(cntValue.([]T), nxt.([]T)...)
		})
		return ToJsonLine[T](reduce)
	default:
		sys.Panic("Not Support This Type")
		return ""
	}
}
