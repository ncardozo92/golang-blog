package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/persistence"
	"github.com/ncardozo92/golang-blog/persistence/relational"
	"golang.org/x/crypto/bcrypt"
)

const secret string = "secret"

var userRepository persistence.UserRepository = relational.UserRepositoryImpl{}

func Login(context *gin.Context) {

	var requestDTO dto.LoginRequestDTO

	// We bind the request body to the DTO to manipulate the data
	requestBindingErr := context.ShouldBindJSON(&requestDTO)

	if requestBindingErr != nil {
		log.Println("Error al recuperar los datos del login:", requestBindingErr.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": "El request no es válido"})
		return
	}

	// We find the user in the database by it's username
	user, findUserErr := userRepository.FindUserByUsername(requestDTO.Username)

	if findUserErr != nil {
		log.Println("Error al recuperar el usuario:", findUserErr.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Hubo un error al momento de realizar el login"})
		return
	}

	// After finding the user, we validate the password is correct
	bcryptCompareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestDTO.Password))

	if bcryptCompareErr != nil {
		log.Println("Bcrypt Error:", bcryptCompareErr.Error())
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Password incorrecto"})
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
		log.Println("Bcrypt Error:", jwtSigningErr.Error())
		context.JSON(http.StatusUnauthorized, gin.H{"error": "No pudimos crear tu sesión"})
		return
	}

	// We make the response and send it to the client
	loginResponse := dto.LoginResponseDTO{Token: token}

	context.JSON(http.StatusOK, loginResponse)
}
