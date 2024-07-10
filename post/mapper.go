package post

import (
	"github.com/ncardozo92/golang-blog/dto"
	"github.com/ncardozo92/golang-blog/entity"
)

func ToDTO(entity entity.Post) dto.PostDTO {
	return dto.PostDTO{
		Id:     entity.Id,
		Title:  entity.Title,
		Body:   entity.Body,
		Author: entity.Author,
		Tags:   entity.Tags,
	}
}

func fromDTO(dto dto.PostDTO) entity.Post {
	return entity.Post{
		Id:     dto.Id,
		Title:  dto.Title,
		Author: dto.Author,
		Tags:   dto.Tags,
		Body:   dto.Body,
	}
}
