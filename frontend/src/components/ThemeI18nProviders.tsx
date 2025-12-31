'use client';

import { ThemeProvider as NextThemesProvider } from 'next-themes';
import { I18nProvider } from '@/lib/i18n';

export function ThemeI18nProviders({ children }: { children: React.ReactNode }) {
  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
      themes={['light', 'dark', 'legendary']}
    >
      <I18nProvider>
        {children}
      </I18nProvider>
    </NextThemesProvider>
  );
}
