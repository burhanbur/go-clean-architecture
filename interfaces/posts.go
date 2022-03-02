package interfaces

import "blog/models"

type IPostService interface {
	CreatePost(Post *models.Posts) error
	GetPosts(offset, limit string) (interface{}, error)
	GetPostById(id string) (*models.Posts, error)
	UpdatePost(Post *models.Posts) error
	DeletePost(id string) error
}

type IPostRepository interface {
	CreatePost(Post *models.Posts) error
	GetPosts() ([]*models.Posts, error)
	GetPostPaginations(offset, limit string) (*models.PostsPagination, error)
	GetPostById(id string) (*models.Posts, error)
	UpdatePost(Post *models.Posts) error
	DeletePost(id string) error
}
