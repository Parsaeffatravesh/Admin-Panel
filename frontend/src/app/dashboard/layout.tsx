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

  const authEndpoint = process.env.NEXT_PUBLIC_API_URL
    ? `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/me`
    : 'http://localhost:8080/api/v1/auth/me';

  try {
    console.log('Checking auth at:', authEndpoint);
    const response = await fetch(authEndpoint, {
      cache: 'no-store',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      },
    });

    console.log('Auth check status:', response.status);
    if (!response.ok) {
      redirect('/login');
    }

    const data = await response.json();
    console.log('Auth check data:', data);
    if (!data.success) {
      redirect('/login');
    }
  } catch (error) {
    console.error('Auth check error:', error);
    redirect('/login');
  }

  return (
    <ClientProviders>
      <DashboardShell>
        {children}
      </DashboardShell>
    </ClientProviders>
  );
}
