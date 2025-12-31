-- Performance Optimization Indexes

-- Composite indexes for common query patterns

-- Users search with tenant and status filter
CREATE INDEX IF NOT EXISTS idx_users_tenant_status ON users(tenant_id, status);

-- Audit logs with date range and tenant filter
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_created ON audit_logs(tenant_id, created_at DESC);

-- Audit logs with action and resource filter
CREATE INDEX IF NOT EXISTS idx_audit_logs_action_resource ON audit_logs(action, resource, created_at DESC);

-- Roles with tenant for quick lookup
CREATE INDEX IF NOT EXISTS idx_roles_tenant_system ON roles(tenant_id, is_system);

-- User roles with role lookup
CREATE INDEX IF NOT EXISTS idx_user_roles_role_user ON user_roles(role_id, user_id);
