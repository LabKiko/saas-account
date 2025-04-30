package handler

import (
	"context"
	"saas-account/model"
	"saas-account/service"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// ApplicationUsageHandler 应用使用记录处理器
type ApplicationUsageHandler struct {
	usageService service.ApplicationUsageService
}

// NewApplicationUsageHandler 创建应用使用记录处理器
func NewApplicationUsageHandler(usageService service.ApplicationUsageService) *ApplicationUsageHandler {
	return &ApplicationUsageHandler{
		usageService: usageService,
	}
}

// Create 创建应用使用记录
func (h *ApplicationUsageHandler) Create(ctx context.Context, c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	var usage model.ApplicationUsage
	if err := c.BindJSON(&usage); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 设置应用ID
	usage.ApplicationId = appID

	// 验证必填字段
	if usage.UsageType == "" || usage.UsageAmount == 0 {
		BadRequest(c, "使用类型和使用量不能为空")
		return
	}

	// 创建使用记录
	if err := h.usageService.Create(ctx, &usage); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, usage)
}

// GetByID 根据ID获取应用使用记录
func (h *ApplicationUsageHandler) GetByID(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的使用记录ID")
		return
	}

	usage, err := h.usageService.GetByID(ctx, uint(id))
	if err != nil {
		NotFound(c, "使用记录不存在")
		return
	}

	Success(c, usage)
}

// List 获取应用使用记录列表
func (h *ApplicationUsageHandler) List(ctx context.Context, c *app.RequestContext) {
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

	// 获取日期范围参数
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time
	var usages []model.ApplicationUsage
	var total int64

	// 如果指定了日期范围，则按日期范围查询
	if startDateStr != "" && endDateStr != "" {
		var err1, err2 error
		startDate, err1 = time.Parse("2006-01-02", startDateStr)
		endDate, err2 = time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			BadRequest(c, "无效的日期格式，请使用YYYY-MM-DD格式")
			return
		}

		// 将结束日期设置为当天的最后一秒
		endDate = endDate.Add(24*time.Hour - time.Second)

		usages, total, err = h.usageService.GetByApplicationAndDateRange(ctx, uint(appID), startDate, endDate, page, pageSize)
	} else {
		// 否则查询所有记录
		usages, total, err = h.usageService.GetByApplication(ctx, uint(appID), page, pageSize)
	}

	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithPagination(c, usages, total, page, pageSize)
}

// GetSummary 获取应用使用统计摘要
func (h *ApplicationUsageHandler) GetSummary(ctx context.Context, c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取日期范围参数
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time

	// 如果未指定日期范围，则默认为过去30天
	if startDateStr == "" || endDateStr == "" {
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -30)
	} else {
		var err1, err2 error
		startDate, err1 = time.Parse("2006-01-02", startDateStr)
		endDate, err2 = time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			BadRequest(c, "无效的日期格式，请使用YYYY-MM-DD格式")
			return
		}

		// 将结束日期设置为当天的最后一秒
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	summary, err := h.usageService.GetSummaryByApplication(ctx, uint(appID), startDate, endDate)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c, map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"summary":    summary,
	})
}

// RecordAPIUsage 记录API使用
func (h *ApplicationUsageHandler) RecordAPIUsage(ctx context.Context, c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取请求参数
	var req struct {
		UserID *uint `json:"user_id"`
		Amount int64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.Amount <= 0 {
		BadRequest(c, "使用量必须大于0")
		return
	}

	// 检查API限制
	allowed, err := h.usageService.CheckAPILimit(ctx, uint(appID))
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	if !allowed {
		Forbidden(c, "已超过API使用限制")
		return
	}

	// 记录API使用
	if err := h.usageService.RecordAPIUsage(ctx, uint(appID), req.UserID, req.Amount); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// RecordStorageUsage 记录存储使用
func (h *ApplicationUsageHandler) RecordStorageUsage(ctx context.Context, c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取请求参数
	var req struct {
		UserID *uint `json:"user_id"`
		Amount int64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.Amount <= 0 {
		BadRequest(c, "使用量必须大于0")
		return
	}

	// 检查存储限制
	allowed, err := h.usageService.CheckStorageLimit(ctx, uint(appID), req.Amount)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	if !allowed {
		Forbidden(c, "已超过存储使用限制")
		return
	}

	// 记录存储使用
	if err := h.usageService.RecordStorageUsage(ctx, uint(appID), req.UserID, req.Amount); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}

// RecordFeatureUsage 记录功能使用
func (h *ApplicationUsageHandler) RecordFeatureUsage(ctx context.Context, c *app.RequestContext) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的应用ID")
		return
	}

	// 获取请求参数
	var req struct {
		UserID      *uint  `json:"user_id"`
		FeatureName string `json:"feature_name"`
		Amount      int64  `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.FeatureName == "" || req.Amount <= 0 {
		BadRequest(c, "功能名称和使用量不能为空且使用量必须大于0")
		return
	}

	// 记录功能使用
	if err := h.usageService.RecordFeatureUsage(ctx, uint(appID), req.UserID, req.FeatureName, req.Amount); err != nil {
		Fail(c, 500, err.Error())
		return
	}

	Success(c, nil)
}
