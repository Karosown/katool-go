package ai

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/karosown/katool-go/ai/aiconfig"
)

// FormatFromStruct 从Go结构体自动生成JSON Schema格式
// 支持结构体、map、JSON字符串等多种输入
// 示例：
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	schema := FormatFromStruct(User{})
func FormatFromStruct(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// 处理指针类型
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			typ = typ.Elem()
			val = reflect.New(typ).Elem()
		} else {
			typ = typ.Elem()
			val = val.Elem()
		}
	}

	// 处理数组/切片类型
	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		schema := map[string]interface{}{
			"type": "array",
		}
		elemType := typ.Elem()
		elemSchema, err := generateFieldSchema(elemType)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema for array element: %v", err)
		}
		schema["items"] = elemSchema
		return schema, nil
	}

	schema := make(map[string]interface{})
	schema["type"] = "object"

	// 添加结构体名称作为 title
	if typ.Kind() == reflect.Struct {
		schema["title"] = typ.Name()
	}

	properties := make(map[string]interface{})
	required := []string{}

	if typ.Kind() == reflect.Struct {
		// 处理结构体级别的 description
		if typ.NumField() > 0 {
			// 检查第一个字段是否有结构体级别的 comment（通过特殊字段）
			// 或者使用包级别的注释
		}

		// 处理结构体字段
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// 跳过未导出的字段
			if !field.IsExported() {
				continue
			}

			// 获取JSON标签
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue // 跳过忽略的字段
			}

			fieldName := field.Name
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				fieldName = parts[0]
				if fieldName == "" {
					fieldName = field.Name
				}
			}

			// 检查是否必需（默认必需，除非有omitempty标签）
			isOptional := strings.Contains(jsonTag, "omitempty")
			if !isOptional {
				required = append(required, fieldName)
			}

			// 生成字段Schema
			fieldSchema, err := generateFieldSchema(field.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to generate schema for field %s: %v", fieldName, err)
			}

			// 添加字段的额外信息
			enrichFieldSchema(fieldSchema, field)

			properties[fieldName] = fieldSchema
		}
	} else if typ.Kind() == reflect.Map {
		// 处理map类型 - 根据map的实际内容生成properties
		// 如果map的key是string类型，则根据实际键值对生成schema
		if typ.Key().Kind() == reflect.String {
			mapKeys := val.MapKeys()
			for _, key := range mapKeys {
				keyStr := key.String()
				value := val.MapIndex(key)

				if !value.IsValid() || value.IsNil() {
					continue
				}

				// 根据值的类型生成schema
				fieldSchema, err := generateFieldSchemaFromValue(value)
				if err != nil {
					// 如果无法推断类型，使用通用的object
					fieldSchema = map[string]interface{}{"type": "object"}
				}

				// 为map的值添加描述（如果值本身是string，使用它作为example）
				if value.Kind() == reflect.String {
					fieldSchema["example"] = value.String()
				}

				properties[keyStr] = fieldSchema
				// map的所有字段默认都是必需的
				required = append(required, keyStr)
			}
		} else {
			// 非string key的map，使用additionalProperties
			schema["additionalProperties"] = true
		}
		return schema, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %s, expected struct or map", typ.Kind())
	}

	if len(properties) > 0 {
		schema["properties"] = properties
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema, nil
}

// enrichFieldSchema 为字段 schema 添加额外的元数据
// 支持的 tag：
// - description/desc: 字段描述
// - title: 字段标题
// - example: 示例值
// - enum: 枚举值（逗号分隔）
// - pattern: 正则表达式模式（字符串类型）
// - minimum: 最小值（数字类型）
// - maximum: 最大值（数字类型）
// - minLength: 最小长度（字符串类型）
// - maxLength: 最大长度（字符串类型）
func enrichFieldSchema(schema map[string]interface{}, field reflect.StructField) {
	// 添加描述
	if desc := field.Tag.Get("description"); desc != "" {
		schema["description"] = desc
	} else if desc := field.Tag.Get("desc"); desc != "" {
		schema["description"] = desc
	}

	// 添加标题
	if title := field.Tag.Get("title"); title != "" {
		schema["title"] = title
	}

	// 添加示例
	if example := field.Tag.Get("example"); example != "" {
		schema["example"] = example
	}

	// 添加枚举值
	if enumStr := field.Tag.Get("enum"); enumStr != "" {
		enums := strings.Split(enumStr, ",")
		enumValues := make([]interface{}, len(enums))
		for i, e := range enums {
			enumValues[i] = strings.TrimSpace(e)
		}
		schema["enum"] = enumValues
	}

	// 字符串类型的额外约束
	if schema["type"] == "string" {
		if pattern := field.Tag.Get("pattern"); pattern != "" {
			schema["pattern"] = pattern
		}
		if minLen := field.Tag.Get("minLength"); minLen != "" {
			schema["minLength"] = minLen
		}
		if maxLen := field.Tag.Get("maxLength"); maxLen != "" {
			schema["maxLength"] = maxLen
		}
		if format := field.Tag.Get("format"); format != "" {
			schema["format"] = format // 如: email, uri, date-time 等
		}
	}

	// 数字类型的额外约束
	if schema["type"] == "integer" || schema["type"] == "number" {
		if min := field.Tag.Get("minimum"); min != "" {
			schema["minimum"] = min
		}
		if max := field.Tag.Get("maximum"); max != "" {
			schema["maximum"] = max
		}
	}

	// 数组类型的额外约束
	if schema["type"] == "array" {
		if minItems := field.Tag.Get("minItems"); minItems != "" {
			schema["minItems"] = minItems
		}
		if maxItems := field.Tag.Get("maxItems"); maxItems != "" {
			schema["maxItems"] = maxItems
		}
	}
}

