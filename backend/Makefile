.PHONY: help deps test test-coverage lint fmt build clean

# 默认目标
help:
	@echo "Available targets:"
	@echo "  deps           - Download dependencies"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  build          - Build all services"
	@echo "  clean          - Clean build artifacts"

# 下载依赖
deps:
	go mod download
	go mod tidy

# 运行测试
test:
	go test -v -race ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	golangci-lint run --timeout=5m ./...

# 格式化代码
fmt:
	go fmt ./...
	gofmt -s -w .

# 构建所有服务
build:
	@echo "Building auth-service..."
	CGO_ENABLED=0 go build -o bin/auth-service ./cmd/auth-service
	@echo "Build complete!"

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache -testcache
