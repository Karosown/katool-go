package wrapper

import (
	"encoding/json"
	"github.com/karosown/katool-go/db/xmongo/mongo_util"
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
	q.query[clomn] = bson.M{
		"$eq": value,
	}
	return q
}

func (q *Query) Ne(clomn string, value any) *Query {
	q.query[clomn] = bson.M{
		"$ne": value,
	}
	return q
}
func (q *Query) Gt(clomn string, value any) *Query {
	q.query[clomn] = bson.M{
		"$gt": value,
	}
	return q
}
func (q *Query) Gte(clomn string, value any) *Query {
	q.query[clomn] = bson.M{
		"$gte": value,
	}
	return q
}
func (q *Query) Lt(clomn string, value any) *Query {
	q.query[clomn] = bson.M{
		"$lt": value,
	}
	return q
}
func (q *Query) Lte(clomn string, value any) *Query {
	q.query[clomn] = bson.M{
		"$lte": value,
	}
	return q
}
func (q *Query) And(query ...*Query) *Query {
	q.query["$and"] = query
	return q
}
func (q *Query) Or(query ...*Query) *Query {
	q.query["$or"] = query
	return q
}
func (q *Query) Build() QueryWrapper {
	return mongo_util.BuildQueryWrapper(q.query)
}
func (q *Query) Origin() QueryWrapper {
	return q.query
}