// generateFieldSchema 生成字段的JSON Schema
func generateFieldSchema(fieldType reflect.Type) (map[string]interface{}, error) {
	schema := make(map[string]interface{})

	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
		// 指针类型表示可选，但我们不在schema中标记，而是在required列表中处理
	}

	switch fieldType.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		elemType := fieldType.Elem()
		elemSchema, err := generateFieldSchema(elemType)
		if err != nil {
			return nil, err
		}
		schema["items"] = elemSchema
	case reflect.Map:
		schema["type"] = "object"
		schema["additionalProperties"] = true
	case reflect.Struct:
		// 递归处理嵌套结构体
		schema["type"] = "object"
		schema["title"] = fieldType.Name() // 添加结构体名称

		properties := make(map[string]interface{})
		required := []string{}

		for i := 0; i < fieldType.NumField(); i++ {
			field := fieldType.Field(i)
			if !field.IsExported() {
				continue
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			fieldName := field.Name
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				fieldName = parts[0]
				if fieldName == "" {
					fieldName = field.Name
				}
			}

			isOptional := strings.Contains(jsonTag, "omitempty")
			if !isOptional {
				required = append(required, fieldName)
			}

			fieldSchema, err := generateFieldSchema(field.Type)
			if err != nil {
				return nil, err
			}

			// 使用统一的 enrichFieldSchema 函数
			enrichFieldSchema(fieldSchema, field)

			properties[fieldName] = fieldSchema
		}

		if len(properties) > 0 {
			schema["properties"] = properties
		}
		if len(required) > 0 {
			schema["required"] = required
		}
	case reflect.Interface:
		// 接口类型，允许任意值
		return nil, fmt.Errorf("interface type not supported for schema generation")
	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return schema, nil
}

// generateFieldSchemaFromValue 从reflect.Value生成字段schema
func generateFieldSchemaFromValue(value reflect.Value) (map[string]interface{}, error) {
	schema := make(map[string]interface{})

	// 处理指针
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("nil pointer value")
		}
		value = value.Elem()
	}

	// 处理接口
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil, fmt.Errorf("nil interface value")
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		if value.Len() > 0 {
			// 从第一个元素推断类型
			firstElem := value.Index(0)
			elemSchema, err := generateFieldSchemaFromValue(firstElem)
			if err == nil {
				schema["items"] = elemSchema
			} else {
				schema["items"] = map[string]interface{}{}
			}
		} else {
			schema["items"] = map[string]interface{}{}
		}
	case reflect.Map:
		schema["type"] = "object"
		// 递归处理map
		if value.Type().Key().Kind() == reflect.String {
			properties := make(map[string]interface{})
			for _, key := range value.MapKeys() {
				keyStr := key.String()
				mapValue := value.MapIndex(key)
				if mapValue.IsValid() {
					propSchema, err := generateFieldSchemaFromValue(mapValue)
					if err == nil {
						properties[keyStr] = propSchema
					}
				}
			}
			if len(properties) > 0 {
				schema["properties"] = properties
			}
		} else {
			schema["additionalProperties"] = true
		}
	case reflect.Struct:
		// 对于结构体，使用类型信息
		return generateFieldSchema(value.Type())
	default:
		return nil, fmt.Errorf("unsupported value kind: %s", value.Kind())
	}

	return schema, nil
}

// FormatFromJSON 从JSON字符串自动生成JSON Schema格式
// 示例：
//
//	jsonStr := `{"name": "John", "age": 30}`
//	schema := FormatFromJSON(jsonStr)
func FormatFromJSON(jsonStr string) (map[string]interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("invalid JSON string: %v", err)
	}

	return FormatFromValue(data)
}

