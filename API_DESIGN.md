# API Design

## 4.1 认证授权 API (auth-service)

### 认证相关
```
POST   /api/v1/auth/login              # 用户登录
POST   /api/v1/auth/logout             # 用户登出
POST   /api/v1/auth/refresh            # 刷新 Token
GET    /api/v1/auth/me                 # 获取当前用户信息
PUT    /api/v1/auth/password           # 修改密码
```

### 租户管理
```
GET    /api/v1/tenants                 # 获取租户列表
POST   /api/v1/tenants                 # 创建租户
GET    /api/v1/tenants/:id             # 获取租户详情
PUT    /api/v1/tenants/:id             # 更新租户
DELETE /api/v1/tenants/:id             # 删除租户
```

### 用户管理
```
GET    /api/v1/users                   # 获取用户列表
POST   /api/v1/users                   # 创建用户
GET    /api/v1/users/:id               # 获取用户详情
PUT    /api/v1/users/:id               # 更新用户
DELETE /api/v1/users/:id               # 删除用户
PUT    /api/v1/users/:id/status        # 修改用户状态
GET    /api/v1/users/:id/roles         # 获取用户角色
PUT    /api/v1/users/:id/roles         # 设置用户角色
```

### 角色管理
```
GET    /api/v1/roles                   # 获取角色列表
POST   /api/v1/roles                   # 创建角色
GET    /api/v1/roles/:id               # 获取角色详情
PUT    /api/v1/roles/:id               # 更新角色
DELETE /api/v1/roles/:id               # 删除角色
GET    /api/v1/roles/:id/permissions   # 获取角色权限
PUT    /api/v1/roles/:id/permissions   # 设置角色权限
```

### 权限管理
```
GET    /api/v1/permissions             # 获取权限列表
POST   /api/v1/permissions             # 创建权限
GET    /api/v1/permissions/:id         # 获取权限详情
PUT    /api/v1/permissions/:id         # 更新权限
DELETE /api/v1/permissions/:id         # 删除权限
```

## 4.2 设备管理 API (device-service)

### 产品管理
```
GET    /api/v1/products                # 获取产品列表
POST   /api/v1/products                # 创建产品
GET    /api/v1/products/:id            # 获取产品详情
PUT    /api/v1/products/:id            # 更新产品
DELETE /api/v1/products/:id            # 删除产品
```

### 设备管理
```
GET    /api/v1/devices                 # 获取设备列表
POST   /api/v1/devices                 # 创建设备
GET    /api/v1/devices/:id             # 获取设备详情
PUT    /api/v1/devices/:id             # 更新设备
DELETE /api/v1/devices/:id             # 删除设备
PUT    /api/v1/devices/:id/status      # 修改设备状态
GET    /api/v1/devices/:id/credentials # 获取设备凭证
PUT    /api/v1/devices/:id/credentials # 重置设备凭证
GET    /api/v1/devices/:id/attributes  # 获取设备属性
PUT    /api/v1/devices/:id/attributes  # 更新设备属性
```

### 设备分组
```
GET    /api/v1/device-groups           # 获取分组列表
POST   /api/v1/device-groups           # 创建分组
GET    /api/v1/device-groups/:id       # 获取分组详情
PUT    /api/v1/device-groups/:id       # 更新分组
DELETE /api/v1/device-groups/:id       # 删除分组
GET    /api/v1/device-groups/:id/devices # 获取分组设备
POST   /api/v1/device-groups/:id/devices # 添加设备到分组
DELETE /api/v1/device-groups/:id/devices/:deviceId # 从分组移除设备
```

## 4.3 数据处理 API (data-service)

### 遥测数据
```
GET    /api/v1/devices/:id/telemetry   # 获取设备遥测数据
POST   /api/v1/devices/:id/telemetry   # 上报遥测数据（设备端）
GET    /api/v1/devices/:id/telemetry/latest # 获取最新遥测数据
GET    /api/v1/devices/:id/telemetry/:key # 获取指定 key 的遥测数据
DELETE /api/v1/devices/:id/telemetry   # 删除遥测数据
```

