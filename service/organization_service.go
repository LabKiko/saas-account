package service

import (
	"context"
	"errors"
	"saas-account/model"
	"saas-account/repository"
)

// OrganizationService 组织服务接口
type OrganizationService interface {
	Create(ctx context.Context, org *model.Organization, creatorID uint) error
	GetByID(ctx context.Context, id uint) (*model.Organization, error)
	GetByOwnerID(ctx context.Context, ownerID uint) ([]model.Organization, error)
	List(ctx context.Context, page, pageSize int) ([]model.Organization, int64, error)
	Update(ctx context.Context, org *model.Organization) error
	Delete(ctx context.Context, id uint) error
	AddMember(ctx context.Context, orgID, userID uint, role string) error
	RemoveMember(ctx context.Context, orgID, userID uint) error
	GetMembers(ctx context.Context, orgID uint, page, pageSize int) ([]model.OrganizationMember, int64, error)
	UpdateMember(ctx context.Context, orgID, userID uint, role string) error
}

// organizationService 组织服务实现
type organizationService struct {
	orgRepo     repository.OrganizationRepository
	orgMemberRepo repository.OrganizationMemberRepository
	userRepo    repository.UserRepository
}

// NewOrganizationService 创建组织服务
func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	orgMemberRepo repository.OrganizationMemberRepository,
	userRepo repository.UserRepository,
) OrganizationService {
	return &organizationService{
		orgRepo:     orgRepo,
		orgMemberRepo: orgMemberRepo,
		userRepo:    userRepo,
	}
}

// Create 创建组织
func (s *organizationService) Create(ctx context.Context, org *model.Organization, creatorID uint) error {
	// 检查创建者是否存在
	creator, err := s.userRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}

	// 设置组织拥有者
	org.OwnerId = creator.ID

	// 设置默认状态
	if org.Status == "" {
		org.Status = "active"
	}

	// 创建组织
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return err
	}

	// 添加创建者为组织成员（拥有者角色）
	member := &model.OrganizationMember{
		OrganizationId: org.ID,
		UserId:         creator.ID,
		Role:           "owner",
		Status:         "active",
	}

	return s.orgMemberRepo.Create(ctx, member)
}

// GetByID 根据ID获取组织
func (s *organizationService) GetByID(ctx context.Context, id uint) (*model.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

// GetByOwnerID 根据拥有者ID获取组织列表
func (s *organizationService) GetByOwnerID(ctx context.Context, ownerID uint) ([]model.Organization, error) {
	return s.orgRepo.GetByOwnerID(ctx, ownerID)
}

// List 获取组织列表
func (s *organizationService) List(ctx context.Context, page, pageSize int) ([]model.Organization, int64, error) {
	return s.orgRepo.List(ctx, page, pageSize)
}

// Update 更新组织
func (s *organizationService) Update(ctx context.Context, org *model.Organization) error {
	// 获取现有组织
	existingOrg, err := s.orgRepo.GetByID(ctx, org.ID)
	if err != nil {
		return err
	}

	// 保留拥有者ID
	org.OwnerId = existingOrg.OwnerId

	// 更新组织
	return s.orgRepo.Update(ctx, org)
}

// Delete 删除组织
func (s *organizationService) Delete(ctx context.Context, id uint) error {
	return s.orgRepo.Delete(ctx, id)
}

// AddMember 添加组织成员
func (s *organizationService) AddMember(ctx context.Context, orgID, userID uint, role string) error {
	// 检查组织是否存在
	_, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// 检查用户是否存在
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 检查用户是否已经是组织成员
	existingMember, err := s.orgMemberRepo.GetByOrganizationAndUser(ctx, orgID, userID)
	if err == nil && existingMember != nil {
		return errors.New("用户已经是组织成员")
	}

	// 验证角色
	if role != "admin" && role != "member" && role != "owner" {
		return errors.New("无效的角色")
	}

	// 添加组织成员
	member := &model.OrganizationMember{
		OrganizationId: orgID,
		UserId:         userID,
		Role:           role,
		Status:         "active",
	}

	return s.orgMemberRepo.Create(ctx, member)
}

// RemoveMember 移除组织成员
func (s *organizationService) RemoveMember(ctx context.Context, orgID, userID uint) error {
	// 检查组织是否存在
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// 不能移除组织拥有者
	if org.OwnerId == userID {
		return errors.New("不能移除组织拥有者")
	}

	// 移除组织成员
	return s.orgMemberRepo.DeleteByOrganizationAndUser(ctx, orgID, userID)
}

// GetMembers 获取组织成员列表
func (s *organizationService) GetMembers(ctx context.Context, orgID uint, page, pageSize int) ([]model.OrganizationMember, int64, error) {
	// 检查组织是否存在
	_, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}

	return s.orgMemberRepo.GetByOrganization(ctx, orgID, page, pageSize)
}

// UpdateMember 更新组织成员
func (s *organizationService) UpdateMember(ctx context.Context, orgID, userID uint, role string) error {
	// 检查组织是否存在
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// 检查成员是否存在
	member, err := s.orgMemberRepo.GetByOrganizationAndUser(ctx, orgID, userID)
	if err != nil {
		return err
	}

	// 不能更改组织拥有者的角色
	if org.OwnerId == userID && role != "owner" {
		return errors.New("不能更改组织拥有者的角色")
	}

	// 验证角色
	if role != "admin" && role != "member" && role != "owner" {
		return errors.New("无效的角色")
	}

	// 更新成员角色
	member.Role = role
	return s.orgMemberRepo.Update(ctx, member)
}
