package dto

type CommentDTO struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
	IdPost  int64  `json:"id_post"`
	IdUser  int64  `json:"id_user"`
}
