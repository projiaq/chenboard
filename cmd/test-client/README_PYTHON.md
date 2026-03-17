# Python 测试客户端

本目录包含两个 Python 测试客户端，用于测试 TSDB 数据库。

## 1. HTTP API 客户端 (test_client.py)

使用 HTTP RESTful API 连接 TSDB。

### 安装依赖

```bash
pip install requests
```

### 运行

```bash
python test_client.py
```

### 功能

- 健康检查
- 创建表
- 插入数据
- 查询数据

## 2. PostgreSQL 协议客户端 (test_pg_client.py)

使用 PostgreSQL 协议连接 TSDB（需要服务器启用 PostgreSQL 协议支持）。

### 安装依赖

```bash
pip install psycopg2-binary
```

### 运行

```bash
python test_pg_client.py
```

### 功能

- 通过 PostgreSQL 协议连接
- 执行标准 SQL 语句
- CREATE TABLE
- INSERT INTO
- SELECT 查询

## 使用示例

### HTTP API 方式

```python
import requests

# 创建表
requests.post("http://localhost:6041/api/create", json={
    "table": "sensors",
    "columns": ["ts", "temperature", "humidity"],
    "types": ["timestamp", "float", "float"]
})

# 插入数据
requests.post("http://localhost:6041/api/insert", json={
    "table": "sensors",
    "values": [{
        "ts": "2026-03-17T10:00:00Z",
        "temperature": 23.5,
        "humidity": 65.2
    }]
})

# 查询数据
resp = requests.post("http://localhost:6041/api/query", json={
    "table": "sensors",
    "start_time": "2026-03-17T00:00:00Z",
    "end_time": "2026-03-17T23:59:59Z",
    "columns": []
})
print(resp.json())
```

### PostgreSQL 协议方式

```python
import psycopg2

# 连接
conn = psycopg2.connect(
    host="localhost",
    port=5432,
    database="tsdb",
    user="admin"
)

cur = conn.cursor()

# 创建表
cur.execute("""
    CREATE TABLE sensors (
        ts TIMESTAMP,
        temperature FLOAT,
        humidity FLOAT
    )
""")

# 插入数据
cur.execute(
    "INSERT INTO sensors VALUES (%s, %s, %s)",
    ('2026-03-17 10:00:00', 23.5, 65.2)
)
conn.commit()

# 查询数据
cur.execute("SELECT * FROM sensors")
rows = cur.fetchall()
for row in rows:
    print(row)

cur.close()
conn.close()
```

## 故障排除

### 连接被拒绝

确保 TSDB 服务器正在运行：

```bash
# 启动服务器
go run ../../cmd/tsdb -data ./data -http-port 6041 -pg-port 5432
```

### 模块未找到

安装所需的 Python 包：

```bash
pip install requests psycopg2-binary
```

### 端口冲突

如果默认端口被占用，可以修改服务器启动参数：

```bash
go run ../../cmd/tsdb -http-port 6042 -pg-port 5433
```

然后修改测试脚本中的端口号。
