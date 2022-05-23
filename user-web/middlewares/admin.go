package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你没有权限请求该接口",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
