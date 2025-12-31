package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/karosown/katool-go/ai/types"
)

// FunctionWrapper 封装单个函数的元数据与实际实现。
type FunctionWrapper struct {
	Name        string                 `json:"name"`        // 函数名称
	Description string                 `json:"description"` // 函数描述
	Function    interface{}            `json:"-"`           // 实际的 Go 函数
	Parameters  map[string]interface{} `json:"parameters"`  // 参数定义（JSON Schema）
	ParamOrder  []string               `json:"-"`           // 参数顺序，用于绑定参数名
}

// FunctionRegistry 负责注册与调用函数。
type FunctionRegistry struct {
	functions map[string]*FunctionWrapper
}

// NewFunctionRegistry 创建新的函数注册表。
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]*FunctionWrapper),
	}
}

// RegisterFunction 按照 Go 函数签名自动生成参数 schema 并注册。
func (r *FunctionRegistry) RegisterFunction(name, description string, fn interface{}) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a Go function, got %s", fnType.Kind())
	}

	parameters, err := r.generateParameters(fnType)
	if err != nil {
		return fmt.Errorf("failed to generate parameters for function %s: %v", name, err)
	}

	paramOrder := make([]string, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramOrder[i] = fmt.Sprintf("param%d", i+1)
	}

	r.functions[name] = &FunctionWrapper{
		Name:        name,
		Description: description,
		Function:    fn,
		Parameters:  parameters,
		ParamOrder:  paramOrder,
	}
	return nil
}

// RegisterFunctionWith 注册函数，使用自定义参数 schema 与参数名顺序。
// 若 paramOrder 为空则使用 param1/param2... 默认顺序。
func (r *FunctionRegistry) RegisterFunctionWith(name, description string, parameters map[string]interface{}, paramOrder []string, fn interface{}) error {
	if parameters == nil {
		return fmt.Errorf("parameters cannot be nil")
	}

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a Go function, got %s", fnType.Kind())
	}

	if len(paramOrder) == 0 {
		paramOrder = make([]string, fnType.NumIn())
		for i := 0; i < fnType.NumIn(); i++ {
			paramOrder[i] = fmt.Sprintf("param%d", i+1)
		}
	}
	if len(paramOrder) != fnType.NumIn() {
		return fmt.Errorf("paramOrder length (%d) does not match function args (%d)", len(paramOrder), fnType.NumIn())
	}

	r.functions[name] = &FunctionWrapper{
		Name:        name,
		Description: description,
		Function:    fn,
		Parameters:  parameters,
		ParamOrder:  paramOrder,
	}
	return nil
}

// generateParameters 根据函数签名生成参数定义。
func (r *FunctionRegistry) generateParameters(fnType reflect.Type) (map[string]interface{}, error) {
	parameters := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := parameters["properties"].(map[string]interface{})
	required := parameters["required"].([]string)

	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i+1)

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

