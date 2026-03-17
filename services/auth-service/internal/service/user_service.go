package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"iot-platform/services/auth-service/internal/domain"
	"iot-platform/services/auth-service/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, tenantID string, req *domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, tenantID string, page, pageSize int) ([]*domain.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, tenantID string, req *domain.CreateUserRequest) (*domain.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(ctx, tenantID, req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		TenantID: tenantID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Status:   "active",
		Metadata: req.Metadata,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Metadata != nil {
		user.Metadata = req.Metadata
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, tenantID string, page, pageSize int) ([]*domain.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(ctx, tenantID, offset, pageSize)
}
