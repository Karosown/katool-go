package test

import (
	"fmt"
	"log"
	"testing"
)

// 由于Go版本问题，这里提供一个独立的演示
// 当Go版本问题解决后，可以直接运行: go run demo.go

func Test_Demo(t *testing.T) {
	fmt.Println("=== 规则引擎演示 ===")
	fmt.Println("注意：当前由于Go版本不匹配无法编译，请先解决Go版本问题")
	fmt.Println()

	// 演示基本用法
	fmt.Println("1. 基本用法示例:")
	fmt.Println(`
	// 创建规则引擎
	engine := NewRuleEngine[int]()
	
	// 注册规则
	engine.RegisterRule("add_ten",
		func(data int, _ any) bool { return data > 0 },
		func(data int, _ any) (int, any, error) {
			return data + 10, "添加了10", nil
		},
	)
	
	// 构建规则链
	_, err := engine.NewBuilder("simple_chain").
		AddRule("add_ten").
		Build()
	
	// 执行规则
	result := engine.Execute("simple_chain", 5)
	// 结果: 15, 信息: "添加了10"
	`)

	fmt.Println("\n2. 复杂用法示例:")
	fmt.Println(`
	type User struct {
		ID       int
		Name     string
		Age      int
		VipLevel int
		Balance  float64
	}
	
	engine := NewRuleEngine[User]()
	
	// 注册多个规则
	engine.RegisterRule("validate_age", ...)
	engine.RegisterRule("calculate_discount", ...)
	engine.RegisterRule("check_balance", ...)
	
	// 构建规则链
	engine.NewBuilder("user_processing").
		AddRule("validate_age").
		AddRule("calculate_discount").
		AddRule("check_balance").
		Build()
	
	// 批量执行
	results := engine.BatchExecute(chainNames, userData)
	`)

	fmt.Println("\n3. 解决Go版本问题的方法:")
	fmt.Println("- 方法1: 重新安装Go 1.23.1")
	fmt.Println("- 方法2: 设置正确的GOROOT")
	fmt.Println("- 方法3: 清理并重建Go工具链")
	fmt.Println()

	fmt.Println("解决后运行命令:")
	fmt.Println("cd ruleengine")
	fmt.Println("go test -v                    # 运行测试")
	fmt.Println("go test -run TestExample      # 运行示例")
	fmt.Println("go test -bench=.              # 性能测试")

	// 模拟规则引擎的核心逻辑（不依赖其他文件）
	fmt.Println("\n=== 核心逻辑演示 ===")
	demoBasicLogic()
}

// 演示核心逻辑（不依赖其他文件）
func demoBasicLogic() {
	// 模拟用户数据
	type User struct {
		Name string
		Age  int
		VIP  bool
	}

	users := []User{
		{"张三", 25, true},
		{"李四", 17, false},
		{"王五", 45, true},
	}

	// 模拟规则处理
	for _, user := range users {
		fmt.Printf("\n处理用户: %s\n", user.Name)

		// 年龄验证规则
		if user.Age >= 18 {
			fmt.Printf("✅ 年龄验证通过: %d岁\n", user.Age)
		} else {
			fmt.Printf("❌ 年龄验证失败: %d岁 (未成年)\n", user.Age)
			continue
		}

		// VIP折扣规则
		if user.VIP {
			fmt.Printf("✅ VIP用户，享受9折优惠\n")
		} else {
			fmt.Printf("ℹ️  普通用户，无折扣\n")
		}

		// 最终结果
		status := "处理完成"
		if user.Age < 18 {
			status = "处理失败（未成年）"
		}
		fmt.Printf("📊 最终状态: %s\n", status)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
