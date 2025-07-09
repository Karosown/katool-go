package pathutil

import (
	"path"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/optional"
)

// Wrapper 路径包装器
// Wrapper is a path wrapper
type Wrapper struct {
	Path string
}

// NewWrapper 创建新的路径包装器
// NewWrapper creates a new path wrapper
func NewWrapper(path string) *Wrapper {
	return &Wrapper{
		Path: path,
	}
}

// Format 格式化路径（确保以/结尾）
// Format formats the path (ensures it ends with /)
func (p *Wrapper) Format() *Wrapper {
	p.Path = optional.IsTrue(p.Path[len(p.Path)-1] == '/', p.Path, p.Path+"/")
	return p
}

// ForceCreate 强制创建目录
// ForceCreate forcibly creates the directory
func (p *Wrapper) ForceCreate() bool {
	p.Format()
	exist := fileutil.IsExist(p.Path)
	if !exist {
		exist = fileutil.CreateDir(p.Path) == nil
	}
	return exist
}

// ForceCreateEmptyFile 强制创建空文件
// ForceCreateEmptyFile forcibly creates an empty file
func (p *Wrapper) ForceCreateEmptyFile(fileName string) bool {
	p.Format()
	p.ForceCreate()
	exist := fileutil.IsExist(p.Path + fileName)
	if !exist {
		exist = fileutil.CreateFile(p.Path + fileName)
	}
	return exist
}

// BeforeLayer 获取上一层目录
// BeforeLayer gets the parent directory
func (p *Wrapper) BeforeLayer() *Wrapper {
	return &Wrapper{
		Path: optional.IsTrue(p.Path[len(p.Path)-1] == '/', path.Dir(p.Path[:len(p.Path)-1]), path.Dir(p.Path)),
	}
}

// AfterLayer 获取子目录
// AfterLayer gets the child directory
func (p *Wrapper) AfterLayer(childlayer string) *Wrapper {
	return &Wrapper{p.Path + "/" + childlayer}
}
