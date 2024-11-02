package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

type TokensManager interface {
	NewAccessToken(userID string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type tokensManager struct {
	accessTokenSigningKey string
	accessTokenTTL        time.Duration
	refreshTokenLength    int
}

func NewTokensManager(cfg *viper.Viper) TokensManager {
	return &tokensManager{
		accessTokenSigningKey: cfg.GetString("auth.tokens.access_token.signing_key"),
		accessTokenTTL:        cfg.GetDuration("auth.tokens.access_token.ttl"),
		refreshTokenLength:    cfg.GetInt("auth.tokens.refresh_token.length"),
	}
}

func (m *tokensManager) NewAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(m.accessTokenTTL).Unix(),
		})

	return token.SignedString([]byte(m.accessTokenSigningKey))
}

func (m *tokensManager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.accessTokenSigningKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "error parsing access token")
	}

	if !token.Valid {
		return "", errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error getting claims from access token")
	}

	return claims["sub"].(string), nil
}

func (m *tokensManager) NewRefreshToken() (string, error) {
	b := make([]byte, m.refreshTokenLength)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
