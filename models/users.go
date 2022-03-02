package models

import "time"

type Users struct {
	Id          string    `json:"id" gorm:"type:varchar(36);primary_key"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UsersPagination struct {
	*Pagination
	Records []Users `json:"records"`
}

func UsersTable() string {
	return "users"
}
