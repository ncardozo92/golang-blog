package relational

import (
	"database/sql"

	"github.com/ncardozo92/golang-blog/entity"
)

type CommentRepositorySQL struct{}

func (repository CommentRepositorySQL) GetAllByPostId(id int64) ([]entity.Comment, error) {

	queryResult, sqlQueryErr := getDatabase().Query("SELECT * FROM comment WHERE id_post = ?", id)

	if sqlQueryErr != nil {
		return nil, sqlQueryErr
	}

	return scanComments(queryResult)
}

func (repository CommentRepositorySQL) GetByUserId(id int64) ([]entity.Comment, error) {

	queryResult, queryErr := getDatabase().Query("SELECT * FROM comment WHERE id_user = ?", id)

	if queryErr != nil {
		return nil, queryErr
	}

	return scanComments(queryResult)
}

func scanComments(rows *sql.Rows) ([]entity.Comment, error) {
	comments := []entity.Comment{}

	for rows.Next() {
		comment := entity.Comment{}

		queryResultScanErr := rows.Scan(&comment.Id, &comment.Content, &comment.IdPost, &comment.IdUser)

		if queryResultScanErr != nil {
			return nil, queryResultScanErr
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
