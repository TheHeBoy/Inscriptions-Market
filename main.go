package main

import (
	"fmt"
	"gohub/cmd"
	"gohub/internal/mqI"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/database"
	"gohub/pkg/eth"
	"gohub/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "Gohub",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {

			// 配置初始化，依赖命令行 --env 参数
			config.SetupConfig(cmd.Env)

			// 初始化 Logger
			logger.SetupLogger()

			// 初始化数据库
			database.SetupDB()

			// 初始化以太坊客户端
			eth.SetupEth()

			// 初始化消息队列
			mqI.SetupMQ()
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		cmd.CmdPlay,
		cmd.CmdMigrate,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
