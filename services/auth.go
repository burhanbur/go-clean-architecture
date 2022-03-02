package services

import (
	"blog/interfaces"
	"blog/models"
	"blog/redis"
	"blog/utils"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	cache "github.com/go-redis/redis"
)

type AuthService struct {
	AuthRepository interfaces.IAuthRepository
	redisClient    redis.IRedisClient
}

func (a AuthService) Login(auth *models.Auth) (*models.AuthResponse, error) {
	var response models.AuthResponse
	var params = map[string]string{}

	params["email"] = auth.Email
	params["password"] = auth.Password

	user, err := a.ValidateUser(params)

	if err != nil {
		return nil, err
	}

	response.Id = user["id"]
	response.Email = user["email"]
	response.Token = user["token"]

	redisKeys := make(map[string]string)
	keyToken := fmt.Sprintf("%v", user["id"])
	redisKeys[keyToken] = user["token"]

	currentTime := time.Now()
	duration, err := strconv.Atoi(os.Getenv("REDIS_DURATION"))
	if err != nil {
		return nil, err
	}

	for key, value := range redisKeys {
		redisExist := true
		_, err = a.redisClient.RetrieveById(key)

		if err != nil {
			if err == cache.Nil {
				redisExist = false
				expiredTime := currentTime.Local().Add(time.Hour*0 +
					time.Minute*time.Duration(duration) +
					time.Second*0)

				err = a.redisClient.SetValue(key, models.RedisInfo{
					Token:        value,
					CreatedToken: currentTime.Format("2006-01-02 15:04:05"),
					ExpiredToken: expiredTime.Format("2006-01-02 15:04:05"),
				})

				if err != nil {
					return nil, err
				}
			}
		}

		if redisExist {
			//delete existing redis
			err := a.redisClient.DelKey(key)
			if err != nil {
				return nil, err
			}

			//insert new token in redis
			err = a.redisClient.SetValue(key, models.RedisInfo{
				Token: value,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &response, err
}

func (a AuthService) Logout(logout *models.Logout) error {
	redisKeys := make([]string, 0)

	keyToken := fmt.Sprintf("%v", logout.Id)
	redisKeys = append(redisKeys, keyToken)

	for _, key := range redisKeys {
		err := a.redisClient.DelKey(key)
		if err != nil {
			return errors.New("Unauthorized access")
		}
	}

	return nil
}

func (a AuthService) ValidateUser(params map[string]string) (map[string]string, error) {
	var returnValue map[string]string
	returnValue = map[string]string{}

	user, err := a.AuthRepository.CheckUser(params["email"])
	if err != nil {
		return nil, err
	}

	userId := user.Id
	email := user.Email
	password := user.Password

	pass := params["password"]

	// compare password as hash to pass as plain
	verify := utils.ComparePasswords(password, []byte(pass))

	if verify == false {
		err := errors.New("Invalid credentials")
		return nil, err
	}

	var data map[string]string
	data = map[string]string{}

	data["id"] = userId
	data["email"] = email

	token, err := GenerateToken(data)

	if err != nil {
		return nil, err
	}

	returnValue["id"] = userId
	returnValue["email"] = email
	returnValue["token"] = token

	return returnValue, nil
}

func GenerateToken(params map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	duration, err := strconv.Atoi(os.Getenv("REDIS_DURATION"))
	if err != nil {
		return "", err
	}

	claims["id"] = params["id"]
	claims["email"] = params["email"]
	claims["hit"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour*0 + time.Minute*time.Duration(duration) + time.Second*0).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewAuthService(AuthRepository interfaces.IAuthRepository, redisClient redis.IRedisClient) interfaces.IAuthService {
	return &AuthService{AuthRepository: AuthRepository, redisClient: redisClient}
}
