package service

import (
	"context"
	"saas-account/model"
	"saas-account/repository"
	"time"
)

// ApplicationUsageService 应用使用记录服务接口
type ApplicationUsageService interface {
	Create(ctx context.Context, usage *model.ApplicationUsage) error
	GetByID(ctx context.Context, id int64) (*model.ApplicationUsage, error)
	GetByApplication(ctx context.Context, appID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetByApplicationAndDateRange(ctx context.Context, appID int64, startDate, endDate time.Time, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error)
	GetSummaryByApplication(ctx context.Context, appID int64, startDate, endDate time.Time) (map[string]int64, error)
	Delete(ctx context.Context, id int64) error
	RecordAPIUsage(ctx context.Context, appID int64, userID *int64, amount int64) error
	RecordStorageUsage(ctx context.Context, appID int64, userID *int64, amount int64) error
	RecordFeatureUsage(ctx context.Context, appID int64, userID *int64, featureName string, amount int64) error
	CheckAPILimit(ctx context.Context, appID int64) (bool, error)
	CheckStorageLimit(ctx context.Context, appID int64, additionalAmount int64) (bool, error)
}

// applicationUsageService 应用使用记录服务实现
type applicationUsageService struct {
	usageRepo repository.ApplicationUsageRepository
	appRepo   repository.OrganizationApplicationRepository
	limitRepo repository.OrganizationApplicationLimitRepository
}

// NewApplicationUsageService 创建应用使用记录服务
func NewApplicationUsageService(
	usageRepo repository.ApplicationUsageRepository,
	appRepo repository.OrganizationApplicationRepository,
	limitRepo repository.OrganizationApplicationLimitRepository,
) ApplicationUsageService {
	return &applicationUsageService{
		usageRepo: usageRepo,
		appRepo:   appRepo,
		limitRepo: limitRepo,
	}
}

// Create 创建应用使用记录
func (s *applicationUsageService) Create(ctx context.Context, usage *model.ApplicationUsage) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, usage.ApplicationId)
	if err != nil {
		return err
	}

	// 设置使用日期

	usage.UsageDate = time.Now().Unix()

	// 创建使用记录
	return s.usageRepo.Create(ctx, usage)
}

// GetByID 根据ID获取应用使用记录
func (s *applicationUsageService) GetByID(ctx context.Context, id int64) (*model.ApplicationUsage, error) {
	return s.usageRepo.GetByID(ctx, id)
}

// GetByApplication 根据应用ID获取使用记录
func (s *applicationUsageService) GetByApplication(ctx context.Context, appID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, 0, err
	}

	return s.usageRepo.GetByApplication(ctx, appID, page, pageSize)
}

// GetByApplicationAndDateRange 根据应用ID和日期范围获取使用记录
func (s *applicationUsageService) GetByApplicationAndDateRange(ctx context.Context, appID int64, startDate, endDate time.Time, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, 0, err
	}

	return s.usageRepo.GetByApplicationAndDateRange(ctx, appID, startDate, endDate, page, pageSize)
}

// GetByUser 根据用户ID获取使用记录
func (s *applicationUsageService) GetByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.ApplicationUsage, int64, error) {
	return s.usageRepo.GetByUser(ctx, userID, page, pageSize)
}

// GetSummaryByApplication 获取应用使用统计摘要
func (s *applicationUsageService) GetSummaryByApplication(ctx context.Context, appID int64, startDate, endDate time.Time) (map[string]int64, error) {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}

	return s.usageRepo.GetSummaryByApplication(ctx, appID, startDate, endDate)
}

// Delete 删除应用使用记录
func (s *applicationUsageService) Delete(ctx context.Context, id int64) error {
	return s.usageRepo.Delete(ctx, id)
}

// RecordAPIUsage 记录API使用
func (s *applicationUsageService) RecordAPIUsage(ctx context.Context, appID int64, userID *int64, amount int64) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 创建使用记录
	usage := &model.ApplicationUsage{
		ApplicationId: appID,
		UserId:        userID,
		UsageType:     "api_call",
		UsageAmount:   amount,
		UsageDate:     time.Now().Unix(),
		Details:       "{}",
	}

	return s.usageRepo.Create(ctx, usage)
}

// RecordStorageUsage 记录存储使用
func (s *applicationUsageService) RecordStorageUsage(ctx context.Context, appID int64, userID *int64, amount int64) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 创建使用记录
	usage := &model.ApplicationUsage{
		ApplicationId: appID,
		UserId:        userID,
		UsageType:     "storage",
		UsageAmount:   amount,
		UsageDate:     time.Now().Unix(),
		Details:       "{}",
	}

	return s.usageRepo.Create(ctx, usage)
}

// RecordFeatureUsage 记录功能使用
func (s *applicationUsageService) RecordFeatureUsage(ctx context.Context, appID int64, userID *int64, featureName string, amount int64) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 创建使用记录
	usage := &model.ApplicationUsage{
		ApplicationId: appID,
		UserId:        userID,
		UsageType:     "feature",
		UsageAmount:   amount,
		UsageDate:     time.Now().Unix(),
		Details:       "{\"feature_name\":\"" + featureName + "\"}",
	}

	return s.usageRepo.Create(ctx, usage)
}

// CheckAPILimit 检查API限制
func (s *applicationUsageService) CheckAPILimit(ctx context.Context, appID int64) (bool, error) {
	// 获取应用限制
	limit, err := s.limitRepo.GetByApplicationID(ctx, appID)
	if err != nil {
		return false, err
	}

	// 获取今日API使用量
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	summary, err := s.usageRepo.GetSummaryByApplication(ctx, appID, today, tomorrow)
	if err != nil {
		return false, err
	}

	// 检查是否超过限制
	apiUsage := summary["api_call"]
	return apiUsage < int64(limit.MaxRequests), nil
}

// CheckStorageLimit 检查存储限制
func (s *applicationUsageService) CheckStorageLimit(ctx context.Context, appID int64, additionalAmount int64) (bool, error) {
	// 获取应用限制
	limit, err := s.limitRepo.GetByApplicationID(ctx, appID)
	if err != nil {
		return false, err
	}

	// 获取存储使用量
	summary, err := s.usageRepo.GetSummaryByApplication(ctx, appID, time.Time{}, time.Now())
	if err != nil {
		return false, err
	}

	// 检查是否超过限制
	storageUsage := summary["storage"]
	return storageUsage+additionalAmount <= limit.MaxStorage, nil
}
