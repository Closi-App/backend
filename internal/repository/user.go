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
	Update(ctx context.Context, id bson.ObjectID, input domain.UserUpdateInput) error
	UpdateSettings(ctx context.Context, id bson.ObjectID, input domain.UserSettingsUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error

	AdjustPoints(ctx context.Context, id bson.ObjectID, pointsAmount int) error
	AddFavorite(ctx context.Context, id, questionID bson.ObjectID) error
	RemoveFavorite(ctx context.Context, id, questionID bson.ObjectID) error
	AddAchievement(ctx context.Context, id, achievementID bson.ObjectID) error
	RemoveAchievement(ctx context.Context, id, achievementID bson.ObjectID) error
	SetSubscription(ctx context.Context, id bson.ObjectID, subscription domain.Subscription) error
	Confirm(ctx context.Context, id bson.ObjectID) error
	Block(ctx context.Context, id bson.ObjectID) error
	Unblock(ctx context.Context, id bson.ObjectID) error

	CreateSession(ctx context.Context, refreshToken string, userID bson.ObjectID, expiration time.Duration) error
	GetSession(ctx context.Context, refreshToken string) (userID bson.ObjectID, err error)
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

func (r *userRepository) Update(ctx context.Context, id bson.ObjectID, input domain.UserUpdateInput) error {
	updateFields := bson.M{}

	if input.Name != nil {
		updateFields["name"] = input.Name
	}
	if input.Username != nil {
		updateFields["username"] = input.Username
	}
	if input.Email != nil {
		updateFields["email"] = input.Email
		updateFields["is_confirmed"] = false
	}
	if input.Password != nil {
		updateFields["password"] = input.Password
	}
	if input.AvatarURL != nil {
		updateFields["avatar_url"] = input.AvatarURL
	}

	updateFields["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})

	return err
}

func (r *userRepository) UpdateSettings(ctx context.Context, id bson.ObjectID, input domain.UserSettingsUpdateInput) error {
	updateFields := bson.M{}

	if input.CountryID != nil {
		updateFields["settings.country_id"] = input.CountryID
	}
	if input.Language != nil {
		updateFields["settings.language"] = input.Language
	}
	if input.Appearance != nil {
		updateFields["settings.appearance"] = input.Appearance
	}
	if input.EmailNotifications != nil {
		updateFields["settings.email_notifications"] = input.EmailNotifications
	}

	updateFields["updated_at"] = time.Now()

	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})

	return err
}

func (r *userRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *userRepository) AdjustPoints(ctx context.Context, id bson.ObjectID, pointsAmount int) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"points": pointsAmount}})

	return err
}

func (r *userRepository) AddFavorite(ctx context.Context, id, questionID bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"favorites": questionID}})

	return err
}

func (r *userRepository) RemoveFavorite(ctx context.Context, id, questionID bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"favorites": questionID}})

	return err
}

func (r *userRepository) AddAchievement(ctx context.Context, id, achievementID bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"achievements": achievementID}})

	return err
}

func (r *userRepository) RemoveAchievement(ctx context.Context, id, achievementID bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"achievements": achievementID}})

	return err
}

func (r *userRepository) SetSubscription(ctx context.Context, id bson.ObjectID, subscription domain.Subscription) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"subscription": subscription}})

	return err
}

func (r *userRepository) Confirm(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_confirmed": true}})

	return err
}

func (r *userRepository) Block(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_blocked": true}})

	return err
}

func (r *userRepository) Unblock(ctx context.Context, id bson.ObjectID) error {
	_, err := r.db.Collection(domain.UserCollectionName).
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_blocked": false}})

	return err
}

func (r *userRepository) CreateSession(ctx context.Context, refreshToken string, userID bson.ObjectID, expiration time.Duration) error {
	key := fmt.Sprintf(dbSessionKeyFormat, refreshToken)
	value := userID.Hex()

	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r *userRepository) GetSession(ctx context.Context, refreshToken string) (bson.ObjectID, error) {
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
