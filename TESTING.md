# 测试指南

## 方法 1: 使用 Go 测试客户端（推荐）

```bash
# 启动服务器
go run ./cmd/tsdb -port 6041 -data ./data

# 在另一个终端运行测试客户端
go run ./cmd/test-client
```

测试客户端会自动执行：
- 健康检查
- 创建表
- 插入 5 条测试数据
- 查询并显示结果

## 方法 2: 使用 Shell 脚本

```bash
# 启动服务器
./tsdb -port 6041 -data ./data

# 运行测试脚本（需要 jq 工具）
bash test.sh
```

## 方法 3: 使用 curl 手动测试

### 1. 健康检查
```bash
curl http://localhost:6041/health
```

### 2. 创建表
```bash
curl -X POST http://localhost:6041/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "columns": ["ts", "temperature", "humidity", "device_id"],
    "types": ["timestamp", "float", "float", "string"]
  }'
```

### 3. 插入数据
```bash
curl -X POST http://localhost:6041/api/insert \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "values": [{
      "ts": "2026-03-17T10:00:00Z",
      "temperature": 23.5,
      "humidity": 65.2,
      "device_id": "sensor001"
    }]
  }'
```

### 4. 查询数据
```bash
curl -X POST http://localhost:6041/api/query \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "start_time": "2026-03-17T00:00:00Z",
    "end_time": "2026-03-17T23:59:59Z",
    "columns": ["ts", "temperature", "humidity"]
  }'
```

## 方法 4: 使用 Postman 或 Insomnia

导入以下 API 端点：

- **健康检查**: `GET http://localhost:6041/health`
- **创建表**: `POST http://localhost:6041/api/create`
- **插入数据**: `POST http://localhost:6041/api/insert`
- **查询数据**: `POST http://localhost:6041/api/query`

## qStudio 连接（需要 PostgreSQL 协议适配器）

目前数据库使用 HTTP API，不支持直接用 qStudio 连接。如需 qStudio 支持，有两个选择：

### 选项 1: 使用 HTTP 客户端
在 qStudio 中使用 HTTP 查询功能（如果支持）

### 选项 2: 添加 PostgreSQL 协议支持
我可以添加 PostgreSQL wire protocol 适配器，让 qStudio 可以像连接 PostgreSQL 一样连接到 TSDB。

需要我实现 PostgreSQL 协议适配器吗？这样就可以用 qStudio、DBeaver、DataGrip 等标准 SQL 客户端连接了。

## 性能测试

```bash
# 运行基准测试
go test -bench=. -benchmem ./internal/storage

# 预期结果示例：
# BenchmarkInsert-8    500000    2500 ns/op    320 B/op    5 allocs/op
# BenchmarkQuery-8     100000   15000 ns/op   1024 B/op   15 allocs/op
```

## 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test -v ./internal/storage
go test -v ./internal/api
```
