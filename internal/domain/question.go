package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var (
	ErrQuestionNotFound = NewError("ERR_QUESTION_NOT_FOUND", "question not found")
)

const (
	QuestionCollectionName = "questions"
)

type Question struct {
	ID          bson.ObjectID `bson:"_id" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Attachments []string      `bson:"attachments" json:"attachments"`
	Points      uint          `bson:"points" json:"points"`
	Location    Location      `bson:"location" json:"location"`
	UserID      bson.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	// TODO: tags
	// TODO: likes/rating for questions and answers
}
