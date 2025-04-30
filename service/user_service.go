package service

import (
	"context"
	"errors"
	"saas-account/model"
	"saas-account/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, user *model.User) error {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("邮箱已被注册")
	}

	// 检查手机号是否已存在
	if user.Phone != "" {
		existingUser, err = s.userRepo.GetByPhone(ctx, user.Phone)
		if err == nil && existingUser != nil {
			return errors.New("手机号已被注册")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置默认状态
	if user.Status == "" {
		user.Status = "active"
	}

	user.ID =

	// 创建用户
	return s.userRepo.Create(ctx, user)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail 根据邮箱获取用户
func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetByPhone 根据手机号获取用户
func (s *userService) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return s.userRepo.GetByPhone(ctx, phone)
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize)
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, user *model.User) error {
	// 获取现有用户
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// 检查邮箱是否已被其他用户使用
	if user.Email != existingUser.Email {
		otherUser, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && otherUser != nil && otherUser.ID != user.ID {
			return errors.New("邮箱已被其他用户注册")
		}
	}

	// 检查手机号是否已被其他用户使用
	if user.Phone != "" && user.Phone != existingUser.Phone {
		otherUser, err := s.userRepo.GetByPhone(ctx, user.Phone)
		if err == nil && otherUser != nil && otherUser.ID != user.ID {
			return errors.New("手机号已被其他用户注册")
		}
	}

	// 保留原密码
	user.Password = existingUser.Password

	// 更新用户
	return s.userRepo.Update(ctx, user)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码不正确")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now().Unix()

	return s.userRepo.Update(ctx, user)
}
