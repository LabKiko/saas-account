package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"saas-account/logger"

	"github.com/cloudwego/hertz/pkg/app"
)

// Recovery 中间件，捕获和处理panic
func Recovery() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求ID
				requestID := GetRequestID(ctx)
				traceID := GetTraceID(ctx)
				spanID := GetSpanID(ctx)

				// 获取请求信息
				method := string(ctx.Request.Method())
				path := string(ctx.Request.URI().Path())

				// 记录panic信息
				logger.Logger.Error("[%s] [%s] [%s] Panic recovered: %s %s, error: %v\nStack trace:\n%s",
					requestID, traceID, spanID, method, path, err, debug.Stack())

				// 返回500错误
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code":       500,
					"message":    "Internal Server Error",
					"request_id": requestID,
				})
				ctx.Abort()
			}
		}()

		// 继续处理请求
		ctx.Next(c)
	}
}
