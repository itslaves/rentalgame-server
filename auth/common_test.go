package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"

	rgRedis "github.com/itslaves/rentalgames-server/common/redis"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
)

func newRouterWithSession() *gin.Engine {
	router := gin.New()

	sessionName := "rg_session"
	sessionStore := rgSessions.NewRedisStore(
		rgRedis.TestClient(),
		&sessions.Options{
			Path:   "/",
			MaxAge: 60,
			Domain: "localhost",
		},
		[]byte("secret"),
	)
	router.Use(rgSessions.Register(sessionName, sessionStore))

	return router
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestOAuthVendorUrls(t *testing.T) {
	router := newRouterWithSession()
	router.GET("/oauth/vendors", OAuthVendorUrls)

	w := performRequest(router, "GET", "/oauth/vendors")
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, 1, len(w.Result().Cookies()))

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "rg_session", cookie.Name)
	assert.Equal(t, 60, cookie.MaxAge)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, "localhost", cookie.Domain)
}
