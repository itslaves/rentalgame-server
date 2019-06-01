package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			c.Redirect(http.StatusFound, loginUrl)
		}

		c.Next()
	}
}
