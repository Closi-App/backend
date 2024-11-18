package domain

import "go.mongodb.org/mongo-driver/v2/bson"

var (
	ErrCountryNotFound = NewError("ERR_COUNTRY_NOT_FOUND", "country not found")
)

const (
	CountryCollectionName = "countries"
)

type Country struct {
	ID    bson.ObjectID     `bson:"_id" json:"id"`
	Names map[string]string `bson:"name" json:"name"`
}
