package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my_blog/common/conf"
	"runtime/debug"
)


func Recover()  gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("PANIC - StackTrace:\n%s%s\n", string(debug.Stack()), r)
				conf.Log.Debug("%s", errMsg)
				fmt.Println(errMsg)
				c.String(500, "internal error")
			}
		}()
		c.Next()
	}
}
