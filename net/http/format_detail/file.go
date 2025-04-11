package remote

import (
	"encoding/json"
	"path"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/xmap"
	"github.com/karosown/katool-go/xlog"
	"github.com/spf13/cast"
)

type FileSaveFormat struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
	FileLockers         xmap.SafeMap[string, *sync.RWMutex]
	FileFullNameBuilder func(data any) string
	FileFilterFunc      func(data any) any
	Append              bool
}

func (c *FileSaveFormat) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (e *FileSaveFormat) Encode(obj any) (any, error) {
	filePath := e.FileFullNameBuilder(obj)
	obj = e.FileFilterFunc(obj)
	fileLock, b := e.FileLockers.Get(filePath)
	if !b {
		store, _ := e.FileLockers.LoadOrStore(filePath, &sync.RWMutex{})
		fileLock = store
	}
	if !fileutil.IsExist(filePath) {
		fileLock.Lock()
		if !fileutil.IsExist(filePath) {
			fileurl := path.Dir(filePath)
			if !fileutil.IsExist(fileurl) {
				pathLock, _ := e.FileLockers.LoadOrStore(fileurl, &sync.RWMutex{})
				pathLock.Lock()
				if !fileutil.IsExist(fileurl) {
					err := fileutil.CreateDir(fileurl)
					if err != nil {
						xlog.KaToolLoggerWrapper.ApplicationDesc("create dir has error").Panic()
					}
				}
				pathLock.Unlock()
			}
			file := fileutil.CreateFile(filePath)
			if file == false {
				xlog.KaToolLoggerWrapper.ApplicationDesc("create file is failed").Panic()
			}
		}
		fileLock.Unlock()
	}
	toString, err := cast.ToStringE(obj)
	if err != nil {
		bytes, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			xlog.KaToolLoggerWrapper.Warn().ApplicationDesc("encode error").Panic()
			return nil, err
		}
		toString, err = string(bytes), nil
	}
	fileLock.Lock()
	fileutil.WriteStringToFile(filePath, toString, e.Append)
	fileLock.Unlock()
	return filePath, nil
}

func (e *FileSaveFormat) Decode(encode any, back any) (any, error) {
	return e.Encode(encode)
}
