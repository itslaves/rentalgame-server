package main

import (
	"github.com/skyoo2003/rentalgames-server/cmd"
	_ "github.com/skyoo2003/rentalgames-server/internal/docs"
	"github.com/skyoo2003/rentalgames-server/internal/third_party/mysql"
	"github.com/skyoo2003/rentalgames-server/internal/third_party/redis"
)

// @title RentalGames Swagger API
// @version 1.0
// @description RentalGames API specification
// @termsOfService http://swagger.io/terms/
// @BasePath /v1
func main() {
	if err := redis.Init(); err != nil {
		panic(err)
	}
	defer redis.Close()
	if err := mysql.Init(); err != nil {
		panic(err)
	}
	defer mysql.Close()

	cmd.Execute()
}
