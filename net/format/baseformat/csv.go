package baseformat

import (
	"fmt"

	"github.com/karosown/katool-go/net/format"

	"github.com/karosown/katool-go/file/file_serialize"
)

type CSVEnDecodeFormat struct {
	format.DefaultEnDeCodeFormat
	format.BytesDecodeFormatValid
}

func (c *CSVEnDecodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (c *CSVEnDecodeFormat) Encode(obj any) (any, error) {
	return nil, fmt.Errorf("current not support csv encode")
}
func (c *CSVEnDecodeFormat) Decode(encode any, back any) (any, error) {
	res := file_serialize.CSVSerializer{}.ReadByBytes(encode.([]byte), back)
	return res, nil
}
