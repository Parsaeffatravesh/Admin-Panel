'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { ApiError, authApi, User, LoginResponse } from './api';

type RefreshTokenStrategy = 'cookie' | 'storage';

// Strategy decision: rely on HttpOnly refresh_token cookies rather than storing them in local storage.
const REFRESH_TOKEN_STRATEGY: RefreshTokenStrategy = 'cookie';
const REFRESH_TOKEN_KEY = 'refresh_token';
const TOKEN_METADATA_KEY = 'token_metadata';

interface TokenMetadata {
  accessTokenExpiresAt: number;
  refreshToken?: string;
}

let refreshPromise: Promise<LoginResponse['tokens'] | null> | null = null;

const persistTokenMetadata = (tokens: LoginResponse['tokens']) => {
  if (typeof window === 'undefined') return;

  const metadata: TokenMetadata = {
    accessTokenExpiresAt: Date.now() + tokens.expires_in * 1000,
  };

  if (REFRESH_TOKEN_STRATEGY === 'storage') {
    metadata.refreshToken = tokens.refresh_token;
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
  } else {
    localStorage.removeItem(REFRESH_TOKEN_KEY);
  }

  localStorage.setItem(TOKEN_METADATA_KEY, JSON.stringify(metadata));
};

const clearTokenMetadata = () => {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(TOKEN_METADATA_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
};

const getStoredRefreshToken = () => {
  if (typeof window === 'undefined') return null;
  if (REFRESH_TOKEN_STRATEGY !== 'storage') return null;
  return localStorage.getItem(REFRESH_TOKEN_KEY);
};

export const refreshAccessTokenSingleFlight = async () => {
  if (typeof window === 'undefined') return null;
  if (refreshPromise) return refreshPromise;

  refreshPromise = (async () => {
    const refreshToken = getStoredRefreshToken();

    if (REFRESH_TOKEN_STRATEGY === 'storage' && !refreshToken) {
      clearTokenMetadata();
      return null;
    }

    const refreshedTokens = await authApi.refresh(refreshToken || undefined);
    persistTokenMetadata(refreshedTokens);
    return refreshedTokens;
  })()
    .catch((error) => {
      clearTokenMetadata();
      throw error;
    })
    .finally(() => {
      refreshPromise = null;
    });

  return refreshPromise;
};

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

  useEffect(() => {
    const initAuth = async () => {
      const storedUser = localStorage.getItem('user');

      if (storedUser) {
        try {
          setUser(JSON.parse(storedUser));
        } catch {
          localStorage.removeItem('user');
        }
      }

      const hasSession = () => {
        const cookieString = typeof document !== 'undefined' ? document.cookie : '';
        const hasAccessCookie = cookieString.split(';').some(cookie => cookie.trim().startsWith('access_token='));
        const hasRefreshCookie = cookieString.split(';').some(cookie => cookie.trim().startsWith('refresh_token='));
        const hasStoredRefresh =
          REFRESH_TOKEN_STRATEGY === 'storage' && !!localStorage.getItem(REFRESH_TOKEN_KEY);

        return hasAccessCookie || hasRefreshCookie || hasStoredRefresh;
      };

      if (!hasSession()) {
        localStorage.removeItem('user');
        setUser(null);
        setIsLoading(false);
        return;
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
      } catch (error) {
        if (error instanceof ApiError && error.status === 401) {
          clearTokenMetadata();
        }
        localStorage.removeItem('user');
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };
    
    initAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const response: LoginResponse = await authApi.login(email, password);

    persistTokenMetadata(response.tokens);
    localStorage.setItem('user', JSON.stringify(response.user));
    setUser(response.user);
    window.location.href = '/dashboard';
  };

  const logout = async () => {
    try {
      await authApi.logout();
    } catch {
    } finally {
      clearTokenMetadata();
      localStorage.removeItem('user');
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
