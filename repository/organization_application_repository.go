package repository

import (
	"saas-account/config"
	"saas-account/model"
)

// OrganizationApplicationRepository 组织应用仓库接口
type OrganizationApplicationRepository interface {
	Create(app *model.OrganizationApplication) error
	GetByID(id uint) (*model.OrganizationApplication, error)
	GetByAppKey(appKey string) (*model.OrganizationApplication, error)
	GetByOrganization(orgID uint, page, pageSize int) ([]model.OrganizationApplication, int64, error)
	List(page, pageSize int) ([]model.OrganizationApplication, int64, error)
	Update(app *model.OrganizationApplication) error
	Delete(id uint) error
}

// organizationApplicationRepository 组织应用仓库实现
type organizationApplicationRepository struct{}

// NewOrganizationApplicationRepository 创建组织应用仓库
func NewOrganizationApplicationRepository() OrganizationApplicationRepository {
	return &organizationApplicationRepository{}
}

// Create 创建组织应用
func (r *organizationApplicationRepository) Create(app *model.OrganizationApplication) error {
	return config.DB.Create(app).Error
}

// GetByID 根据ID获取组织应用
func (r *organizationApplicationRepository) GetByID(id uint) (*model.OrganizationApplication, error) {
	var app model.OrganizationApplication
	err := config.DB.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByAppKey 根据AppKey获取组织应用
func (r *organizationApplicationRepository) GetByAppKey(appKey string) (*model.OrganizationApplication, error) {
	var app model.OrganizationApplication
	err := config.DB.Where("app_key = ?", appKey).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByOrganization 根据组织ID获取应用列表
func (r *organizationApplicationRepository) GetByOrganization(orgID uint, page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	var apps []model.OrganizationApplication
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.Model(&model.OrganizationApplication{}).Where("organization_id = ?", orgID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.Where("organization_id = ?", orgID).
		Offset(offset).Limit(pageSize).
		Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// List 获取应用列表
func (r *organizationApplicationRepository) List(page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	var apps []model.OrganizationApplication
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.Model(&model.OrganizationApplication{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.Offset(offset).Limit(pageSize).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// Update 更新组织应用
func (r *organizationApplicationRepository) Update(app *model.OrganizationApplication) error {
	return config.DB.Save(app).Error
}

// Delete 删除组织应用（软删除）
func (r *organizationApplicationRepository) Delete(id uint) error {
	return config.DB.Delete(&model.OrganizationApplication{}, id).Error
}
