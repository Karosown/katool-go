package convert

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)

func Convert[T any, R any](datas []T, vacuumMachine func(agent T) R) []R {
	res := make([]R, 0, len(datas))
	for _, data := range datas {
		res = append(res, vacuumMachine(data))
	}
	return res
}

func CopyProperties[T, R any](src T, dest R) (R, error) {
	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	destType := reflect.TypeOf(dest)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
	}
	ret := destType
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	restPtr := reflect.New(destType)
	rest := restPtr.Elem()
	destFieldNum := destType.NumField()
	for i := 0; i < destFieldNum; i++ {
		dstField := destType.Field(i)
		currentFieldName := dstField.Name
		dstFieldValue := rest.FieldByName(currentFieldName)
		orginFieldValue := srcValue.FieldByName(currentFieldName)
		orginFieldType, e := srcType.FieldByName(currentFieldName)
		if !e {
			continue
		}
		if orginFieldType.Type != dstFieldValue.Type() {
			return restPtr.Interface().(R), errors.New("field type not match")
		}
		dstFieldValue.Set(orginFieldValue)
	}
	if ret != destType {
		return restPtr.Interface().(R), nil
	}
	return rest.Interface().(R), nil
}

func fieldSetValue(field *reflect.Value, fieldValue reflect.Value) error {
	field.Set(fieldValue.Convert(field.Type()))
	return nil
}

func ToString(source any) string {
	var str string
	if source == nil {
		return str
	}
	// vt := source.(type)
	if reflect.ValueOf(source).Kind() == reflect.Ptr {
		source = reflect.ValueOf(source).Elem().Interface()
	}
	switch source.(type) {
	case float64:
		ft := source.(float64)
		str = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := source.(float32)
		str = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := source.(int)
		str = strconv.Itoa(it)
	case uint:
		it := source.(uint)
		str = strconv.Itoa(int(it))
	case int8:
		it := source.(int8)
		str = strconv.Itoa(int(it))
	case uint8:
		it := source.(uint8)
		str = strconv.Itoa(int(it))
	case int16:
		it := source.(int16)
		str = strconv.Itoa(int(it))
	case uint16:
		it := source.(uint16)
		str = strconv.Itoa(int(it))
	case int32:
		it := source.(int32)
		str = strconv.Itoa(int(it))
	case uint32:
		it := source.(uint32)
		str = strconv.Itoa(int(it))
	case int64:
		it := source.(int64)
		str = strconv.FormatInt(it, 10)
	case uint64:
		it := source.(uint64)
		str = strconv.FormatUint(it, 10)
	case string:
		str = source.(string)
	case []byte:
		str = string(source.([]byte))
	case time.Time:
		t := source.(time.Time)
		jsone, _ := json.Marshal(t)
		str = string(jsone)
		str = str[1 : len(str)-1]
	case bool:
		b := source.(bool)
		if b {
			str = "true"
		} else {
			str = "false"
		}
	default:
		newValue, _ := json.Marshal(source)
		str = string(newValue)
	}
	return str
}
