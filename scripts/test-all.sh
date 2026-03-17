#!/bin/bash

set -e

echo "Running tests for all services..."

# 测试 pkg
echo "Testing pkg..."
cd pkg
go test -v -race -coverprofile=coverage.out ./...
cd ..

# 测试 auth-service
echo "Testing auth-service..."
cd services/auth-service
go test -v -race -coverprofile=coverage.out ./...
cd ../..

echo "All tests passed successfully"
