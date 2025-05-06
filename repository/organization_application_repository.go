package repository

import (
	"context"
	"saas-account/config"
	"saas-account/model"
)

// OrganizationApplicationRepository 组织应用仓库接口
type OrganizationApplicationRepository interface {
	Create(ctx context.Context, app *model.OrganizationApplication) error
	GetByID(ctx context.Context, id int64) (*model.OrganizationApplication, error)
	GetByAppKey(ctx context.Context, appKey string) (*model.OrganizationApplication, error)
	GetByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]model.OrganizationApplication, int64, error)
	List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplication, int64, error)
	Update(ctx context.Context, app *model.OrganizationApplication) error
	Delete(ctx context.Context, id int64) error
}

// organizationApplicationRepository 组织应用仓库实现
type organizationApplicationRepository struct{}

// NewOrganizationApplicationRepository 创建组织应用仓库
func NewOrganizationApplicationRepository() OrganizationApplicationRepository {
	return &organizationApplicationRepository{}
}

// Create 创建组织应用
func (r *organizationApplicationRepository) Create(ctx context.Context, app *model.OrganizationApplication) error {
	return config.DB.WithContext(ctx).Create(app).Error
}

// GetByID 根据ID获取组织应用
func (r *organizationApplicationRepository) GetByID(ctx context.Context, id int64) (*model.OrganizationApplication, error) {
	var app model.OrganizationApplication
	err := config.DB.WithContext(ctx).First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByAppKey 根据AppKey获取组织应用
func (r *organizationApplicationRepository) GetByAppKey(ctx context.Context, appKey string) (*model.OrganizationApplication, error) {
	var app model.OrganizationApplication
	err := config.DB.WithContext(ctx).Where("app_key = ?", appKey).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByOrganization 根据组织ID获取应用列表
func (r *organizationApplicationRepository) GetByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	var apps []model.OrganizationApplication
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.OrganizationApplication{}).Where("organization_id = ?", orgID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Where("organization_id = ?", orgID).
		Offset(offset).Limit(pageSize).
		Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// List 获取应用列表
func (r *organizationApplicationRepository) List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	var apps []model.OrganizationApplication
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.OrganizationApplication{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// Update 更新组织应用
func (r *organizationApplicationRepository) Update(ctx context.Context, app *model.OrganizationApplication) error {
	return config.DB.WithContext(ctx).Save(app).Error
}

// Delete 删除组织应用（软删除）
func (r *organizationApplicationRepository) Delete(ctx context.Context, id int64) error {
	return config.DB.WithContext(ctx).Delete(&model.OrganizationApplication{}, id).Error
}
