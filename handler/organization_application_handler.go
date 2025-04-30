package handler

import (
	"context"
	"saas-account/model"
	"saas-account/service"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// OrganizationApplicationHandler 组织应用处理器
type OrganizationApplicationHandler struct {
	appService service.OrganizationApplicationService
}

// NewOrganizationApplicationHandler 创建组织应用处理器
func NewOrganizationApplicationHandler(appService service.OrganizationApplicationService) *OrganizationApplicationHandler {
	return &OrganizationApplicationHandler{
		appService: appService,
	}
}

// Create 创建组织应用
func (h *OrganizationApplicationHandler) Create(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	var app model.OrganizationApplication
	if err := c.BindJSON(&app); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置组织ID
	app.OrganizationId = uint(orgID)

	// 验证必填字段
	if app.Name == "" || app.Type == "" {
		BadRequest(c, "应用名称和类型不能为空")
		return
	}

	// 创建应用
	if err := h.appService.Create(context.Background(), &app); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, app)
}

// GetByID 根据ID获取组织应用
func (h *OrganizationApplicationHandler) GetByID(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	_, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	app, err := h.appService.GetByID(context.Background(), uint(appID))
	if err != nil {
		NotFound(c, "应用不存在")
		return
	}

	// 隐藏应用密钥
	app.AppSecret = ""

	Success(c, app)
}

// List 获取组织应用列表
func (h *OrganizationApplicationHandler) List(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 获取应用列表
	apps, total, err := h.appService.GetByOrganization(context.Background(), uint(orgID), page, pageSize)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// 隐藏应用密钥
	for i := range apps {
		apps[i].AppSecret = ""
	}

	SuccessWithPagination(c, apps, total, page, pageSize)
}

// Update 更新组织应用
func (h *OrganizationApplicationHandler) Update(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	_, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	var app model.OrganizationApplication
	if err := c.BindJSON(&app); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置应用ID
	app.ID = uint(appID)

	// 更新应用
	if err := h.appService.Update(context.Background(), &app); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	// 获取更新后的应用
	updatedApp, err := h.appService.GetByID(context.Background(), uint(appID))
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// 隐藏应用密钥
	updatedApp.AppSecret = ""

	Success(c, updatedApp)
}

// Delete 删除组织应用
func (h *OrganizationApplicationHandler) Delete(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	_, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	if err := h.appService.Delete(context.Background(), uint(appID)); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// RegenerateAppSecret 重新生成应用密钥
func (h *OrganizationApplicationHandler) RegenerateAppSecret(c *app.RequestContext) {
	orgIDStr := c.Param("org_id")
	_, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	appSecret, err := h.appService.RegenerateAppSecret(context.Background(), uint(appID))
	if err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, map[string]string{
		"app_secret": appSecret,
	})
}

// GetMembers 获取应用成员列表
func (h *OrganizationApplicationHandler) GetMembers(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 获取应用成员列表
	members, total, err := h.appService.GetMembers(context.Background(), uint(appID), page, pageSize)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithPagination(c, members, total, page, pageSize)
}

// AddMember 添加应用成员
func (h *OrganizationApplicationHandler) AddMember(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取请求参数
	var req struct {
		UserID      uint   `json:"user_id"`
		Role        string `json:"role"`
		Permissions string `json:"permissions"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.UserID == 0 {
		BadRequest(c, "用户ID不能为空")
		return
	}

	// 添加应用成员
	if err := h.appService.AddMember(context.Background(), uint(appID), req.UserID, req.Role, req.Permissions); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// UpdateMember 更新应用成员
func (h *OrganizationApplicationHandler) UpdateMember(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	// 获取请求参数
	var req struct {
		Role        string `json:"role"`
		Permissions string `json:"permissions"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 更新应用成员
	if err := h.appService.UpdateMember(context.Background(), uint(appID), uint(userID), req.Role, req.Permissions); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// RemoveMember 移除应用成员
func (h *OrganizationApplicationHandler) RemoveMember(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	// 移除应用成员
	if err := h.appService.RemoveMember(context.Background(), uint(appID), uint(userID)); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// GetLimit 获取应用限制
func (h *OrganizationApplicationHandler) GetLimit(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	limit, err := h.appService.GetLimit(context.Background(), uint(appID))
	if err != nil {
		NotFound(c, "应用限制不存在")
		return
	}

	Success(c, limit)
}

// SetLimit 设置应用限制
func (h *OrganizationApplicationHandler) SetLimit(c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	var limit model.OrganizationApplicationLimit
	if err := c.BindJSON(&limit); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置应用ID
	limit.OrganizationApplicationId = uint(appID)

	// 设置应用限制
	if err := h.appService.SetLimit(context.Background(), &limit); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, limit)
}
