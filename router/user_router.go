package router

import (
	"saas-account/handler"
	"saas-account/repository"
	"saas-account/service"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerUserRoutes 注册用户相关路由
func registerUserRoutes(group *route.RouterGroup) {
	// 创建依赖
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	users := group.Group("/users")

	// 创建用户
	users.POST("", userHandler.Create)

	// 获取用户列表
	users.GET("", userHandler.List)

	// 获取单个用户
	users.GET("/:id", userHandler.GetByID)

	// 更新用户
	users.PUT("/:id", userHandler.Update)

	// 删除用户
	users.DELETE("/:id", userHandler.Delete)

	// 修改密码
	users.POST("/:id/change-password", userHandler.ChangePassword)
}
