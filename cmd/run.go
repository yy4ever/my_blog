package main

import (
	"fmt"
	"my_blog/common/conf"
	"my_blog/common/router"
	"my_blog/control/user/user_manager"
)

func init() {
	conf.DefaultInit()
	user_manager.InitData()
}

func main() {
	engine := router.InitRouter()
	fmt.Print("Service starting.")
	engine.Run(conf.Cnf.PORT)
	//gin.DebugPrintRouteFunc()
}
