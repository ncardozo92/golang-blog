// package error offers a helper function to handle every error in the API
package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/dto"
)

// ValidateError takes the error to validate, if err != nil then it will send a HTTP response
// with the specified status code, the posible message and details
func BuildError(
	context *gin.Context,
	err error,
	statusCode int,
	errorMessage string) {

	var responseDTO dto.GenericErrorDTO
	responseDTO.Message = errorMessage
	// todo: analizar si vamos a usar los details
	//responseDTO.Details = details

	log.Println("Error:", err.Error())

	context.AbortWithStatusJSON(statusCode, responseDTO)

}
