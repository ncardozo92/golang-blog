package relational

import (
	"github.com/ncardozo92/golang-blog/entity"
)

type UserRepositoryImpl struct {
}

func (repository UserRepositoryImpl) FindUserByUsername(username string) (entity.User, error) {

	foundUser := entity.User{}

	userData := getDatabase().QueryRow("SELECT * FROM blog_user WHERE username = ?", username)
	err := userData.Scan(&foundUser.Id, &foundUser.Username, &foundUser.Password)

	return foundUser, err
}
