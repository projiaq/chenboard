package domain

import "time"

// User 用户
type User struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Password    string                 `json:"-"`
	FullName    string                 `json:"full_name,omitempty"`
	Phone       string                 `json:"phone,omitempty"`
	Status      string                 `json:"status"`
	LastLoginAt *time.Time             `json:"last_login_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
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
