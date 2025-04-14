package mongoutil

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StructToBSONM 将结构体转换为 bson.M
// 参数：
//   - v: 要转换的结构体
//   - ignoreZeroValue: 是否忽略零值字段
//   - tagName: 要使用的标签名称，通常是 "bson"
func StructToBSONM(v interface{}, ignoreZeroValue bool, tagName string) bson.M {
	result := bson.M{}

	// 如果传入的是nil，返回空的bson.M
	if v == nil {
		return result
	}

	val := reflect.ValueOf(v)

	// 如果传入的是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
	}

	// 确保我们处理的是结构体
	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()

	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出的字段
		if !fieldType.IsExported() {
			continue
		}

		// 解析标签
		tags := parseTag(fieldType.Tag.Get(tagName))
		if tags.name == "-" {
			continue // 跳过标记为"-"的字段
		}

		// 处理嵌入式结构体
		if fieldType.Anonymous && tags.name == "" {
			// 如果是内联嵌入式结构体
			if tags.inline {
				embedM := StructToBSONM(field.Interface(), ignoreZeroValue, tagName)
				for k, v := range embedM {
					result[k] = v
				}
				continue
			}
		}

		fieldName := tags.name
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		// 处理零值字段
		if ignoreZeroValue && isZeroValue(field) && !tags.force {
			continue
		}
		// 特殊处理 ObjectID
		if field.Type().String() == "primitive.ObjectID" {
			objectID := field.Interface().(primitive.ObjectID)
			if ignoreZeroValue && objectID.IsZero() {
				continue
			}
		}
		// 根据字段类型进行处理
		fieldValue := getFieldValue(field, ignoreZeroValue, tagName)

		// 如果未设置忽略零值或字段值不为零值，则添加到结果中
		if fieldValue != nil {
			result[fieldName] = fieldValue
		}
	}

	return result
}

// parseTag 解析字段标签
func parseTag(tag string) struct {
	name   string
	force  bool
	inline bool
} {
	if tag == "" {
		return struct {
			name   string
			force  bool
			inline bool
		}{}
	}

	parts := strings.Split(tag, ",")
	name := parts[0]

	// 处理选项
	force := false
	inline := false
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			// 不强制包含空值
		} else if opt == "!omitempty" || opt == "force" {
			force = true
		} else if opt == "inline" {
			inline = true
		}
	}

	return struct {
		name   string
		force  bool
		inline bool
	}{name, force, inline}
}

// isZeroValue 检查字段是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Struct:
		// 对于时间类型特殊处理
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		// 对于其他结构体，检查每个字段
		return false // 简化处理，实际中可能需要更复杂的逻辑
	case reflect.Map, reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// getFieldValue 获取字段的值（支持复杂类型和嵌套）
func getFieldValue(field reflect.Value, ignoreZeroValue bool, tagName string) interface{} {
	// 处理指针类型
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil
		}
		return getFieldValue(field.Elem(), ignoreZeroValue, tagName)
	}

	switch field.Kind() {
	case reflect.Struct:
		// 处理特殊类型
		if t, ok := field.Interface().(primitive.ObjectID); ok {
			return t
		}
		if t, ok := field.Interface().(time.Time); ok {
			return t
		}

		// 对普通结构体递归转换
		return StructToBSONM(field.Interface(), ignoreZeroValue, tagName)

	case reflect.Slice, reflect.Array:
		if field.Len() == 0 {
			return nil
		}

		// 处理原始字节数组
		if field.Type().Elem().Kind() == reflect.Uint8 {
			return field.Interface()
		}

		// 处理切片类型
		result := make([]interface{}, field.Len())
		for i := 0; i < field.Len(); i++ {
			result[i] = getFieldValue(field.Index(i), ignoreZeroValue, tagName)
		}
		return result

	case reflect.Map:
		if field.Len() == 0 {
			return nil
		}

		result := bson.M{}
		iter := field.MapRange()
		for iter.Next() {
			k := iter.Key().String() // 简化处理，仅支持字符串键
			v := getFieldValue(iter.Value(), ignoreZeroValue, tagName)
			if v != nil {
				result[k] = v
			}
		}
		return result

	default:
		// 基本类型直接返回
		return field.Interface()
	}
}

// 辅助函数：将结构体转化为用于更新的 bson.D（带有 $set 操作符）
func StructToUpdateBSON(v interface{}, ignoreZeroValue bool) bson.D {
	m := StructToBSONM(v, ignoreZeroValue, "bson")
	// 删除不可修改的字段
	delete(m, "_id") // 移除 _id 字段
	if len(m) == 0 {
		return bson.D{}
	}
	return bson.D{{"$set", m}}
}
