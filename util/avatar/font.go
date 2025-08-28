package avatar

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FontSource 字体来源类型
type FontSource int

const (
	FontSourceSystem FontSource = iota // 系统字体
	FontSourceFile                     // 本地文件
	FontSourceURL                      // 网络URL
	FontSourceData                     // 原始数据
	FontSourceGoogle                   // Google Fonts
)

// FontConfig 字体配置
type FontConfig struct {
	FontFamily   string            // 字体族名称
	FontWeight   string            // 字体粗细 (normal, bold, 100-900)
	FontStyle    string            // 字体样式 (normal, italic, oblique)
	FontSource   FontSource        // 字体来源
	FontData     []byte            // 字体文件数据
	FontFormat   string            // 字体格式 (woff, woff2, truetype, opentype)
	FontPath     string            // 字体文件路径
	FontURL      string            // 字体URL
	EmbedFont    bool              // 是否嵌入字体到SVG中
	Fallback     string            // 后备字体
	GoogleParams map[string]string // Google Fonts参数
}

// FontRegistry 字体注册表
type FontRegistry struct {
	fonts map[string]*FontConfig
}

// NewFontRegistry 创建字体注册表
func NewFontRegistry() *FontRegistry {
	return &FontRegistry{
		fonts: make(map[string]*FontConfig),
	}
}

// RegisterFont 注册字体
func (fr *FontRegistry) RegisterFont(name string, config *FontConfig) {
	fr.fonts[name] = config
}

// GetFont 获取字体配置
func (fr *FontRegistry) GetFont(name string) (*FontConfig, bool) {
	config, exists := fr.fonts[name]
	return config, exists
}

// ListFonts 列出所有字体
func (fr *FontRegistry) ListFonts() []string {
	var fonts []string
	for name := range fr.fonts {
		fonts = append(fonts, name)
	}
	return fonts
}

// 全局字体注册表
var GlobalFontRegistry = NewFontRegistry()

// 预定义字体配置
func init() {
	// 系统字体
	GlobalFontRegistry.RegisterFont("default", &FontConfig{
		FontFamily: "Arial, 'Microsoft YaHei', 'SimHei', sans-serif",
		FontWeight: "normal",
		FontStyle:  "normal",
		FontSource: FontSourceSystem,
		Fallback:   "sans-serif",
	})

	GlobalFontRegistry.RegisterFont("chinese", &FontConfig{
		FontFamily: "'Microsoft YaHei', 'PingFang SC', 'Hiragino Sans GB', 'Source Han Sans CN', sans-serif",
		FontWeight: "normal",
		FontStyle:  "normal",
		FontSource: FontSourceSystem,
		Fallback:   "sans-serif",
	})

	GlobalFontRegistry.RegisterFont("serif", &FontConfig{
		FontFamily: "Georgia, 'Times New Roman', 'SimSun', 'Source Han Serif CN', serif",
		FontWeight: "normal",
		FontStyle:  "normal",
		FontSource: FontSourceSystem,
		Fallback:   "serif",
	})

	GlobalFontRegistry.RegisterFont("monospace", &FontConfig{
		FontFamily: "'Courier New', 'Consolas', 'Monaco', 'Source Code Pro', monospace",
		FontWeight: "normal",
		FontStyle:  "normal",
		FontSource: FontSourceSystem,
		Fallback:   "monospace",
	})

	// Google Fonts 示例
	GlobalFontRegistry.RegisterFont("roboto", &FontConfig{
		FontFamily: "'Roboto', sans-serif",
		FontWeight: "400",
		FontStyle:  "normal",
		FontSource: FontSourceGoogle,
		FontURL:    "https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap",
		Fallback:   "sans-serif",
		GoogleParams: map[string]string{
			"family":  "Roboto:wght@400;700",
			"display": "swap",
		},
	})

	// Google Sans 风格
	GlobalFontRegistry.RegisterFont("google-sans", &FontConfig{
		FontFamily: "'Google Sans', 'Product Sans', 'Roboto', 'Helvetica Neue', Arial, sans-serif",
		FontWeight: "500", // Google使用中等粗细
		FontStyle:  "normal",
		FontSource: FontSourceGoogle,
		FontURL:    "https://fonts.googleapis.com/css2?family=Google+Sans:wght@400;500;600&display=swap",
		Fallback:   "sans-serif",
		GoogleParams: map[string]string{
			"family":  "Google+Sans:wght@400;500;600",
			"display": "swap",
		},
	})

	// Roboto 风格（Android/旧版Google服务）
	GlobalFontRegistry.RegisterFont("old-roboto", &FontConfig{
		FontFamily: "'Roboto', 'Arial', sans-serif",
		FontWeight: "500",
		FontStyle:  "normal",
		FontSource: FontSourceGoogle,
		FontURL:    "https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap",
		Fallback:   "sans-serif",
		GoogleParams: map[string]string{
			"family":  "Roboto:wght@400;500;700",
			"display": "swap",
		},
	})

	// Product Sans（Google品牌字体，较难获取）
	GlobalFontRegistry.RegisterFont("product-sans", &FontConfig{
		FontFamily: "'Product Sans', 'Google Sans', 'Roboto', Arial, sans-serif",
		FontWeight: "400",
		FontStyle:  "normal",
		FontSource: FontSourceSystem,
		Fallback:   "sans-serif",
	})
}

