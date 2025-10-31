package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateOauthStateString gera um 'state' aleatório
func GenerateOauthStateString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
