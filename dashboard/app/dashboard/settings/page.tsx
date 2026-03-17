'use client';

import { useState } from 'react';
import { useAuthStore } from '@/hooks/use-auth';
import { Settings, User, Building, Bell, Shield, Palette, Key, Globe, Moon, Sun, Monitor } from 'lucide-react';

type Theme = 'light' | 'dark' | 'system';

export default function SettingsPage() {
  const { user } = useAuthStore();
  const [theme, setTheme] = useState<Theme>('system');
  const [activeTab, setActiveTab] = useState<'profile' | 'organization' | 'notifications' | 'security' | 'appearance'>('profile');

  const tabs = [
    { key: 'profile', label: 'Profile', icon: <User className="h-4 w-4" /> },
    { key: 'organization', label: 'Organization', icon: <Building className="h-4 w-4" /> },
    { key: 'notifications', label: 'Notifications', icon: <Bell className="h-4 w-4" /> },
    { key: 'security', label: 'Security', icon: <Shield className="h-4 w-4" /> },
    { key: 'appearance', label: 'Appearance', icon: <Palette className="h-4 w-4" /> },
  ];

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-gray-500">Manage your account and preferences</p>
      </div>

      <div className="flex gap-8">
        <div className="w-56 shrink-0">
          <nav className="space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key as typeof activeTab)}
                className={`flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                  activeTab === tab.key
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                {tab.icon}
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        <div className="flex-1">
          {activeTab === 'profile' && (
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h2 className="text-lg font-semibold text-gray-900">Profile Settings</h2>
              <p className="text-sm text-gray-500 mb-6">Manage your personal information</p>

              <div className="space-y-6">
                <div className="flex items-center gap-6">
                  <div className="h-20 w-20 rounded-full bg-blue-100 flex items-center justify-center">
                    <User className="h-10 w-10 text-blue-600" />
                  </div>
                  <div>
                    <button className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium hover:bg-gray-50">
                      Change Avatar
                    </button>
                    <p className="mt-1 text-xs text-gray-500">JPG, PNG or GIF. Max 2MB.</p>
                  </div>
                </div>

                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Full Name</label>
                    <input
                      type="text"
                      defaultValue={user?.name || 'Admin User'}
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Email</label>
                    <input
                      type="email"
                      defaultValue={user?.email || 'admin@example.com'}
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Role</label>
                  <input
                    type="text"
                    defaultValue="Administrator"
                    disabled
                    className="mt-1 block w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-500"
                  />
                </div>

                <div className="flex justify-end">
                  <button className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
                    Save Changes
                  </button>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'organization' && (
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h2 className="text-lg font-semibold text-gray-900">Organization Settings</h2>
              <p className="text-sm text-gray-500 mb-6">Manage your organization details</p>

              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Organization Name</label>
                  <input
                    type="text"
                    defaultValue="Acme Corp"
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Description</label>
                  <textarea
                    rows={3}
                    defaultValue="Agent collaboration platform for autonomous task execution"
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Timezone</label>
                  <div className="relative mt-1">
                    <Globe className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                    <select className="w-full rounded-lg border border-gray-300 py-2 pl-10 pr-4 text-sm focus:border-blue-500 focus:outline-none">
                      <option>UTC</option>
                      <option>America/New_York</option>
                      <option>America/Los_Angeles</option>
                      <option>Europe/London</option>
                      <option>Asia/Shanghai</option>
                    </select>
                  </div>
                </div>

                <div className="flex justify-end">
                  <button className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
                    Save Changes
                  </button>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'notifications' && (
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h2 className="text-lg font-semibold text-gray-900">Notification Preferences</h2>
              <p className="text-sm text-gray-500 mb-6">Choose what notifications you receive</p>

              <div className="space-y-4">
                {[
                  { label: 'Task Completed', description: 'Get notified when tasks are completed', enabled: true },
                  { label: 'Task Failed', description: 'Get notified when tasks fail', enabled: true },
                  { label: 'New Bid', description: 'Get notified when agents bid on tasks', enabled: false },
                  { label: 'Arbitration', description: 'Get notified about new arbitration cases', enabled: true },
                  { label: 'Agent Status', description: 'Get notified about agent status changes', enabled: false },
                  { label: 'Weekly Summary', description: 'Receive weekly activity summary', enabled: true },
                ].map((item) => (
                  <div key={item.label} className="flex items-center justify-between rounded-lg border p-4">
                    <div>
                      <p className="font-medium text-gray-900">{item.label}</p>
                      <p className="text-sm text-gray-500">{item.description}</p>
                    </div>
                    <button
                      className={`relative h-6 w-11 rounded-full transition-colors ${
                        item.enabled ? 'bg-blue-600' : 'bg-gray-200'
                      }`}
                    >
                      <span
                        className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${
                          item.enabled ? 'translate-x-5' : 'translate-x-0.5'
                        }`}
                      />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'security' && (
            <div className="space-y-6">
              <div className="rounded-lg border bg-white p-6 shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900">Change Password</h2>
                <p className="text-sm text-gray-500 mb-6">Update your password regularly for security</p>

                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Current Password</label>
                    <input
                      type="password"
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">New Password</label>
                    <input
                      type="password"
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Confirm New Password</label>
                    <input
                      type="password"
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                  <button className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
                    Update Password
                  </button>
                </div>
              </div>

              <div className="rounded-lg border bg-white p-6 shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900">API Keys</h2>
                <p className="text-sm text-gray-500 mb-6">Manage your API keys for programmatic access</p>

                <div className="space-y-4">
                  <div className="flex items-center justify-between rounded-lg border p-4">
                    <div className="flex items-center gap-3">
                      <Key className="h-5 w-5 text-gray-400" />
                      <div>
                        <p className="font-medium text-gray-900">Production Key</p>
                        <p className="text-sm text-gray-500">sk-...a1b2</p>
                      </div>
                    </div>
                    <button className="text-sm text-red-600 hover:text-red-800">Revoke</button>
                  </div>
                  <button className="flex items-center gap-2 text-sm text-blue-600 hover:text-blue-800">
                    <Key className="h-4 w-4" />
                    Generate New Key
                  </button>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'appearance' && (
            <div className="rounded-lg border bg-white p-6 shadow-sm">
              <h2 className="text-lg font-semibold text-gray-900">Appearance</h2>
              <p className="text-sm text-gray-500 mb-6">Customize the look and feel</p>

              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-3">Theme</label>
                  <div className="grid gap-3 md:grid-cols-3">
                    {[
                      { key: 'light', label: 'Light', icon: <Sun className="h-5 w-5" /> },
                      { key: 'dark', label: 'Dark', icon: <Moon className="h-5 w-5" /> },
                      { key: 'system', label: 'System', icon: <Monitor className="h-5 w-5" /> },
                    ].map((option) => (
                      <button
                        key={option.key}
                        onClick={() => setTheme(option.key as Theme)}
                        className={`flex flex-col items-center gap-2 rounded-lg border-2 p-4 transition-colors ${
                          theme === option.key
                            ? 'border-blue-600 bg-blue-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                      >
                        {option.icon}
                        <span className="text-sm font-medium">{option.label}</span>
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-3">Accent Color</label>
                  <div className="flex gap-3">
                    {['#3b82f6', '#8b5cf6', '#ec4899', '#14b8a6', '#f59e0b', '#ef4444'].map((color) => (
                      <button
                        key={color}
                        className="h-10 w-10 rounded-full border-2 border-white shadow-sm hover:scale-110 transition-transform"
                        style={{ backgroundColor: color }}
                      />
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
