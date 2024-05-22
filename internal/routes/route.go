package routes

import (
	"gohub/internal/controller/app"
	"gohub/internal/routes/middlewares"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {

	// 注册全局中间件
	registerGlobalMiddleWare(router)

	//  注册 API 路由
	RegisterAPIRoutes(router)

	//  配置 404 路由
	setup404Handler(router)
}

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	appGroup := r.Group("/app")
	//api.Use(middlewares.LimitIP("200-H"))
	appRoutes(appGroup)

	adminGroup := r.Group("/admin")
	//admin.Use(middlewares.LimitIP("200-H"))
	adminRoutes(adminGroup)
}

func appRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	lgc := new(app.LoginController)
	authGroup.GET("/message", middlewares.GuestJWT(), lgc.GetMessageAuth)
	authGroup.POST("/login", middlewares.GuestJWT(), lgc.LoginBySignatureAuth)
	authGroup.POST("/refresh-token", lgc.RefreshToken)

	tokenGroup := r.Group("/tokens")
	tcl := new(app.TokenController)
	tokenGroup.GET("/page", tcl.PageTokens)
	tokenGroup.GET("/page-listing", tcl.PageListingToken)
	tokenGroup.GET("/:address", tcl.GetTokensByAddress)

	orderGroup := r.Group("/orders")
	ocl := new(app.OrderController)
	orderGroup.PUT("/create", ocl.CreateOrder)
	orderGroup.PUT("/sign", ocl.SignOrder)
	orderGroup.GET("/listing", ocl.GetListingOrderByTick)
	orderGroup.GET("/:address", ocl.PageBySeller)

	msc20Group := r.Group("/msc20")
	m20c := new(app.Msc20Controller)
	msc20Group.GET("/:address", m20c.GetMsc20ByAddress)

	inscriptionGroup := r.Group("/inscriptions")
	itc := new(app.InscriptionController)
	inscriptionGroup.GET("/latest", itc.GetLatest)
}

func adminRoutes(r *gin.RouterGroup) {}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.ForceUA(),
		middlewares.Cors(),
	)
}

func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}
