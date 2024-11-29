package file_serialize

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// TextFileSerializer 实现 SerializedInterface 的 txt 文件解析器
type TextFileSerializer struct {
	// MapToStruct 解析时将数据映射到结构体
	// fields 是字段名，headers 是标题,backDao 是映射到的结构体
	MapToStruct func(fields []string, headers []string, backDao any) any
	Convert     func(backDao any) any
}

// ReadByBytes 从字节切片解析内容
func (t TextFileSerializer) ReadByBytes(buf []byte, backDao any) []any {
	reader := bufio.NewReader(strings.NewReader(string(buf)))
	defer reader.UnreadByte()
	return t.read(reader, backDao)
}

// ReadByFile 通过文件句柄解析内容
func (t TextFileSerializer) ReadByFile(file *os.File, backDao any) []any {
	reader := bufio.NewReader(file)
	return t.read(reader, backDao)
}

// ReadByPath 通过文件路径解析内容
func (t TextFileSerializer) ReadByPath(path string, backDao any) []any {
	file, err := os.Open(path)
	if err != nil {
		errors.New("error opening file: " + err.Error())
		return nil
	}
	defer file.Close()
	return t.ReadByFile(file, backDao)
}

// Write 将数据写入到指定路径的 txt 文件
func (t TextFileSerializer) Write(path string, sourceDao any) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	switch v := sourceDao.(type) {
	case []any:
		for _, record := range v {
			if recordStr, ok := record.(string); ok {
				_, err := writer.WriteString(recordStr + "\n")
				if err != nil {
					return err
				}
			}
		}
		writer.Flush()
		return nil
	default:
		return errors.New("unsupported sourceDao type for writing")
	}
}

// read 是解析通用逻辑，支持从 reader 中读取数据
func (t TextFileSerializer) read(reader *bufio.Reader, backDao any) []any {
	var result []any
	var headers []string
	isFirstLine := true

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		// 可能需要跳过的情况，表头注释
		if len(fields) <= 1 {
			continue
		}
		// 如果是第一行，解析表头
		if isFirstLine {
			headers = fields
			isFirstLine = false
			continue
		}
		newRecord := reflect.New(reflect.TypeOf(backDao).Elem())
		if newRecord.Kind() != reflect.Ptr || newRecord.Elem().Kind() != reflect.Struct {
			errors.New("backDao must be a pointer to a struct")
			return nil
		}
		// 如果提供了backDao，按结构体字段填充
		if backDao != nil {
			if t.MapToStruct == nil {
				t.MapToStruct = func(fields []string, headers []string, backDao any) any {
					val := reflect.ValueOf(backDao)
					elem := val.Elem()

					for i, header := range headers {
						var field reflect.Value
						// 如果有tag:text，那么按照tag来查找
						for j := 0; j < elem.NumField(); j++ {
							if elem.Type().Field(j).Tag.Get("text") == header {
								field = elem.Field(j)
								break
							}
						}
						// 如果没有tag:text，那么按照字段名来查找
						if !field.IsValid() {
							field = elem.FieldByName(header)
						}
						if !field.IsValid() || !field.CanSet() {
							continue // 跳过未找到的或不可设置的字段
						}
						if i >= len(fields) {
							fields = append(fields, "")
						}
						fieldValue := fields[i]
						switch field.Kind() {
						case reflect.String:
							field.SetString(strings.TrimSpace(fieldValue))
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							intValue, _ := strconv.ParseInt(fieldValue, 10, 64)
							field.SetInt(intValue)
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							uintValue, _ := strconv.ParseUint(fieldValue, 10, 64)
							field.SetUint(uintValue)
						case reflect.Float32, reflect.Float64:
							floatValue, _ := strconv.ParseFloat(fieldValue, 64)
							field.SetFloat(floatValue)
						case reflect.Bool:
							boolValue, _ := strconv.ParseBool(fieldValue)
							field.SetBool(boolValue)
						case reflect.Slice:
							if field.Type().Elem().Kind() == reflect.Uint8 {
								data := []byte(fieldValue)
								field.Set(reflect.ValueOf(data))
							}
						default:
							errors.New("unsupported field type: " + field.Type().String())

						}
					}
					return val.Interface()
				}
			}
			record := t.MapToStruct(fields, headers, newRecord.Interface())
			// 如果有需要转换的字段，可以通过convert进行字段转换
			// 支持不同结构体转换
			if t.Convert != nil {
				record = t.Convert(record)
			}
			result = append(result, record)
		} else {
			result = append(result, fields)
		}
	}

	return result
}
