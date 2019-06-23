package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/itslaves/rentalgames-server/common/sessions"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func NaverOAuthCallback(c *gin.Context) {
	clientID := viper.GetString("oauth_naver_client_id")
	clientSecret := viper.GetString("oauth_naver_client_secret")
	authURL := viper.GetString("oauth.naver.authorizeURL")
	tokenURL := viper.GetString("oauth.naver.tokenURL")
	redirectURL := viper.GetString("oauth.naver.redirectURL")

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURL,
	}

	token, err := oauthConfig.Exchange(context.TODO(), c.Query("code"))
	if err != nil {
		// TODO: 에러 유형 파라미터로 전달
		location := fmt.Sprintf("http://%s/error", WEB_SERVER_ADDR)
		c.Redirect(http.StatusFound, location)
	}
	if !token.Valid() {
		// TODO: 에러 유형 파라미터로 전달
		location := fmt.Sprintf("http://%s/error", WEB_SERVER_ADDR)
		c.Redirect(http.StatusFound, location)
	}

	// TODO: MySQL DB 연동 처리
	userID := "itslaves"
	alreadyExists := true

	session := sessions.Session(c)

	if alreadyExists {
		value := OAuthSessionValue{
			UserID:       userID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		jsonBytes, _ := json.Marshal(&value)
		session.Values[userID] = jsonBytes
		if serr := session.Save(c.Request, c.Writer); err != nil {
			fmt.Println(serr)
		}

		location := fmt.Sprintf("http://%s/", WEB_SERVER_ADDR)
		c.Redirect(http.StatusFound, location)
	} else {
		userProfileURL := viper.GetString("oauth.naver.userProfileURL")

		client := oauthConfig.Client(context.TODO(), token)
		resp, err := client.Get(userProfileURL)
		if err != nil {
			// TODO: 에러 유형 파라미터로 전달
			location := fmt.Sprintf("http://%s/error", WEB_SERVER_ADDR)
			c.Redirect(http.StatusFound, location)
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

		location := fmt.Sprintf("http://%s/join?%s", WEB_SERVER_ADDR, params.Encode())
		c.Redirect(http.StatusFound, location)
	}
}
