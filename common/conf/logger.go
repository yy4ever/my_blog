package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"my_blog/basic/utils"
	"os"
	"path"
	"syscall"
)

type Logger struct {
	entry *logrus.Entry
}

var Log Logger

func (l Logger) Debug(_fmt string, msg interface{}) {
	l.entry.Debug(fmt.Sprintf(_fmt, msg))
}

func (l Logger) Info(_fmt string, msg interface{}) {
	l.entry.Info(fmt.Sprintf(_fmt, msg))
}

func (l Logger) Warn(_fmt string, msg interface{}) {
	l.entry.Warn(fmt.Sprintf(_fmt, msg))
}

func (l Logger) Error(_fmt string, msg ...interface{}) {
	l.entry.Error(fmt.Sprintf(_fmt, msg))
}

func InitLogger(module string) {
	logPath := path.Join(Cnf.LogFilePath, fmt.Sprintf("%s.log", module))
	_, err := utils.PathExists(Cnf.LogFilePath)
	if err != nil {
		os.MkdirAll(Cnf.LogFilePath, os.ModePerm)
	}

	src, err := os.OpenFile(logPath, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})
	entry := logger.WithFields(logrus.Fields{})
	Log = Logger{entry: entry}
}
