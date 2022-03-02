package config

import (
	"fmt"
	"os"
	"blog/utils"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		utils.Log{}.Error(err.Error())
	} else {
		utils.Log{}.Info("Redis successfully connected!")
	}

	return client
}
