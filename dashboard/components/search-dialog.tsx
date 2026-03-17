'use client';

import { useEffect, useState } from 'react';
import { Command } from 'cmdk';
import { useRouter } from 'next/navigation';
import { Search, LayoutDashboard, Target, CheckSquare, Users, FileBox, Brain, Shield, Settings, User } from 'lucide-react';

interface SearchDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const items = [
  { category: 'Pages', name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { category: 'Pages', name: 'Intents', href: '/dashboard/intents', icon: Target },
  { category: 'Pages', name: 'Tasks', href: '/dashboard/tasks', icon: CheckSquare },
  { category: 'Pages', name: 'Agents', href: '/dashboard/agents', icon: Users },
  { category: 'Pages', name: 'Artifacts', href: '/dashboard/artifacts', icon: FileBox },
  { category: 'Pages', name: 'Memory', href: '/dashboard/memory', icon: Brain },
  { category: 'Pages', name: 'Governance', href: '/dashboard/governance', icon: Shield },
  { category: 'Pages', name: 'Settings', href: '/dashboard/settings', icon: Settings },
  { category: 'Actions', name: 'Create New Intent', href: '/dashboard/intents/new', icon: Target },
  { category: 'Actions', name: 'View Profile', href: '/dashboard/settings', icon: User },
];

export function SearchDialog({ open, onOpenChange }: SearchDialogProps) {
  const router = useRouter();
  const [search, setSearch] = useState('');

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        onOpenChange(!open);
      }
    };

    document.addEventListener('keydown', down);
    return () => document.removeEventListener('keydown', down);
  }, [open, onOpenChange]);

  const handleSelect = (href: string) => {
    router.push(href);
    onOpenChange(false);
    setSearch('');
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center bg-black/50 pt-[20vh]">
      <div className="w-full max-w-xl rounded-xl bg-white shadow-2xl">
        <Command className="rounded-xl">
          <div className="flex items-center border-b px-3">
            <Search className="h-5 w-5 text-gray-400" />
            <Command.Input
              value={search}
              onValueChange={setSearch}
              placeholder="Search pages, actions..."
              className="flex-1 border-0 bg-transparent py-4 pl-3 pr-4 text-sm outline-none placeholder:text-gray-400"
            />
            <kbd className="hidden rounded border bg-gray-100 px-1.5 py-0.5 text-xs text-gray-400 sm:inline-block">
              ESC
            </kbd>
          </div>
          <Command.List className="max-h-72 overflow-y-auto p-2">
            <Command.Empty className="py-6 text-center text-sm text-gray-500">
              No results found.
            </Command.Empty>
            {['Pages', 'Actions'].map((category) => (
              <Command.Group key={category} heading={category}>
                {items
                  .filter((item) => item.category === category)
                  .filter((item) =>
                    item.name.toLowerCase().includes(search.toLowerCase())
                  )
                  .map((item) => (
                    <Command.Item
                      key={item.href}
                      onSelect={() => handleSelect(item.href)}
                      className="flex cursor-pointer items-center gap-3 rounded-lg px-3 py-2 text-sm aria-selected:bg-blue-50 aria-selected:text-blue-700"
                    >
                      <item.icon className="h-4 w-4 text-gray-400" />
                      {item.name}
                    </Command.Item>
                  ))}
              </Command.Group>
            ))}
          </Command.List>
        </Command>
      </div>
    </div>
  );
}
