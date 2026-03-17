package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testdb")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want 9090", cfg.Server.Port)
	}

	if cfg.Database.Host != "testdb" {
		t.Errorf("Database.Host = %s, want testdb", cfg.Database.Host)
	}

	if cfg.JWT.Secret != "test-secret" {
		t.Errorf("JWT.Secret = %s, want test-secret", cfg.JWT.Secret)
	}
}

func TestGetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	if dsn := cfg.GetDSN(); dsn != expected {
		t.Errorf("GetDSN() = %s, want %s", dsn, expected)
	}
}

func TestGetEnvDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// 测试默认值
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %s, want 0.0.0.0", cfg.Server.Host)
	}

	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Server.ReadTimeout = %v, want 30s", cfg.Server.ReadTimeout)
	}

	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Database.MaxOpenConns = %d, want 25", cfg.Database.MaxOpenConns)
	}
}
