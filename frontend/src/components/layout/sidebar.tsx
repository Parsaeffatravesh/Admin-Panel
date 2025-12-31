'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { cn } from '@/lib/utils';
import { useState, useEffect, useTransition } from 'react';
import {
  LayoutDashboard,
  Users,
  Shield,
  FileText,
  Settings,
  LogOut,
  Sun,
  Moon,
  Sparkles,
  Menu,
  X,
  Languages,
} from 'lucide-react';
import { useAuth } from '@/lib/auth';
import { useTheme, Theme } from '../../hooks/useTheme';
import { useI18n, Language } from '@/lib/i18n';

const navigationKeys = [
  { key: 'nav.dashboard', href: '/dashboard', icon: LayoutDashboard },
  { key: 'nav.users', href: '/dashboard/users', icon: Users },
  { key: 'nav.roles', href: '/dashboard/roles', icon: Shield },
  { key: 'nav.auditLogs', href: '/dashboard/audit', icon: FileText },
  { key: 'nav.settings', href: '/dashboard/settings', icon: Settings },
];

const themes: { name: Theme; labelKey: string; icon: typeof Sun }[] = [
  { name: 'light', labelKey: 'theme.light', icon: Sun },
  { name: 'dark', labelKey: 'theme.dark', icon: Moon },
  { name: 'legendary', labelKey: 'theme.legendary', icon: Sparkles },
];

const languages: { code: Language; label: string; flag: string }[] = [
  { code: 'en', label: 'EN', flag: '🇺🇸' },
  { code: 'fa', label: 'FA', flag: '🇮🇷' },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const { logout, user } = useAuth();
  const { theme, setTheme, mounted: themeMounted } = useTheme();
  const { t, language, setLanguage, isRTL, mounted: i18nMounted } = useI18n();
  const [isPending, startTransition] = useTransition();
  const [navigatingTo, setNavigatingTo] = useState<string | null>(null);

  useEffect(() => {
    if (!isPending && navigatingTo) {
      setNavigatingTo(null);
    }
  }, [isPending, navigatingTo]);

  const handleNavigation = (href: string) => {
    if (pathname === href) return;
    setNavigatingTo(href);
    startTransition(() => {
      router.push(href);
    });
  };

  const sidebarContent = (
    <>
      <div className="flex h-16 items-center justify-center border-b border-sidebar-border px-4">
        <h1 className="text-xl font-bold text-sidebar-foreground">{t('app.title')}</h1>
      </div>

      <nav className="flex-1 space-y-1 px-3 py-4 overflow-y-auto">
        {navigationKeys.map((item) => {
          const isActive = pathname === item.href || 
            (item.href !== '/dashboard' && pathname.startsWith(item.href));
          const isNavigating = navigatingTo === item.href;
          return (
            <button
              key={item.key}
              onClick={() => handleNavigation(item.href)}
              disabled={isNavigating}
              className={cn(
                'group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium',
                'transition-all duration-200',
                isNavigating
                  ? 'bg-sidebar-accent/30 text-sidebar-foreground/50 cursor-wait'
                  : isActive
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground shadow-sm'
                    : 'text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground'
              )}
            >
              <item.icon
                className={cn(
                  'h-5 w-5 flex-shrink-0 transition-all duration-200',
                  isNavigating
                    ? 'animate-pulse'
                    : isActive 
                      ? 'stroke-[2.5px]' 
                      : 'stroke-[1.5px] group-hover:stroke-[2px]'
                )}
              />
              <span className={cn(isNavigating && 'animate-pulse')}>
                {t(item.key)}
              </span>
              {isNavigating && (
                <span className="ltr:ml-auto rtl:mr-auto">
                  <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                </span>
              )}
            </button>
          );
        })}
      </nav>

      <div className="border-t border-sidebar-border p-3 space-y-3">
        <div>
          <p className="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2 px-1">
            {t('language.title')}
          </p>
          <div className="flex gap-1">
            {i18nMounted && languages.map((lang) => (
              <button
                key={lang.code}
                onClick={() => setLanguage(lang.code)}
                className={cn(
                  'flex-1 flex items-center justify-center gap-1.5 py-1.5 px-2 rounded-md text-xs font-medium',
                  'transition-all duration-200',
                  language === lang.code
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                    : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/30 hover:text-sidebar-foreground'
                )}
                title={lang.label}
              >
                <span>{lang.flag}</span>
                <span>{lang.label}</span>
              </button>
            ))}
          </div>
        </div>

        <div>
          <p className="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2 px-1">
            {t('theme.title')}
          </p>
          <div className="flex gap-1">
            {themeMounted && themes.map((t_item) => (
              <button
                key={t_item.name}
                onClick={() => setTheme(t_item.name)}
                className={cn(
                  'flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md text-xs font-medium',
                  'transition-all duration-200',
                  theme === t_item.name
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                    : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/30 hover:text-sidebar-foreground'
                )}
                title={t(t_item.labelKey)}
              >
                <t_item.icon className="h-3.5 w-3.5" />
              </button>
            ))}
          </div>
        </div>

        <div className="px-1">
          <p className="text-sm font-medium text-sidebar-foreground truncate">
            {user?.email}
          </p>
          <p className="text-xs text-sidebar-foreground/50 truncate">
            {user?.first_name} {user?.last_name}
          </p>
        </div>

        <button
          onClick={logout}
          className={cn(
            'group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium',
            'text-sidebar-foreground/70 hover:bg-destructive/10 hover:text-destructive',
            'transition-all duration-200'
          )}
        >
          <LogOut className="h-5 w-5 flex-shrink-0" />
          {t('auth.logout')}
        </button>
      </div>
    </>
  );

  return (
    <div className="hidden lg:flex h-screen w-64 flex-col bg-sidebar border-r border-sidebar-border rtl:border-r-0 rtl:border-l">
      {sidebarContent}
    </div>
  );
}

