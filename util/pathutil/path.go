package pathutil

import (
	"path"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/optional"
)

type Wrapper struct {
	Path string
}

func NewWrapper(path string) *Wrapper {
	return &Wrapper{
		Path: path,
	}
}
func (p *Wrapper) Format() *Wrapper {
	p.Path = optional.IsTrue(p.Path[len(p.Path)-1] == '/', p.Path, p.Path+"/")
	return p
}
func (p *Wrapper) ForceCreate() bool {
	p.Format()
	exist := fileutil.IsExist(p.Path)
	if !exist {
		exist = fileutil.CreateDir(p.Path) == nil
	}
	return exist
}
func (p *Wrapper) ForceCreateEmptyFile(fileName string) bool {
	p.Format()
	p.ForceCreate()
	exist := fileutil.IsExist(p.Path + fileName)
	if !exist {
		exist = fileutil.CreateFile(p.Path + fileName)
	}
	return exist
}
func (p *Wrapper) BeforeLayer() *Wrapper {
	return &Wrapper{
		Path: optional.IsTrue(p.Path[len(p.Path)-1] == '/', path.Dir(p.Path[:len(p.Path)-1]), path.Dir(p.Path)),
	}
}
func (p *Wrapper) AfterLayer(childlayer string) *Wrapper {
	return &Wrapper{p.Path + "/" + childlayer}
}
