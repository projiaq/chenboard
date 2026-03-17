# IoT Platform - CI/CD Ready Monorepo

一个完全适配 GitHub Actions 的物联网平台项目，所有代码都可以在 CI 环境中直接运行。

## 🎯 工程约束保证

✅ 所有代码都能在 GitHub Actions Ubuntu runner 中直接执行
✅ 不依赖本地 IDE、手工安装、手工改路径
✅ 提供完整的构建和测试命令
✅ 所有依赖版本明确，避免 CI 不稳定
✅ 优先保证可构建、可测试、可打包

## 📁 项目结构

```
iot-platform/
├── backend/                    # Go 后端（单一模块）
│   ├── cmd/                    # 服务入口
│   │   └── auth-service/
│   ├── internal/               # 私有代码
│   │   └── auth/
│   ├── pkg/                    # 共享库
│   ├── migrations/             # 数据库迁移
│   ├── go.mod                  # Go 依赖
│   ├── Makefile                # 构建脚本
│   └── .golangci.yml           # Lint 配置
├── deployments/                # 部署配置
│   ├── docker/                 # Dockerfile
│   └── docker-compose/         # Docker Compose
├── scripts/                    # 自动化脚本
│   ├── build/
│   ├── test/
│   └── db/
└── .github/workflows/          # CI/CD 配置
    └── backend-ci.yml
```

## 🚀 快速开始

### 本地开发

#### 1. 前置要求

- Go 1.23.0+
- Docker & Docker Compose
- Make

#### 2. 启动基础设施

```bash
cd deployments/docker-compose
docker-compose up -d postgres redis
```

#### 3. 运行数据库迁移

```bash
# 方式 1: 使用脚本
chmod +x scripts/db/migrate-up.sh
./scripts/db/migrate-up.sh

# 方式 2: 手动执行
cd backend/migrations/auth
migrate -path . -database "postgresql://postgres:postgres@localhost:5432/iot_platform?sslmode=disable" up
```

#### 4. 启动服务

```bash
cd backend
make run
```

或者直接运行：

```bash
cd backend
go run ./cmd/auth-service
```

#### 5. 测试 API

```bash
# 健康检查
curl http://localhost:8081/health

# 登录（默认账号：admin/admin123）
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123",
    "tenant_id": "default"
  }'
```

### 使用 Docker Compose

```bash
# 启动所有服务
cd deployments/docker-compose
docker-compose up -d

# 查看日志
docker-compose logs -f auth-service

# 停止服务
docker-compose down
```

## 🔨 构建和测试

### 后端

```bash
cd backend

# 下载依赖
make deps

# 运行测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 代码检查
make lint

# 格式化代码
make fmt

# 构建所有服务
make build

# 清理
make clean
```

### 使用脚本

```bash
# 构建后端
chmod +x scripts/build/build-backend.sh
./scripts/build/build-backend.sh

# 测试后端
chmod +x scripts/test/test-backend.sh
./scripts/test/test-backend.sh
```

## 🐳 Docker 构建

### 构建单个服务

```bash
# 从项目根目录执行
docker build -f deployments/docker/auth-service.Dockerfile -t iot-platform/auth-service:latest .
```

### 运行容器

```bash
docker run -p 8081:8081 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=iot_platform \
  -e JWT_SECRET=your-secret-key \
  iot-platform/auth-service:latest
```

## 🔄 GitHub Actions CI/CD

### Backend CI 工作流

触发条件：
- Push 到 `main` 或 `develop` 分支
- Pull Request 到 `main` 或 `develop` 分支
- 修改了 `backend/**` 路径下的文件

执行步骤：
1. **Lint**: 代码检查（golangci-lint）
2. **Test**: 单元测试（带覆盖率）
3. **Build**: 构建所有服务

查看配置：[.github/workflows/backend-ci.yml](.github/workflows/backend-ci.yml)

### 本地模拟 CI 环境

```bash
# 模拟 CI 环境运行测试
cd backend
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 模拟 CI 环境构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/auth-service ./cmd/auth-service
```

## 📦 依赖版本

### 后端依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| Go | 1.23.0 | 编程语言 |
| github.com/golang-jwt/jwt/v5 | v5.2.1 | JWT 认证 |
| github.com/google/uuid | v1.6.0 | UUID 生成 |
| github.com/gorilla/mux | v1.8.1 | HTTP 路由 |
| github.com/lib/pq | v1.10.9 | PostgreSQL 驱动 |
| github.com/redis/go-redis/v9 | v9.5.1 | Redis 客户端 |
| go.uber.org/zap | v1.27.0 | 结构化日志 |
| golang.org/x/crypto | v0.31.0 | 密码加密 |

### 基础设施版本

| 服务 | 版本 | 镜像 |
|------|------|------|
| PostgreSQL | 16.1 | postgres:16.1-alpine |
| Redis | 7.2 | redis:7.2-alpine |
| Alpine Linux | 3.19 | alpine:3.19 |

## 🧪 测试

### 运行所有测试

```bash
cd backend
go test ./...
```

### 运行特定包的测试

```bash
cd backend
go test ./pkg/jwt
go test ./pkg/config
go test ./pkg/errors
```

### 查看测试覆盖率

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 环境变量

复制 `.env.example` 并修改：

```bash
cp backend/.env.example backend/.env
```

主要配置项：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 8081 | 服务端口 |
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 5432 | 数据库端口 |
| DB_USER | postgres | 数据库用户 |
| DB_PASSWORD | postgres | 数据库密码 |
| DB_NAME | iot_platform | 数据库名称 |
| JWT_SECRET | change-me | JWT 密钥（生产环境必须修改） |
| LOG_LEVEL | info | 日志级别 |

## 🔐 默认账号

系统初始化后会创建默认管理员账号：

- **用户名**: admin
- **密码**: admin123
- **租户**: default

⚠️ **生产环境请立即修改默认密码！**

## 📊 API 文档

### 认证相关

#### 登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123",
  "tenant_id": "default"
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "00000000-0000-0000-0000-000000000001",
      "tenant_id": "00000000-0000-0000-0000-000000000001",
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "System Administrator",
      "status": "active"
    }
  },
  "timestamp": 1234567890
}
```

#### 刷新令牌

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 健康检查

```http
GET /health
```

响应：`OK`

## 🐛 故障排查

### 数据库连接失败

```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 检查数据库连接
psql -h localhost -U postgres -d iot_platform

# 查看数据库日志
docker logs iot-postgres
```

### 服务启动失败

```bash
# 查看服务日志
docker logs iot-auth-service

# 检查环境变量
docker exec iot-auth-service env | grep DB_
```

### 测试失败

```bash
# 清理缓存
cd backend
go clean -cache -testcache

# 重新运行测试
go test -v ./...
```

## 📚 相关文档

- [Monorepo 设计](MONOREPO_DESIGN.md) - 详细的目录结构设计
- [数据库设计](DATABASE_DESIGN.md) - 数据库表结构
- [API 设计](API_DESIGN.md) - 完整的 API 清单
- [开发路线图](ROADMAP.md) - 迭代计划

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
