package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/persistence"
	"github.com/ncardozo92/golang-blog/persistence/relational"
	"github.com/ncardozo92/golang-blog/utils"
)

var postRepository persistence.PostRepository = relational.PostRepositorySQL{}

func GetAllTags(context *gin.Context) {
	tags, getAllTagsErr := postRepository.GetAllTags()

	if getAllTagsErr != nil {
		utils.BuildError(context, getAllTagsErr, http.StatusInternalServerError, "Hubo un problema al recuperar las etiquetas")
		return
	}

	context.JSON(http.StatusOK, tags)
}

func GetAllPosts(context *gin.Context) {
	posts, postsErr := postRepository.GetAllPosts()

	if postsErr != nil {
		utils.BuildError(context, postsErr, http.StatusInternalServerError,
			"No se pudo recuperar el listado de etiquetas")
		return
	}

	context.JSON(http.StatusOK, posts)
}
