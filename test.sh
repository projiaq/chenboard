#!/bin/bash

# TSDB 测试脚本

BASE_URL="http://localhost:6041"

echo "=== 1. 检查服务健康状态 ==="
curl -s ${BASE_URL}/health | jq .

echo -e "\n=== 2. 创建传感器数据表 ==="
curl -s -X POST ${BASE_URL}/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "table": "sensors",
    "columns": ["ts", "temperature", "humidity", "device_id"],
    "types": ["timestamp", "float", "float", "string"]
  }' | jq .

echo -e "\n=== 3. 插入测试数据 ==="
for i in {1..5}; do
  temp=$(echo "20 + $i * 0.5" | bc)
  humidity=$(echo "60 + $i * 2" | bc)
  timestamp=$(date -u -d "+${i} seconds" +"%Y-%m-%dT%H:%M:%SZ")

  curl -s -X POST ${BASE_URL}/api/insert \
    -H "Content-Type: application/json" \
    -d "{
      \"table\": \"sensors\",
      \"values\": [{
        \"ts\": \"${timestamp}\",
        \"temperature\": ${temp},
        \"humidity\": ${humidity},
        \"device_id\": \"sensor00${i}\"
      }]
    }" | jq .

  echo "插入数据 ${i}/5"
done

echo -e "\n=== 4. 查询数据 ==="
start_time=$(date -u -d "-1 hour" +"%Y-%m-%dT%H:%M:%SZ")
end_time=$(date -u -d "+1 hour" +"%Y-%m-%dT%H:%M:%SZ")

curl -s -X POST ${BASE_URL}/api/query \
  -H "Content-Type: application/json" \
  -d "{
    \"table\": \"sensors\",
    \"start_time\": \"${start_time}\",
    \"end_time\": \"${end_time}\",
    \"columns\": [\"ts\", \"temperature\", \"humidity\", \"device_id\"]
  }" | jq .

echo -e "\n=== 测试完成 ==="
