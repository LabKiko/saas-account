package main

import (
	"log"
	"saas-account/config"
	"saas-account/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 创建Hertz服务器
	h := server.Default(server.WithHostPorts(":8080"))

	// 注册路由
	router.RegisterRoutes(h)

	// 启动服务器
	log.Println("Server starting on :8080")
	h.Spin()
}
