package persistence

import "github.com/ncardozo92/golang-blog/entity"

type CommentRepository interface {
	GetAllByPostId(id int64) ([]entity.Comment, error)
	GetByUserId(id int64) ([]entity.Comment, error)
	Save(comment entity.Comment) error
	Delete(id int64) (bool, error)
}
