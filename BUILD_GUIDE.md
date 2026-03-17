# 如何编译和运行

## 方式 1: 使用 Docker（推荐，无需安装 Go）

### Windows 用户

1. 确保已安装 Docker Desktop
2. 双击运行：`scripts\build\docker-build.bat`
3. 或者在命令行执行：

```cmd
cd scripts\build
docker-build.bat
```

### Linux/Mac 用户

```bash
chmod +x scripts/build/docker-build.sh
./scripts/build/docker-build.sh
```

### 运行服务

```bash
cd deployments/docker-compose
docker-compose up -d
```

### 测试 API

```bash
# 健康检查
curl http://localhost:8081/health

# 登录
curl -X POST http://localhost:8081/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"admin\",\"password\":\"admin123\",\"tenant_id\":\"default\"}"
```

## 方式 2: 本地编译（需要安装 Go）

### 安装 Go

1. 下载 Go 1.23.0+: https://go.dev/dl/
2. 安装并配置环境变量

### 编译步骤

```bash
# 1. 进入 backend 目录
cd backend

# 2. 下载依赖
go mod download

# 3. 运行测试
go test -v ./...

# 4. 编译
go build -o bin/auth-service.exe ./cmd/auth-service

# 5. 运行
bin\auth-service.exe
```

## 方式 3: 使用 GitHub Actions（自动编译）

1. 推送代码到 GitHub
2. GitHub Actions 会自动：
   - 运行测试
   - 编译所有服务
   - 构建 Docker 镜像

查看工作流：`.github/workflows/backend-ci.yml`

## 验证编译成功

### Docker 方式

```bash
# 查看镜像
docker images | grep iot-platform

# 应该看到：
# iot-platform/auth-service   latest   xxx   xxx   xxx MB
```

### 本地编译方式

```bash
# 查看可执行文件
ls -lh backend/bin/

# 应该看到：
# auth-service.exe (Windows)
# auth-service (Linux/Mac)
```

## 常见问题

### Q: Docker 构建失败

```bash
# 清理 Docker 缓存
docker system prune -a

# 重新构建
docker build -f deployments/docker/auth-service.Dockerfile -t iot-platform/auth-service:latest .
```

### Q: Go 编译失败

```bash
# 清理缓存
go clean -cache -modcache

# 重新下载依赖
go mod download

# 重新编译
go build -o bin/auth-service ./cmd/auth-service
```

### Q: 端口被占用

```bash
# Windows
netstat -ano | findstr :8081

# Linux/Mac
lsof -i :8081

# 修改端口
# 编辑 .env 文件，修改 SERVER_PORT
```

## 快速验证脚本

### Windows (PowerShell)

```powershell
# 保存为 verify.ps1
Write-Host "Verifying IoT Platform Build..." -ForegroundColor Green

# 检查 Docker
if (Get-Command docker -ErrorAction SilentlyContinue) {
    Write-Host "✓ Docker installed" -ForegroundColor Green
    docker --version
} else {
    Write-Host "✗ Docker not installed" -ForegroundColor Red
}

# 检查镜像
$image = docker images -q iot-platform/auth-service:latest
if ($image) {
    Write-Host "✓ Docker image built" -ForegroundColor Green
} else {
    Write-Host "✗ Docker image not found" -ForegroundColor Yellow
}

# 检查服务
$response = Invoke-WebRequest -Uri http://localhost:8081/health -UseBasicParsing -ErrorAction SilentlyContinue
if ($response.StatusCode -eq 200) {
    Write-Host "✓ Service is running" -ForegroundColor Green
} else {
    Write-Host "✗ Service not running" -ForegroundColor Yellow
}
```

### Linux/Mac (Bash)

```bash
#!/bin/bash
# 保存为 verify.sh

echo "Verifying IoT Platform Build..."

# 检查 Docker
if command -v docker &> /dev/null; then
    echo "✓ Docker installed"
    docker --version
else
    echo "✗ Docker not installed"
fi

# 检查镜像
if docker images | grep -q "iot-platform/auth-service"; then
    echo "✓ Docker image built"
else
    echo "✗ Docker image not found"
fi

# 检查服务
if curl -f http://localhost:8081/health &> /dev/null; then
    echo "✓ Service is running"
else
    echo "✗ Service not running"
fi
```
