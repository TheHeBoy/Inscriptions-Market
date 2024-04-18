// Package routes 注册路由
package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controllers/api/auth"
	"gohub/app/http/middlewares"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	api := r.Group("/api")
	api.Use(middlewares.LimitIP("200-H"))
	apiRoutes(api)

	admin := r.Group("/admin")
	admin.Use(middlewares.LimitIP("200-H"))
	adminRoutes(admin)
}

func apiRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	lgc := new(auth.LoginController)
	authGroup.GET("/message", middlewares.GuestJWT(), lgc.GetMessage)
	authGroup.POST("/login", middlewares.GuestJWT(), lgc.LoginBySignature)
	authGroup.POST("/refresh-token", lgc.RefreshToken)
}

func adminRoutes(r *gin.RouterGroup) {

}
