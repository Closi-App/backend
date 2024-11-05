package domain

import "time"

type SubscriptionType string

// TODO: storing subscriptions in db

const (
	FreeSubscription      SubscriptionType = "free"
	MonthlySubscription   SubscriptionType = "monthly"
	QuarterlySubscription SubscriptionType = "quarterly"
	AnnualSubscription    SubscriptionType = "annual"
)

type Subscription struct {
	Type      SubscriptionType `bson:"type" json:"type"`
	ExpiresAt time.Time        `bson:"expires_at" json:"expires_at"`
}

func NewSubscription(subscriptionType SubscriptionType) Subscription {
	subscription := Subscription{
		Type: subscriptionType,
	}

	switch subscription.Type {
	case FreeSubscription:
		subscription.ExpiresAt = time.Time{}
	case MonthlySubscription:
		subscription.ExpiresAt = time.Now().AddDate(0, 1, 0)
	case QuarterlySubscription:
		subscription.ExpiresAt = time.Now().AddDate(0, 3, 0)
	case AnnualSubscription:
		subscription.ExpiresAt = time.Now().AddDate(1, 0, 0)
	}

	return subscription
}

func (s Subscription) IsActive() bool {
	return time.Now().Before(s.ExpiresAt) || s.Type == FreeSubscription
}
