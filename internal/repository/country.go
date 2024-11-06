package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CountryRepository interface {
	Create(ctx context.Context, country domain.Country) error
	Get(ctx context.Context) ([]domain.Country, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Country, error)
}

type countryRepository struct {
	*Repository
}

func NewCountryRepository(repository *Repository) CountryRepository {
	return &countryRepository{
		Repository: repository,
	}
}

func (r *countryRepository) Create(ctx context.Context, country domain.Country) error {
	_, err := r.db.Collection(domain.CountryCollectionName).
		InsertOne(ctx, country)

	return err
}

func (r *countryRepository) Get(ctx context.Context) ([]domain.Country, error) {
	cursor, err := r.db.Collection(domain.CountryCollectionName).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var countries []domain.Country

	if err := cursor.All(ctx, &countries); err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.Country, error) {
	var country domain.Country

	err := r.db.Collection(domain.CountryCollectionName).
		FindOne(ctx, bson.M{"_id": id}).Decode(&country)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Country{}, domain.ErrCountryNotFound
		}
		return domain.Country{}, err
	}

	return country, nil
}
