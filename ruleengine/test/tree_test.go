package test

import (
	"errors"
	"github.com/karosown/katool-go/ruleengine"
	"testing"
)

// 测试用的简单数据类型
type TestData struct {
	Value int
}

func TestRuleNodeCreation(t *testing.T) {
	// 创建测试节点
	valid := func(data TestData, _ any) bool {
		return data.Value > 0
	}
	exec := func(data TestData, _ any) (TestData, any, error) {
		return TestData{Value: data.Value * 2}, data.Value, nil
	}

	node := ruleengine.NewRuleNode[TestData](valid, exec, nil)

	if node == nil {
		t.Fatal("Failed to create RuleNode")
	}

	if node.Valid == nil {
		t.Error("Valid function is nil")
	}

	if node.Exec == nil {
		t.Error("Exec function is nil")
	}
}

func TestRuleTreeExecution(t *testing.T) {
	// 创建叶子节点
	leafNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value > 5 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 10}, data.Value, nil
		},
		nil,
	)

	// 创建中间节点
	midNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value > 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * 2}, data.Value, nil
		},
		[]*ruleengine.RuleNode[TestData]{leafNode}...,
	)

	// 创建根节点
	rootNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return data, data, nil
		},
		[]*ruleengine.RuleNode[TestData]{midNode}...,
	)

	// 创建规则树
	tree := ruleengine.NewRuleTree[TestData](rootNode)

	// 测试用例1: 正常执行
	testData := TestData{Value: 3}
	result, _, err := tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 期望结果: 3 -> 6 (中间节点) -> 不满足叶子节点条件
	if result.Value != 16 {
		t.Errorf("Expected result value 6, got %d", result.Value)
	}

	// 测试用例2: 满足所有条件
	testData = TestData{Value: 6}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 期望结果: 6 -> 12 (中间节点) -> 22 (叶子节点)
	if result.Value != 22 {
		t.Errorf("Expected result value 22, got %d", result.Value)
	}
}

func TestRuleTreeWithError(t *testing.T) {
	// 创建会返回错误的节点
	errorNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{}, nil, errors.New("test error")
		},
		nil,
	)

	// 创建根节点
	rootNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return data, data, nil
		},
		[]*ruleengine.RuleNode[TestData]{errorNode}...,
	)

	// 创建规则树
	tree := ruleengine.NewRuleTree[TestData](rootNode)

	// 测试错误处理
	testData := TestData{Value: 1}
	_, _, err := tree.Run(testData)

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestRuleLayerLen(t *testing.T) {
	// 创建测试节点
	node1 := ruleengine.NewRuleNode[TestData](nil, nil, nil)
	node2 := ruleengine.NewRuleNode[TestData](nil, nil, nil)
	node3 := ruleengine.NewRuleNode[TestData](nil, nil, nil)

	// 创建复杂的层次结构
	node2.NxtLayer = []*ruleengine.RuleNode[TestData]{node3}
	layer := ruleengine.RuleLayer[TestData]{node1, node2}

	count := layer.Len()
	expectedCount := 3 // node1, node2, node3

	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}
}

