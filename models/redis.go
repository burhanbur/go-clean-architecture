package models

type RedisInfo struct {
	Token        string `json:"token"`
	CreatedToken string `json:"created_token"`
	ExpiredToken string `json:"expired_token"`
}
