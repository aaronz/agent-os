'use client';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Artifact } from '@/types';
import { Search, FileText, Image, Code, Download, Eye, Trash2, ExternalLink } from 'lucide-react';

const typeIcons: Record<string, React.ReactNode> = {
  text: <FileText className="h-5 w-5" />,
  image: <Image className="h-5 w-5" />,
  code: <Code className="h-5 w-5" />,
  default: <FileText className="h-5 w-5" />,
};

const statusConfig: Record<string, { color: string; bg: string }> = {
  pending_review: { color: 'text-yellow-700', bg: 'bg-yellow-100' },
  approved: { color: 'text-green-700', bg: 'bg-green-100' },
  rejected: { color: 'text-red-700', bg: 'bg-red-100' },
  deprecated: { color: 'text-gray-700', bg: 'bg-gray-100' },
};

export default function ArtifactsPage() {
  const [selectedArtifact, setSelectedArtifact] = useState<Artifact | null>(null);
  const [filter, setFilter] = useState<string>('all');
  const [search, setSearch] = useState('');

  const { data: artifacts, isLoading } = useQuery<Artifact[]>({
    queryKey: ['artifacts', filter],
    queryFn: () => api.get(filter !== 'all' ? `/api/v1/artifacts?status=${filter}` : '/api/v1/artifacts'),
    refetchInterval: 30000,
  });

  const filteredArtifacts = artifacts?.filter(artifact => {
    const matchesSearch = artifact.title.toLowerCase().includes(search.toLowerCase()) ||
      artifact.description.toLowerCase().includes(search.toLowerCase());
    const matchesFilter = filter === 'all' || artifact.status === filter;
    return matchesSearch && matchesFilter;
  });

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Artifacts</h1>
        <p className="text-gray-500">Browse and manage task outputs and deliverables</p>
      </div>

      <div className="mb-6 flex gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search artifacts..."
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
          <option value="all">All Status</option>
          <option value="pending_review">Pending Review</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
          <option value="deprecated">Deprecated</option>
        </select>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {isLoading ? (
          Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="animate-pulse rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-4">
                <div className="h-10 w-10 rounded bg-gray-200" />
                <div className="flex-1">
                  <div className="h-4 w-32 rounded bg-gray-200" />
                  <div className="mt-2 h-3 w-24 rounded bg-gray-200" />
                </div>
              </div>
            </div>
          ))
        ) : filteredArtifacts?.length === 0 ? (
          <div className="col-span-full p-8 text-center text-gray-500">No artifacts found</div>
        ) : (
          filteredArtifacts?.map((artifact) => {
            const icon = typeIcons[artifact.type] || typeIcons.default;
            const status = statusConfig[artifact.status] || statusConfig.pending_review;
            return (
              <div
                key={artifact.id}
                onClick={() => setSelectedArtifact(artifact)}
                className="cursor-pointer rounded-lg border bg-white p-6 shadow-sm hover:shadow-md"
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                      {icon}
                    </div>
                    <div>
                      <h3 className="font-medium text-gray-900">{artifact.title}</h3>
                      <p className="text-xs text-gray-500">{artifact.type}</p>
                    </div>
                  </div>
                  <span className={`rounded-full px-2 py-1 text-xs font-medium ${status.bg} ${status.color}`}>
                    {artifact.status.replace('_', ' ')}
                  </span>
                </div>

                <p className="mt-3 text-sm text-gray-600 line-clamp-2">{artifact.description}</p>

                <div className="mt-4 flex items-center justify-between border-t pt-4">
                  <span className="text-xs text-gray-400">
                    v{artifact.version} · {formatRelativeTime(artifact.created_at)}
                  </span>
                  <div className="flex gap-2">
                    <button className="rounded p-1 hover:bg-gray-100" title="View">
                      <Eye className="h-4 w-4 text-gray-500" />
                    </button>
                    <button className="rounded p-1 hover:bg-gray-100" title="Download">
                      <Download className="h-4 w-4 text-gray-500" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>

      {selectedArtifact && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="mx-4 max-h-[90vh] w-full max-w-2xl overflow-auto rounded-lg bg-white p-6 shadow-xl">
            <div className="mb-6 flex items-start justify-between">
              <div>
                <h2 className="text-xl font-bold text-gray-900">{selectedArtifact.title}</h2>
                <p className="text-sm text-gray-500 mt-1">ID: {selectedArtifact.id}</p>
              </div>
              <button onClick={() => setSelectedArtifact(null)} className="text-gray-400 hover:text-gray-600">
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="text-sm font-medium text-gray-700">Type</label>
                  <p className="mt-1 text-sm text-gray-600 capitalize">{selectedArtifact.type}</p>
                </div>
                <div className="flex-1">
                  <label className="text-sm font-medium text-gray-700">Version</label>
                  <p className="mt-1 text-sm text-gray-600">v{selectedArtifact.version}</p>
                </div>
                <div className="flex-1">
                  <label className="text-sm font-medium text-gray-700">Status</label>
                  <div className="mt-1">
                    {(() => {
                      const s = statusConfig[selectedArtifact.status] || statusConfig.pending_review;
                      return (
                        <span className={`rounded-full px-2 py-1 text-xs font-medium ${s.bg} ${s.color}`}>
                          {selectedArtifact.status.replace('_', ' ')}
                        </span>
                      );
                    })()}
                  </div>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Description</label>
                <p className="mt-1 text-sm text-gray-600">{selectedArtifact.description}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Content Reference</label>
                <div className="mt-1 flex items-center gap-2">
                  <code className="flex-1 rounded bg-gray-100 p-2 text-xs">{selectedArtifact.content_ref}</code>
                  <button className="rounded p-1 hover:bg-gray-100">
                    <ExternalLink className="h-4 w-4 text-gray-500" />
                  </button>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Content Hash</label>
                <p className="mt-1 font-mono text-xs text-gray-500">{selectedArtifact.content_hash}</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-700">Created By</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedArtifact.created_by}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-700">Created At</label>
                  <p className="mt-1 text-sm text-gray-600">{new Date(selectedArtifact.created_at).toLocaleString()}</p>
                </div>
              </div>

              {selectedArtifact.dependencies.length > 0 && (
                <div>
                  <label className="text-sm font-medium text-gray-700">Dependencies</label>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedArtifact.dependencies.map((dep) => (
                      <span key={dep} className="rounded bg-gray-100 px-2 py-1 text-xs text-gray-700">
                        {dep.slice(0, 8)}...
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
