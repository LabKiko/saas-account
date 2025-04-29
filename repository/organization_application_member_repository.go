package repository

import (
	"context"
	"saas-account/config"
	"saas-account/model"
)

// OrganizationApplicationMemberRepository 组织应用成员仓库接口
type OrganizationApplicationMemberRepository interface {
	Create(ctx context.Context, member *model.OrganizationApplicationMember) error
	GetByID(ctx context.Context, id uint) (*model.OrganizationApplicationMember, error)
	GetByApplicationAndUser(ctx context.Context, appID, userID uint) (*model.OrganizationApplicationMember, error)
	GetByApplication(ctx context.Context, appID uint, page, pageSize int) ([]model.OrganizationApplicationMember, int64, error)
	GetByUser(ctx context.Context, userID uint) ([]model.OrganizationApplicationMember, error)
	Update(ctx context.Context, member *model.OrganizationApplicationMember) error
	Delete(ctx context.Context, id uint) error
	DeleteByApplicationAndUser(ctx context.Context, appID, userID uint) error
}

// organizationApplicationMemberRepository 组织应用成员仓库实现
type organizationApplicationMemberRepository struct{}

// NewOrganizationApplicationMemberRepository 创建组织应用成员仓库
func NewOrganizationApplicationMemberRepository() OrganizationApplicationMemberRepository {
	return &organizationApplicationMemberRepository{}
}

// Create 创建组织应用成员
func (r *organizationApplicationMemberRepository) Create(ctx context.Context, member *model.OrganizationApplicationMember) error {
	return config.DB.WithContext(ctx).Create(member).Error
}

// GetByID 根据ID获取组织应用成员
func (r *organizationApplicationMemberRepository) GetByID(ctx context.Context, id uint) (*model.OrganizationApplicationMember, error) {
	var member model.OrganizationApplicationMember
	err := config.DB.WithContext(ctx).First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByApplicationAndUser 根据应用ID和用户ID获取组织应用成员
func (r *organizationApplicationMemberRepository) GetByApplicationAndUser(ctx context.Context, appID, userID uint) (*model.OrganizationApplicationMember, error) {
	var member model.OrganizationApplicationMember
	err := config.DB.WithContext(ctx).Where("organization_application_id = ? AND user_id = ?", appID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByApplication 根据应用ID获取成员列表
func (r *organizationApplicationMemberRepository) GetByApplication(ctx context.Context, appID uint, page, pageSize int) ([]model.OrganizationApplicationMember, int64, error) {
	var members []model.OrganizationApplicationMember
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.OrganizationApplicationMember{}).Where("organization_application_id = ?", appID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Where("organization_application_id = ?", appID).
		Offset(offset).Limit(pageSize).
		Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// GetByUser 根据用户ID获取所属应用列表
func (r *organizationApplicationMemberRepository) GetByUser(ctx context.Context, userID uint) ([]model.OrganizationApplicationMember, error) {
	var members []model.OrganizationApplicationMember
	err := config.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// Update 更新组织应用成员
func (r *organizationApplicationMemberRepository) Update(ctx context.Context, member *model.OrganizationApplicationMember) error {
	return config.DB.WithContext(ctx).Save(member).Error
}

// Delete 删除组织应用成员（软删除）
func (r *organizationApplicationMemberRepository) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&model.OrganizationApplicationMember{}, id).Error
}

// DeleteByApplicationAndUser 根据应用ID和用户ID删除组织应用成员
func (r *organizationApplicationMemberRepository) DeleteByApplicationAndUser(ctx context.Context, appID, userID uint) error {
	return config.DB.WithContext(ctx).Where("organization_application_id = ? AND user_id = ?", appID, userID).
		Delete(&model.OrganizationApplicationMember{}).Error
}
