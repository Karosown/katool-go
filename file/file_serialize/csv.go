package file_serialize

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type CSVSerializer struct{}

func (c CSVSerializer) ReadByBytes(bytes []byte, backDao any) []any {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReadByBytes panic:", err)
		}
	}()
	reader := csv.NewReader(strings.NewReader(string(bytes)))
	var result []any

	// 获取backDao的类型信息
	backDaoType := reflect.TypeOf(backDao)
	if backDaoType.Kind() != reflect.Ptr {
		panic("backDao必须是指针类型")
	}
	elemType := backDaoType.Elem()
	if elemType.Kind() != reflect.Struct {
		panic("backDao必须是结构体指针")
	}

	// 生成预期的标题行字段名列表
	var expectedHeaders []string
	fieldIndexMap := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fieldName := field.Name

		// 如果使用了标签，可以在这里处理，例如：
		// tagName := field.Tag.Get("csv")
		// if tagName != "" {
		//     fieldName = tagName
		// }

		expectedHeaders = append(expectedHeaders, fieldName)
		fieldIndexMap[fieldName] = i
	}

	// 新增：创建一个缓冲读取器，以逐行读取文件内容
	bufReader := bufio.NewReader(strings.NewReader(string(bytes)))

	var headers []string
	for {
		// 逐行读取
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				panic("未找到CSV标题行")
			}
			panic(fmt.Errorf("读取文件出错: %v", err))
		}

		// 去除行首行尾的空白字符
		line = strings.TrimSpace(line)

		// 判断是否是空行
		if line == "" {
			continue
		}

		// 尝试解析当前行
		record, err := csv.NewReader(strings.NewReader(line)).Read()
		if err != nil {
			// 如果解析失败，继续读取下一行
			continue
		}

		// 去除字段名的首尾引号
		for i := range record {
			record[i] = strings.Trim(record[i], "\"")
		}

		// 检查解析后的字段名是否与预期的标题行匹配
		if len(record) != len(expectedHeaders) {
			continue
		}

		matched := true
		for i := range record {
			if record[i] != expectedHeaders[i] {
				matched = false
				break
			}
		}

		if matched {
			// 找到了标题行
			headers = record
			break
		}
	}

	// 设置 CSV 读取器的位置，从当前位置开始读取
	reader = csv.NewReader(bufReader)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		newStructPtr := reflect.New(elemType)
		newStruct := newStructPtr.Elem()

		for i, value := range record {
			if i >= len(headers) {
				continue
			}
			fieldName := headers[i]

			fieldIdx, exists := fieldIndexMap[fieldName]
			if !exists {
				continue
			}
			fieldValue := newStruct.Field(fieldIdx)
			if !fieldValue.CanSet() {
				continue
			}

			// 去除值的首尾引号
			value = strings.Trim(value, "\"")

			// 根据字段类型进行转换
			newStruct.Field(fieldIdx).SetString(value)
		}

		result = append(result, newStructPtr.Interface())
	}
	return result
}
func (c CSVSerializer) ReadByFile(file *os.File, backDao any) []any {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("%s", err)
		}
	}()
	reader := csv.NewReader(file)
	var result []any

	// 获取backDao的类型信息
	backDaoType := reflect.TypeOf(backDao)
	if backDaoType.Kind() != reflect.Ptr {
		panic("backDao必须是指针类型")
	}
	elemType := backDaoType.Elem()
	if elemType.Kind() != reflect.Struct {
		panic("backDao必须是结构体指针")
	}

	// 生成预期的标题行字段名列表
	var expectedHeaders []string
	fieldIndexMap := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fieldName := field.Name

		// 如果使用了标签，可以在这里处理，例如：
		// tagName := field.Tag.Get("csv")
		// if tagName != "" {
		//     fieldName = tagName
		// }

		expectedHeaders = append(expectedHeaders, fieldName)
		fieldIndexMap[fieldName] = i
	}

	// 新增：创建一个缓冲读取器，以逐行读取文件内容
	bufReader := bufio.NewReader(file)

	var headers []string
	for {
		// 逐行读取
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				panic("未找到CSV标题行")
			}
			panic(fmt.Errorf("读取文件出错: %v", err))
		}

		// 去除行首行尾的空白字符
		line = strings.TrimSpace(line)

		// 判断是否是空行
		if line == "" {
			continue
		}

		// 尝试解析当前行
		record, err := csv.NewReader(strings.NewReader(line)).Read()
		if err != nil {
			// 如果解析失败，继续读取下一行
			continue
		}

		// 去除字段名的首尾引号
		for i := range record {
			record[i] = strings.Trim(record[i], "\"")
		}

		// 检查解析后的字段名是否与预期的标题行匹配
		if len(record) != len(expectedHeaders) {
			continue
		}

		matched := true
		for i := range record {
			if record[i] != expectedHeaders[i] {
				matched = false
				break
			}
		}

		if matched {
			// 找到了标题行
			headers = record
			break
		}
	}

	// 设置 CSV 读取器的位置，从当前位置开始读取
	reader = csv.NewReader(bufReader)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		newStructPtr := reflect.New(elemType)
		newStruct := newStructPtr.Elem()

		for i, value := range record {
			if i >= len(headers) {
				continue
			}
			fieldName := headers[i]

			fieldIdx, exists := fieldIndexMap[fieldName]
			if !exists {
				continue
			}
			fieldValue := newStruct.Field(fieldIdx)
			if !fieldValue.CanSet() {
				continue
			}

			// 去除值的首尾引号
			value = strings.Trim(value, "\"")

			// 根据字段类型进行转换
			newStruct.Field(fieldIdx).SetString(value)
		}

		result = append(result, newStructPtr.Interface())
	}
	return result
}

func (c CSVSerializer) ReadByPath(path string, backDao any) []any {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	return c.ReadByFile(file, backDao)
}

func (c CSVSerializer) Write(path string, sourceDao any) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("%s", err)
		}
	}()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	sourceValue := reflect.ValueOf(sourceDao)
	if sourceValue.Kind() != reflect.Slice {
		panic("sourceDao必须是切片类型")
	}

	if sourceValue.Len() == 0 {
		return nil
	}

	elemType := sourceValue.Index(0).Type()

	if elemType.Kind() != reflect.Struct {
		panic("sourceDao的元素必须是结构体类型")
	}

	// 写入标题行
	var headers []string
	for i := 0; i < elemType.NumField(); i++ {
		headers = append(headers, elemType.Field(i).Name)
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	// 写入数据行
	for i := 0; i < sourceValue.Len(); i++ {
		elem := sourceValue.Index(i)
		var record []string
		for j := 0; j < elem.NumField(); j++ {
			fieldValue := elem.Field(j)
			record = append(record, fieldValue.String())
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
