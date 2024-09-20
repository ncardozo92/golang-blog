package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/auth"
	"github.com/ncardozo92/golang-blog/utils"
)

// IsAuthorized validates the JWT provided at Authorization header
// and verify the signature method
func IsAuthorized() gin.HandlerFunc {
	return func(context *gin.Context) {

		if !excludeFromValidation(context.FullPath()) {
			authToken := context.GetHeader("Authorization")

			// if the token is not sent, then the response should be Forbidden
			if authToken == "" {
				utils.BuildError(
					context,
					errors.New("forbidden request"),
					http.StatusForbidden,
					`Header "Authorization" no suministrado`)
			} else {
				JWTValidationErr := utils.ValidateJWT(authToken)

				// if the token is sent but there is an error
				if JWTValidationErr != nil {
					utils.BuildError(
						context,
						JWTValidationErr,
						http.StatusUnauthorized,
						"El JWT provisto no es válido") // acá puede fallar por invalid signature
				} else {
					context.Next()
				}
			}
		} else {
			context.Next()
		}
	}
}

func excludeFromValidation(path string) bool {
	return strings.EqualFold(path, auth.LOGIN_PATH)
}
