package test

import (
	"fmt"
	"log"
	"strconv"

	"github.com/karosown/katool-go/ruleengine"
)

// User 示例用户数据结构
type User struct {
	ID       int
	Name     string
	Age      int
	Email    string
	VipLevel int
	Balance  float64
}

// ExampleUsage 演示规则引擎的使用方法
func ExampleUsage() {
	// 创建规则引擎
	engine := ruleengine.NewRuleEngine[User]()

	// 1. 注册基础规则
	engine.RegisterRule("validate_age",
		func(user User, _ any) bool {
			return user.Age > 0 && user.Age < 150
		},
		func(user User, _ any) (User, any, error) {
			if user.Age < 18 {
				return user, "未成年用户", nil
			} else if user.Age >= 60 {
				return user, "老年用户", nil
			}
			return user, "成年用户", nil
		},
	)

	engine.RegisterRule("validate_email",
		func(user User, _ any) bool {
			return user.Email != ""
		},
		func(user User, _ any) (User, any, error) {
			// 简单的邮箱验证
			if len(user.Email) > 5 && fmt.Sprintf("%s", user.Email)[len(user.Email)-4:] == ".com" {
				return user, "邮箱格式正确", nil
			}
			return user, "邮箱格式错误", fmt.Errorf("无效的邮箱格式")
		},
	)

	engine.RegisterRule("calculate_discount",
		func(user User, _ any) bool {
			return user.VipLevel > 0
		},
		func(user User, _ any) (User, any, error) {
			discount := float64(user.VipLevel) * 0.1
			if discount > 0.5 {
				discount = 0.5 // 最大50%折扣
			}
			return user, fmt.Sprintf("VIP折扣: %.0f%%", discount*100), nil
		},
	)

	engine.RegisterRule("check_balance",
		func(user User, _ any) bool {
			return true // 总是执行
		},
		func(user User, _ any) (User, any, error) {
			if user.Balance < 0 {
				return user, "余额不足", fmt.Errorf("用户余额为负")
			} else if user.Balance > 10000 {
				return user, "高价值客户", nil
			}
			return user, "普通余额", nil
		},
	)

	// 2. 添加中间件
	engine.AddMiddleware(func(user User, next func(User) (User, any, error)) (User, any, error) {
		fmt.Printf("[中间件] 开始处理用户: %s (ID: %d)\n", user.Name, user.ID)
		result, data, err := next(user)
		fmt.Printf("[中间件] 处理完成用户: %s\n", user.Name)
		return result, data, err
	})

	// 3. 构建规则链
	// 基础验证链
	_, err := engine.NewBuilder("basic_validation").
		AddRule("validate_age").
		AddRule("validate_email").
		Build()
	if err != nil {
		log.Printf("构建基础验证链失败: %v", err)
		return
	}

	// VIP处理链
	_, err = engine.NewBuilder("vip_processing").
		AddRule("calculate_discount").
		AddRule("check_balance").
		Build()
	if err != nil {
		log.Printf("构建VIP处理链失败: %v", err)
		return
	}

	// 完整处理链
	_, err = engine.NewBuilder("complete_processing").
		AddRule("validate_age").
		AddRule("validate_email").
		AddCustomRule(
			func(user User, _ any) bool {
				return user.Age >= 18 // 只有成年人才能继续
			},
			func(user User, _ any) (User, any, error) {
				return user, "成年验证通过", nil
			},
		).
		AddRule("calculate_discount").
		AddRule("check_balance").
		Build()
	if err != nil {
		log.Printf("构建完整处理链失败: %v", err)
		return
	}

	// 4. 测试数据
	users := []User{
		{ID: 1, Name: "张三", Age: 25, Email: "zhangsan@example.com", VipLevel: 2, Balance: 1500.0},
		{ID: 2, Name: "李四", Age: 17, Email: "lisi@test.com", VipLevel: 0, Balance: 300.0},
		{ID: 3, Name: "王五", Age: 45, Email: "wangwu@company.com", VipLevel: 5, Balance: 15000.0},
		{ID: 4, Name: "赵六", Age: 30, Email: "invalid-email", VipLevel: 1, Balance: -50.0},
	}

	// 5. 执行测试
	fmt.Println("=== 规则引擎使用演示 ===")

	// 显示引擎统计信息
	stats := engine.Stats()
	fmt.Printf("引擎统计: 规则数量=%d, 规则链数量=%d, 中间件数量=%d\n\n",
		stats["rules_count"], stats["chains_count"], stats["middleware_count"])

	// 单个规则链执行
	fmt.Println("--- 单个规则链执行测试 ---")
	for _, user := range users {
		fmt.Printf("\n处理用户: %s\n", user.Name)
		result := engine.Execute("complete_processing", user)
		if result.Error != nil {
			fmt.Printf("❌ 执行失败: %v\n", result.Error)
		} else {
			fmt.Printf("✅ 执行成功: %v\n", result.Result)
		}
	}

	// 批量执行
	fmt.Println("\n--- 批量执行测试 ---")
	user := users[0] // 使用第一个用户测试
	batchResults := engine.BatchExecute([]string{"basic_validation", "vip_processing"}, user)
	for chainName, result := range batchResults {
		fmt.Printf("规则链 [%s]: ", chainName)
		if result.Error != nil {
			fmt.Printf("❌ %v\n", result.Error)
		} else {
			fmt.Printf("✅ %v\n", result.Result)
		}
	}

	// 执行所有规则链
	fmt.Println("\n--- 执行所有规则链 ---")
	allResults := engine.ExecuteAll(users[2]) // 使用第三个用户测试
	for chainName, result := range allResults {
		fmt.Printf("规则链 [%s]: ", chainName)
		if result.Error != nil {
			fmt.Printf("❌ %v\n", result.Error)
		} else {
			fmt.Printf("✅ %v\n", result.Result)
		}
	}

	// 动态管理
	fmt.Println("\n--- 动态管理测试 ---")
	fmt.Printf("当前规则: %v\n", engine.ListRules())
	fmt.Printf("当前规则链: %v\n", engine.ListChains())

	// 添加新规则
	engine.RegisterRule("new_rule",
		func(user User, _ any) bool { return true },
		func(user User, _ any) (User, any, error) {
			return user, "新规则执行", nil
		},
	)
	fmt.Printf("添加新规则后: %v\n", engine.ListRules())
}

