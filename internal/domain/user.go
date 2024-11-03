package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var (
	ErrUserAlreadyExists = NewError("ERR_USER_ALREADY_EXISTS", "user already exists")
	ErrUserNotFound      = NewError("ERR_USER_NOT_FOUND", "user not found")
)

const (
	DefaultUserPoints = 10
)

type User struct {
	ID                      bson.ObjectID           `bson:"_id" json:"id"`
	Name                    string                  `bson:"name" json:"name"`
	Username                string                  `bson:"username" json:"username"`
	Email                   string                  `bson:"email" json:"email"`
	Password                string                  `bson:"password" json:"password"`
	AvatarURL               string                  `bson:"avatar_url" json:"avatar_url"`
	Points                  uint                    `bson:"points" json:"points"`
	Favorites               []bson.ObjectID         `bson:"favorites" json:"favorites"`
	Subscription            Subscription            `bson:"subscription" json:"subscription"`
	NotificationPreferences NotificationPreferences `bson:"notification_preferences" json:"notification_preferences"`
	CreatedAt               time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time               `bson:"updated_at" json:"updated_at"`
}

type NotificationPreferences struct {
	Email bool `bson:"email" json:"email"`
	Push  bool `bson:"push" json:"push"`
}
