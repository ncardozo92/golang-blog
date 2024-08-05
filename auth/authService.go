package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/persistence"
	"github.com/ncardozo92/golang-blog/persistence/relational"
	"github.com/ncardozo92/golang-blog/utils"
	"golang.org/x/crypto/bcrypt"
)

const LOGIN_PATH string = "/login"

var userRepository persistence.UserRepository = relational.UserRepositorySQL{}

func Login(context *gin.Context) {

	var requestDTO dto.LoginRequestDTO

	// We bind the request body to the DTO to manipulate the data
	requestBindingErr := context.ShouldBindJSON(&requestDTO)

	if requestBindingErr != nil {
		utils.BuildError(
			context,
			requestBindingErr,
			http.StatusBadRequest,
			"El request no es válido, por favor siga la documentación")
		return
	}

	// We find the user in the database by it's username
	user, findUserErr := userRepository.FindUserByUsername(requestDTO.Username)

	if findUserErr != nil {
		if findUserErr == sql.ErrNoRows {
			utils.BuildError(
				context,
				findUserErr,
				http.StatusNotFound,
				"usuario no encontrado")
			return
		}
		utils.BuildError(
			context,
			findUserErr,
			http.StatusInternalServerError,
			"Hubo un error al buscar el usuario")
		return
	}

	// After finding the user, we validate the password is correct
	bcryptCompareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestDTO.Password))

	if bcryptCompareErr != nil {
		utils.BuildError(
			context,
			bcryptCompareErr,
			http.StatusUnauthorized,
			"Password incorrecto")
		return
	}

	// We generate the JWT

	token, jwtSigningErr := utils.GenerateJWT(user)

	if jwtSigningErr != nil {
		utils.BuildError(
			context,
			jwtSigningErr,
			http.StatusInternalServerError,
			"No pudimos generar el token")
		return
	}

	// We make the response and send it to the client
	loginResponse := dto.LoginResponseDTO{Token: token}

	context.JSON(http.StatusOK, loginResponse)
}
