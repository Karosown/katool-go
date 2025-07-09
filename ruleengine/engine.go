package ruleengine

import (
	"fmt"
	"sync"
)

// RuleEngine 规则引擎管理器
// RuleEngine is a rule engine manager
type RuleEngine[T any] struct {
	rules      map[string]*RuleNode[T] // 已注册的规则节点
	chains     map[string]*RuleTree[T] // 已构建的规则链
	mutex      sync.RWMutex            // 读写锁保护并发访问
	middleware []MiddlewareFunc[T]     // 中间件函数
}

// MiddlewareFunc 中间件函数类型
// MiddlewareFunc is the type for middleware functions
type MiddlewareFunc[T any] func(data T, next func(T) (T, any, error)) (T, any, error)

// RuleBuilder 规则构建器
// RuleBuilder is a rule builder
type RuleBuilder[T any] struct {
	engine *RuleEngine[T]
	nodes  []*RuleNode[T]
	name   string
}

// ExecuteResult 执行结果
// ExecuteResult represents the execution result
type ExecuteResult[T any] struct {
	Data   T
	Result any
	Error  error
	Chain  string
}

// NewRuleEngine 创建新的规则引擎
// NewRuleEngine creates a new rule engine
func NewRuleEngine[T any]() *RuleEngine[T] {
	return &RuleEngine[T]{
		rules:      make(map[string]*RuleNode[T]),
		chains:     make(map[string]*RuleTree[T]),
		middleware: make([]MiddlewareFunc[T], 0),
	}
}

// RegisterRule 注册单个规则节点
// RegisterRule registers a single rule node
func (e *RuleEngine[T]) RegisterRule(name string, valid func(T, any) bool, exec func(T, any) (T, any, error)) *RuleEngine[T] {
	return e.RegisterRuleNode(name, NewRuleNode(valid, exec))
}

// RegisterRuleNode 注册规则节点
// RegisterRuleNode registers a rule node
func (e *RuleEngine[T]) RegisterRuleNode(name string, node *RuleNode[T]) *RuleEngine[T] {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.rules[name] = node
	return e
}

// GetRule 获取已注册的规则
// GetRule gets a registered rule
func (e *RuleEngine[T]) GetRule(name string) (*RuleNode[T], bool) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	rule, exists := e.rules[name]
	return rule, exists
}

// NewBuilder 创建规则构建器
// NewBuilder creates a rule builder
func (e *RuleEngine[T]) NewBuilder(chainName string) *RuleBuilder[T] {
	return &RuleBuilder[T]{
		engine: e,
		nodes:  make([]*RuleNode[T], 0),
		name:   chainName,
	}
}

// AddRule 向构建器添加规则（通过名称）
// AddRule adds a rule to the builder (by name)
func (b *RuleBuilder[T]) AddRule(ruleName string) *RuleBuilder[T] {
	if rule, exists := b.engine.GetRule(ruleName); exists {
		b.nodes = append(b.nodes, rule)
	}
	return b
}

// AddCustomRule 向构建器添加自定义规则
// AddCustomRule adds a custom rule to the builder
func (b *RuleBuilder[T]) AddCustomRule(valid func(T, any) bool, exec func(T, any) (T, any, error)) *RuleBuilder[T] {
	return b.AddCustomRuleNode(NewRuleNode(valid, exec))
}

// AddCustomRuleNode 向构建器添加自定义规则节点
// AddCustomRuleNode adds a custom rule node to the builder
func (b *RuleBuilder[T]) AddCustomRuleNode(node *RuleNode[T]) *RuleBuilder[T] {
	b.nodes = append(b.nodes, node)
	return b
}

// AddConditionalChain 添加条件分支链
// AddConditionalChain adds a conditional branch chain
func (b *RuleBuilder[T]) AddConditionalChain(condition func(T, any) bool, trueChain, falseChain []*RuleNode[T]) *RuleBuilder[T] {
	conditionNode := NewRuleNode(
		condition,
		func(data T, convertData any) (T, any, error) {
			if condition(data, convertData) {
				// 执行 true 分支
				for _, node := range trueChain {
					var err error
					data, convertData, err = node.Exec(data, convertData)
					if err != nil {
						return data, convertData, err
					}
				}
			} else {
				// 执行 false 分支
				for _, node := range falseChain {
					var err error
					data, convertData, err = node.Exec(data, convertData)
					if err != nil {
						return data, convertData, err
					}
				}
			}
			return data, convertData, nil
		},
	)
	b.nodes = append(b.nodes, conditionNode)
	return b
}

