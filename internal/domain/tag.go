package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrTagNotFound = NewError("ERR_TAG_NOT_FOUND", "tag not found")
)

const (
	TagCollectionName = "tags"
)

type Tag struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	CountryID bson.ObjectID `bson:"country_id" json:"country_id"`
}
