package handler

import (
	"context"
	"saas-account/logger"
	"saas-account/middleware"
	"saas-account/model"
	"saas-account/service"
	"saas-account/utils"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create 创建用户
func (h *UserHandler) Create(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	logger.Logger.InfoWithContext(reqCtx, "开始创建用户")

	var user model.User
	if err := c.BindAndValidate(&user); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "解析请求参数失败: %v", err)
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if user.Name == "" || user.Email == "" || user.Password == "" {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了空的必填字段")
		BadRequest(c, "姓名、邮箱和密码不能为空")
		return
	}

	// 验证邮箱格式
	if !utils.IsValidEmail(user.Email) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了无效的邮箱格式: %s", user.Email)
		BadRequest(c, "邮箱格式无效")
		return
	}

	// 验证手机号格式
	if user.Phone != "" && !utils.IsValidPhone(user.Phone) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了无效的手机号格式: %s", user.Phone)
		BadRequest(c, "手机号格式无效")
		return
	}

	// 验证密码强度
	if !utils.IsStrongPassword(user.Password) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了弱密码")
		BadRequest(c, "密码强度不足，密码应包含大小写字母、数字和特殊字符，长度至少8位")
		return
	}

	// 创建用户
	logger.Logger.InfoWithContext(reqCtx, "开始创建用户: %s (%s)", user.Name, user.Email)
	if err := h.userService.Create(reqCtx, &user); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "创建用户失败: %v", err)
		Fail(c, 500, err.Error())
		return
	}

	// 清除敏感信息
	user.Password = ""

	logger.Logger.InfoWithContext(reqCtx, "用户创建成功: ID=%d, 名称=%s", user.ID, user.Name)
	Success(c, user)
}

// GetByID 根据ID获取用户
func (h *UserHandler) GetByID(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	idStr := c.Param("id")
	logger.Logger.InfoWithContext(reqCtx, "开始获取用户信息: ID=%s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "无效的用户ID: %s, 错误: %v", idStr, err)
		BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetByID(reqCtx, uint(id))
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "用户不存在: ID=%d, 错误: %v", id, err)
		NotFound(c, "用户不存在")
		return
	}

	// 清除敏感信息
	user.Password = ""

	logger.Logger.InfoWithContext(reqCtx, "成功获取用户信息: ID=%d, 名称=%s", user.ID, user.Name)
	Success(c, user)
}

// List 获取用户列表
func (h *UserHandler) List(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	logger.Logger.InfoWithContext(reqCtx, "开始获取用户列表: page=%s, page_size=%s", pageStr, pageSizeStr)

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Logger.WarnWithContext(reqCtx, "无效的页码: %s, 使用默认值1", pageStr)
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		logger.Logger.WarnWithContext(reqCtx, "无效的页面大小: %s, 使用默认值10", pageSizeStr)
		pageSize = 10
	}

	// 获取用户列表
	logger.Logger.InfoWithContext(reqCtx, "查询用户列表: page=%d, page_size=%d", page, pageSize)
	users, total, err := h.userService.List(reqCtx, page, pageSize)
	if err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "获取用户列表失败: %v", err)
		InternalServerError(c, err.Error())
		return
	}

	// 清除敏感信息
	for i := range users {
		users[i].Password = ""
	}

	logger.Logger.InfoWithContext(reqCtx, "成功获取用户列表: 总数=%d, 当前页=%d, 每页数量=%d", total, page, pageSize)
	SuccessWithPagination(c, users, total, page, pageSize)
}

// Update 更新用户
func (h *UserHandler) Update(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	idStr := c.Param("id")
	logger.Logger.InfoWithContext(reqCtx, "开始更新用户: ID=%s", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "无效的用户ID: %s, 错误: %v", idStr, err)
		BadRequest(c, "无效的用户ID")
		return
	}

	var user model.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "解析请求参数失败: %v", err)
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置用户ID
	user.ID = id

	// 验证邮箱格式
	if user.Email != "" && !utils.IsValidEmail(user.Email) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了无效的邮箱格式: %s", user.Email)
		BadRequest(c, "邮箱格式无效")
		return
	}

	// 验证手机号格式
	if user.Phone != "" && !utils.IsValidPhone(user.Phone) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了无效的手机号格式: %s", user.Phone)
		BadRequest(c, "手机号格式无效")
		return
	}

	// 更新用户
	logger.Logger.InfoWithContext(reqCtx, "开始更新用户信息: ID=%d", id)
	if err := h.userService.Update(reqCtx, &user); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "更新用户失败: %v", err)
		Fail(c, 500, err.Error())
		return
	}

	// 获取更新后的用户
	updatedUser, err := h.userService.GetByID(reqCtx, uint(id))
	if err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "获取更新后的用户失败: %v", err)
		InternalServerError(c, err.Error())
		return
	}

	// 清除敏感信息
	updatedUser.Password = ""

	logger.Logger.InfoWithContext(reqCtx, "用户更新成功: ID=%d, 名称=%s", updatedUser.ID, updatedUser.Name)
	Success(c, updatedUser)
}

// Delete 删除用户
func (h *UserHandler) Delete(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	idStr := c.Param("id")
	logger.Logger.InfoWithContext(reqCtx, "开始删除用户: ID=%s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "无效的用户ID: %s, 错误: %v", idStr, err)
		BadRequest(c, "无效的用户ID")
		return
	}

	// 先获取用户信息，以便记录日志
	user, err := h.userService.GetByID(reqCtx, uint(id))
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "要删除的用户不存在: ID=%d, 错误: %v", id, err)
		NotFound(c, "用户不存在")
		return
	}

	logger.Logger.InfoWithContext(reqCtx, "开始删除用户: ID=%d, 名称=%s", id, user.Name)
	if err := h.userService.Delete(reqCtx, uint(id)); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "删除用户失败: %v", err)
		Fail(c, 500, err.Error())
		return
	}

	logger.Logger.InfoWithContext(reqCtx, "用户删除成功: ID=%d, 名称=%s", id, user.Name)
	Success(c, nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(ctx context.Context, c *app.RequestContext) {
	// 创建带请求ID的上下文
	reqCtx := middleware.WithRequestContext(ctx, c)

	idStr := c.Param("id")
	logger.Logger.InfoWithContext(reqCtx, "开始修改用户密码: ID=%s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Logger.WarnWithContext(reqCtx, "无效的用户ID: %s, 错误: %v", idStr, err)
		BadRequest(c, "无效的用户ID")
		return
	}

	// 获取请求参数
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&req); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "解析请求参数失败: %v", err)
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.OldPassword == "" || req.NewPassword == "" {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了空的密码")
		BadRequest(c, "旧密码和新密码不能为空")
		return
	}

	// 验证新密码强度
	if !utils.IsStrongPassword(req.NewPassword) {
		logger.Logger.WarnWithContext(reqCtx, "用户提交了弱密码")
		BadRequest(c, "密码强度不足，密码应包含大小写字母、数字和特殊字符，长度至少8位")
		return
	}

	// 修改密码
	logger.Logger.InfoWithContext(reqCtx, "开始修改用户密码: ID=%d", id)
	if err := h.userService.ChangePassword(reqCtx, uint(id), req.OldPassword, req.NewPassword); err != nil {
		logger.Logger.ErrorWithContext(reqCtx, "修改密码失败: %v", err)
		Fail(c, 500, err.Error())
		return
	}

	logger.Logger.InfoWithContext(reqCtx, "用户密码修改成功: ID=%d", id)
	Success(c, nil)
}
