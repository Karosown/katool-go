package wrapper

import (
	"encoding/json"

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

// addOperatorCondition 统一处理操作符条件添加
func (q *Query) addOperatorCondition(column string, operator string, value interface{}) *Query {
	// MongoDB 原生支持点记法，无需特殊处理嵌套字段
	wrapper := q.getOperatorWrapper(column, operator)
	wrapper[operator] = value
	return q
}

// getOperatorWrapper 安全地获取或创建字段的操作符包装器
func (q *Query) getOperatorWrapper(column string, operator string) bson.M {
	// 检查字段是否已存在
	if val, exists := q.query[column]; exists {
		switch v := val.(type) {
		case bson.M:
			// 检查是否已存在相同操作符
			if _, opExists := v[operator]; opExists && needsMultipleValues(operator) {
				return q.handleDuplicateOperator(column, operator, v)
			}
			return v
		case map[string]interface{}:
			wrapper := bson.M(v)
			if _, opExists := wrapper[operator]; opExists && needsMultipleValues(operator) {
				return q.handleDuplicateOperator(column, operator, wrapper)
			}
			return wrapper
		case QueryWrapper:
			wrapper := bson.M(v)
			if _, opExists := wrapper[operator]; opExists && needsMultipleValues(operator) {
				return q.handleDuplicateOperator(column, operator, wrapper)
			}
			return wrapper
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

// needsMultipleValues 判断操作符是否需要处理多值情况
func needsMultipleValues(operator string) bool {
	switch operator {
	case "$ne", "$gt", "$gte", "$lt", "$lte":
		return true
	default:
		return false
	}
}

// handleDuplicateOperator 处理重复操作符，转换为$and结构
func (q *Query) handleDuplicateOperator(column, operator string, existingWrapper bson.M) bson.M {
	// 创建$and条件来容纳多个相同操作符
	existingCondition := QueryWrapper{column: existingWrapper}

	// 初始化$and数组
	var andConditions []QueryWrapper
	if q.query["$and"] != nil {
		andConditions = q.query["$and"].([]QueryWrapper)
	}

	// 添加现有条件到$and中
	andConditions = append(andConditions, existingCondition)
	q.query["$and"] = andConditions

	// 从原字段中移除
	delete(q.query, column)

	// 返回新的wrapper用于当前操作
	return bson.M{}
}

// SortDirection defines the sort direction type
type SortDirection int

const (
	ASC  SortDirection = 1  // 升序
	DESC SortDirection = -1 // 降序
)

// Sort contains sorting information
type Sort struct {
	Field     string
	Direction SortDirection
}

// OrderBy sets the sorting order for query results
func (q *Query) OrderBy(column string, direction SortDirection) *Query {
	if q.query["$sort"] == nil {
		q.query["$sort"] = []Sort{}
	}

	sorts := q.query["$sort"].([]Sort)
	sorts = append(sorts, Sort{Field: column, Direction: direction})
	q.query["$sort"] = sorts
	return q
}

// ClearOrder removes all sorting criteria
func (q *Query) ClearOrder() *Query {
	delete(q.query, "$sort")
	return q
}

// GetSortBson returns the sort specification in BSON format
func (q *Query) GetSortBson() bson.D {
	if q.query["$sort"] == nil {
		return bson.D{}
	}

	sorts := q.query["$sort"].([]Sort)
	sortDoc := bson.D{}

	for _, sort := range sorts {
		sortDoc = append(sortDoc, bson.E{Key: sort.Field, Value: sort.Direction})
	}

	return sortDoc
}

// 简化的查询方法 - 使用统一的处理逻辑

// Eq 等于条件
func (q *Query) Eq(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$eq", value)
}

// Ne 不等于条件
func (q *Query) Ne(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$ne", value)
}

// Gt 大于条件
func (q *Query) Gt(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$gt", value)
}

// Gte 大于等于条件
func (q *Query) Gte(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$gte", value)
}

// Lt 小于条件
func (q *Query) Lt(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$lt", value)
}

// Lte 小于等于条件
func (q *Query) Lte(column string, value interface{}) *Query {
	return q.addOperatorCondition(column, "$lte", value)
}

// Exists 字段存在性条件
func (q *Query) Exists(column string, exists bool) *Query {
	return q.addOperatorCondition(column, "$exists", exists)
}

// In 包含条件
func (q *Query) In(column string, values ...interface{}) *Query {
	return q.addOperatorCondition(column, "$in", values)
}

// Nin 不包含条件
func (q *Query) Nin(column string, values ...interface{}) *Query {
	return q.addOperatorCondition(column, "$nin", values)
}

// Regex 正则表达式匹配
func (q *Query) Regex(column string, pattern string, options ...string) *Query {
	regexObj := bson.M{"$regex": pattern}
	if len(options) > 0 && options[0] != "" {
		regexObj["$options"] = options[0]
	}
	q.query[column] = regexObj
	return q
}

// And 逻辑与操作
func (q *Query) And(queries ...*Query) *Query {
	var conditions []QueryWrapper

	// 添加现有的$and条件
	if q.query["$and"] != nil {
		conditions = q.query["$and"].([]QueryWrapper)
	}

	// 添加新的查询条件
	for _, query := range queries {
		conditions = append(conditions, query.query)
	}

	q.query["$and"] = conditions
	return q
}

// Or 逻辑或操作
func (q *Query) Or(queries ...*Query) *Query {
	var conditions []QueryWrapper

	// 添加现有的$or条件
	if q.query["$or"] != nil {
		conditions = q.query["$or"].([]QueryWrapper)
	}

	// 添加新的查询条件
	for _, query := range queries {
		conditions = append(conditions, query.query)
	}

	q.query["$or"] = conditions
	return q
}

// Build 构建最终查询条件，包含软删除过滤
func (q *Query) Build(deletedFields ...string) QueryWrapper {
	return BuildQueryWrapper(q.query, deletedFields...)
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
		wrapper[field] = bson.M{"$or": []bson.M{
			{"$exists": false},
			{"$eq": nil},
		}}
	}

	return wrapper
}

// BuildQueryWrapper 构建查询包装器，合并基础过滤条件
func BuildQueryWrapper(queryWrapperMap QueryWrapper, deletedFields ...string) QueryWrapper {
	filter := BaseFilter(deletedFields...)

	if queryWrapperMap["$and"] == nil {
		return maputil.Merge(queryWrapperMap, filter)
	}

	// 如果已有$and条件，添加到现有数组中
	wrappers := queryWrapperMap["$and"].([]QueryWrapper)
	for k, v := range filter {
		wrappers = append(wrappers, QueryWrapper{k: v})
	}
	queryWrapperMap["$and"] = wrappers
	return queryWrapperMap
}
