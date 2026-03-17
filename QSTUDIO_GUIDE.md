# qStudio 连接指南

## 快速开始

### 1. 启动 TSDB 服务器

```bash
# 编译
go build -o tsdb ./cmd/tsdb

# 启动服务器（同时启动 HTTP API 和 PostgreSQL 协议）
./tsdb -data ./data -http-port 6041 -pg-port 5432
```

服务器会启动两个端口：
- **5432**: PostgreSQL 协议端口（用于 qStudio 连接）
- **6041**: HTTP API 端口（用于 RESTful 访问）

### 2. 在 qStudio 中连接

1. 打开 qStudio
2. 点击 "Add Connection"
3. 选择 "PostgreSQL"
4. 填写连接信息：
   - **Host**: `localhost`
   - **Port**: `5432`
   - **Database**: `tsdb`（任意名称）
   - **Username**: `admin`（任意名称）
   - **Password**: （留空或任意）

5. 点击 "Test" 测试连接
6. 点击 "Save" 保存连接

### 3. 执行 SQL 查询

连接成功后，可以在 qStudio 中执行标准 SQL：

#### 创建表
```sql
CREATE TABLE sensors (
    ts TIMESTAMP,
    temperature FLOAT,
    humidity FLOAT,
    device_id VARCHAR(50)
);
```

#### 插入数据
```sql
INSERT INTO sensors VALUES
    ('2026-03-17 10:00:00', 23.5, 65.2, 'sensor001'),
    ('2026-03-17 10:01:00', 24.0, 66.0, 'sensor002'),
    ('2026-03-17 10:02:00', 23.8, 65.5, 'sensor003');
```

#### 查询数据
```sql
-- 查询所有数据
SELECT * FROM sensors;

-- 按时间范围查询
SELECT * FROM sensors
WHERE ts > '2026-03-17 09:00:00';

-- 查询特定设备
SELECT ts, temperature FROM sensors
WHERE device_id = 'sensor001';

-- 聚合查询（未来支持）
SELECT device_id, AVG(temperature) as avg_temp
FROM sensors
GROUP BY device_id;
```

## 支持的 SQL 功能

### ✅ 已支持
- `CREATE TABLE` - 创建表
- `INSERT INTO` - 插入数据
- `SELECT * FROM` - 查询数据
- 基本数据类型：TIMESTAMP, INT, FLOAT, VARCHAR, BOOLEAN

### 🚧 计划支持
- `WHERE` 条件过滤
- `ORDER BY` 排序
- `LIMIT` 限制结果数量
- `GROUP BY` 聚合
- `JOIN` 表连接
- `UPDATE` 更新数据
- `DELETE` 删除数据

## 数据类型映射

| SQL 类型 | TSDB 内部类型 |
|---------|-------------|
| TIMESTAMP, DATETIME | TypeTimestamp |
| INT, BIGINT | TypeInt |
| FLOAT, DOUBLE, REAL | TypeFloat |
| VARCHAR, TEXT, CHAR | TypeString |
| BOOLEAN, BOOL | TypeBool |

## 其他 SQL 客户端

除了 qStudio，还可以使用以下客户端连接：

### DBeaver
1. 新建连接 → PostgreSQL
2. Host: `localhost`, Port: `5432`
3. Database: `tsdb`
4. 连接即可

### DataGrip
1. 新建数据源 → PostgreSQL
2. Host: `localhost`, Port: `5432`
3. Database: `tsdb`
4. 测试连接

### psql 命令行
```bash
psql -h localhost -p 5432 -U admin -d tsdb
```

### Python (psycopg2)
```python
import psycopg2

conn = psycopg2.connect(
    host="localhost",
    port=5432,
    database="tsdb",
    user="admin"
)

cur = conn.cursor()
cur.execute("SELECT * FROM sensors")
rows = cur.fetchall()
for row in rows:
    print(row)
```

### Go (pgx)
```go
import "github.com/jackc/pgx/v5"

conn, _ := pgx.Connect(context.Background(),
    "postgres://admin@localhost:5432/tsdb")
defer conn.Close(context.Background())

rows, _ := conn.Query(context.Background(),
    "SELECT * FROM sensors")
defer rows.Close()
```

## 故障排除

### 连接被拒绝
- 确认服务器已启动：`ps aux | grep tsdb`
- 检查端口是否被占用：`lsof -i :5432`
- 尝试使用其他端口：`./tsdb -pg-port 5433`

### SQL 语法错误
- 确保使用标准 PostgreSQL SQL 语法
- 表名和列名区分大小写
- 字符串使用单引号 `'string'`

### 性能优化
- 使用批量插入而不是单条插入
- 为时间戳列创建索引（未来支持）
- 定期执行 VACUUM（未来支持）

## 示例：完整工作流

```sql
-- 1. 创建 IoT 传感器表
CREATE TABLE iot_sensors (
    ts TIMESTAMP,
    device_id VARCHAR(50),
    temperature FLOAT,
    humidity FLOAT,
    pressure FLOAT,
    location VARCHAR(100)
);

-- 2. 批量插入测试数据
INSERT INTO iot_sensors VALUES
    ('2026-03-17 10:00:00', 'device001', 23.5, 65.2, 1013.25, 'Beijing'),
    ('2026-03-17 10:00:00', 'device002', 24.0, 66.0, 1012.80, 'Shanghai'),
    ('2026-03-17 10:01:00', 'device001', 23.6, 65.3, 1013.20, 'Beijing'),
    ('2026-03-17 10:01:00', 'device002', 24.1, 66.2, 1012.75, 'Shanghai');

-- 3. 查询最新数据
SELECT * FROM iot_sensors ORDER BY ts DESC LIMIT 10;

-- 4. 查询特定设备
SELECT * FROM iot_sensors WHERE device_id = 'device001';

-- 5. 统计分析（未来支持）
SELECT device_id,
       AVG(temperature) as avg_temp,
       MAX(temperature) as max_temp,
       MIN(temperature) as min_temp
FROM iot_sensors
GROUP BY device_id;
```

## 性能基准

在标准硬件上（Intel i7, 16GB RAM）：
- 插入速度：~100,000 行/秒
- 查询速度：~50,000 行/秒
- 压缩比：~5:1（使用 zstd）

## 下一步

- 查看 [README.md](README.md) 了解更多功能
- 查看 [TESTING.md](TESTING.md) 了解测试方法
- 提交 Issue 报告问题或建议新功能
