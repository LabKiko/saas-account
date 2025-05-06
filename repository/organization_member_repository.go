package repository

import (
	"context"
	"saas-account/config"
	"saas-account/model"
)

// OrganizationMemberRepository 组织成员仓库接口
type OrganizationMemberRepository interface {
	Create(ctx context.Context, member *model.OrganizationMember) error
	GetByID(ctx context.Context, id int64) (*model.OrganizationMember, error)
	GetByOrganizationAndUser(ctx context.Context, orgID, userID int64) (*model.OrganizationMember, error)
	GetByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]model.OrganizationMember, int64, error)
	GetByUser(ctx context.Context, userID int64) ([]model.OrganizationMember, error)
	Update(ctx context.Context, member *model.OrganizationMember) error
	Delete(ctx context.Context, id int64) error
	DeleteByOrganizationAndUser(ctx context.Context, orgID, userID int64) error
}

// organizationMemberRepository 组织成员仓库实现
type organizationMemberRepository struct{}

// NewOrganizationMemberRepository 创建组织成员仓库
func NewOrganizationMemberRepository() OrganizationMemberRepository {
	return &organizationMemberRepository{}
}

// Create 创建组织成员
func (r *organizationMemberRepository) Create(ctx context.Context, member *model.OrganizationMember) error {
	return config.DB.WithContext(ctx).Create(member).Error
}

// GetByID 根据ID获取组织成员
func (r *organizationMemberRepository) GetByID(ctx context.Context, id int64) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	err := config.DB.WithContext(ctx).First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByOrganizationAndUser 根据组织ID和用户ID获取组织成员
func (r *organizationMemberRepository) GetByOrganizationAndUser(ctx context.Context, orgID, userID int64) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	err := config.DB.WithContext(ctx).Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByOrganization 根据组织ID获取成员列表
func (r *organizationMemberRepository) GetByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]model.OrganizationMember, int64, error) {
	var members []model.OrganizationMember
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := config.DB.WithContext(ctx).Model(&model.OrganizationMember{}).Where("organization_id = ?", orgID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = config.DB.WithContext(ctx).Where("organization_id = ?", orgID).
		Offset(offset).Limit(pageSize).
		Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// GetByUser 根据用户ID获取所属组织列表
func (r *organizationMemberRepository) GetByUser(ctx context.Context, userID int64) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	err := config.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// Update 更新组织成员
func (r *organizationMemberRepository) Update(ctx context.Context, member *model.OrganizationMember) error {
	return config.DB.WithContext(ctx).Save(member).Error
}

// Delete 删除组织成员（软删除）
func (r *organizationMemberRepository) Delete(ctx context.Context, id int64) error {
	return config.DB.WithContext(ctx).Delete(&model.OrganizationMember{}, id).Error
}

// DeleteByOrganizationAndUser 根据组织ID和用户ID删除组织成员
func (r *organizationMemberRepository) DeleteByOrganizationAndUser(ctx context.Context, orgID, userID int64) error {
	return config.DB.WithContext(ctx).Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&model.OrganizationMember{}).Error
}
