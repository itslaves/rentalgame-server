package route

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"

	"github.com/itslaves/rentalgames-server/auth"
	"github.com/itslaves/rentalgames-server/common/redis"
	"github.com/itslaves/rentalgames-server/common/sessions"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	webAddr := viper.GetString("web.addr")
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

	sessionName := viper.GetString("session.name")
	sessionStore := sessions.NewRedisStore()
	r.Use(sessions.Register(sessionName, sessionStore))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth := r.Group("/oauth")
	{
		oauth.GET("/vendors", auth.OAuthAuthUrls)
		oauth.GET("/callback/kakao", auth.KakaoOAuthCallback)
		oauth.GET("/callback/naver", auth.NaverOAuthCallback)
		oauth.GET("/callback/google", auth.GoogleOAuthCallback)
	}

	v1 := r.Group("/v1")
	v1.Use(auth.Authenticate())
	// TODO: v1 API 라우트 추가
	// {
	// }

	debug := r.Group("/debug")
	{
		debug.POST("/redis", func(c *gin.Context) {
			redisClient := redis.Client()
			val := c.Query("val")
			key := c.Query("key")
			cmd := redisClient.Set(key, val, 0)
			c.JSON(200, gin.H{
				"message": cmd.String(),
			})
		})
		debug.GET("/redis", func(c *gin.Context) {
			redisClient := redis.Client()
			key := c.Query("key")
			c.JSON(200, gin.H{
				"message": redisClient.Get(key).String(),
			})
		})
	}

	return r
}
