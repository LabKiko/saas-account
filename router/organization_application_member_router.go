package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerOrganizationApplicationMemberRoutes 注册组织应用成员相关路由
func registerOrganizationApplicationMemberRoutes(group *route.RouterGroup) {
	// 创建依赖
	appRepo := repository.NewOrganizationApplicationRepository()
	appMemberRepo := repository.NewOrganizationApplicationMemberRepository()
	appLimitRepo := repository.NewOrganizationApplicationLimitRepository()
	orgRepo := repository.NewOrganizationRepository()
	userRepo := repository.NewUserRepository()
	appService := service.NewOrganizationApplicationService(appRepo, appMemberRepo, appLimitRepo, orgRepo, userRepo)
	appHandler := handler.NewOrganizationApplicationHandler(appService)

	apps := group.Group("/applications/:app_id")

	// 获取应用成员列表
	apps.GET("/members", appHandler.GetMembers)

	// 添加应用成员
	apps.POST("/members", appHandler.AddMember)

	// 更新应用成员
	apps.PUT("/members/:user_id", appHandler.UpdateMember)

	// 删除应用成员
	apps.DELETE("/members/:user_id", appHandler.RemoveMember)
}
