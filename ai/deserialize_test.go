package ai

import (
	"testing"
)

type mockDeserializeProvider struct {
	chatResp       *ChatResponse
	chatErr        error
	streamRespChan chan *ChatResponse
	streamErr      error
}

func (m *mockDeserializeProvider) Chat(req *ChatRequest) (*ChatResponse, error) {
	return m.chatResp, m.chatErr
}

func (m *mockDeserializeProvider) ChatStream(req *ChatRequest) (<-chan *ChatResponse, error) {
	return m.streamRespChan, m.streamErr
}

func (m *mockDeserializeProvider) ChatWithTools(req *ChatRequest, tools []Tool) (*ChatResponse, error) {
	// 简化：测试里不关心 tools，直接复用 Chat
	return m.Chat(req)
}

func (m *mockDeserializeProvider) GetName() string       { return "mock" }
func (m *mockDeserializeProvider) GetModels() []string   { return []string{"mock-model"} }
func (m *mockDeserializeProvider) ValidateConfig() error { return nil }

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestChatWithDeserialize_MessageContent(t *testing.T) {
	p := &mockDeserializeProvider{
		chatResp: &ChatResponse{
			Choices: []Choice{
				{Message: Message{Role: "assistant", Content: `{"name":"张三","age":25}`}},
			},
		},
	}

	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	got, err := ChatWithDeserialize[person](c, &ChatRequest{Model: "mock-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "张三" || got.Age != 25 {
		t.Fatalf("unexpected result: %+v", *got)
	}
}

func TestChatWithDeserialize_ToolCallArguments(t *testing.T) {
	p := &mockDeserializeProvider{
		chatResp: &ChatResponse{
			Choices: []Choice{
				{
					Message: Message{
						Role: "assistant",
						ToolCalls: []ToolCall{
							{
								ID:   "1",
								Type: "function",
								Function: ToolCallFunction{
									Name:      "extract_structured_data",
									Arguments: `{"name":"李四","age":30}`,
								},
							},
						},
					},
				},
			},
		},
	}

	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	got, err := ChatWithDeserialize[person](c, &ChatRequest{Model: "mock-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "李四" || got.Age != 30 {
		t.Fatalf("unexpected result: %+v", *got)
	}
}

func TestChatStreamWithDeserialize_AccumulateAndFinalize(t *testing.T) {
	ch := make(chan *ChatResponse, 10)
	ch <- &ChatResponse{
		Choices: []Choice{
			{Delta: Message{Role: "assistant", Content: `{"name":"`}},
		},
	}
	ch <- &ChatResponse{
		Choices: []Choice{
			{Delta: Message{Content: `王五","age":`}},
		},
	}
	ch <- &ChatResponse{
		Choices: []Choice{
			{Delta: Message{Content: `40}`}, FinishReason: "stop"},
		},
	}
	close(ch)

	p := &mockDeserializeProvider{streamRespChan: ch}
	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	out, err := ChatStreamWithDeserialize[person](c, &ChatRequest{Model: "mock-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var finals []*StreamResult[person]
	for r := range out {
		if r.IsComplete() {
			finals = append(finals, r)
		}
	}
	if len(finals) == 0 {
		t.Fatalf("expected a final result")
	}
	final := finals[len(finals)-1]
	if final.Err != nil {
		t.Fatalf("unexpected final error: %v, accumulated=%s", final.Err, final.Accumulated)
	}
	if final.Data.Name != "王五" || final.Data.Age != 40 {
		t.Fatalf("unexpected final data: %+v", final.Data)
	}
}

func TestChatStreamWithDeserialize_FinalOnlyOnce(t *testing.T) {
	ch := make(chan *ChatResponse, 10)
	ch <- &ChatResponse{
		Choices: []Choice{
			{Delta: Message{Role: "assistant", Content: `{"name":"`}},
		},
	}
	ch <- &ChatResponse{
		Choices: []Choice{
			{Delta: Message{Content: `王五","age":40}`}, FinishReason: "stop"},
		},
	}
	close(ch)

	p := &mockDeserializeProvider{streamRespChan: ch}
	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	out, err := ChatStreamWithDeserialize[person](c, &ChatRequest{Model: "mock-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	finalCount := 0
	for r := range out {
		if r.IsComplete() {
			finalCount++
		}
	}
	if finalCount != 1 {
		t.Fatalf("expected exactly 1 final item, got %d", finalCount)
	}
}

func TestChatWithDeserialize_UnwrapItemsIntoSlice(t *testing.T) {
	type qa struct {
		Q string `json:"q"`
		A string `json:"a"`
	}

	p := &mockDeserializeProvider{
		chatResp: &ChatResponse{
			Choices: []Choice{
				{Message: Message{Role: "assistant", Content: `{"items":[{"q":"什么是五险一金？","a":"..."}]}`}},
			},
		},
	}

	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	got, err := ChatWithDeserialize[[]qa](c, &ChatRequest{Model: "mock-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*got) != 1 || (*got)[0].Q == "" {
		t.Fatalf("unexpected result: %+v", *got)
	}
}

func TestChatUnmarshalInto(t *testing.T) {
	p := &mockDeserializeProvider{
		chatResp: &ChatResponse{
			Choices: []Choice{
				{Message: Message{Role: "assistant", Content: `{"name":"赵六","age":18}`}},
			},
		},
	}

	c := &Client{
		providers: map[ProviderType]AIProvider{
			ProviderOpenAI: p,
		},
		currentProvider: ProviderOpenAI,
	}

	var got person
	if err := c.ChatUnmarshalInto(&ChatRequest{Model: "mock-model"}, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "赵六" || got.Age != 18 {
		t.Fatalf("unexpected result: %+v", got)
	}
}
