package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerOrganizationApplicationLimitRoutes 注册组织应用限制相关路由
func registerOrganizationApplicationLimitRoutes(group *route.RouterGroup) {
	// 创建依赖
	appRepo := repository.NewOrganizationApplicationRepository()
	appMemberRepo := repository.NewOrganizationApplicationMemberRepository()
	appLimitRepo := repository.NewOrganizationApplicationLimitRepository()
	orgRepo := repository.NewOrganizationRepository()
	userRepo := repository.NewUserRepository()
	appService := service.NewOrganizationApplicationService(appRepo, appMemberRepo, appLimitRepo, orgRepo, userRepo)
	appHandler := handler.NewOrganizationApplicationHandler(appService)

	apps := group.Group("/applications/:app_id")

	// 获取应用限制
	apps.GET("/limits", appHandler.GetLimit)

	// 创建或更新应用限制
	apps.PUT("/limits", appHandler.SetLimit)
}
