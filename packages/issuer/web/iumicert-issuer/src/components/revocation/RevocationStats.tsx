"use client";

import { useState, useEffect } from "react";

interface RevocationStatsProps {
  refreshTrigger?: number;
}

export function RevocationStats({ refreshTrigger }: RevocationStatsProps) {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, [refreshTrigger]);

  const loadStats = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const response = await fetch(`${apiUrl}/api/issuer/revocations/stats`);
      const data = await response.json();

      if (data.success) {
        setStats(data.data);
      }
    } catch (error) {
      console.error("Failed to load stats:", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="bg-white rounded-xl shadow-lg p-6 border border-gray-100 animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
            <div className="h-8 bg-gray-200 rounded w-1/2"></div>
          </div>
        ))}
      </div>
    );
  }

  const statCards = [
    {
      label: "Approved (Pending Processing)",
      value: stats?.approved_requests || 0,
      color: "green",
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
      description: "Will be processed during next term publication",
    },
    {
      label: "Processed",
      value: stats?.processed_requests || 0,
      color: "blue",
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
        </svg>
      ),
      description: "Successfully processed in past batches",
    },
    {
      label: "Pending Review",
      value: stats?.pending_requests || 0,
      color: "yellow",
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
      description: "Awaiting registrar approval",
    },
    {
      label: "Total Batches",
      value: stats?.total_batches || 0,
      color: "purple",
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      ),
      description: "Revocation batches processed",
    },
  ];

  const colorClasses: Record<string, { bg: string; text: string; border: string }> = {
    green: { bg: "bg-green-100", text: "text-green-600", border: "border-green-200" },
    blue: { bg: "bg-blue-100", text: "text-blue-600", border: "border-blue-200" },
    yellow: { bg: "bg-yellow-100", text: "text-yellow-600", border: "border-yellow-200" },
    purple: { bg: "bg-purple-100", text: "text-purple-600", border: "border-purple-200" },
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
      {statCards.map((card) => {
        const colors = colorClasses[card.color];
        return (
          <div
            key={card.label}
            className={`bg-white rounded-xl shadow-lg p-6 border ${colors.border} hover:shadow-xl transition-shadow`}
          >
            <div className="flex items-start justify-between mb-4">
              <div className={`w-12 h-12 ${colors.bg} rounded-xl flex items-center justify-center ${colors.text}`}>
                {card.icon}
              </div>
              {card.value > 0 && (
                <div className={`px-2 py-1 ${colors.bg} rounded-full`}>
                  <span className={`text-xs font-bold ${colors.text}`}>
                    {card.value > 99 ? "99+" : card.value}
                  </span>
                </div>
              )}
            </div>

            <h3 className="text-sm font-medium text-gray-600 mb-1">{card.label}</h3>
            <p className="text-3xl font-bold text-gray-900 mb-2">{card.value}</p>
            <p className="text-xs text-gray-500">{card.description}</p>
          </div>
        );
      })}
    </div>
  );
}
