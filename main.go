package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ncardozo92/golang-blog/auth"
	"github.com/ncardozo92/golang-blog/comment"
	"github.com/ncardozo92/golang-blog/middleware"
	"github.com/ncardozo92/golang-blog/post"
)

func main() {
	router := gin.Default()

	// Middlewares to use
	router.Use(middleware.IsAuthorized())

	// Authentication
	router.POST(auth.LOGIN_PATH, auth.Login)

	// Posts and tags
	router.GET("/tags", post.GetAllTags)
	router.GET("/posts", post.GetAllPosts)
	router.GET("/posts/:id", post.GetById)
	router.POST("/posts", post.Create)
	router.PUT("/posts/:id", post.Update)
	router.DELETE("/posts/:id", post.Delete)

	// Comments
	router.GET("/comments/:postId", comment.GetAllByPostId)
	router.GET("/comments", comment.GetByUser)
	router.POST("/comments", comment.Create)
	router.DELETE("/comments/:id", comment.Delete)

	router.Run(":8080")
}
