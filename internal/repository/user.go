package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const userCollectionName = "users"

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error
	SetSession(ctx context.Context, id bson.ObjectID, session domain.Session) error
}

type userRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Collection(userCollectionName).
		InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error) {
	var user domain.User

	err := r.db.Collection(userCollectionName).
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error) {
	var user domain.User

	err := r.db.Collection(userCollectionName).
		FindOne(ctx, bson.M{"$or": []interface{}{
			bson.M{"username": usernameOrEmail},
			bson.M{"email": usernameOrEmail},
		}}).
		Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

type UserUpdateInput struct {
	Name                    string
	Username                string
	Email                   string
	Password                string
	AvatarURL               string
	NotificationPreferences *domain.NotificationPreferences
}

func (r *userRepository) Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error {
	updateQuery := bson.M{}

	if input.Name != "" {
		updateQuery["name"] = input.Name
	}
	if input.Username != "" {
		updateQuery["username"] = input.Username
	}
	if input.Email != "" {
		updateQuery["email"] = input.Email
	}
	if input.Password != "" {
		updateQuery["password"] = input.Password
	}
	if input.AvatarURL != "" {
		updateQuery["avatar_url"] = input.AvatarURL
	}
	if input.NotificationPreferences != nil {
		updateQuery["notification_preferences"] = input.NotificationPreferences
	}

	_, err := r.db.Collection(userCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateQuery})

	return err
}

func (r *userRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(userCollectionName).
		DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *userRepository) SetSession(ctx context.Context, id bson.ObjectID, session domain.Session) error {
	_, err := r.db.Collection(userCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session}})

	return err
}
