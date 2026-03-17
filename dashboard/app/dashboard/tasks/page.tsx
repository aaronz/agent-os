'use client';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Task } from '@/types';
import { Search, Filter, Clock, User, CheckCircle, XCircle, AlertCircle, Eye } from 'lucide-react';

const statusConfig: Record<string, { color: string; bg: string; icon: React.ReactNode }> = {
  pending: { color: 'text-gray-600', bg: 'bg-gray-100', icon: <Clock className="h-4 w-4" /> },
  open: { color: 'text-blue-600', bg: 'bg-blue-100', icon: <AlertCircle className="h-4 w-4" /> },
  bidding: { color: 'text-yellow-600', bg: 'bg-yellow-100', icon: <AlertCircle className="h-4 w-4" /> },
  assigned: { color: 'text-purple-600', bg: 'bg-purple-100', icon: <User className="h-4 w-4" /> },
  executing: { color: 'text-orange-600', bg: 'bg-orange-100', icon: <Clock className="h-4 w-4" /> },
  reviewing: { color: 'text-indigo-600', bg: 'bg-indigo-100', icon: <Eye className="h-4 w-4" /> },
  completed: { color: 'text-green-600', bg: 'bg-green-100', icon: <CheckCircle className="h-4 w-4" /> },
  failed: { color: 'text-red-600', bg: 'bg-red-100', icon: <XCircle className="h-4 w-4" /> },
  cancelled: { color: 'text-gray-500', bg: 'bg-gray-100', icon: <XCircle className="h-4 w-4" /> },
};

const priorityConfig: Record<string, string> = {
  high: 'text-red-600 bg-red-50 border-red-200',
  medium: 'text-yellow-600 bg-yellow-50 border-yellow-200',
  low: 'text-green-600 bg-green-50 border-green-200',
};

export default function TasksPage() {
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [filter, setFilter] = useState<string>('all');

  const { data: tasks, isLoading } = useQuery<Task[]>({
    queryKey: ['tasks', filter],
    queryFn: () => api.get(filter !== 'all' ? `/api/v1/tasks?status=${filter}` : '/api/v1/tasks'),
    refetchInterval: 30000,
  });

  const filteredTasks = tasks?.filter(task => {
    if (filter === 'all') return true;
    return task.status === filter;
  });

  return (
    <div className="p-8">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Tasks</h1>
          <p className="text-gray-500">Monitor and manage task execution</p>
        </div>
        <div className="flex gap-2">
          {['all', 'pending', 'open', 'executing', 'completed', 'failed'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                filter === status
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </button>
          ))}
        </div>
      </div>

      <div className="mb-6">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search tasks by title or description..."
            className="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 text-sm focus:border-blue-500 focus:outline-none"
          />
        </div>
      </div>

      <div className="overflow-hidden rounded-lg border bg-white shadow-sm">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                Task
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                Priority
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                Assigned
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                Deadline
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 bg-white">
            {isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <tr key={i}>
                  <td colSpan={6} className="px-6 py-4">
                    <div className="animate-pulse flex items-center gap-4">
                      <div className="h-4 w-48 rounded bg-gray-200" />
                      <div className="h-4 w-24 rounded bg-gray-200" />
                    </div>
                  </td>
                </tr>
              ))
            ) : filteredTasks?.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                  No tasks found
                </td>
              </tr>
            ) : (
              filteredTasks?.map((task) => {
                const status = statusConfig[task.status] || statusConfig.pending;
                return (
                  <tr
                    key={task.id}
                    className="hover:bg-gray-50 cursor-pointer"
                    onClick={() => setSelectedTask(task)}
                  >
                    <td className="px-6 py-4">
                      <div>
                        <div className="font-medium text-gray-900">{task.title}</div>
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {task.description}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${status.bg} ${status.color}`}>
                        {status.icon}
                        {task.status}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded border px-2 py-0.5 text-xs font-medium ${priorityConfig[task.priority]}`}>
                        {task.priority}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {task.assigned_agent_id ? (
                        <span className="text-sm text-gray-900">{task.assigned_agent_id.slice(0, 8)}...</span>
                      ) : (
                        <span className="text-sm text-gray-400">Unassigned</span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-sm text-gray-500">
                        {task.deadline_at ? formatRelativeTime(task.deadline_at) : 'No deadline'}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                        View
                      </button>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* Task Detail Modal */}
      {selectedTask && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="mx-4 max-h-[90vh] w-full max-w-2xl overflow-auto rounded-lg bg-white p-6 shadow-xl">
            <div className="mb-6 flex items-start justify-between">
              <div>
                <h2 className="text-xl font-bold text-gray-900">{selectedTask.title}</h2>
                <p className="text-sm text-gray-500 mt-1">ID: {selectedTask.id}</p>
              </div>
              <button
                onClick={() => setSelectedTask(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium text-gray-700">Description</label>
                <p className="mt-1 text-sm text-gray-600">{selectedTask.description}</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-700">Status</label>
                  <div className="mt-1">
                    {(() => {
                      const status = statusConfig[selectedTask.status] || statusConfig.pending;
                      return (
                        <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${status.bg} ${status.color}`}>
                          {status.icon}
                          {selectedTask.status}
                        </span>
                      );
                    })()}
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-700">Priority</label>
                  <div className="mt-1">
                    <span className={`inline-flex rounded border px-2 py-0.5 text-xs font-medium ${priorityConfig[selectedTask.priority]}`}>
                      {selectedTask.priority}
                    </span>
                  </div>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Required Capabilities</label>
                <div className="mt-2 flex flex-wrap gap-2">
                  {selectedTask.required_capabilities.map((cap) => (
                    <span key={cap} className="rounded-full bg-blue-50 px-2 py-1 text-xs text-blue-700">
                      {cap}
                    </span>
                  ))}
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Dependencies</label>
                <div className="mt-2">
                  {selectedTask.dependencies.length > 0 ? (
                    <div className="flex flex-wrap gap-2">
                      {selectedTask.dependencies.map((dep) => (
                        <span key={dep} className="rounded bg-gray-100 px-2 py-1 text-xs text-gray-700">
                          {dep.slice(0, 8)}...
                        </span>
                      ))}
                    </div>
                  ) : (
                    <span className="text-sm text-gray-400">No dependencies</span>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-700">Estimated Duration</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedTask.estimated_duration_min} min</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-700">Max Execution Time</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedTask.max_execution_time_min} min</p>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Acceptance Criteria</label>
                <ul className="mt-2 list-inside list-disc text-sm text-gray-600">
                  {selectedTask.acceptance_criteria.map((criteria, i) => (
                    <li key={i}>{criteria}</li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
