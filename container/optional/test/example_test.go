package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/karosown/katool-go/container/optional"
)

func TestOptionalUsageExamples(t *testing.T) {
	// 基础用法示例
	testBasicUsage(t)

	// 字符串处理示例
	testStringProcessing(t)

	// 链式操作示例
	testChainOperations(t)
}

func testBasicUsage(t *testing.T) {
	fmt.Println("=== 基础用法示例 ===")

	// 创建Optional
	opt := optional.Of("Hello World")
	emptyOpt := optional.Empty[string]()

	// 安全检查和获取
	if opt.IsPresent() {
		value := opt.Get()
		fmt.Println("值存在:", value)
	}

	// 提供默认值
	defaultValue := emptyOpt.OrElse("默认值")
	fmt.Println("默认值:", defaultValue)

	// 条件执行
	opt.IfPresent(func(v string) {
		fmt.Println("处理值:", v)
	})

	fmt.Println()
}

func testStringProcessing(t *testing.T) {
	fmt.Println("=== 字符串处理示例 ===")

	// 字符串处理链式操作 - 方法1：使用MapTyped
	input := "  hello  "
	result1 := optional.MapTyped(optional.Of(input), strings.TrimSpace).
		Filter(func(s string) bool { return len(s) > 0 }).
		OrElse("空字符串")
	fmt.Println("处理结果1:", result1)

	// 字符串处理链式操作 - 方法2：分步处理
	opt := optional.Of(input)

	// 先处理空格
	var trimmed optional.Optional[string]
	if opt.IsPresent() {
		trimmed = optional.Of(strings.TrimSpace(opt.Get()))
	} else {
		trimmed = optional.Empty[string]()
	}

	// 再过滤
	filtered := trimmed.Filter(func(s string) bool { return len(s) > 0 })

	// 获取结果
	result2 := filtered.OrElse("空字符串")
	fmt.Println("处理结果2:", result2)

	fmt.Println()
}

func testChainOperations(t *testing.T) {
	fmt.Println("=== 链式操作示例 ===")

	// 工具函数示例
	condition := true
	enabled := optional.IsTrue(condition, "启用", "禁用")
	fmt.Println("状态:", enabled)

	enabledFunc := func() string { return "功能已启用" }
	disabledFunc := func() string { return "功能已禁用" }
	result := optional.IsTrueByFunc(condition, enabledFunc, disabledFunc)
	fmt.Println("功能状态:", result)

	// 复杂的Optional操作
	userInput := "  USER123  "
	processedUser := processUserInput(userInput)
	fmt.Println("处理后的用户输入:", processedUser)

	fmt.Println()
}

// 辅助函数：处理用户输入
func processUserInput(input string) string {
	// 创建Optional
	opt := optional.Of(input)

	// 去除空格
	trimmed := opt.Map(func(s string) any {
		return strings.TrimSpace(s)
	})

	// 转换为小写并检查长度
	result := "无效输入"
	if trimmed.IsPresent() {
		if str, ok := trimmed.Get().(string); ok && len(str) > 0 {
			result = strings.ToLower(str)
		}
	}

	return result
}

// 专门为字符串设计的Optional操作
type StringOptional struct {
	optional.Optional[string]
}

// NewStringOptional 创建字符串Optional
func NewStringOptional(value string) StringOptional {
	return StringOptional{optional.Of(value)}
}

// EmptyStringOptional 创建空的字符串Optional
func EmptyStringOptional() StringOptional {
	return StringOptional{optional.Empty[string]()}
}

// TrimSpace 去除字符串两端空格
func (s StringOptional) TrimSpace() StringOptional {
	if !s.IsPresent() {
		return EmptyStringOptional()
	}
	return NewStringOptional(strings.TrimSpace(s.Get()))
}

// FilterNonEmpty 过滤空字符串
func (s StringOptional) FilterNonEmpty() StringOptional {
	return StringOptional{s.Filter(func(str string) bool { return len(str) > 0 })}
}

func TestStringOptional(t *testing.T) {
	fmt.Println("=== 专用字符串Optional示例 ===")

	// 使用专用的字符串Optional
	result := NewStringOptional("  hello  ").
		TrimSpace().
		FilterNonEmpty().
		OrElse("空字符串")

	fmt.Println("字符串处理结果:", result)
}

func TestIsNillable(t *testing.T) {
	A := struct {
		A string
	}{}
	print(optional.ToNillable(A) == nil)
}
