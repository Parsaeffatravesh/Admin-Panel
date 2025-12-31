'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { authApi, User, LoginResponse } from './api';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  const setAuthCookies = (accessToken: string, expiresInSeconds?: number) => {
    if (typeof document === 'undefined') return;
    const accessMaxAge = expiresInSeconds ?? 900;
    const secureFlag = window.location.protocol === 'https:' ? '; Secure' : '';
    document.cookie = `access_token=${encodeURIComponent(accessToken)}; Path=/; Max-Age=${accessMaxAge}; SameSite=Lax${secureFlag}`;
  };

  const clearAuthCookies = () => {
    if (typeof document === 'undefined') return;
    const secureFlag = window.location.protocol === 'https:' ? '; Secure' : '';
    document.cookie = `access_token=; Path=/; Max-Age=0; SameSite=Lax${secureFlag}`;
  };

  useEffect(() => {
    const initAuth = async () => {
      const storedUser = localStorage.getItem('user');
      const storedAccessToken = localStorage.getItem('access_token');
      const storedExpiresAt = localStorage.getItem('token_expires_at');
      
      if (storedUser) {
        try {
          setUser(JSON.parse(storedUser));
        } catch {
          localStorage.removeItem('user');
        }
      }

      if (storedAccessToken) {
        const expiresInSeconds = storedExpiresAt
          ? Math.max(0, Math.floor((Number(storedExpiresAt) - Date.now()) / 1000))
          : undefined;
        setAuthCookies(storedAccessToken, expiresInSeconds);
      }
      
      try {
        const meData = await authApi.me();
        if (meData) {
          const userData: User = {
            id: meData.user_id,
            tenant_id: meData.tenant_id,
            email: meData.email,
            first_name: '',
            last_name: '',
            status: 'active',
            created_at: '',
            updated_at: '',
          };
          setUser(userData);
          localStorage.setItem('user', JSON.stringify(userData));
        }
      } catch {
        localStorage.removeItem('user');
        setUser(null);
      }
      
      setIsLoading(false);
    };
    
    initAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const response: LoginResponse = await authApi.login(email, password);
    
    // Store tokens for subsequent requests
    if (response.tokens) {
      localStorage.setItem('access_token', response.tokens.access_token);
      localStorage.setItem(
        'token_expires_at',
        String(Date.now() + (response.tokens.expires_in || 0) * 1000)
      );
      setAuthCookies(response.tokens.access_token, response.tokens.expires_in);
    }

    localStorage.setItem('user', JSON.stringify(response.user));
    setUser(response.user);
    window.location.href = '/dashboard';
  };

  const logout = async () => {
    try {
      await authApi.logout();
    } catch {
    } finally {
      localStorage.removeItem('user');
      localStorage.removeItem('access_token');
      localStorage.removeItem('token_expires_at');
      clearAuthCookies();
      setUser(null);
      router.push('/login');
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
