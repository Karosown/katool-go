package tool

import (
	"encoding/json"
	"fmt"
	"github.com/karosown/katool-go/ai/aiconfig"
	"reflect"
	"strings"
)

// FunctionWrapper 函数调用封装器
type FunctionWrapper struct {
	Name        string                 `json:"name"`        // 函数名称
	Description string                 `json:"description"` // 函数描述
	Function    interface{}            `json:"-"`           // 实际的Go函数
	Parameters  map[string]interface{} `json:"parameters"`  // 参数定义
}

// FunctionRegistry 函数注册表
type FunctionRegistry struct {
	functions map[string]*FunctionWrapper
}

// NewFunctionRegistry 创建新的函数注册表
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]*FunctionWrapper),
	}
}

// RegisterFunction 注册函数
func (r *FunctionRegistry) RegisterFunction(name, description string, fn interface{}) error {
	// 验证函数类型
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a Go function, got %s", fnType.Kind())
	}

	// 生成参数定义
	parameters, err := r.generateParameters(fnType)
	if err != nil {
		return fmt.Errorf("failed to generate parameters for function %s: %v", name, err)
	}

	// 注册函数
	r.functions[name] = &FunctionWrapper{
		Name:        name,
		Description: description,
		Function:    fn,
		Parameters:  parameters,
	}

	return nil
}

// generateParameters 生成函数参数定义
func (r *FunctionRegistry) generateParameters(fnType reflect.Type) (map[string]interface{}, error) {
	parameters := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := parameters["properties"].(map[string]interface{})
	required := parameters["required"].([]string)

	// 遍历函数参数
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i+1)

		// 生成参数定义
		paramDef, err := r.generateParameterDefinition(paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to generate parameter definition for %s: %v", paramName, err)
		}

		properties[paramName] = paramDef
		required = append(required, paramName)
	}

	parameters["required"] = required
	return parameters, nil
}

// generateParameterDefinition 生成单个参数定义
func (r *FunctionRegistry) generateParameterDefinition(paramType reflect.Type) (map[string]interface{}, error) {
	paramDef := map[string]interface{}{}

	switch paramType.Kind() {
	case reflect.String:
		paramDef["type"] = "string"
		paramDef["description"] = "字符串参数"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		paramDef["type"] = "integer"
		paramDef["description"] = "整数参数"
	case reflect.Float32, reflect.Float64:
		paramDef["type"] = "number"
		paramDef["description"] = "数字参数"
	case reflect.Bool:
		paramDef["type"] = "boolean"
		paramDef["description"] = "布尔参数"
	case reflect.Slice:
		// 处理切片类型
		elemType := paramType.Elem()
		paramDef["type"] = "array"
		paramDef["description"] = "数组参数"

		// 生成元素类型定义
		elemDef, err := r.generateParameterDefinition(elemType)
		if err != nil {
			return nil, err
		}
		paramDef["items"] = elemDef
	case reflect.Map:
		// 处理map类型
		paramDef["type"] = "object"
		paramDef["description"] = "对象参数"
		paramDef["additionalProperties"] = true
	case reflect.Struct:
		// 处理结构体类型
		paramDef["type"] = "object"
		paramDef["description"] = "对象参数"

		// 生成结构体字段定义
		properties := make(map[string]interface{})
		required := []string{}

		for i := 0; i < paramType.NumField(); i++ {
			field := paramType.Field(i)
			fieldName := field.Name

			// 检查json标签
			if jsonTag := field.Tag.Get("json"); jsonTag != "" {
				if jsonTag != "-" {
					parts := strings.Split(jsonTag, ",")
					fieldName = parts[0]
				}
			}

			// 生成字段定义
			fieldDef, err := r.generateParameterDefinition(field.Type)
			if err != nil {
				return nil, err
			}

			properties[fieldName] = fieldDef
			required = append(required, fieldName)
		}

		paramDef["properties"] = properties
		paramDef["required"] = required
	default:
		return nil, fmt.Errorf("unsupported parameter type: %s", paramType.Kind())
	}

	return paramDef, nil
}

// GetTools 获取所有注册的工具
func (r *FunctionRegistry) GetTools() []aiconfig.Tool {
	tools := make([]aiconfig.Tool, 0, len(r.functions))

	for _, wrapper := range r.functions {
		tools = append(tools, aiconfig.Tool{
			Type: "function",
			Function: aiconfig.ToolFunction{
				Name:        wrapper.Name,
				Description: wrapper.Description,
				Parameters:  wrapper.Parameters,
			},
		})
	}

	return tools
}

