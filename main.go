package main

import (
	"log"
	"net/http"
	"os"

	"blog/config"
	"blog/controllers"
	"blog/middlewares"
	"blog/redis"
	"blog/repositories"
	"blog/services"
	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	start := "Initializing Server ..."
	utils.Log{}.Info(start)
	err := godotenv.Load()

	if err != nil {
		utils.Log{}.Error("Loading .env file not found")
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(utils.CORSMiddleware())

	// handler for if no route available
	r.NoRoute(func(c *gin.Context) {
		url := c.Request.Host + c.Request.URL.Path

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  false,
				"message": "route not found",
				"url":     url,
			},
		)
	})

	db := config.InitDB()
	redisClient := config.InitRedis()

	v1 := r.Group("/api")

	// error handling
	v1.Use(middlewares.ErrorHandle())

	// redis
	redisRepo := redis.NewRedisClient(redisClient)
	redisMiddleware := middlewares.RedisClient{Client: redisRepo}
	v1.Use(redisMiddleware.TokenValidation())

	// auth
	IAuthRepository := repositories.NewAuthRepository(db)
	IAuthService := services.NewAuthService(IAuthRepository, redisRepo)
	AuthController := controllers.Auth{Service: IAuthService}
	AuthController.Auth(r)

	//user
	IUserRepository := repositories.NewUserRepository(db)
	IUserService := services.NewUserService(IUserRepository)
	UserController := controllers.Users{Service: IUserService}
	UserController.Users(v1)

	// post
	IPostRepository := repositories.NewPostRepository(db)
	IPostService := services.NewPostService(IPostRepository)
	PostController := controllers.Posts{Service: IPostService}
	PostController.Posts(v1)

	gin.SetMode(gin.ReleaseMode)

	connected := "Connected to port " + os.Getenv("APP_PORT")
	utils.Log{}.Info(connected)
	log.Fatal(r.Run(os.Getenv("APP_PORT")))
}
