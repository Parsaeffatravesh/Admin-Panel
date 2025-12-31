'use client';

import { useState } from 'react';
import { useAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { toast } from 'sonner';
import { useI18n, Language } from '@/lib/i18n';
import { cn } from '@/lib/utils';
import { Copy, Check, Eye, EyeOff } from 'lucide-react';

const languages: { code: Language; label: string; flag: string }[] = [
  { code: 'en', label: 'EN', flag: '🇺🇸' },
  { code: 'fa', label: 'FA', flag: '🇮🇷' },
];

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [copiedEmail, setCopiedEmail] = useState(false);
  const [copiedPassword, setCopiedPassword] = useState(false);
  const { login } = useAuth();
  const { t, language, setLanguage, mounted } = useI18n();

  const copyToClipboard = async (text: string, type: 'email' | 'password') => {
    try {
      await navigator.clipboard.writeText(text);
      if (type === 'email') {
        setCopiedEmail(true);
        setTimeout(() => setCopiedEmail(false), 2000);
      } else {
        setCopiedPassword(true);
        setTimeout(() => setCopiedPassword(false), 2000);
      }
      toast.success(language === 'fa' ? 'کپی شد' : 'Copied');
    } catch {
      toast.error(language === 'fa' ? 'خطا در کپی' : 'Failed to copy');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      // Trim inputs to avoid whitespace issues
      await login(email.trim(), password.trim());
      toast.success(language === 'fa' ? 'ورود موفق' : 'Logged in successfully');
    } catch (err) {
      const message = err instanceof Error ? err.message : (language === 'fa' ? 'ورود ناموفق' : 'Login failed');
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background transition-colors duration-200 p-4">
      <div className="fixed top-4 ltr:right-4 rtl:left-4 flex gap-1">
        {mounted && languages.map((lang) => (
          <button
            key={lang.code}
            onClick={() => setLanguage(lang.code)}
            className={cn(
              'flex items-center gap-1.5 py-1.5 px-3 rounded-lg text-sm font-medium',
              'transition-all duration-200 border',
              language === lang.code
                ? 'bg-primary text-primary-foreground border-primary'
                : 'bg-card text-card-foreground border-border hover:bg-accent'
            )}
          >
            <span>{lang.flag}</span>
            <span>{lang.label}</span>
          </button>
        ))}
      </div>

      <Card className="w-full max-w-md animate-in fade-in slide-in-from-bottom-2">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-bold">{t('app.title')}</CardTitle>
          <CardDescription>{t('app.signIn')}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive border border-destructive/20">
                {error}
              </div>
            )}
            <div className="space-y-2">
              <label htmlFor="email" className="text-sm font-medium text-foreground">
                {t('auth.email')}
              </label>
              <Input
                id="email"
                type="email"
                placeholder="admin@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="transition-colors duration-200"
                dir="ltr"
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium text-foreground">
                {t('auth.password')}
              </label>
              <div className="relative group">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder={language === 'fa' ? 'رمز عبور را وارد کنید' : 'Enter your password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="transition-colors duration-200 pr-10 rtl:pl-10 rtl:pr-3"
                  dir="ltr"
                />
                {password.length > 0 && (
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute inset-y-0 right-0 flex items-center pr-3 rtl:left-0 rtl:right-auto rtl:pl-3 text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {showPassword ? (
                      <EyeOff className="h-4 w-4" />
                    ) : (
                      <Eye className="h-4 w-4" />
                    )}
                  </button>
                )}
              </div>
            </div>
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading 
                ? (language === 'fa' ? 'در حال ورود...' : 'Signing in...') 
                : t('auth.signIn')
              }
            </Button>
            <div className="text-center text-sm text-muted-foreground space-y-2">
              <p className="font-medium">{t('auth.demo')}</p>
              <div className="flex items-center justify-center gap-2 flex-wrap">
                <button
                  type="button"
                  onClick={() => copyToClipboard('admin@example.com', 'email')}
                  className={cn(
                    'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-mono',
                    'bg-muted hover:bg-muted/80 transition-colors duration-150',
                    'border border-border'
                  )}
                >
                  admin@example.com
                  {copiedEmail ? (
                    <Check className="h-3 w-3 text-green-500" />
                  ) : (
                    <Copy className="h-3 w-3 opacity-50" />
                  )}
                </button>
                <span>/</span>
                <button
                  type="button"
                  onClick={() => copyToClipboard('Admin123!', 'password')}
                  className={cn(
                    'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-mono',
                    'bg-muted hover:bg-muted/80 transition-colors duration-150',
                    'border border-border'
                  )}
                >
                  Admin123!
                  {copiedPassword ? (
                    <Check className="h-3 w-3 text-green-500" />
                  ) : (
                    <Copy className="h-3 w-3 opacity-50" />
                  )}
                </button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
