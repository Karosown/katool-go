package mongo_util

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FormatPrimitiveTimeOfTimeZone(pt primitive.DateTime, zone string) time.Time {
	return FormTimeOfTimeZone(pt.Time(), zone)
}

func FormTimeOfTimeZone(pt time.Time, zone string) time.Time {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return pt
	}
	return pt.In(loc)
}
