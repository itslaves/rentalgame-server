package auth

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	rgErrors "github.com/skyoo2003/rentalgames-server/internal/errors"
	rgSessions "github.com/skyoo2003/rentalgames-server/internal/sessions"
	rgMySQL "github.com/skyoo2003/rentalgames-server/internal/third_party/mysql"
	rgUser "github.com/skyoo2003/rentalgames-server/internal/user"
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

				c.JSON(http.StatusUnauthorized, rgErrors.ErrorResponse(rgErrors.JoinRequired))
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, rgErrors.ErrorResponse(rgErrors.LoginRequired))
			c.Abort()
		}
		c.Next()
	}
}
