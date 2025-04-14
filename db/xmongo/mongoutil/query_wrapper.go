package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIDFromHex(id string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return hex
}
