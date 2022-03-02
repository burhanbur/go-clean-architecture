package controllers

import (
	"blog/interfaces"
	"blog/models"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	Service interfaces.IAuthService
}

func (a Auth) Auth(route *gin.Engine) {
	route.POST(utils.LOGIN, a.Login)
	route.POST(utils.LOGOUT, a.Logout)
}

func (a Auth) Login(context *gin.Context) {
	var auth models.Auth
	err := context.ShouldBindJSON(&auth)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := a.Service.Login(&auth)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

func (a Auth) Logout(context *gin.Context) {
	var logout models.Logout

	err := context.ShouldBindJSON(&logout)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	err = a.Service.Logout(&logout)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": true, "message": "Success logout"})
}
