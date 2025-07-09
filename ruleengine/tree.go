package ruleengine

import (
	"errors"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/optional"
)

// RuleErr 规则错误类型，用于规则执行过程中的控制流管理
// RuleErr is a rule error type used for control flow management during rule execution
type RuleErr error

// 规则执行控制常量
// Rule execution control constants
var (
	// EOF 规则树结束标记，表示应该立即终止整个规则树的执行
	// EOF is a rule tree end marker indicating that the entire rule tree execution should be terminated immediately
	EOF RuleErr = errors.New("RULETREE EOF")

	// FALLTHROUGH 规则穿透标记，表示跳过当前节点但继续执行规则树
	// FALLTHROUGH is a rule fallthrough marker indicating to skip the current node but continue executing the rule tree
	FALLTHROUGH RuleErr = errors.New("RULETREE FALL THROUGH")
)

// RuleNodeMeta 规则节点元数据，如果需要做类型转换，可以使用SourceType
// RuleNodeMeta represents rule node metadata. SourceType can be used if type conversion is needed
type RuleNodeMeta[T any] struct {
	// SourceTypeData 源类型数据，存储原始的输入数据
	// SourceTypeData stores the original input data of source type
	SourceTypeData T

	// ConvertData 转换数据，存储经过类型转换或处理后的数据
	// ConvertData stores data after type conversion or processing
	ConvertData any

	// Valid 验证函数，用于判断是否应该执行当前规则节点
	// Valid is a validation function to determine whether the current rule node should be executed
	Valid func(T, any) bool

	// Exec 执行函数，实际的业务逻辑处理函数
	// Exec is the execution function that contains the actual business logic processing
	Exec func(T, any) (T, any, error)
}

// RuleLayer 规则层，包含同一层级的规则节点切片
// RuleLayer represents a rule layer containing a slice of nodes at the same level
type RuleLayer[T any] []*RuleNode[T]

