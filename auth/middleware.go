package auth

import (
	"github.com/itslaves/rentalgames-server/common/errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	rgMySQL "github.com/itslaves/rentalgames-server/common/mysql"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
	rgUser "github.com/itslaves/rentalgames-server/user"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := rgSessions.Session(c)

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

				c.JSON(http.StatusUnauthorized, errors.ErrorResponse(errors.Unauthorized))
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse(errors.LoginRequired))
			c.Abort()
		}
		c.Next()
	}
}
