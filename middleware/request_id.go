package middleware

import (
	"context"
	"saas-account/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// 上下文键
const (
	RequestIDKey = "X-Request-ID"
	TraceIDKey   = "X-Trace-ID"
	SpanIDKey    = "X-Span-ID"
)

// RequestID 中间件，为每个请求生成唯一的请求ID
func RequestID() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 检查请求头中是否已有请求ID
		requestID := string(ctx.Request.Header.Peek(RequestIDKey))
		if requestID == "" {
			// 如果没有，生成一个新的请求ID
			requestID = utils.GenerateStringID()
			ctx.Request.Header.Set(RequestIDKey, requestID)
		}

		// 生成跟踪ID和跨度ID
		traceID := string(ctx.Request.Header.Peek(TraceIDKey))
		if traceID == "" {
			traceID = utils.GenerateStringID()
			ctx.Request.Header.Set(TraceIDKey, traceID)
		}

		spanID := utils.GenerateStringID()
		ctx.Request.Header.Set(SpanIDKey, spanID)

		// 设置响应头
		ctx.Response.Header.Set(RequestIDKey, requestID)
		ctx.Response.Header.Set(TraceIDKey, traceID)
		ctx.Response.Header.Set(SpanIDKey, spanID)

		// 将请求ID存储在上下文中
		ctx.Set(RequestIDKey, requestID)
		ctx.Set(TraceIDKey, traceID)
		ctx.Set(SpanIDKey, spanID)

		// 继续处理请求
		ctx.Next(c)
	}
}

// GetRequestID 从上下文中获取请求ID
func GetRequestID(ctx *app.RequestContext) string {
	if requestID, exists := ctx.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}

// GetTraceID 从上下文中获取跟踪ID
func GetTraceID(ctx *app.RequestContext) string {
	if traceID, exists := ctx.Get(TraceIDKey); exists {
		return traceID.(string)
	}
	return ""
}

// GetSpanID 从上下文中获取跨度ID
func GetSpanID(ctx *app.RequestContext) string {
	if spanID, exists := ctx.Get(SpanIDKey); exists {
		return spanID.(string)
	}
	return ""
}

// WithRequestContext 将请求上下文信息添加到context.Context中
func WithRequestContext(ctx context.Context, reqCtx *app.RequestContext) context.Context {
	ctx = context.WithValue(ctx, RequestIDKey, GetRequestID(reqCtx))
	ctx = context.WithValue(ctx, TraceIDKey, GetTraceID(reqCtx))
	ctx = context.WithValue(ctx, SpanIDKey, GetSpanID(reqCtx))
	return ctx
}

// GetRequestIDFromContext 从context.Context中获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetTraceIDFromContext 从context.Context中获取跟踪ID
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GetSpanIDFromContext 从context.Context中获取跨度ID
func GetSpanIDFromContext(ctx context.Context) string {
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok {
		return spanID
	}
	return ""
}
