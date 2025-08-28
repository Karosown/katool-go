package avatar

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// MemoryFile 内存文件结构（保持不变）
type MemoryFile struct {
	*bytes.Reader
	name    string
	content []byte
	closed  bool
}

func NewMemoryFile(name string, content []byte) *MemoryFile {
	return &MemoryFile{
		Reader:  bytes.NewReader(content),
		name:    name,
		content: content,
		closed:  false,
	}
}

func (mf *MemoryFile) Name() string    { return mf.name }
func (mf *MemoryFile) Close() error    { mf.closed = true; return nil }
func (mf *MemoryFile) Content() []byte { return mf.content }

func (mf *MemoryFile) Stat() (os.FileInfo, error) {
	return &memoryFileInfo{name: mf.name, size: int64(len(mf.content))}, nil
}

func (mf *MemoryFile) WriteToFile(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	return os.WriteFile(filePath, mf.content, 0644)
}

type memoryFileInfo struct {
	name string
	size int64
}

func (mfi *memoryFileInfo) Name() string       { return mfi.name }
func (mfi *memoryFileInfo) Size() int64        { return mfi.size }
func (mfi *memoryFileInfo) Mode() os.FileMode  { return 0644 }
func (mfi *memoryFileInfo) ModTime() time.Time { return time.Now() }
func (mfi *memoryFileInfo) IsDir() bool        { return false }
func (mfi *memoryFileInfo) Sys() interface{}   { return nil }

// AvatarGenerator 头像生成器
type AvatarGenerator struct {
	Size         int           // 头像尺寸
	FontSize     int           // 字体大小
	TextColor    string        // 文字颜色
	BgColor      string        // 背景颜色
	OutputDir    string        // 输出目录
	UseHSL       bool          // 是否使用HSL颜色
	FontConfig   *FontConfig   // 字体配置
	FontRegistry *FontRegistry // 字体注册表
}

// FileResult 文件结果接口
type FileResult interface {
	io.ReadSeeker
	io.Closer
	Name() string
	Content() []byte
}
