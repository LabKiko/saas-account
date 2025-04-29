package repository

import (
	"context"
	"saas-account/config"
	"saas-account/model"
)

// OrganizationApplicationLimitRepository 组织应用限制仓库接口
type OrganizationApplicationLimitRepository interface {
	Create(ctx context.Context, limit *model.OrganizationApplicationLimit) error
	GetByID(ctx context.Context, id uint) (*model.OrganizationApplicationLimit, error)
	GetByApplicationID(ctx context.Context, appID uint) (*model.OrganizationApplicationLimit, error)
	List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplicationLimit, int64, error)
	Update(ctx context.Context, limit *model.OrganizationApplicationLimit) error
	Delete(ctx context.Context, id uint) error
}

// organizationApplicationLimitRepository 组织应用限制仓库实现
type organizationApplicationLimitRepository struct{}

// NewOrganizationApplicationLimitRepository 创建组织应用限制仓库
func NewOrganizationApplicationLimitRepository() OrganizationApplicationLimitRepository {
	return &organizationApplicationLimitRepository{}
}

// Create 创建组织应用限制
func (r *organizationApplicationLimitRepository) Create(ctx context.Context, limit *model.OrganizationApplicationLimit) error {
	return config.DB.WithContext(ctx).Create(limit).Error
}

// GetByID 根据ID获取组织应用限制
func (r *organizationApplicationLimitRepository) GetByID(ctx context.Context, id uint) (*model.OrganizationApplicationLimit, error) {
	var limit model.OrganizationApplicationLimit
	err := config.DB.WithContext(ctx).First(&limit, id).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

// GetByApplicationID 根据应用ID获取组织应用限制
func (r *organizationApplicationLimitRepository) GetByApplicationID(ctx context.Context, appID uint) (*model.OrganizationApplicationLimit, error) {
	var limit model.OrganizationApplicationLimit
	err := config.DB.WithContext(ctx).Where("application_id = ?", appID).First(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

// List 获取组织应用限制列表
func (r *organizationApplicationLimitRepository) List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplicationLimit, int64, error) {
	var limits []model.OrganizationApplicationLimit
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.OrganizationApplicationLimit{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&limits).Error
	if err != nil {
		return nil, 0, err
	}

	return limits, total, nil
}

// Update 更新组织应用限制
func (r *organizationApplicationLimitRepository) Update(ctx context.Context, limit *model.OrganizationApplicationLimit) error {
	return config.DB.WithContext(ctx).Save(limit).Error
}

// Delete 删除组织应用限制（软删除）
func (r *organizationApplicationLimitRepository) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&model.OrganizationApplicationLimit{}, id).Error
}
