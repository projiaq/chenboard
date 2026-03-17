package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"iot-platform/pkg/errors"
	"iot-platform/pkg/response"
	"iot-platform/services/auth-service/internal/domain"
	"iot-platform/services/auth-service/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create 创建用户
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, "Invalid request body", err))
		return
	}

	user, err := h.userService.Create(r.Context(), tenantID, &req)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, err.Error(), err))
		return
	}

	response.Success(w, user)
}

// GetByID 获取用户详情
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeResourceNotFound, "User not found", err))
		return
	}

	response.Success(w, user)
}

// Update 更新用户
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, "Invalid request body", err))
		return
	}

	user, err := h.userService.Update(r.Context(), id, &req)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeInvalidParams, err.Error(), err))
		return
	}

	response.Success(w, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.userService.Delete(r.Context(), id); err != nil {
		response.Error(w, errors.Wrap(errors.CodeInternalError, "Failed to delete user", err))
		return
	}

	response.Success(w, map[string]string{"message": "User deleted successfully"})
}

// List 获取用户列表
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.userService.List(r.Context(), tenantID, page, pageSize)
	if err != nil {
		response.Error(w, errors.Wrap(errors.CodeInternalError, "Failed to list users", err))
		return
	}

	pageData := response.NewPageData(users, total, page, pageSize)
	response.Success(w, pageData)
}
