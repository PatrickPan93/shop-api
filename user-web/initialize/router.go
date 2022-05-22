package initialize

import (
	"github.com/gin-gonic/gin"
	router2 "shop-api/user-web/router"
)

func Routers() *gin.Engine {

	router := gin.Default()
	ApiGroup := router.Group("/v1")
	router2.InitUserRouter(ApiGroup)

	return router
}
