package repository

import (
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	log *logger.Logger
	db  *mongo.Database
	rdb *redis.Client
}

func NewRepository(log *logger.Logger, db *mongo.Database, rdb *redis.Client) *Repository {
	return &Repository{
		log: log,
		db:  db,
		rdb: rdb,
	}
}
