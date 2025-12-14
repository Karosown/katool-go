# 高德地图 MCP 测试程序

这个程序演示如何使用 `katool-go/ai/agent` 模块连接和测试高德地图的 MCP (Model Context Protocol) 服务。

## 功能特性

程序提供了以下高德地图相关功能：

1. **地理编码** (`geocode`) - 将地址转换为经纬度坐标
2. **逆地理编码** (`reverse_geocode`) - 将经纬度坐标转换为地址信息
3. **路径规划** (`route_planning`) - 规划两点之间的路径（支持驾车、步行、骑行）
4. **地点搜索** (`place_search`) - 根据关键词搜索地点信息
5. **天气查询** (`weather_query`) - 查询指定位置的天气信息

## 使用方式

### 方式1: 使用 SimpleMCPClient（模拟服务，用于测试）

这是默认方式，使用内置的 `SimpleMCPClient` 模拟高德地图服务，无需外部依赖。

```bash
cd ai/agent/examples
go run amap_mcp_test.go
```

### 方式2: 使用 mark3labs/mcp-go 连接真实服务器

如果需要连接真实的高德地图 MCP 服务器：

1. 安装依赖：
```bash
go get github.com/mark3labs/mcp-go
```

2. 设置环境变量（高德地图MCP服务器地址）：
```bash
export AMAP_MCP_ENDPOINT=https://mcp.amap.com/sse
```

3. 取消注释 `testWithMark3LabsMCP` 函数中的代码并运行

### 方式3: 使用官方 SDK 连接真实服务器

如果需要使用官方的 MCP SDK：

1. 安装依赖：
```bash
go get github.com/modelcontextprotocol/go-sdk
```

2. 设置环境变量：
```bash
export AMAP_MCP_ENDPOINT=https://mcp.amap.com/sse
```

3. 取消注释 `testWithOfficialSDK` 函数中的代码并运行

## 测试内容

程序会执行以下测试：

1. **工具列表展示** - 显示所有可用的高德地图工具
2. **直接工具调用测试**：
   - 地理编码：查询"北京市天安门广场"的坐标
   - 路径规划：规划从天安门到故宫的驾车路线
   - 地点搜索：搜索"天安门"相关地点
3. **Agent 自动执行测试**：
   - 使用 AI Agent 自动调用工具完成任务
   - 测试多轮对话和工具调用

## 示例输出

```
=== 高德地图 MCP 测试程序 ===

--- 方式1: 使用 SimpleMCPClient 模拟高德地图服务 ---

可用工具数量: 5
  - geocode: 将地址转换为经纬度坐标
  - reverse_geocode: 将经纬度坐标转换为地址信息
  - route_planning: 规划两点之间的路径
  - place_search: 根据关键词搜索地点信息
  - weather_query: 查询指定位置的天气信息

--- 测试工具调用 ---

1. 测试地理编码工具...
地理编码结果: map[formatted_address:北京市东城区天安门广场 location:116.397428,39.90923]

2. 测试路径规划工具...
路径规划结果: map[distance:15.2公里 duration:25分钟 ...]

3. 测试地点搜索工具...
地点搜索结果: map[results:[...]]

--- 测试Agent执行任务 ---

任务1: 查询'北京市天安门广场'的坐标
响应: 根据查询结果，北京市天安门广场的坐标是 116.397428,39.90923...

任务2: 规划从天安门到故宫的驾车路线
响应: 已为您规划好从天安门到故宫的驾车路线，距离15.2公里，预计25分钟...

任务3: 搜索附近的餐厅
响应: 已找到天安门附近的餐厅...
```

## 注意事项

1. **模拟服务**：默认使用 `SimpleMCPClient` 模拟服务，返回的是模拟数据
2. **真实服务**：连接真实的高德地图 MCP 服务器需要：
   - 高德地图开放平台的 API Key
   - 正确配置的 MCP 服务器地址
   - 网络连接
3. **AI 模型配置**：程序中使用 `gpt-4` 作为默认模型，你可以根据实际情况修改

## 自定义配置

### 修改 AI 模型

在 `testAgentExecution` 函数中修改：

```go
agent.WithAgentConfig(&agent.AgentConfig{
    Model: "your-model-name",  // 修改这里
    MaxToolCallRounds: 5,
})
```

### 添加更多工具

在 `registerAmapTools` 函数中添加新的工具注册：

```go
client.RegisterTool(agent.MCPTool{
    Name:        "your_tool_name",
    Description: "工具描述",
    Parameters: map[string]interface{}{
        // 参数定义
    },
}, func(ctx context.Context, args string) (interface{}, error) {
    // 工具实现
    return result, nil
})
```

## 故障排除

1. **编译错误**：确保已安装所有依赖
2. **运行时错误**：检查日志输出，确认工具是否正确注册
3. **连接失败**：如果使用真实服务器，检查网络连接和服务器地址配置
