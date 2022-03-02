package models

type Auth struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type VerifyToken struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type Logout struct {
	Id string `json:"id"`
}
