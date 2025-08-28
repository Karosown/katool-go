package avatar

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/adler32"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

// AvatarResult 头像生成结果
type AvatarResult struct {
	Base64Data string     // Base64数据URL
	URLEncoded string     // URL编码版本
	RawBase64  string     // 纯Base64数据
	File       FileResult // 文件对象
	FileName   string     // 文件名
	SVGContent string     // SVG内容
	IsMemory   bool       // 是否为内存文件
}

// GetDataURL 获取不同格式的data URL
func (ar *AvatarResult) GetDataURL(format string) string {
	switch format {
	case "base64":
		return ar.Base64Data
	case "url":
		return ar.URLEncoded
	case "raw":
		return ar.RawBase64
	default:
		return ar.Base64Data
	}
}

// ValidateBase64 验证Base64数据是否有效
func (ar *AvatarResult) ValidateBase64() error {
	if ar.RawBase64 == "" {
		return fmt.Errorf("Base64数据为空")
	}

	decoded, err := base64.StdEncoding.DecodeString(ar.RawBase64)
	if err != nil {
		return fmt.Errorf("Base64解码失败: %v", err)
	}

	svgStr := string(decoded)
	if !strings.HasPrefix(svgStr, "<svg") {
		return fmt.Errorf("解码后的内容不是有效的SVG")
	}

	return nil
}

// NewAvatarGenerator 创建头像生成器
func NewAvatarGenerator() *AvatarGenerator {
	defaultFont, _ := GlobalFontRegistry.GetFont("default")
	return &AvatarGenerator{
		Size:         100,
		FontSize:     50,
		TextColor:    "#ffffff",
		BgColor:      "rgb(255,250,250)",
		OutputDir:    "",
		UseHSL:       false,
		FontConfig:   defaultFont,
		FontRegistry: GlobalFontRegistry,
	}
}

// SetSize 设置头像尺寸
func (ag *AvatarGenerator) SetSize(size int) *AvatarGenerator {
	ag.Size = size
	ag.FontSize = size / 2
	return ag
}

// SetColors 设置颜色
func (ag *AvatarGenerator) SetColors(textColor, bgColor string) *AvatarGenerator {
	ag.TextColor = textColor
	ag.BgColor = bgColor
	return ag
}

// SetOutputDir 设置输出目录
func (ag *AvatarGenerator) SetOutputDir(dir string) *AvatarGenerator {
	ag.OutputDir = dir
	return ag
}

// SetFont 设置字体（使用注册的字体）
func (ag *AvatarGenerator) SetFont(fontName string) *AvatarGenerator {
	if config, exists := ag.FontRegistry.GetFont(fontName); exists {
		ag.FontConfig = config
	}
	return ag
}

// SetCustomFont 设置自定义字体配置
func (ag *AvatarGenerator) SetCustomFont(config *FontConfig) *AvatarGenerator {
	ag.FontConfig = config
	return ag
}

// SetFontBuilder 使用字体构建器设置字体
func (ag *AvatarGenerator) SetFontBuilder(builder *FontBuilder) *AvatarGenerator {
	ag.FontConfig = builder.Build()
	return ag
}

// EnableHSLColor 启用HSL颜色
func (ag *AvatarGenerator) EnableHSLColor() *AvatarGenerator {
	ag.UseHSL = true
	return ag
}

// generateColors 根据文本生成颜色
func (ag *AvatarGenerator) generateColors(text string) (string, string) {
	if !ag.UseHSL {
		return ag.TextColor, ag.BgColor
	}

	hasher := adler32.New()
	hasher.Write([]byte(text))
	total := hasher.Sum32()

	hue := total % 360
	bgColor := fmt.Sprintf("hsl(%d, 60%%, 75%%)", hue)
	textColor := "#ffffff"

	if hue > 40 && hue < 200 {
		textColor = "#333333"
	}

	return textColor, bgColor
}

// getFirstCharUpper 获取字符串的第一个字符并转为大写
func (ag *AvatarGenerator) getFirstCharUpper(text string) string {
	if text == "" {
		return "?"
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return "?"
	}

	r, _ := utf8.DecodeRuneInString(text)
	return strings.ToUpper(string(r))
}

