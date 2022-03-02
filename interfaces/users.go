package interfaces

import "blog/models"

type IUserService interface {
	CreateUser(User *models.Users) error
	GetUsers(offset, limit string) (interface{}, error)
	GetUserById(id string) (*models.Users, error)
	UpdateUser(User *models.Users) error
	DeleteUser(id string) error
}

type IUserRepository interface {
	CreateUser(User *models.Users) error
	GetUsers() ([]*models.Users, error)
	GetUserPaginations(offset, limit string) (*models.UsersPagination, error)
	GetUserById(id string) (*models.Users, error)
	UpdateUser(User *models.Users) error
	DeleteUser(id string) error
}
