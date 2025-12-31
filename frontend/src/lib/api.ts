const API_URL = process.env.NEXT_PUBLIC_API_URL || '';

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const isAuthEndpoint = endpoint.includes('/auth/login') || endpoint.includes('/auth/refresh');
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include',
    });

    if (response.status === 401 && !isAuthEndpoint) {
      handleAuthFailure();
      throw new Error('Unauthorized. Please sign in again.');
    }

    const contentType = response.headers.get('content-type') || '';

    if (contentType.includes('application/json')) {
      try {
        const data: ApiResponse<T> = await response.json();
        if (!response.ok || !data.success) {
          throw new ApiError(response.status, data.error?.message || `Server returned ${response.status}`);
        }
        return data.data as T;
      } catch (err) {
        if (err instanceof ApiError) {
          throw err;
        }
        throw new ApiError(response.status, 'Failed to parse server response');
      }
    }

    await response.text();
    throw new ApiError(response.status, `Server returned ${response.status}`);
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const api = new ApiClient(API_URL);

const handleAuthFailure = () => {
  if (typeof window === 'undefined') return;
  localStorage.removeItem('user');
  window.alert('نشست شما منقضی شده است. لطفاً دوباره وارد شوید.');
  window.location.href = '/login';
};

export interface User {
  id: string;
  tenant_id: string;
  email: string;
  first_name: string;
  last_name: string;
  status: 'active' | 'inactive' | 'suspended';
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

export interface Role {
  id: string;
  tenant_id: string;
  name: string;
  description: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
  permissions?: string[];
}

export interface Permission {
  id: string;
  name: string;
  resource: string;
  action: string;
  description: string;
  created_at: string;
}

export interface AuditLog {
  id: string;
  tenant_id: string;
  user_id?: string;
  action: string;
  resource: string;
  resource_id?: string;
  old_value?: string;
  new_value?: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface LoginResponse {
  user: User;
  tokens: {
    access_token: string;
    refresh_token: string;
    expires_in: number;
  };
}

export interface DashboardStats {
  total_users: number;
  active_users: number;
  total_roles: number;
  recent_logins: number;
  users_by_status: Record<string, number>;
  recent_activity: Array<{
    action: string;
    resource: string;
    user_email: string;
    created_at: string;
  }>;
}

export const authApi = {
  login: (email: string, password: string) =>
    api.post<LoginResponse>('/api/v1/auth/login', { email, password }),
  logout: () => api.post('/api/v1/auth/logout'),
  me: () => api.get<{ user_id: string; tenant_id: string; email: string }>('/api/v1/auth/me'),
  refresh: (refreshToken?: string) =>
    api.post<{ access_token: string; refresh_token: string; expires_in: number }>(
      '/api/v1/auth/refresh',
      refreshToken ? { refresh_token: refreshToken } : undefined
    ),
};

export const usersApi = {
  list: (params?: { page?: number; per_page?: number; search?: string; status?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params?.search) searchParams.set('search', params.search);
    if (params?.status) searchParams.set('status', params.status);
    return api.get<PaginatedResponse<User>>(`/api/v1/users?${searchParams}`);
  },
  get: (id: string) => api.get<User>(`/api/v1/users/${id}`),
  create: (data: { email: string; password: string; first_name: string; last_name: string; role_ids?: string[] }) =>
    api.post<User>('/api/v1/users', data),
  update: (id: string, data: Partial<{ email: string; first_name: string; last_name: string; status: string; role_ids: string[] }>) =>
    api.put<User>(`/api/v1/users/${id}`, data),
  delete: (id: string) => api.delete(`/api/v1/users/${id}`),
  resetPassword: (id: string, newPassword: string) =>
    api.post(`/api/v1/users/${id}/reset-password`, { new_password: newPassword }),
  getRoles: (id: string) => api.get<Role[]>(`/api/v1/users/${id}/roles`),
};

export const rolesApi = {
  list: (params?: { page?: number; per_page?: number; search?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params?.search) searchParams.set('search', params.search);
    return api.get<PaginatedResponse<Role>>(`/api/v1/roles?${searchParams}`);
  },
  get: (id: string) => api.get<Role>(`/api/v1/roles/${id}`),
  create: (data: { name: string; description: string; permission_ids?: string[] }) =>
    api.post<Role>('/api/v1/roles', data),
  update: (id: string, data: Partial<{ name: string; description: string; permission_ids: string[] }>) =>
    api.put<Role>(`/api/v1/roles/${id}`, data),
  delete: (id: string) => api.delete(`/api/v1/roles/${id}`),
  getPermissions: (id: string) => api.get<Permission[]>(`/api/v1/roles/${id}/permissions`),
};

export const permissionsApi = {
  list: () => api.get<Permission[]>('/api/v1/permissions'),
};

export const auditApi = {
  list: (params?: { page?: number; per_page?: number; search?: string; action?: string; resource?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params?.search) searchParams.set('search', params.search);
    if (params?.action) searchParams.set('action', params.action);
    if (params?.resource) searchParams.set('resource', params.resource);
    return api.get<PaginatedResponse<AuditLog>>(`/api/v1/audit-logs?${searchParams}`);
  },
};

export const dashboardApi = {
  getStats: () => api.get<DashboardStats>('/api/v1/dashboard/stats'),
};

export interface FeatureFlag {
  id: string;
  tenant_id: string;
  key: string;
  name: string;
  description: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export const featureFlagsApi = {
  list: (params?: { page?: number; per_page?: number; search?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params?.search) searchParams.set('search', params.search);
    return api.get<PaginatedResponse<FeatureFlag>>(`/api/v1/feature-flags?${searchParams}`);
  },
  get: (id: string) => api.get<FeatureFlag>(`/api/v1/feature-flags/${id}`),
  create: (data: { key: string; name: string; description?: string; enabled?: boolean }) =>
    api.post<FeatureFlag>('/api/v1/feature-flags', data),
  update: (id: string, data: Partial<{ name: string; description: string; enabled: boolean }>) =>
    api.put<FeatureFlag>(`/api/v1/feature-flags/${id}`, data),
  delete: (id: string) => api.delete(`/api/v1/feature-flags/${id}`),
  toggle: (id: string) => api.post<FeatureFlag>(`/api/v1/feature-flags/${id}/toggle`),
};

export const adminApi = {
  setAdmin: (userId: string, data: { is_admin: boolean; admin_password?: string }) =>
    api.post(`/api/v1/users/${userId}/set-admin`, data),
  getAdminStatus: (userId: string) =>
    api.get<{ is_admin: boolean }>(`/api/v1/users/${userId}/admin-status`),
};

export const getExportUrl = (endpoint: string) => {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || '';
  return `${baseUrl}${endpoint}`;
};
