package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/karosown/katool-go/ai/types"
	"github.com/karosown/katool-go/xlog"
)

// MultiMCPAdapter 多MCP适配器，用于合并多个MCP服务器的工具
type MultiMCPAdapter struct {
	// 多个MCP适配器
	adapters []*MCPAdapter

	// 工具缓存（合并后的）
	toolsCache []types.Tool
	toolsMap   map[string]*types.Tool
	toolSource map[string]int // 工具名称 -> 适配器索引

	// 日志记录器
	logger xlog.Logger

	// 互斥锁
	mu  sync.RWMutex
	ctx context.Context
}

func (m *MultiMCPAdapter) Context() context.Context {
	return m.ctx
}
func (m *MultiMCPAdapter) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// NewMultiMCPAdapter 创建多MCP适配器
func NewMultiMCPAdapter(ctx context.Context, logger xlog.Logger) *MultiMCPAdapter {
	return &MultiMCPAdapter{
		adapters:   make([]*MCPAdapter, 0),
		toolsCache: make([]types.Tool, 0),
		toolsMap:   make(map[string]*types.Tool),
		toolSource: make(map[string]int),
		logger:     logger,
		ctx:        ctx,
	}
}

// AddAdapter 添加MCP适配器
func (m *MultiMCPAdapter) AddAdapter(adapter *MCPAdapter) error {
	if adapter == nil {
		return fmt.Errorf("MCP adapter cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.adapters = append(m.adapters, adapter)
	adapterIndex := len(m.adapters) - 1

	// 合并工具
	tools := adapter.GetTools()
	for _, tool := range tools {
		toolName := tool.Function.Name

		// 检查工具名称冲突
		if existingIndex, exists := m.toolSource[toolName]; exists {
			m.logger.Warnf("Tool name conflict: %s already exists from adapter %d, skipping from adapter %d", toolName, existingIndex, adapterIndex)
			continue
		}

		// 添加工具
		m.toolsCache = append(m.toolsCache, tool)
		m.toolsMap[toolName] = &tool
		m.toolSource[toolName] = adapterIndex
	}

	m.logger.Infof("Added MCP adapter %d, total tools: %d", adapterIndex, len(m.toolsCache))
	return nil
}

// RemoveAdapter 移除MCP适配器
func (m *MultiMCPAdapter) RemoveAdapter(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if index < 0 || index >= len(m.adapters) {
		return fmt.Errorf("adapter index out of range: %d", index)
	}

	// 移除该适配器的所有工具
	toolsToRemove := make(map[string]bool)
	for toolName, adapterIndex := range m.toolSource {
		if adapterIndex == index {
			toolsToRemove[toolName] = true
		}
	}

	// 重建工具列表
	newToolsCache := make([]types.Tool, 0)
	newToolsMap := make(map[string]*types.Tool)
	newToolSource := make(map[string]int)

	for i, tool := range m.toolsCache {
		toolName := tool.Function.Name
		if !toolsToRemove[toolName] {
			newToolsCache = append(newToolsCache, tool)
			newToolsMap[toolName] = &m.toolsCache[i]
			// 更新适配器索引（移除的适配器之后的索引需要减1）
			oldIndex := m.toolSource[toolName]
			if oldIndex > index {
				newToolSource[toolName] = oldIndex - 1
			} else {
				newToolSource[toolName] = oldIndex
			}
		}
	}

	// 移除适配器
	m.adapters = append(m.adapters[:index], m.adapters[index+1:]...)
	m.toolsCache = newToolsCache
	m.toolsMap = newToolsMap
	m.toolSource = newToolSource

	m.logger.Infof("Removed MCP adapter %d, remaining tools: %d", index, len(m.toolsCache))
	return nil
}

// GetAdapters 获取所有适配器
func (m *MultiMCPAdapter) GetAdapters() []*MCPAdapter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	adapters := make([]*MCPAdapter, len(m.adapters))
	copy(adapters, m.adapters)
	return adapters
}

// GetAdapterCount 获取适配器数量
func (m *MultiMCPAdapter) GetAdapterCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.adapters)
}

// GetTools 获取所有工具（合并后的）
func (m *MultiMCPAdapter) GetTools() []types.Tool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tools := make([]types.Tool, len(m.toolsCache))
	copy(tools, m.toolsCache)
	return tools
}

// HasTool 检查工具是否存在
func (m *MultiMCPAdapter) HasTool(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.toolsMap[name]
	return exists
}

// GetToolSource 获取工具来源的适配器索引
func (m *MultiMCPAdapter) GetToolSource(toolName string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	index, exists := m.toolSource[toolName]
	return index, exists
}

// CallTool 调用工具
func (m *MultiMCPAdapter) CallTool(ctx context.Context, name string, arguments string) (interface{}, error) {
	m.mu.RLock()

	// 查找工具来源
	adapterIndex, exists := m.toolSource[name]
	if !exists {
		m.mu.RUnlock()
		return nil, fmt.Errorf("tool %s not found", name)
	}

	// 获取对应的适配器
	if adapterIndex < 0 || adapterIndex >= len(m.adapters) {
		m.mu.RUnlock()
		return nil, fmt.Errorf("invalid adapter index: %d", adapterIndex)
	}

	adapter := m.adapters[adapterIndex]
	m.mu.RUnlock()

	// 调用工具
	return adapter.CallTool(ctx, name, arguments)
}

// RefreshTools 刷新所有适配器的工具列表
func (m *MultiMCPAdapter) RefreshTools(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清空缓存
	m.toolsCache = make([]types.Tool, 0)
	m.toolsMap = make(map[string]*types.Tool)
	m.toolSource = make(map[string]int)

	// 重新合并所有适配器的工具
	for adapterIndex, adapter := range m.adapters {
		if err := adapter.RefreshTools(ctx); err != nil {
			m.logger.Warnf("Failed to refresh tools for adapter %d: %v", adapterIndex, err)
			continue
		}

		tools := adapter.GetTools()
		for _, tool := range tools {
			toolName := tool.Function.Name

			// 检查工具名称冲突
			if _, exists := m.toolSource[toolName]; exists {
				m.logger.Warnf("Tool name conflict: %s already exists, skipping from adapter %d", toolName, adapterIndex)
				continue
			}

			// 添加工具
			m.toolsCache = append(m.toolsCache, tool)
			m.toolsMap[toolName] = &tool
			m.toolSource[toolName] = adapterIndex
		}
	}

	m.logger.Infof("Refreshed tools from %d adapters, total tools: %d", len(m.adapters), len(m.toolsCache))
	return nil
}

// GetToolCount 获取工具数量
func (m *MultiMCPAdapter) GetToolCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.toolsCache)
}

// GetToolCountByAdapter 获取每个适配器的工具数量
func (m *MultiMCPAdapter) GetToolCountByAdapter() map[int]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	counts := make(map[int]int)
	for _, adapterIndex := range m.toolSource {
		counts[adapterIndex]++
	}
	return counts
}
