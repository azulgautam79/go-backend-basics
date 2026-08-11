package user

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
