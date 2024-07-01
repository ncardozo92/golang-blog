package persistence

import (
	"github.com/ncardozo92/golang-blog/entity"
)

type PostRepository interface {
	GetAllPosts() ([]entity.Post, error)
	GetAllTags() ([]entity.Tag, error)
}
