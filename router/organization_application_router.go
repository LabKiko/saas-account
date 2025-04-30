package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerOrganizationApplicationRoutes 注册组织应用相关路由
func registerOrganizationApplicationRoutes(group *route.RouterGroup) {
	// 创建依赖
	appRepo := repository.NewOrganizationApplicationRepository()
	appMemberRepo := repository.NewOrganizationApplicationMemberRepository()
	appLimitRepo := repository.NewOrganizationApplicationLimitRepository()
	orgRepo := repository.NewOrganizationRepository()
	userRepo := repository.NewUserRepository()
	appService := service.NewOrganizationApplicationService(appRepo, appMemberRepo, appLimitRepo, orgRepo, userRepo)
	appHandler := handler.NewOrganizationApplicationHandler(appService)

	orgs := group.Group("/organizations/:org_id")

	// 获取组织应用列表
	orgs.GET("/applications", appHandler.List)

	// 创建组织应用
	orgs.POST("/applications", appHandler.Create)

	// 获取单个组织应用
	orgs.GET("/applications/:id", appHandler.GetByID)

	// 更新组织应用
	orgs.PUT("/applications/:id", appHandler.Update)

	// 删除组织应用
	orgs.DELETE("/applications/:id", appHandler.Delete)

	// 重新生成应用密钥
	orgs.POST("/applications/:id/regenerate-secret", appHandler.RegenerateAppSecret)
}
