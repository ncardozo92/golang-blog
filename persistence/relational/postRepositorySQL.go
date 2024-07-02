package relational

import "github.com/ncardozo92/golang-blog/entity"

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
	for _, post := range postsList {
		tagsIds := []int64{}

		tagsRows, tagsQueryErr := getDatabase().Query("SELECT id_tag FROM post_tag WHERE id_post = ?", post.Id)

		if tagsQueryErr != nil {
			return nil, tagsQueryErr
		}

		for tagsRows.Next() {
			var tag int64
			tagScanErr := tagsRows.Scan(&tag)

			if tagScanErr != nil {
				return nil, tagScanErr
			}
			tagsIds = append(tagsIds, tag)
		}

		post.Tags = tagsIds
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
