# IoT Platform

一个基于 Go 的原生物联网管理平台，参考 ThingsBoard 的产品能力和架构思想。

## 特性

- 🏢 多租户 / 多项目 / 多角色权限
- 📱 设备管理：产品、设备、分组、标签、凭证
- 🔌 协议接入：MQTT、HTTP（可扩展 CoAP/LwM2M）
- 📊 数据能力：telemetry、attributes、events、commands、alarms
- ⚙️ 规则引擎：filter、transform、route、alarm、webhook
- 🎛️ 管理后台：用户、权限、设备、规则、告警、Dashboard
- 📺 大屏客户端：实时图表、告警面板、地图、设备状态墙

## 技术栈

### 后端
- Go 1.23+
- PostgreSQL 16+
- TimescaleDB 2.14+
- Redis 7+
- NATS 2.10+
- EMQX 5+

### 前端
- React 18+
- TypeScript 5+
- Ant Design 5+
- ECharts 5+

### 部署
- Docker & Docker Compose
- GitHub Actions
- Kubernetes (可选)

## 快速开始

### 前置要求

- Go 1.23+
- Docker & Docker Compose
- Node.js 20+ (前端开发)
- Make

### 本地开发

1. 克隆仓库
```bash
git clone <repository-url>
cd iot-platform
```

2. 启动基础设施
```bash
docker-compose up -d postgres redis
```

3. 初始化数据库
```bash
cd services/auth-service
make migrate-up
```

4. 启动服务
```bash
# 启动 auth-service
cd services/auth-service
make run

# 启动 api-gateway
cd services/api-gateway
make run
```

5. 访问服务
- API Gateway: http://localhost:8080
- Auth Service: http://localhost:8081

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 项目结构

详见 [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)

## API 文档

详见 [API_DESIGN.md](./API_DESIGN.md)

## 数据库设计

详见 [DATABASE_DESIGN.md](./DATABASE_DESIGN.md)

## 开发指南

### 构建所有服务
```bash
make build-all
```

### 运行测试
```bash
make test-all
```

### 代码检查
```bash
make lint-all
```

## 部署

### Docker 镜像构建
```bash
# 构建所有服务镜像
make docker-build-all

# 构建单个服务
cd services/auth-service
make docker-build
```

### Kubernetes 部署
```bash
kubectl apply -f deployments/k8s/
```

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可证

MIT License
