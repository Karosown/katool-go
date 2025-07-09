package ruleengine

import (
	"errors"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/optional"
)

// RuleErr 规则错误类型
// RuleErr is a rule error type
type RuleErr error

// 规则执行控制常量
// Rule execution control constants
var (
	EOF         RuleErr = errors.New("RULETREE EOF")          // 规则树结束标记 / Rule tree end marker
	FALLTHROUGH RuleErr = errors.New("RULETREE FALL THROUGH") // 规则穿透标记 / Rule fallthrough marker
)

// RuleNodeMeta 如果需要做类型转换，可以使用SourceType
// RuleNodeMeta can use SourceType if type conversion is needed
type RuleNodeMeta[T any] struct {
	SourceTypeData T                            // 源类型数据 / Source type data
	ConvertData    any                          // 转换数据 / Converted data
	Valid          func(T, any) bool            // 验证函数 / Validation function
	Exec           func(T, any) (T, any, error) // 执行函数 / Execution function
}

// RuleLayer 规则层，包含同一层级的规则节点
// RuleLayer represents a rule layer containing nodes at the same level
type RuleLayer[T any] []*RuleNode[T]

// Len 计算规则层的节点总数（使用BFS遍历）
// Len calculates the total number of nodes in the rule layer (using BFS traversal)
func (r *RuleLayer[T]) Len() int {
	// bfs计算大小
	queue := make([]*RuleNode[T], 0)
	visited := make(map[*RuleNode[T]]bool)
	count := 0

	// Add all nodes of this layer to the initial queue
	for _, node := range *r {
		queue = append(queue, node)
		visited[node] = true
	}

	// Perform BFS
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		count++
		for _, nxt := range node.NxtLayer {
			if !visited[nxt] {
				queue = append(queue, nxt)
				visited[nxt] = true
			}
		}
	}
	return count
}

// RuleNode 规则节点，包含规则元数据和下一层规则
// RuleNode represents a rule node containing rule metadata and next layer rules
type RuleNode[T any] struct {
	RuleNodeMeta[T]
	NxtLayer RuleLayer[T] // 下一层规则节点 / Next layer rule nodes
}

// ToQueue 将当前节点加入队列（如果验证通过）
// ToQueue adds current node to queue (if validation passes)
func (r *RuleNode[T]) ToQueue(queue chan *RuleNode[T]) bool {
	if r.Valid(r.SourceTypeData, r.ConvertData) {
		queue <- r
	}
	return true
}

// LayerToQueue 将下一层的所有节点加入队列
// LayerToQueue adds all nodes in the next layer to the queue
func (r *RuleNode[T]) LayerToQueue(queue chan *RuleNode[T]) bool {
	for _, item := range r.NxtLayer {
		if !item.ToQueue(queue) {
			return false
		}
	}
	return true
}

// Else 创建相反条件的规则节点
// Else creates a rule node with opposite condition
func (r *RuleNode[T]) Else(exec func(T, any) (T, any, error)) *RuleNode[T] {
	return NewRuleNode[T](func(t T, a any) bool {
		return !r.Valid(t, a)
	}, exec)
}

// LayerData 为下一层节点设置数据
// LayerData sets data for nodes in the next layer
func (r *RuleNode[T]) LayerData(data T, cvntData any) bool {
	for _, item := range r.NxtLayer {
		if nil == item {
			return true
		}
		item.SourceTypeData = data
		item.ConvertData = cvntData
	}
	return true
}

// RuleTree 规则树，用于管理和执行规则链
// RuleTree manages and executes rule chains
type RuleTree[T any] struct {
	Root      *RuleNode[T]      // 根节点 / Root node
	waitQueue chan *RuleNode[T] // 等待队列 / Wait queue
}

// NewRuleNode 创建新的规则节点
// NewRuleNode creates a new rule node
func NewRuleNode[T any](valid func(T, any) bool, exec func(T, any) (T, any, error), nxtlayer ...*RuleNode[T]) *RuleNode[T] {
	return &RuleNode[T]{
		RuleNodeMeta: RuleNodeMeta[T]{
			Valid: valid,
			Exec:  exec,
		},
		NxtLayer: optional.IsTrue(cutil.IsBlank[*[]*RuleNode[T]](&nxtlayer), nil, nxtlayer),
	}
}

// NewRuleTree 创建新的规则树
// NewRuleTree creates a new rule tree
func NewRuleTree[T any](root *RuleNode[T]) *RuleTree[T] {
	return &RuleTree[T]{
		// 预留10个位置
		waitQueue: make(chan *RuleNode[T], len(root.NxtLayer)+0x11),
		Root:      root,
	}
}

// Run 执行规则树，处理数据流转
// Run executes the rule tree and processes data flow
func (r *RuleTree[T]) Run(data T) (T, any, error) {
	r.Root.SourceTypeData, r.Root.ConvertData = data, data
	r.Root.ToQueue(r.waitQueue)
	var orginData T
	var cvtDara any
	var err error
	for {
		select {
		case node := <-r.waitQueue:
			orginData, cvtDara, err = node.Exec(node.SourceTypeData, r.Root.ConvertData)
			if err != nil {
				if errors.Is(err, EOF) {
					close(r.waitQueue)
					err = nil
					r.waitQueue = make(chan *RuleNode[T], r.Root.NxtLayer.Len()+0x11)
					return orginData, cvtDara, err
				}
				if errors.Is(err, FALLTHROUGH) {
					err = nil
					r.Root.SourceTypeData = orginData
					r.Root.ConvertData = cvtDara
				}
				break
			}
			// 这里直接进行状态转移，把这一层的状态移交给下一层，从而避免了不同逻辑之间的影响，也不用管使用深度遍历还是广度遍历
			node.LayerData(orginData, cvtDara)
			node.LayerToQueue(r.waitQueue)
		default:
			close(r.waitQueue)
			r.waitQueue = make(chan *RuleNode[T], r.Root.NxtLayer.Len()+0x11)
			return orginData, cvtDara, err
		}
	}
}
