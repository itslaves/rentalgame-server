package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type OAuthSessionValue struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func RandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
