package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	rgRedis "github.com/itslaves/rentalgames-server/common/redis"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCallbackHandlerOkPage(t *testing.T) {
	h := &callbackHandler{webAddr: "localhost"}
	assert.Equal(t, "http://localhost/", h.okPage())
}

func TestCallbackHandlerErrPage(t *testing.T) {
	h := &callbackHandler{webAddr: "localhost"}
	assert.Equal(t, "http://localhost/error", h.errPage())
}

func TestCallbackHandlerValidateState(t *testing.T) {
	viper.Set("web.addr", "localhost")

	router := newRouterWithSession()
	router.GET("/success", func(ctx *gin.Context) {
		session := rgSessions.Session(ctx)
		session.Values[State] = "1234"

		h := NewCallbackHandler(ctx, VendorGoogle, GoogleOAuthConfig(), nil)
		assert.NoError(t, h.validateState())
	})
	performRequest(router, "GET", "/success?state=1234")

	router.GET("/fail", func(ctx *gin.Context) {
		session := rgSessions.Session(ctx)
		session.Values[State] = "4321"

		h := NewCallbackHandler(ctx, VendorGoogle, GoogleOAuthConfig(), nil)
		assert.Error(t, h.validateState())
	})
	performRequest(router, "GET", "/fail?state=1234")
}

func TestCallbackHandlerLoadToken(t *testing.T) {
	viper.Set("web.addr", "localhost")

	router := newRouterWithSession()
	router.GET("/success", func(ctx *gin.Context) {
		config := new(oauthConfigMock)
		config.On("Exchange", mock.Anything, mock.Anything, mock.Anything).
			Return(&oauth2.Token{AccessToken: "at", Expiry: time.Now().Add(60 * time.Second)}, nil)

		h := NewCallbackHandler(ctx, VendorGoogle, config, nil)
		assert.NoError(t, h.loadToken())
	})
	performRequest(router, "GET", "/success?code=1234")

	router.GET("/fail", func(ctx *gin.Context) {
		config := new(oauthConfigMock)
		config.On("Exchange", mock.Anything, mock.Anything, mock.Anything).
			Return(&oauth2.Token{AccessToken: "at", Expiry: time.Now()}, nil)

		h := NewCallbackHandler(ctx, VendorGoogle, config, nil)
		assert.Error(t, h.loadToken())
	})
	performRequest(router, "GET", "/fail?code=1234")
}

func TestCallbackHandlerLoadUserProfile(t *testing.T) {
	viper.Set("web.addr", "localhost")
	viper.Set("oauth.google.userProfileURL", "http://google/userProfile")

	router := newRouterWithSession()
	router.GET("/success", func(ctx *gin.Context) {
		output := []byte(`{"name": "user"}`)
		config := new(oauthConfigMock)
		config.On("Client", mock.Anything, mock.Anything).
			Return(&http.Client{
				Transport: &transportMock{
					h: func(req *http.Request) *http.Response {
						rec := httptest.NewRecorder()
						rec.Write(output)
						return rec.Result()
					},
				},
			})

		h := NewCallbackHandler(ctx, VendorGoogle, config, func(resp []byte) *UserProfile {
			assert.Equal(t, output, resp)
			return nil
		})
		assert.NoError(t, h.loadUserProfile())
	})
	performRequest(router, "GET", "/success?code=1234")

	// router.GET("/fail", func(ctx *gin.Context) {
	// 	config := new(oauthConfigMock)
	// 	config.On("Exchange", mock.Anything, mock.Anything, mock.Anything).
	// 		Return(&oauth2.Token{AccessToken: "at", Expiry: time.Now()}, nil)

	// 	h := NewCallbackHandler(ctx, VendorGoogle, config, nil)
	// 	assert.Error(t, h.loadToken())
	// })
	// performRequest(router, "GET", "/fail?code=1234")
}
