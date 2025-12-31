package ai

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/helper/jsonhp"

	"github.com/karosown/katool-go/ai/providers"
	"github.com/karosown/katool-go/xlog"
)

// mockProvider 模拟的AI提供者，用于测试
type mockProvider struct {
	name   string
	models []string
	config *Config
}

func (m *mockProvider) Chat(req *ChatRequest) (*ChatResponse, error) {
	return &ChatResponse{
		ID:      "test-id",
		Model:   req.Model,
		Choices: []Choice{{Message: Message{Role: "assistant", Content: "Hello from mock"}}},
	}, nil
}

func (m *mockProvider) ChatStream(req *ChatRequest) (<-chan *ChatResponse, error) {
	ch := make(StreamChatResponse, 1)
	ch <- &ChatResponse{
		ID:      "test-id",
		Model:   req.Model,
		Choices: []Choice{{Delta: Message{Role: "assistant", Content: "Hello"}}},
	}
	ch.Close(nil)
	return ch, nil
}

func (m *mockProvider) ChatWithTools(req *ChatRequest, tools []Tool) (*ChatResponse, error) {
	return m.Chat(req)
}

func (m *mockProvider) GetName() string {
	return m.name
}

func (m *mockProvider) GetModels() []string {
	return m.models
}

func (m *mockProvider) ValidateConfig() error {
	if m.config == nil {
		return nil
	}
	if m.name == "invalid" {
		return fmt.Errorf("invalid config")
	}
	return nil
}

// TestNewClientWithProvider 测试使用指定提供者创建客户端
func TestNewClientWithProvider(t *testing.T) {
	config := &Config{
		APIKey:     "test-key",
		BaseURL:    "https://api.test.com/v1",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Headers:    make(map[string]string),
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if client.GetProvider() != ProviderOllama {
		t.Errorf("Expected provider Ollama, got %s", client.GetProvider())
	}

	if !client.HasProvider(ProviderOllama) {
		t.Error("Client should have Ollama provider")
	}
}

// TestNewClientWithProvider_NilConfig 测试使用nil配置创建客户端
func TestNewClientWithProvider_NilConfig(t *testing.T) {
	client, err := NewClientWithProvider(ProviderOllama, nil)
	if err != nil {
		t.Fatalf("Failed to create client with nil config: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}
}

// TestNewClientWithProvider_InvalidProvider 测试使用无效的提供者类型
func TestNewClientWithProvider_InvalidProvider(t *testing.T) {
	config := &Config{
		APIKey:  "test-key",
		BaseURL: "https://api.test.com/v1",
	}

	_, err := NewClientWithProvider(ProviderType("invalid"), config)
	if err == nil {
		t.Error("Expected error for invalid provider type")
	}
}

// TestSetProvider 测试切换提供者
func TestSetProvider(t *testing.T) {
	config1 := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config1)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 尝试切换到不存在的提供者
	err = client.SetProvider(ProviderOpenAI)
	if err == nil {
		t.Error("Expected error when switching to non-existent provider")
	}

	// 添加第二个提供者（通过手动创建）
	config2 := &Config{
		APIKey:  "test-key",
		BaseURL: "https://api.openai.com/v1",
	}

	// 手动添加提供者到map（用于测试）
	provider2 := providers.NewOpenAIProvider(config2)
	client.providers[ProviderOpenAI] = provider2

	// 现在应该可以切换了
	err = client.SetProvider(ProviderOpenAI)
	if err != nil {
		t.Fatalf("Failed to switch provider: %v", err)
	}

	if client.GetProvider() != ProviderOpenAI {
		t.Errorf("Expected provider OpenAI, got %s", client.GetProvider())
	}
}

// TestGetProvider 测试获取当前提供者
func TestGetProvider(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	provider := client.GetProvider()
	if provider != ProviderOllama {
		t.Errorf("Expected provider Ollama, got %s", provider)
	}
}

// TestHasProvider 测试检查提供者是否存在
func TestHasProvider(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if !client.HasProvider(ProviderOllama) {
		t.Error("Client should have Ollama provider")
	}

	if client.HasProvider(ProviderOpenAI) {
		t.Error("Client should not have OpenAI provider")
	}
}

// TestListProviders 测试列出所有提供者
func TestListProviders(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	providers := client.ListProviders()
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(providers))
	}

	if providers[0] != ProviderOllama {
		t.Errorf("Expected provider Ollama, got %s", providers[0])
	}
}

