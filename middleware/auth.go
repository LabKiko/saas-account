package middleware

import (
	"context"
	"net/http"
	"saas-account/utils"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// Auth 中间件，验证JWT令牌
func Auth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取Authorization头
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":       401,
				"message":    "Authorization header is required",
				"request_id": GetRequestID(ctx),
			})
			ctx.Abort()
			return
		}

		// 检查Authorization头格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":       401,
				"message":    "Authorization header format must be Bearer {token}",
				"request_id": GetRequestID(ctx),
			})
			ctx.Abort()
			return
		}

		// 解析JWT令牌
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":       401,
				"message":    "Invalid or expired token",
				"request_id": GetRequestID(ctx),
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储在上下文中
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)

		// 继续处理请求
		ctx.Next(c)
	}
}

// AdminAuth 中间件，验证用户是否为管理员
func AdminAuth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 先执行Auth中间件
		Auth()(c, ctx)
		if ctx.IsAborted() {
			return
		}

		// 检查用户角色
		role, exists := ctx.Get("role")
		if !exists || role.(string) != "admin" {
			ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code":       403,
				"message":    "Admin access required",
				"request_id": GetRequestID(ctx),
			})
			ctx.Abort()
			return
		}

		// 继续处理请求
		ctx.Next(c)
	}
}
