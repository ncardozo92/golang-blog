package comment

import (
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/entity"
)

func ToDTO(c entity.Comment) dto.CommentDTO {
	return dto.CommentDTO{
		Id:      c.Id,
		IdPost:  c.IdPost,
		IdUser:  c.IdUser,
		Content: c.Content,
	}
}

func FromDTO(dto dto.CommentDTO) entity.Comment {
	return entity.Comment{
		Id:      dto.Id,
		IdPost:  dto.IdPost,
		IdUser:  dto.IdUser,
		Content: dto.Content,
	}
}
