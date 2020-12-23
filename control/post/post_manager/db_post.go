package post_manager

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"my_blog/common/conf"
	e "my_blog/common/entity"
	"my_blog/common/models"
	"my_blog/control/user/user_manager"
)

var log conf.Logger

func init() {
	log = conf.Log
}

// todo limit call speed
func Add(data *e.Post, currentUserID int) error {
	sql := "INSERT INTO post (body, body_html, author_id) VALUES (:body, :body_html, :author_id)"
	maybeUnsafeHTML := markdown.ToHTML([]byte(data.Body), nil, nil)
	// clean danger chars
	bodyHtml := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)
	args := map[string]interface{}{
		"body":      data.Body,
		"body_html": string(bodyHtml),
		"author_id": currentUserID,
	}
	_, err := conf.DB.NamedExec(sql, args)
	if err != nil {
		log.Error("Db insert error.\nsql: %s\nerr: %s", sql, err)
	}
	return err
}

func List(userID int, offset, limit string) (gin.H, error) {
	var posts []models.Post
	var args []interface{}
	args = append(args, userID)
	sql := "SELECT * FROM post WHERE author_id = ?"
	limits := ""
	if offset != "" && limit != ""{
		limits = " LIMIT ?, ?"
		args = append(args, offset, limit)
	}
	if limits != "" {
		sql += limits
	}
	err := conf.DB.Select(&posts, sql, args...)
	if err != nil {
		log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
	}
	var userIDs []int
	for _, c := range posts {
		userIDs = append(userIDs, c.AuthorID)
	}
	users, _ := user_manager.GetUserByIDs(userIDs...)
	for i, c := range posts {
		for _, u := range users {
			if u.ID == c.AuthorID {
				posts[i].Author = gin.H{"name": u.Name, "uuid": u.Uuid}
				break
			}
		}
	}
	var cnt int
	sql = "SELECT COUNT(*) AS cnt FROM post WHERE author_id = ?"
	conf.DB.QueryRow(sql, userID).Scan(&cnt)
	return gin.H{"rows": posts, "total": cnt}, err
}

func Remove(postID int) error {
	sql := "UPDATE post SET deleted = 1 WHERE id = ?"
	_, err := conf.DB.Exec(sql, postID)
	if err != nil {
		log.Error("Db update error.\nsql: %s\nerr: %s", sql, err)
	}
	return err
}

func DisableComment(PostID, CommentID int) error {
	sql := "UPDATE comment SET disabled = 1 WHERE id = ?"
	_, err := conf.DB.Exec(sql, CommentID)
	if err != nil {
		log.Error("Db update error.\nsql: %s\nerr: %s", sql, err)
	}
	return err
}

func Get(PostID int) (models.Post, error) {
	var post models.Post
	sql := "SELECT * FROM post WHERE id = ?"
	err := conf.DB.Get(&post, sql, PostID)
	if err != nil {
		log.Error("Db insert error.\nsql: %s\nerr: %s", sql, err)
	}
	return post, err
}

func AddComment(data e.Comment, currentUserID, PostID int) error {
	_, err := Get(PostID)
	if err != nil {
		return err
	}
	sql := "INSERT INTO comment (body, body_html, author_id, post_id) VALUES (:body, :body_html, :author_id, :post_id)"
	maybeUnsafeHTML := markdown.ToHTML([]byte(data.Body), nil, nil)
	// clean danger chars
	bodyHtml := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)
	args := map[string]interface{}{
		"body":      data.Body,
		"body_html": string(bodyHtml),
		"author_id": currentUserID,
		"post_id":   PostID,
	}
	tx, err := conf.DB.Beginx()
	_, err1 := tx.NamedExec(sql, args)
	sql = "UPDATE post SET comment_count = comment_count + 1 WHERE id = ?"
	_, err2 := tx.Exec(sql, PostID)
	err3 := tx.Commit()
	if err1 != nil{
		log.Error("Db insert error.\nsql: %s\nerr1: %s", sql, err1)
	}
	if err2 != nil{
		log.Error("Db insert error.\nsql: %s\nerr1: %s", sql, err2)
	}
	if err3 != nil{
		log.Error("Db insert error.\nsql: %s\nerr1: %s", sql, err3)
	}
	return err
}

func ListComment(postID, userID int) ([]models.Comment, error) {
	var comments []models.Comment
	sql := "SELECT * FROM comment WHERE author_id = ? and post_id = ?"
	err := conf.DB.Select(&comments, sql, userID, postID)
	if err != nil {
		log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
	}
	var userIDs []int
	for _, c := range comments {
		userIDs = append(userIDs, c.AuthorID)
	}
	users, _ := user_manager.GetUserByIDs(userIDs...)
	for i, c := range comments {
		for _, u := range users {
			if u.ID == c.AuthorID {
				comments[i].Author = gin.H{"name": u.Name, "uuid": u.Uuid}
				break
			}
		}
	}
	return comments, err
}

