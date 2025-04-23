package ruleengine

import (
	"errors"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/optional"
)

type RuleErr error

var (
	EOF RuleErr = errors.New("RULETREE EOF")
)

// RuleNodeMeta 如果需要做类型转换，可以使用SourceType
type RuleNodeMeta[T any] struct {
	SourceTypeData T
	ConvertData    any
	Valid          func(T, any) bool
	Exec           func(T, any) (T, any, error)
}
type RuleLayer[T any] []*RuleNode[T]

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

type RuleNode[T any] struct {
	RuleNodeMeta[T]
	NxtLayer RuleLayer[T]
}

func (r *RuleNode[T]) ToQueue(queue chan *RuleNode[T]) bool {
	if r.Valid(r.SourceTypeData, r.ConvertData) {
		queue <- r
	}
	return true
}
func (r *RuleNode[T]) LayerToQueue(queue chan *RuleNode[T]) bool {
	for _, item := range r.NxtLayer {
		if !item.ToQueue(queue) {
			return false
		}
	}
	return true
}
func (r *RuleNode[T]) Else(exec func(T, any) (T, any, error)) *RuleNode[T] {
	return NewRuleNode[T](func(t T, a any) bool {
		return !r.Valid(t, a)
	}, exec)
}
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

type RuleTree[T any] struct {
	Root      *RuleNode[T]
	waitQueue chan *RuleNode[T]
}

func NewRuleNode[T any](valid func(T, any) bool, exec func(T, any) (T, any, error), nxtlayer ...*RuleNode[T]) *RuleNode[T] {
	return &RuleNode[T]{
		RuleNodeMeta: RuleNodeMeta[T]{
			Valid: valid,
			Exec:  exec,
		},
		NxtLayer: optional.IsTrue(cutil.IsBlank[*[]*RuleNode[T]](&nxtlayer), nil, nxtlayer),
	}
}
func NewRuleTree[T any](root *RuleNode[T]) *RuleTree[T] {
	return &RuleTree[T]{
		// 预留10个位置
		waitQueue: make(chan *RuleNode[T], len(root.NxtLayer)+0x11),
		Root:      root,
	}
}

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
					r.waitQueue = make(chan *RuleNode[T], r.Root.NxtLayer.Len()+0x11)
					return orginData, cvtDara, nil
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
