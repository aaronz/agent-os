'use client';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import { Arbitration } from '@/types';
import { Scale, Gavel, AlertTriangle, CheckCircle, Clock, Shield, Users, Activity } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

const statusConfig: Record<string, { color: string; bg: string; icon: React.ReactNode }> = {
  pending: { color: 'text-yellow-700', bg: 'bg-yellow-100', icon: <Clock className="h-4 w-4" /> },
  ruled: { color: 'text-blue-700', bg: 'bg-blue-100', icon: <Gavel className="h-4 w-4" /> },
  closed: { color: 'text-gray-700', bg: 'bg-gray-100', icon: <CheckCircle className="h-4 w-4" /> },
};

const mockReputationData = [
  { name: 'High (>80)', value: 45, color: '#22c55e' },
  { name: 'Medium (50-80)', value: 30, color: '#eab308' },
  { name: 'Low (<50)', value: 25, color: '#ef4444' },
];

const mockActivityData = [
  { day: 'Mon', arbitrations: 3, disputes: 1 },
  { day: 'Tue', arbitrations: 5, disputes: 2 },
  { day: 'Wed', arbitrations: 2, disputes: 0 },
  { day: 'Thu', arbitrations: 8, disputes: 3 },
  { day: 'Fri', arbitrations: 4, disputes: 1 },
  { day: 'Sat', arbitrations: 1, disputes: 0 },
  { day: 'Sun', arbitrations: 2, disputes: 1 },
];

