package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserService interface {
	SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type userService struct {
	*Service
	repository      repository.UserRepository
	passwordHasher  auth.PasswordHasher
	tokensManager   auth.TokensManager
	refreshTokenTTL time.Duration
}

func NewUserService(
	service *Service,
	cfg *viper.Viper,
	repository repository.UserRepository,
	passwordHasher auth.PasswordHasher,
	tokensManager auth.TokensManager,
) UserService {
	return &userService{
		Service:         service,
		repository:      repository,
		passwordHasher:  passwordHasher,
		tokensManager:   tokensManager,
		refreshTokenTTL: cfg.GetDuration("auth.tokens.refresh_token.ttl"),
	}
}

type UserSignUpInput struct {
	Name         string
	Username     string
	Email        string
	Password     string
	Location     domain.Location
	Language     domain.Language
	ReferrerCode string
}

func (s *userService) SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error) {
	id := bson.NewObjectID()

	hashedPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	referralCode, err := domain.NewReferralCode()
	if err != nil {
		return Tokens{}, err
	}

	input.Language = domain.ParseLanguage(input.Language)

	if err := s.repository.Create(ctx, domain.User{
		ID:           id,
		Name:         input.Name,
		Username:     input.Username,
		Email:        input.Email,
		Password:     hashedPassword,
		AvatarURL:    "",
		Points:       domain.UserDefaultPoints,
		Favorites:    nil,
		ReferralCode: referralCode,
		Subscription: domain.NewSubscription(domain.FreeSubscription),
		Settings: domain.UserSettings{
			Location:           input.Location,
			Language:           input.Language,
			EmailNotifications: true,
		},
		IsConfirmed: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return Tokens{}, err
	}

	if input.ReferrerCode != "" {
		referrer, err := s.repository.GetByReferralCode(ctx, input.ReferrerCode)
		if err == nil {
			if err := s.repository.AddPoints(ctx, referrer.ID, domain.UserReferralPoints); err != nil {
				return Tokens{}, err
			}
			if err := s.repository.AddPoints(ctx, id, domain.UserReferralPoints); err != nil {
				return Tokens{}, err
			}
		}
	}

	// TODO: sending confirmation email

	return s.createSession(ctx, id)
}

type UserSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *userService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	user, err := s.repository.GetByUsernameOrEmail(ctx, input.UsernameOrEmail)
	if err != nil {
		return Tokens{}, err
	}

	if !s.passwordHasher.Check(user.Password, input.Password) {
		return Tokens{}, domain.ErrUserNotFound
	}

	return s.createSession(ctx, user.ID)
}

func (s *userService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	userID, err := s.repository.GetSessionUserID(ctx, refreshToken)
	if err != nil {
		return Tokens{}, domain.ErrUnauthorized
	}

	return s.createSession(ctx, userID)
}

// TODO: user confirmation

func (s *userService) GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error) {
	return s.repository.GetByID(ctx, id)
}

type UserUpdateInput struct {
	Name      string
	Username  string
	Email     string
	Password  string
	AvatarURL string
	Settings  domain.UserSettings
}

func (s *userService) Update(ctx context.Context, id bson.ObjectID, input UserUpdateInput) error {
	var (
		hashedPassword string
		err            error
	)

	if input.Password != "" {
		hashedPassword, err = s.passwordHasher.Hash(input.Password)
		if err != nil {
			return err
		}
	}

	if input.Settings.Language != "" {
		input.Settings.Language = domain.ParseLanguage(input.Settings.Language)
	}

	// TODO: sending confirmation email if email was updated

	return s.repository.Update(ctx, id, repository.UserUpdateInput{
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
		AvatarURL: input.AvatarURL,
		Settings:  &input.Settings,
	})
}

func (s *userService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Delete(ctx, id)
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *userService) createSession(ctx context.Context, id bson.ObjectID) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.AccessToken, err = s.tokensManager.NewAccessToken(id.Hex())
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = s.tokensManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	err = s.repository.CreateSession(ctx, tokens.RefreshToken, id, s.refreshTokenTTL)

	return tokens, err
}
