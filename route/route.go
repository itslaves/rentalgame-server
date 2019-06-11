package route

import (
	"encoding/json"
	"fmt"
	"github.com/itslaves/rentalgames-server/article"
	"github.com/itslaves/rentalgames-server/auth"
	"github.com/itslaves/rentalgames-server/common/redis"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	SessionKey    = "sessionKey"
	SessionSecret = "sessionSecret"
)

func Route() {
	r := gin.Default()

	err := redis.Init()
	if err != nil {
		panic(err)
	}

	defer redis.Close()

	store := cookie.NewStore([]byte(SessionSecret))
	r.Use(sessions.Sessions(SessionKey, store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/login", auth.LoginView)
	r.GET("/auth/callback", auth.Authenticate)

	r.GET("logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
	})

	r.GET("/kakao", func(context *gin.Context) {
		context.HTML(http.StatusOK, "kakao.tmpl", nil)
	})

	r.GET("/kakao-user", func(context *gin.Context) {
		token := context.Query("accessToken")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded;charset=utf-8")
		res, _ := client.Do(req)
		userInfo, err := ioutil.ReadAll(res.Body)

		if err != nil {
			context.JSON(http.StatusInternalServerError, nil)
			return
		}

		var u map[string]map[string]interface{}
		json.Unmarshal(userInfo, &u)

		fmt.Println(u)

		context.JSON(http.StatusOK, gin.H{
			"nickname": u["properties"]["nickname"],
		})
	})

	r.GET("/index", func(c *gin.Context) {
		user := auth.CurrentUser(c)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"name":  user.Name,
			"email": user.Email,
		})
	})

	r.GET("/articles", article.Retrieve)
	r.POST("/articles", article.Create)
	r.PUT("/articles/:id", article.Update)

	r.POST("/redis", func(c *gin.Context) {
		val := c.Query("val")
		key := c.Query("key")

		redisClient := redis.Client()
		cmd := redisClient.Set(key, val, 0)

		c.JSON(200, gin.H{
			"message": cmd.String(),
		})
	})

	r.GET("/redis", func(c *gin.Context) {
		key := c.Query("key")

		redisClient := redis.Client()

		c.JSON(200, gin.H{
			"message": redisClient.Get(key).String(),
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
