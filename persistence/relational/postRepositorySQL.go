package relational

import "github.com/ncardozo92/golang-blog/entity"

type PostRepositorySQL struct{}

func (repository PostRepositorySQL) GetAllPosts() ([]entity.Post, error) {

	postsList := []entity.Post{}

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

		// todo: recuperar los tags

		postsList = append(postsList, post)
	}

	return postsList, nil
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
