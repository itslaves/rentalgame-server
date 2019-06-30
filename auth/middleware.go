package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/itslaves/rentalgames-server/common/sessions"
	"github.com/spf13/viper"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Session(c)
		webAddr := viper.GetString("web.addr")
		if userID, ok := session.Values[UserID]; ok {
			// TODO: MySQL DB 연동 부분
			var err error // FIXME: DB 연동 후 제거
			// _, err := db.Select(fmt.Sprintf("SELECT userID FROM users WHERE userID = '%s'", userID))
			if err != nil {
				params := url.Values{}
				params.Set(UserID, userID.(string))
				params.Set(Nickname, session.Values[Nickname].(string))
				params.Set(ProfileImage, session.Values[ProfileImage].(string))
				params.Set(Gender, session.Values[Gender].(string))
				params.Set(Email, session.Values[Email].(string))

				joinPage := fmt.Sprintf("http://%s/join?%s", webAddr, params.Encode())
				c.Redirect(http.StatusTemporaryRedirect, joinPage)
			}
		} else {
			loginPage := fmt.Sprintf("http://%s/login", webAddr)
			c.Redirect(http.StatusTemporaryRedirect, loginPage)
		}
		c.Next()
	}
}