// TestRegisterFunction 测试注册函数
func TestRegisterFunction(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 注册一个简单的函数
	err = client.RegisterFunction("add", "两数相加", func(a, b int) int {
		return a + b
	})

	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	functions := client.GetRegisteredFunctions()
	if len(functions) != 1 {
		t.Errorf("Expected 1 registered function, got %d", len(functions))
	}

	if functions[0] != "add" {
		t.Errorf("Expected function name 'add', got '%s'", functions[0])
	}
}

// TestRegisterFunction_Multiple 测试注册多个函数
func TestRegisterFunction_Multiple(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 注册多个函数
	functions := []struct {
		name        string
		description string
		fn          interface{}
	}{
		{"add", "两数相加", func(a, b int) int { return a + b }},
		{"multiply", "两数相乘", func(a, b int) int { return a * b }},
		{"greet", "问候", func(name string) string { return "Hello, " + name }},
	}

	for _, f := range functions {
		err := client.RegisterFunction(f.name, f.description, f.fn)
		if err != nil {
			t.Fatalf("Failed to register function %s: %v", f.name, err)
		}
	}

	registered := client.GetRegisteredFunctions()
	if len(registered) != len(functions) {
		t.Errorf("Expected %d registered functions, got %d", len(functions), len(registered))
	}
}

// TestChat 测试基本聊天功能（需要真实API，跳过）
func TestChat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires real API")
	}

	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &ChatRequest{
		Model: "llama2",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	_, err = client.Chat(req)
	if err != nil {
		t.Logf("Chat failed (this is expected if Ollama is not running): %v", err)
	}
}

// TestChat 测试基本聊天功能（需要真实API，跳过）
func TestChatRes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires real API")
	}

	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Ollama 的 format 参数只接受 "json" 字符串，不接受 JSON Schema 对象
	// 需要在 prompt 中描述期望的 JSON 结构
	req := &ChatRequest{
		//Model: "llama3.1",
		Model: "Qwen2",
		Messages: []Message{
			NewMessage(RoleSystem, `你是一个智能的研究人员。我给你一个关键词，请你返回给我可能找到结果的网址。注意一定要高相关性！可以是二级域名返回给我，格式如下
`+jsonhp.ToJSON([]struct {
				Link string `json:"Link" description:"Link"`
			}{{}})),
			NewMessage(RoleUser, "电子科大计算机考研通知"),
		},
		Format: optional.Must(FormatArrayOf([]struct {
			A string `json:"a" description:"Link"`
		}{})), // Ollama 只接受 "json" 字符串
	}

	res, err := ChatStreamWithDeserialize[[]struct {
		A string `json:"a" description:"Link"`
	}](client, req)
	if err != nil {
		t.Logf("Chat failed (this is expected if Ollama is not running): %v", err)
	}
	for {
		select {
		case a, ok := <-res:
			if ok {
				if a.IsComplete() {
					fmt.Println(a.Data)
				}
			} else {
				return
			}
		}
	}
}

