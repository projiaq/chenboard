# IoT Platform Project Structure

```
iot-platform/
├── .github/
│   └── workflows/
│       ├── ci.yml                    # CI 工作流
│       └── docker.yml                # Docker 镜像构建
├── services/
│   ├── api-gateway/                  # API 网关
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── router/
│   │   ├── Dockerfile
│   │   ├── Makefile
│   │   └── go.mod
│   ├── auth-service/                 # 认证授权服务
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/              # 领域模型
│   │   │   ├── repository/          # 数据访问
│   │   │   ├── service/             # 业务逻辑
│   │   │   └── handler/             # HTTP 处理
│   │   ├── migrations/              # 数据库迁移
│   │   ├── Dockerfile
│   │   ├── Makefile
│   │   └── go.mod
│   ├── device-service/               # 设备管理服务
│   ├── data-service/                 # 数据处理服务
│   ├── rule-engine/                  # 规则引擎
│   ├── alarm-service/                # 告警服务
│   ├── notification-service/         # 通知服务
│   ├── protocol-mqtt/                # MQTT 协议适配器
│   └── protocol-http/                # HTTP 协议适配器
├── web/
│   ├── admin/                        # 管理后台
│   │   ├── src/
│   │   ├── public/
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── Dockerfile
│   └── dashboard/                    # 大屏客户端
│       ├── src/
│       ├── public/
│       ├── package.json
│       ├── tsconfig.json
│       └── Dockerfile
├── pkg/                              # 共享库
│   ├── config/                       # 配置管理
│   ├── database/                     # 数据库连接
│   ├── redis/                        # Redis 客户端
│   ├── nats/                         # NATS 客户端
│   ├── jwt/                          # JWT 工具
│   ├── logger/                       # 日志工具
│   ├── errors/                       # 错误处理
│   ├── middleware/                   # 通用中间件
│   └── models/                       # 共享模型
├── deployments/
│   ├── docker-compose.yml            # 本地开发环境
│   ├── docker-compose.prod.yml       # 生产环境
│   └── k8s/                          # Kubernetes 配置
├── scripts/
│   ├── init-db.sh                    # 数据库初始化
│   ├── build-all.sh                  # 构建所有服务
│   └── test-all.sh                   # 测试所有服务
├── docs/
│   ├── architecture.md               # 架构文档
│   ├── api/                          # API 文档
│   └── deployment.md                 # 部署文档
├── go.work                           # Go workspace
├── .gitignore
└── README.md
```

## 目录说明

### services/
每个服务独立的 Go 模块，包含：
- `cmd/`: 入口文件
- `internal/`: 内部实现（不对外暴露）
- `migrations/`: 数据库迁移脚本
- `Dockerfile`: 容器化配置
- `Makefile`: 构建脚本
- `go.mod`: 依赖管理

### pkg/
跨服务共享的 Go 包，所有服务都可以引用

### web/
前端项目，使用 React + TypeScript

### deployments/
部署配置，包括 Docker Compose 和 Kubernetes

### .github/workflows/
CI/CD 配置，自动化测试、构建、部署