// FormatFromValue 从任意值（map、slice等）自动生成JSON Schema格式
// 示例：
//
//	data := map[string]interface{}{
//	    "name": "John",
//	    "age": 30,
//	}
//	schema := FormatFromValue(data)
func FormatFromValue(v interface{}) (map[string]interface{}, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		return formatFromMap(val)
	case []interface{}:
		return formatFromArray(val)
	case string:
		return map[string]interface{}{"type": "string"}, nil
	case int, int8, int16, int32, int64:
		return map[string]interface{}{"type": "integer"}, nil
	case uint, uint8, uint16, uint32, uint64:
		return map[string]interface{}{"type": "integer"}, nil
	case float32, float64:
		return map[string]interface{}{"type": "number"}, nil
	case bool:
		return map[string]interface{}{"type": "boolean"}, nil
	case nil:
		return nil, fmt.Errorf("cannot generate schema from nil value")
	default:
		// 尝试使用反射处理
		return FormatFromStruct(v)
	}
}

// formatFromMap 从map生成schema
func formatFromMap(data map[string]interface{}) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := schema["properties"].(map[string]interface{})
	required := []string{}

	for key, value := range data {
		propSchema, err := FormatFromValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema for key %s: %v", key, err)
		}
		properties[key] = propSchema
		// 所有字段都设为必需
		required = append(required, key)
	}

	schema["required"] = required
	return schema, nil
}

// formatFromArray 从数组生成schema
func formatFromArray(data []interface{}) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "array",
	}

	if len(data) > 0 {
		// 使用第一个元素作为示例
		itemSchema, err := FormatFromValue(data[0])
		if err != nil {
			return nil, err
		}
		schema["items"] = itemSchema
	} else {
		// 空数组，使用通用类型
		schema["items"] = map[string]interface{}{"type": "string"}
	}

	return schema, nil
}

// FormatFromType 根据类型创建Format（类型参数方式）
// 示例：
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	schema := FormatFromType[User]()
func FormatFromType[T any]() (map[string]interface{}, error) {
	var zero T
	return FormatFromStruct(zero)
}

// FormatArrayOf 从结构体或map生成数组格式的Schema
// 示例：
//
//	type Item struct {
//	    Q string `json:"Q"`
//	    A string `json:"A"`
//	}
//	schema := FormatArrayOf(Item{})
//	结果: {"type": "array", "items": {...Item的schema...}}
func FormatArrayOf(v interface{}) (map[string]interface{}, error) {
	itemSchema, err := FormatFromStruct(v)
	if err != nil {
		// 如果结构体失败，尝试从值生成
		itemSchema, err = FormatFromValue(v)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"type":  "array",
		"items": itemSchema,
	}, nil
}

// FormatArrayOfType 使用泛型方式生成数组格式的Schema
// 示例：
//
//	type Item struct {
//	    Q string `json:"Q"`
//	    A string `json:"A"`
//	}
//	schema := FormatArrayOfType[Item]()
func FormatArrayOfType[T any]() (map[string]interface{}, error) {
	var zero T
	return FormatArrayOf(zero)
}

// SetFormatFromStruct 从结构体自动生成并设置Format
// 便捷方法：直接为请求设置从结构体生成的格式
func SetFormatFromStruct(req *aiconfig.ChatRequest, v interface{}) (*aiconfig.ChatRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	schema, err := FormatFromStruct(v)
	if err != nil {
		return nil, err
	}
	req.Format = schema
	return req, nil
}

// SetFormatFromJSON 从JSON字符串自动生成并设置Format
func SetFormatFromJSON(req *aiconfig.ChatRequest, jsonStr string) (*aiconfig.ChatRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	schema, err := FormatFromJSON(jsonStr)
	if err != nil {
		return nil, err
	}
	req.Format = schema
	return req, nil
}

// WithFormatFrom 链式调用构建器
func WithFormatFrom(req *aiconfig.ChatRequest) *FormatFromBuilder {
	return &FormatFromBuilder{req: req}
}

// FormatFromBuilder 格式自动生成构建器
type FormatFromBuilder struct {
	req *aiconfig.ChatRequest
}

// Struct 从结构体设置格式
func (b *FormatFromBuilder) Struct(v interface{}) (*aiconfig.ChatRequest, error) {
	return SetFormatFromStruct(b.req, v)
}

// JSON 从JSON字符串设置格式
func (b *FormatFromBuilder) JSON(jsonStr string) (*aiconfig.ChatRequest, error) {
	return SetFormatFromJSON(b.req, jsonStr)
}

// Value 从值设置格式
func (b *FormatFromBuilder) Value(v interface{}) (*aiconfig.ChatRequest, error) {
	if b.req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	schema, err := FormatFromValue(v)
	if err != nil {
		return nil, err
	}
	b.req.Format = schema
	return b.req, nil
}
