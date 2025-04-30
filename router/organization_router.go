package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerOrganizationRoutes 注册组织相关路由
func registerOrganizationRoutes(group *route.RouterGroup) {
	// 创建依赖
	orgRepo := repository.NewOrganizationRepository()
	orgMemberRepo := repository.NewOrganizationMemberRepository()
	userRepo := repository.NewUserRepository()
	orgService := service.NewOrganizationService(orgRepo, orgMemberRepo, userRepo)
	orgHandler := handler.NewOrganizationHandler(orgService)

	orgs := group.Group("/organizations")

	// 创建组织
	orgs.POST("", orgHandler.Create)

	// 获取组织列表
	orgs.GET("", orgHandler.List)

	// 获取单个组织
	orgs.GET("/:id", orgHandler.GetByID)

	// 更新组织
	orgs.PUT("/:id", orgHandler.Update)

	// 删除组织
	orgs.DELETE("/:id", orgHandler.Delete)

	// 获取组织成员列表
	orgs.GET("/:id/members", orgHandler.GetMembers)

	// 添加组织成员
	orgs.POST("/:id/members", orgHandler.AddMember)

	// 更新组织成员
	orgs.PUT("/:id/members/:user_id", orgHandler.UpdateMember)

	// 移除组织成员
	orgs.DELETE("/:id/members/:user_id", orgHandler.RemoveMember)
}
