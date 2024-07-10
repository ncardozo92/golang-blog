package relational

import (
	"github.com/ncardozo92/golang-blog/entity"
)

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
	/*
		for _, post := range postsList {
			tagsIds, tagsIdsErr := repository.getAllPostTags(post.Id)

			if tagsIdsErr != nil {
				return postsList, tagsIdsErr
			}

			// por acá debe estar el error de los tags
			post.Tags = tagsIds
		}
	*/

	for i := 0; i < len(postsList); i++ {
		tagsIds, tagsIdsErr := repository.getAllPostTags(postsList[i].Id)

		if tagsIdsErr != nil {
			return postsList, tagsIdsErr
		}

		// por acá debe estar el error de los tags
		postsList[i].Tags = tagsIds
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

	tagsIds, tagsErr := repository.getAllPostTags(foundPost.Id)

	if tagsErr != nil {
		return foundPost, tagsErr
	}

	foundPost.Tags = tagsIds

	return foundPost, nil
}

func (repository PostRepositorySQL) getAllPostTags(postId int64) ([]int64, error) {
	tagsIds := []int64{}

	tagsRows, tagsQueryErr := getDatabase().Query("SELECT id_tag FROM post_tag WHERE id_post = ?", postId)

	if tagsQueryErr != nil {
		return tagsIds, tagsQueryErr
	}

	for tagsRows.Next() {
		var tag int64
		tagScanErr := tagsRows.Scan(&tag)

		if tagScanErr != nil {
			return tagsIds, tagScanErr
		}
		tagsIds = append(tagsIds, tag)
	}

	return tagsIds, nil
}

func (repository PostRepositorySQL) CreatePost(post entity.Post) error {

	sqlStatement := "INSERT INTO post(title, body, author) VALUES(?,?,?)"

	// We do this we need to use a transaction to prevent incosistencies
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
	for _, tagId := range post.Tags {
		_, tagAssignErr := transaction.Exec("INSERT INTO post_tag(id_post, id_tag) VALUES(?,?)", idInsertedPost, tagId)

		if tagAssignErr != nil {
			transaction.Rollback()
			return tagAssignErr
		}
	}

	transactionCommitErr := transaction.Commit()

	if transactionCommitErr != nil {
		return transactionCommitErr
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
