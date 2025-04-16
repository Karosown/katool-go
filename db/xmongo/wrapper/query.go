package wrapper

import (
	"encoding/json"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/convert"
	"go.mongodb.org/mongo-driver/bson"
)

type QueryWrapper bson.M

func (q *QueryWrapper) ToJSON() string {
	marshal, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return ""
	}
	return string(marshal)
}

type Query struct {
	query QueryWrapper
}

func NewQuery() *Query {
	return &Query{
		query: QueryWrapper{},
	}
}

func (q *Query) Eq(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$eq": value,
	})
	return q
}

func (q *Query) Ne(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$ne": value,
	})
	return q
}

func (q *Query) Gt(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$gt": value,
	})
	return q
}

func (q *Query) Gte(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$gte": value,
	})
	return q
}

func (q *Query) validWrapper(column string) map[string]any {
	if w, ok := q.query[column].(map[string]any); ok {
		return w
	}
	return QueryWrapper{}
}

func (q *Query) Lt(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$lt": value,
	})
	return q
}

func (q *Query) Lte(column string, value any) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$lte": value,
	})
	return q
}

func (q *Query) And(queries ...*Query) *Query {
	allQueries := append([]*Query{q}, queries...) // Corrected the order
	newQuery := NewQuery()
	newQuery.query["$and"] = convert.FromAnySlice[QueryWrapper](stream.ToStream(&allQueries).Map(func(item *Query) any {
		return item.Origin()
	}).ToList())
	return newQuery
}

func (q *Query) Or(queries ...*Query) *Query {
	allQueries := append([]*Query{q}, queries...) // Corrected the order
	newQuery := NewQuery()
	newQuery.query["$or"] = convert.FromAnySlice[QueryWrapper](stream.ToStream(&allQueries).Map(func(item *Query) any {
		return item.Origin()
	}).ToList())
	return newQuery
}

func (q *Query) Exists(column string, exists bool) *Query {
	wrapper := q.validWrapper(column)
	q.query[column] = maputil.Merge(wrapper, QueryWrapper{
		"$exists": exists,
	})
	return q
}

func (q *Query) Build(deletedField ...string) QueryWrapper {
	return BuildQueryWrapper(q.query, deletedField...)
}

func (q *Query) Origin() QueryWrapper {
	return q.query
}

var DeletedField = "delete_at"

var BaseFilter = func(deletedField ...string) QueryWrapper {
	wrapper := QueryWrapper{}
	if cutil.IsEmpty(deletedField) {
		deletedField = append(deletedField, DeletedField)
	}

	for _, field := range deletedField {
		wrapper[field] = QueryWrapper{"$exists": false}
	}
	return wrapper
}

func BuildQueryWrapper(queryWrapperMap QueryWrapper, deletedField ...string) QueryWrapper {
	m := QueryWrapper{}
	queryWrapperMap = maputil.Merge(queryWrapperMap, BaseFilter(deletedField...))
	for k, v := range queryWrapperMap {
		m[k] = v
	}
	return m
}
