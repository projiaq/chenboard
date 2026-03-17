package jwt

import (
	"testing"
	"time"
)

func TestManager_GenerateAccessToken(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateAccessToken("user123", "tenant123", "testuser")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}
}

func TestManager_ValidateToken(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	// 生成令牌
	token, err := manager.GenerateAccessToken("user123", "tenant123", "testuser")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	// 验证令牌
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("claims.UserID = %s, want user123", claims.UserID)
	}

	if claims.TenantID != "tenant123" {
		t.Errorf("claims.TenantID = %s, want tenant123", claims.TenantID)
	}

	if claims.Username != "testuser" {
		t.Errorf("claims.Username = %s, want testuser", claims.Username)
	}
}

func TestManager_ValidateToken_Invalid(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid token", "invalid.token.here"},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tt.token)
			if err == nil {
				t.Error("ValidateToken() expected error, got nil")
			}
		})
	}
}

func TestManager_ValidateToken_Expired(t *testing.T) {
	// 创建一个过期时间很短的管理器
	manager := NewManager("test-secret", 1*time.Millisecond, 1*time.Millisecond)

	token, err := manager.GenerateAccessToken("user123", "tenant123", "testuser")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	// 等待令牌过期
	time.Sleep(10 * time.Millisecond)

	_, err = manager.ValidateToken(token)
	if err != ErrExpiredToken {
		t.Errorf("ValidateToken() error = %v, want %v", err, ErrExpiredToken)
	}
}

func TestManager_GenerateRefreshToken(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateRefreshToken("user123", "tenant123", "testuser")
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}

	// 验证刷新令牌
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("claims.UserID = %s, want user123", claims.UserID)
	}
}
