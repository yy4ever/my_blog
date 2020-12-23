package post

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_blog/common/define/permission"
	"my_blog/common/entity"
	"my_blog/common/models"
	"my_blog/common/utils"
	"my_blog/control/post/post_manager"
	"my_blog/control/user/user_manager"
	"strconv"
)

type IPost interface {
	Add(c *gin.Context)
	List(c *gin.Context)
	Remove(c *gin.Context)
	AddComment(c *gin.Context)
	ListComment(c *gin.Context)
	DisableComment(c *gin.Context)
}

type Post struct {
}

func Handlers() IPost {
	return &Post{}
}

func (p Post) Add(c *gin.Context) {
	var data entity.Post
	errRes := utils.ErrRes{Ctx:c}
	if err := c.ShouldBindJSON(&data); err != nil {
		errRes.Response(400, 401, "Invalid params")
		return
	}
	sess := sessions.Default(c)
	userID, _ := sess.Get("User").(int)
	err := post_manager.Add(&data, userID)
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(201, "")
}

func (p Post) List(c *gin.Context) {
	errRes := utils.ErrRes{Ctx:c}
	sess := sessions.Default(c)
	userID, _ := sess.Get("User").(int)
	posts, err := post_manager.List(userID, c.Query("offset"), c.Query("limit"))
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(200, posts)
}

func (p Post) Remove(c *gin.Context) {
	errRes := utils.ErrRes{Ctx:c}
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		errRes.Response(400, 400, "Invalid params")
		return
	}
	err = post_manager.Remove(postID)
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(200, "")
}

func (p Post) DisableComment(c *gin.Context) {
	errRes := utils.ErrRes{Ctx:c}
	commentID, err1 := strconv.Atoi(c.Param("comment_id"))
	postID, err2 := strconv.Atoi(c.Param("comment_id"))
	if err1 != nil || err2 != nil {
		errRes.Response(400, 400, "Invalid params")
		return
	}
	post, err := post_manager.Get(postID)
	if err != nil {
		errRes.Response(400, 400, "Invalid params")
		return
	}
	currentUser, _ := c.Get("current_user")
	user, _ := currentUser.(models.User)
	if user_manager.UserCan(&user, permission.MODERATE) && post.AuthorID != user.ID {
		errRes.Response(401, 401, "Not allowed")
		return
	}
	err = post_manager.DisableComment(postID, commentID)
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(200, "")
}

func (p Post) AddComment(c *gin.Context) {
	var data entity.Comment
	errRes := utils.ErrRes{Ctx:c}
	if err := c.ShouldBindJSON(&data); err != nil {
		errRes.Response(400, 401, "Invalid params")
		return
	}
	sess := sessions.Default(c)
	userID, _ := sess.Get("User").(int)
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		errRes.Response(400, 400, "Invalid params")
		return
	}
	err = post_manager.AddComment(data, userID, postID)
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(201, "")
}

func (p Post) ListComment(c *gin.Context) {
	errRes := utils.ErrRes{Ctx:c}
	sess := sessions.Default(c)
	userID, _ := sess.Get("User").(int)
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		errRes.Response(400, 400, "Invalid params")
		return
	}
	comments, err := post_manager.ListComment(postID, userID)
	if err != nil {
		errRes.Response(500, 500, "Internal error")
		return
	}
	c.JSON(200, comments)
}
