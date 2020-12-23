package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"my_blog/basic/utils"
	config "my_blog/common/conf"
	"net/http"
	"os"
	"path"
	"syscall"
	"time"
)

var timeFormat = "02/Jan/2006:15:04:05 -0700"

func AccessLogger() gin.HandlerFunc {
	logPath := path.Join(config.Cnf.LogFilePath, "access.log")
	_, err := utils.PathExists(config.Cnf.LogFilePath)
	if err != nil {
		os.MkdirAll(config.Cnf.LogFilePath, os.ModePerm)
	}

	src, err := os.OpenFile(logPath, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	//logger.SetFormatter(&logrus.TextFormatter{})
	return func(c *gin.Context) {
		//hostname, err := os.Hostname()
		//if err != nil {
		//	hostname = "unknow"
		//}
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		//clientUserAgent := c.Request.UserAgent()
		//referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		entry := logger.WithFields(logrus.Fields{})

		msg := fmt.Sprintf("%s [%s %s] => generated %d bytes in %dms (%s  %d) %d header",
			clientIP, c.Request.Method, path, dataLength, latency, c.Request.Proto, statusCode, len(c.Request.Header))
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String() + msg)
		} else {
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
