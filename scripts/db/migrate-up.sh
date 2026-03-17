#!/bin/bash

set -e

echo "=== Database Migration Script ==="

# 默认值
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-iot_platform}
DB_SSLMODE=${DB_SSLMODE:-disable}

# 检查 migrate 是否安装
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/migrate
    chmod +x /usr/local/bin/migrate
fi

# 构建数据库连接字符串
DB_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# 进入 migrations 目录
cd "$(dirname "$0")/../backend/migrations/auth"

echo ""
echo "Running migrations..."
echo "Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"

migrate -path . -database "${DB_URL}" up

echo ""
echo "=== Migration Complete ==="
