# Database Design

## 3.1 认证授权模块 (auth-service)

### tenants - 租户表
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, suspended, deleted
    max_devices INTEGER DEFAULT 1000,
    max_users INTEGER DEFAULT 100,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_code ON tenants(code);
CREATE INDEX idx_tenants_status ON tenants(status);
```

### users - 用户表
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, locked
    last_login_at TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, username),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

### roles - 角色表
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE, -- 系统内置角色不可删除
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
```

### permissions - 权限表
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL, -- device, user, rule, alarm, etc.
    action VARCHAR(50) NOT NULL,    -- create, read, update, delete, execute
    code VARCHAR(200) UNIQUE NOT NULL, -- device:create, user:read, etc.
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
```

### role_permissions - 角色权限关联表
```sql
CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);
```

### user_roles - 用户角色关联表
```sql
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
```

## 3.2 设备管理模块 (device-service)

### device_products - 设备产品表
```sql
CREATE TABLE device_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    protocol VARCHAR(50) NOT NULL, -- mqtt, http, coap
    data_format VARCHAR(50) NOT NULL, -- json, protobuf, custom
    attributes JSONB, -- 产品属性定义
    telemetry JSONB,  -- 遥测数据定义
    commands JSONB,   -- 命令定义
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_device_products_tenant_id ON device_products(tenant_id);
```

### devices - 设备表
```sql
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES device_products(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    device_key VARCHAR(100) UNIQUE NOT NULL, -- 设备唯一标识
    secret VARCHAR(255) NOT NULL, -- 设备密钥
    status VARCHAR(20) NOT NULL DEFAULT 'inactive', -- active, inactive, disabled
    online_status VARCHAR(20) NOT NULL DEFAULT 'offline', -- online, offline
    last_online_at TIMESTAMP,
    last_offline_at TIMESTAMP,
    attributes JSONB, -- 当前属性值
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_devices_tenant_id ON devices(tenant_id);
CREATE INDEX idx_devices_product_id ON devices(product_id);
CREATE INDEX idx_devices_device_key ON devices(device_key);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_online_status ON devices(online_status);
```

### device_groups - 设备分组表
```sql
CREATE TABLE device_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES device_groups(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_groups_tenant_id ON device_groups(tenant_id);
CREATE INDEX idx_device_groups_parent_id ON device_groups(parent_id);
```

### device_group_members - 设备分组成员表
```sql
CREATE TABLE device_group_members (
    group_id UUID NOT NULL REFERENCES device_groups(id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (group_id, device_id)
);
```

## 3.3 数据处理模块 (data-service)

### device_telemetry - 设备遥测数据表 (TimescaleDB)
```sql
CREATE TABLE device_telemetry (
    time TIMESTAMPTZ NOT NULL,
    device_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    key VARCHAR(100) NOT NULL,
    value_type VARCHAR(20) NOT NULL, -- string, number, boolean, json
    value_string TEXT,
    value_number DOUBLE PRECISION,
    value_boolean BOOLEAN,
    value_json JSONB,
    PRIMARY KEY (device_id, key, time)
);

-- 转换为 TimescaleDB 超表
SELECT create_hypertable('device_telemetry', 'time');

-- 创建索引
CREATE INDEX idx_device_telemetry_device_id ON device_telemetry(device_id, time DESC);
CREATE INDEX idx_device_telemetry_tenant_id ON device_telemetry(tenant_id, time DESC);
CREATE INDEX idx_device_telemetry_key ON device_telemetry(key, time DESC);
```

### device_events - 设备事件表
```sql
CREATE TABLE device_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL, -- connection, disconnection, error, custom
    severity VARCHAR(20) NOT NULL, -- info, warning, error, critical
    message TEXT,
    data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_events_device_id ON device_events(device_id, created_at DESC);
CREATE INDEX idx_device_events_tenant_id ON device_events(tenant_id, created_at DESC);
CREATE INDEX idx_device_events_type ON device_events(event_type, created_at DESC);
```

## 3.4 规则引擎模块 (rule-engine)

### rules - 规则表
```sql
CREATE TABLE rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT TRUE,
    trigger_type VARCHAR(50) NOT NULL, -- telemetry, event, schedule
    trigger_config JSONB NOT NULL, -- 触发条件配置
    actions JSONB NOT NULL, -- 动作配置数组
    priority INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rules_tenant_id ON rules(tenant_id);
CREATE INDEX idx_rules_enabled ON rules(enabled);
```

### rule_executions - 规则执行记录表
```sql
CREATE TABLE rule_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    device_id UUID,
    status VARCHAR(20) NOT NULL, -- success, failed, partial
    input_data JSONB,
    output_data JSONB,
    error_message TEXT,
    execution_time_ms INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rule_executions_rule_id ON rule_executions(rule_id, created_at DESC);
CREATE INDEX idx_rule_executions_device_id ON rule_executions(device_id, created_at DESC);
```

## 3.5 告警模块 (alarm-service)

### alarms - 告警表
```sql
CREATE TABLE alarms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES rules(id) ON DELETE SET NULL,
    severity VARCHAR(20) NOT NULL, -- info, warning, error, critical
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, acknowledged, cleared
    title VARCHAR(255) NOT NULL,
    message TEXT,
    data JSONB,
    acknowledged_by UUID REFERENCES users(id) ON DELETE SET NULL,
    acknowledged_at TIMESTAMP,
    cleared_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_alarms_tenant_id ON alarms(tenant_id, created_at DESC);
CREATE INDEX idx_alarms_device_id ON alarms(device_id, created_at DESC);
CREATE INDEX idx_alarms_status ON alarms(status, created_at DESC);
CREATE INDEX idx_alarms_severity ON alarms(severity, created_at DESC);
```

## 数据库初始化顺序

1. 创建 PostgreSQL 数据库
2. 安装 TimescaleDB 扩展
3. 执行认证授权模块表
4. 插入系统初始数据（默认租户、管理员用户、系统角色、权限）
5. 执行设备管理模块表
6. 执行数据处理模块表（包括 TimescaleDB 超表）
7. 执行规则引擎模块表
8. 执行告警模块表
