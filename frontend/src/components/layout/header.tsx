'use client';

import { Bell, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-card px-6 transition-colors duration-200">
      <h2 className="text-xl font-semibold text-foreground">{title}</h2>
      
      <div className="flex items-center gap-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search..."
            className="w-64 pl-9 transition-colors duration-200"
          />
        </div>
        
        <Button variant="ghost" size="icon" className="transition-colors duration-200">
          <Bell className="h-5 w-5" />
        </Button>
      </div>
    </header>
  );
}
