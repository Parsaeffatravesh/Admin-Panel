-- Admin Panel Initial Schema
-- Multi-tenant ready with proper indexes for performance

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tenants table (for multi-tenant support)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_name ON roles(name);

-- Permissions table (global permissions)
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);

-- User-Role junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Role-Permission junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Audit logs table (immutable)
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);

-- Settings table (for feature flags and configuration)
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'string' CHECK (type IN ('string', 'number', 'boolean', 'json')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

CREATE INDEX idx_settings_tenant_id ON settings(tenant_id);
CREATE INDEX idx_settings_key ON settings(key);

-- Insert default permissions
INSERT INTO permissions (id, name, resource, action, description) VALUES
    (uuid_generate_v4(), 'users:read', 'users', 'read', 'View users'),
    (uuid_generate_v4(), 'users:create', 'users', 'create', 'Create users'),
    (uuid_generate_v4(), 'users:update', 'users', 'update', 'Update users'),
    (uuid_generate_v4(), 'users:delete', 'users', 'delete', 'Delete users'),
    (uuid_generate_v4(), 'roles:read', 'roles', 'read', 'View roles'),
    (uuid_generate_v4(), 'roles:create', 'roles', 'create', 'Create roles'),
    (uuid_generate_v4(), 'roles:update', 'roles', 'update', 'Update roles'),
    (uuid_generate_v4(), 'roles:delete', 'roles', 'delete', 'Delete roles'),
    (uuid_generate_v4(), 'audit:read', 'audit', 'read', 'View audit logs'),
    (uuid_generate_v4(), 'settings:read', 'settings', 'read', 'View settings'),
    (uuid_generate_v4(), 'settings:update', 'settings', 'update', 'Update settings'),
    (uuid_generate_v4(), 'dashboard:read', 'dashboard', 'read', 'View dashboard')
ON CONFLICT DO NOTHING;

-- Insert default tenant
INSERT INTO tenants (id, name, slug, status) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Default Organization', 'default', 'active')
ON CONFLICT DO NOTHING;

-- Insert default admin role
INSERT INTO roles (id, tenant_id, name, description, is_system) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Super Admin', 'Full system access', true)
ON CONFLICT DO NOTHING;

-- Assign all permissions to Super Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT DO NOTHING;

-- Insert default admin user (password: Admin123!)
INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, status) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'admin@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4.LT7mVwE1GRBhYq', 'Admin', 'User', 'active')
ON CONFLICT DO NOTHING;

-- Assign Super Admin role to admin user
INSERT INTO user_roles (user_id, role_id) VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;
