'use client';

import { useQuery } from '@tanstack/react-query';
import { dashboardApi } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Shield, Activity, LogIn, TrendingUp } from 'lucide-react';

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
      iconBg: 'bg-primary/10 text-primary',
    },
    {
      title: 'Active Users',
      value: stats?.active_users ?? 0,
      icon: Activity,
      iconBg: 'bg-green-500/10 text-green-600 dark:text-green-400',
    },
    {
      title: 'Total Roles',
      value: stats?.total_roles ?? 0,
      icon: Shield,
      iconBg: 'bg-purple-500/10 text-purple-600 dark:text-purple-400',
    },
    {
      title: 'Recent Logins (24h)',
      value: stats?.recent_logins ?? 0,
      icon: LogIn,
      iconBg: 'bg-orange-500/10 text-orange-600 dark:text-orange-400',
    },
  ];

  return (
    <div className="animate-in fade-in">
      <Header title="Dashboard" />
      <div className="p-6">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {statCards.map((stat) => (
            <Card key={stat.title} className="hover:shadow-md transition-shadow duration-200">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <div className={`rounded-lg p-2 ${stat.iconBg}`}>
                  <stat.icon className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-foreground">
                  {isLoading ? (
                    <div className="h-8 w-16 bg-muted animate-pulse rounded" />
                  ) : (
                    stat.value
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg font-semibold">Users by Status</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="space-y-3">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="h-6 bg-muted animate-pulse rounded" />
                  ))}
                </div>
              ) : (
                <div className="space-y-4">
                  {Object.entries(stats?.users_by_status ?? {}).map(([status, count]) => (
                    <div key={status} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <div className={`h-2 w-2 rounded-full ${
                          status === 'active' ? 'bg-green-500' :
                          status === 'inactive' ? 'bg-muted-foreground' :
                          'bg-destructive'
                        }`} />
                        <span className="capitalize text-foreground">{status}</span>
                      </div>
                      <span className="font-semibold text-foreground">{count as number}</span>
                    </div>
                  ))}
                  {Object.keys(stats?.users_by_status ?? {}).length === 0 && (
                    <p className="text-muted-foreground">No data available</p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-lg font-semibold">Recent Activity</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="space-y-3">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <div key={i} className="h-6 bg-muted animate-pulse rounded" />
                  ))}
                </div>
              ) : (
                <div className="space-y-3">
                  {(stats?.recent_activity ?? []).slice(0, 5).map((activity, i) => (
                    <div key={i} className="flex items-center justify-between text-sm border-b border-border/50 pb-2 last:border-0">
                      <div className="flex items-center gap-2">
                        <TrendingUp className="h-3 w-3 text-muted-foreground" />
                        <span className="text-foreground">
                          <span className="font-medium">{activity.action}</span>
                          <span className="text-muted-foreground"> on </span>
                          <span className="text-muted-foreground">{activity.resource}</span>
                        </span>
                      </div>
                      <span className="text-xs text-muted-foreground">
                        {new Date(activity.created_at).toLocaleTimeString()}
                      </span>
                    </div>
                  ))}
                  {(stats?.recent_activity?.length ?? 0) === 0 && (
                    <p className="text-muted-foreground text-center py-4">No recent activity</p>
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
