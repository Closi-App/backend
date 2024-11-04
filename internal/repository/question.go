package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type QuestionRepository interface {
	Create(ctx context.Context, question domain.Question) error
	GetAll(ctx context.Context) ([]domain.Question, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error)
	GetByLocation(ctx context.Context, location domain.Location) ([]domain.Question, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input QuestionUpdateInput) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
}

type questionRepository struct {
	*Repository
}

func NewQuestionRepository(repository *Repository) QuestionRepository {
	return &questionRepository{
		Repository: repository,
	}
}

func (r *questionRepository) Create(ctx context.Context, question domain.Question) error {
	_, err := r.db.Collection(domain.QuestionCollectionName).
		InsertOne(ctx, question)
	if err != nil {
		return err
	}

	return nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]domain.Question, error) {
	cursor, err := r.db.Collection(domain.QuestionCollectionName).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var questions []domain.Question

	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *questionRepository) GetByLocation(ctx context.Context, location domain.Location) ([]domain.Question, error) {
	cursor, err := r.db.Collection(domain.QuestionCollectionName).
		Find(ctx, bson.M{"location.country": location.Country})
	if err != nil {
		return nil, err
	}

	var questions []domain.Question

	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *questionRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error) {
	var question domain.Question

	err := r.db.Collection(domain.QuestionCollectionName).
		FindOne(ctx, bson.M{"_id": id}).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Question{}, domain.ErrQuestionNotFound
		}
		return domain.Question{}, err
	}

	return question, nil
}

type QuestionUpdateInput struct {
	Title       string
	Description string
	Attachments []string
	Points      uint
	Location    *domain.Location
}

func (r *questionRepository) Update(ctx context.Context, id, userID bson.ObjectID, input QuestionUpdateInput) error {
	updateQuery := bson.M{}

	if input.Title != "" {
		updateQuery["title"] = input.Title
	}
	if input.Description != "" {
		updateQuery["description"] = input.Description
	}
	if input.Attachments != nil {
		var attachments []interface{}

		for _, a := range input.Attachments {
			attachments = append(attachments, a)
		}

		updateQuery["attachments"] = attachments
	}
	if input.Points != 0 {
		updateQuery["points"] = input.Points
	}
	if input.Location != nil {
		updateQuery["location"] = input.Location
	}

	updateQuery["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.QuestionCollectionName).
		UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": updateQuery})

	return err
}

func (r *questionRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := r.db.Collection(domain.QuestionCollectionName).
		DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})

	return err
}
