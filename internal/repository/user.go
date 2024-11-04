package repository

import (
	"context"
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error
	SetSession(ctx context.Context, id bson.ObjectID, session domain.Session) error
}

type userRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) UserRepository {
	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	usernameIndex := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := repository.db.Collection(domain.UserCollectionName).
		Indexes().CreateMany(context.Background(), []mongo.IndexModel{emailIndex, usernameIndex}); err != nil {
		panic("error creating user indexes: " + err.Error())
	}

	return &userRepository{
		Repository: repository,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Collection(domain.UserCollectionName).
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

	err := r.db.Collection(domain.UserCollectionName).
		FindOne(ctx, bson.M{"_id": id}).Decode(&user)
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

	err := r.db.Collection(domain.UserCollectionName).
		FindOne(ctx, bson.M{"$or": []interface{}{
			bson.M{"username": usernameOrEmail},
			bson.M{"email": usernameOrEmail},
		}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User

	err := r.db.Collection(domain.UserCollectionName).
		FindOne(ctx, bson.M{
			"session.refresh_token": refreshToken,
			"session.expires_at":    bson.M{"$gt": time.Now()},
		}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
	}

	return user, nil
}

type UserUpdateInput struct {
	Name                    string
	Username                string
	Email                   string
	Password                string
	AvatarURL               string
	Location                *domain.Location
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
	if input.Location != nil {
		updateQuery["location"] = input.Location
	}
	if input.NotificationPreferences != nil {
		updateQuery["notification_preferences"] = input.NotificationPreferences
	}

	updateQuery["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateQuery})

	return err
}

func (r *userRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *userRepository) SetSession(ctx context.Context, id bson.ObjectID, session domain.Session) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session}})

	return err
}
