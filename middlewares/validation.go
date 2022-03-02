package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"blog/models"
	"blog/redis"
	"blog/utils"

	"github.com/gin-gonic/gin"
)

type RedisClient struct {
	Client redis.IRedisClient
}

func (r RedisClient) VerifyToken(token, userId string) bool {
	redisKeys := make([]string, 0)

	keyToken := fmt.Sprintf("%v", userId)
	redisKeys = append(redisKeys, keyToken)

	for _, key := range redisKeys {
		log.Println(key)
		value, err := r.Client.RetrieveById(key)

		if err != nil {
			return false
		}

		if value != nil {
			var redisInfo models.RedisInfo
			castValue := value.(string)
			err := json.Unmarshal([]byte(castValue), &redisInfo)

			if err != nil {
				log.Println(err.Error())
				return false
			}

			if redisInfo.Token == token {
				return true
			}
		}
	}

	return false
}

func (r RedisClient) TokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Header["Authorization"]) == 0 && len(c.Request.Header["User-Id"]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized access",
			})
			c.Abort()

			return
		} else {
			token := c.Request.Header["Authorization"][0]
			userId := c.Request.Header["User-Id"][0]
			validity := r.VerifyToken(token, userId)

			if validity {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Unauthorized access",
				})
				c.Abort()

				return
			}
		}
	}
}

// ini gw yg bikin iseng
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		utils.ErrorOutput(c, err)
	}
}
