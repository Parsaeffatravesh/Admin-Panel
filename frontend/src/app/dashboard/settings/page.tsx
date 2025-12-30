'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Header } from '@/components/layout/header';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { featureFlagsApi, FeatureFlag } from '@/lib/api';
import { cn } from '@/lib/utils';

function Switch({ checked, onCheckedChange, disabled }: { 
  checked: boolean; 
  onCheckedChange: (checked: boolean) => void;
  disabled?: boolean;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={() => onCheckedChange(!checked)}
      className={cn(
        "relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
        checked ? "bg-blue-600" : "bg-gray-200",
        disabled && "opacity-50 cursor-not-allowed"
      )}
    >
      <span
        className={cn(
          "inline-block h-4 w-4 transform rounded-full bg-white transition-transform",
          checked ? "translate-x-6" : "translate-x-1"
        )}
      />
    </button>
  );
}

function CreateFlagDialog({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
  const [key, setKey] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const createMutation = useMutation({
    mutationFn: () => featureFlagsApi.create({ key, name, description, enabled: false }),
    onSuccess: () => {
      onSuccess();
      onClose();
    },
  });

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 className="text-lg font-semibold mb-4">Create Feature Flag</h3>
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium">Key</label>
            <Input 
              value={key} 
              onChange={(e) => setKey(e.target.value)} 
              placeholder="e.g., enable_2fa"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Name</label>
            <Input 
              value={name} 
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., Two-Factor Authentication"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Description</label>
            <Input 
              value={description} 
              onChange={(e) => setDescription(e.target.value)}
              placeholder="e.g., Require 2FA for all users"
            />
          </div>
        </div>
        <div className="flex justify-end gap-2 mt-6">
          <Button variant="outline" onClick={onClose}>Cancel</Button>
          <Button 
            onClick={() => createMutation.mutate()}
            disabled={!key || !name || createMutation.isPending}
          >
            {createMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function SettingsPage() {
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const queryClient = useQueryClient();

  const { data: flagsData, isLoading } = useQuery({
    queryKey: ['feature-flags'],
    queryFn: () => featureFlagsApi.list({ per_page: 100 }),
  });

  const toggleMutation = useMutation({
    mutationFn: (id: string) => featureFlagsApi.toggle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feature-flags'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => featureFlagsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feature-flags'] });
    },
  });

  const flags = flagsData?.data ?? [];

  return (
    <div>
      <Header title="Settings" />
      <div className="p-6">
        <div className="grid gap-6">
          <Card>
            <CardHeader>
              <CardTitle>General Settings</CardTitle>
              <CardDescription>
                Configure general application settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="text-sm font-medium">Organization Name</label>
                <Input defaultValue="Default Organization" />
              </div>
              <div>
                <label className="text-sm font-medium">Contact Email</label>
                <Input type="email" defaultValue="admin@example.com" />
              </div>
              <Button>Save Changes</Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Feature Flags</CardTitle>
                <CardDescription>
                  Enable or disable features across the application
                </CardDescription>
              </div>
              <Button onClick={() => setShowCreateDialog(true)}>Add Flag</Button>
            </CardHeader>
            <CardContent className="space-y-4">
              {isLoading ? (
                <div className="text-center py-4 text-gray-500">Loading flags...</div>
              ) : flags.length === 0 ? (
                <div className="text-center py-4 text-gray-500">
                  No feature flags configured. Click "Add Flag" to create one.
                </div>
              ) : (
                flags.map((flag: FeatureFlag) => (
                  <div key={flag.id} className="flex items-center justify-between py-2 border-b last:border-0">
                    <div className="flex-1">
                      <p className="font-medium">{flag.name}</p>
                      <p className="text-sm text-gray-500">{flag.description}</p>
                      <p className="text-xs text-gray-400 mt-1">Key: {flag.key}</p>
                    </div>
                    <div className="flex items-center gap-3">
                      <Switch
                        checked={flag.enabled}
                        onCheckedChange={() => toggleMutation.mutate(flag.id)}
                        disabled={toggleMutation.isPending}
                      />
                      <Button
                        variant="ghost"
                        size="sm"
                        className="text-red-600 hover:text-red-700 hover:bg-red-50"
                        onClick={() => {
                          if (confirm(`Delete feature flag "${flag.name}"?`)) {
                            deleteMutation.mutate(flag.id);
                          }
                        }}
                      >
                        Delete
                      </Button>
                    </div>
                  </div>
                ))
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Security</CardTitle>
              <CardDescription>
                Configure security settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="text-sm font-medium">Session Timeout (minutes)</label>
                <Input type="number" defaultValue="30" />
              </div>
              <div>
                <label className="text-sm font-medium">Max Login Attempts</label>
                <Input type="number" defaultValue="5" />
              </div>
              <Button>Save Changes</Button>
            </CardContent>
          </Card>
        </div>
      </div>

      {showCreateDialog && (
        <CreateFlagDialog
          onClose={() => setShowCreateDialog(false)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ['feature-flags'] })}
        />
      )}
    </div>
  );
}
