package fixhelper

import (
	"encoding/json"

	"github.com/kaptinlin/jsonrepair"
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
func JsonUnMarshal[T any](unfixJson string) *T {
	fixedJson := FixJson(unfixJson)
	back := new(T)
	err := json.Unmarshal([]byte(fixedJson), back)
	if err != nil {
		return nil
	}
	return back
}
