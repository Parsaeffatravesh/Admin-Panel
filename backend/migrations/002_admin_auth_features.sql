-- Admin Authentication and Feature Flags Migration

-- Admin Auth table (separate from user password for admin panel access)
CREATE TABLE IF NOT EXISTS admin_auth (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    admin_password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    enabled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_auth_is_admin ON admin_auth(is_admin) WHERE is_admin = TRUE;

-- Feature Flags table
CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

CREATE INDEX idx_feature_flags_tenant_id ON feature_flags(tenant_id);
CREATE INDEX idx_feature_flags_key ON feature_flags(key);
CREATE INDEX idx_feature_flags_enabled ON feature_flags(enabled);

-- Add admin_auth permission
INSERT INTO permissions (id, name, resource, action, description) VALUES
    (uuid_generate_v4(), 'admin:manage', 'admin', 'manage', 'Manage admin access')
ON CONFLICT DO NOTHING;

-- Add feature_flags permissions
INSERT INTO permissions (id, name, resource, action, description) VALUES
    (uuid_generate_v4(), 'feature_flags:read', 'feature_flags', 'read', 'View feature flags'),
    (uuid_generate_v4(), 'feature_flags:create', 'feature_flags', 'create', 'Create feature flags'),
    (uuid_generate_v4(), 'feature_flags:update', 'feature_flags', 'update', 'Update feature flags'),
    (uuid_generate_v4(), 'feature_flags:delete', 'feature_flags', 'delete', 'Delete feature flags')
ON CONFLICT DO NOTHING;

-- Assign new permissions to Super Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions WHERE resource IN ('admin', 'feature_flags')
ON CONFLICT DO NOTHING;

-- Create admin auth entry for existing admin user
INSERT INTO admin_auth (user_id, admin_password_hash, is_admin, enabled_at)
SELECT id, password_hash, TRUE, NOW() FROM users WHERE email = 'admin@example.com'
ON CONFLICT DO NOTHING;

-- Add GIN index on audit_logs metadata for filtering
CREATE INDEX IF NOT EXISTS idx_audit_logs_metadata ON audit_logs USING GIN (old_value);
CREATE INDEX IF NOT EXISTS idx_audit_logs_new_value ON audit_logs USING GIN (new_value);
