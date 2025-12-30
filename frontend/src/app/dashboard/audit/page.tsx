'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { auditApi, AuditLog } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { DataTable } from '@/components/data-table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

export default function AuditPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [isExporting, setIsExporting] = useState(false);

  const { data, isLoading } = useQuery({
    queryKey: ['audit-logs', page, search],
    queryFn: () => auditApi.list({ page, per_page: 20, search }),
  });

  const handleExport = async () => {
    setIsExporting(true);
    try {
      const token = localStorage.getItem('access_token');
      const baseUrl = process.env.NEXT_PUBLIC_API_URL || '';
      const searchParams = new URLSearchParams();
      if (search) searchParams.set('search', search);
      
      const response = await fetch(`${baseUrl}/api/v1/audit-logs/export?${searchParams}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      
      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `audit_logs_${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
      }
    } catch (error) {
      console.error('Export failed:', error);
    } finally {
      setIsExporting(false);
    }
  };

  const columns = [
    {
      key: 'action',
      header: 'Action',
      render: (log: AuditLog) => (
        <Badge variant="outline">{log.action}</Badge>
      ),
    },
    {
      key: 'resource',
      header: 'Resource',
    },
    {
      key: 'ip_address',
      header: 'IP Address',
    },
    {
      key: 'user_agent',
      header: 'User Agent',
      render: (log: AuditLog) => (
        <span className="max-w-xs truncate block" title={log.user_agent}>
          {log.user_agent?.substring(0, 50)}...
        </span>
      ),
    },
    {
      key: 'created_at',
      header: 'Date',
      render: (log: AuditLog) => new Date(log.created_at).toLocaleString(),
    },
  ];

  return (
    <div>
      <Header title="Audit Logs" />
      <div className="p-6">
        <div className="mb-6 flex justify-between items-start">
          <div>
            <h3 className="text-lg font-medium">System Audit Trail</h3>
            <p className="text-sm text-gray-500">
              View all system activities and changes
            </p>
          </div>
          <Button 
            onClick={handleExport} 
            disabled={isExporting}
            variant="outline"
          >
            {isExporting ? 'Exporting...' : 'Export CSV'}
          </Button>
        </div>

        <DataTable
          columns={columns}
          data={data?.data ?? []}
          total={data?.total ?? 0}
          page={page}
          perPage={20}
          onPageChange={setPage}
          onSearch={setSearch}
          isLoading={isLoading}
        />
      </div>
    </div>
  );
}
