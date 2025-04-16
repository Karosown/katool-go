package wrapper

import (
	"encoding/json"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/stream"
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
func (q *Query) Eq(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$eq": value,
	})
	return q
}

func (q *Query) Ne(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$ne": value,
	})
	return q
}
func (q *Query) Gt(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$gt": value,
	})
	return q
}
func (q *Query) Gte(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$gte": value,
	})
	return q
}
func (q *Query) validWrapper(clomn string) map[string]any {
	switch q.query[clomn].(type) {
	case QueryWrapper:
		return q.query[clomn].(QueryWrapper)
	case map[string]interface{}:
		return q.query[clomn].(map[string]interface{})
	case bson.M:
		return q.query[clomn].(bson.M)
	default:
		return QueryWrapper{}
	}
}
func (q *Query) Lt(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$lt": value,
	})
	return q
}
func (q *Query) Lte(clomn string, value any) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$lte": value,
	})
	return q
}
func (q *Query) And(query ...*Query) *Query {
	query = append(query, q)
	q.query["$and"] = func(query ...*Query) []any {
		return stream.ToStream(&query).Map(func(item *Query) any {
			return item.Origin()
		}).ToList()
	}(query...)
	for s, _ := range q.query {
		delete(q.query, s)
	}
	return q
}
func (q *Query) Or(query ...*Query) *Query {
	query = append(query, q)
	q.query["$or"] = func(query ...*Query) []any {
		return stream.ToStream(&query).Map(func(item *Query) any {
			return item.Origin()
		}).ToList()
	}(query...)
	for s, _ := range q.query {
		delete(q.query, s)
	}
	return q
}
func (q *Query) Exists(clomn string, est bool) *Query {
	wrapper := q.validWrapper(clomn)
	q.query[clomn] = maputil.Merge(wrapper, QueryWrapper{
		"$exists": est,
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
