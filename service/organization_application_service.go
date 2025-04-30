package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"saas-account/model"
	"saas-account/repository"
	"time"
)

// OrganizationApplicationService 组织应用服务接口
type OrganizationApplicationService interface {
	Create(ctx context.Context, app *model.OrganizationApplication) error
	GetByID(ctx context.Context, id uint) (*model.OrganizationApplication, error)
	GetByAppKey(ctx context.Context, appKey string) (*model.OrganizationApplication, error)
	GetByOrganization(ctx context.Context, orgID uint, page, pageSize int) ([]model.OrganizationApplication, int64, error)
	List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplication, int64, error)
	Update(ctx context.Context, app *model.OrganizationApplication) error
	Delete(ctx context.Context, id uint) error
	RegenerateAppSecret(ctx context.Context, id uint) (string, error)
	AddMember(ctx context.Context, appID, userID uint, role string, permissions string) error
	RemoveMember(ctx context.Context, appID, userID uint) error
	GetMembers(ctx context.Context, appID uint, page, pageSize int) ([]model.OrganizationApplicationMember, int64, error)
	UpdateMember(ctx context.Context, appID, userID uint, role string, permissions string) error
	SetLimit(ctx context.Context, limit *model.OrganizationApplicationLimit) error
	GetLimit(ctx context.Context, appID uint) (*model.OrganizationApplicationLimit, error)
}

// organizationApplicationService 组织应用服务实现
type organizationApplicationService struct {
	appRepo     repository.OrganizationApplicationRepository
	appMemberRepo repository.OrganizationApplicationMemberRepository
	appLimitRepo repository.OrganizationApplicationLimitRepository
	orgRepo     repository.OrganizationRepository
	userRepo    repository.UserRepository
}

// NewOrganizationApplicationService 创建组织应用服务
func NewOrganizationApplicationService(
	appRepo repository.OrganizationApplicationRepository,
	appMemberRepo repository.OrganizationApplicationMemberRepository,
	appLimitRepo repository.OrganizationApplicationLimitRepository,
	orgRepo repository.OrganizationRepository,
	userRepo repository.UserRepository,
) OrganizationApplicationService {
	return &organizationApplicationService{
		appRepo:     appRepo,
		appMemberRepo: appMemberRepo,
		appLimitRepo: appLimitRepo,
		orgRepo:     orgRepo,
		userRepo:    userRepo,
	}
}

// generateAppKey 生成应用密钥
func generateAppKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateAppSecret 生成应用密钥
func generateAppSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create 创建组织应用
func (s *organizationApplicationService) Create(ctx context.Context, app *model.OrganizationApplication) error {
	// 检查组织是否存在
	_, err := s.orgRepo.GetByID(ctx, app.OrganizationId)
	if err != nil {
		return err
	}

	// 生成应用密钥
	appKey, err := generateAppKey()
	if err != nil {
		return err
	}
	app.AppKey = appKey

	// 生成应用密钥
	appSecret, err := generateAppSecret()
	if err != nil {
		return err
	}
	app.AppSecret = appSecret

	// 设置默认状态
	if app.Status == "" {
		app.Status = "active"
	}

	// 创建应用
	if err := s.appRepo.Create(ctx, app); err != nil {
		return err
	}

	// 创建默认应用限制
	limit := &model.OrganizationApplicationLimit{
		OrganizationApplicationId: app.ID,
		PlanName:                 "free",
		MaxUsers:                 5,
		MaxStorage:               1073741824, // 1GB
		MaxRequests:              10000,      // 每天10000请求
		Features:                 "{}",       // 空特性
		ExpiresAt:                nil,        // 永不过期
		AutoRenew:                false,
	}

	return s.appLimitRepo.Create(ctx, limit)
}

// GetByID 根据ID获取组织应用
func (s *organizationApplicationService) GetByID(ctx context.Context, id uint) (*model.OrganizationApplication, error) {
	return s.appRepo.GetByID(ctx, id)
}

// GetByAppKey 根据AppKey获取组织应用
func (s *organizationApplicationService) GetByAppKey(ctx context.Context, appKey string) (*model.OrganizationApplication, error) {
	return s.appRepo.GetByAppKey(ctx, appKey)
}

// GetByOrganization 根据组织ID获取应用列表
func (s *organizationApplicationService) GetByOrganization(ctx context.Context, orgID uint, page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	// 检查组织是否存在
	_, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}

	return s.appRepo.GetByOrganization(ctx, orgID, page, pageSize)
}

