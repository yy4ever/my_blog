package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_blog/control/user/user_manager"
)

func Permission(perms ...int) gin.HandlerFunc {
	return func (c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("User")
		if userID == nil {
			c.JSON(401, "Unauthorized")
			c.Abort()
		} else {
			v, _ := userID.(int)
			user, _ := user_manager.GetUserByID(v)
			for _, p := range perms {
				if !user_manager.UserCan(&user, p) {
					c.JSON(401, "Unauthorized")
					c.Abort()
				}
			}
			c.Set("current_user", user)
			c.Next()
		}
	}
}
