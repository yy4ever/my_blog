package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"my_blog/common/conf"
	perm "my_blog/common/define/permission"
	m "my_blog/common/middleware"
	"my_blog/common/utils"
	"my_blog/control/post"
	"my_blog/control/user"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(utils.AccessLogger(), m.Recover())
	store, _ := redis.NewStore(10, "tcp", fmt.Sprintf("%s:6379", conf.Cnf.RedisHost),
		conf.Cnf.RedisPwd, []byte(conf.Cnf.SecretKey))
	r.Use(sessions.Sessions("blog.session", store))

	loginHandlers := user.Handlers()
	u := r.Group("/u")
	{
		u.POST("", loginHandlers.Register)
		u.GET("", loginHandlers.Current)
		u.POST("/login", loginHandlers.Login)
		u.POST("/logout", m.Permission(), loginHandlers.LogOut)
		u.POST("/follow", m.Permission(), loginHandlers.Follow)
		u.POST("/unfollow", m.Permission(), loginHandlers.UnFollow)
		u.GET("/followers", m.Permission(), loginHandlers.ListFollowers)
		u.GET("/following", m.Permission(), loginHandlers.ListFollowing)
	}
	main := r.Group("")
	{
		main.GET("/users", m.Permission(), loginHandlers.List)
	}

	postHandlers := post.Handlers()
	p := r.Group("/post")
	{
		p.POST("", m.Permission(), postHandlers.Add)
		p.GET("", m.Permission(), postHandlers.List)
		p.DELETE("/:post_id", m.Permission(perm.REMOVE), postHandlers.List)
		p.POST("/:post_id/comment", m.Permission(), postHandlers.AddComment)
		p.GET("/:post_id/comment", m.Permission(), postHandlers.ListComment)
		p.PUT("/:post_id/comment/:comment_id", m.Permission(), postHandlers.DisableComment)
	}
	return r
}
