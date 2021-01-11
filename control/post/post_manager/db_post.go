package post_manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"my_blog/common/conf"
	e "my_blog/common/entity"
	"my_blog/common/models"
	"my_blog/control/post/define"
	"my_blog/control/user/user_manager"
	"time"
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
				posts[i].Url = fmt.Sprintf("/post/%d", c.ID)
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
		log.Error("Db select error.\nsql: %s\nerr: %s", sql, err)
	}
	author, err := user_manager.GetUserByID(post.AuthorID)
	if err != nil {
		return models.Post{}, err
	}
	post.Author = gin.H{"name": author.Name, "uuid": author.Uuid}
	post.Url = fmt.Sprintf("/post/%d", post.ID)
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

func CachePost(post *models.Post) {
	postKey := fmt.Sprintf("post:%d", post.ID)
	exist, err := conf.Redis.Exist(postKey)
	if err != nil {
		conf.Log.Error("Redis exist cmd error, err: %s", err)
		return
	}
	if exist {
		return
	}
	err = conf.Redis.SetJson(postKey, post)
	if err != nil {
		conf.Log.Error("Failed to cache post, err: %s", err)
	}
	score := time.Now().Unix()+define.VOTE_SCORE
	reply, err := conf.Redis.Do("ZADD", "score:", score, postKey)
	conf.Log.Info("Cache post (%s), reply (%#v)", postKey, reply)
}

func GetCachePost(postID int) (models.Post, error) {
	var post models.Post
	err := conf.Redis.GetJson(fmt.Sprintf("post:%d", postID), &post)
	if err != nil {
		return models.Post{}, err
	}
	return post, err
}

func IsVoted(postID, userID int) (bool, error) {
	key := fmt.Sprintf("voted:%d", postID)
	voted, err := conf.Redis.Do("SADD", key, userID)
	if err != nil {
		conf.Log.Error("Failed to exec: SADD. err: %s", err)
		return false, err
	}
	return voted != 1, err
}

func Vote(postID, userID int) (err error) {
	postKey := fmt.Sprintf("post:%d", postID)
	isVoted, err := IsVoted(postID, userID)
	if isVoted {
		return
	}
	_, err = conf.Redis.Do("ZINCRBY", "score:", postKey, define.VOTE_SCORE)
	if err != nil {
		conf.Log.Error("Failed to exec: ZINCRBY. err: %s", err)
		return
	}
	voteKey := fmt.Sprintf("vote_post:%d", postID)
	_, err = conf.Redis.Do("INCR", voteKey)
	if err != nil {
		conf.Log.Error("Failed to exec: INCR. err: %s", err)
		return
	}
	return
}

