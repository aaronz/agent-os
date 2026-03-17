'use client';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Memory } from '@/types';
import { Search, BookOpen, Lightbulb, AlertTriangle, CheckCircle, Filter, Brain } from 'lucide-react';

const typeConfig: Record<string, { color: string; bg: string; icon: React.ReactNode }> = {
  knowledge: { color: 'text-blue-700', bg: 'bg-blue-100', icon: <BookOpen className="h-4 w-4" /> },
  project: { color: 'text-purple-700', bg: 'bg-purple-100', icon: <Brain className="h-4 w-4" /> },
  task: { color: 'text-green-700', bg: 'bg-green-100', icon: <CheckCircle className="h-4 w-4" /> },
  failure: { color: 'text-red-700', bg: 'bg-red-100', icon: <AlertTriangle className="h-4 w-4" /> },
  best_practice: { color: 'text-yellow-700', bg: 'bg-yellow-100', icon: <Lightbulb className="h-4 w-4" /> },
  review: { color: 'text-indigo-700', bg: 'bg-indigo-100', icon: <CheckCircle className="h-4 w-4" /> },
};

const validityConfig: Record<string, { color: string; bg: string }> = {
  valid: { color: 'text-green-700', bg: 'bg-green-100' },
  invalid: { color: 'text-red-700', bg: 'bg-red-100' },
};

export default function MemoryPage() {
  const [selectedMemory, setSelectedMemory] = useState<Memory | null>(null);
  const [filter, setFilter] = useState<string>('all');
  const [search, setSearch] = useState('');

  const { data: memories, isLoading } = useQuery<Memory[]>({
    queryKey: ['memories', filter],
    queryFn: () => api.get(filter !== 'all' ? `/api/v1/memory?type=${filter}` : '/api/v1/memory'),
    refetchInterval: 60000,
  });

  const filteredMemories = memories?.filter(memory => {
    const matchesSearch = memory.title.toLowerCase().includes(search.toLowerCase()) ||
      memory.content.toLowerCase().includes(search.toLowerCase());
    return matchesSearch;
  });

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Memory</h1>
        <p className="text-gray-500">Agent knowledge base and learned patterns</p>
      </div>

      <div className="mb-6 flex gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search memories..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 text-sm focus:border-blue-500 focus:outline-none"
          />
        </div>
        <select
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-blue-500 focus:outline-none"
        >
          <option value="all">All Types</option>
          <option value="knowledge">Knowledge</option>
          <option value="project">Project</option>
          <option value="task">Task</option>
          <option value="failure">Failure</option>
          <option value="best_practice">Best Practice</option>
          <option value="review">Review</option>
        </select>
      </div>

      <div className="space-y-4">
        {isLoading ? (
          Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="animate-pulse rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-4">
                <div className="h-10 w-10 rounded bg-gray-200" />
                <div className="flex-1">
                  <div className="h-4 w-48 rounded bg-gray-200" />
                  <div className="mt-2 h-3 w-full rounded bg-gray-200" />
                </div>
              </div>
            </div>
          ))
        ) : filteredMemories?.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No memories found</div>
        ) : (
          filteredMemories?.map((memory) => {
            const type = typeConfig[memory.type] || typeConfig.knowledge;
            const validity = validityConfig[memory.validity] || validityConfig.valid;
            return (
              <div
                key={memory.id}
                onClick={() => setSelectedMemory(memory)}
                className="cursor-pointer rounded-lg border bg-white p-6 shadow-sm hover:shadow-md"
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className={`flex h-10 w-10 items-center justify-center rounded-lg ${type.bg} ${type.color}`}>
                      {type.icon}
                    </div>
                    <div>
                      <h3 className="font-medium text-gray-900">{memory.title}</h3>
                      <div className="mt-1 flex items-center gap-2">
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${type.bg} ${type.color}`}>
                          {memory.type.replace('_', ' ')}
                        </span>
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${validity.bg} ${validity.color}`}>
                          {memory.validity}
                        </span>
                      </div>
                    </div>
                  </div>
                  <span className="text-sm text-gray-400">
                    {formatRelativeTime(memory.created_at)}
                  </span>
                </div>

                <p className="mt-3 text-sm text-gray-600 line-clamp-2">{memory.content}</p>

                <div className="mt-4 flex items-center justify-between border-t pt-4">
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-gray-500">Source:</span>
                    <span className="text-xs font-medium text-gray-700">{memory.source}</span>
                  </div>
                  <span className="text-xs text-gray-400">
                    Last retrieved: {memory.last_retrieved_at ? formatRelativeTime(memory.last_retrieved_at) : 'Never'}
                  </span>
                </div>
              </div>
            );
          })
        )}
      </div>

      {selectedMemory && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="mx-4 max-h-[90vh] w-full max-w-2xl overflow-auto rounded-lg bg-white p-6 shadow-xl">
            <div className="mb-6 flex items-start justify-between">
              <div>
                <h2 className="text-xl font-bold text-gray-900">{selectedMemory.title}</h2>
                <p className="text-sm text-gray-500 mt-1">ID: {selectedMemory.id}</p>
              </div>
              <button onClick={() => setSelectedMemory(null)} className="text-gray-400 hover:text-gray-600">
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="text-sm font-medium text-gray-700">Type</label>
                  <div className="mt-1">
                    {(() => {
                      const t = typeConfig[selectedMemory.type] || typeConfig.knowledge;
                      return (
                        <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${t.bg} ${t.color}`}>
                          {t.icon}
                          {selectedMemory.type.replace('_', ' ')}
                        </span>
                      );
                    })()}
                  </div>
                </div>
                <div className="flex-1">
                  <label className="text-sm font-medium text-gray-700">Validity</label>
                  <div className="mt-1">
                    {(() => {
                      const v = validityConfig[selectedMemory.validity] || validityConfig.valid;
                      return (
                        <span className={`rounded-full px-2 py-1 text-xs font-medium ${v.bg} ${v.color}`}>
                          {selectedMemory.validity}
                        </span>
                      );
                    })()}
                  </div>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Content</label>
                <p className="mt-1 whitespace-pre-wrap text-sm text-gray-600">{selectedMemory.content}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Related Entities</label>
                <p className="mt-1 text-sm text-gray-600">{selectedMemory.related_entities || 'None'}</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-700">Source</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedMemory.source}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-700">Created At</label>
                  <p className="mt-1 text-sm text-gray-600">{new Date(selectedMemory.created_at).toLocaleString()}</p>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Last Retrieved</label>
                <p className="mt-1 text-sm text-gray-600">
                  {selectedMemory.last_retrieved_at ? formatRelativeTime(selectedMemory.last_retrieved_at) : 'Never'}
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
