package file_downloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// FileDownloader 是一个文件下载工具类
// FileDownloader is a file download utility class
type FileDownloader struct {
	Client *http.Client
}

// NewFileDownloader 创建一个新的 FileDownloader 实例
// NewFileDownloader creates a new FileDownloader instance
func NewFileDownloader() *FileDownloader {
	return &FileDownloader{
		Client: &http.Client{},
	}
}

// DownloadFile 下载单个文件
// DownloadFile downloads a single file
func (fd *FileDownloader) DownloadFile(url, destPath string) error {
	// 创建目标文件
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	byties, err := fd.DownloadFileBytes(url)

	// 将响应体写入文件
	_, err = io.Copy(out, bytes.NewReader(byties))
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
func (fd *FileDownloader) DownloadFileBytes(url string) ([]byte, error) {
	resp, err := fd.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败，服务器返回状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// DownloadFiles 并发下载多个文件，请自行保证url不同，保证并发安全
// DownloadFiles downloads multiple files concurrently, ensure URLs are different for concurrency safety
func (fd *FileDownloader) DownloadFiles(urls []string, destDir string) []error {
	var wg sync.WaitGroup
	errors := make([]error, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			fileName := filepath.Base(url)
			destPath := filepath.Join(destDir, fileName)
			errors[i] = fd.DownloadFile(url, destPath)
		}(i, url)
	}

	wg.Wait()
	return errors
}
