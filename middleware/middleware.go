package middleware

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterMiddlewares 注册所有中间件
func RegisterMiddlewares(h *server.Hertz) {
	// 注册全局中间件
	h.Use(
		Recovery(),   // 恢复中间件，必须放在最前面
		RequestID(),  // 请求ID中间件
		Logger(),     // 日志中间件
		ErrorLogger(), // 错误日志中间件
		CORS(),       // CORS中间件
	)
}
