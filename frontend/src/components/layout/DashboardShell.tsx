'use client';

import { Sidebar, MobileHeader } from '@/components/layout/sidebar';

export function DashboardShell({ children }: { children: React.ReactNode }) {
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
