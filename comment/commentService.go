package comment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/persistence"
	"github.com/ncardozo92/golang-blog/persistence/relational"
	"github.com/ncardozo92/golang-blog/utils"
)

var commentRepository persistence.CommentRepository = relational.CommentRepositorySQL{}

const BY_USER_QUERY_PARAM_NAME string = "userId"

func GetAllByPostId(context *gin.Context) {
	postId, idConvertErr := strconv.Atoi(context.Param("postId"))
	dtoList := []dto.CommentDTO{}

	if idConvertErr != nil {
		// todo: implementar esta alternativa
		//context.AbortWithStatusJSON(400, ErrorDto)
		utils.BuildError(context, idConvertErr, http.StatusBadRequest, "El id debe ser un valor numérico")
		return
	}

	comments, getCommentsErr := commentRepository.GetAllByPostId(int64(postId))

	if getCommentsErr != nil {
		utils.BuildError(context, getCommentsErr,
			http.StatusInternalServerError, "Hubo un error al recuperar los comentarios del post")
		return
	}

	for _, comment := range comments {
		dtoList = append(dtoList, ToDTO(comment))
	}

	context.JSON(http.StatusOK, dtoList)
}

func GetByUser(context *gin.Context) {
	dtoList := []dto.CommentDTO{}
	authHeader := context.GetHeader("Authorization")
	userId, getUserIdErr := utils.GetUserId(strings.Replace(authHeader, "Bearer ", "", 1))

	if getUserIdErr != nil {
		utils.BuildError(context,
			getUserIdErr,
			http.StatusBadRequest,
			"No se pudo recuperar el ID de usuario desde el JWT")
		return
	}

	comments, getCommentsErr := commentRepository.GetByUserId(userId)

	if getCommentsErr != nil {
		utils.BuildError(context,
			getCommentsErr,
			http.StatusInternalServerError,
			"Hubo un error al recuperar los comentarios")
		return
	}

	for _, comment := range comments {
		dtoList = append(dtoList, ToDTO(comment))
	}

	context.JSON(200, dtoList)
}

func Create(context *gin.Context) {
	dto := dto.CommentDTO{}

	dtoBindingErr := context.ShouldBindJSON(&dto)

	if dtoBindingErr != nil {
		utils.BuildError(context, dtoBindingErr, http.StatusBadRequest, "El JSON enviado no es válido")
	}

	saveErr := commentRepository.Save(FromDTO(dto))

	if saveErr != nil {
		utils.BuildError(context, saveErr, http.StatusInternalServerError, "No se pudo guardar el comentario")
		return
	}
}

func Delete(context *gin.Context) {

	idComment, idCommentConvertErr := strconv.Atoi(context.Param("id"))

	if idCommentConvertErr != nil {
		utils.BuildError(context, idCommentConvertErr, http.StatusBadRequest, "El id debe ser un valor numérico")
		return
	}

	deleted, deleteErr := commentRepository.Delete(int64(idComment))

	if deleteErr != nil {
		utils.BuildError(context, deleteErr, http.StatusInternalServerError, "Hubo un error al eliminar el comentario")
		return
	}

	if !deleted {
		utils.BuildError(context,
			fmt.Errorf("comment not found"),
			http.StatusNotFound,
			"No se encontró el comentario con el ID especificado")
		return
	}

	context.Status(http.StatusOK)
}
