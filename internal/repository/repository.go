package repository

import (
	"github.com/Closi-App/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	log *logger.Logger
	db  *mongo.Database
}

func NewRepository(log *logger.Logger, db *mongo.Database) *Repository {
	return &Repository{
		log: log,
		db:  db,
	}
}
