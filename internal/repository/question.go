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
	GetAll(ctx context.Context, filter ...domain.QuestionGetAllFilter) ([]domain.Question, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input domain.QuestionUpdateInput) error
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

	return err
}

func (r *questionRepository) GetAll(ctx context.Context, filter ...domain.QuestionGetAllFilter) ([]domain.Question, error) {
	filterFields := bson.M{}

	if len(filter) > 0 {
		f := filter[0]

		if f.Title != nil {
			filterFields["title"] = f.Title
		}
		if f.Tag != nil {
			filterFields["tags"] = bson.M{"$in": f.Tag}
		}
		if f.CountryID != nil {
			filterFields["country_id"] = f.CountryID
		}
		if f.UserID != nil {
			filterFields["user_id"] = f.UserID
		}
	}

	cursor, err := r.db.Collection(domain.QuestionCollectionName).
		Find(ctx, filterFields)
	if err != nil {
		return nil, err
	}

	var questions []domain.Question

	if err = cursor.All(ctx, &questions); err != nil {
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

func (r *questionRepository) Update(ctx context.Context, id, userID bson.ObjectID, input domain.QuestionUpdateInput) error {
	updateFields := bson.M{}

	if input.Title != nil {
		updateFields["title"] = input.Title
	}
	if input.Description != nil {
		updateFields["description"] = input.Description
	}
	if input.AttachmentsURL != nil {
		updateFields["attachments_url"] = input.AttachmentsURL
	}
	if input.Tags != nil {
		updateFields["tags"] = input.Tags
	}
	if input.Points != nil {
		updateFields["points"] = input.Points
	}

	updateFields["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.QuestionCollectionName).
		UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": updateFields})

	return err
}

func (r *questionRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := r.db.Collection(domain.QuestionCollectionName).
		DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})

	return err
}
