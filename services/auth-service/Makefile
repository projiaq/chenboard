.PHONY: build run test clean docker-build docker-run migrate-up migrate-down

# 变量
SERVICE_NAME=auth-service
DOCKER_IMAGE=iot-platform/$(SERVICE_NAME)
DOCKER_TAG=latest

# 构建
build:
	go build -o bin/$(SERVICE_NAME) ./cmd

# 运行
run:
	go run ./cmd

# 测试
test:
	go test -v -race -coverprofile=coverage.out ./...

# 测试覆盖率
coverage:
	go tool cover -html=coverage.out

# 代码检查
lint:
	golangci-lint run

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out

# Docker 构建
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker 运行
docker-run:
	docker run -p 8081:8081 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# 数据库迁移（向上）
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/iot_platform?sslmode=disable" up

# 数据库迁移（向下）
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/iot_platform?sslmode=disable" down

# 创建新的迁移文件
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# 安装依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...
