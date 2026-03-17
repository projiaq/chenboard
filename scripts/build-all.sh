#!/bin/bash

set -e

echo "Building all services..."

# 构建 pkg
echo "Building pkg..."
cd pkg
go mod download
go build ./...
cd ..

# 构建 auth-service
echo "Building auth-service..."
cd services/auth-service
go mod download
go build -o bin/auth-service ./cmd
cd ../..

echo "All services built successfully"
