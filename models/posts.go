package models

import "time"

type Posts struct {
	Id        string    `json:"id" gorm:"type:varchar(36);primary_key"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Slug      string    `json:"slug"`
	AuthorId  string    `json:"author_id"`
	Thumbnail string    `json:"thumbnail" gorm:"default:NULL"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostsPagination struct {
	*Pagination
	Records []Posts `json:"records"`
}

func PostsTable() string {
	return "posts"
}
