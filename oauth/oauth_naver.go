package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	AuthorizeURL   = "https://nid.naver.com/oauth2.0/authorize"
	TokenURL       = "https://nid.naver.com/oauth2.0/token"
	UserProfileURL = "https://openapi.naver.com/v1/nid/me"
)

func NaverOAuthCallback(c *gin.Context) {
	clientID := viper.GetString("oauth.naver.clientID")
	clientSecret := viper.GetString("oauth.naver.clientSecret")
	redirectURL := viper.GetString("oauth.naver.redirectURL")

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthorizeURL,
			TokenURL: TokenURL,
		},
		RedirectURL: redirectURL,
	}

	token, err := oauthConfig.Exchange(context.TODO(), c.Query("code"))
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/error")
	}
	if !token.Valid() {
		// TODO: 에러 유형 파라미터로 전달
		c.Redirect(http.StatusBadRequest, "/error")
	}

	// TODO: MySQL DB 연동 처리
	userID := "itslaves"
	alreadyExists := true

	sessionKey := viper.GetString("session.key")
	session := sessions.Default(c)

	if alreadyExists {
		value := OAuthSessionValue{
			UserID:       userID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		jsonBytes, _ := json.Marshal(&value)
		session.Set(sessionKey, jsonBytes)
		session.Save()
		c.Redirect(http.StatusFound, "/")
	} else {
		client := oauthConfig.Client(context.TODO(), token)
		resp, err := client.Get(UserProfileURL)
		if err != nil {
			// TODO: 에러 유형 파라미터로 전달
			c.Redirect(http.StatusBadRequest, "/error")
		}
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		// result = {"resultcode":"00","message":"success","response":{"id":"13078265","nickname":"Zicprit","profile_image":"https:\/\/phinf.pstatic.net\/contact\/20180210_208\/1518260323373OPQLs_PNG\/avatar_profile.png","gender":"M","email":"skyoo2003@naver.com","birthday":"09-22"}}

		nickname, _ := jsonparser.GetString(result, "response", "nickname")
		profileImage, _ := jsonparser.GetString(result, "response", "profile_image")
		gender, _ := jsonparser.GetString(result, "response", "gender")
		if gender == "M" {
			gender = "male"
		} else {
			gender = "female"
		}
		email, _ := jsonparser.GetString(result, "response", "email")

		params := url.Values{}
		params.Set("nickname", nickname)
		params.Set("profileImage", profileImage)
		params.Set("gender", gender)
		params.Set("email", email)

		urlPath := fmt.Sprintf("/join?%s", params.Encode())

		c.Redirect(http.StatusFound, urlPath)
	}
}
