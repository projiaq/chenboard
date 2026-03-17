package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		format  string
		wantErr bool
	}{
		{"json format", "info", "json", false},
		{"console format", "debug", "console", false},
		{"warn level", "warn", "json", false},
		{"error level", "error", "json", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.level, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				Sync()
			}
		})
	}
}

func TestGet(t *testing.T) {
	logger := Get()
	if logger == nil {
		t.Error("Get() returned nil")
	}
}

func TestLogFunctions(t *testing.T) {
	// 初始化日志
	if err := Init("info", "json"); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer Sync()

	// 测试各种日志函数（不会报错，只是确保不 panic）
	Debug("debug message", zap.String("key", "value"))
	Info("info message", zap.Int("count", 1))
	Warn("warn message")
	Error("error message")
}
