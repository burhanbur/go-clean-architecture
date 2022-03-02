package services

import (
	"blog/interfaces"
	"blog/models"

	"github.com/google/uuid"
)

type PostService struct {
	PostRepository interfaces.IPostRepository
}

func (d PostService) CreatePost(post *models.Posts) error {
	post.Id = uuid.New().String()
	return d.PostRepository.CreatePost(post)
}

func (d PostService) GetPosts(offset, limit string) (interface{}, error) {
	if offset != "" && limit != "" {
		Posts, err := d.PostRepository.GetPostPaginations(offset, limit)
		if err != nil {
			return nil, err
		}

		return Posts, nil
	}

	Posts, err := d.PostRepository.GetPosts()
	if err != nil {
		return nil, err
	}

	return Posts, nil
}

func (d PostService) GetPostById(id string) (*models.Posts, error) {
	post, err := d.PostRepository.GetPostById(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (d PostService) UpdatePost(post *models.Posts) error {
	return d.PostRepository.UpdatePost(post)
}

func (d PostService) DeletePost(id string) error {
	return d.PostRepository.DeletePost(id)
}

func NewPostService(PostRepository interfaces.IPostRepository) interfaces.IPostService {
	return PostService{PostRepository: PostRepository}
}
