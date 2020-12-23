package user

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_blog/common/conf"
	"my_blog/common/entity"
	"my_blog/common/errors"
	"my_blog/common/models"
	"my_blog/common/utils"
	"my_blog/control/user/user_manager"
)

type IUser interface {
	Register(c *gin.Context)
	Current(c *gin.Context)
	Login(c *gin.Context)
	LogOut(c *gin.Context)
	List(c *gin.Context)
	Follow(c *gin.Context)
	UnFollow(c *gin.Context)
	ListFollowers(c *gin.Context)
	ListFollowing(c *gin.Context)
}

type User struct {
}

func Handlers() IUser {
	return &User{}
}

func (a User) Register(c *gin.Context) {
	var data entity.UserRegister
	errRes := utils.ErrRes{c}
	if err := c.ShouldBindJSON(&data); err != nil {
		errRes.Response(400, 401, "Invalid params")
		return
	}
	err := user_manager.AddUser(data)
	if err == nil {
		c.JSON(201, "")
	}
	if _, ok := err.(*errors.Err); ok {
		c.JSON(400, err)
	} else {
		c.JSON(400, "Internal error")
	}
}

func (a User) Current(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("User")
	v, _ := userID.(int)
	user, err := user_manager.GetUserByID(v)
	if err != nil {
		c.JSON(500, "Internal error")
		return
	}
	c.JSON(200, user)
}

func (a User) List(c *gin.Context) {
	ret, err := user_manager.List(c.Query("offset"), c.Query("limit"))
	if err != nil {
		c.JSON(500, "Internal error")
	} else {
		c.JSON(200, ret)
	}
}

func (a User) Login(c *gin.Context) {
	var login entity.UserLogin
	errRes := utils.ErrRes{c}
	if err := c.ShouldBindJSON(&login); err != nil {
		errRes.Response(400, 401, "unauthorized")
		return
	}
	user, err := user_manager.GetUserByName(login.Name)
	if err != nil {
		fmt.Printf("Failed user user (name %s)", login.Name)
		errRes.Response(400, 401, "unauthorized")
	}
	if !user_manager.CheckPasswordHash(user.PasswordHash, login.Password) {
		conf.Log.Error("Authorize failed")
	}
	session := sessions.Default(c)
	session.Set("User", user.ID)
	if login.Remember == true {
		session.Options(sessions.Options{MaxAge: conf.Cnf.LoginRememberSeconds})
	}
	session.Save()
	c.String(200, fmt.Sprintf("hello %s.", user.Name))
}

func (a User) LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.Redirect(301,"/user/login")
}

func (a User) Follow(c *gin.Context) {
	errRes := utils.ErrRes{c}
	var f entity.Follow
	err := c.ShouldBindJSON(&f)
	if err != nil {
		errRes.Response(400, 400, "")
		return
	}
	currentUser, _ := c.Get("current_user")
	user, _ := currentUser.(models.User)
	err = user_manager.Follow(user.ID, f.UserID)
	if err != nil {
		errRes.Response(400, 400, err.Error())
		return
	}
	c.JSON(200, "")
}

func (a User) UnFollow(c *gin.Context) {
	errRes := utils.ErrRes{c}
	var f entity.Follow
	err := c.ShouldBindJSON(&f)
	if err != nil {
		errRes.Response(400, 400, "")
		return
	}
	currentUser, _ := c.Get("current_user")
	user, _ := currentUser.(models.User)
	err = user_manager.UnFollow(user.ID, f.UserID)
	if err != nil {
		errRes.Response(400, 400, err.Error())
		return
	}
	c.JSON(200, "")
}


func (a User) ListFollowers(c *gin.Context) {
	currentUser, _ := c.Get("current_user")
	user, _ := currentUser.(models.User)
	users, _ := user_manager.ListFollows(user.ID, "follower")
	c.JSON(200, users)
}

func (a User) ListFollowing(c *gin.Context) {
	currentUser, _ := c.Get("current_user")
	user, _ := currentUser.(models.User)
	users, _ := user_manager.ListFollows(user.ID, "following")
	c.JSON(200, users)
}
