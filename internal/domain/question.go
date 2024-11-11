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
	ID             bson.ObjectID   `bson:"_id" json:"id"`
	Title          string          `bson:"title" json:"title"`
	Description    string          `bson:"description" json:"description"`
	AttachmentsURL []string        `bson:"attachments_url" json:"attachments_url"`
	Tags           []bson.ObjectID `bson:"tags" json:"tags"`
	Points         uint            `bson:"points" json:"points"`
	CountryID      bson.ObjectID   `bson:"country_id" json:"country_id"`
	UserID         bson.ObjectID   `bson:"user_id" json:"user_id"`
	CreatedAt      time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `bson:"updated_at" json:"updated_at"`
}

type QuestionGetAllFilter struct {
	Title     *string
	Tag       *bson.ObjectID
	CountryID *bson.ObjectID
	UserID    *bson.ObjectID
}

type QuestionUpdateInput struct {
	Title          *string
	Description    *string
	AttachmentsURL []string
	Tags           []bson.ObjectID
	Points         *uint
}
