package relational

import (
	"database/sql"
	"fmt"

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

func (repository CommentRepositorySQL) Save(comment entity.Comment) error {
	db := getDatabase()

	result, insertErr := db.Exec("INSERT INTO comment(id_post,id_user,content) VALUES(?,?,?)",
		comment.IdPost, comment.IdUser, comment.Content)

	rowsAffected, getRowsAffectedErr := result.RowsAffected()

	if insertErr != nil {
		return insertErr
	} else if getRowsAffectedErr != nil {
		return getRowsAffectedErr
	} else if rowsAffected != 1 {
		return fmt.Errorf("comment not saved")
	}

	return nil
}
