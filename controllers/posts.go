package controllers

import (
	"blog/interfaces"
	"blog/models"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Posts struct {
	Service interfaces.IPostService
}

func (r Posts) Posts(route *gin.RouterGroup) {
	route.POST(utils.CREATE_POST, r.CreatePost)
	route.GET(utils.READ_POSTS, r.ReadPosts)
	route.GET(utils.READ_POST_BY_ID, r.ReadPostById)
	route.PUT(utils.UPDATE_POST, r.UpdatePost)
	route.DELETE(utils.DELETE_POST, r.DeletePost)
}

func (r Posts) CreatePost(context *gin.Context) {
	var Post models.Posts
	err := context.ShouldBindJSON(&Post)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	err = r.Service.CreatePost(&Post)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    Post,
		},
	)
}

func (r Posts) ReadPosts(context *gin.Context) {
	offset := context.Query("page")
	limit := context.Query("limit")

	Posts, err := r.Service.GetPosts(offset, limit)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	if offset != "" && limit != "" {
		context.JSON(
			http.StatusOK, gin.H{
				"status":  false,
				"message": "success",
				"data":    Posts.(*models.PostsPagination),
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    Posts.([]*models.Posts),
		},
	)
}

func (r Posts) ReadPostById(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "bad request",
		})
		return
	}

	Post, err := r.Service.GetPostById(id)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(http.StatusOK,
		gin.H{
			"status":  false,
			"message": "success",
			"data":    Post,
		})
}

func (r Posts) UpdatePost(context *gin.Context) {
	var Post models.Posts

	err := context.ShouldBindJSON(&Post)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	if Post.Id == "" {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  false,
				"message": "id not found",
			})
		return
	}

	err = r.Service.UpdatePost(&Post)
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
		"data":    Post,
	})
}

func (r Posts) DeletePost(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "bad request",
		})
		return
	}

	err := r.Service.DeletePost(id)
	if err != nil {
		utils.ErrorOutput(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "success delete Post",
	})
}
