'use client';

import * as React from 'react';
import { Moon, Sun, Sparkles } from 'lucide-react';
import { useTheme } from 'next-themes';

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = React.useState(false);

  React.useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return <div className="w-9 h-9" />;

  const cycleTheme = () => {
    if (theme === 'light') setTheme('dark');
    else if (theme === 'dark') setTheme('legendary');
    else setTheme('light');
  };

  return (
    <button
      onClick={cycleTheme}
      className="relative flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-background hover:bg-accent hover:text-accent-foreground transition-colors"
      title="Change theme"
    >
      <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0 legendary:-rotate-90 legendary:scale-0" />
      <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100 legendary:scale-0" />
      <Sparkles className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all legendary:rotate-0 legendary:scale-100 text-primary" />
      <span className="sr-only">Toggle theme</span>
    </button>
  );
}
