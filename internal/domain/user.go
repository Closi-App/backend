package domain

import (
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"math/rand"
	"time"
)

var (
	ErrUserAlreadyExists      = NewError("ERR_USER_ALREADY_EXISTS", "user already exists")
	ErrUserNotFound           = NewError("ERR_USER_NOT_FOUND", "user not found")
	ErrUserInsufficientPoints = NewError("ERR_USER_INSUFFICIENT_POINTS", "insufficient points")
)

const (
	UserCollectionName = "users"

	UserDefaultPoints  = 10
	UserReferralPoints = 50

	userReferralCodeLength = 4
)

type User struct {
	ID           bson.ObjectID   `bson:"_id" json:"id"`
	Name         string          `bson:"name" json:"name"`
	Username     string          `bson:"username" json:"username"`
	Email        string          `bson:"email" json:"email"`
	Password     string          `bson:"password" json:"password"`
	AvatarURL    string          `bson:"avatar_url" json:"avatar_url"`
	Points       int             `bson:"points" json:"points"`
	Favorites    []bson.ObjectID `bson:"favorites" json:"favorites"`
	ReferralCode string          `bson:"referral_code" json:"referral_code"`
	Subscription Subscription    `bson:"subscription" json:"subscription"`
	Settings     UserSettings    `bson:"settings" json:"settings"`
	IsConfirmed  bool            `bson:"is_confirmed" json:"is_confirmed"`
	CreatedAt    time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `bson:"updated_at" json:"updated_at"`
	// TODO: otp
}

type UserSettings struct {
	Location           Location `bson:"location" json:"location"`
	Language           Language `bson:"language" json:"language"`
	EmailNotifications bool     `bson:"email_notifications" json:"email_notifications"`
	// TODO: changing appearance
}

func NewReferralCode() (string, error) {
	b := make([]byte, userReferralCodeLength)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
