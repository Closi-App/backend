package repository

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const userCollectionName = "users"

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error
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
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

type UserUpdateInput struct {
	Name                    string
	Username                string
	Email                   string
	Password                string
	AvatarURL               string
	NotificationPreferences domain.NotificationPreferences
}

func (r *userRepository) Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	//TODO implement me
	panic("implement me")
}