// generateParameterDefinition 生成单个参数定义。
func (r *FunctionRegistry) generateParameterDefinition(paramType reflect.Type) (map[string]interface{}, error) {
	paramDef := map[string]interface{}{}

	switch paramType.Kind() {
	case reflect.String:
		paramDef["type"] = "string"
		paramDef["description"] = "string parameter"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		paramDef["type"] = "integer"
		paramDef["description"] = "integer parameter"
	case reflect.Float32, reflect.Float64:
		paramDef["type"] = "number"
		paramDef["description"] = "number parameter"
	case reflect.Bool:
		paramDef["type"] = "boolean"
		paramDef["description"] = "boolean parameter"
	case reflect.Slice:
		elemType := paramType.Elem()
		paramDef["type"] = "array"
		paramDef["description"] = "array parameter"

		elemDef, err := r.generateParameterDefinition(elemType)
		if err != nil {
			return nil, err
		}
		paramDef["items"] = elemDef
	case reflect.Map:
		paramDef["type"] = "object"
		paramDef["description"] = "object parameter"
		paramDef["additionalProperties"] = true
	case reflect.Struct:
		paramDef["type"] = "object"
		paramDef["description"] = "object parameter"

		properties := make(map[string]interface{})
		required := []string{}

		for i := 0; i < paramType.NumField(); i++ {
			field := paramType.Field(i)
			fieldName := field.Name

			if jsonTag := field.Tag.Get("json"); jsonTag != "" {
				if jsonTag != "-" {
					parts := strings.Split(jsonTag, ",")
					fieldName = parts[0]
				}
			}

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

// GetTools 获取所有注册的工具定义。
func (r *FunctionRegistry) GetTools() []types.Tool {
	tools := make([]types.Tool, 0, len(r.functions))

	for _, wrapper := range r.functions {
		tools = append(tools, types.Tool{
			Type: "function",
			Function: types.ToolFunction{
				Name:        wrapper.Name,
				Description: wrapper.Description,
				Parameters:  wrapper.Parameters,
			},
		})
	}

	return tools
}

// CallFunction 调用函数。
func (r *FunctionRegistry) CallFunction(name string, arguments string) (interface{}, error) {
	wrapper, exists := r.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	var params map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %v", err)
	}

	fnType := reflect.TypeOf(wrapper.Function)
	fnValue := reflect.ValueOf(wrapper.Function)

	// 单参数且为 map：直接把完整参数对象传入
	if fnType.NumIn() == 1 && fnType.In(0).Kind() == reflect.Map {
		convertedValue, err := r.convertParameter(params, fnType.In(0))
		if err != nil {
			return nil, fmt.Errorf("failed to convert map parameter: %v", err)
		}
		results := fnValue.Call([]reflect.Value{convertedValue})
		switch len(results) {
		case 0:
			return nil, nil
		case 1:
			return results[0].Interface(), nil
		default:
			values := make([]interface{}, len(results))
			for i, result := range results {
				values[i] = result.Interface()
			}
			return values, nil
		}
	}

	// 按参数顺序准备参数
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		defaultParamName := fmt.Sprintf("param%d", i+1)
		paramName := defaultParamName
		if len(wrapper.ParamOrder) > i && wrapper.ParamOrder[i] != "" {
			paramName = wrapper.ParamOrder[i]
		}

		paramValue, exists := params[paramName]
		if !exists && paramName != defaultParamName {
			if fallbackValue, ok := params[defaultParamName]; ok {
				paramValue = fallbackValue
				exists = true
				paramName = defaultParamName
			}
		}
		if !exists {
			return nil, fmt.Errorf("missing parameter: %s", paramName)
		}

		convertedValue, err := r.convertParameter(paramValue, paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert parameter %s: %v", paramName, err)
		}

		args[i] = convertedValue
	}

	results := fnValue.Call(args)

	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		return results[0].Interface(), nil
	}

	values := make([]interface{}, len(results))
	for i, result := range results {
		values[i] = result.Interface()
	}
	return values, nil
}

// convertParameter 将参数转换为目标类型。
func (r *FunctionRegistry) convertParameter(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	if reflect.TypeOf(value) == targetType {
		return reflect.ValueOf(value), nil
	}

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
		if reflect.TypeOf(value).Kind() == reflect.Map {
			mapValue := reflect.ValueOf(value)
			targetMap := reflect.MakeMap(targetType)
			for _, key := range mapValue.MapKeys() {
				v := mapValue.MapIndex(key)
				convertedKey, err := r.convertParameter(key.Interface(), targetType.Key())
				if err != nil {
					return reflect.Value{}, err
				}
				convertedValue, err := r.convertParameter(v.Interface(), targetType.Elem())
				if err != nil {
					return reflect.Value{}, err
				}
				targetMap.SetMapIndex(convertedKey, convertedValue)
			}
			return targetMap, nil
		}
		return reflect.Value{}, fmt.Errorf("cannot convert %T to map", value)
	case reflect.Struct:
		if reflect.TypeOf(value).Kind() == reflect.Map {
			mapValue := reflect.ValueOf(value)
			targetStruct := reflect.New(targetType).Elem()
			for i := 0; i < targetType.NumField(); i++ {
				field := targetType.Field(i)
				fieldName := field.Name

				if jsonTag := field.Tag.Get("json"); jsonTag != "" {
					if jsonTag != "-" {
						parts := strings.Split(jsonTag, ",")
						fieldName = parts[0]
					}
				}

				for _, key := range mapValue.MapKeys() {
					if key.String() == fieldName {
						val := mapValue.MapIndex(key)
						convertedValue, err := r.convertParameter(val.Interface(), field.Type)
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

// GetFunctionNames 获取已注册函数名称。
func (r *FunctionRegistry) GetFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// HasFunction 检查函数是否已注册。
func (r *FunctionRegistry) HasFunction(name string) bool {
	_, exists := r.functions[name]
	return exists
}
