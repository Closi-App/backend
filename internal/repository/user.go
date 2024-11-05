package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

const dbSessionKeyFormat = "session:%s"

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error)
	GetByReferralCode(ctx context.Context, referralCode string) (domain.User, error)
	AddPoints(ctx context.Context, id bson.ObjectID, pointsAmount int) error
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error

	CreateSession(ctx context.Context, refreshToken string, userID bson.ObjectID, expiration time.Duration) error
	GetSessionUserID(ctx context.Context, refreshToken string) (bson.ObjectID, error)
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

func (r *userRepository) GetByReferralCode(ctx context.Context, referralCode string) (domain.User, error) {
	var user domain.User

	err := r.db.Collection(domain.UserCollectionName).
		FindOne(ctx, bson.M{"referral_code": referralCode}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) AddPoints(ctx context.Context, id bson.ObjectID, pointsAmount int) error {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Points+pointsAmount < 0 {
		return domain.ErrUserInsufficientPoints
	}

	_, err = r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"points": pointsAmount}})

	return err
}

type UserUpdateInput struct {
	Name      string
	Username  string
	Email     string
	Password  string
	AvatarURL string
	Settings  *domain.UserSettings
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
		updateQuery["is_confirmed"] = false
	}
	if input.Password != "" {
		updateQuery["password"] = input.Password
	}
	if input.AvatarURL != "" {
		updateQuery["avatar_url"] = input.AvatarURL
	}
	if input.Settings != nil {
		updateQuery["settings"] = input.Settings
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

func (r *userRepository) CreateSession(ctx context.Context, refreshToken string, userID bson.ObjectID, expiration time.Duration) error {
	key := fmt.Sprintf(dbSessionKeyFormat, refreshToken)
	value := userID.Hex()

	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r *userRepository) GetSessionUserID(ctx context.Context, refreshToken string) (bson.ObjectID, error) {
	key := fmt.Sprintf(dbSessionKeyFormat, refreshToken)

	value, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return bson.ObjectID{}, err
	}

	userID, err := bson.ObjectIDFromHex(value)
	if err != nil {
		return bson.ObjectID{}, err
	}

	return userID, nil
}
