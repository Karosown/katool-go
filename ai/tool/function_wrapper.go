package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/karosown/katool-go/ai/types"
)

// FunctionWrapper 鍑芥暟璋冪敤灏佽鍣?
type FunctionWrapper struct {
	Name        string                 `json:"name"`        // 鍑芥暟鍚嶇О
	Description string                 `json:"description"` // 鍑芥暟鎻忚堪
	Function    interface{}            `json:"-"`           // 瀹為檯鐨凣o鍑芥暟
	Parameters  map[string]interface{} `json:"parameters"`  // 鍙傛暟瀹氫箟
	ParamOrder  []string               `json:"-"`           // 鍙傛暟鍚嶇О瀵瑰簲椤哄簭
}

// FunctionRegistry 鍑芥暟娉ㄥ唽琛?
type FunctionRegistry struct {
	functions map[string]*FunctionWrapper
}

// NewFunctionRegistry 鍒涘缓鏂扮殑鍑芥暟娉ㄥ唽琛?
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]*FunctionWrapper),
	}
}

// RegisterFunction 娉ㄥ唽鍑芥暟
func (r *FunctionRegistry) RegisterFunction(name, description string, fn interface{}) error {
	// 楠岃瘉鍑芥暟绫诲瀷
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a Go function, got %s", fnType.Kind())
	}

	// 鐢熸垚鍙傛暟瀹氫箟
	parameters, err := r.generateParameters(fnType)
	if err != nil {
		return fmt.Errorf("failed to generate parameters for function %s: %v", name, err)
	}

	paramOrder := make([]string, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramOrder[i] = fmt.Sprintf("param%d", i+1)
	}

	// 娉ㄥ唽鍑芥暟
	r.functions[name] = &FunctionWrapper{
		Name:        name,
		Description: description,
		Function:    fn,
		Parameters:  parameters,
		ParamOrder:  paramOrder,
	}

	return nil
}

// RegisterFunctionWith 注册函数，允许手动指定参数 schema 与参数名顺序
func (r *FunctionRegistry) RegisterFunctionWith(name, description string, parameters map[string]interface{}, paramOrder []string, fn interface{}) error {
	if parameters == nil {
		return fmt.Errorf("parameters cannot be nil")
	}

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("function must be a Go function, got %s", fnType.Kind())
	}

	// 默认参数名
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

// generateParameters 鐢熸垚鍑芥暟鍙傛暟瀹氫箟
func (r *FunctionRegistry) generateParameters(fnType reflect.Type) (map[string]interface{}, error) {
	parameters := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := parameters["properties"].(map[string]interface{})
	required := parameters["required"].([]string)

	// 閬嶅巻鍑芥暟鍙傛暟
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i+1)

		// 鐢熸垚鍙傛暟瀹氫箟
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
		// 处理切片类型
		elemType := paramType.Elem()
		paramDef["type"] = "array"
		paramDef["description"] = "array parameter"

		// 生成元素类型定义
		elemDef, err := r.generateParameterDefinition(elemType)
		if err != nil {
			return nil, err
		}
		paramDef["items"] = elemDef
	case reflect.Map:
		// 处理map类型
		paramDef["type"] = "object"
		paramDef["description"] = "object parameter"
		paramDef["additionalProperties"] = true
	case reflect.Struct:
		// 处理结构体类型
		paramDef["type"] = "object"
		paramDef["description"] = "object parameter"

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

// GetTools 鑾峰彇鎵€鏈夋敞鍐岀殑宸ュ叿
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

// CallFunction 璋冪敤鍑芥暟
func (r *FunctionRegistry) CallFunction(name string, arguments string) (interface{}, error) {
	wrapper, exists := r.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// 瑙ｆ瀽鍙傛暟
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %v", err)
	}

	// 鑾峰彇鍑芥暟绫诲瀷
	fnType := reflect.TypeOf(wrapper.Function)
	fnValue := reflect.ValueOf(wrapper.Function)

	// 鍑嗗鍑芥暟鍙傛暟
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

		// 杞崲鍙傛暟绫诲瀷
		convertedValue, err := r.convertParameter(paramValue, paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert parameter %s: %v", paramName, err)
		}

		args[i] = convertedValue
	}

	// 璋冪敤鍑芥暟
	results := fnValue.Call(args)

	// 澶勭悊杩斿洖鍊?
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		// 澶氫釜杩斿洖鍊硷紝杩斿洖鍒囩墖
		values := make([]interface{}, len(results))
		for i, result := range results {
			values[i] = result.Interface()
		}
		return values, nil
	}
}

// convertParameter 杞崲鍙傛暟绫诲瀷
func (r *FunctionRegistry) convertParameter(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	// 濡傛灉宸茬粡鏄洰鏍囩被鍨嬶紝鐩存帴杩斿洖
	if reflect.TypeOf(value) == targetType {
		return reflect.ValueOf(value), nil
	}

	// 鏍规嵁鐩爣绫诲瀷杩涜杞崲
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
		// 澶勭悊鍒囩墖绫诲瀷
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
		// 澶勭悊map绫诲瀷
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
		// 澶勭悊缁撴瀯浣撶被鍨?
		if reflect.TypeOf(value).Kind() == reflect.Map {
			mapValue := reflect.ValueOf(value)
			targetStruct := reflect.New(targetType).Elem()

			for i := 0; i < targetType.NumField(); i++ {
				field := targetType.Field(i)
				fieldName := field.Name

				// 妫€鏌son鏍囩
				if jsonTag := field.Tag.Get("json"); jsonTag != "" {
					if jsonTag != "-" {
						parts := strings.Split(jsonTag, ",")
						fieldName = parts[0]
					}
				}

				// 鏌ユ壘瀵瑰簲鐨勫€?
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

// GetFunctionNames 鑾峰彇鎵€鏈夋敞鍐岀殑鍑芥暟鍚嶇О
func (r *FunctionRegistry) GetFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// HasFunction 妫€鏌ュ嚱鏁版槸鍚﹀凡娉ㄥ唽
func (r *FunctionRegistry) HasFunction(name string) bool {
	_, exists := r.functions[name]
	return exists
}