// List 获取应用列表
func (s *organizationApplicationService) List(ctx context.Context, page, pageSize int) ([]model.OrganizationApplication, int64, error) {
	return s.appRepo.List(ctx, page, pageSize)
}

// Update 更新组织应用
func (s *organizationApplicationService) Update(ctx context.Context, app *model.OrganizationApplication) error {
	// 获取现有应用
	existingApp, err := s.appRepo.GetByID(ctx, app.ID)
	if err != nil {
		return err
	}

	// 保留不可修改的字段
	app.OrganizationId = existingApp.OrganizationId
	app.AppKey = existingApp.AppKey
	app.AppSecret = existingApp.AppSecret

	// 更新应用
	return s.appRepo.Update(ctx, app)
}

// Delete 删除组织应用
func (s *organizationApplicationService) Delete(ctx context.Context, id uint) error {
	return s.appRepo.Delete(ctx, id)
}

// RegenerateAppSecret 重新生成应用密钥
func (s *organizationApplicationService) RegenerateAppSecret(ctx context.Context, id uint) (string, error) {
	// 获取现有应用
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	// 生成新的应用密钥
	appSecret, err := generateAppSecret()
	if err != nil {
		return "", err
	}
	app.AppSecret = appSecret

	// 更新应用
	if err := s.appRepo.Update(ctx, app); err != nil {
		return "", err
	}

	return appSecret, nil
}

// AddMember 添加应用成员
func (s *organizationApplicationService) AddMember(ctx context.Context, appID, userID uint, role string, permissions string) error {
	// 检查应用是否存在
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 检查用户是否存在
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 检查用户是否已经是应用成员
	existingMember, err := s.appMemberRepo.GetByApplicationAndUser(ctx, appID, userID)
	if err == nil && existingMember != nil {
		return errors.New("用户已经是应用成员")
	}

	// 验证角色
	if role != "admin" && role != "user" && role != "guest" {
		return errors.New("无效的角色")
	}

	// 检查应用成员数量限制
	limit, err := s.appLimitRepo.GetByApplicationID(ctx, appID)
	if err == nil && limit != nil {
		members, total, err := s.appMemberRepo.GetByApplication(ctx, appID, 1, 0)
		if err == nil && total >= int64(limit.MaxUsers) {
			return errors.New("已达到应用成员数量限制")
		}
	}

	// 添加应用成员
	member := &model.OrganizationApplicationMember{
		OrganizationApplicationId: appID,
		UserId:                   userID,
		Role:                     role,
		Status:                   "active",
		Permissions:              permissions,
	}

	return s.appMemberRepo.Create(ctx, member)
}

// RemoveMember 移除应用成员
func (s *organizationApplicationService) RemoveMember(ctx context.Context, appID, userID uint) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 移除应用成员
	return s.appMemberRepo.DeleteByApplicationAndUser(ctx, appID, userID)
}

// GetMembers 获取应用成员列表
func (s *organizationApplicationService) GetMembers(ctx context.Context, appID uint, page, pageSize int) ([]model.OrganizationApplicationMember, int64, error) {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, 0, err
	}

	return s.appMemberRepo.GetByApplication(ctx, appID, page, pageSize)
}

// UpdateMember 更新应用成员
func (s *organizationApplicationService) UpdateMember(ctx context.Context, appID, userID uint, role string, permissions string) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	// 检查成员是否存在
	member, err := s.appMemberRepo.GetByApplicationAndUser(ctx, appID, userID)
	if err != nil {
		return err
	}

	// 验证角色
	if role != "admin" && role != "user" && role != "guest" {
		return errors.New("无效的角色")
	}

	// 更新成员
	member.Role = role
	member.Permissions = permissions
	return s.appMemberRepo.Update(ctx, member)
}

// SetLimit 设置应用限制
func (s *organizationApplicationService) SetLimit(ctx context.Context, limit *model.OrganizationApplicationLimit) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, limit.OrganizationApplicationId)
	if err != nil {
		return err
	}

	// 检查是否已存在限制
	existingLimit, err := s.appLimitRepo.GetByApplicationID(ctx, limit.OrganizationApplicationId)
	if err == nil && existingLimit != nil {
		// 更新现有限制
		limit.ID = existingLimit.ID
		return s.appLimitRepo.Update(ctx, limit)
	}

	// 创建新限制
	return s.appLimitRepo.Create(ctx, limit)
}

// GetLimit 获取应用限制
func (s *organizationApplicationService) GetLimit(ctx context.Context, appID uint) (*model.OrganizationApplicationLimit, error) {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}

	return s.appLimitRepo.GetByApplicationID(ctx, appID)
}
