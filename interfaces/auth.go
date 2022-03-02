package interfaces

import "blog/models"

type IAuthService interface {
	Login(auth *models.Auth) (*models.AuthResponse, error)
	Logout(logout *models.Logout) error
	ValidateUser(params map[string]string) (map[string]string, error)
}

type IAuthRepository interface {
	CheckUser(params string) (*models.Auth, error)
}
