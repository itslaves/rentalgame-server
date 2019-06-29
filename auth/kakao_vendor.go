package auth

import (
	"strconv"

	"github.com/buger/jsonparser"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func KakaoOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("oauth_kakao_client_id"),
		ClientSecret: viper.GetString("oauth_kakao_client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("oauth.kakao.authorizeURL"),
			TokenURL: viper.GetString("oauth.kakao.tokenURL"),
		},
		RedirectURL: viper.GetString("oauth.kakao.redirectURL"),
	}
}

func KakaoOAuthCallback(c *gin.Context) {
	session := rgSessions.Session(c)
	oauthConfig := KakaoOAuthConfig()
	userProfileParser := func(result []byte) *UserProfile {
		id, _ := jsonparser.GetInt(result, "id")
		nickname, _ := jsonparser.GetString(result, "properties", "nickname")
		profileImage, _ := jsonparser.GetString(result, "properties", "profile_image")
		gender, _ := jsonparser.GetString(result, "kakao_account", "gender")
		email, _ := jsonparser.GetString(result, "kakao_account", "email")

		return &UserProfile{
			UserID:       strconv.Itoa(int(id)),
			Nickname:     nickname,
			ProfileImage: profileImage,
			Gender:       gender,
			Email:        email,
		}
	}
	handler := NewCallbackHandler("kakao", session, oauthConfig, c, userProfileParser)
	handler.Handle()
}
