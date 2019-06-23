package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itslaves/rentalgames-server/common/sessions"
	"github.com/spf13/viper"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Session(c)
		if _, exist := session.Values["userID"]; !exist {
			webAddr := viper.GetString("web.addr")
			location := fmt.Sprintf("http://%s/login", webAddr)
			c.Redirect(http.StatusFound, location)
		}
		c.Next()
	}
}
