package wrapper

import (
	"encoding/json"
	"strings"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool-go/container/cutil"
	"go.mongodb.org/mongo-driver/bson"
)

// QueryWrapper 封装bson.M以提供更好的类型支持
type QueryWrapper bson.M

// ToJSON 将查询条件转换为格式化的JSON字符串
func (q QueryWrapper) ToJSON() string {
	marshal, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return ""
	}
	return string(marshal)
}

// Query 查询构建器
type Query struct {
	query QueryWrapper
}

// NewQuery 创建新的查询构建器
func NewQuery() *Query {
	return &Query{
		query: QueryWrapper{},
	}
}

// getOperatorWrapper 安全地获取或创建字段的操作符包装器
func (q *Query) getOperatorWrapper(column string) bson.M {
	// 处理嵌套字段
	if strings.Contains(column, ".") {
		return q.createNestedFieldQuery(column, nil, "")
	}

	// 检查字段是否已存在
	if val, exists := q.query[column]; exists {
		switch v := val.(type) {
		case bson.M:
			return v
		case map[string]interface{}:
			return bson.M(v)
		case QueryWrapper:
			return bson.M(v)
		default:
			// 字段存在但不是操作符结构，需要重置
			delete(q.query, column)
		}
	}

	// 创建新的操作符包装器
	wrapper := bson.M{}
	q.query[column] = wrapper
	return wrapper
}

// createNestedFieldQuery 创建嵌套字段查询
func (q *Query) createNestedFieldQuery(path string, value interface{}, operator string) bson.M {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		if operator == "" {
			return bson.M{path: value}
		}
		return bson.M{path: bson.M{operator: value}}
	}

	field := parts[0]
	remainingPath := strings.Join(parts[1:], ".")

	var nestedQuery bson.M
	if operator == "" {
		nestedQuery = q.createNestedFieldQuery(remainingPath, value, operator)
	} else {
		nestedQuery = q.createNestedFieldQuery(remainingPath, value, operator)
	}

	return bson.M{field: nestedQuery}
}

// Eq 等于条件
func (q *Query) Eq(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$eq")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$eq"] = value
	q.query[column] = wrapper
	return q
}

// Ne 不等于条件
func (q *Query) Ne(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$ne")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$ne"] = value
	q.query[column] = wrapper
	return q
}

// Gt 大于条件
func (q *Query) Gt(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$gt")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$gt"] = value
	q.query[column] = wrapper
	return q
}

// Gte 大于等于条件
func (q *Query) Gte(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$gte")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$gte"] = value
	q.query[column] = wrapper
	return q
}

// Lt 小于条件
func (q *Query) Lt(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$lt")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$lt"] = value
	q.query[column] = wrapper
	return q
}

// Lte 小于等于条件
func (q *Query) Lte(column string, value any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, value, "$lte")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$lte"] = value
	q.query[column] = wrapper
	return q
}

// Exists 字段存在性条件
func (q *Query) Exists(column string, exists bool) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, exists, "$exists")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$exists"] = exists
	q.query[column] = wrapper
	return q
}

// In 包含条件
func (q *Query) In(column string, values ...any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, values, "$in")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$in"] = values
	q.query[column] = wrapper
	return q
}

// Nin 不包含条件
func (q *Query) Nin(column string, values ...any) *Query {
	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, values, "$nin")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	wrapper := q.getOperatorWrapper(column)
	wrapper["$nin"] = values
	q.query[column] = wrapper
	return q
}

// Regex 正则表达式匹配
func (q *Query) Regex(column string, pattern string, options string) *Query {
	regexObj := bson.M{"$regex": pattern}
	if options != "" {
		regexObj["$options"] = options
	}

	if strings.Contains(column, ".") {
		nestedQuery := q.createNestedFieldQuery(column, regexObj, "")
		q.query = maputil.Merge(q.query, QueryWrapper(nestedQuery))
		return q
	}

	q.query[column] = regexObj
	return q
}

// And 逻辑与操作
func (q *Query) And(queries ...*Query) *Query {
	qar := []QueryWrapper{}
	// 添加其他查询条件
	if q.query["$and"] != nil {
		qar = append(qar, q.query["$and"].([]QueryWrapper)...)
	}
	for _, query := range queries {
		qar = append(qar, query.query)
	}
	q.query["$and"] = qar
	return q
}

// Or 逻辑或操作
func (q *Query) Or(queries ...*Query) *Query {
	qar := []QueryWrapper{}
	// 添加其他查询条件
	if q.query["$or"] != nil {
		qar = append(qar, q.query["$or"].([]QueryWrapper)...)
	}
	for _, query := range queries {
		qar = append(qar, query.query)
	}
	q.query["$or"] = qar
	return q
}

// Build 构建最终查询条件，包含软删除过滤
func (q *Query) Build(deletedFields ...string) bson.M {
	return bson.M(BuildQueryWrapper(q.query, deletedFields...))
}

// Origin 获取原始查询条件
func (q *Query) Origin() QueryWrapper {
	return q.query
}

// DeletedField 软删除字段名称
var DeletedField = "deleted_at"

// BaseFilter 创建基础过滤条件（排除已删除文档）
func BaseFilter(deletedFields ...string) QueryWrapper {
	wrapper := QueryWrapper{}

	if cutil.IsEmpty(deletedFields) {
		deletedFields = []string{DeletedField}
	}

	for _, field := range deletedFields {
		wrapper[field] = bson.M{"$exists": false}
	}

	return wrapper
}

// BuildQueryWrapper 构建查询包装器，合并基础过滤条件
func BuildQueryWrapper(queryWrapperMap QueryWrapper, deletedFields ...string) QueryWrapper {
	// 直接合并用户查询和软删除过滤器
	return maputil.Merge(queryWrapperMap, BaseFilter(deletedFields...))
}
