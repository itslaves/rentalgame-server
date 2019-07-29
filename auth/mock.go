package auth

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type transportMock struct {
	h func(req *http.Request) *http.Response
}

func (t *transportMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.h(req), nil
}

type oauthConfigMock struct {
	mock.Mock
}

func (c *oauthConfigMock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	args := c.Called(state, opts)
	return args.String(0)
}

func (c *oauthConfigMock) PasswordCredentialsToken(ctx context.Context, username, password string) (*oauth2.Token, error) {
	args := c.Called(ctx, username, password)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (c *oauthConfigMock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	args := c.Called(ctx, code, opts)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (c *oauthConfigMock) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	args := c.Called(ctx, t)
	return args.Get(0).(*http.Client)
}

func (c *oauthConfigMock) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	args := c.Called(ctx, t)
	return args.Get(0).(oauth2.TokenSource)
}
