package ai

// SetFormat 为请求设置结构化输出格式（JSON Schema）
// format 可以是 map[string]interface{} 或其他可以序列化为 JSON Schema 的类型
func SetFormat(req *ChatRequest, format interface{}) *ChatRequest {
	if req == nil {
		return nil
	}
	req.Format = format
	return req
}

// WithFormat 创建一个带格式设置的请求构建器（链式调用）
func WithFormat(req *ChatRequest) *FormatBuilder {
	return &FormatBuilder{req: req}
}

// FormatBuilder 格式构建器
type FormatBuilder struct {
	req *ChatRequest
}

// Set 设置格式
func (b *FormatBuilder) Set(format interface{}) *ChatRequest {
	return SetFormat(b.req, format)
}

// CountrySchema 示例：国家信息的JSON Schema
var CountrySchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"name": map[string]interface{}{
			"type": "string",
		},
		"capital": map[string]interface{}{
			"type": "string",
		},
		"languages": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "string",
			},
		},
	},
	"required": []string{"name", "capital", "languages"},
}

// PetSchema 示例：宠物信息的JSON Schema
var PetSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"name": map[string]interface{}{
			"type": "string",
		},
		"animal": map[string]interface{}{
			"type": "string",
		},
		"age": map[string]interface{}{
			"type": "integer",
		},
		"color": map[string]interface{}{
			"type":        "string",
			"description": "Pet color (optional)",
		},
		"favorite_toy": map[string]interface{}{
			"type":        "string",
			"description": "Favorite toy (optional)",
		},
	},
	"required": []string{"name", "animal", "age"},
}

// PetListSchema 示例：宠物列表的JSON Schema
var PetListSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"pets": map[string]interface{}{
			"type":  "array",
			"items": PetSchema,
		},
	},
	"required": []string{"pets"},
}

// NewJSONSchema 创建新的JSON Schema对象
// 这是一个便捷函数，用于创建标准的JSON Schema结构
func NewJSONSchema(properties map[string]interface{}, required []string) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// NewArraySchema 创建数组类型的JSON Schema
func NewArraySchema(itemSchema map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":  "array",
		"items": itemSchema,
	}
}

// NewPropertySchema 创建属性Schema的便捷函数
func NewPropertySchema(propType string, description ...string) map[string]interface{} {
	schema := map[string]interface{}{
		"type": propType,
	}
	if len(description) > 0 && description[0] != "" {
		schema["description"] = description[0]
	}
	return schema
}
