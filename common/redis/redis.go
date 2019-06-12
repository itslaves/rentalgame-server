package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var redisClient redis.UniversalClient

func Init() error {
	if redisClient != nil {
		return errors.New("RedisClient has already been initialized.")
	}

	// https://godoc.org/github.com/go-redis/redis#UniversalOptions
	options := redis.UniversalOptions{
		Addrs: viper.GetStringSlice("redis.addrs"),
		Password: viper.GetString("redis.password"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	}

	// masterName이 있을 경우는 sentinel 모드로 동작
	if v := viper.GetString("redis.masterName"); v != "" {
		options.MasterName = v
	}

	redisClient = redis.NewUniversalClient(&options)

	if _, err := redisClient.Ping().Result(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return redisClient.Close()
}

func Client() redis.UniversalClient {
	return redisClient
}