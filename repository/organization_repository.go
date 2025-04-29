package repository

import (
	"context"
	"saas-account/config"
	"saas-account/model"
)

// OrganizationRepository 组织仓库接口
type OrganizationRepository interface {
	Create(ctx context.Context, org *model.Organization) error
	GetByID(ctx context.Context, id uint) (*model.Organization, error)
	GetByOwnerID(ctx context.Context, ownerID uint) ([]model.Organization, error)
	List(ctx context.Context, page, pageSize int) ([]model.Organization, int64, error)
	Update(ctx context.Context, org *model.Organization) error
	Delete(ctx context.Context, id uint) error
}

// organizationRepository 组织仓库实现
type organizationRepository struct{}

// NewOrganizationRepository 创建组织仓库
func NewOrganizationRepository() OrganizationRepository {
	return &organizationRepository{}
}

// Create 创建组织
func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return config.DB.WithContext(ctx).Create(org).Error
}

// GetByID 根据ID获取组织
func (r *organizationRepository) GetByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	err := config.DB.WithContext(ctx).First(&org, id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByOwnerID 根据拥有者ID获取组织列表
func (r *organizationRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]model.Organization, error) {
	var orgs []model.Organization
	err := config.DB.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// List 获取组织列表
func (r *organizationRepository) List(ctx context.Context, page, pageSize int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.Organization{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&orgs).Error
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// Update 更新组织
func (r *organizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return config.DB.WithContext(ctx).Save(org).Error
}

// Delete 删除组织（软删除）
func (r *organizationRepository) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&model.Organization{}, id).Error
}
