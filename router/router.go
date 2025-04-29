package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz) {
	// API版本前缀
	api := h.Group("/api/v1")

	// 注册用户相关路由
	registerUserRoutes(api)

	// 注册组织相关路由
	registerOrganizationRoutes(api)

	// 注册组织成员相关路由
	registerOrganizationMemberRoutes(api)

	// 注册组织应用相关路由
	registerOrganizationApplicationRoutes(api)

	// 注册组织应用成员相关路由
	registerOrganizationApplicationMemberRoutes(api)

	// 注册组织应用限制相关路由
	registerOrganizationApplicationLimitRoutes(api)

	// 注册应用使用记录相关路由
	registerApplicationUsageRoutes(api)
}
