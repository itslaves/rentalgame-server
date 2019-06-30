package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func Register(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := sessions.GetRegistry(c.Request)
		s, _ := r.Get(store, name)
		c.Set(ContextKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}