### 设备事件
```
GET    /api/v1/devices/:id/events      # 获取设备事件
POST   /api/v1/devices/:id/events      # 创建设备事件
GET    /api/v1/devices/:id/events/:eventId # 获取事件详情
```

### 设备命令
```
POST   /api/v1/devices/:id/commands    # 发送命令到设备
GET    /api/v1/devices/:id/commands    # 获取命令历史
GET    /api/v1/devices/:id/commands/:commandId # 获取命令详情
```

## 4.4 规则引擎 API (rule-engine)

### 规则管理
```
GET    /api/v1/rules                   # 获取规则列表
POST   /api/v1/rules                   # 创建规则
GET    /api/v1/rules/:id               # 获取规则详情
PUT    /api/v1/rules/:id               # 更新规则
DELETE /api/v1/rules/:id               # 删除规则
PUT    /api/v1/rules/:id/enable        # 启用规则
PUT    /api/v1/rules/:id/disable       # 禁用规则
POST   /api/v1/rules/:id/test          # 测试规则
GET    /api/v1/rules/:id/executions    # 获取规则执行历史
```

## 4.5 告警管理 API (alarm-service)

### 告警管理
```
GET    /api/v1/alarms                  # 获取告警列表
GET    /api/v1/alarms/:id              # 获取告警详情
PUT    /api/v1/alarms/:id/acknowledge  # 确认告警
PUT    /api/v1/alarms/:id/clear        # 清除告警
DELETE /api/v1/alarms/:id              # 删除告警
GET    /api/v1/alarms/statistics       # 获取告警统计
```

## 4.6 通知服务 API (notification-service)

### 通知配置
```
GET    /api/v1/notifications/configs   # 获取通知配置列表
POST   /api/v1/notifications/configs   # 创建通知配置
GET    /api/v1/notifications/configs/:id # 获取通知配置详情
PUT    /api/v1/notifications/configs/:id # 更新通知配置
DELETE /api/v1/notifications/configs/:id # 删除通知配置
```

### 通知历史
```
GET    /api/v1/notifications/history   # 获取通知历史
GET    /api/v1/notifications/history/:id # 获取通知详情
```

## 4.7 WebSocket API

### 实时数据推送
```
WS     /api/v1/ws/telemetry            # 订阅设备遥测数据
WS     /api/v1/ws/events               # 订阅设备事件
WS     /api/v1/ws/alarms               # 订阅告警
WS     /api/v1/ws/devices/status       # 订阅设备在线状态
```

## 4.8 设备端 API (protocol-http)

### HTTP 协议接入
```
POST   /api/v1/device/telemetry        # 设备上报遥测数据
POST   /api/v1/device/attributes       # 设备上报属性
POST   /api/v1/device/events           # 设备上报事件
GET    /api/v1/device/commands         # 设备获取命令
POST   /api/v1/device/commands/:id/ack # 设备确认命令
```

## API 通用规范

### 请求头
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
X-Tenant-ID: <tenant_id>  # 可选，多租户场景
```

### 响应格式
```json
{
  "code": 0,              // 0 表示成功，非 0 表示错误
  "message": "success",   // 响应消息
  "data": {},             // 响应数据
  "timestamp": 1234567890 // 时间戳
}
```

### 分页参数
```
?page=1&page_size=20&sort=created_at&order=desc
```

### 过滤参数
```
?status=active&search=keyword&start_time=xxx&end_time=xxx
```

### 错误码
```
0     - 成功
1000  - 参数错误
1001  - 资源不存在
1002  - 资源已存在
2000  - 未认证
2001  - 无权限
2002  - Token 过期
3000  - 内部错误
3001  - 数据库错误
3002  - 外部服务错误
```