export default function GovernancePage() {
  const [activeTab, setActiveTab] = useState<'overview' | 'arbitrations' | 'rules'>('overview');
  const [selectedArbitration, setSelectedArbitration] = useState<Arbitration | null>(null);

  const { data: arbitrations, isLoading } = useQuery<Arbitration[]>({
    queryKey: ['arbitrations'],
    queryFn: () => api.get('/api/v1/arbitrations'),
    refetchInterval: 30000,
  });

  const { data: agents } = useQuery({
    queryKey: ['agents-summary'],
    queryFn: () => api.get('/api/v1/agents'),
  });

  const pendingCount = arbitrations?.filter(a => a.status === 'pending').length || 0;
  const totalAgents = agents?.length || 0;
  const avgReputation = agents?.reduce((acc, a) => acc + a.reputation, 0) / (totalAgents || 1) || 0;

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Governance</h1>
        <p className="text-gray-500">Monitor system health, disputes, and agent reputation</p>
      </div>

      <div className="mb-6 flex gap-2 border-b">
        {[
          { key: 'overview', label: 'Overview', icon: <Activity className="h-4 w-4" /> },
          { key: 'arbitrations', label: 'Arbitrations', icon: <Scale className="h-4 w-4" /> },
          { key: 'rules', label: 'Rules', icon: <Shield className="h-4 w-4" /> },
        ].map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key as typeof activeTab)}
            className={`flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors ${
              activeTab === tab.key
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {tab.icon}
            {tab.label}
            {tab.key === 'arbitrations' && pendingCount > 0 && (
              <span className="ml-1 rounded-full bg-red-100 px-2 py-0.5 text-xs text-red-600">
                {pendingCount}
              </span>
            )}
          </button>
        ))}
      </div>

      {activeTab === 'overview' && (
        <div className="space-y-6">
          <div className="grid gap-6 md:grid-cols-4">
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="rounded-lg bg-blue-100 p-2">
                  <Scale className="h-5 w-5 text-blue-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Pending Cases</p>
                  <p className="text-2xl font-bold text-gray-900">{pendingCount}</p>
                </div>
              </div>
            </div>
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="rounded-lg bg-green-100 p-2">
                  <Users className="h-5 w-5 text-green-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Total Agents</p>
                  <p className="text-2xl font-bold text-gray-900">{totalAgents}</p>
                </div>
              </div>
            </div>
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="rounded-lg bg-purple-100 p-2">
                  <Shield className="h-5 w-5 text-purple-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Avg Reputation</p>
                  <p className="text-2xl font-bold text-gray-900">{avgReputation.toFixed(1)}</p>
                </div>
              </div>
            </div>
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="rounded-lg bg-yellow-100 p-2">
                  <AlertTriangle className="h-5 w-5 text-yellow-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Active Rules</p>
                  <p className="text-2xl font-bold text-gray-900">12</p>
                </div>
              </div>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h3 className="mb-4 text-lg font-semibold text-gray-900">Arbitration Activity</h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={mockActivityData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="day" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="arbitrations" fill="#3b82f6" name="Arbitrations" />
                    <Bar dataKey="disputes" fill="#f59e0b" name="Disputes" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h3 className="mb-4 text-lg font-semibold text-gray-900">Agent Reputation Distribution</h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={mockReputationData}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={80}
                      paddingAngle={5}
                      dataKey="value"
                      label={({ name, value }) => `${name}: ${value}%`}
                    >
                      {mockReputationData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'arbitrations' && (
        <div className="space-y-4">
          {isLoading ? (
            Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="animate-pulse rounded-lg border bg-white p-6 shadow-sm">
                <div className="h-4 w-48 rounded bg-gray-200" />
              </div>
            ))
          ) : arbitrations?.length === 0 ? (
            <div className="p-8 text-center text-gray-500">No arbitrations found</div>
          ) : (
            arbitrations?.map((arb) => {
              const status = statusConfig[arb.status] || statusConfig.pending;
              return (
                <div
                  key={arb.id}
                  onClick={() => setSelectedArbitration(arb)}
                  className="cursor-pointer rounded-lg border bg-white p-6 shadow-sm hover:shadow-md"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <div className="flex items-center gap-2">
                        <h3 className="font-medium text-gray-900">{arb.type} Dispute</h3>
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${status.bg} ${status.color}`}>
                          {status.icon}
                          {arb.status}
                        </span>
                      </div>
                      <p className="mt-1 text-sm text-gray-500">ID: {arb.id}</p>
                    </div>
                    <span className="text-sm text-gray-400">{formatRelativeTime(arb.created_at)}</span>
                  </div>

                  <div className="mt-4">
                    <label className="text-xs font-medium text-gray-500">Claim</label>
                    <p className="mt-1 text-sm text-gray-700 line-clamp-2">{arb.claim}</p>
                  </div>

                  <div className="mt-4 flex items-center justify-between border-t pt-4">
                    <div className="flex gap-4">
                      <div>
                        <span className="text-xs text-gray-500">Applicant</span>
                        <p className="text-sm font-medium text-gray-700">{arb.applicant_id.slice(0, 8)}...</p>
                      </div>
                      <div>
                        <span className="text-xs text-gray-500">Respondent</span>
                        <p className="text-sm font-medium text-gray-700">{arb.respondent_id.slice(0, 8)}...</p>
                      </div>
                    </div>
                    {arb.status === 'ruled' && (
                      <span className={`text-sm font-medium ${arb.is_applicant_win ? 'text-green-600' : 'text-red-600'}`}>
                        {arb.is_applicant_win ? 'Applicant Wins' : 'Respondent Wins'}
                      </span>
                    )}
                  </div>
                </div>
              );
            })
          )}
        </div>
      )}

      {activeTab === 'rules' && (
        <div className="rounded-lg border bg-white p-6 shadow-sm">
          <div className="mb-6 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900">Governance Rules</h3>
            <button className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
              Add Rule
            </button>
          </div>

          <div className="space-y-4">
            {[
              { name: 'Bidding Window', description: 'Maximum time agents have to bid on tasks', enabled: true },
              { name: 'Minimum Reputation', description: 'Minimum reputation score required to bid', enabled: true },
              { name: 'Auto-fail Timeout', description: 'Automatically fail tasks that exceed max execution time', enabled: true },
              { name: 'Retry Limit', description: 'Maximum number of retries for failed tasks', enabled: true },
              { name: 'Review Required', description: 'Artifacts must be reviewed before approval', enabled: true },
              { name: 'Arbitration Threshold', description: 'Minimum confidence required for auto-ruling', enabled: false },
            ].map((rule) => (
              <div key={rule.name} className="flex items-center justify-between rounded-lg border p-4">
                <div>
                  <h4 className="font-medium text-gray-900">{rule.name}</h4>
                  <p className="text-sm text-gray-500">{rule.description}</p>
                </div>
                <button
                  className={`relative h-6 w-11 rounded-full transition-colors ${
                    rule.enabled ? 'bg-blue-600' : 'bg-gray-200'
                  }`}
                >
                  <span
                    className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${
                      rule.enabled ? 'translate-x-5' : 'translate-x-0.5'
                    }`}
                  />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {selectedArbitration && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="mx-4 max-h-[90vh] w-full max-w-2xl overflow-auto rounded-lg bg-white p-6 shadow-xl">
            <div className="mb-6 flex items-start justify-between">
              <div>
                <h2 className="text-xl font-bold text-gray-900">{selectedArbitration.type} Dispute</h2>
                <p className="text-sm text-gray-500 mt-1">ID: {selectedArbitration.id}</p>
              </div>
              <button onClick={() => setSelectedArbitration(null)} className="text-gray-400 hover:text-gray-600">
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div>
                {(() => {
                  const s = statusConfig[selectedArbitration.status] || statusConfig.pending;
                  return (
                    <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${s.bg} ${s.color}`}>
                      {s.icon}
                      {selectedArbitration.status}
                    </span>
                  );
                })()}
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Claim</label>
                <p className="mt-1 text-sm text-gray-600">{selectedArbitration.claim}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-700">Evidence</label>
                <p className="mt-1 text-sm text-gray-600">{selectedArbitration.evidence}</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-700">Applicant</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedArbitration.applicant_id}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-700">Respondent</label>
                  <p className="mt-1 text-sm text-gray-600">{selectedArbitration.respondent_id}</p>
                </div>
              </div>

              {selectedArbitration.status === 'ruled' && (
                <>
                  <div>
                    <label className="text-sm font-medium text-gray-700">Ruling</label>
                    <p className="mt-1 text-sm text-gray-600">{selectedArbitration.ruling}</p>
                  </div>

                  <div className="flex items-center gap-2">
                    <span className={`text-sm font-medium ${selectedArbitration.is_applicant_win ? 'text-green-600' : 'text-red-600'}`}>
                      {selectedArbitration.is_applicant_win ? 'Applicant Wins' : 'Respondent Wins'}
                    </span>
                  </div>

                  {selectedArbitration.penalty_decision.length > 0 && (
                    <div>
                      <label className="text-sm font-medium text-gray-700">Penalties</label>
                      <div className="mt-2 flex flex-wrap gap-2">
                        {selectedArbitration.penalty_decision.map((penalty, i) => (
                          <span key={i} className="rounded-full bg-red-100 px-2 py-1 text-xs text-red-700">
                            {penalty}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
