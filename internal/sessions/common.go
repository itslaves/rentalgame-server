package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const ContextKey = "context.session"

func Session(c *gin.Context) *sessions.Session {
	return c.MustGet(ContextKey).(*sessions.Session)
}
