package post

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/entity"
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

// GetAllPosts returns all the posts in the database
func GetAllPosts(context *gin.Context) {
	response := []dto.PostDTO{}
	posts, postsErr := postRepository.GetAllPosts()

	if postsErr != nil {
		utils.BuildError(context, postsErr, http.StatusInternalServerError,
			"No se pudieron recuperar los posts")
		return
	}

	for _, post := range posts {
		response = append(response, ToDTO(post))
	}

	context.JSON(http.StatusOK, response)
}

func GetById(context *gin.Context) {

	idParam, okId := context.Params.Get("id")

	id, idConvertErr := strconv.ParseInt(idParam, 10, 64)

	if idConvertErr != nil || !okId {
		utils.BuildError(context, idConvertErr, http.StatusBadRequest, "El id debe ser un valor numérico")
		return
	}

	foundPost, queryPostErr := postRepository.GetById(id)

	if queryPostErr != nil {
		var statusCode int

		if isPresent(foundPost) {
			statusCode = http.StatusInternalServerError
		} else {
			statusCode = http.StatusNotFound
		}
		utils.BuildError(
			context,
			queryPostErr,
			statusCode,
			"Hubo un error al recuperar el post")
		return
	}

	context.JSON(http.StatusOK, ToDTO(foundPost))
}

func Create(context *gin.Context) {
	dto := dto.PostDTO{}
	DTOBindingErr := context.BindJSON(&dto)

	jwt := context.GetHeader("Authorization")

	if DTOBindingErr != nil {
		utils.BuildError(context, DTOBindingErr, http.StatusBadRequest, "El cuerpo de la solicitud no es válido")
		return
	}

	userId, userIdErr := utils.GetUserId(jwt)

	if userIdErr != nil {
		utils.BuildError(context, userIdErr, http.StatusBadRequest, "No se pudo crear el Post")
		return
	}

	dto.Author = userId

	createErr := postRepository.CreatePost(fromDTO(dto))

	if createErr != nil {
		utils.BuildError(context, createErr, http.StatusInternalServerError, "Hubo un error al crear el post")
		return
	}

	context.Status(http.StatusOK)
}

func Update(context *gin.Context) {
	var dto dto.PostDTO
	postId, postIdConvertErr := strconv.Atoi(context.Param("id"))

	if postIdConvertErr != nil {
		utils.BuildError(context, postIdConvertErr, http.StatusBadRequest, "El ID debe ser un valor numérico")
		return
	}

	DTOBindingErr := context.BindJSON(&dto)

	if DTOBindingErr != nil {
		utils.BuildError(context, DTOBindingErr, http.StatusBadRequest, "El cuerpo de la solicitud no es válido")
		return
	}

	existsPost, postUpdateErr := postRepository.UpdatePost(int64(postId), fromDTO(dto))

	if postUpdateErr != nil {
		var statusCode int
		var message string
		if !existsPost {
			statusCode = http.StatusNotFound
			message = "El post no existe"
		} else {
			statusCode = http.StatusInternalServerError
			message = "Hubo un error al actualizar los datos del post"
		}

		utils.BuildError(context, postUpdateErr, statusCode, message)
		return
	}

	context.Status(http.StatusOK)
}

func isPresent(post entity.Post) bool {
	return post.Id != 0
}

// Deletes a post by it´s ID
func Delete(context *gin.Context) {

	id, idConvertErr := strconv.Atoi(context.Param("id"))

	if idConvertErr != nil {
		utils.BuildError(context, idConvertErr, http.StatusBadRequest, "El ID debe ser un valor numérico")
		return
	}

	deleted, deleteErr := postRepository.Delete(int64(id))

	if deleteErr != nil {
		utils.BuildError(context, deleteErr, http.StatusInternalServerError, "Hubo un error al eliminar el post")
		return
	} else if !deleted {
		utils.BuildError(context,
			errors.New("post not found"),
			http.StatusInternalServerError, "No existe un Post con el id proporcionado")
	}
}