// TestChatWithProvider 测试使用指定提供者聊天
func TestChatWithProvider(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &ChatRequest{
		Model: "llama2",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	// 使用当前提供者
	_, err = client.ChatWithProvider(ProviderOllama, req)
	if err != nil {
		// 这是预期的，因为可能没有运行Ollama
		t.Logf("Chat failed (expected if Ollama is not running): %v", err)
	}

	// 使用不存在的提供者
	_, err = client.ChatWithProvider(ProviderOpenAI, req)
	if err == nil {
		t.Error("Expected error when using non-existent provider")
	}
}

// TestChatWithFallback 测试自动降级功能
func TestChatWithFallback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires real API")
	}

	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &ChatRequest{
		Model: "Qwen2",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	providers := []ProviderType{
		ProviderOllama,
	}

	_, err = client.ChatWithFallback(providers, req)
	if err != nil {
		// 这是预期的，因为可能没有运行Ollama
		t.Logf("ChatWithFallback failed (expected if Ollama is not running): %v", err)
	}
}

// TestSetLogger 测试设置日志记录器
func TestSetLogger(t *testing.T) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logger := client.GetLogger()
	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// 设置新的日志记录器
	newLogger := &xlog.LogrusAdapter{}
	client.SetLogger(newLogger)

	currentLogger := client.GetLogger()
	if currentLogger != newLogger {
		t.Error("Logger should be updated")
	}
}

// TestNewClientFromEnv 测试从环境变量创建客户端
func TestNewClientFromEnv(t *testing.T) {
	// 保存原始环境变量
	originalKey := os.Getenv("OLLAMA_BASE_URL")

	// 设置测试环境变量
	os.Setenv("OLLAMA_BASE_URL", "http://localhost:11434/v1")
	defer func() {
		// 恢复原始环境变量
		if originalKey != "" {
			os.Setenv("OLLAMA_BASE_URL", originalKey)
		} else {
			os.Unsetenv("OLLAMA_BASE_URL")
		}
	}()

	client, err := NewClientFromEnv(ProviderOllama)
	if err != nil {
		t.Logf("Failed to create client from env (this is expected if config is invalid): %v", err)
		return
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if client.GetProvider() != ProviderOllama {
		t.Errorf("Expected provider Ollama, got %s", client.GetProvider())
	}
}

// TestNewClient_NoProviders 测试没有提供者时创建客户端
func TestNewClient_NoProviders(t *testing.T) {
	// 保存原始环境变量
	originalKeys := map[string]string{
		"OPENAI_API_KEY":   os.Getenv("OPENAI_API_KEY"),
		"DEEPSEEK_API_KEY": os.Getenv("DEEPSEEK_API_KEY"),
		"CLAUDE_API_KEY":   os.Getenv("CLAUDE_API_KEY"),
		"OLLAMA_BASE_URL":  os.Getenv("OLLAMA_BASE_URL"),
		"LOCALAI_BASE_URL": os.Getenv("LOCALAI_BASE_URL"),
	}

	// 清除所有环境变量
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("DEEPSEEK_API_KEY")
	os.Unsetenv("CLAUDE_API_KEY")
	os.Unsetenv("OLLAMA_BASE_URL")
	os.Unsetenv("LOCALAI_BASE_URL")

	defer func() {
		// 恢复原始环境变量
		for key, value := range originalKeys {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	_, err := NewClient()
	if err == nil {
		t.Error("Expected error when no providers are available")
	}
}

// TestGetConfigFromEnv 测试从环境变量获取配置
func TestGetConfigFromEnv(t *testing.T) {
	// 测试Ollama（不需要API密钥）
	config := getConfigFromEnv(ProviderOllama)
	if config == nil {
		t.Error("Ollama config should not be nil (uses default URL)")
	}

	if config.BaseURL == "" {
		t.Error("Ollama BaseURL should not be empty")
	}
}

// BenchmarkClientCreation 基准测试客户端创建
func BenchmarkClientCreation(b *testing.B) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewClientWithProvider(ProviderOllama, config)
		if err != nil {
			b.Fatalf("Failed to create client: %v", err)
		}
	}
}

// BenchmarkRegisterFunction 基准测试函数注册
func BenchmarkRegisterFunction(b *testing.B) {
	config := &Config{
		BaseURL: "http://localhost:11434/v1",
	}

	client, err := NewClientWithProvider(ProviderOllama, config)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.RegisterFunction("test", "test function", func() {})
		if err != nil {
			b.Fatalf("Failed to register function: %v", err)
		}
	}
}

