package remote

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipEnDecodeFormat struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
}

func (c *GzipEnDecodeFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (g *GzipEnDecodeFormat) Decode(encode any, back any) (any, error) {
	// 主要用于解压使用，具体的解析走下一个工具
	reader, err := gzip.NewReader(bytes.NewReader(encode.([]byte)))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, err
	}
	return decompressedData.Bytes(), nil
}
