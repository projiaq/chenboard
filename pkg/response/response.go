package response

import (
	"encoding/json"
	"net/http"
	"time"

	"iot-platform/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code      errors.ErrorCode `json:"code"`
	Message   string           `json:"message"`
	Data      interface{}      `json:"data,omitempty"`
	Timestamp int64            `json:"timestamp"`
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Code:      errors.CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	var statusCode int
	var code errors.ErrorCode
	var message string

	// 判断是否为应用错误
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
		code = e.Code
		message = e.Message
	} else {
		code = errors.CodeInternalError
		message = err.Error()
	}

	// 根据错误码确定 HTTP 状态码
	switch code {
	case errors.CodeInvalidParams, errors.CodeResourceExists:
		statusCode = http.StatusBadRequest
	case errors.CodeResourceNotFound:
		statusCode = http.StatusNotFound
	case errors.CodeUnauthorized, errors.CodeTokenExpired, errors.CodeInvalidToken:
		statusCode = http.StatusUnauthorized
	case errors.CodeForbidden:
		statusCode = http.StatusForbidden
	default:
		statusCode = http.StatusInternalServerError
	}

	JSON(w, statusCode, Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})

	// 记录错误
	if appErr != nil && appErr.Err != nil {
		// 这里可以添加日志记录
	}
}

// JSON 发送 JSON 响应
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// PageData 分页数据
type PageData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPageData 创建分页数据
func NewPageData(items interface{}, total int64, page, pageSize int) *PageData {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
