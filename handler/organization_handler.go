package handler

import (
	"context"
	"saas-account/model"
	"saas-account/service"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// OrganizationHandler 组织处理器
type OrganizationHandler struct {
	orgService service.OrganizationService
}

// NewOrganizationHandler 创建组织处理器
func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// Create 创建组织
func (h *OrganizationHandler) Create(ctx context.Context, c *app.RequestContext) {
	var org model.Organization
	if err := c.BindJSON(&org); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if org.Name == "" {
		BadRequest(c, "组织名称不能为空")
		return
	}

	// 获取创建者ID
	creatorIDStr := c.GetString("user_id")
	creatorID, err := strconv.ParseUint(creatorIDStr, 10, 64)
	if err != nil {
		Unauthorized(c, "未授权")
		return
	}

	// 创建组织
	if err := h.orgService.Create(context.Background(), &org, uint(creatorID)); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, org)
}

// GetByID 根据ID获取组织
func (h *OrganizationHandler) GetByID(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	org, err := h.orgService.GetByID(context.Background(), uint(id))
	if err != nil {
		NotFound(c, "组织不存在")
		return
	}

	Success(c, org)
}

// List 获取组织列表
func (h *OrganizationHandler) List(ctx context.Context, c *app.RequestContext) {
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

	// 获取组织列表
	orgs, total, err := h.orgService.List(ctx, page, pageSize)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithPagination(c, orgs, total, page, pageSize)
}

// Update 更新组织
func (h *OrganizationHandler) Update(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	var org model.Organization
	if err := c.BindJSON(&org); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置组织ID
	org.ID = id

	// 更新组织
	if err := h.orgService.Update(ctx, &org); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	// 获取更新后的组织
	updatedOrg, err := h.orgService.GetByID(ctx, uint(id))
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c, updatedOrg)
}

// Delete 删除组织
func (h *OrganizationHandler) Delete(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	if err := h.orgService.Delete(ctx, uint(id)); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// GetMembers 获取组织成员列表
func (h *OrganizationHandler) GetMembers(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
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

	// 获取组织成员列表
	members, total, err := h.orgService.GetMembers(ctx, uint(id), page, pageSize)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithPagination(c, members, total, page, pageSize)
}

// AddMember 添加组织成员
func (h *OrganizationHandler) AddMember(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	// 获取请求参数
	var req struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
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

	// 添加组织成员
	if err := h.orgService.AddMember(ctx, uint(id), req.UserID, req.Role); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// UpdateMember 更新组织成员
func (h *OrganizationHandler) UpdateMember(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
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
		Role string `json:"role"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 更新组织成员
	if err := h.orgService.UpdateMember(ctx, uint(id), uint(userID), req.Role); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// RemoveMember 移除组织成员
func (h *OrganizationHandler) RemoveMember(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的组织ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	// 移除组织成员
	if err := h.orgService.RemoveMember(ctx, uint(id), uint(userID)); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}