// CallFunction 调用函数
func (r *FunctionRegistry) CallFunction(name string, arguments string) (interface{}, error) {
	wrapper, exists := r.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// 解析参数
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 获取函数类型
	fnType := reflect.TypeOf(wrapper.Function)
	fnValue := reflect.ValueOf(wrapper.Function)

	// 准备函数参数
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i+1)

		paramValue, exists := params[paramName]
		if !exists {
			return nil, fmt.Errorf("missing parameter: %s", paramName)
		}

		// 转换参数类型
		convertedValue, err := r.convertParameter(paramValue, paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert parameter %s: %v", paramName, err)
		}

		args[i] = convertedValue
	}

	// 调用函数
	results := fnValue.Call(args)

	// 处理返回值
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		// 多个返回值，返回切片
		values := make([]interface{}, len(results))
		for i, result := range results {
			values[i] = result.Interface()
		}
		return values, nil
	}
}

// convertParameter 转换参数类型
func (r *FunctionRegistry) convertParameter(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	// 如果已经是目标类型，直接返回
	if reflect.TypeOf(value) == targetType {
		return reflect.ValueOf(value), nil
	}

	// 根据目标类型进行转换
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(fmt.Sprintf("%v", value)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case float64:
			return reflect.ValueOf(int64(v)).Convert(targetType), nil
		case int:
			return reflect.ValueOf(v).Convert(targetType), nil
		case int64:
			return reflect.ValueOf(v).Convert(targetType), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %T to %s", value, targetType.Kind())
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float64:
			return reflect.ValueOf(v).Convert(targetType), nil
		case int:
			return reflect.ValueOf(float64(v)).Convert(targetType), nil
		case int64:
			return reflect.ValueOf(float64(v)).Convert(targetType), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %T to %s", value, targetType.Kind())
		}
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			return reflect.ValueOf(v), nil
		case string:
			return reflect.ValueOf(v == "true"), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot convert %T to bool", value)
		}
	case reflect.Slice:
		// 处理切片类型
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			sliceValue := reflect.ValueOf(value)
			targetSlice := reflect.MakeSlice(targetType, sliceValue.Len(), sliceValue.Len())

			for i := 0; i < sliceValue.Len(); i++ {
				elemValue, err := r.convertParameter(sliceValue.Index(i).Interface(), targetType.Elem())
				if err != nil {
					return reflect.Value{}, err
				}
				targetSlice.Index(i).Set(elemValue)
			}

			return targetSlice, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %T to slice", value)
	case reflect.Map:
		// 处理map类型
		if reflect.TypeOf(value).Kind() == reflect.Map {
			mapValue := reflect.ValueOf(value)
			targetMap := reflect.MakeMap(targetType)

			for _, key := range mapValue.MapKeys() {
				value := mapValue.MapIndex(key)
				convertedKey, err := r.convertParameter(key.Interface(), targetType.Key())
				if err != nil {
					return reflect.Value{}, err
				}
				convertedValue, err := r.convertParameter(value.Interface(), targetType.Elem())
				if err != nil {
					return reflect.Value{}, err
				}
				targetMap.SetMapIndex(convertedKey, convertedValue)
			}

			return targetMap, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %T to map", value)
	case reflect.Struct:
		// 处理结构体类型
		if reflect.TypeOf(value).Kind() == reflect.Map {
			mapValue := reflect.ValueOf(value)
			targetStruct := reflect.New(targetType).Elem()

			for i := 0; i < targetType.NumField(); i++ {
				field := targetType.Field(i)
				fieldName := field.Name

				// 检查json标签
				if jsonTag := field.Tag.Get("json"); jsonTag != "" {
					if jsonTag != "-" {
						parts := strings.Split(jsonTag, ",")
						fieldName = parts[0]
					}
				}

				// 查找对应的值
				for _, key := range mapValue.MapKeys() {
					if key.String() == fieldName {
						value := mapValue.MapIndex(key)
						convertedValue, err := r.convertParameter(value.Interface(), field.Type)
						if err != nil {
							return reflect.Value{}, err
						}
						targetStruct.Field(i).Set(convertedValue)
						break
					}
				}
			}

			return targetStruct, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %T to struct", value)
	default:
		return reflect.Value{}, fmt.Errorf("unsupported conversion to %s", targetType.Kind())
	}
}

// GetFunctionNames 获取所有注册的函数名称
func (r *FunctionRegistry) GetFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// HasFunction 检查函数是否已注册
func (r *FunctionRegistry) HasFunction(name string) bool {
	_, exists := r.functions[name]
	return exists
}
