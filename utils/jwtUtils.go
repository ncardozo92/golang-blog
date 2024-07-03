package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ncardozo92/golang-blog/entity"
)

const bearerJWTPrefix string = "Bearer "

var secret string = os.Getenv("JWT_SECRET_TOKEN")

func GenerateJWT(user entity.User) (string, error) {

	now := time.Now().Unix()

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":    "golang-blog",
			"sub":    user.Username,
			"userId": user.Id,
			"iat":    now,
			"exp":    now + (5 * 3600), // the token is valid for 5 minutes
		})

	return tokenGenerator.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) error {
	_, tokenParsingErr := jwt.Parse(strings.Replace(tokenString, bearerJWTPrefix, "", 1), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("El método del JWT recibido no es el mismo que usamos para firmar los JWT")
		}

		return []byte(secret), nil
	})

	return tokenParsingErr
}
