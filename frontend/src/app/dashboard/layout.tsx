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
    : '/api/v1/auth/me';

  try {
    const response = await fetch(authEndpoint, {
      cache: 'no-store',
      credentials: 'include',
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (!response.ok) {
      cookieStore.set('access_token', '', {
        expires: new Date(0),
        path: '/',
      });
      redirect('/login');
    }
  } catch {
    cookieStore.set('access_token', '', {
      expires: new Date(0),
      path: '/',
    });
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
