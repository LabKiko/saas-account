package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerApplicationUsageRoutes 注册应用使用记录相关路由
func registerApplicationUsageRoutes(group *route.RouterGroup) {
	// 创建依赖
	usageRepo := repository.NewApplicationUsageRepository()
	appRepo := repository.NewOrganizationApplicationRepository()
	limitRepo := repository.NewOrganizationApplicationLimitRepository()
	usageService := service.NewApplicationUsageService(usageRepo, appRepo, limitRepo)
	usageHandler := handler.NewApplicationUsageHandler(usageService)

	apps := group.Group("/applications/:app_id")

	// 获取应用使用记录列表
	apps.GET("/usages", usageHandler.List)

	// 创建应用使用记录
	apps.POST("/usages", usageHandler.Create)

	// 获取应用使用统计
	apps.GET("/usages/summary", usageHandler.GetSummary)

	// 记录API使用
	apps.POST("/usages/api", usageHandler.RecordAPIUsage)

	// 记录存储使用
	apps.POST("/usages/storage", usageHandler.RecordStorageUsage)

	// 记录功能使用
	apps.POST("/usages/feature", usageHandler.RecordFeatureUsage)
}
