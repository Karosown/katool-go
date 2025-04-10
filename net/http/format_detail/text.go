package remote

import (
	"fmt"

	"github.com/karosown/katool-go/file/file_serialize"
)

type TextEnDecodeFormat struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
	MapToStruct func(fields []string, headers []string, backDao any) any
	Convert     func(backDao any) any
	BackDao     any
}

func (c *TextEnDecodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (c *TextEnDecodeFormat) Encode(obj any) (any, error) {
	return nil, fmt.Errorf("This Function Not Support. If you wangt to use this function, please complate the function Encode(obj any) (any,error)")
}
func (c *TextEnDecodeFormat) Decode(encode any, back any) (any, error) {
	if c.BackDao != nil {
		back = c.BackDao
	}
	res := file_serialize.TextFileSerializer{
		MapToStruct: c.MapToStruct,
		Convert:     c.Convert,
	}.ReadByBytes(encode.([]byte), back)
	return res, nil
}
