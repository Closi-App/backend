package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	SignUp(ctx context.Context, input UserSignUpInput) (string, error)
	SignIn(ctx context.Context, input UserSignInInput) (string, error)
	GetByID(ctx context.Context, id bson.ObjectID)
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type userService struct {
	*Service
	repository repository.UserRepository
}

func NewUserService(service *Service, repository repository.UserRepository) UserService {
	return &userService{
		Service:    service,
		repository: repository,
	}
}

type UserSignUpInput struct {
	Name     string
	Username string
	Email    string
	Password string
}

func (s *userService) SignUp(ctx context.Context, input UserSignUpInput) (string, error) {
	//TODO implement me
	panic("implement me")
}

type UserSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *userService) SignIn(ctx context.Context, input UserSignInInput) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) GetByID(ctx context.Context, id bson.ObjectID) {
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

func (s *userService) Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error {
	//TODO implement me
	panic("implement me")
}

func (s *userService) Delete(ctx context.Context, id bson.ObjectID) error {
	//TODO implement me
	panic("implement me")
}
