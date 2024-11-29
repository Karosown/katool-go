package file_serialize

import "os"

type SerializedInterface interface {
	ReadByBytes(buf []byte, backDao any) []any
	ReadByFile(file *os.File, backDao any) []any
	ReadByPath(path string, backDao any) []any
	Write(path string, sourceDao any) error
}
