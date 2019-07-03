package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	rgMySQL "github.com/itslaves/rentalgames-server/common/mysql"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
	rgUser "github.com/itslaves/rentalgames-server/user"
	"github.com/spf13/viper"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := rgSessions.Session(c)
		webAddr := viper.GetString("web.addr")

		if userID, ok := session.Values[UserID]; ok {
			vendor := session.Values[Vendor]
			db := rgMySQL.Client()
			condTemplate := "oauth_vendor = ? AND oauth_id = ?"
			if err := db.Where(condTemplate, vendor, userID).First(&rgUser.User{}).Error; err != nil {
				params := url.Values{}
				params.Set(UserID, userID.(string))
				params.Set(Nickname, session.Values[Nickname].(string))
				params.Set(ProfileImage, session.Values[ProfileImage].(string))
				params.Set(Gender, session.Values[Gender].(string))
				params.Set(Email, session.Values[Email].(string))

				joinPage := fmt.Sprintf("http://%s/join?%s", webAddr, params.Encode())
				c.JSON(http.StatusUnauthorized, gin.H{
					"redirect": joinPage,
				})
				c.Abort()
			}
		} else {
			loginPage := fmt.Sprintf("http://%s/login", webAddr)
			c.JSON(http.StatusUnauthorized, gin.H{
				"redirect": loginPage,
			})
			c.Abort()
		}
		c.Next()
	}
}
