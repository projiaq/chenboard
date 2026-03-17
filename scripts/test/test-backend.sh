#!/bin/bash

set -e

echo "=== IoT Platform Backend Test Script ==="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# 进入 backend 目录
cd "$(dirname "$0")/../backend"

echo ""
echo "Step 1: Download dependencies..."
go mod download

echo ""
echo "Step 2: Run unit tests..."
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

echo ""
echo "Step 3: Generate coverage report..."
go tool cover -func=coverage.out

echo ""
echo "=== Test Complete ==="
