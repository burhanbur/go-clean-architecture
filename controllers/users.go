package controllers

import (
	"blog/interfaces"
	"blog/models"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Users struct {
	Service interfaces.IUserService
}

func (r Users) Users(route *gin.RouterGroup) {
	route.POST(utils.CREATE_USER, r.CreateUser)
	route.GET(utils.READ_USERS, r.ReadUsers)
	route.GET(utils.READ_USER_BY_ID, r.ReadUserById)
	route.PUT(utils.UPDATE_USER, r.UpdateUser)
	route.DELETE(utils.DELETE_USER, r.DeleteUser)
}

func (r Users) CreateUser(context *gin.Context) {
	var User models.Users
	err := context.ShouldBindJSON(&User)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	err = r.Service.CreateUser(&User)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    User,
		},
	)
}

func (r Users) ReadUsers(context *gin.Context) {
	offset := context.Query("page")
	limit := context.Query("limit")

	Users, err := r.Service.GetUsers(offset, limit)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	if offset != "" && limit != "" {
		context.JSON(
			http.StatusOK, gin.H{
				"status":  false,
				"message": "success",
				"data":    Users.(*models.UsersPagination),
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    Users.([]*models.Users),
		},
	)
}

func (r Users) ReadUserById(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "bad request",
		})
		return
	}

	User, err := r.Service.GetUserById(id)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    User,
		})
}

func (r Users) UpdateUser(context *gin.Context) {
	var User models.Users

	err := context.ShouldBindJSON(&User)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	if User.Id == "" {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  false,
				"message": "id not found",
			})
		return
	}

	err = r.Service.UpdateUser(&User)
	if err != nil {
		if err.Error() == "Department data not found" {
			context.JSON(http.StatusNoContent, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "success",
		"data":    User,
	})
}

func (r Users) DeleteUser(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "bad request",
		})
		return
	}

	err := r.Service.DeleteUser(id)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "success delete User",
	})
}
