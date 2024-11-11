package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var (
	ErrAnswerNotFound = NewError("ERR_ANSWER_NOT_FOUND", "answer not found")
)

const (
	AnswerCollectionName = "answers"
)

type Answer struct {
	ID         bson.ObjectID `bson:"_id" json:"id"`
	Text       string        `bson:"text" json:"text"`
	Likes      uint          `bson:"likes" json:"likes"`
	QuestionID bson.ObjectID `bson:"question_id" json:"question_id"`
	UserID     bson.ObjectID `bson:"user_id" json:"user_id"`
	IsVerified bool          `bson:"is_verified" json:"is_verified"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

type AnswerGetAllFilter struct {
	QuestionID *bson.ObjectID
	UserID     *bson.ObjectID
}

type AnswerUpdateInput struct {
	Text *string
}
