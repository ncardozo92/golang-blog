package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ncardozo92/golang-blog/entity"
)

const (
	bearerJWTPrefix string = "Bearer "
	userIdField     string = "userId"
)

type CustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

var secret string = os.Getenv("JWT_SECRET_TOKEN")

// GenerateJWT generates the token for the user
func GenerateJWT(user entity.User) (string, error) {

	now := time.Now().Unix()

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":       "golang-blog",
			"sub":       user.Username,
			userIdField: user.Id,
			"iat":       now,
			"exp":       now + (5 * 60), // the token is valid for 5 minutes
		})

	return tokenGenerator.SignedString([]byte(secret))
}

// Validates the signature method and the sing of the JWT
func ValidateJWT(tokenString string) error {
	_, err := parseJWT(tokenString)

	return err
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(strings.Replace(tokenString, bearerJWTPrefix, "", 1), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("el algoritmo del JWT recibido no es el mismo que usamos para firmar los JWT")
		}

		return []byte(secret), nil
	})
}

// returns the user ID from the JWT
func GetUserId(tokenString string) (int64, error) {
	// We recover the payload of the JWT
	encodedPayload := strings.Split(tokenString, ".")[1]
	rawPayload := make(map[string]any)

	// Now we decode the payload
	decodedPayload, decodingErr := base64.RawURLEncoding.DecodeString(encodedPayload)

	if decodingErr != nil {
		return 0, decodingErr
	}

	unmarshallErr := json.Unmarshal(decodedPayload, &rawPayload)

	if unmarshallErr != nil {
		return 0, unmarshallErr
	}

	userId, userIdOk := rawPayload[userIdField]

	if !userIdOk {
		return 0, fmt.Errorf("user ID is not present at JWT")
	}

	return int64(userId.(float64)), nil // In this line I convert an interface value to a int64 value
}
