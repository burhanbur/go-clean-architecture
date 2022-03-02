package redis

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
}

func (r redisClient) RetrieveById(id string) (interface{}, error) {
	value, err := r.client.Get(id).Result()

	if err != nil {
		return nil, err
	}

	return value, err
}

func (r redisClient) SetValue(key string, value interface{}) error {
	serialize, _ := json.Marshal(value)
	duration, err := strconv.Atoi(os.Getenv("REDIS_DURATION"))

	if err != nil {
		return err
	}

	err = r.client.Set(key, string(serialize), time.Duration(duration)*time.Minute).Err()

	return err
}

func (r redisClient) DelKey(key string) error {
	return r.client.Del(key).Err()
}

func NewRedisClient(client *redis.Client) IRedisClient {
	return redisClient{client: client}
}
