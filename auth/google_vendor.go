package auth

import (
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("oauth_google_client_id"),
		ClientSecret: viper.GetString("oauth_google_client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("oauth.google.authorizeURL"),
			TokenURL: viper.GetString("oauth.google.tokenURL"),
		},
		RedirectURL: viper.GetString("oauth.google.redirectURL"),
	}
}

func GoogleOAuthCallback(ctx *gin.Context) {
	oauthConfig := GoogleOAuthConfig()
	userProfileParser := func(result []byte) *UserProfile {
		id, _ := jsonparser.GetString(result, "sub")
		nickname, _ := jsonparser.GetString(result, "name")
		profileImage, _ := jsonparser.GetString(result, "picture")
		email, _ := jsonparser.GetString(result, "email")

		return &UserProfile{
			UserID:       id,
			Nickname:     nickname,
			ProfileImage: profileImage,
			Gender:       "",
			Email:        email,
		}
	}
	handler := NewCallbackHandler(ctx, VendorGoogle, oauthConfig, userProfileParser)
	handler.Handle()
}
