package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/hertz-contrib/cors"
)

// CORS 中间件，处理跨域请求
func CORS() app.HandlerFunc {
	return cors.New(cors.DefaultConfig())
}
