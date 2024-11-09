package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// TODO: admins

var (
	ErrUserAlreadyExists      = NewError("ERR_USER_ALREADY_EXISTS", "user already exists")
	ErrUserNotFound           = NewError("ERR_USER_NOT_FOUND", "user not found")
	ErrUserInsufficientPoints = NewError("ERR_USER_INSUFFICIENT_POINTS", "insufficient points")
)

const (
	UserCollectionName = "users"

	UserDefaultPoints  = 10
	UserReferralPoints = 50

	UserReferralCodeLength = 8
)

type User struct {
	ID           bson.ObjectID   `bson:"_id" json:"id"`
	Name         string          `bson:"name" json:"name"`
	Username     string          `bson:"username" json:"username"`
	Email        string          `bson:"email" json:"email"`
	Password     string          `bson:"password" json:"password"`
	AvatarURL    string          `bson:"avatar_url" json:"avatar_url"`
	Points       uint            `bson:"points" json:"points"`
	Favorites    []bson.ObjectID `bson:"favorites" json:"favorites"`
	Achievements []bson.ObjectID `bson:"achievements" json:"achievements"`
	SocialLinks  []string        `bson:"social_links" json:"social_links"`
	ReferralCode string          `bson:"referral_code" json:"referral_code"`
	Subscription Subscription    `bson:"subscription" json:"subscription"`
	Settings     UserSettings    `bson:"settings" json:"settings"`
	IsConfirmed  bool            `bson:"is_confirmed" json:"is_confirmed"`
	IsBlocked    bool            `bson:"is_blocked" json:"is_blocked"`
	CreatedAt    time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `bson:"updated_at" json:"updated_at"`
	// TODO: promotions, promo codes for subscription
}

type UserSettings struct {
	CountryID          bson.ObjectID `bson:"country_id" json:"country_id"`
	Language           Language      `bson:"language" json:"language"`
	Appearance         Appearance    `bson:"appearance" json:"appearance"`
	EmailNotifications bool          `bson:"email_notifications" json:"email_notifications"`
}

type UserUpdateInput struct {
	Name               *string
	Username           *string
	Email              *string
	Password           *string
	AvatarURL          *string
	SocialLinks        []string
	CountryID          *bson.ObjectID
	Language           *Language
	Appearance         *Appearance
	EmailNotifications *bool
}
