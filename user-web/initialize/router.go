package initialize

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-web/middlewares"
	router2 "shop-api/user-web/router"
)

func Routers() *gin.Engine {

	router := gin.Default()
	// 支持跨域
	router.Use(middlewares.Cors())
	ApiGroup := router.Group("/v1")
	// 用户路由
	router2.InitUserRouter(ApiGroup)
	// 基础路由
	router2.InitBaseRouter(ApiGroup)

	return router
}
