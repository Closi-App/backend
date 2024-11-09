package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func NewReferralCode(length int) (string, error) {
	b := make([]byte, length)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	code := fmt.Sprintf("%x", b)

	return strings.ToUpper(code), nil
}
