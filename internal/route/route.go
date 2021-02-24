package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/skyoo2003/rentalgames-server/internal/auth"
	"github.com/skyoo2003/rentalgames-server/internal/sessions"
	"github.com/skyoo2003/rentalgames-server/internal/third_party/redis"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	oauth := r.Group("/oauth")
	{
		oauth.GET("/vendors", auth.OAuthVendorUrls)
		oauth.GET("/callback/kakao", auth.KakaoOAuthCallback)
		oauth.GET("/callback/naver", auth.NaverOAuthCallback)
		oauth.GET("/callback/google", auth.GoogleOAuthCallback)
	}

	v1 := r.Group("/v1")
	v1.Use(auth.Authenticate())
	// TODO: v1 API 라우트 추가
	{
		v1.GET("/authentication", func(c *gin.Context) {
			c.JSON(http.StatusOK, nil)
		})
	}

	debug := r.Group("/debug")
	{
		debug.POST("/redis", func(c *gin.Context) {
			redisClient := redis.Client()
			val := c.Query("val")
			key := c.Query("key")
			cmd := redisClient.Set(key, val, 0)
			c.JSON(http.StatusOK, gin.H{
				"message": cmd.String(),
			})
		})
		debug.GET("/redis", func(c *gin.Context) {
			redisClient := redis.Client()
			key := c.Query("key")
			c.JSON(http.StatusOK, gin.H{
				"message": redisClient.Get(key).String(),
			})
		})
	}

	return r
}
