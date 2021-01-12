package models

import "github.com/gin-gonic/gin"

type Post struct {
	ID           int       `db:"id" json:"id"`
	Body         string    `db:"body" json:"-"`
	BodyHTML     string    `db:"body_html" json:"body"`
	AuthorID     int       `db:"author_id" json:"-"`
	Author       gin.H     `db:"-" json:"author"`
	Comments     []Comment `db:"-" json:"comments"`
	CreatedAt    JSONTime  `db:"created_at" json:"created_at"`
	CommentCount int       `db:"comment_count" json:"comment_count"`
	Deleted      int       `db:"deleted" json:"deleted"`
	Url			 string 	   `db:"-" json:"url"`
}

type Comment struct {
	ID        int      `db:"id" json:"id"`
	Body      string   `db:"body" json:"body"`
	BodyHTML  string   `db:"body_html" json:"body_html"`
	AuthorID  int      `db:"author_id" json:"-"`
	Author    gin.H    `db:"-" json:"author"`
	Disabled  bool     `db:"disabled" json:"disabled"`
	PostID    int      `db:"post_id" json:"post_id"`
	CreatedAt JSONTime `db:"created_at" json:"created_at"`
}
