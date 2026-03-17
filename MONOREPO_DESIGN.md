# IoT Platform Monorepo 结构设计

## 1. 完整目录树

```
iot-platform/
├── .github/
│   ├── workflows/
│   │   ├── backend-ci.yml           # 后端 CI（lint、test、build）
│   │   ├── frontend-ci.yml          # 前端 CI（lint、test、build）
│   │   ├── docker-build.yml         # Docker 镜像构建
│   │   └── deploy.yml               # 部署流程
│   └── dependabot.yml               # 依赖更新配置
│
├── backend/                         # 后端 Go 代码根目录
│   ├── cmd/                         # 可执行程序入口
│   │   ├── auth-service/
│   │   │   └── main.go
│   │   ├── device-service/
│   │   │   └── main.go
│   │   ├── data-service/
│   │   │   └── main.go
│   │   ├── rule-engine/
│   │   │   └── main.go
│   │   ├── alarm-service/
│   │   │   └── main.go
│   │   ├── notification-service/
│   │   │   └── main.go
│   │   ├── api-gateway/
│   │   │   └── main.go
│   │   ├── protocol-mqtt/
│   │   │   └── main.go
│   │   └── protocol-http/
│   │       └── main.go
│   │
│   ├── internal/                    # 私有应用代码（不可被外部引用）
│   │   ├── auth/                    # auth-service 内部实现
│   │   │   ├── domain/              # 领域模型
│   │   │   ├── repository/          # 数据访问层
│   │   │   ├── service/             # 业务逻辑层
│   │   │   ├── handler/             # HTTP 处理层
│   │   │   └── middleware/          # 中间件
│   │   ├── device/                  # device-service 内部实现
│   │   ├── data/                    # data-service 内部实现
│   │   ├── rule/                    # rule-engine 内部实现
│   │   ├── alarm/                   # alarm-service 内部实现
│   │   ├── notification/            # notification-service 内部实现
│   │   ├── gateway/                 # api-gateway 内部实现
│   │   └── protocol/                # 协议适配器共享代码
│   │       ├── mqtt/
│   │       └── http/
│   │
│   ├── pkg/                         # 公共库（可被外部引用）
│   │   ├── config/                  # 配置管理
│   │   ├── database/                # 数据库连接
│   │   ├── redis/                   # Redis 客户端
│   │   ├── nats/                    # NATS 客户端
│   │   ├── jwt/                     # JWT 工具
│   │   ├── logger/                  # 日志工具
│   │   ├── errors/                  # 错误处理
│   │   ├── middleware/              # 通用中间件
│   │   ├── response/                # 统一响应
│   │   └── models/                  # 共享数据模型
│   │
│   ├── migrations/                  # 数据库迁移脚本
│   │   ├── auth/
│   │   │   ├── 000001_init.up.sql
│   │   │   └── 000001_init.down.sql
│   │   ├── device/
│   │   ├── data/
│   │   └── alarm/
│   │
│   ├── test/                        # 集成测试和 E2E 测试
│   │   ├── integration/
│   │   └── e2e/
│   │
│   ├── go.mod                       # Go 模块定义
│   ├── go.sum                       # Go 依赖锁定
│   ├── Makefile                     # 构建脚本
│   └── .golangci.yml                # Go lint 配置
│
├── frontend/                        # 前端代码根目录
│   ├── admin/                       # 管理后台
│   │   ├── src/
│   │   │   ├── components/
│   │   │   ├── pages/
│   │   │   ├── services/
│   │   │   ├── utils/
│   │   │   ├── App.tsx
│   │   │   └── main.tsx
│   │   ├── public/
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── vite.config.ts
│   │   └── .eslintrc.js
│   │
│   ├── dashboard/                   # 大屏客户端
│   │   ├── src/
│   │   ├── public/
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── vite.config.ts
│   │   └── .eslintrc.js
│   │
│   └── shared/                      # 前端共享代码
│       ├── components/              # 共享组件
│       ├── hooks/                   # 共享 Hooks
│       ├── utils/                   # 工具函数
│       ├── types/                   # TypeScript 类型定义
│       └── package.json
│
├── deployments/                     # 部署配置
│   ├── docker/                      # Dockerfile 集合
│   │   ├── auth-service.Dockerfile
│   │   ├── device-service.Dockerfile
│   │   ├── api-gateway.Dockerfile
│   │   ├── admin-web.Dockerfile
│   │   └── dashboard-web.Dockerfile
│   │
│   ├── docker-compose/              # Docker Compose 配置
│   │   ├── docker-compose.yml       # 开发环境
│   │   ├── docker-compose.prod.yml  # 生产环境
│   │   └── docker-compose.test.yml  # 测试环境
│   │
│   ├── kubernetes/                  # Kubernetes 配置
│   │   ├── base/                    # 基础配置
│   │   ├── overlays/                # 环境覆盖
│   │   │   ├── dev/
│   │   │   ├── staging/
│   │   │   └── prod/
│   │   └── kustomization.yaml
│   │
│   └── helm/                        # Helm Charts
│       └── iot-platform/
│           ├── Chart.yaml
│           ├── values.yaml
│           └── templates/
│
├── scripts/                         # 脚本工具
│   ├── build/
│   │   ├── build-backend.sh         # 构建所有后端服务
│   │   ├── build-frontend.sh        # 构建所有前端应用
│   │   └── build-all.sh             # 构建全部
│   │
│   ├── test/
│   │   ├── test-backend.sh          # 测试后端
│   │   ├── test-frontend.sh         # 测试前端
│   │   └── test-all.sh              # 测试全部
│   │
│   ├── deploy/
│   │   ├── deploy-dev.sh            # 部署到开发环境
│   │   ├── deploy-staging.sh        # 部署到预发布环境
│   │   └── deploy-prod.sh           # 部署到生产环境
│   │
│   └── db/
│       ├── init-db.sh               # 初始化数据库
│       ├── migrate-up.sh            # 执行迁移
│       └── migrate-down.sh          # 回滚迁移
│
├── docs/                            # 文档
│   ├── architecture/                # 架构文档
│   │   ├── overview.md
│   │   ├── backend.md
│   │   ├── frontend.md
│   │   └── deployment.md
│   │
│   ├── api/                         # API 文档
│   │   ├── auth-service.md
│   │   ├── device-service.md
│   │   └── openapi.yaml
│   │
│   ├── development/                 # 开发文档
│   │   ├── setup.md
│   │   ├── coding-standards.md
│   │   └── contributing.md
│   │
│   └── operations/                  # 运维文档
│       ├── deployment.md
│       ├── monitoring.md
│       └── troubleshooting.md
│
├── .gitignore                       # Git 忽略配置
├── .editorconfig                    # 编辑器配置
├── README.md                        # 项目说明
├── LICENSE                          # 许可证
└── CHANGELOG.md                     # 变更日志
```

