'use client';

import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from 'react';

export type Theme = 'light' | 'dark' | 'legendary';

const THEME_STORAGE_KEY = 'admin-panel-theme';

interface ThemeContextType {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  mounted: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const applyTheme = (newTheme: Theme) => {
  const root = document.documentElement;
  root.classList.remove('light', 'dark', 'legendary');
  if (newTheme !== 'light') {
    root.classList.add(newTheme);
  }
  root.setAttribute('data-theme', newTheme);
};

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>('light');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem(THEME_STORAGE_KEY) as Theme | null;
    if (stored && ['light', 'dark', 'legendary'].includes(stored)) {
      setThemeState(stored);
      applyTheme(stored);
    }
    setMounted(true);
  }, []);

  const setTheme = useCallback((newTheme: Theme) => {
    const root = document.documentElement;

    if (!root.classList.contains('theme-transition')) {
      root.classList.add('theme-transition');
    }

    applyTheme(newTheme);
    localStorage.setItem(THEME_STORAGE_KEY, newTheme);
    setThemeState(newTheme);

    setTimeout(() => {
      root.classList.remove('theme-transition');
    }, 350);
  }, []);

  return (
    <ThemeContext.Provider value={{ theme, setTheme, mounted }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    return { theme: 'light' as Theme, setTheme: () => {}, mounted: false };
  }
  return context;
}
