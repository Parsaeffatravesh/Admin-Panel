'use client';

import { Bell, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useI18n } from '@/lib/i18n';

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  const { t, isRTL } = useI18n();
  
  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-card px-4 sm:px-6 transition-colors duration-200">
      <h2 className="text-lg sm:text-xl font-semibold text-foreground">{title}</h2>
      
      <div className="flex items-center gap-2 sm:gap-4">
        <div className="relative hidden sm:block">
          <Search className="absolute ltr:left-3 rtl:right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder={t('common.search')}
            className="w-48 sm:w-64 ltr:pl-9 rtl:pr-9 transition-colors duration-200"
            dir="ltr"
          />
        </div>
        
        <Button variant="ghost" size="icon" className="transition-colors duration-200 sm:hidden">
          <Search className="h-5 w-5" />
        </Button>
        
        <Button variant="ghost" size="icon" className="transition-colors duration-200">
          <Bell className="h-5 w-5" />
        </Button>
      </div>
    </header>
  );
}