// Len 计算规则层的节点总数（使用BFS遍历）
// 该方法会递归遍历所有子层，计算整个规则树的总节点数
// Len calculates the total number of nodes in the rule layer (using BFS traversal)
// This method recursively traverses all sub-layers to calculate the total number of nodes in the entire rule tree
func (r *RuleLayer[T]) Len() int {
	// 使用BFS算法计算规则树的总大小
	// Use BFS algorithm to calculate the total size of the rule tree
	queue := make([]*RuleNode[T], 0)
	visited := make(map[*RuleNode[T]]bool)
	count := 0

	// 将当前层的所有节点添加到初始队列中
	// Add all nodes of this layer to the initial queue
	for _, node := range *r {
		queue = append(queue, node)
		visited[node] = true
	}

	// 执行广度优先搜索遍历
	// Perform breadth-first search traversal
	for len(queue) > 0 {
		// 从队列头部取出一个节点
		// Dequeue a node from the front of the queue
		node := queue[0]
		queue = queue[1:]
		count++

		// 将该节点的下一层所有未访问节点加入队列
		// Add all unvisited nodes from the next layer to the queue
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
// 每个规则节点代表规则树中的一个执行单元
// RuleNode represents a rule node containing rule metadata and next layer rules
// Each rule node represents an execution unit in the rule tree
type RuleNode[T any] struct {
	RuleNodeMeta[T]
	// NxtLayer 下一层规则节点，形成树状结构
	// NxtLayer contains the next layer rule nodes, forming a tree structure
	NxtLayer RuleLayer[T]
}

// ToQueue 将当前节点加入队列（如果验证通过）
// 只有当验证函数返回true时，节点才会被加入执行队列
// ToQueue adds current node to queue (if validation passes)
// The node is added to the execution queue only when the validation function returns true
func (r *RuleNode[T]) ToQueue(queue chan *RuleNode[T]) bool {
	// 验证节点是否满足执行条件
	// Check if the node meets the execution criteria
	if r.Valid(r.SourceTypeData, r.ConvertData) {
		// 将节点发送到执行队列
		// Send the node to the execution queue
		queue <- r
	}
	return true
}

// LayerToQueue 将下一层的所有节点加入队列
// 遍历当前节点的所有子节点，并尝试将它们加入执行队列
// LayerToQueue adds all nodes in the next layer to the queue
// Iterate through all child nodes of the current node and attempt to add them to the execution queue
func (r *RuleNode[T]) LayerToQueue(queue chan *RuleNode[T]) bool {
	// 遍历下一层的所有节点
	// Iterate through all nodes in the next layer
	for _, item := range r.NxtLayer {
		// 尝试将每个节点加入队列，如果失败则返回false
		// Try to add each node to the queue, return false if failed
		if !item.ToQueue(queue) {
			return false
		}
	}
	return true
}

// Else 创建相反条件的规则节点
// 基于当前节点的验证条件创建一个条件相反的新节点
// Else creates a rule node with opposite condition
// Creates a new node with opposite condition based on the current node's validation condition
func (r *RuleNode[T]) Else(exec func(T, any) (T, any, error)) *RuleNode[T] {
	// 创建一个验证条件相反的新节点
	// Create a new node with opposite validation condition
	return NewRuleNode[T](func(t T, a any) bool {
		// 返回当前节点验证条件的相反结果
		// Return the opposite result of the current node's validation condition
		return !r.Valid(t, a)
	}, exec)
}

// LayerData 为下一层节点设置数据
// 将处理后的数据传递给下一层的所有节点
// LayerData sets data for nodes in the next layer
// Pass processed data to all nodes in the next layer
func (r *RuleNode[T]) LayerData(data T, cvntData any) bool {
	// 遍历下一层的所有节点
	// Iterate through all nodes in the next layer
	for _, item := range r.NxtLayer {
		// 检查节点是否为空
		// Check if the node is nil
		if nil == item {
			return true
		}
		// 设置节点的源数据和转换数据
		// Set the node's source data and converted data
		item.SourceTypeData = data
		item.ConvertData = cvntData
	}
	return true
}

// RuleTree 规则树，用于管理和执行规则链
// 规则树是规则执行的核心结构，使用队列进行规则的异步执行
// RuleTree manages and executes rule chains
// RuleTree is the core structure for rule execution, using queues for asynchronous rule execution
type RuleTree[T any] struct {
	// Root 根节点，规则树执行的起始点
	// Root is the root node, the starting point of rule tree execution
	Root *RuleNode[T]

	// waitQueue 等待执行的节点队列，用于异步处理规则节点
	// waitQueue is a queue of nodes waiting to be executed, used for asynchronous processing of rule nodes
	waitQueue chan *RuleNode[T]
}

// NewRuleNode 创建新的规则节点
// 构造函数用于创建包含验证逻辑和执行逻辑的规则节点
// NewRuleNode creates a new rule node
// Constructor function to create a rule node containing validation logic and execution logic
func NewRuleNode[T any](valid func(T, any) bool, exec func(T, any) (T, any, error), nxtlayer ...*RuleNode[T]) *RuleNode[T] {
	return &RuleNode[T]{
		RuleNodeMeta: RuleNodeMeta[T]{
			Valid: valid, // 验证函数 / Validation function
			Exec:  exec,  // 执行函数 / Execution function
		},
		// 如果没有提供下一层节点，则设置为nil，否则使用提供的节点
		// If no next layer nodes are provided, set to nil, otherwise use the provided nodes
		NxtLayer: optional.IsTrue(cutil.IsBlank[*[]*RuleNode[T]](&nxtlayer), nil, nxtlayer),
	}
}

// NewRuleTree 创建新的规则树
// 基于提供的根节点创建规则树，并初始化执行队列
// NewRuleTree creates a new rule tree
// Creates a rule tree based on the provided root node and initializes the execution queue
func NewRuleTree[T any](root *RuleNode[T]) *RuleTree[T] {
	return &RuleTree[T]{
		// 预留足够的队列空间，避免阻塞（根节点下一层数量 + 16个额外空间）
		// Reserve sufficient queue space to avoid blocking (next layer count + 16 extra spaces)
		waitQueue: make(chan *RuleNode[T], len(root.NxtLayer)+0x11),
		Root:      root,
	}
}

// Run 执行规则树，处理数据流转
// 这是规则树的主要执行方法，使用队列和select语句实现异步处理
// Run executes the rule tree and processes data flow
// This is the main execution method of the rule tree, using queues and select statements for asynchronous processing
func (r *RuleTree[T]) Run(data T) (T, any, error) {
	// 初始化根节点的数据
	// Initialize the root node's data
	r.Root.SourceTypeData, r.Root.ConvertData = data, data

	// 将根节点加入执行队列
	// Add the root node to the execution queue
	r.Root.ToQueue(r.waitQueue)

	// 定义执行过程中的状态变量
	// Define state variables during execution
	var orginData T // 原始数据 / Original data
	var cvtDara any // 转换数据 / Converted data
	var err error   // 错误信息 / Error information

	// 主执行循环
	// Main execution loop
	for {
		select {
		// 从队列中取出待执行的节点
		// Dequeue a node to be executed from the queue
		case node := <-r.waitQueue:
			// 执行节点的业务逻辑
			// Execute the node's business logic
			orginData, cvtDara, err = node.Exec(node.SourceTypeData, r.Root.ConvertData)

			// 处理执行过程中的错误
			// Handle errors during execution
			if err != nil {
				// 如果是EOF错误，表示规则树执行结束
				// If it's an EOF error, it indicates the rule tree execution is finished
				if errors.Is(err, EOF) {
					close(r.waitQueue) // 关闭队列 / Close the queue
					err = nil          // 清除错误状态 / Clear the error state
					// 重新创建队列以供下次使用
					// Recreate the queue for next use
					r.waitQueue = make(chan *RuleNode[T], r.Root.NxtLayer.Len()+0x11)
					return orginData, cvtDara, err
				}

				// 如果是FALLTHROUGH错误，表示跳过当前节点但继续执行
				// If it's a FALLTHROUGH error, it indicates to skip the current node but continue execution
				if errors.Is(err, FALLTHROUGH) {
					err = nil // 清除错误状态 / Clear the error state
					// 更新根节点的数据状态
					// Update the root node's data state
					r.Root.SourceTypeData = orginData
					r.Root.ConvertData = cvtDara
				}
				break // 跳出当前处理，继续下一轮循环 / Break current processing, continue to next loop
			}

			// 执行成功后，进行状态转移
			// After successful execution, perform state transition
			// 这里直接进行状态转移，把这一层的状态移交给下一层，从而避免了不同逻辑之间的影响，也不用管使用深度遍历还是广度遍历
			// State transition is performed directly here, transferring the state of this layer to the next layer,
			// thus avoiding interference between different logic and eliminating the need to consider depth-first or breadth-first traversal
			node.LayerData(orginData, cvtDara)

			// 将下一层的节点加入执行队列
			// Add nodes from the next layer to the execution queue
			node.LayerToQueue(r.waitQueue)

		// 队列为空时，表示所有节点都已处理完毕
		// When the queue is empty, it indicates all nodes have been processed
		default:
			close(r.waitQueue) // 关闭队列 / Close the queue
			// 重新创建队列以供下次使用
			// Recreate the queue for next use
			r.waitQueue = make(chan *RuleNode[T], r.Root.NxtLayer.Len()+0x11)
			return orginData, cvtDara, err
		}
	}
}
