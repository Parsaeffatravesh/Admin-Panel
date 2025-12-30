'use client';

import { useQuery } from '@tanstack/react-query';
import { dashboardApi } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Shield, Activity, LogIn } from 'lucide-react';

export default function DashboardPage() {
  const { data: stats, isLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: dashboardApi.getStats,
  });

  const statCards = [
    {
      title: 'Total Users',
      value: stats?.total_users ?? 0,
      icon: Users,
      color: 'bg-blue-500',
    },
    {
      title: 'Active Users',
      value: stats?.active_users ?? 0,
      icon: Activity,
      color: 'bg-green-500',
    },
    {
      title: 'Total Roles',
      value: stats?.total_roles ?? 0,
      icon: Shield,
      color: 'bg-purple-500',
    },
    {
      title: 'Recent Logins (24h)',
      value: stats?.recent_logins ?? 0,
      icon: LogIn,
      color: 'bg-orange-500',
    },
  ];

  return (
    <div>
      <Header title="Dashboard" />
      <div className="p-6">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {statCards.map((stat) => (
            <Card key={stat.title}>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-gray-500">
                  {stat.title}
                </CardTitle>
                <div className={`rounded-full p-2 ${stat.color}`}>
                  <stat.icon className="h-4 w-4 text-white" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {isLoading ? '...' : stat.value}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Users by Status</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <p>Loading...</p>
              ) : (
                <div className="space-y-4">
                  {Object.entries(stats?.users_by_status ?? {}).map(([status, count]) => (
                    <div key={status} className="flex items-center justify-between">
                      <span className="capitalize">{status}</span>
                      <span className="font-medium">{count}</span>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <p>Loading...</p>
              ) : (
                <div className="space-y-4">
                  {(stats?.recent_activity ?? []).slice(0, 5).map((activity, i) => (
                    <div key={i} className="flex items-center justify-between text-sm">
                      <span>
                        <span className="font-medium">{activity.action}</span> on{' '}
                        <span className="text-gray-500">{activity.resource}</span>
                      </span>
                      <span className="text-gray-400">
                        {new Date(activity.created_at).toLocaleTimeString()}
                      </span>
                    </div>
                  ))}
                  {(stats?.recent_activity?.length ?? 0) === 0 && (
                    <p className="text-gray-500">No recent activity</p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
