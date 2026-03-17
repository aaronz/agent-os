'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Intent } from '@/types';
import { Plus, Filter, Search } from 'lucide-react';
import Link from 'next/link';

const statusColors: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-800',
  open: 'bg-blue-100 text-blue-800',
  planning: 'bg-purple-100 text-purple-800',
  executing: 'bg-yellow-100 text-yellow-800',
  paused: 'bg-orange-100 text-orange-800',
  completed: 'bg-green-100 text-green-800',
  failed: 'bg-red-100 text-red-800',
  cancelled: 'bg-gray-100 text-gray-600',
};

const priorityColors: Record<string, string> = {
  high: 'text-red-600',
  medium: 'text-yellow-600',
  low: 'text-green-600',
};

export default function IntentsPage() {
  const { data: intents, isLoading } = useQuery<Intent[]>({
    queryKey: ['intents'],
    queryFn: () => api.get('/api/v1/intents'),
    refetchInterval: 30000,
  });

  return (
    <div className="p-8">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Intents</h1>
          <p className="text-gray-500">Manage and track intents</p>
        </div>
        <Link
          href="/intents/new"
          className="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <Plus className="h-4 w-4" />
          New Intent
        </Link>
      </div>

      <div className="mb-6 flex gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search intents..."
            className="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 text-sm focus:border-blue-500 focus:outline-none"
          />
        </div>
        <button className="flex items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium hover:bg-gray-50">
          <Filter className="h-4 w-4" />
          Filter
        </button>
      </div>

      <div className="rounded-lg border bg-white shadow-sm">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : intents?.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No intents found</div>
        ) : (
          <table className="w-full">
            <thead className="border-b bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Title</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Priority</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Created</th>
                <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Updated</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {intents?.map((intent) => (
                <tr key={intent.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <Link href={`/intents/${intent.id}`} className="font-medium text-blue-600 hover:underline">
                      {intent.title}
                    </Link>
                    <p className="text-sm text-gray-500">{intent.description.slice(0, 60)}...</p>
                  </td>
                  <td className={`px-6 py-4 font-medium ${priorityColors[intent.priority]}`}>
                    {intent.priority}
                  </td>
                  <td className="px-6 py-4">
                    <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColors[intent.status]}`}>
                      {intent.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {formatRelativeTime(intent.created_at)}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {formatRelativeTime(intent.updated_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
