package auth

import (
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	rgSessions "github.com/skyoo2003/rentalgames-server/internal/sessions"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func NaverOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("oauth_naver_client_id"),
		ClientSecret: viper.GetString("oauth_naver_client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("oauth.naver.authorizeURL"),
			TokenURL: viper.GetString("oauth.naver.tokenURL"),
		},
		RedirectURL: viper.GetString("oauth.naver.redirectURL"),
	}
}

func NaverOAuthCallback(c *gin.Context) {
	session := rgSessions.Session(c)
	oauthConfig := NaverOAuthConfig()
	userProfileParser := func(result []byte) *UserProfile {
		id, _ := jsonparser.GetString(result, "response", "id")
		nickname, _ := jsonparser.GetString(result, "response", "nickname")
		profileImage, _ := jsonparser.GetString(result, "response", "profile_image")
		gender, _ := jsonparser.GetString(result, "response", "gender")
		if gender == "M" {
			gender = "male"
		} else {
			gender = "female"
		}
		email, _ := jsonparser.GetString(result, "response", "email")

		return &UserProfile{
			UserID:       id,
			Nickname:     nickname,
			ProfileImage: profileImage,
			Gender:       gender,
			Email:        email,
		}
	}
	handler := NewCallbackHandler("naver", session, oauthConfig, c, userProfileParser)
	handler.Handle()
}