// TestSetFormat 测试设置格式
func TestSetFormat(t *testing.T) {
	req := &ChatRequest{
		Model: "llama3.1",
		Messages: []Message{
			{Role: "user", Content: "Test"},
		},
	}

	// 设置格式
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
		},
	}

	req = SetFormat(req, schema)
	if req.Format == nil {
		t.Error("Format should be set")
	}
}

// TestNewJSONSchema 测试创建JSON Schema
func TestNewJSONSchema(t *testing.T) {
	properties := map[string]interface{}{
		"name": NewPropertySchema("string", "Name"),
		"age":  NewPropertySchema("integer", "Age"),
	}

	schema := NewJSONSchema(properties, []string{"name", "age"})

	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	if schema["required"] == nil {
		t.Error("Required fields should be set")
	}
}

// TestNewPropertySchema 测试创建属性Schema
func TestNewPropertySchema(t *testing.T) {
	prop := NewPropertySchema("string", "Test description")

	if prop["type"] != "string" {
		t.Error("Property type should be string")
	}

	if prop["description"] != "Test description" {
		t.Error("Property description should be set")
	}
}

// TestNewArraySchema 测试创建数组Schema
func TestNewArraySchema(t *testing.T) {
	itemSchema := NewPropertySchema("string")
	arraySchema := NewArraySchema(itemSchema)

	if arraySchema["type"] != "array" {
		t.Error("Array schema type should be array")
	}

	if arraySchema["items"] == nil {
		t.Error("Array schema should have items")
	}
}

// TestFormatFromStruct 测试从结构体生成Format
func TestFormatFromStruct(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email,omitempty"`
	}

	schema, err := FormatFromStruct(User{})
	if err != nil {
		t.Fatalf("Failed to generate schema from struct: %v", err)
	}

	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties should be a map")
	}

	if properties["name"] == nil {
		t.Error("Name property should exist")
	}

	if properties["age"] == nil {
		t.Error("Age property should exist")
	}
}

// TestFormatFromJSON 测试从JSON字符串生成Format
func TestFormatFromJSON(t *testing.T) {
	jsonStr := `{"name": "John", "age": 30, "email": "john@example.com"}`

	schema, err := FormatFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to generate schema from JSON: %v", err)
	}

	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties should be a map")
	}

	if properties["name"] == nil {
		t.Error("Name property should exist")
	}
}

// TestFormatFromValue 测试从值生成Format
func TestFormatFromValue(t *testing.T) {
	// 测试map
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	schema, err := FormatFromValue(data)
	if err != nil {
		t.Fatalf("Failed to generate schema from map: %v", err)
	}

	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	// 测试数组
	arrayData := []interface{}{"a", "b", "c"}
	arraySchema, err := FormatFromValue(arrayData)
	if err != nil {
		t.Fatalf("Failed to generate schema from array: %v", err)
	}

	if arraySchema["type"] != "array" {
		t.Error("Array schema type should be array")
	}
}

// TestFormatFromType 测试泛型方法
func TestFormatFromType(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	schema, err := FormatFromType[User]()
	if err != nil {
		t.Fatalf("Failed to generate schema from type: %v", err)
	}

	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}
}

// TestSetFormatFromStruct 测试设置Format从结构体
func TestSetFormatFromStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	req := &ChatRequest{
		Model: "test",
		Messages: []Message{
			{Role: "user", Content: "test"},
		},
	}

	req, err := SetFormatFromStruct(req, User{})
	if err != nil {
		t.Fatalf("Failed to set format from struct: %v", err)
	}

	if req.Format == nil {
		t.Error("Format should be set")
	}
}
