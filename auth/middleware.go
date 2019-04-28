package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
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