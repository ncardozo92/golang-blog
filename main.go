package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/auth"
	"github.com/ncardozo92/golang-blog/post"
)

func main() {
	router := gin.Default()

	// Authentication
	router.POST("/login", auth.Login)

	// Posts and tags
	router.GET("/tags", post.GetAllTags)

	router.Run(":8080")
}
