'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth';
import { Sidebar, MobileHeader } from '@/components/layout/sidebar';
import { useI18n } from '@/lib/i18n';
import { DashboardSkeleton } from '@/components/ui/skeleton';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const { t } = useI18n();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen">
        <Sidebar />
        <MobileHeader />
        <main className="flex-1 bg-background lg:ltr:ml-0 lg:rtl:mr-0 pt-14 lg:pt-0">
          <DashboardSkeleton />
        </main>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <MobileHeader />
      <main className="flex-1 bg-background lg:ltr:ml-0 lg:rtl:mr-0 pt-14 lg:pt-0">
        <div className="page-enter">
          {children}
        </div>
      </main>
    </div>
  );
}
