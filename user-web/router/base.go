package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-web/api"
)

// InitBaseRouter 基础路由
func InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	{
		baseRouter.GET("captcha", api.GetCaptcha)
	}

}