export function MobileHeader() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { t, isRTL } = useI18n();

  return (
    <>
      <header className="lg:hidden fixed top-0 left-0 right-0 z-30 bg-sidebar border-b border-sidebar-border h-14 flex items-center px-4">
        <button
          onClick={() => setSidebarOpen(true)}
          className="p-2 rounded-lg hover:bg-sidebar-accent/50 text-sidebar-foreground/70"
        >
          <Menu className="h-5 w-5" />
        </button>
        <h1 className="text-lg font-bold text-sidebar-foreground mx-auto">{t('app.title')}</h1>
        <div className="w-9" />
      </header>
      <MobileSidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
    </>
  );
}

function MobileSidebar({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) {
  const pathname = usePathname();
  const router = useRouter();
  const { logout, user } = useAuth();
  const { theme, setTheme, mounted: themeMounted } = useTheme();
  const { t, language, setLanguage, isRTL, mounted: i18nMounted } = useI18n();
  const [isPending, startTransition] = useTransition();
  const [navigatingTo, setNavigatingTo] = useState<string | null>(null);

  useEffect(() => {
    if (!isPending && navigatingTo) {
      setNavigatingTo(null);
    }
  }, [isPending, navigatingTo]);

  const handleNavigation = (href: string) => {
    if (pathname === href) return;
    setNavigatingTo(href);
    startTransition(() => {
      router.push(href);
    });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <>
      <div 
        className="fixed inset-0 bg-black/50 z-40 lg:hidden"
        onClick={onClose}
      />
      <div className={cn(
        "fixed inset-y-0 z-50 w-72 bg-sidebar flex flex-col lg:hidden",
        "transition-transform duration-300 ease-out",
        isRTL ? "right-0" : "left-0"
      )}>
        <div className="flex h-16 items-center justify-between border-b border-sidebar-border px-4">
          <h1 className="text-xl font-bold text-sidebar-foreground">{t('app.title')}</h1>
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-sidebar-accent/50 text-sidebar-foreground/70"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 space-y-1 px-3 py-4 overflow-y-auto">
          {navigationKeys.map((item) => {
            const isActive = pathname === item.href || 
              (item.href !== '/dashboard' && pathname.startsWith(item.href));
            const isNavigating = navigatingTo === item.href;
            return (
              <button
                key={item.key}
                onClick={() => handleNavigation(item.href)}
                disabled={isNavigating}
                className={cn(
                  'group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium',
                  'transition-all duration-200',
                  isNavigating
                    ? 'bg-sidebar-accent/30 text-sidebar-foreground/50 cursor-wait'
                    : isActive
                      ? 'bg-sidebar-accent text-sidebar-accent-foreground shadow-sm'
                      : 'text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground'
                )}
              >
                <item.icon
                  className={cn(
                    'h-5 w-5 flex-shrink-0 transition-all duration-200',
                    isNavigating
                      ? 'animate-pulse'
                      : isActive 
                        ? 'stroke-[2.5px]' 
                        : 'stroke-[1.5px] group-hover:stroke-[2px]'
                  )}
                />
                <span className={cn(isNavigating && 'animate-pulse')}>
                  {t(item.key)}
                </span>
                {isNavigating && (
                  <span className="ltr:ml-auto rtl:mr-auto">
                    <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  </span>
                )}
              </button>
            );
          })}
        </nav>

        <div className="border-t border-sidebar-border p-3 space-y-3">
          <div>
            <p className="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2 px-1">
              {t('language.title')}
            </p>
            <div className="flex gap-1">
              {i18nMounted && languages.map((lang) => (
                <button
                  key={lang.code}
                  onClick={() => setLanguage(lang.code)}
                  className={cn(
                    'flex-1 flex items-center justify-center gap-1.5 py-1.5 px-2 rounded-md text-xs font-medium',
                    'transition-all duration-200',
                    language === lang.code
                      ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                      : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/30 hover:text-sidebar-foreground'
                  )}
                  title={lang.label}
                >
                  <span>{lang.flag}</span>
                  <span>{lang.label}</span>
                </button>
              ))}
            </div>
          </div>

          <div>
            <p className="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2 px-1">
              {t('theme.title')}
            </p>
            <div className="flex gap-1">
              {themeMounted && themes.map((t_item) => (
                <button
                  key={t_item.name}
                  onClick={() => setTheme(t_item.name)}
                  className={cn(
                    'flex-1 flex items-center justify-center gap-1 py-1.5 px-2 rounded-md text-xs font-medium',
                    'transition-all duration-200',
                    theme === t_item.name
                      ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                      : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/30 hover:text-sidebar-foreground'
                  )}
                  title={t(t_item.labelKey)}
                >
                  <t_item.icon className="h-3.5 w-3.5" />
                </button>
              ))}
            </div>
          </div>

          <div className="px-1">
            <p className="text-sm font-medium text-sidebar-foreground truncate">
              {user?.email}
            </p>
            <p className="text-xs text-sidebar-foreground/50 truncate">
              {user?.first_name} {user?.last_name}
            </p>
          </div>

          <button
            onClick={logout}
            className={cn(
              'group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium',
              'text-sidebar-foreground/70 hover:bg-destructive/10 hover:text-destructive',
              'transition-all duration-200'
            )}
          >
            <LogOut className="h-5 w-5 flex-shrink-0" />
            {t('auth.logout')}
          </button>
        </div>
      </div>
    </>
  );
}
