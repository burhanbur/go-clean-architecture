package repositories

import (
	"blog/config"
	"blog/interfaces"
	"blog/models"
	"blog/utils"
	"strconv"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (d UserRepository) CreateUser(User *models.Users) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.UsersTable()).Create(User).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func (d UserRepository) GetUsers() ([]*models.Users, error) {
	Users := make([]*models.Users, 0)
	err := d.db.Debug().Table(models.UsersTable()).Find(&Users).Error
	if err != nil {
		return nil, err
	}

	return Users, nil
}

func (d UserRepository) GetUserPaginations(offset, limit string) (*models.UsersPagination, error) {
	var UserPagination models.UsersPagination
	User := make([]*models.Users, 0)
	var pagination models.Pagination
	var count int
	offsetLimit, err := utils.OffsetLimit(offset, limit)
	if err != nil {
		return nil, err
	}

	newOffset, _ := strconv.Atoi(offset)
	err = d.db.Debug().Model(models.Users{}).Limit(offsetLimit["limit"]).Offset(offsetLimit["offset"]).Scan(&UserPagination.Records).Error
	if err != nil {
		return nil, err
	}

	err = d.db.Debug().Find(&User).Count(&count).Error
	if err != nil {
		return nil, err
	}

	p, err := pagination.PageData(newOffset, offsetLimit["limit"], count)
	if err != nil {
		return nil, err
	}

	UserPagination.Pagination = p

	return &UserPagination, nil
}

func (d UserRepository) GetUserById(id string) (*models.Users, error) {
	var User models.Users
	err := d.db.Debug().Table(models.UsersTable()).Where("id = ?", id).Find(&User).Error
	if err != nil {
		return nil, err
	}

	return &User, nil
}

func (d UserRepository) UpdateUser(User *models.Users) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.UsersTable()).Where("id = ?", User.Id).Update(User).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func (d UserRepository) DeleteUser(id string) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.UsersTable()).Where("id = ?", id).Delete(&models.Users{}).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func NewUserRepository(db *config.DB) interfaces.IUserRepository {
	return UserRepository{db: db.SQL}
}
