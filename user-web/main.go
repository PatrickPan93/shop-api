package main

import (
	"fmt"
	"go.uber.org/zap"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
)

func main() {

	// 1. 初始化logger
	initialize.InitLogger()

	// 2. 初始化配置
	initialize.InitConfig()

	// 3. 初始化routers
	router := initialize.Routers()

	// 4.  初始化验证器
	initialize.InitTrans("zh")

	// 5. 注册自定义mobile验证器
	initialize.InitCustomValidator()

	zap.S().Debugf("going to start web server, listen port %d\n", 8080)
	zap.S().Fatal(router.Run(fmt.Sprintf("%s:%d", global.ServerConfig.UserWebInfo.Host, global.ServerConfig.UserWebInfo.Port)))
}
