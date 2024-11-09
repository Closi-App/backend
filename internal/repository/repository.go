package repository

import (
	"github.com/Closi-App/backend/pkg/logger"
	imgbb "github.com/JohnNON/ImgBB"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	log   *logger.Logger
	db    *mongo.Database
	rdb   *redis.Client
	imgbb *imgbb.Client
}

func NewRepository(
	log *logger.Logger,
	db *mongo.Database,
	rdb *redis.Client,
	imgbb *imgbb.Client,
) *Repository {
	return &Repository{
		log:   log,
		db:    db,
		rdb:   rdb,
		imgbb: imgbb,
	}
}
