package test

import (
	"fmt"
	"testing"

	"github.com/karosown/katool-go/ruleengine"
)

// TestNewRuleEngine 测试规则引擎创建
func TestNewRuleEngine(t *testing.T) {
	engine := ruleengine.NewRuleEngine[string]()

	if engine == nil {
		t.Fatal("规则引擎创建失败")
	}

	stats := engine.Stats()
	if stats["rules_count"] != 0 {
		t.Errorf("新创建的引擎应该有0个规则，实际有 %v", stats["rules_count"])
	}

	if stats["chains_count"] != 0 {
		t.Errorf("新创建的引擎应该有0个规则链，实际有 %v", stats["chains_count"])
	}
}

// TestRegisterRule 测试规则注册
func TestRegisterRule(t *testing.T) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("test_rule",
		func(data int, _ any) bool { return data > 0 },
		func(data int, _ any) (int, any, error) { return data * 2, "doubled", nil },
	)

	// 验证规则存在
	rule, exists := engine.GetRule("test_rule")
	if !exists {
		t.Fatal("规则注册失败")
	}

	if rule == nil {
		t.Fatal("获取的规则为nil")
	}

	// 验证统计信息
	stats := engine.Stats()
	if stats["rules_count"] != 1 {
		t.Errorf("应该有1个规则，实际有 %v", stats["rules_count"])
	}
}

// TestRuleBuilder 测试规则构建器
func TestRuleBuilder(t *testing.T) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("add_one",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data + 1, "added_one", nil },
	)

	engine.RegisterRule("multiply_two",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data * 2, "multiplied_two", nil },
	)

	// 构建规则链
	tree, err := engine.NewBuilder("test_chain").
		AddRule("add_one").
		AddRule("multiply_two").
		Build()

	if err != nil {
		t.Fatalf("构建规则链失败: %v", err)
	}

	if tree == nil {
		t.Fatal("构建的规则树为nil")
	}

	// 验证规则链存在
	chains := engine.ListChains()
	found := false
	for _, chain := range chains {
		if chain == "test_chain" {
			found = true
			break
		}
	}

	if !found {
		t.Error("规则链未正确注册")
	}
}

// TestExecute 测试规则执行
func TestExecute(t *testing.T) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("add_ten",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data + 10, "added_ten", nil },
	)

	// 构建规则链
	_, err := engine.NewBuilder("simple_chain").
		AddRule("add_ten").
		Build()

	if err != nil {
		t.Fatalf("构建规则链失败: %v", err)
	}

	// 执行规则
	result := engine.Execute("simple_chain", 5)

	if result.Error != nil {
		t.Fatalf("规则执行失败: %v", result.Error)
	}

	if result.Data != 15 {
		t.Errorf("期望结果为15，实际为 %v", result.Data)
	}

	if result.Result != "added_ten" {
		t.Errorf("期望结果信息为'added_ten'，实际为 %v", result.Result)
	}
}

// TestMiddleware 测试中间件
func TestMiddleware(t *testing.T) {
	engine := ruleengine.NewRuleEngine[string]()
	middlewareExecuted := false

	// 添加中间件
	engine.AddMiddleware(func(data string, next func(string) (string, any, error)) (string, any, error) {
		middlewareExecuted = true
		return next(data + "_middleware")
	})

	// 注册规则
	engine.RegisterRule("append_rule",
		func(data string, _ any) bool { return true },
		func(data string, _ any) (string, any, error) { return data + "_rule", "appended", nil },
	)

	// 构建规则链
	_, err := engine.NewBuilder("middleware_chain").
		AddRule("append_rule").
		Build()

	if err != nil {
		t.Fatalf("构建规则链失败: %v", err)
	}

	// 执行规则
	result := engine.Execute("middleware_chain", "test")

	if result.Error != nil {
		t.Fatalf("规则执行失败: %v", result.Error)
	}

	if !middlewareExecuted {
		t.Error("中间件未执行")
	}

	if result.Data != "test_middleware_rule" {
		t.Errorf("期望结果为'test_middleware_rule'，实际为 %v", result.Data)
	}
}

// TestBatchExecute 测试批量执行
func TestBatchExecute(t *testing.T) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("add_one",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data + 1, "added_one", nil },
	)

	engine.RegisterRule("multiply_two",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data * 2, "multiplied_two", nil },
	)

	// 构建两个规则链
	_, err1 := engine.NewBuilder("chain1").AddRule("add_one").Build()
	_, err2 := engine.NewBuilder("chain2").AddRule("multiply_two").Build()

	if err1 != nil || err2 != nil {
		t.Fatalf("构建规则链失败: %v, %v", err1, err2)
	}

	// 批量执行
	results := engine.BatchExecute([]string{"chain1", "chain2"}, 5)

	if len(results) != 2 {
		t.Errorf("期望2个结果，实际得到 %v", len(results))
	}

	if results["chain1"].Data != 6 {
		t.Errorf("chain1期望结果为6，实际为 %v", results["chain1"].Data)
	}

	if results["chain2"].Data != 10 {
		t.Errorf("chain2期望结果为10，实际为 %v", results["chain2"].Data)
	}
}

// TestRemoveRule 测试规则移除
func TestRemoveRule(t *testing.T) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("temp_rule",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data, "temp", nil },
	)

	// 验证规则存在
	_, exists := engine.GetRule("temp_rule")
	if !exists {
		t.Fatal("规则注册失败")
	}

	// 移除规则
	removed := engine.RemoveRule("temp_rule")
	if !removed {
		t.Error("规则移除失败")
	}

	// 验证规则不存在
	_, exists = engine.GetRule("temp_rule")
	if exists {
		t.Error("规则移除后仍然存在")
	}
}

// TestExample 运行基础示例
func TestExample(t *testing.T) {
	fmt.Println("运行基础示例...")
	ExampleUsage()
}

// TestAdvancedExample 运行高级示例
func TestAdvancedExample(t *testing.T) {
	fmt.Println("运行高级示例...")
	ExampleAdvanced()
}

// BenchmarkRuleExecution 性能测试
func BenchmarkRuleExecution(b *testing.B) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册规则
	engine.RegisterRule("simple_rule",
		func(data int, _ any) bool { return true },
		func(data int, _ any) (int, any, error) { return data + 1, "incremented", nil },
	)

	// 构建规则链
	_, err := engine.NewBuilder("benchmark_chain").
		AddRule("simple_rule").
		Build()

	if err != nil {
		b.Fatalf("构建规则链失败: %v", err)
	}

	// 性能测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := engine.Execute("benchmark_chain", i)
		if result.Error != nil {
			b.Fatalf("执行失败: %v", result.Error)
		}
	}
}

// BenchmarkBatchExecution 批量执行性能测试
func BenchmarkBatchExecution(b *testing.B) {
	engine := ruleengine.NewRuleEngine[int]()

	// 注册多个规则
	for i := 0; i < 5; i++ {
		ruleName := fmt.Sprintf("rule_%d", i)
		engine.RegisterRule(ruleName,
			func(data int, _ any) bool { return true },
			func(data int, _ any) (int, any, error) { return data + i, "processed", nil },
		)

		chainName := fmt.Sprintf("chain_%d", i)
		_, err := engine.NewBuilder(chainName).AddRule(ruleName).Build()
		if err != nil {
			b.Fatalf("构建规则链失败: %v", err)
		}
	}

	chainNames := []string{"chain_0", "chain_1", "chain_2", "chain_3", "chain_4"}

	// 性能测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := engine.BatchExecute(chainNames, i)
		if len(results) != 5 {
			b.Fatalf("期望5个结果，实际得到 %v", len(results))
		}
	}
}
