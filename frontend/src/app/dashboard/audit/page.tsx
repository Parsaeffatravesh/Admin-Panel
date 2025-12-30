'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { auditApi, AuditLog } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { DataTable } from '@/components/data-table';
import { Badge } from '@/components/ui/badge';

export default function AuditPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['audit-logs', page, search],
    queryFn: () => auditApi.list({ page, per_page: 20, search }),
  });

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
        <div className="mb-6">
          <h3 className="text-lg font-medium">System Audit Trail</h3>
          <p className="text-sm text-gray-500">
            View all system activities and changes
          </p>
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
