'use client';

import { useState } from 'react';
import { useAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { toast } from 'sonner';
import { useI18n, Language } from '@/lib/i18n';
import { cn } from '@/lib/utils';

const languages: { code: Language; label: string; flag: string }[] = [
  { code: 'en', label: 'EN', flag: '🇺🇸' },
  { code: 'fa', label: 'FA', flag: '🇮🇷' },
];

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const { t, language, setLanguage, mounted } = useI18n();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
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
              <Input
                id="password"
                type="password"
                placeholder={language === 'fa' ? 'رمز عبور را وارد کنید' : 'Enter your password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="transition-colors duration-200"
                dir="ltr"
              />
            </div>
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading 
                ? (language === 'fa' ? 'در حال ورود...' : 'Signing in...') 
                : t('auth.signIn')
              }
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              {t('auth.demo')} admin@example.com / Admin123!
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
