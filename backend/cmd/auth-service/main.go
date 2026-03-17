package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"iot-platform/backend/internal/auth/handler"
	"iot-platform/backend/internal/auth/repository"
	"iot-platform/backend/internal/auth/service"
	"iot-platform/backend/pkg/config"
	"iot-platform/backend/pkg/database"
	"iot-platform/backend/pkg/jwt"
	"iot-platform/backend/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting auth-service...")

	// 连接数据库
	db, err := database.Connect(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", logger.Get().Error(err))
	}
	defer db.Close()

	logger.Info("Connected to database")

	// 创建 JWT 管理器
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// 初始化仓储
	userRepo := repository.NewUserRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, jwtManager)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)

	// 创建路由
	r := mux.NewRouter()

	// 健康检查
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API 路由
	api := r.PathPrefix("/api/v1").Subrouter()

	// 认证相关（无需认证）
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info(fmt.Sprintf("Server listening on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Get().Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Get().Error(err))
	}

	logger.Info("Server exited")
}
