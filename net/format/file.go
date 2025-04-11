package format

import (
	"encoding/json"
	"path"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/xmap"
	"github.com/karosown/katool-go/lock"
	"github.com/karosown/katool-go/sys"
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
		lock.Synchronized(fileLock, func() {
			if !fileutil.IsExist(filePath) {
				fileurl := path.Dir(filePath)
				if !fileutil.IsExist(fileurl) {
					pathLock, _ := e.FileLockers.LoadOrStore(fileurl, &sync.RWMutex{})
					lock.Synchronized(pathLock, func() {
						if !fileutil.IsExist(fileurl) {
							err := fileutil.CreateDir(fileurl)
							if err != nil {
								sys.Panic("create dir has error")
							}
						}
					})
				}
				file := fileutil.CreateFile(filePath)
				if file == false {
					sys.Panic("create file is failed")
				}
			}
		})
	}
	toString, err := cast.ToStringE(obj)
	if err != nil {
		bytes, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			sys.Warn("encode error")
			return nil, err
		}
		toString, err = string(bytes), nil
	}
	lock.Synchronized(fileLock, func() {
		fileutil.WriteStringToFile(filePath, toString, e.Append)
	})
	return filePath, nil
}

func (e *FileSaveFormat) Decode(encode any, back any) (any, error) {
	return e.Encode(encode)
}
