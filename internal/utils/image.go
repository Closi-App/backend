package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const imageQuality = 40

func NewImageName(fileBytes []byte) string {
	hash := md5.Sum(fileBytes)
	return hex.EncodeToString(hash[:])
}