// generateSVG 生成SVG内容（修复版本 - 确保格式正确）
func (ag *AvatarGenerator) generateSVG(text string) string {
	textColor, bgColor := ag.generateColors(text)
	first := ag.getFirstCharUpper(text)
	fontStyle := ag.generateFontStyle()
	fontDefs := ag.generateFontDefs()

	// 生成紧凑的SVG，避免不必要的空格和换行
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" version="1.1" height="%d" width="%d">%s<rect fill="%s" x="0" y="0" width="%d" height="%d"/><text x="%d" y="%d" font-size="%d" %s fill="%s" text-anchor="middle" dominant-baseline="central">%s</text></svg>`,
		ag.Size, ag.Size,
		fontDefs,
		bgColor, ag.Size, ag.Size,
		ag.Size/2, ag.Size/2, ag.FontSize, fontStyle, textColor, first)

	return svg
}

// generateFileName 生成文件名
func (ag *AvatarGenerator) generateFileName(text string, customName ...string) string {
	if len(customName) > 0 && customName[0] != "" {
		fileName := customName[0]
		if !strings.HasSuffix(fileName, ".svg") {
			fileName += ".svg"
		}
		return fileName
	}

	hasher := adler32.New()
	hasher.Write([]byte(text + time.Now().String()))
	return fmt.Sprintf("avatar_%d.svg", hasher.Sum32())
}

// GenerateAvatar 生成头像（修复Base64格式）
func (ag *AvatarGenerator) GenerateAvatar(text string) (*AvatarResult, error) {
	svgContent := ag.generateSVG(text)

	// 生成Base64数据
	base64Data := base64.StdEncoding.EncodeToString([]byte(svgContent))
	dataURL := "data:image/svg+xml;base64," + base64Data

	// 生成URL编码版本作为备选
	urlEncodedSVG := strings.ReplaceAll(svgContent, "#", "%23")
	urlEncodedSVG = strings.ReplaceAll(urlEncodedSVG, " ", "%20")
	urlEncodedSVG = strings.ReplaceAll(urlEncodedSVG, "<", "%3C")
	urlEncodedSVG = strings.ReplaceAll(urlEncodedSVG, ">", "%3E")
	urlEncodedSVG = strings.ReplaceAll(urlEncodedSVG, "\"", "%22")
	dataURLEncoded := "data:image/svg+xml;charset=utf-8," + urlEncodedSVG

	return &AvatarResult{
		Base64Data: dataURL,
		URLEncoded: dataURLEncoded,
		RawBase64:  base64Data,
		SVGContent: svgContent,
	}, nil
}

