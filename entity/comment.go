package entity

type Comment struct {
	Id      int64
	Content string
	IdPost  int64
	IdUser  int64
}