## 2. 目录职责说明

### 2.1 根目录
- **.github/**: GitHub 配置，包括 CI/CD 工作流
- **backend/**: 所有后端 Go 代码
- **frontend/**: 所有前端代码
- **deployments/**: 所有部署配置
- **scripts/**: 自动化脚本
- **docs/**: 项目文档

### 2.2 Backend 目录
- **cmd/**: 每个服务的入口文件（main.go）
  - 每个子目录对应一个可执行程序
  - 只包含启动逻辑，不包含业务代码

- **internal/**: 私有应用代码
  - 按服务划分子目录
  - 每个服务的内部实现（domain、repository、service、handler）
  - Go 编译器强制：internal 下的代码不能被外部引用

- **pkg/**: 公共库
  - 可以被 internal 引用
  - 可以被外部项目引用
  - 包含通用工具和基础设施代码

- **migrations/**: 数据库迁移脚本
  - 按服务划分子目录
  - 使用 golang-migrate 格式

- **test/**: 集成测试和 E2E 测试
  - 单元测试放在各自的包内（*_test.go）
  - 集成测试放在这里

### 2.3 Frontend 目录
- **admin/**: 管理后台应用
  - 独立的 React 项目
  - 有自己的 package.json

- **dashboard/**: 大屏客户端应用
  - 独立的 React 项目
  - 有自己的 package.json

- **shared/**: 前端共享代码
  - 可以被 admin 和 dashboard 引用
  - 作为 npm workspace 的一部分

### 2.4 Deployments 目录
- **docker/**: 所有 Dockerfile
  - 每个服务一个 Dockerfile
  - 统一命名规范：<service-name>.Dockerfile

- **docker-compose/**: Docker Compose 配置
  - 按环境分离

- **kubernetes/**: K8s 配置
  - 使用 Kustomize 管理多环境

- **helm/**: Helm Charts
  - 用于生产环境部署

### 2.5 Scripts 目录
- **build/**: 构建脚本
- **test/**: 测试脚本
- **deploy/**: 部署脚本
- **db/**: 数据库管理脚本

## 3. CI/CD 参与规则

### 3.1 参与 CI 的目录

#### Backend CI (`backend-ci.yml`)
触发路径：
```yaml
paths:
  - 'backend/**'
  - '.github/workflows/backend-ci.yml'
```

执行内容：
- Lint: `backend/` 下所有 Go 代码
- Test: `backend/` 下所有单元测试和集成测试
- Build: `backend/cmd/` 下所有服务

#### Frontend CI (`frontend-ci.yml`)
触发路径：
```yaml
paths:
  - 'frontend/**'
  - '.github/workflows/frontend-ci.yml'
```

执行内容：
- Lint: `frontend/admin/` 和 `frontend/dashboard/`
- Test: 前端单元测试
- Build: 前端生产构建

#### Docker Build (`docker-build.yml`)
触发路径：
```yaml
paths:
  - 'backend/**'
  - 'frontend/**'
  - 'deployments/docker/**'
  - '.github/workflows/docker-build.yml'
```

执行内容：
- 构建所有服务的 Docker 镜像
- 推送到容器镜像仓库

### 3.2 不参与 CI 的目录
- `docs/`: 文档变更不触发 CI
- `scripts/`: 脚本变更不触发 CI（除非明确需要）
- `deployments/kubernetes/`: K8s 配置变更不触发 CI
- `deployments/helm/`: Helm 配置变更不触发 CI

## 4. Docker 构建规则

### 4.1 后端服务 Dockerfile 结构

**位置**: `deployments/docker/<service-name>.Dockerfile`

**构建上下文**: 仓库根目录

**示例** (`deployments/docker/auth-service.Dockerfile`):
```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /app/bin/auth-service ./cmd/auth-service

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/auth-service .
COPY backend/migrations/auth ./migrations
EXPOSE 8081
CMD ["./auth-service"]
```

**构建命令**:
```bash
docker build -f deployments/docker/auth-service.Dockerfile -t iot-platform/auth-service:latest .
```

### 4.2 前端应用 Dockerfile 结构

**位置**: `deployments/docker/admin-web.Dockerfile`

**构建上下文**: 仓库根目录

**示例**:
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY frontend/admin/package*.json ./
RUN npm ci
COPY frontend/admin/ ./
COPY frontend/shared/ ../shared/
RUN npm run build

# Runtime stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY deployments/docker/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 4.3 参与 Docker 构建的目录
- `backend/`: 所有后端服务
- `frontend/admin/`: 管理后台
- `frontend/dashboard/`: 大屏客户端
- `frontend/shared/`: 前端共享代码
- `deployments/docker/`: Dockerfile 和配置文件

## 5. 共享代码组织

### 5.1 后端共享代码

**位置**: `backend/pkg/`

**引用方式**:
```go
import "iot-platform/backend/pkg/logger"
import "iot-platform/backend/pkg/jwt"
```

**go.mod**:
```go
module iot-platform/backend

go 1.23
```

**特点**:
- 单一 Go 模块
- 所有服务共享同一个 go.mod
- 简化依赖管理
- 加快构建速度

### 5.2 前端共享代码

**位置**: `frontend/shared/`

**引用方式** (使用 npm workspace):

**根目录 package.json**:
```json
{
  "name": "iot-platform-frontend",
  "private": true,
  "workspaces": [
    "frontend/admin",
    "frontend/dashboard",
    "frontend/shared"
  ]
}
```

**admin/package.json**:
```json
{
  "dependencies": {
    "@iot-platform/shared": "workspace:*"
  }
}
```

**特点**:
- 使用 npm workspaces
- 共享代码作为内部包
- 统一依赖管理
- 支持 TypeScript 类型共享

## 6. GitHub Actions 路径触发策略

### 6.1 Backend CI 触发规则

```yaml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'
      - '.github/workflows/backend-ci.yml'
  pull_request:
    branches: [main, develop]
    paths:
      - 'backend/**'
      - '.github/workflows/backend-ci.yml'

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Lint
        run: |
          cd backend
          golangci-lint run ./...

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - name: Test
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...

  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service:
          - auth-service
          - device-service
          - api-gateway
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - name: Build
        run: |
          cd backend
          go build -o bin/${{ matrix.service }} ./cmd/${{ matrix.service }}
```

### 6.2 Frontend CI 触发规则

```yaml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/**'
      - '.github/workflows/frontend-ci.yml'
  pull_request:
    branches: [main, develop]
    paths:
      - 'frontend/**'
      - '.github/workflows/frontend-ci.yml'

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        app: [admin, dashboard]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        run: npm ci
      - name: Lint
        run: |
          cd frontend/${{ matrix.app }}
          npm run lint
      - name: Test
        run: |
          cd frontend/${{ matrix.app }}
          npm run test
      - name: Build
        run: |
          cd frontend/${{ matrix.app }}
          npm run build
```

### 6.3 Docker Build 触发规则

```yaml
name: Docker Build

on:
  push:
    branches: [main]
    tags: ['v*']
    paths:
      - 'backend/**'
      - 'frontend/**'
      - 'deployments/docker/**'

jobs:
  build-backend:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service:
          - auth-service
          - device-service
          - api-gateway
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: deployments/docker/${{ matrix.service }}.Dockerfile
          push: true
          tags: ghcr.io/${{ github.repository }}/${{ matrix.service }}:latest
```

## 7. 推荐的 Monorepo 组织方案

### 7.1 核心原则

1. **按技术栈分离**: backend、frontend、deployments 清晰分离
2. **单一 Go 模块**: 所有后端服务共享一个 go.mod，简化依赖管理
3. **npm workspaces**: 前端使用 workspaces 管理共享代码
4. **统一 Dockerfile 位置**: 所有 Dockerfile 放在 `deployments/docker/`
5. **根目录构建上下文**: Docker 构建使用根目录作为上下文
6. **路径触发 CI**: 根据变更路径智能触发 CI

### 7.2 优势

1. **清晰的边界**: 后端、前端、部署配置各自独立
2. **高效的 CI**: 只测试和构建变更的部分
3. **简化的依赖管理**: 后端单一模块，前端 workspaces
4. **灵活的部署**: 支持 Docker Compose、K8s、Helm
5. **易于扩展**: 新增服务只需在 cmd/ 和 internal/ 添加目录
6. **符合最佳实践**: 遵循 Go 标准项目布局和前端工程规范

### 7.3 与传统 Monorepo 的区别

**不推荐的方案**（每个服务独立模块）:
```
services/
  auth-service/
    go.mod          ❌ 每个服务独立 go.mod
    go.sum
    cmd/
    internal/
  device-service/
    go.mod          ❌ 依赖管理复杂
    go.sum
    cmd/
    internal/
```

**推荐的方案**（单一模块）:
```
backend/
  go.mod            ✅ 单一 go.mod
  go.sum
  cmd/
    auth-service/
    device-service/
  internal/
    auth/
    device/
  pkg/              ✅ 共享代码
```

### 7.4 迁移路径

如果未来需要拆分为独立仓库：

1. **拆分后端服务**:
   ```bash
   # 将 backend/internal/auth/ 移到新仓库
   # 将 backend/cmd/auth-service/ 移到新仓库
   # 将 backend/pkg/ 发布为独立库
   ```

2. **拆分前端应用**:
   ```bash
   # 将 frontend/admin/ 移到新仓库
   # 将 frontend/shared/ 发布为 npm 包
   ```

3. **保持部署配置**:
   ```bash
   # deployments/ 可以保留在主仓库
   # 或者移到独立的 infra 仓库
   ```

## 8. 实施检查清单

- [ ] 创建目录结构
- [ ] 配置 Go 模块（backend/go.mod）
- [ ] 配置 npm workspaces（根目录 package.json）
- [ ] 创建 Dockerfile 模板
- [ ] 配置 GitHub Actions 工作流
- [ ] 编写构建脚本
- [ ] 编写测试脚本
- [ ] 配置 golangci-lint
- [ ] 配置 ESLint
- [ ] 编写 README 和文档

## 9. 总结

这个 monorepo 结构：
- ✅ 清晰分层（后端、前端、部署、文档、脚本）
- ✅ 支持 CI 独立 lint/test/build
- ✅ 支持未来拆分为多个服务
- ✅ 符合 Go 和前端工程实践
- ✅ 适合 GitHub Actions
- ✅ 易于维护和扩展

**这是唯一推荐的方案，不需要其他选项。**
