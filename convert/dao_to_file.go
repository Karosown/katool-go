package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
)

func StructToCSV[T any](datas []T, fullPath string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("StructToCSV panic: " + err.(error).Error())
		}
	}()
	if datas == nil || len(datas) == 0 {
		return nil
	}
	records := make([][]string, 0)
	tags := make([]string, 0)
	reflectType := reflect.TypeOf(datas[0]).Elem()
	headers := make([]string, 0)
	tagsNums := reflectType.NumField()
	for i := 0; i < tagsNums; i++ {
		field := reflectType.Field(i)
		tags = append(tags, field.Tag.Get("csv"))
		headers = append(headers, field.Name)
	}
	records = append(records, tags) // 第一行写入表头
	for _, data := range datas {
		record := make([]string, 0)
		reflectValue := reflect.ValueOf(data).Elem()
		for i := 0; i < tagsNums; i++ {
			field := reflectValue.FieldByName(headers[i])
			// 针对不同类型的数据进行转换
			switch field.Kind() {
			case reflect.String:
				record = append(record, field.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				record = append(record, strconv.FormatInt(field.Int(), 10))
			case reflect.Float32, reflect.Float64:
				record = append(record, strconv.FormatFloat(field.Float(), 'f', -1, 64))
			case reflect.Bool:
				record = append(record, strconv.FormatBool(field.Bool()))
			case reflect.Struct:
				if field.Type().String() == "time.Time" {
					elems := field.Interface().(time.Time)
					// 将时间格式化为指定格式
					record = append(record, elems.Format(time.DateTime))
				} else if field.Type().String() == "[]uint8" && field.Len() > 0 {
					record = append(record, string(field.Bytes()))
				} else {
					return errors.New("unsupported field type: " + field.Type().String())
				}
			}
		}
		records = append(records, record)
	}
	if !fileutil.IsExist(fullPath) {
		dir := filepath.Dir(fullPath)
		if !fileutil.IsExist(dir) {
			err := fileutil.CreateDir(dir)
			if err != nil {
				panic(err)
			}
		}
		status := fileutil.CreateFile(fullPath)
		if !status {
			panic(errors.New("create file error"))
		}
	}
	err := fileutil.WriteCsvFile(fullPath, records, false)
	return err
}

func StructToJsonFlatLineFile[T any](datas []T, fullPath string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("StructToJsonFlatLineFile panic: " + err.(error).Error())
		}
	}()
	if datas == nil || len(datas) == 0 {
		slog.Info("datas is null")
		return nil
	}
	jsons := make([]string, 0)
	for _, data := range datas {
		strline, err := json.Marshal(data)
		if err != nil {
			return err
		}
		jsons = append(jsons, string(strline))
	}
	res := strings.Join(jsons, "")
	if !fileutil.IsExist(fullPath) {
		if !fileutil.IsExist(filepath.Dir(fullPath)) {
			err := fileutil.CreateDir(filepath.Dir(fullPath))
			if err != nil {
				return err
			}
		}
		status := fileutil.CreateFile(fullPath)
		if !status {
			return errors.New("create file error")
		}
	}
	return fileutil.WriteStringToFile(fullPath, res, false)
}

func StructToJsonFile[T any](datas []T, fullPath string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("StructToJsonFlatLineFile panic: " + err.(error).Error())
		}
	}()
	if datas == nil || len(datas) == 0 {
		slog.Info("StructToJsonFile ==> datas is null")
		return nil
	}
	strline, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	if !fileutil.IsExist(fullPath) {
		if !fileutil.IsExist(filepath.Dir(fullPath)) {
			err = fileutil.CreateDir(filepath.Dir(fullPath))
			if err != nil {
				return err
			}
		}
		status := fileutil.CreateFile(fullPath)
		if !status {
			return errors.New("create file error")
		}
	}
	return fileutil.WriteStringToFile(fullPath, string(strline), false)
}
