package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserService interface {
	SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
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

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *userService) SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error) {
	id := bson.NewObjectID()
	hashedPassword := input.Password // TODO

	err := s.repository.Create(ctx, domain.User{
		ID:                      id,
		Name:                    input.Name,
		Username:                input.Username,
		Email:                   input.Email,
		Password:                hashedPassword,
		AvatarURL:               "",
		Points:                  0,
		Favorites:               nil,
		Subscription:            domain.NewSubscription(domain.FreeSubscription),
		NotificationPreferences: domain.NotificationPreferences{Email: true, Push: true},
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	})
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, id)
}

type UserSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *userService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	user, err := s.repository.GetByCredentials(ctx, input.UsernameOrEmail, input.Password)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *userService) GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error) {
	return s.repository.GetByID(ctx, id)
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
	var hashedPassword string
	if input.Password != "" {
		hashedPassword = input.Password // TODO
	}

	err := s.repository.Update(ctx, id, repository.UserUpdateInput{
		Name:                    input.Name,
		Username:                input.Username,
		Email:                   input.Email,
		Password:                hashedPassword,
		AvatarURL:               input.AvatarURL,
		NotificationPreferences: &input.NotificationPreferences,
	})

	return err
}

func (s *userService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Delete(ctx, id)
}

func (s *userService) createSession(ctx context.Context, id bson.ObjectID) (Tokens, error) {
	//TODO implement me
	panic("implement me")
}
