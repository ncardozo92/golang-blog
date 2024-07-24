package persistence

import (
	"github.com/ncardozo92/golang-blog/entity"
)

type PostRepository interface {
	GetAllPosts() ([]entity.Post, error)
	GetAllTags() ([]entity.Tag, error)
	GetById(id int64) (entity.Post, error)
	CreatePost(post entity.Post) error
	UpdatePost(id int64, post entity.Post) (bool, error)
	Delete(id int64) (bool, error)
}
