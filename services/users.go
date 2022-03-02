package services

import (
	"blog/interfaces"
	"blog/models"
	"blog/utils"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepository interfaces.IUserRepository
}

func (d UserService) CreateUser(User *models.Users) error {
	var password string
	User.Id = uuid.New().String()
	password = utils.HashAndSalt([]byte(User.Password))
	User.Password = password
	return d.UserRepository.CreateUser(User)
}

func (d UserService) GetUsers(offset, limit string) (interface{}, error) {
	if offset != "" && limit != "" {
		Users, err := d.UserRepository.GetUserPaginations(offset, limit)
		if err != nil {
			return nil, err
		}

		return Users, nil
	}

	Users, err := d.UserRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	return Users, nil
}

func (d UserService) GetUserById(id string) (*models.Users, error) {
	User, err := d.UserRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return User, nil
}

func (d UserService) UpdateUser(User *models.Users) error {
	if User.Password != "" {
		var password string
		password = utils.HashAndSalt([]byte(User.Password))
		User.Password = password
	}

	return d.UserRepository.UpdateUser(User)
}

func (d UserService) DeleteUser(id string) error {
	return d.UserRepository.DeleteUser(id)
}

func NewUserService(UserRepository interfaces.IUserRepository) interfaces.IUserService {
	return UserService{UserRepository: UserRepository}
}
