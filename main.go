package main

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"saas-account/config"
	"saas-account/logger"
	"saas-account/middleware"
	"saas-account/router"
)

func main() {
	// 加载配置
	appConfig := config.GetConfig()

	// 初始化日志
	logger.InitLogger(appConfig.LogLevel, appConfig.LogOutput)
	logger := logger.GetLogger()
	defer logger.Close()

	// 初始化数据库
	config.InitDB()

	// 创建Hertz服务器
	serverAddr := fmt.Sprintf("%s:%d", appConfig.ServerHost, appConfig.ServerPort)
	h := server.Default(server.WithHostPorts(serverAddr))

	// 注册中间件
	middleware.RegisterMiddlewares(h)

	// 注册路由
	router.RegisterRoutes(h)

	// 启动服务器
	logger.Info("Server starting on %s", serverAddr)
	h.Spin()
}
