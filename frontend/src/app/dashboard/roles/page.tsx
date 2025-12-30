'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { rolesApi, Role } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { DataTable } from '@/components/data-table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Plus, Pencil, Trash2, X } from 'lucide-react';

export default function RolesPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['roles', page, search],
    queryFn: () => rolesApi.list({ page, per_page: 10, search }),
  });

  const createMutation = useMutation({
    mutationFn: (data: { name: string; description: string }) =>
      rolesApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      setIsCreateOpen(false);
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Role> }) =>
      rolesApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      setEditingRole(null);
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => rolesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] });
    },
  });

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
    });
  };

  const handleCreate = () => {
    createMutation.mutate(formData);
  };

  const handleUpdate = () => {
    if (editingRole) {
      updateMutation.mutate({
        id: editingRole.id,
        data: formData,
      });
    }
  };

  const handleEdit = (role: Role) => {
    setEditingRole(role);
    setFormData({
      name: role.name,
      description: role.description,
    });
  };

  const columns = [
    {
      key: 'name',
      header: 'Name',
    },
    {
      key: 'description',
      header: 'Description',
    },
    {
      key: 'is_system',
      header: 'Type',
      render: (role: Role) => (
        <Badge variant={role.is_system ? 'default' : 'secondary'}>
          {role.is_system ? 'System' : 'Custom'}
        </Badge>
      ),
    },
    {
      key: 'created_at',
      header: 'Created',
      render: (role: Role) => new Date(role.created_at).toLocaleDateString(),
    },
  ];

  return (
    <div>
      <Header title="Roles" />
      <div className="p-6">
        <div className="mb-6 flex items-center justify-between">
          <h3 className="text-lg font-medium">Manage Roles</h3>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Role
          </Button>
        </div>

        {(isCreateOpen || editingRole) && (
          <Card className="mb-6">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>{editingRole ? 'Edit Role' : 'Create Role'}</CardTitle>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => {
                  setIsCreateOpen(false);
                  setEditingRole(null);
                  resetForm();
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="text-sm font-medium">Name</label>
                  <Input
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="Role name"
                    disabled={editingRole?.is_system}
                  />
                </div>
                <div>
                  <label className="text-sm font-medium">Description</label>
                  <Input
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    placeholder="Role description"
                  />
                </div>
              </div>
              <div className="mt-4 flex justify-end gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    setIsCreateOpen(false);
                    setEditingRole(null);
                    resetForm();
                  }}
                >
                  Cancel
                </Button>
                <Button
                  onClick={editingRole ? handleUpdate : handleCreate}
                  disabled={createMutation.isPending || updateMutation.isPending}
                >
                  {editingRole ? 'Update' : 'Create'}
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        <DataTable
          columns={columns}
          data={data?.data ?? []}
          total={data?.total ?? 0}
          page={page}
          perPage={10}
          onPageChange={setPage}
          onSearch={setSearch}
          isLoading={isLoading}
          actions={(role) => (
            <div className="flex gap-2">
              <Button variant="ghost" size="icon" onClick={() => handleEdit(role)}>
                <Pencil className="h-4 w-4" />
              </Button>
              {!role.is_system && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => {
                    if (confirm('Are you sure you want to delete this role?')) {
                      deleteMutation.mutate(role.id);
                    }
                  }}
                >
                  <Trash2 className="h-4 w-4 text-red-500" />
                </Button>
              )}
            </div>
          )}
        />
      </div>
    </div>
  );
}
