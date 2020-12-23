package utils

import (
	"github.com/gin-gonic/gin"
)


type ErrRes struct {
	Ctx *gin.Context
}


func (e *ErrRes) Response (statusCode, code int, message string)  {
	e.Ctx.JSON(
		statusCode,
		gin.H{
			"error": gin.H{
				"code": code,
				"message": message,
			},
		},
	)
	return
}
