package main

import (
	"github.com/itslaves/rentalgames-server/cmd"
	"github.com/itslaves/rentalgames-server/common/mysql"
	"github.com/itslaves/rentalgames-server/common/redis"
)

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
