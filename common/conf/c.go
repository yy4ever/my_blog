package conf

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

var (
	DB  *sqlx.DB
	Cnf *Conf
)

func DefaultInit() {
	InitConf()
	InitLogger("blog")
	InitDB()
}

func InitConf() {
	var files []string
	pwd, _ := filepath.Abs(".")
	files, _ = filepath.Glob("./env.*.yaml")
	var filename string
	for _, f := range files {
		switch f {
		case "env.dev.yaml":
			filename = f
		case "env.prod.yaml":
			filename = f
		default:
			continue
		}
	}
	if filename == "" {
		fmt.Println("Environment config file not found!")
		return
	}
	yamlFile, err := ioutil.ReadFile(filepath.Join(pwd, filename))
	if err != nil {
		fmt.Printf("Open file error: %s", err)
		return
	}
	var cnf Conf
	err = yaml.Unmarshal(yamlFile, &cnf)
	Cnf = &cnf
	return
}

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		Cnf.DbUser, Cnf.DbPassword, Cnf.DbHost, Cnf.DbPort, Cnf.DbDataBase)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		Log.Error("Failed to connect to mysql server. %s", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	DB = db
	timer := time.NewTicker(time.Minute * 10)
	go func(DB *sqlx.DB) {
		for _ = range timer.C {
			if err = DB.Ping(); err != nil {
				autoConnectMysql()
			}
		}
	}(DB)
	return nil
}

func autoConnectMysql() {
	retryTimes := 5
	for retryTimes > 0 {
		retryTimes--
		if err := DB.Ping(); err != nil {
			fmt.Sprintf("Failed to connect to mysql, retry %d times", 5-retryTimes)
			continue
		}
		fmt.Println("Connecting to mysql successfully")
		break
	}
}
