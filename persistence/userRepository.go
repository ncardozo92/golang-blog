package persistence

import "github.com/ncardozo92/golang-blog/entity"

type UserRepository interface {
	FindUserByUsername(username string) (entity.User, error)
}
