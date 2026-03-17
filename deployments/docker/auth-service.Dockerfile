# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git make

# 复制 go.mod 和 go.sum
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend/ ./

# 构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/auth-service ./cmd/auth-service

# Runtime stage
FROM alpine:3.19

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/auth-service .

# 复制迁移文件
COPY backend/migrations/auth ./migrations

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# 暴露端口
EXPOSE 8081

# 运行
CMD ["./auth-service"]
