package remote

import (
	"bytes"
	"encoding/xml"

	remote "github.com/karosown/katool/net/http"
)

type XMLEnDeCodeFormat struct {
	remote.DefaultEnDeCodeFormat
	remote.BytesDecodeFormatValid
}

func (c *XMLEnDeCodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (e *XMLEnDeCodeFormat) Encode(obj any) (any, error) {
	marshal, err := xml.Marshal(obj)
	if err == nil {
		s := bytes.NewBuffer(marshal).String()
		return &s, nil
	}
	return nil, err
}

func (e *XMLEnDeCodeFormat) Decode(encode any, back any) (any, error) {
	err := xml.Unmarshal(encode.([]byte), back)
	return nil, err
}
