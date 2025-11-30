"use client";

import { useState, useEffect } from "react";

interface RevocationRequest {
  RequestID: string;
  StudentID: string;
  TermID: string;
  CourseID: string;
  Reason: string;
  RequestedBy: string;
  Status: string;
  CreatedAt: string;
  Notes?: string;
}

interface RevocationListProps {
  refreshTrigger?: number;
}

export function RevocationList({ refreshTrigger }: RevocationListProps) {
  const [requests, setRequests] = useState<RevocationRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<string>("approved");

  useEffect(() => {
    loadRequests();
  }, [refreshTrigger, filter]);

  const loadRequests = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const url = filter
        ? `${apiUrl}/api/issuer/revocations?status=${filter}`
        : `${apiUrl}/api/issuer/revocations`;

      const response = await fetch(url);
      const data = await response.json();

      if (data.success) {
        setRequests(data.data.requests || []);
      }
    } catch (error) {
      console.error("Failed to load revocation requests:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (requestID: string) => {
    if (!confirm("Are you sure you want to delete this revocation request?")) {
      return;
    }

    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const response = await fetch(`${apiUrl}/api/issuer/revocations/${requestID}`, {
        method: "DELETE",
      });

      const data = await response.json();

      if (data.success) {
        loadRequests(); // Refresh list
      } else {
        alert("Failed to delete request");
      }
    } catch (error) {
      alert("Network error");
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      approved: "bg-green-100 text-green-800 border-green-200",
      pending: "bg-yellow-100 text-yellow-800 border-yellow-200",
      processed: "bg-blue-100 text-blue-800 border-blue-200",
      rejected: "bg-gray-100 text-gray-800 border-gray-200",
    };

    return (
      <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${styles[status] || styles.pending}`}>
        {status.toUpperCase()}
      </span>
    );
  };

  return (
    <div className="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Revocation Requests</h2>
            <p className="text-sm text-gray-500">{requests.length} total</p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => loadRequests()}
            disabled={loading}
            className="p-2 text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors disabled:opacity-50"
            title="Refresh list"
          >
            <svg
              className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">All</option>
            <option value="approved">Approved</option>
            <option value="pending">Pending</option>
            <option value="processed">Processed</option>
            <option value="rejected">Rejected</option>
          </select>
        </div>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-500 border-t-transparent"></div>
          <p className="mt-4 text-gray-500">Loading requests...</p>
        </div>
      ) : requests.length === 0 ? (
        <div className="text-center py-12">
          <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p className="text-gray-500">No revocation requests found</p>
        </div>
      ) : (
        <div className="space-y-4 max-h-[600px] overflow-y-auto">
          {requests.map((request) => (
            <div
              key={request.RequestID}
              className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="font-semibold text-gray-900">
                      {request.StudentID} / {request.CourseID}
                    </h3>
                    {getStatusBadge(request.Status)}
                  </div>
                  <p className="text-sm text-gray-600">{request.TermID}</p>
                </div>

                {request.Status === "approved" && (
                  <button
                    onClick={() => handleDelete(request.RequestID)}
                    className="text-red-600 hover:text-red-800 p-2"
                    title="Delete request"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                )}
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex items-start gap-2">
                  <span className="font-medium text-gray-700">Reason:</span>
                  <span className="text-gray-600 flex-1">{request.Reason}</span>
                </div>

                {request.Notes && (
                  <div className="flex items-start gap-2">
                    <span className="font-medium text-gray-700">Notes:</span>
                    <span className="text-gray-600 flex-1">{request.Notes}</span>
                  </div>
                )}

                <div className="flex items-center gap-4 text-xs text-gray-500 pt-2 border-t border-gray-100">
                  <span>By: {request.RequestedBy}</span>
                  <span>•</span>
                  <span>{formatDate(request.CreatedAt)}</span>
                  <span>•</span>
                  <span className="font-mono text-xs">{request.RequestID.slice(0, 20)}...</span>
                </div>
              </div>

              {request.Status === "approved" && (
                <div className="mt-3 pt-3 border-t border-amber-100 bg-amber-50 -mx-4 -mb-4 px-4 py-3 rounded-b-lg">
                  <p className="text-xs text-amber-800 flex items-center gap-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="font-medium">Pending Processing:</span>
                    Will be processed automatically during next term publication
                  </p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
