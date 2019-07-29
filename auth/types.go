package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthConfig interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	PasswordCredentialsToken(ctx context.Context, username, password string) (*oauth2.Token, error)
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}