// ExampleAdvanced 高级用法示例
func ExampleAdvanced() {
	engine := ruleengine.NewRuleEngine[map[string]interface{}]()

	// 数据转换规则
	engine.RegisterRule("string_to_int",
		func(data map[string]interface{}, _ any) bool {
			_, exists := data["number_string"]
			return exists
		},
		func(data map[string]interface{}, _ any) (map[string]interface{}, any, error) {
			if str, ok := data["number_string"].(string); ok {
				if num, err := strconv.Atoi(str); err == nil {
					data["number"] = num
					return data, "字符串转数字成功", nil
				}
			}
			return data, nil, fmt.Errorf("字符串转数字失败")
		},
	)

	// 数据验证规则
	engine.RegisterRule("validate_range",
		func(data map[string]interface{}, _ any) bool {
			_, exists := data["number"]
			return exists
		},
		func(data map[string]interface{}, _ any) (map[string]interface{}, any, error) {
			if num, ok := data["number"].(int); ok {
				if num >= 1 && num <= 100 {
					data["valid"] = true
					return data, "数值在有效范围内", nil
				}
				data["valid"] = false
				return data, "数值超出范围", nil
			}
			return data, nil, fmt.Errorf("数据类型错误")
		},
	)

	// 条件分支示例
	trueChain := []*ruleengine.RuleNode[map[string]interface{}]{
		ruleengine.NewRuleNode(
			func(data map[string]interface{}, _ any) bool { return true },
			func(data map[string]interface{}, _ any) (map[string]interface{}, any, error) {
				data["result"] = "有效数据处理"
				return data, "有效分支执行", nil
			},
		),
	}

	falseChain := []*ruleengine.RuleNode[map[string]interface{}]{
		ruleengine.NewRuleNode(
			func(data map[string]interface{}, _ any) bool { return true },
			func(data map[string]interface{}, _ any) (map[string]interface{}, any, error) {
				data["result"] = "无效数据处理"
				return data, "无效分支执行", nil
			},
		),
	}

	// 构建带条件分支的规则链
	_, err := engine.NewBuilder("data_processing").
		AddRule("string_to_int").
		AddRule("validate_range").
		AddConditionalChain(
			func(data map[string]interface{}, _ any) bool {
				if valid, exists := data["valid"]; exists {
					return valid.(bool)
				}
				return false
			},
			trueChain,
			falseChain,
		).
		Build()

	if err != nil {
		log.Printf("构建数据处理链失败: %v", err)
		return
	}

	// 测试数据
	testData := []map[string]interface{}{
		{"number_string": "50"},
		{"number_string": "150"},
		{"number_string": "abc"},
		{"number_string": "25"},
	}

	fmt.Println("\n=== 高级用法演示 ===")
	for i, data := range testData {
		fmt.Printf("\n测试数据 %d: %v\n", i+1, data)
		result := engine.Execute("data_processing", data)
		fmt.Printf("处理结果: %v\n", result.Data)
		if result.Error != nil {
			fmt.Printf("错误: %v\n", result.Error)
		} else {
			fmt.Printf("状态: %v\n", result.Result)
		}
	}
}
