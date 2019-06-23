package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	WEB_SERVER_ADDR = "www.it-slaves.com:8000"
)

type OAuthSessionValue struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
