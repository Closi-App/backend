package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer domain.Answer) error
	GetAll(ctx context.Context, filter ...domain.AnswerGetAllFilter) ([]domain.Answer, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Answer, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input domain.AnswerUpdateInput) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error

	AddLike(ctx context.Context, id bson.ObjectID) error
	RemoveLike(ctx context.Context, id bson.ObjectID) error
	Verify(ctx context.Context, id bson.ObjectID) error
}

type answerRepository struct {
	*Repository
}

func NewAnswerRepository(repository *Repository) AnswerRepository {
	return &answerRepository{
		Repository: repository,
	}
}

func (r *answerRepository) Create(ctx context.Context, answer domain.Answer) error {
	_, err := r.db.Collection(domain.AnswerCollectionName).
		InsertOne(ctx, answer)

	return err
}

func (r *answerRepository) GetAll(ctx context.Context, filter ...domain.AnswerGetAllFilter) ([]domain.Answer, error) {
	filterFields := bson.M{}

	if len(filter) > 0 {
		f := filter[0]

		if f.QuestionID != nil {
			filterFields["question_id"] = f.QuestionID
		}
		if f.UserID != nil {
			filterFields["user_id"] = f.UserID
		}
	}

	cursor, err := r.db.Collection(domain.AnswerCollectionName).
		Find(ctx, filterFields)
	if err != nil {
		return nil, err
	}

	var answers []domain.Answer

	if err = cursor.All(ctx, &answers); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *answerRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.Answer, error) {
	var answer domain.Answer

	err := r.db.Collection(domain.AnswerCollectionName).
		FindOne(ctx, bson.M{"_id": id}).Decode(&answer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Answer{}, domain.ErrAnswerNotFound
		}

		return domain.Answer{}, err
	}

	return answer, nil
}

func (r *answerRepository) Update(ctx context.Context, id, userID bson.ObjectID, input domain.AnswerUpdateInput) error {
	updateFields := bson.M{}

	if input.Text != nil {
		updateFields["text"] = input.Text
	}

	updateFields["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.AnswerCollectionName).
		UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": updateFields})

	return err
}

func (r *answerRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := r.db.Collection(domain.AnswerCollectionName).
		DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})

	return err
}

func (r *answerRepository) AddLike(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.AnswerCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"likes": 1}})

	return err
}

func (r *answerRepository) RemoveLike(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.AnswerCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"likes": -1}})

	return err
}

func (r *answerRepository) Verify(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.AnswerCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_verified": true}})

	return err
}
