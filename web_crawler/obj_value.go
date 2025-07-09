package web_crawler

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// AttrValue 属性值包装器
// AttrValue is an attribute value wrapper
type AttrValue struct {
	WebReaderValue
	WebReaderError
}

// ReadArray 读取数组数据
// ReadArray reads array data
func ReadArray(url, js string, renderFunc func(*rod.Page)) AttrValue {
	var gen func() AttrValue
	tryNum := 7
	gen = func() AttrValue {
		code := readArray(url, js, renderFunc)
		if code.IsErr() {
			if tryNum != 0 {
				tryNum--
				if tryNum == 0 {
					WebChrome.ReStart()
				} else {
					time.Sleep(time.Duration(7-tryNum+1) * time.Second)
				}
				return gen()
			}
			return AttrValue{
				WebReaderValue: nil,
				WebReaderError: WebReaderError{
					errors.New("the AttrValue Get Error"),
				},
			}
		}
		return code
	}
	return gen()
}

// readArray 内部读取数组方法
// readArray is the internal method for reading arrays
func readArray(url, js string, rendorFunc func(*rod.Page)) AttrValue {
	results, err := execFun(url, js, rendorFunc)
	if err != nil {
		return AttrValue{
			WebReaderValue{}, WebReaderError{err},
		}
	}
	res := results.Value.String()[1:len(results.Value.String())]
	result := strings.Split(res, " ")
	return AttrValue{
		NewWebReaderValue(result...), WebReaderError{err},
	}
}

// JsonValue JSON值包装器
// JsonValue is a JSON value wrapper
type JsonValue[T any] struct {
	Value *T
	WebReaderError
}

// ReadToJson 读取JSON数据并反序列化为指定类型
// ReadToJson reads JSON data and deserializes to specified type
func ReadToJson[T any](url string, obj T, js string, renderFunc func(*rod.Page)) JsonValue[T] {
	var gen func() JsonValue[T]
	tryNum := 7
	gen = func() JsonValue[T] {
		code := readToJson(url, obj, js, renderFunc)
		if code.IsErr() {
			if tryNum != 0 {
				tryNum--
				if tryNum == 0 {
					WebChrome.ReStart()
				} else {
					time.Sleep(time.Duration(7-tryNum+1) * time.Second)
				}
				return gen()
			}
			return JsonValue[T]{
				Value: nil,
				WebReaderError: WebReaderError{
					errors.New("the Json Get Error:" + code.error.Error()),
				},
			}
		}
		return code
	}
	return gen()
}

// readToJson 内部读取JSON方法
// readToJson is the internal method for reading JSON
func readToJson[T any](url string, obj T, js string, rendorFunc func(*rod.Page)) JsonValue[T] {
	result, err := execFun(url, js, rendorFunc)
	if result == nil || err != nil {
		return JsonValue[T]{
			nil, WebReaderError{err},
		}
	}
	// 根据T来创建对象
	value := reflect.New(reflect.TypeOf(obj)).Interface()

	err = json.Unmarshal([]byte(result.Value.String()), value)
	return JsonValue[T]{
		value.(*T), WebReaderError{err},
	}
}
