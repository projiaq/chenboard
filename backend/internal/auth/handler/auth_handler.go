package handler

import (
	"encoding/json"
	"net/http"

	"iot-platform/backend/internal/auth/domain"
	"iot-platform/backend/internal/auth/service"
	"iot-platform/backend/pkg/errors"
	"iot-platform/backend/pkg/response"
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