// GenerateAvatarFile 生成头像文件
func (ag *AvatarGenerator) GenerateAvatarFile(text string, filename ...string) (*AvatarResult, error) {
	// 先生成基础数据
	result, err := ag.GenerateAvatar(text)
	if err != nil {
		return nil, err
	}

	fileName := ag.generateFileName(text, filename...)
	result.FileName = fileName

	// 内存模式
	if ag.OutputDir == "" {
		memFile := NewMemoryFile(fileName, []byte(result.SVGContent))
		result.File = memFile
		result.IsMemory = true
		return result, nil
	}

	// 磁盘模式
	if err := os.MkdirAll(ag.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %v", err)
	}

	filePath := filepath.Join(ag.OutputDir, fileName)
	if err := os.WriteFile(filePath, []byte(result.SVGContent), 0644); err != nil {
		return nil, fmt.Errorf("写入文件失败: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	result.File = &osFileWrapper{file}
	result.IsMemory = false

	return result, nil
}

type osFileWrapper struct {
	*os.File
}

func (ofw *osFileWrapper) Content() []byte {
	content, _ := io.ReadAll(ofw.File)
	ofw.File.Seek(0, 0)
	return content
}

// GenerateMultipartFile 生成multipart文件
func (ag *AvatarGenerator) GenerateMultipartFile(text string, fieldName string, filename ...string) (*bytes.Buffer, string, error) {
	result, err := ag.GenerateAvatarFile(text, filename...)
	if err != nil {
		return nil, "", err
	}
	defer result.File.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, result.FileName)
	if err != nil {
		return nil, "", fmt.Errorf("创建multipart字段失败: %v", err)
	}

	if _, err := io.Copy(part, result.File); err != nil {
		return nil, "", fmt.Errorf("写入multipart内容失败: %v", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("关闭multipart writer失败: %v", err)
	}

	return &buf, writer.FormDataContentType(), nil
}

// BatchGenerate 批量生成头像
func (ag *AvatarGenerator) BatchGenerate(texts []string) ([]*AvatarResult, error) {
	results := make([]*AvatarResult, 0, len(texts))

	for _, text := range texts {
		result, err := ag.GenerateAvatarFile(text)
		if err != nil {
			return nil, fmt.Errorf("生成头像失败 [%s]: %v", text, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GenerateRandomAvatar 生成随机头像
func (ag *AvatarGenerator) GenerateRandomAvatar() (*AvatarResult, error) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 1)
	rand.Read(b)
	randomChar := string(chars[int(b[0])%len(chars)])

	return ag.GenerateAvatarFile(randomChar)
}

// generateFontDefs 生成字体定义（修复版本 - 移除无效的外部字体导入）
func (ag *AvatarGenerator) generateFontDefs() string {
	if ag.FontConfig == nil {
		return ""
	}

	var defs strings.Builder

	// 只处理嵌入字体数据，移除外部字体导入
	if ag.FontConfig.EmbedFont && len(ag.FontConfig.FontData) > 0 {
		fontBase64 := base64.StdEncoding.EncodeToString(ag.FontConfig.FontData)

		defs.WriteString("<defs><style type=\"text/css\">")
		defs.WriteString("@font-face {")
		defs.WriteString(fmt.Sprintf("font-family: '%s';", ag.FontConfig.FontFamily))
		defs.WriteString(fmt.Sprintf("font-weight: %s;", ag.FontConfig.FontWeight))
		defs.WriteString(fmt.Sprintf("font-style: %s;", ag.FontConfig.FontStyle))
		defs.WriteString(fmt.Sprintf("src: url(data:font/%s;base64,%s) format('%s');",
			ag.FontConfig.FontFormat, fontBase64, ag.FontConfig.FontFormat))
		defs.WriteString("}")
		defs.WriteString("</style></defs>")
	}

	// 移除外部字体导入部分，因为在SVG文件和Base64中无效
	// 注释掉原来的外部字体导入代码
	/*
		if ag.FontConfig.FontSource == FontSourceGoogle && ag.FontConfig.FontURL != "" && !ag.FontConfig.EmbedFont {
			defs.WriteString("<defs><style type=\"text/css\">")
			defs.WriteString(fmt.Sprintf("@import url('%s');", ag.FontConfig.FontURL))
			defs.WriteString("</style></defs>")
		}
	*/

	return defs.String()
}

// generateFontStyle 生成字体样式（优化版本 - 使用系统字体后备）
func (ag *AvatarGenerator) generateFontStyle() string {
	if ag.FontConfig == nil {
		return `font-family="Arial, sans-serif"`
	}

	var styles []string

	// 处理字体族 - 如果是外部字体且未嵌入，则使用后备字体
	fontFamily := ag.FontConfig.FontFamily
	if ag.FontConfig.FontSource == FontSourceGoogle && !ag.FontConfig.EmbedFont {
		// 对于Google字体，如果没有嵌入，使用系统后备字体
		fontFamily = ag.getSystemFallbackFont(ag.FontConfig.FontFamily)
	}

	if fontFamily != "" {
		styles = append(styles, fmt.Sprintf(`font-family="%s"`, fontFamily))
	}

	if ag.FontConfig.FontWeight != "" && ag.FontConfig.FontWeight != "normal" {
		styles = append(styles, fmt.Sprintf(`font-weight="%s"`, ag.FontConfig.FontWeight))
	}

	if ag.FontConfig.FontStyle != "" && ag.FontConfig.FontStyle != "normal" {
		styles = append(styles, fmt.Sprintf(`font-style="%s"`, ag.FontConfig.FontStyle))
	}

	return strings.Join(styles, " ")
}

// getSystemFallbackFont 获取系统后备字体
func (ag *AvatarGenerator) getSystemFallbackFont(originalFont string) string {
	// 根据原始字体返回合适的系统后备字体
	originalFont = strings.ToLower(originalFont)

	if strings.Contains(originalFont, "roboto") {
		return "Arial, 'Helvetica Neue', Helvetica, sans-serif"
	}
	if strings.Contains(originalFont, "open sans") {
		return "Arial, 'Helvetica Neue', Helvetica, sans-serif"
	}
	if strings.Contains(originalFont, "lato") {
		return "Arial, 'Helvetica Neue', Helvetica, sans-serif"
	}
	if strings.Contains(originalFont, "montserrat") {
		return "'Helvetica Neue', Helvetica, Arial, sans-serif"
	}
	if strings.Contains(originalFont, "source sans") {
		return "Arial, 'Helvetica Neue', Helvetica, sans-serif"
	}
	if strings.Contains(originalFont, "poppins") {
		return "'Helvetica Neue', Helvetica, Arial, sans-serif"
	}
	if strings.Contains(originalFont, "inter") {
		return "'SF Pro Display', -apple-system, BlinkMacSystemFont, Arial, sans-serif"
	}

	// 默认后备
	if ag.FontConfig.Fallback != "" {
		return ag.FontConfig.Fallback
	}

	return "Arial, 'Helvetica Neue', Helvetica, sans-serif"
}
