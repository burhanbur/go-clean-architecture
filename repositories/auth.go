package repositories

import (
	"blog/config"
	"blog/interfaces"
	"blog/models"

	"github.com/jinzhu/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (a AuthRepository) CheckUser(email string) (*models.Auth, error) {
	var users models.Auth

	err := a.db.Debug().Table(models.UsersTable()).First(&users, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func NewAuthRepository(db *config.DB) interfaces.IAuthRepository {
	return &AuthRepository{db: db.SQL}
}
