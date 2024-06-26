package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/auth"
)

func main() {
	router := gin.Default()

	router.POST("/login", auth.Login)

	router.Run(":8080")
}
