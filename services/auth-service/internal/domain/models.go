package domain

import (
	"time"
)

// Tenant 租户
type Tenant struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Code       string                 `json:"code"`
	Status     string                 `json:"status"`
	MaxDevices int                    `json:"max_devices"`
	MaxUsers   int                    `json:"max_users"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// User 用户
type User struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Password    string                 `json:"-"` // 不序列化密码
	FullName    string                 `json:"full_name,omitempty"`
	Phone       string                 `json:"phone,omitempty"`
	Status      string                 `json:"status"`
	LastLoginAt *time.Time             `json:"last_login_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Role 角色
type Role struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission 权限
type Permission struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	FullName string                 `json:"full_name,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	FullName string                 `json:"full_name,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
