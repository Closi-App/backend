package domain

import "time"

type Session struct {
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at" json:"expires_at"`
}
