package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:6041"

type Response struct {
	Status string                   `json:"status"`
	Data   []map[string]interface{} `json:"data,omitempty"`
}

func main() {
	fmt.Println("=== TSDB 测试客户端 ===")
	fmt.Println()

	// 1. 健康检查
	fmt.Println("1. 检查服务健康状态...")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("✓ 服务运行正常 (状态码: %d)\n\n", resp.StatusCode)

	// 2. 创建表
	fmt.Println("2. 创建传感器数据表...")
	createTable := map[string]interface{}{
		"table":   "sensors",
		"columns": []string{"ts", "temperature", "humidity", "device_id"},
		"types":   []string{"timestamp", "float", "float", "string"},
	}
	if err := postJSON("/api/create", createTable); err != nil {
		fmt.Printf("创建表失败: %v\n", err)
	} else {
		fmt.Println("✓ 表创建成功")
		fmt.Println()
	}

	// 3. 插入数据
	fmt.Println("3. 插入测试数据...")
	for i := 1; i <= 5; i++ {
		insertData := map[string]interface{}{
			"table": "sensors",
			"values": []map[string]interface{}{
				{
					"ts":          time.Now().Add(time.Duration(i) * time.Second).Format(time.RFC3339),
					"temperature": 20.0 + float64(i)*0.5,
					"humidity":    60.0 + float64(i)*2.0,
					"device_id":   fmt.Sprintf("sensor%03d", i),
				},
			},
		}
		if err := postJSON("/api/insert", insertData); err != nil {
			fmt.Printf("插入数据失败: %v\n", err)
		} else {
			fmt.Printf("✓ 插入数据 %d/5\n", i)
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()

	// 4. 查询数据
	fmt.Println("4. 查询数据...")
	queryData := map[string]interface{}{
		"table":      "sensors",
		"start_time": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"end_time":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"columns":    []string{"ts", "temperature", "humidity", "device_id"},
	}

	respData, err := postJSONWithResponse("/api/query", queryData)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	var result Response
	if err := json.Unmarshal(respData, &result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 查询成功，返回 %d 条记录:\n", len(result.Data))
	for i, row := range result.Data {
		fmt.Printf("  [%d] 设备: %v, 温度: %.1f°C, 湿度: %.1f%%\n",
			i+1, row["device_id"], row["temperature"], row["humidity"])
	}

	fmt.Println("\n=== 测试完成 ===")
}

func postJSON(path string, data interface{}) error {
	_, err := postJSONWithResponse(path, data)
	return err
}

func postJSONWithResponse(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
