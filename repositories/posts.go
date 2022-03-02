package repositories

import (
	"blog/config"
	"blog/interfaces"
	"blog/models"
	"blog/utils"
	"strconv"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func (d PostRepository) CreatePost(post *models.Posts) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.PostsTable()).Create(post).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func (d PostRepository) GetPosts() ([]*models.Posts, error) {
	Posts := make([]*models.Posts, 0)
	err := d.db.Debug().Table(models.PostsTable()).Find(&Posts).Error
	if err != nil {
		return nil, err
	}

	return Posts, nil
}

func (d PostRepository) GetPostPaginations(offset, limit string) (*models.PostsPagination, error) {
	var PostPagination models.PostsPagination
	Post := make([]*models.Posts, 0)
	var pagination models.Pagination
	var count int
	offsetLimit, err := utils.OffsetLimit(offset, limit)
	if err != nil {
		return nil, err
	}

	newOffset, _ := strconv.Atoi(offset)
	err = d.db.Debug().Model(models.Posts{}).Limit(offsetLimit["limit"]).Offset(offsetLimit["offset"]).Scan(&PostPagination.Records).Error
	if err != nil {
		return nil, err
	}

	err = d.db.Debug().Find(&Post).Count(&count).Error
	if err != nil {
		return nil, err
	}

	p, err := pagination.PageData(newOffset, offsetLimit["limit"], count)
	if err != nil {
		return nil, err
	}

	PostPagination.Pagination = p

	return &PostPagination, nil
}

func (d PostRepository) GetPostById(id string) (*models.Posts, error) {
	var Post models.Posts
	err := d.db.Debug().Table(models.PostsTable()).Where("id = ?", id).Find(&Post).Error
	if err != nil {
		return nil, err
	}

	return &Post, nil
}

func (d PostRepository) UpdatePost(post *models.Posts) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.PostsTable()).Where("id = ?", post.Id).Update(post).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func (d PostRepository) DeletePost(id string) error {
	tx := d.db.Begin()

	err := tx.Debug().Table(models.PostsTable()).Where("id = ?", id).Delete(&models.Posts{}).Error
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func NewPostRepository(db *config.DB) interfaces.IPostRepository {
	return PostRepository{db: db.SQL}
}
