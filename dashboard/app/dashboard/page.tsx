'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import {
  Users,
  Target,
  CheckSquare,
  AlertCircle,
  Activity,
  TrendingUp,
} from 'lucide-react';

interface Metrics {
  requests_total: number;
  tasks_created: number;
  tasks_completed: number;
  tasks_failed: number;
  intents_created: number;
  intents_completed: number;
  agents_active: number;
  error_rate: number;
}

export default function DashboardPage() {
  const { data: metrics, isLoading: metricsLoading } = useQuery<Metrics>({
    queryKey: ['metrics'],
    queryFn: () => api.get('/api/v1/metrics'),
    refetchInterval: 30000,
  });

  const { data: agents } = useQuery({
    queryKey: ['agents'],
    queryFn: () => api.get('/api/v1/agents'),
  });

  const { data: intents } = useQuery({
    queryKey: ['intents'],
    queryFn: () => api.get('/api/v1/intents'),
  });

  const statCards = [
    {
      name: 'Active Agents',
      value: metrics?.agents_active || 0,
      icon: Users,
      color: 'text-blue-600',
    },
    {
      name: 'Open Intents',
      value: intents?.length || 0,
      icon: Target,
      color: 'text-purple-600',
    },
    {
      name: 'Tasks Completed',
      value: metrics?.tasks_completed || 0,
      icon: CheckSquare,
      color: 'text-green-600',
    },
    {
      name: 'Error Rate',
      value: `${metrics?.error_rate.toFixed(1) || 0}%`,
      icon: AlertCircle,
      color: 'text-red-600',
    },
  ];

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500">System overview and monitoring</p>
      </div>

      <div className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => (
          <div
            key={stat.name}
            className="rounded-lg border bg-white p-6 shadow-sm"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-500">{stat.name}</p>
                {metricsLoading ? (
                  <div className="h-8 w-16 animate-pulse rounded bg-gray-200" />
                ) : (
                  <p className="text-2xl font-semibold">{stat.value}</p>
                )}
              </div>
              <stat.icon className={`h-8 w-8 ${stat.color}`} />
            </div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold">Recent Activity</h2>
          <div className="space-y-4">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-100">
                  <Activity className="h-5 w-5 text-blue-600" />
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium">Task completed</p>
                  <p className="text-xs text-gray-500">{formatRelativeTime(new Date().toISOString())}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-lg border bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold">System Health</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">API Status</span>
              <span className="flex items-center gap-2 text-sm font-medium text-green-600">
                <span className="h-2 w-2 rounded-full bg-green-600" />
                Operational
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Database</span>
              <span className="flex items-center gap-2 text-sm font-medium text-green-600">
                <span className="h-2 w-2 rounded-full bg-green-600" />
                Connected
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Message Queue</span>
              <span className="flex items-center gap-2 text-sm font-medium text-green-600">
                <span className="h-2 w-2 rounded-full bg-green-600" />
                Active
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Vector Store</span>
              <span className="flex items-center gap-2 text-sm font-medium text-yellow-600">
                <span className="h-2 w-2 rounded-full bg-yellow-600" />
                Degraded
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
