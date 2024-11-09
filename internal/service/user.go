package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/Closi-App/backend/internal/utils"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserService interface {
	SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, input domain.UserUpdateInput) error
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

	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
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
	CountryID    bson.ObjectID
	Language     domain.Language
	ReferrerCode string
}

func (s *userService) SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error) {
	id := bson.NewObjectID()

	hashedPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	referralCode, err := utils.NewReferralCode(domain.UserReferralCodeLength)
	if err != nil {
		return Tokens{}, err
	}

	if err := s.repository.Create(ctx, domain.User{
		ID:           id,
		Name:         input.Name,
		Username:     input.Username,
		Email:        input.Email,
		Password:     hashedPassword,
		AvatarURL:    "",
		Points:       domain.UserDefaultPoints,
		Favorites:    nil,
		Achievements: nil,
		SocialLinks:  nil,
		ReferralCode: referralCode,
		Subscription: domain.NewSubscription(domain.FreeSubscription),
		Settings: domain.UserSettings{
			CountryID:          input.CountryID,
			Language:           input.Language,
			EmailNotifications: true,
			Appearance:         domain.LightAppearance,
		},
		IsConfirmed: false,
		IsBlocked:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return Tokens{}, err
	}

	if input.ReferrerCode != "" {
		referrer, err := s.repository.GetByReferralCode(ctx, input.ReferrerCode)
		if err == nil {
			if err := s.repository.AdjustPoints(ctx, referrer.ID, domain.UserReferralPoints); err != nil {
				return Tokens{}, err
			}
			if err := s.repository.AdjustPoints(ctx, id, domain.UserReferralPoints); err != nil {
				return Tokens{}, err
			}
		}
	}

	// TODO: sending confirmation email with confirmation link

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

func (s *userService) GetByID(ctx context.Context, id bson.ObjectID) (domain.User, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id bson.ObjectID, input domain.UserUpdateInput) error {
	if input.Password != nil {
		hashedPassword, err := s.passwordHasher.Hash(*input.Password)
		if err != nil {
			return err
		}

		input.Password = &hashedPassword
	}

	// TODO: sending confirmation email with confirmation link if email was updated

	return s.repository.Update(ctx, id, input)
}

func (s *userService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Delete(ctx, id)
}

func (s *userService) AdjustPoints(ctx context.Context, id bson.ObjectID, pointsAmount int) error {
	return s.repository.AdjustPoints(ctx, id, pointsAmount)
}

func (s *userService) AddFavorite(ctx context.Context, id, questionID bson.ObjectID) error {
	return s.repository.AddFavorite(ctx, id, questionID)
}

func (s *userService) RemoveFavorite(ctx context.Context, id, questionID bson.ObjectID) error {
	return s.repository.RemoveFavorite(ctx, id, questionID)
}

func (s *userService) AddAchievement(ctx context.Context, id, achievementID bson.ObjectID) error {
	return s.repository.AddAchievement(ctx, id, achievementID)
}

func (s *userService) RemoveAchievement(ctx context.Context, id, achievementID bson.ObjectID) error {
	return s.repository.RemoveAchievement(ctx, id, achievementID)
}

func (s *userService) SetSubscription(ctx context.Context, id bson.ObjectID, subscription domain.Subscription) error {
	return s.repository.SetSubscription(ctx, id, subscription)
}

func (s *userService) Confirm(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Confirm(ctx, id)
}

func (s *userService) Block(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Block(ctx, id)
}

func (s *userService) Unblock(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Unblock(ctx, id)
}

func (s *userService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	userID, err := s.repository.GetSession(ctx, refreshToken)
	if err != nil {
		return Tokens{}, domain.ErrUnauthorized
	}

	return s.createSession(ctx, userID)
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