func TestQueueOperations(t *testing.T) {
	// 创建测试节点
	node1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return data, data, nil
		},
		nil,
	)

	node2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return false }, // 不会被加入队列
		func(data TestData, _ any) (TestData, any, error) {
			return data, data, nil
		},
		nil,
	)

	// 创建队列并测试
	queue := make(chan *ruleengine.RuleNode[TestData], 5)
	node1.SourceTypeData = TestData{Value: 1}
	node2.SourceTypeData = TestData{Value: 2}

	// 测试ToQueue
	success := node1.ToQueue(queue)
	if !success {
		t.Error("node1.ToQueue should return true")
	}

	success = node2.ToQueue(queue)
	if !success {
		t.Error("node2.ToQueue should return true even when valid returns false")
	}

	// 检查队列长度
	count := len(queue)
	if count != 1 {
		t.Errorf("Expected queue length 1, got %d", count)
	}
}
func TestMultiBranchRuleTree(t *testing.T) {
	// Create branch nodes with different conditions
	branch1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value < 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * -1}, "branch1", nil
		},
		nil,
	)

	branch2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value > 0 && data.Value < 10 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 5}, "branch2", nil
		},
		nil,
	)

	branch3 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value >= 10 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * 2}, "branch3", nil
		},
		nil,
	)

	// Create root node with multiple branches
	rootNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return data, "root", nil
		},
		[]*ruleengine.RuleNode[TestData]{branch1, branch2, branch3}...,
	)

	// Create the rule tree
	tree := ruleengine.NewRuleTree[TestData](rootNode)

	// Test case 1: Negative value (should go through branch1)
	testData := TestData{Value: -5}
	result, _, err := tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 5 {
		t.Errorf("Expected result value 5, got %d", result.Value)
	}

	// Test case 2: Small positive value (should go through branch2)
	testData = TestData{Value: 7}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 12 {
		t.Errorf("Expected result value 12, got %d", result.Value)
	}

	// Test case 3: Large positive value (should go through branch3)
	testData = TestData{Value: 15}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 30 {
		t.Errorf("Expected result value 30, got %d", result.Value)
	}

	// Test case 4: Value that matches no branches (should only process root)
	testData = TestData{Value: 0}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 0 {
		t.Errorf("Expected result value 0, got %d", result.Value)
	}
}
func TestMultiBranchMultiLevelRuleTree(t *testing.T) {
	// Create the deepest level nodes (level 3) +100
	b1_1_1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value%2 == 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 100}, "b1_1_1", nil
		},
		nil,
	)
	// +200
	b1_1_2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value%2 != 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 200}, "b1_1_2", nil
		},
		nil,
	)
	// *3
	b3_1_1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * 3}, "b3_1_1", nil
		},
		nil,
	)

	// Create level 2 nodes *2
	b1_1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value < 50 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * 2}, "b1_1", nil
		},
		[]*ruleengine.RuleNode[TestData]{b1_1_1, b1_1_2}...,
	)
	// -20
	b1_2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value >= 50 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value - 20}, "b1_2", nil
		},
		nil,
	)
	// +10
	b2_1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 10}, "b2_1", nil
		},
		nil,
	)
	// +5
	b3_1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value%5 == 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 5}, "b3_1", nil
		},
		[]*ruleengine.RuleNode[TestData]{b3_1_1}...,
	)
	// -5
	b3_2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value%5 != 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value - 5}, "b3_2", nil
		},
		nil,
	)

	// Create level 1 branch nodes
	branch1 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value > 0 && data.Value < 100 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value + 1}, "branch1", nil
		},
		[]*ruleengine.RuleNode[TestData]{b1_1, b1_2}...,
	)

	branch2 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value < 0 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value * -1}, "branch2", nil
		},
		[]*ruleengine.RuleNode[TestData]{b2_1}...,
	)

	branch3 := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return data.Value >= 100 },
		func(data TestData, _ any) (TestData, any, error) {
			return TestData{Value: data.Value / 2}, "branch3", nil
		},
		[]*ruleengine.RuleNode[TestData]{b3_1, b3_2}...,
	)

	// Create root node
	rootNode := ruleengine.NewRuleNode[TestData](func(data TestData, _ any) bool { return true },
		func(data TestData, _ any) (TestData, any, error) {
			return data, "root", nil
		},
		[]*ruleengine.RuleNode[TestData]{branch1, branch2, branch3}...,
	)

	// Create the rule tree
	tree := ruleengine.NewRuleTree[TestData](rootNode)

	// Test cases to follow different paths through the tree

	// Test case 1: Path root -> branch1 -> b1_1 -> b1_1_1
	// 25 -> 26 -> 52 -> 152
	testData := TestData{Value: 25}
	result, _, err := tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 1: Expected no error, got %v", err)
	}

	if result.Value != 152 {
		t.Errorf("Test case 1: Expected result value 152, got %d", result.Value)
	}

	// Test case 2: Path root -> branch1 -> b1_1 -> b1_1_2
	// 26 -> 27 -> 54 -> 154
	testData = TestData{Value: 26}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 2: Expected no error, got %v", err)
	}

	if result.Value != 154 {
		t.Errorf("Test case 2: Expected result value 254, got %d", result.Value)
	}

	// Test case 3: Path root -> branch1 -> b1_2
	// 60 -> 61 -> 41
	testData = TestData{Value: 60}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 3: Expected no error, got %v", err)
	}

	if result.Value != 41 {
		t.Errorf("Test case 3: Expected result value 41, got %d", result.Value)
	}

	// Test case 4: Path root -> branch2 -> b2_1
	// -10 -> 10 -> 20
	testData = TestData{Value: -10}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 4: Expected no error, got %v", err)
	}

	if result.Value != 20 {
		t.Errorf("Test case 4: Expected result value 20, got %d", result.Value)
	}

	// Test case 5: Path root -> branch3 -> b3_1 -> b3_1_1
	// 100 -> 50 -> 55 -> 165
	testData = TestData{Value: 100}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 5: Expected no error, got %v", err)
	}

	if result.Value != 165 {
		t.Errorf("Test case 5: Expected result value 165, got %d", result.Value)
	}

	// Test case 6: Path root -> branch3 -> b3_2
	// 102 -> 51 -> 46
	testData = TestData{Value: 102}
	result, _, err = tree.Run(testData)

	if err != nil {
		t.Errorf("Test case 6: Expected no error, got %v", err)
	}

	if result.Value != 46 {
		t.Errorf("Test case 6: Expected result value 46, got %d", result.Value)
	}
}
