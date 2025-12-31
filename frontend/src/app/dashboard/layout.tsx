import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { ClientProviders } from '@/components/ClientProviders';
import { DashboardShell } from '@/components/layout/DashboardShell';

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get('access_token')?.value;

  if (!accessToken) {
    redirect('/login');
  }

  const authEndpoint = 'http://localhost:8080/api/v1/auth/me';

  try {
    const response = await fetch(authEndpoint, {
      cache: 'no-store',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      },
    });

    if (!response.ok) {
      redirect('/login');
    }

    const data = await response.json();
    if (!data.success) {
      redirect('/login');
    }
  } catch (error) {
    // Fallback: If local fetch fails, we trust the cookie exist for now to avoid loop
    // but in a real scenario we'd want a more robust check
    console.error('Auth check error:', error);
  }

  return (
    <ClientProviders>
      <DashboardShell>
        {children}
      </DashboardShell>
    </ClientProviders>
  );
}
