package repository

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	log logger.Logger
	db  *mongo.Database
}

func New(log logger.Logger, db *mongo.Database) *Repository {
	return &Repository{
		log: log,
		db:  db,
	}
}

func NewDB(cfg *config.Config, log logger.Logger) *mongo.Database {
	opts := options.Client().ApplyURI(cfg.Mongo.URI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("error connecting to mongo db")
	}

	return client.Database(cfg.Mongo.DBName)
}
