package mongo_util

import (
	"github.com/duke-git/lancet/v2/maputil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BaseFilter = bson.M{
	"delete_at": bson.M{"$exists": false}, // Field doesn't exist
}

func BuildQueryWrapper(queryWrapperMap map[string]any) bson.M {
	m := bson.M{}
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
