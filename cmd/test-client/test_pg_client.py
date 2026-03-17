#!/usr/bin/env python3
"""
TSDB PostgreSQL 协议测试客户端
使用 psycopg2 通过 PostgreSQL 协议连接 TSDB
"""

import psycopg2
from datetime import datetime
import time

# 连接配置
CONN_PARAMS = {
    "host": "localhost",
    "port": 5432,
    "database": "tsdb",
    "user": "admin"
}


def test_postgresql_protocol():
    """测试 PostgreSQL 协议连接"""
    print("=== TSDB PostgreSQL 协议测试 ===\n")

    try:
        # 1. 连接数据库
        print("1. 连接到 TSDB...")
        conn = psycopg2.connect(**CONN_PARAMS)
        cur = conn.cursor()
        print("✓ 连接成功\n")

        # 2. 创建表
        print("2. 创建传感器数据表...")
        try:
            cur.execute("""
                CREATE TABLE sensors (
                    ts TIMESTAMP,
                    temperature FLOAT,
                    humidity FLOAT,
                    device_id VARCHAR(50)
                )
            """)
            conn.commit()
            print("✓ 表创建成功\n")
        except Exception as e:
            if "already exists" in str(e):
                print("✓ 表已存在\n")
                conn.rollback()
            else:
                raise

        # 3. 插入数据
        print("3. 插入测试数据...")
        for i in range(1, 6):
            ts = datetime.now()
            temperature = 20.0 + i * 0.5
            humidity = 60.0 + i * 2.0
            device_id = f"sensor{i:03d}"

            cur.execute(
                "INSERT INTO sensors VALUES (%s, %s, %s, %s)",
                (ts, temperature, humidity, device_id)
            )
            print(f"✓ 插入数据 {i}/5")
            time.sleep(0.1)

        conn.commit()
        print()

        # 4. 查询数据
        print("4. 查询数据...")
        cur.execute("SELECT * FROM sensors")
        rows = cur.fetchall()

        print(f"✓ 查询成功，返回 {len(rows)} 条记录:")
        for i, row in enumerate(rows, 1):
            ts, temperature, humidity, device_id = row
            print(f"  [{i}] 时间: {ts}, 设备: {device_id}, 温度: {temperature:.1f}°C, 湿度: {humidity:.1f}%")

        # 5. 关闭连接
        cur.close()
        conn.close()

        print("\n=== 测试完成 ===")

    except psycopg2.Error as e:
        print(f"数据库错误: {e}")
    except Exception as e:
        print(f"错误: {e}")


if __name__ == "__main__":
    test_postgresql_protocol()
