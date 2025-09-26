package csv

import (
	"bufio"
	"encoding"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ReadBytes 通过字节切片读取 CSV，返回 []T
func ReadBytes[T any](b []byte) ([]T, error) {
	return Read[T](strings.NewReader(string(b)))
}

// ReadFile 从已打开的文件读取 CSV，返回 []T
func ReadFile[T any](f *os.File) ([]T, error) {
	return Read[T](f)
}

// ReadPath 从路径读取 CSV，返回 []T
func ReadPath[T any](path string) ([]T, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read[T](f)
}

// Read 从任意 io.Reader 读取 CSV，返回 []T
func Read[T any](r io.Reader) ([]T, error) {
	var zero []T

	// 目标类型
	tType := reflect.TypeOf((*T)(nil)).Elem()
	if tType.Kind() != reflect.Struct {
		return zero, fmt.Errorf("T must be a struct, got %s", tType.Kind())
	}

	// 构建字段映射: headerName -> fieldInfo
	fields, headersExpected := extractFields(tType)

	// 逐行扫描，寻找标题行
	br := bufio.NewReader(r)

	var headers []string
	for {
		line, err := br.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return zero, fmt.Errorf("read line: %w", err)
		}
		line = strings.TrimSpace(line)
		if line != "" {
			rec, err2 := csv.NewReader(strings.NewReader(line)).Read()
			if err2 == nil {
				for i := range rec {
					rec[i] = strings.Trim(rec[i], "\"")
				}
				if headerMatch(rec, headersExpected) {
					headers = rec
					break
				}
			}
		}
		if errors.Is(err, io.EOF) {
			return zero, errors.New("csv header not found")
		}
	}

	// 继续用同一个缓冲区之后的内容读取数据
	cr := csv.NewReader(br)

	// 建立 header 索引 -> 字段描述
	colToField := make([]fieldInfo, len(headers))
	for i, h := range headers {
		if fi, ok := fields[h]; ok {
			colToField[i] = fi
		} else {
			// 没有匹配到字段则留空（zero fieldInfo）
		}
	}

	var out []T
	var err error
	var line int64
	line = 1
	for {
		var rec []string
		rec, err = cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			err = errors.Join(err, fmt.Errorf("has error in %d line when read line: %w", line, err))
			line++
			continue
		}
		v := reflect.New(tType).Elem()
		for i, raw := range rec {
			if i >= len(colToField) || colToField[i].index == -1 {
				continue
			}
			raw = strings.Trim(raw, "\"")
			if err1 := setFieldValue(v.Field(colToField[i].index), raw); err1 != nil {
				err = errors.Join(err, fmt.Errorf("has error in %d line when read line: %w", line, err1))
				continue
			}
		}
		out = append(out, v.Interface().(T))
		line++
	}

	return out, err
}

// WritePath 将切片数据写入 CSV 文件
func WritePath[T any](path string, data []T) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write[T](f, data)
}

// Write 将切片数据写入到任意 io.Writer
func Write[T any](w io.Writer, data []T) error {
	cr := csv.NewWriter(w)
	defer cr.Flush()

	var tType reflect.Type
	if len(data) == 0 {
		// 对空数据也要写 header？这里无法得知类型字段，只能返回 nil 或错误
		return nil
	} else {
		tType = reflect.TypeOf(data[0])
	}

	if tType.Kind() != reflect.Struct {
		return fmt.Errorf("T must be a struct, got %s", tType.Kind())
	}

	fields, headers := extractFields(tType)

	// 写 header
	if err := cr.Write(headers); err != nil {
		return err
	}

	// 写数据
	for _, row := range data {
		rv := reflect.ValueOf(row)
		record := make([]string, len(headers))
		for i, h := range headers {
			fi := fields[h]
			record[i] = getFieldString(rv.Field(fi.index))
		}
		if err := cr.Write(record); err != nil {
			return err
		}
	}
	return cr.Error()
}

// ----------------- 辅助方法与类型 -----------------

type fieldInfo struct {
	index int
	name  string // header name
}

func extractFields(t reflect.Type) (map[string]fieldInfo, []string) {
	fields := make(map[string]fieldInfo)
	var headers []string

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" { // 未导出字段跳过
			continue
		}
		name := sf.Tag.Get("csv")
		if name == "" || name == "-" {
			if name == "-" {
				continue
			}
			name = sf.Name
		}

		fields[name] = fieldInfo{
			index: i,
			name:  name,
		}
		headers = append(headers, name)
	}

	return fields, headers
}

func headerMatch(actual []string, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

func setFieldValue(fv reflect.Value, s string) error {
	if !fv.CanSet() {
		return nil
	}
	ft := fv.Type()

	// 支持 TextUnmarshaler
	if fv.Addr().Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
		return fv.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
	}

	switch ft.Kind() {
	case reflect.String:
		fv.SetString(s)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s == "" {
			fv.SetInt(0)
			return nil
		}
		i, err := strconv.ParseInt(s, 10, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if s == "" {
			fv.SetUint(0)
			return nil
		}
		u, err := strconv.ParseUint(s, 10, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		if s == "" {
			fv.SetFloat(0)
			return nil
		}
		f, err := strconv.ParseFloat(s, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Pointer:
		// 支持 *T 的简单场景
		if s == "" {
			fv.Set(reflect.Zero(ft))
			return nil
		}
		elem := reflect.New(ft.Elem())
		if err := setFieldValue(elem.Elem(), s); err != nil {
			return err
		}
		fv.Set(elem)
	default:
		// 其他类型可考虑 json 反序列化或实现 TextUnmarshaler
		// 这里先返回可读错误
		return fmt.Errorf("unsupported field type: %s", ft.String())
	}
	return nil
}

func getFieldString(fv reflect.Value) string {
	if fv.Kind() == reflect.Pointer {
		if fv.IsNil() {
			return ""
		}
		fv = fv.Elem()
	}

	// 支持 TextMarshaler
	if fv.CanAddr() && fv.Addr().Type().Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) {
		b, err := fv.Addr().Interface().(encoding.TextMarshaler).MarshalText()
		if err == nil {
			return string(b)
		}
	}

	switch fv.Kind() {
	case reflect.String:
		return fv.String()
	case reflect.Bool:
		if fv.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(fv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fv.Float(), 'f', -1, fv.Type().Bits())
	default:
		// 未支持类型返回空字符串或可考虑 JSON
		return ""
	}
}
