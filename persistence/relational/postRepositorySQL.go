package relational

import (
	"database/sql"
	"errors"

	"github.com/ncardozo92/golang-blog/entity"
)

const query_delete_tags string = "DELETE FROM post_tag WHERE id_post = ?"

type PostRepositorySQL struct{}

func (repository PostRepositorySQL) GetAllPosts() ([]entity.Post, error) {

	postsList := []entity.Post{}

	// first we recover all posts
	postRows, postsQueryErr := getDatabase().Query("SELECT * FROM post")

	if postsQueryErr != nil {
		return nil, postsQueryErr
	}

	for postRows.Next() {
		post := entity.Post{}

		postScanErr := postRows.Scan(&post.Id, &post.Title, &post.Body, &post.Author)

		if postScanErr != nil {
			return nil, postScanErr
		}

		postsList = append(postsList, post)
	}

	// second we recover all tags to every post

	for i := 0; i < len(postsList); i++ {
		tagsIdsErr := repository.getAllPostTags(&postsList[i])

		if tagsIdsErr != nil {
			return postsList, tagsIdsErr
		}
	}

	return postsList, nil
}

func (repository PostRepositorySQL) GetById(id int64) (entity.Post, error) {

	foundPost := entity.Post{}

	row := getDatabase().QueryRow("SELECT * FROM post WHERE id = ?", id)

	postScanErr := row.Scan(&foundPost.Id, &foundPost.Title, &foundPost.Body, &foundPost.Author)

	if postScanErr != nil {
		return foundPost, postScanErr
	}

	tagsErr := repository.getAllPostTags(&foundPost)

	if tagsErr != nil {
		return foundPost, tagsErr
	}

	return foundPost, nil
}

func (repository PostRepositorySQL) getAllPostTags(post *entity.Post) error {
	tagsIds := []int64{}

	tagsRows, tagsQueryErr := getDatabase().Query("SELECT id_tag FROM post_tag WHERE id_post = ?", post.Id)

	if tagsQueryErr != nil {
		return tagsQueryErr
	}

	for tagsRows.Next() {
		var tag int64
		tagScanErr := tagsRows.Scan(&tag)

		if tagScanErr != nil {
			return tagScanErr
		}
		tagsIds = append(tagsIds, tag)
	}

	post.Tags = tagsIds

	return nil
}

func (repository PostRepositorySQL) CreatePost(post entity.Post) error {

	sqlStatement := "INSERT INTO post(title, body, author) VALUES(?,?,?)"

	// We do this because we need to use a transaction to prevent incosistencies
	transaction, getTransactionErr := getDatabase().Begin()

	if getTransactionErr != nil {
		return getTransactionErr
	}

	insertResult, insertErr := transaction.Exec(sqlStatement, post.Title, post.Body, post.Author)

	if insertErr != nil {
		transaction.Rollback()
		return insertErr
	}

	idInsertedPost, getLastInsertIdErr := insertResult.LastInsertId()

	if getLastInsertIdErr != nil {
		transaction.Rollback()
		return getLastInsertIdErr
	}

	// associating all tags to the post
	TagsAssociationErr := repository.associateTags(transaction, idInsertedPost, post.Tags)

	if TagsAssociationErr != nil {
		transaction.Rollback()
		return TagsAssociationErr
	}

	transactionCommitErr := transaction.Commit()

	if transactionCommitErr != nil {
		transaction.Rollback()
		return transactionCommitErr
	}

	return nil
}

func (repository PostRepositorySQL) associateTags(transaction *sql.Tx, idPost int64, tagsIds []int64) error {
	// If we are updating a post, then we need to disassociate the actual tags it may have
	_, deleteErr := transaction.Exec(query_delete_tags, idPost)

	if deleteErr != nil {
		return deleteErr
	}

	// Now we do the association
	for _, tagId := range tagsIds {
		_, insertErr := transaction.Exec("INSERT INTO post_tag(id_post, id_tag) VALUES(?,?)", idPost, tagId)

		if insertErr != nil {
			return insertErr
		}
	}

	return nil
}

func (repository PostRepositorySQL) GetAllTags() ([]entity.Tag, error) {

	tagsList := []entity.Tag{}

	tagsRows, tagQueryErr := getDatabase().Query("SELECT * FROM tag")

	if tagQueryErr != nil {
		return nil, tagQueryErr
	}

	for tagsRows.Next() {
		tag := entity.Tag{}

		tagScanErr := tagsRows.Scan(&tag.Id, &tag.Description)

		if tagScanErr != nil {
			return nil, tagScanErr
		}

		tagsList = append(tagsList, tag)
	}
	return tagsList, nil
}

func (repository PostRepositorySQL) UpdatePost(id int64, updatedData entity.Post) (bool, error) {

	// We start a transaction to update the data
	transaction, getTransactionErr := getDatabase().Begin()

	if getTransactionErr != nil {
		return false, getTransactionErr
	}
	// We realize the update on the post table
	updateResult, updateErr := transaction.Exec("UPDATE post SET "+
		"title = ?, "+"body = ? "+
		"WHERE id = ?", updatedData.Title, updatedData.Body, id)

	if updateErr != nil {
		transaction.Rollback()
		return false, updateErr
	}

	rowsAffected, getRowsAffectedErr := updateResult.RowsAffected()

	if getRowsAffectedErr != nil {
		return false, getRowsAffectedErr
	}

	if rowsAffected < 1 {
		return false, errors.New("post Does not Exists")
	}
	// We re assing the posts
	tagsUpdatingErr := repository.associateTags(transaction, id, updatedData.Tags)

	if tagsUpdatingErr != nil {
		return true, tagsUpdatingErr
	}

	commitErr := transaction.Commit()

	if commitErr != nil {
		return true, commitErr
	}

	return true, nil
}

func (repository PostRepositorySQL) Delete(id int64) (bool, error) {

	transaction, getTransactionErr := getDatabase().Begin()

	if getTransactionErr != nil {
		return false, getTransactionErr
	}

	// first we delete all tags associated to the post
	_, deleteTagsErr := transaction.Exec(query_delete_tags, id)

	if deleteTagsErr != nil {
		transaction.Rollback()
		return false, deleteTagsErr
	}

	// we delete all comments associated to the post
	_, deleteCommentsErr := transaction.Exec("DELETE FROM comment WHERE id_post = ?", id)

	if deleteCommentsErr != nil {
		transaction.Rollback()
		return false, deleteCommentsErr
	}

	// we delete the post
	deletePostResult, deletePostErr := transaction.Exec("DELETE FROM post WHERE id = ?", id)

	if deletePostErr != nil {
		transaction.Rollback()
		return false, deletePostErr
	}

	deletedPosts, getDeletedPostErr := deletePostResult.RowsAffected()

	if getDeletedPostErr != nil {
		transaction.Rollback()
		return true, getDeletedPostErr
	}

	if deletedPosts < 0 {
		transaction.Rollback()
		return false, nil
	}

	transaction.Commit()
	return true, nil
}
