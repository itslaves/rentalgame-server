package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	rgAuth "github.com/itslaves/rentalgames-server/auth"
	rgRedis "github.com/itslaves/rentalgames-server/common/redis"
	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"
)

func Route() *gin.Engine {
	webAddr := viper.GetString("web.addr")
	sessionName := viper.GetString("session.name")
	sessionMaxAge := viper.GetInt("session.maxAge")
	sessionDomain := viper.GetString("session.domain")
	sessionHashKey := []byte(viper.GetString("session_hash_key"))
	sessionBlockKey := []byte(viper.GetString("session_block_key"))

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{fmt.Sprintf("http://%s", webAddr)},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	sessionStore := rgSessions.NewRedisStore(
		rgRedis.Client(),
		&sessions.Options{
			Path:   "/",
			MaxAge: sessionMaxAge,
			Domain: sessionDomain,
		},
		sessionHashKey,
		sessionBlockKey,
	)
	r.Use(rgSessions.Register(sessionName, sessionStore))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	oauth := r.Group("/oauth")
	{
		oauth.GET("/vendors", rgAuth.OAuthVendorUrls)
		oauth.GET("/callback/kakao", rgAuth.KakaoOAuthCallback)
		oauth.GET("/callback/naver", rgAuth.NaverOAuthCallback)
		oauth.GET("/callback/google", rgAuth.GoogleOAuthCallback)
	}

	v1 := r.Group("/v1")
	v1.Use(rgAuth.Authenticate())
	// TODO: v1 API 라우트 추가
	{
		v1.GET("/authentication", func(c *gin.Context) {
			c.JSON(http.StatusOK, nil)
		})
	}

	debug := r.Group("/debug")
	{
		debug.POST("/redis", func(c *gin.Context) {
			redisClient := rgRedis.Client()
			val := c.Query("val")
			key := c.Query("key")
			cmd := redisClient.Set(key, val, 0)
			c.JSON(http.StatusOK, gin.H{
				"message": cmd.String(),
			})
		})
		debug.GET("/redis", func(c *gin.Context) {
			redisClient := rgRedis.Client()
			key := c.Query("key")
			c.JSON(http.StatusOK, gin.H{
				"message": redisClient.Get(key).String(),
			})
		})
	}

	return r
}
