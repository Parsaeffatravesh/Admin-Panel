'use client';

import { useQuery } from '@tanstack/react-query';
import { dashboardApi } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Shield, Activity, LogIn, TrendingUp, ArrowUpRight } from 'lucide-react';
import { useI18n } from '@/lib/i18n';

export default function DashboardPage() {
  const { data: stats, isLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: dashboardApi.getStats,
  });

  const { t, language } = useI18n();

  const statCards = [
    {
      titleKey: 'dashboard.totalUsers',
      value: stats?.total_users ?? 0,
      icon: Users,
      color: 'text-blue-600 dark:text-blue-400',
      bgColor: 'bg-blue-500/10',
      trend: '+12%',
    },
    {
      titleKey: 'dashboard.activeUsers',
      value: stats?.active_users ?? 0,
      icon: Activity,
      color: 'text-emerald-600 dark:text-emerald-400',
      bgColor: 'bg-emerald-500/10',
      trend: '+5%',
    },
    {
      titleKey: 'dashboard.totalRoles',
      value: stats?.total_roles ?? 0,
      icon: Shield,
      color: 'text-violet-600 dark:text-violet-400',
      bgColor: 'bg-violet-500/10',
      trend: 'Static',
    },
    {
      titleKey: 'dashboard.auditLogs',
      value: stats?.recent_logins ?? 0,
      icon: LogIn,
      color: 'text-orange-600 dark:text-orange-400',
      bgColor: 'bg-orange-500/10',
      trend: 'Recent',
    },
  ];

  const statusLabels: Record<string, { en: string; fa: string }> = {
    active: { en: 'Active', fa: 'فعال' },
    inactive: { en: 'Inactive', fa: 'غیرفعال' },
    suspended: { en: 'Suspended', fa: 'معلق' },
  };

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <Header title={t('dashboard.title')} />
      
      <div className="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto space-y-8">
        <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
          {statCards.map((stat, idx) => (
            <Card 
              key={stat.titleKey} 
              className="group relative overflow-hidden border-border/50 bg-card/50 backdrop-blur-sm hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5 shadow-sm"
              style={{ animationDelay: `${idx * 100}ms` }}
            >
              <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                <CardTitle className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">
                  {t(stat.titleKey)}
                </CardTitle>
                <div className={`rounded-xl p-2.5 ${stat.bgColor} ${stat.color} transition-transform duration-300 group-hover:scale-110`}>
                  <stat.icon className="h-5 w-5" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="flex items-baseline space-x-2 rtl:space-x-reverse">
                  <div className="text-3xl font-bold tracking-tight text-foreground">
                    {isLoading ? (
                      <div className="h-9 w-20 bg-muted/50 animate-pulse rounded-lg" />
                    ) : (
                      stat.value
                    )}
                  </div>
                  {!isLoading && (
                    <span className="text-xs font-medium text-muted-foreground/80 flex items-center">
                      <ArrowUpRight className="h-3 w-3 mr-1 rtl:ml-1" />
                      {stat.trend}
                    </span>
                  )}
                </div>
              </CardContent>
              <div className="absolute inset-x-0 bottom-0 h-1 bg-gradient-to-r from-transparent via-primary/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
            </Card>
          ))}
        </div>

        <div className="grid gap-6 grid-cols-1 lg:grid-cols-2">
          <Card className="border-border/50 bg-card/50 backdrop-blur-sm shadow-sm overflow-hidden">
            <CardHeader className="border-b border-border/50 bg-muted/30">
              <CardTitle className="text-lg font-semibold flex items-center gap-2">
                <Activity className="h-5 w-5 text-primary" />
                {language === 'fa' ? 'وضعیت کاربران' : 'User Status'}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              {isLoading ? (
                <div className="space-y-4">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="h-10 bg-muted/50 animate-pulse rounded-lg" />
                  ))}
                </div>
              ) : (
                <div className="space-y-6">
                  {Object.entries(stats?.users_by_status ?? {}).map(([status, count]) => (
                    <div key={status} className="group flex items-center justify-between p-3 rounded-xl hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-3">
                        <div className={`h-3 w-3 rounded-full shadow-[0_0_8px_rgba(0,0,0,0.1)] ${
                          status === 'active' ? 'bg-emerald-500 shadow-emerald-500/20' :
                          status === 'inactive' ? 'bg-slate-400 shadow-slate-400/20' :
                          'bg-rose-500 shadow-rose-500/20'
                        }`} />
                        <span className="font-medium text-foreground/90">
                          {statusLabels[status]?.[language] || status}
                        </span>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className="text-lg font-bold text-foreground">{count as number}</span>
                        <div className="h-1.5 w-24 bg-muted rounded-full overflow-hidden">
                          <div 
                            className={`h-full rounded-full transition-all duration-1000 ${
                              status === 'active' ? 'bg-emerald-500' :
                              status === 'inactive' ? 'bg-slate-400' :
                              'bg-rose-500'
                            }`}
                            style={{ width: `${Math.min(((count as number) / (stats?.total_users || 1)) * 100, 100)}%` }}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                  {Object.keys(stats?.users_by_status ?? {}).length === 0 && (
                    <div className="text-center py-8 text-muted-foreground italic">
                      {language === 'fa' ? 'داده‌ای یافت نشد' : 'No data found'}
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          <Card className="border-border/50 bg-card/50 backdrop-blur-sm shadow-sm overflow-hidden">
            <CardHeader className="border-b border-border/50 bg-muted/30">
              <CardTitle className="text-lg font-semibold flex items-center gap-2">
                <TrendingUp className="h-5 w-5 text-primary" />
                {language === 'fa' ? 'آخرین فعالیت‌ها' : 'Recent Activity'}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              {isLoading ? (
                <div className="space-y-4">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <div key={i} className="h-12 bg-muted/50 animate-pulse rounded-lg" />
                  ))}
                </div>
              ) : (
                <div className="space-y-1">
                  {(stats?.recent_activity ?? []).slice(0, 6).map((activity, i) => (
                    <div 
                      key={i} 
                      className="flex items-center justify-between p-3 rounded-xl hover:bg-muted/50 transition-all duration-200 border-l-2 border-transparent hover:border-primary/50 group"
                    >
                      <div className="flex items-center gap-4">
                        <div className="h-8 w-8 rounded-full bg-primary/5 flex items-center justify-center text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
                          <Activity className="h-4 w-4" />
                        </div>
                        <div className="flex flex-col">
                          <span className="text-sm font-semibold text-foreground">
                            {activity.action}
                          </span>
                          <span className="text-xs text-muted-foreground">
                            {activity.resource}
                          </span>
                        </div>
                      </div>
                      <span className="text-[10px] uppercase tracking-wider font-bold text-muted-foreground/60 bg-muted px-2 py-1 rounded-md">
                        {new Date(activity.created_at).toLocaleTimeString(language === 'fa' ? 'fa-IR' : 'en-US', { hour: '2-digit', minute: '2-digit' })}
                      </span>
                    </div>
                  ))}
                  {(stats?.recent_activity?.length ?? 0) === 0 && (
                    <div className="text-center py-12 text-muted-foreground flex flex-col items-center gap-3">
                      <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
                        <Activity className="h-6 w-6 opacity-20" />
                      </div>
                      {language === 'fa' ? 'هنوز فعالیتی ثبت نشده است' : 'No activity recorded yet'}
                    </div>
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
