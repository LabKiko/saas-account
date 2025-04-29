package repository

import (
	"context"
	"time"

	"saas-account/config"
	"saas-account/model"
)

// ApplicationUsageRepository 应用使用记录仓库接口
type ApplicationUsageRepository interface {
	Create(ctx context.Context, usage *model.ApplicationUsage) error
	GetByID(ctx context.Context, id int64) (*model.ApplicationUsage, error)
	GetByApplication(ctx context.Context, appID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetByApplicationAndDateRange(ctx context.Context, appID int64, startDate, endDate time.Time, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetSummaryByApplication(ctx context.Context, appID int64, startDate, endDate time.Time) (map[string]int64, error)
	Delete(ctx context.Context, id int64) error
}

// applicationUsageRepository 应用使用记录仓库实现
type applicationUsageRepository struct{}

// NewApplicationUsageRepository 创建应用使用记录仓库
func NewApplicationUsageRepository() ApplicationUsageRepository {
	return &applicationUsageRepository{}
}

// Create 创建应用使用记录
func (r *applicationUsageRepository) Create(ctx context.Context, usage *model.ApplicationUsage) error {
	return config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Create(usage).Error
}

// GetByID 根据ID获取应用使用记录
func (r *applicationUsageRepository) GetByID(ctx context.Context, id int64) (*model.ApplicationUsage, error) {
	var usage model.ApplicationUsage
	err := config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).First(&usage, id).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetByApplication 根据应用ID获取使用记录
func (r *applicationUsageRepository) GetByApplication(ctx context.Context, appID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	var usages []model.ApplicationUsage
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Where("application_id = ?", appID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Where("application_id = ?", appID).
		Order("usage_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&usages).Error
	if err != nil {
		return nil, 0, err
	}

	return usages, total, nil
}

// GetByApplicationAndDateRange 根据应用ID和日期范围获取使用记录
func (r *applicationUsageRepository) GetByApplicationAndDateRange(ctx context.Context, appID int64, startDate, endDate time.Time, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	var usages []model.ApplicationUsage
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).
		Where("application_id = ? AND usage_date BETWEEN ? AND ?", appID, startDate, endDate).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Where("application_id = ? AND usage_date BETWEEN ? AND ?", appID, startDate, endDate).
		Order("usage_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&usages).Error
	if err != nil {
		return nil, 0, err
	}

	return usages, total, nil
}

// GetByUser 根据用户ID获取使用记录
func (r *applicationUsageRepository) GetByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	var usages []model.ApplicationUsage
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.Model(&model.ApplicationUsage{}).
		WithContext(ctx).Where("user_id = ?", userID).
		Order("usage_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&usages).Error
	if err != nil {
		return nil, 0, err
	}

	return usages, total, nil
}

// GetSummaryByApplication 获取应用使用统计摘要
func (r *applicationUsageRepository) GetSummaryByApplication(ctx context.Context, appID int64, startDate, endDate time.Time) (map[string]int64, error) {
	type Result struct {
		UsageType   string
		TotalAmount int64
	}

	var results []Result
	err := config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).
		Select("usage_type, SUM(usage_amount) as total_amount").
		Where("application_id = ? AND usage_date BETWEEN ? AND ?", appID, startDate, endDate).
		Group("usage_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]int64)
	for _, result := range results {
		summary[result.UsageType] = result.TotalAmount
	}

	return summary, nil
}

// Delete 删除应用使用记录（软删除）
func (r *applicationUsageRepository) Delete(ctx context.Context, id int64) error {
	return config.DB.Model(&model.ApplicationUsage{}).WithContext(ctx).Delete(&model.ApplicationUsage{}, id).Error
}
