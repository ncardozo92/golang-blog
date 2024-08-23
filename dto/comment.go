package dto

type CommentDTO struct {
	Id      int64
	Content string
	IdPost  int64
	IdUser  int64
}
