package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/persistence"
	"github.com/ncardozo92/golang-blog/persistence/relational"
	"github.com/ncardozo92/golang-blog/utils"
	"golang.org/x/crypto/bcrypt"
)

var secret string = os.Getenv("JWT_SECRET_TOKEN")

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

	// We generate the token with the library found

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":    "golang-blog",
			"sub":    user.Username,
			"userId": user.Id,
			"iat":    time.Now().Unix(),
			"exp":    time.Now().Unix() + (5 * 3600), // the token is valid for 5 minutes
		})

	token, jwtSigningErr := tokenGenerator.SignedString([]byte(secret))

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
