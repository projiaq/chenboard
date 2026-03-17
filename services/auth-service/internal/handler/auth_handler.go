package handler

import (
	"encoding/json"
	"net/http"

	"iot-platform/pkg/errors"
	"iot-platform/pkg/response"
	"iot-platform/services/auth-service/internal/domain"
	"iot-platform/services/auth-service/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, "Invalid request body", err))
		return
	}

	// 如果没有提供租户ID，使用默认租户
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeUnauthorized, err.Error(), err))
		return
	}

	response.Success(w, resp)
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, "Invalid request body", err))
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeUnauthorized, err.Error(), err))
		return
	}

	response.Success(w, resp)
}

// GetMe 获取当前用户信息
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID（由中间件设置）
	userID := r.Context().Value("user_id").(string)

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeResourceNotFound, "User not found", err))
		return
	}

	response.Success(w, user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req domain.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, "Invalid request body", err))
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, &req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, err.Error(), err))
		return
	}

	response.Success(w, map[string]string{"message": "Password changed successfully"})
}
