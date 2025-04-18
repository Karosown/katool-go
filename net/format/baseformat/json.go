package baseformat

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/karosown/katool-go/net/format"
)

type JSONEnDeCodeFormat struct {
	format.DefaultEnDeCodeFormat
	format.BytesDecodeFormatValid
}

func (c *JSONEnDeCodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (e *JSONEnDeCodeFormat) Encode(obj any) (any, error) {
	marshal, err := json.Marshal(obj)
	if err == nil {
		s := bytes.NewBuffer(marshal).String()
		return &s, nil
	}
	return nil, err
}

func (e *JSONEnDeCodeFormat) Decode(encode any, back any) (any, error) {
	err := json.Unmarshal(encode.([]byte), back)
	return back, err
}

type JSONArrayEnDeCodeFormat struct {
	format.DefaultEnDeCodeFormat
	format.BytesDecodeFormatValid
}

func (c *JSONArrayEnDeCodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (e *JSONArrayEnDeCodeFormat) Encode(obj any) (any, error) {
	marshal, err := json.Marshal(obj)
	if err == nil {
		s := bytes.NewBuffer(marshal).String()
		return &s, nil
	}
	return nil, err
}

func (e *JSONArrayEnDeCodeFormat) Decode(encode any, back any) (any, error) {
	anyArray := reflect.SliceOf(reflect.TypeOf(back))
	err := json.Unmarshal(encode.([]byte), anyArray)
	return back, err
}
