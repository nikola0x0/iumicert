"use client";

import { useState, useEffect } from "react";
import { useAccount } from "wagmi";
import { RevocationForm } from "@/components/revocation/RevocationForm";
import { RevocationList } from "@/components/revocation/RevocationList";
import { RevocationStats } from "@/components/revocation/RevocationStats";

export default function RevocationsPage() {
  const { isConnected } = useAccount();
  const [mounted, setMounted] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  useEffect(() => {
    setMounted(true);
  }, []);

  const handleRevocationCreated = () => {
    // Trigger refresh of list and stats
    setRefreshTrigger((prev) => prev + 1);
  };

  const handleRefreshAll = () => {
    // Manual refresh trigger for both stats and list
    setRefreshTrigger((prev) => prev + 1);
  };

  if (!mounted) {
    return null;
  }

  if (!isConnected) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="text-center max-w-md">
          <div className="w-20 h-20 bg-red-100 rounded-2xl flex items-center justify-center mx-auto mb-6">
            <svg
              className="w-10 h-10 text-red-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-slate-900 mb-4">
            Wallet Connection Required
          </h2>
          <p className="text-gray-600">
            Please connect your wallet to access revocation management.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="relative overflow-hidden bg-gradient-to-br from-red-500 via-orange-600 to-amber-600 rounded-3xl p-10 shadow-2xl shadow-red-500/20">
        <div className="absolute top-0 right-0 w-72 h-72 bg-white/5 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-12 -left-12 w-64 h-64 bg-white/5 rounded-full blur-3xl"></div>

        <div className="relative z-10">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 mb-4">
            <svg
              className="w-4 h-4 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <span className="text-xs font-semibold text-white/90">
              CREDENTIAL CORRECTIONS
            </span>
          </div>
          <div className="flex items-center justify-between">
            <h1 className="text-5xl font-extrabold text-white mb-3 tracking-tight">
              Revocation Management
            </h1>
            <button
              onClick={handleRefreshAll}
              className="flex items-center gap-2 px-4 py-2 bg-white/10 hover:bg-white/20 backdrop-blur-sm rounded-xl border border-white/20 text-white transition-all"
              title="Refresh all data"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span className="text-sm font-medium">Refresh</span>
            </button>
          </div>
          <p className="text-red-100 text-lg max-w-2xl leading-relaxed">
            Manage credential corrections and revocations. Approved requests
            will be automatically processed during the next term publication.
          </p>
        </div>
      </div>

      {/* Statistics Dashboard */}
      <RevocationStats refreshTrigger={refreshTrigger} />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Create New Revocation Request */}
        <div>
          <RevocationForm onRevocationCreated={handleRevocationCreated} />
        </div>

        {/* Revocation Requests List */}
        <div>
          <RevocationList refreshTrigger={refreshTrigger} />
        </div>
      </div>
    </div>
  );
}
