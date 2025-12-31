'use client';

import { ThemeProvider } from '@/hooks/useTheme';
import { I18nProvider } from '@/lib/i18n';

export function ThemeI18nProviders({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <I18nProvider>
        {children}
      </I18nProvider>
    </ThemeProvider>
  );
}
