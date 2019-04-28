package main

import (
	"encoding/json"
	"fmt"
	"gin-sample/article"
	"gin-sample/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const (
	SessionKey = "sessionKey"
	SessionSecret = "sessionSecret"
)

func main() {
	r := gin.Default()

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
			"name": user.Name,
			"email": user.Email,
		})
	})

	r.GET("/articles", article.Retrieve)
	r.POST("/articles", article.Create)
	r.PUT("/articles/:id", article.Update)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}