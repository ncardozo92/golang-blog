package dto

type PostDTO struct {
	Id     int64   `json:"id"`
	Title  string  `json:"title"`
	Body   string  `json:"body"`
	Author int64   `json:"author"`
	Tags   []int64 `json:"tags"`
}
