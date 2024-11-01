package mongo

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongo(cfg *viper.Viper) *mongo.Database {
	opts := options.Client().ApplyURI(cfg.GetString("mongo.uri"))

	client, err := mongo.Connect(opts)
	if err != nil {
		panic("error connecting to mongo database")
	}

	return client.Database(cfg.GetString("mongo.database"))
}
