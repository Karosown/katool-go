// ziputil.go
package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/karosown/katool-go/util/pathutil"
)

// CompressToZip 压缩文件或目录到ZIP文件
func CompressToZip(src, dest string) error {
	pathutil.NewWrapper(dest).BeforeLayer().ForceCreate()
	// 创建ZIP文件
	zipFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建ZIP文件失败: %w", err)
	}
	defer zipFile.Close()

	// 创建ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 获取源文件/目录信息
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	// 遍历文件并添加到ZIP
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建ZIP文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, src))
		} else {
			header.Name = strings.TrimPrefix(path, filepath.Dir(src)+string(os.PathSeparator))
		}

		// 规范化路径分隔符
		header.Name = filepath.ToSlash(header.Name)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 创建writer
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 写入文件内容
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

// ExtractZip 解压ZIP文件到指定目录
func ExtractZip(src, dest string) error {
	// 打开ZIP文件
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	// 创建目标目录
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 解压文件
	for _, file := range reader.File {
		err := extractFile(file, dest)
		if err != nil {
			return fmt.Errorf("解压文件 %s 失败: %w", file.Name, err)
		}
	}

	return nil
}

// extractFile 解压单个文件
func extractFile(file *zip.File, destDir string) error {
	// 构建目标路径
	path := filepath.Join(destDir, file.Name)

	// 检查路径安全性，防止目录遍历攻击
	if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("无效的文件路径: %s", file.Name)
	}

	// 如果是目录
	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, file.FileInfo().Mode())
	}

	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// 打开ZIP中的文件
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	// 创建目标文件
	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// 复制内容
	_, err = io.Copy(targetFile, fileReader)
	return err
}

// ListZipContents 列出ZIP文件内容
func ListZipContents(zipPath string) ([]FileInfo, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	var files []FileInfo
	for _, file := range reader.File {
		info := FileInfo{
			Name:           file.Name,
			Size:           file.FileInfo().Size(),
			CompressedSize: int64(file.CompressedSize64),
			ModTime:        file.FileInfo().ModTime(),
			IsDir:          file.FileInfo().IsDir(),
			Method:         file.Method,
		}
		files = append(files, info)
	}

	return files, nil
}

// FileInfo ZIP文件信息结构体
type FileInfo struct {
	Name           string
	Size           int64
	CompressedSize int64
	ModTime        time.Time
	IsDir          bool
	Method         uint16
}

// String 实现Stringer接口
func (f FileInfo) String() string {
	if f.IsDir {
		return fmt.Sprintf("%-50s [DIR] %s", f.Name, f.ModTime.Format("2006-01-02 15:04:05"))
	}
	ratio := float64(f.CompressedSize) / float64(f.Size) * 100
	return fmt.Sprintf("%-50s %8d -> %8d (%.1f%%) %s",
		f.Name, f.Size, f.CompressedSize, ratio, f.ModTime.Format("2006-01-02 15:04:05"))
}

// ExtractFile 从ZIP中提取单个文件
func ExtractFile(zipPath, fileName, destPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if file.Name == fileName {
			return extractFile(file, filepath.Dir(destPath))
		}
	}

	return fmt.Errorf("文件 %s 在ZIP中不存在", fileName)
}

// AddToZip 向现有ZIP文件添加文件
func AddToZip(zipPath, filePath string) error {
	// 先读取现有ZIP内容
	var existingFiles []*zip.File
	if _, err := os.Stat(zipPath); err == nil {
		reader, err := zip.OpenReader(zipPath)
		if err != nil {
			return fmt.Errorf("打开现有ZIP文件失败: %w", err)
		}
		existingFiles = reader.File
		defer reader.Close()
	}

	// 创建新的ZIP文件
	newZip, err := os.Create(zipPath + ".tmp")
	if err != nil {
		return fmt.Errorf("创建临时ZIP文件失败: %w", err)
	}
	defer newZip.Close()

	zipWriter := zip.NewWriter(newZip)
	defer zipWriter.Close()

	// 复制现有文件
	for _, file := range existingFiles {
		err := copyZipFile(file, zipWriter)
		if err != nil {
			return fmt.Errorf("复制现有文件失败: %w", err)
		}
	}

	// 添加新文件
	err = addFileToZip(filePath, zipWriter)
	if err != nil {
		return fmt.Errorf("添加新文件失败: %w", err)
	}

	zipWriter.Close()
	newZip.Close()

	// 替换原文件
	return os.Rename(zipPath+".tmp", zipPath)
}

// copyZipFile 复制ZIP文件中的文件
func copyZipFile(file *zip.File, writer *zip.Writer) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	fileWriter, err := writer.CreateHeader(&file.FileHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, fileReader)
	return err
}

// addFileToZip 添加文件到ZIP
func addFileToZip(filePath string, writer *zip.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	fileWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	return err
}

// GetZipInfo 获取ZIP文件基本信息
func GetZipInfo(zipPath string) (*ZipInfo, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	info := &ZipInfo{
		Path:      zipPath,
		FileCount: len(reader.File),
	}

	for _, file := range reader.File {
		info.TotalSize += file.FileInfo().Size()
		info.CompressedSize += int64(file.CompressedSize64)
		if file.FileInfo().IsDir() {
			info.DirCount++
		}
	}

	return info, nil
}

// ZipInfo ZIP文件信息
type ZipInfo struct {
	Path           string
	FileCount      int
	DirCount       int
	TotalSize      int64
	CompressedSize int64
}

// String 实现Stringer接口
func (z ZipInfo) String() string {
	ratio := float64(z.CompressedSize) / float64(z.TotalSize) * 100
	return fmt.Sprintf(`ZIP文件信息:
  路径: %s
  文件数: %d
  目录数: %d  
  原始大小: %d 字节
  压缩大小: %d 字节
  压缩率: %.1f%%`,
		z.Path, z.FileCount, z.DirCount, z.TotalSize, z.CompressedSize, ratio)
}
