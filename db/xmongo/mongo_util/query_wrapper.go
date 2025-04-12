package mongo_util

import (
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DeletedField = "delete_at"
var BaseFilter = func() wrapper.QueryWrapper {
	return wrapper.QueryWrapper{
		DeletedField: wrapper.QueryWrapper{"$exists": false}, // Field doesn't exist
	}
}()

func BuildQueryWrapper(queryWrapperMap map[string]any) wrapper.QueryWrapper {
	m := wrapper.QueryWrapper{}
	queryWrapperMap = maputil.Merge(queryWrapperMap, BaseFilter)
	for k, v := range queryWrapperMap {
		m[k] = v
	}
	return m
}

func ObjectIDFromHex(id string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return hex
}