// Build 构建规则链
// Build builds the rule chain
func (b *RuleBuilder[T]) Build() (*RuleTree[T], error) {
	if len(b.nodes) == 0 {
		return nil, fmt.Errorf("规则链 '%s' 为空", b.name)
	}

	// 构建树形结构
	root := b.nodes[0]
	current := root

	for i := 1; i < len(b.nodes); i++ {
		current.NxtLayer = append(current.NxtLayer, b.nodes[i])
		current = b.nodes[i]
	}

	tree := NewRuleTree(root)

	// 注册到引擎
	b.engine.mutex.Lock()
	b.engine.chains[b.name] = tree
	b.engine.mutex.Unlock()

	return tree, nil
}

// Execute 执行指定的规则链
// Execute executes the specified rule chain
func (e *RuleEngine[T]) Execute(chainName string, data T) *ExecuteResult[T] {
	e.mutex.RLock()
	chain, exists := e.chains[chainName]
	e.mutex.RUnlock()

	if !exists {
		return &ExecuteResult[T]{
			Data:  data,
			Error: fmt.Errorf("规则链 '%s' 不存在", chainName),
			Chain: chainName,
		}
	}

	// 应用中间件
	finalData, result, err := e.applyMiddleware(data, func(d T) (T, any, error) {
		return chain.Run(d)
	})

	return &ExecuteResult[T]{
		Data:   finalData,
		Result: result,
		Error:  err,
		Chain:  chainName,
	}
}

// ExecuteAll 执行所有规则链
// ExecuteAll executes all rule chains
func (e *RuleEngine[T]) ExecuteAll(data T) map[string]*ExecuteResult[T] {
	e.mutex.RLock()
	chains := make(map[string]*RuleTree[T])
	for name, chain := range e.chains {
		chains[name] = chain
	}
	e.mutex.RUnlock()

	results := make(map[string]*ExecuteResult[T])

	for name := range chains {
		results[name] = e.Execute(name, data)
	}

	return results
}

// BatchExecute 批量执行规则链（并发）
// BatchExecute executes rule chains in batch (concurrent)
func (e *RuleEngine[T]) BatchExecute(chainNames []string, data T) map[string]*ExecuteResult[T] {
	results := make(map[string]*ExecuteResult[T])
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, chainName := range chainNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			result := e.Execute(name, data)

			mutex.Lock()
			results[name] = result
			mutex.Unlock()
		}(chainName)
	}

	wg.Wait()
	return results
}

// AddMiddleware 添加中间件
// AddMiddleware adds middleware
func (e *RuleEngine[T]) AddMiddleware(middleware MiddlewareFunc[T]) *RuleEngine[T] {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.middleware = append(e.middleware, middleware)
	return e
}

// applyMiddleware 应用中间件
// applyMiddleware applies middleware
func (e *RuleEngine[T]) applyMiddleware(data T, final func(T) (T, any, error)) (T, any, error) {
	if len(e.middleware) == 0 {
		return final(data)
	}

	// 递归应用中间件
	var applyNext func(int, T) (T, any, error)
	applyNext = func(index int, d T) (T, any, error) {
		if index >= len(e.middleware) {
			return final(d)
		}

		return e.middleware[index](d, func(nextData T) (T, any, error) {
			return applyNext(index+1, nextData)
		})
	}

	return applyNext(0, data)
}

// ListRules 列出所有已注册的规则
// ListRules lists all registered rules
func (e *RuleEngine[T]) ListRules() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	rules := make([]string, 0, len(e.rules))
	for name := range e.rules {
		rules = append(rules, name)
	}
	return rules
}

// ListChains 列出所有已构建的规则链
// ListChains lists all built rule chains
func (e *RuleEngine[T]) ListChains() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	chains := make([]string, 0, len(e.chains))
	for name := range e.chains {
		chains = append(chains, name)
	}
	return chains
}

// RemoveRule 移除规则
// RemoveRule removes a rule
func (e *RuleEngine[T]) RemoveRule(name string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, exists := e.rules[name]; exists {
		delete(e.rules, name)
		return true
	}
	return false
}

// RemoveChain 移除规则链
// RemoveChain removes a rule chain
func (e *RuleEngine[T]) RemoveChain(name string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, exists := e.chains[name]; exists {
		delete(e.chains, name)
		return true
	}
	return false
}

// Clear 清空所有规则和规则链
// Clear clears all rules and rule chains
func (e *RuleEngine[T]) Clear() *RuleEngine[T] {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.rules = make(map[string]*RuleNode[T])
	e.chains = make(map[string]*RuleTree[T])
	e.middleware = make([]MiddlewareFunc[T], 0)

	return e
}

// Stats 获取引擎统计信息
// Stats gets engine statistics
func (e *RuleEngine[T]) Stats() map[string]interface{} {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return map[string]interface{}{
		"rules_count":      len(e.rules),
		"chains_count":     len(e.chains),
		"middleware_count": len(e.middleware),
	}
}
