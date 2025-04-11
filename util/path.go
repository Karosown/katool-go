package util

import (
	"path"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/optional"
)

type PathUtil struct {
	Path string
}

func NewPathUtil(path string) *PathUtil {
	return &PathUtil{
		Path: path,
	}
}
func (p *PathUtil) Format() *PathUtil {
	p.Path = optional.IsTrue(p.Path[len(p.Path)-1] == '/', p.Path, p.Path+"/")
	return p
}
func (p *PathUtil) ForceCreate() bool {
	p.Format()
	exist := fileutil.IsExist(p.Path)
	if !exist {
		exist = fileutil.CreateDir(p.Path) == nil
	}
	return exist
}
func (p *PathUtil) ForceCreateEmptyFile(fileName string) bool {
	p.Format()
	p.ForceCreate()
	exist := fileutil.IsExist(p.Path + fileName)
	if !exist {
		exist = fileutil.CreateFile(p.Path + fileName)
	}
	return exist
}
func (p *PathUtil) BeforeLayer() *PathUtil {
	return &PathUtil{
		Path: optional.IsTrue(p.Path[len(p.Path)-1] == '/', path.Dir(p.Path[:len(p.Path)-1]), path.Dir(p.Path)),
	}
}
func (p *PathUtil) AfterLayer(childlayer string) *PathUtil {
	return &PathUtil{p.Path + "/" + childlayer}
}
