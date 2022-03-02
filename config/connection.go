package config

import (
	"blog/models"
	"blog/utils"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

func InitDB() *DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("APP_DATABASE_USER"),
		os.Getenv("APP_DATABASE_PASSWORD"),
		os.Getenv("APP_DATABASE_HOST"),
		os.Getenv("APP_DATABASE_PORT"),
		os.Getenv("APP_DATABASE_NAME"),
	)

	db, err := gorm.Open(os.Getenv("APP_DATABASE_DRIVER"), dataSource)

	if err != nil {
		utils.Log{}.Error(err.Error())

		panic(err)
	}

	utils.Log{}.Info("Database successfully connected!")

	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)
	db.DB().SetConnMaxLifetime(5)

	dbConn.SQL = db

	db.AutoMigrate(
		&models.Users{},
		&models.Posts{},
	)

	return dbConn
}
