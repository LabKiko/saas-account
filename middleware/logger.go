package middleware

import (
	"context"
	"saas-account/logger"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// Logger 中间件，记录请求和响应信息
func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取请求开始时间
		start := time.Now()

		// 获取请求ID
		requestID := GetRequestID(ctx)
		traceID := GetTraceID(ctx)
		spanID := GetSpanID(ctx)

		// 获取请求信息
		method := string(ctx.Request.Method())
		path := string(ctx.Request.URI().Path())
		ip := ctx.ClientIP()
		userAgent := string(ctx.Request.Header.UserAgent())

		// 记录请求开始
		logger.Logger.Info("[%s] [%s] [%s] Request started: %s %s from %s, User-Agent: %s",
			requestID, traceID, spanID, method, path, ip, userAgent)

		// 继续处理请求
		ctx.Next(c)

		// 获取响应信息
		statusCode := ctx.Response.StatusCode()
		latency := time.Since(start)

		// 记录请求结束
		logger.Logger.Info("[%s] [%s] [%s] Request completed: %s %s, status: %d, latency: %v",
			requestID, traceID, spanID, method, path, statusCode, latency)
	}
}

// ErrorLogger 中间件，记录错误信息
func ErrorLogger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 继续处理请求
		ctx.Next(c)

		// 检查是否有错误
		if len(ctx.Errors) > 0 {
			// 获取请求ID
			requestID := GetRequestID(ctx)
			traceID := GetTraceID(ctx)
			spanID := GetSpanID(ctx)

			// 获取请求信息
			method := string(ctx.Request.Method())
			path := string(ctx.Request.URI().Path())

			// 记录错误

			for _, err := range ctx.Errors {
				logger.Logger.Error("[%s] [%s] [%s] Request error: %s %s, error: %v",
					requestID, traceID, spanID, method, path, err.Err)
			}
		}
	}
}
