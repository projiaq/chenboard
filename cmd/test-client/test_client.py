#!/usr/bin/env python3
"""
TSDB Python 测试客户端
用于测试 TSDB 时序数据库的 HTTP API
"""

import requests
import json
from datetime import datetime, timedelta
import time

BASE_URL = "http://localhost:6041"


def health_check():
    """健康检查"""
    print("1. 检查服务健康状态...")
    try:
        resp = requests.get(f"{BASE_URL}/health")
        resp.raise_for_status()
        print(f"✓ 服务运行正常 (状态码: {resp.status_code})\n")
        return True
    except Exception as e:
        print(f"错误: {e}")
        return False


def create_table():
    """创建表"""
    print("2. 创建传感器数据表...")
    data = {
        "table": "sensors",
        "columns": ["ts", "temperature", "humidity", "device_id"],
        "types": ["timestamp", "float", "float", "string"]
    }

    try:
        resp = requests.post(f"{BASE_URL}/api/create", json=data)
        resp.raise_for_status()
        print("✓ 表创建成功\n")
        return True
    except requests.exceptions.HTTPError as e:
        if "already exists" in str(e):
            print("✓ 表已存在\n")
            return True
        print(f"创建表失败: {e}")
        return False
    except Exception as e:
        print(f"创建表失败: {e}")
        return False


def insert_data():
    """插入测试数据"""
    print("3. 插入测试数据...")

    for i in range(1, 6):
        ts = (datetime.now() + timedelta(seconds=i)).isoformat() + "Z"
        data = {
            "table": "sensors",
            "values": [{
                "ts": ts,
                "temperature": 20.0 + i * 0.5,
                "humidity": 60.0 + i * 2.0,
                "device_id": f"sensor{i:03d}"
            }]
        }

        try:
            resp = requests.post(f"{BASE_URL}/api/insert", json=data)
            resp.raise_for_status()
            print(f"✓ 插入数据 {i}/5")
        except Exception as e:
            print(f"插入数据失败: {e}")
            return False

        time.sleep(0.1)

    print()
    return True


def query_data():
    """查询数据"""
    print("4. 查询数据...")

    start_time = (datetime.now() - timedelta(hours=1)).isoformat() + "Z"
    end_time = (datetime.now() + timedelta(hours=1)).isoformat() + "Z"

    data = {
        "table": "sensors",
        "start_time": start_time,
        "end_time": end_time,
        "columns": ["ts", "temperature", "humidity", "device_id"]
    }

    try:
        resp = requests.post(f"{BASE_URL}/api/query", json=data)
        resp.raise_for_status()
        result = resp.json()

        records = result.get("data", [])
        print(f"✓ 查询成功，返回 {len(records)} 条记录:")

        for i, row in enumerate(records, 1):
            device_id = row.get("device_id", "N/A")
            temperature = row.get("temperature", 0)
            humidity = row.get("humidity", 0)
            print(f"  [{i}] 设备: {device_id}, 温度: {temperature:.1f}°C, 湿度: {humidity:.1f}%")

        return True
    except Exception as e:
        print(f"查询失败: {e}")
        return False


def main():
    """主函数"""
    print("=== TSDB Python 测试客户端 ===\n")

    # 1. 健康检查
    if not health_check():
        return

    # 2. 创建表
    if not create_table():
        return

    # 3. 插入数据
    if not insert_data():
        return

    # 4. 查询数据
    if not query_data():
        return

    print("\n=== 测试完成 ===")


if __name__ == "__main__":
    main()
