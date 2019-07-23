package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	rgRedis "github.com/itslaves/rentalgames-server/common/redis"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
)

type header struct {
	Key   string
	Value string
}

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

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestOAuthVendorUrls(t *testing.T) {
	router := newRouterWithSession()
	router.GET("/oauth/vendors", OAuthVendorUrls)

	w := performRequest(router, "GET", "/oauth/vendors")
	assert.Equal(t, http.StatusOK, w.Code)

	h := w.Header()

	cookieMap := make(map[string]interface{})
	if values, ok := h["Set-Cookie"]; ok {
		for _, item := range strings.Split(values[0], ";") {
			kv := strings.SplitN(item, "=", 2)
			kv[0] = strings.TrimSpace(kv[0])
			kv[1] = strings.TrimSpace(kv[1])
			cookieMap[kv[0]] = kv[1]
		}
	} else {
		t.Fatal("'Set-Cookie' does not exist in the header")
	}

	if value, ok := cookieMap["Max-Age"]; ok {
		if value != "60" {
			t.Fatal("'Max-Age' was unexpected: ", value)
		}
	} else {
		t.Fatal("'Max-Age' does not exist in the cookie")
	}

	if value, ok := cookieMap["Path"]; ok {
		if value != "/" {
			t.Fatal("'Path' was unexpected: ", value)
		}
	} else {
		t.Fatal("'Path' does not exist in the cookie")
	}

	if value, ok := cookieMap["Domain"]; ok {
		if value != "localhost" {
			t.Fatal("'Domain' was unexpected: ", value)
		}
	} else {
		t.Fatal("'Domain' does not exist in the cookie")
	}

	if _, ok := cookieMap["rg_session"]; !ok {
		t.Fatal("'rg_session' does not exist in the cookie")
	}
}
