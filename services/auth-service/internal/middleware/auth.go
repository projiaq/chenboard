package middleware

import (
	"context"
	"net/http"
	"strings"

	"iot-platform/pkg/errors"
	"iot-platform/pkg/jwt"
	"iot-platform/pkg/response"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取 Authorization 头
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, errors.ErrUnauthorized)
				return
			}

			// 解析 Bearer Token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, errors.ErrInvalidToken)
				return
			}

			tokenString := parts[1]

			// 验证令牌
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					response.Error(w, errors.ErrTokenExpired)
				} else {
					response.Error(w, errors.ErrInvalidToken)
				}
				return
			}

			// 将用户信息存入上下文
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "tenant_id", claims.TenantID)
			ctx = context.WithValue(ctx, "username", claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
