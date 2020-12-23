package entity

type Post struct {
	Body      string       `json:"body" binding:"required"`
}

type Comment struct {
	Body      string      `json:"body" binding:"required"`
}
