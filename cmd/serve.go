package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gohub/internal/ethI"
	"gohub/internal/routes"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	mode := gin.DebugMode
	if app.IsProduction() {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// gin 实例
	router := gin.New()

	// 初始化路由绑定
	routes.SetupRoute(router)

	// 初始化 Redis
	redis.SetupRedis()

	// 监听订单日志
	ethI.ListeningOrderLog()

	// 运行服务器
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.Error("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}