// FontBuilder 字体构建器
type FontBuilder struct {
	config *FontConfig
}

// NewFontBuilder 创建字体构建器
func NewFontBuilder() *FontBuilder {
	return &FontBuilder{
		config: &FontConfig{
			FontWeight: "normal",
			FontStyle:  "normal",
			EmbedFont:  false,
			Fallback:   "sans-serif",
		},
	}
}

// SetFamily 设置字体族
func (fb *FontBuilder) SetFamily(family string) *FontBuilder {
	fb.config.FontFamily = family
	return fb
}

// SetWeight 设置字体粗细
func (fb *FontBuilder) SetWeight(weight string) *FontBuilder {
	fb.config.FontWeight = weight
	return fb
}

// SetStyle 设置字体样式
func (fb *FontBuilder) SetStyle(style string) *FontBuilder {
	fb.config.FontStyle = style
	return fb
}

// SetFallback 设置后备字体
func (fb *FontBuilder) SetFallback(fallback string) *FontBuilder {
	fb.config.Fallback = fallback
	return fb
}

// EnableEmbed 启用字体嵌入
func (fb *FontBuilder) EnableEmbed() *FontBuilder {
	fb.config.EmbedFont = true
	return fb
}

// FromFile 从文件加载字体
func (fb *FontBuilder) FromFile(filePath string) *FontBuilder {
	fb.config.FontSource = FontSourceFile
	fb.config.FontPath = filePath

	// 根据扩展名确定格式
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".woff":
		fb.config.FontFormat = "woff"
	case ".woff2":
		fb.config.FontFormat = "woff2"
	case ".ttf":
		fb.config.FontFormat = "truetype"
	case ".otf":
		fb.config.FontFormat = "opentype"
	default:
		fb.config.FontFormat = "truetype"
	}

	return fb
}

// FromURL 从URL加载字体
func (fb *FontBuilder) FromURL(url string) *FontBuilder {
	fb.config.FontSource = FontSourceURL
	fb.config.FontURL = url
	return fb
}

// FromData 从数据加载字体
func (fb *FontBuilder) FromData(data []byte, format string) *FontBuilder {
	fb.config.FontSource = FontSourceData
	fb.config.FontData = data
	fb.config.FontFormat = format
	return fb
}

// FromGoogleFonts 从Google Fonts加载字体
func (fb *FontBuilder) FromGoogleFonts(family string, weights ...string) *FontBuilder {
	fb.config.FontSource = FontSourceGoogle
	fb.config.FontFamily = fmt.Sprintf("'%s', sans-serif", family)

	weightStr := "400"
	if len(weights) > 0 {
		weightStr = strings.Join(weights, ";")
	}

	fb.config.FontURL = fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:wght@%s&display=swap",
		strings.ReplaceAll(family, " ", "+"), weightStr)

	fb.config.GoogleParams = map[string]string{
		"family":  fmt.Sprintf("%s:wght@%s", family, weightStr),
		"display": "swap",
	}

	return fb
}

// Build 构建字体配置
func (fb *FontBuilder) Build() *FontConfig {
	// 加载字体数据（如果需要）
	if fb.config.EmbedFont && fb.config.FontSource == FontSourceFile && len(fb.config.FontData) == 0 {
		if data, err := os.ReadFile(fb.config.FontPath); err == nil {
			fb.config.FontData = data
		}
	}

	if fb.config.EmbedFont && fb.config.FontSource == FontSourceURL && len(fb.config.FontData) == 0 {
		if data, err := downloadFont(fb.config.FontURL); err == nil {
			fb.config.FontData = data
		}
	}

	return fb.config
}

// downloadFont 下载字体文件
func downloadFont(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载字体失败，状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
