package auth

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(hashedPassword, password string) bool
}

type passwordHasher struct {
	salt string
}

func NewPasswordHasher(cfg *viper.Viper) PasswordHasher {
	return &passwordHasher{
		salt: cfg.GetString("auth.password.salt"),
	}
}

func (h *passwordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+h.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "error hashing password")
	}

	return string(hashedPassword), nil
}

func (h *passwordHasher) Check(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+h.salt)) == nil
}
