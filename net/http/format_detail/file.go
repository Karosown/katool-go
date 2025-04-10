package remote

import (
	"encoding/json"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/xmap"
	"github.com/karosown/katool-go/xlog"
	"github.com/spf13/cast"
)

type FileSaveFormat[T any] struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
	FileLockers         xmap.SafeMap[string, *sync.RWMutex]
	FileFullNameBuilder func(data T) string
	Status              int
}

func (c *FileSaveFormat[T]) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (e *FileSaveFormat[T]) Encode(obj any) (any, error) {
	filePath := e.FileFullNameBuilder(obj)
	get, b := e.FileLockers.Get(filePath)
	if !b {
		store, b2 := e.FileLockers.LoadOrStore(filePath, &sync.RWMutex{})
		if !b2 {
			xlog.KaToolLoggerWrapper.ApplicationDesc("get lock error").Panic()
		}
		get = store
	}
	if fileutil.IsExist(filePath) {
		get.Lock()
		if fileutil.IsExist(filePath) {
			err := fileutil.CreateDir(filePath)
			if err != nil {
				xlog.KaToolLoggerWrapper.ApplicationDesc("create dir error").Panic()
			}
		}
		get.Unlock()
	}
	toString, err := cast.ToStringE(obj)
	if err != nil {
		bytes, err := json.Marshal(obj)
		if err != nil {
			xlog.KaToolLoggerWrapper.ApplicationDesc("encode error").Panic()
		}
		toString = string(bytes)
	}
	get.Lock()
	fileutil.WriteStringToFile(filePath, toString, true)
	get.Unlock()
	return nil, err
}

func (e *FileSaveFormat[T]) Decode(encode any, back any) (any, error) {
	return e.Encode(encode)
}
