'use client';

import { useI18n } from '@/lib/i18n';
import { ThemeToggle } from '@/components/ThemeToggle';
import { LanguageToggle } from '@/components/LanguageToggle';

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  const { language } = useI18n();
  
  return (
    <header className="sticky top-0 z-30 flex h-16 w-full items-center justify-between border-b border-border/50 bg-background/80 px-4 backdrop-blur-md sm:px-6 lg:px-8 transition-all duration-300">
      <h1 className="text-xl font-bold tracking-tight text-foreground sm:text-2xl">
        {title}
      </h1>
      <div className="flex items-center gap-2 sm:gap-4">
        <div className="flex items-center bg-muted/50 rounded-lg p-1 gap-1">
          <ThemeToggle />
          <LanguageToggle />
        </div>
      </div>
    </header>
  );
}
