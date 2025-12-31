'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usersApi, User } from '@/lib/api';
import { Header } from '@/components/layout/header';
import { DataTable } from '@/components/data-table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Plus, Pencil, Trash2, X } from 'lucide-react';

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [formData, setFormData] = useState<{
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    status: 'active' | 'inactive' | 'suspended';
  }>({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    status: 'active',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['users', page, search],
    queryFn: () => usersApi.list({ page, per_page: 10, search }),
  });

  const createMutation = useMutation({
    mutationFn: (data: { email: string; password: string; first_name: string; last_name: string }) =>
      usersApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setIsCreateOpen(false);
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<User> }) =>
      usersApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setEditingUser(null);
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => usersApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  const resetForm = () => {
    setFormData({
      email: '',
      password: '',
      first_name: '',
      last_name: '',
      status: 'active',
    });
  };

  const handleCreate = () => {
    createMutation.mutate(formData);
  };

  const handleUpdate = () => {
    if (editingUser) {
      updateMutation.mutate({
        id: editingUser.id,
        data: {
          email: formData.email,
          first_name: formData.first_name,
          last_name: formData.last_name,
          status: formData.status,
        },
      });
    }
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setFormData({
      email: user.email,
      password: '',
      first_name: user.first_name,
      last_name: user.last_name,
      status: user.status,
    });
  };

  const columns = [
    {
      key: 'email',
      header: 'Email',
    },
    {
      key: 'first_name',
      header: 'First Name',
    },
    {
      key: 'last_name',
      header: 'Last Name',
    },
    {
      key: 'status',
      header: 'Status',
      render: (user: User) => (
        <Badge
          variant={
            user.status === 'active'
              ? 'success'
              : user.status === 'inactive'
              ? 'secondary'
              : 'destructive'
          }
        >
          {user.status}
        </Badge>
      ),
    },
    {
      key: 'created_at',
      header: 'Created',
      render: (user: User) => new Date(user.created_at).toLocaleDateString(),
    },
  ];

  return (
    <div>
      <Header title="Users" />
      <div className="p-6">
        <div className="mb-6 flex items-center justify-between">
          <h3 className="text-lg font-medium">Manage Users</h3>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add User
          </Button>
        </div>

        {(isCreateOpen || editingUser) && (
          <Card className="mb-6">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>{editingUser ? 'Edit User' : 'Create User'}</CardTitle>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => {
                  setIsCreateOpen(false);
                  setEditingUser(null);
                  resetForm();
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="text-sm font-medium">Email</label>
                  <Input
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    placeholder="user@example.com"
                  />
                </div>
                {!editingUser && (
                  <div>
                    <label className="text-sm font-medium">Password</label>
                    <Input
                      type="password"
                      value={formData.password}
                      onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                      placeholder="Password"
                    />
                  </div>
                )}
                <div>
                  <label className="text-sm font-medium">First Name</label>
                  <Input
                    value={formData.first_name}
                    onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                    placeholder="John"
                  />
                </div>
                <div>
                  <label className="text-sm font-medium">Last Name</label>
                  <Input
                    value={formData.last_name}
                    onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                    placeholder="Doe"
                  />
                </div>
                {editingUser && (
                  <div>
                    <label className="text-sm font-medium">Status</label>
                    <select
                      value={formData.status}
                      onChange={(e) => setFormData({ ...formData, status: e.target.value as 'active' | 'inactive' | 'suspended' })}
                      className="w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm"
                    >
                      <option value="active">Active</option>
                      <option value="inactive">Inactive</option>
                      <option value="suspended">Suspended</option>
                    </select>
                  </div>
                )}
              </div>
              <div className="mt-4 flex justify-end gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    setIsCreateOpen(false);
                    setEditingUser(null);
                    resetForm();
                  }}
                >
                  Cancel
                </Button>
                <Button
                  onClick={editingUser ? handleUpdate : handleCreate}
                  disabled={createMutation.isPending || updateMutation.isPending}
                >
                  {editingUser ? 'Update' : 'Create'}
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
          actions={(user) => (
            <div className="flex gap-2">
              <Button variant="ghost" size="icon" onClick={() => handleEdit(user)}>
                <Pencil className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => {
                  if (confirm('Are you sure you want to delete this user?')) {
                    deleteMutation.mutate(user.id);
                  }
                }}
              >
                <Trash2 className="h-4 w-4 text-red-500" />
              </Button>
            </div>
          )}
        />
      </div>
    </div>
  );
}
