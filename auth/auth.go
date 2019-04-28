package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

const (
	loginUrl = "http://it-slaves.com:8080/login"

	callbackUrl = "http://it-slaves.com:8080/auth/callback"
	userInfoApi = "https://www.googleapis.com/oauth2/v3/userinfo"

	scopeEmail = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile = "https://www.googleapis.com/auth/userinfo.profile"

	userSessionKey = "userSessionKey"
	stateSessionKey = "stateSessionKey"
)

type Config struct {
	oauth2.Config
	ProviderName string
	UserInfoApi string
}

var config []Config
func init() {
	config = []Config {
		{
			oauth2.Config {
				ClientID:     "489302416615-1ila1jvd9vrhqam10r8p17iqmqp78b3m.apps.googleusercontent.com",
				ClientSecret: "bb3wtKIkEsQ0aGLQqPDgc1xi",
				RedirectURL:  callbackUrl,
				Scopes:       []string{scopeEmail, scopeProfile},
			},
			"google",
			"https://www.googleapis.com/oauth2/v3/userinfo",
		},
	}
}

var OAuthConf *oauth2.Config

func init() {
	OAuthConf = &oauth2.Config{
		ClientID:     "489302416615-1ila1jvd9vrhqam10r8p17iqmqp78b3m.apps.googleusercontent.com",
		ClientSecret: "bb3wtKIkEsQ0aGLQqPDgc1xi",
		RedirectURL:  callbackUrl,
		Scopes:       []string{scopeEmail, scopeProfile},
		Endpoint:     google.Endpoint,
	}
}

func GetLoginUrl(state string) string {
	return OAuthConf.AuthCodeURL(state)
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

type User struct {
	Uid	string `json:"sub"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func CurrentUser(c *gin.Context) *User {
	session := sessions.Default(c)
	v := session.Get(userSessionKey)

	if v == nil {
		return nil
	}

	data := v.([]byte)

	var user User
	json.Unmarshal(data, &user)
	return &user
}

func LoginView(c *gin.Context) {
	session := sessions.Default(c)
	// TODO 세션 유지기간 관련 설정 필요
	state := RandToken()
	session.Set(stateSessionKey, state)
	c.HTML(http.StatusOK, "login.html", gin.H {
		"loginUrl" : GetLoginUrl(state),
	})
}

func Authenticate(c *gin.Context) {
	session := sessions.Default(c)
	// TODO state값 비교 / 제거 추가

	token, err := OAuthConf.Exchange(context.TODO(), c.Query("code"))
	if err != nil {
		c.Error(err)
		return
	}

	client := OAuthConf.Client(context.TODO(), token)
	userResp, err := client.Get(userInfoApi)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer userResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user User
	json.Unmarshal(userInfo, &user)
	userJson, err  := json.Marshal(&user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Set(userSessionKey, userJson)
	session.Save()

	c.Redirect(http.StatusFound, "/index")
}