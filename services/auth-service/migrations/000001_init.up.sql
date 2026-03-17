-- 创建租户表
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    max_devices INTEGER DEFAULT 1000,
    max_users INTEGER DEFAULT 100,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_code ON tenants(code);
CREATE INDEX idx_tenants_status ON tenants(status);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
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

-- 创建角色表
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);

-- 创建权限表
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    code VARCHAR(200) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_resource ON permissions(resource);

-- 创建角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

-- 插入默认租户
INSERT INTO tenants (id, name, code, status, max_devices, max_users)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Tenant', 'default', 'active', 10000, 1000)
ON CONFLICT (code) DO NOTHING;

-- 插入默认管理员用户（密码：admin123）
INSERT INTO users (id, tenant_id, username, email, password_hash, full_name, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'admin',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- admin123
    'System Administrator',
    'active'
)
ON CONFLICT (tenant_id, username) DO NOTHING;

-- 插入系统角色
INSERT INTO roles (id, tenant_id, name, code, description, is_system)
VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Super Admin', 'super_admin', 'System super administrator', true),
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Admin', 'admin', 'Tenant administrator', true),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'User', 'user', 'Normal user', true)
ON CONFLICT (tenant_id, code) DO NOTHING;

-- 插入基础权限
INSERT INTO permissions (resource, action, code, description)
VALUES
    ('user', 'create', 'user:create', 'Create user'),
    ('user', 'read', 'user:read', 'Read user'),
    ('user', 'update', 'user:update', 'Update user'),
    ('user', 'delete', 'user:delete', 'Delete user'),
    ('device', 'create', 'device:create', 'Create device'),
    ('device', 'read', 'device:read', 'Read device'),
    ('device', 'update', 'device:update', 'Update device'),
    ('device', 'delete', 'device:delete', 'Delete device'),
    ('rule', 'create', 'rule:create', 'Create rule'),
    ('rule', 'read', 'rule:read', 'Read rule'),
    ('rule', 'update', 'rule:update', 'Update rule'),
    ('rule', 'delete', 'rule:delete', 'Delete rule')
ON CONFLICT (code) DO NOTHING;

-- 给管理员用户分配超级管理员角色
INSERT INTO user_roles (user_id, role_id)
VALUES ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;
