package entity

type Post struct {
	Id     int64
	Title  string
	Body   string
	Author int64
	Tags   []int64
}
