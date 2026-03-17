'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Agent } from '@/types';
import { Search, Filter, Shield, UserCircle } from 'lucide-react';

const statusColors: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  idle: 'bg-gray-100 text-gray-800',
  busy: 'bg-yellow-100 text-yellow-800',
  disabled: 'bg-red-100 text-red-800',
  terminated: 'bg-gray-100 text-gray-600',
};

export default function AgentsPage() {
  const { data: agents, isLoading } = useQuery<Agent[]>({
    queryKey: ['agents'],
    queryFn: () => api.get('/api/v1/agents'),
    refetchInterval: 30000,
  });

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Agents</h1>
        <p className="text-gray-500">Manage and monitor agents</p>
      </div>

      <div className="mb-6 flex gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search agents..."
            className="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 text-sm focus:border-blue-500 focus:outline-none"
          />
        </div>
        <button className="flex items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium hover:bg-gray-50">
          <Filter className="h-4 w-4" />
          Filter
        </button>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {isLoading ? (
          Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="animate-pulse rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-4">
                <div className="h-12 w-12 rounded-full bg-gray-200" />
                <div className="flex-1">
                  <div className="h-4 w-24 rounded bg-gray-200" />
                  <div className="mt-2 h-3 w-16 rounded bg-gray-200" />
                </div>
              </div>
            </div>
          ))
        ) : agents?.length === 0 ? (
          <div className="col-span-full p-8 text-center text-gray-500">No agents found</div>
        ) : (
          agents?.map((agent) => (
            <div key={agent.id} className="rounded-lg border bg-white p-6 shadow-sm hover:shadow-md">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
                    <UserCircle className="h-8 w-8 text-blue-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold">{agent.name}</h3>
                    <p className="text-sm text-gray-500">{agent.role}</p>
                  </div>
                </div>
                <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColors[agent.status]}`}>
                  {agent.status}
                </span>
              </div>

              <div className="mt-4 flex items-center justify-between border-t pt-4">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4 text-gray-400" />
                  <span className="text-sm font-medium">Reputation</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="h-2 w-24 rounded-full bg-gray-200">
                    <div
                      className="h-2 rounded-full bg-green-500"
                      style={{ width: `${agent.reputation}%` }}
                    />
                  </div>
                  <span className="text-sm font-medium">{agent.reputation}</span>
                </div>
              </div>

              <div className="mt-3 flex flex-wrap gap-1">
                {agent.capabilities.slice(0, 3).map((cap) => (
                  <span key={cap} className="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600">
                    {cap}
                  </span>
                ))}
                {agent.capabilities.length > 3 && (
                  <span className="text-xs text-gray-500">+{agent.capabilities.length - 3}</span>
                )}
              </div>

              <p className="mt-3 text-xs text-gray-400">
                Last active: {formatRelativeTime(agent.last_active_at)}
              </p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
