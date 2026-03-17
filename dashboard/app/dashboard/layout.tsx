'use client';

import { useState } from 'react';
import { Sidebar } from "@/components/sidebar";
import { SearchDialog } from "@/components/search-dialog";
import { Search } from "lucide-react";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [searchOpen, setSearchOpen] = useState(false);

  return (
    <div className="flex h-screen">
      <Sidebar />
      <main className="flex-1 overflow-auto bg-gray-50">
        <div className="sticky top-0 z-40 flex h-16 items-center justify-between border-b bg-white px-6">
          <div />
          <button
            onClick={() => setSearchOpen(true)}
            className="flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-1.5 text-sm text-gray-400 hover:border-gray-300 hover:bg-gray-100"
          >
            <Search className="h-4 w-4" />
            <span>Search...</span>
            <kbd className="ml-2 rounded border bg-white px-1.5 py-0.5 text-xs">
              ⌘K
            </kbd>
          </button>
        </div>
        {children}
      </main>
      <SearchDialog open={searchOpen} onOpenChange={setSearchOpen} />
    </div>
  );
}
