package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TagRepository interface {
	Create(ctx context.Context, tag domain.Tag) error
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Tag, error)
	GetAll(ctx context.Context) ([]domain.Tag, error)
	GetAllByCountryID(ctx context.Context, countryID bson.ObjectID) ([]domain.Tag, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

type tagRepository struct {
	*Repository
}

func NewTagRepository(repository *Repository) TagRepository {
	return &tagRepository{
		Repository: repository,
	}
}

func (r *tagRepository) Create(ctx context.Context, tag domain.Tag) error {
	_, err := r.db.Collection(domain.TagCollectionName).
		InsertOne(ctx, tag)

	return err
}

func (r *tagRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.Tag, error) {
	var tag domain.Tag

	err := r.db.Collection(domain.TagCollectionName).
		FindOne(ctx, bson.M{"_id": id}).Decode(&tag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Tag{}, domain.ErrTagNotFound
		}

		return domain.Tag{}, err
	}

	return tag, nil
}

func (r *tagRepository) GetAll(ctx context.Context) ([]domain.Tag, error) {
	cursor, err := r.db.Collection(domain.TagCollectionName).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var tags []domain.Tag

	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) GetAllByCountryID(ctx context.Context, countryID bson.ObjectID) ([]domain.Tag, error) {
	cursor, err := r.db.Collection(domain.TagCollectionName).
		Find(ctx, bson.M{"country_id": countryID})
	if err != nil {
		return nil, err
	}

	var tags []domain.Tag

	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.TagCollectionName).
		DeleteOne(ctx, bson.M{"_id": id})

	return err
}
