package conf

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func InitSession(r *gin.Engine) {
	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:6379", Cnf.RedisHost),
		Cnf.RedisPwd, []byte(Cnf.SecretKey))
	if err != nil {
		Log.Error("Failed to init redis store, err: %s", err)
	}
	r.Use(sessions.Sessions("blog.session", store))
}
