"use client";

import { useState } from "react";

interface RevocationFormProps {
  onRevocationCreated?: () => void;
}

export function RevocationForm({ onRevocationCreated }: RevocationFormProps) {
  const [formData, setFormData] = useState({
    student_id: "",
    term_id: "",
    course_id: "",
    reason: "",
    requested_by: "",
    notes: "",
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage(null);

    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const response = await fetch(`${apiUrl}/api/issuer/revocations`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (data.success) {
        setMessage({ type: "success", text: data.data.message || "Revocation request created successfully" });
        setFormData({
          student_id: "",
          term_id: "",
          course_id: "",
          reason: "",
          requested_by: "",
          notes: "",
        });
        if (onRevocationCreated) {
          onRevocationCreated();
        }
      } else {
        setMessage({ type: "error", text: data.error || "Failed to create revocation request" });
      }
    } catch (error) {
      setMessage({ type: "error", text: "Network error - please try again" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
          <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Create Revocation Request</h2>
          <p className="text-sm text-gray-500">Enter validated credential corrections</p>
        </div>
      </div>

      {message && (
        <div className={`mb-6 p-4 rounded-lg ${
          message.type === "success" ? "bg-green-50 border border-green-200" : "bg-red-50 border border-red-200"
        }`}>
          <p className={`text-sm ${message.type === "success" ? "text-green-800" : "text-red-800"}`}>
            {message.text}
          </p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Student ID
          </label>
          <input
            type="text"
            value={formData.student_id}
            onChange={(e) => setFormData({ ...formData, student_id: e.target.value })}
            placeholder="ITITIU00001"
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Term ID
            </label>
            <input
              type="text"
              value={formData.term_id}
              onChange={(e) => setFormData({ ...formData, term_id: e.target.value })}
              placeholder="Semester_1_2023"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Course ID
            </label>
            <input
              type="text"
              value={formData.course_id}
              onChange={(e) => setFormData({ ...formData, course_id: e.target.value })}
              placeholder="IT089IU"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
              required
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Reason for Revocation
          </label>
          <textarea
            value={formData.reason}
            onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
            placeholder="e.g., Grade correction verified by professor - should be grade B"
            rows={3}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Requested By (Your Name/Email)
          </label>
          <input
            type="text"
            value={formData.requested_by}
            onChange={(e) => setFormData({ ...formData, requested_by: e.target.value })}
            placeholder="registrar@iu.edu.vn"
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Additional Notes (Optional)
          </label>
          <textarea
            value={formData.notes}
            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            placeholder="Additional context or validation details..."
            rows={2}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
          />
        </div>

        <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <svg className="w-5 h-5 text-amber-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div className="text-sm text-amber-800">
              <p className="font-medium mb-1">Automatic Processing</p>
              <p className="text-amber-700">
                This request will be automatically processed during the next term publication.
                Students will receive updated credentials at that time.
              </p>
            </div>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-gradient-to-r from-red-600 to-orange-600 text-white font-semibold py-4 px-6 rounded-lg hover:from-red-700 hover:to-orange-700 transition-all shadow-lg shadow-red-500/30 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Creating Request..." : "Create Revocation Request"}
        </button>
      </form>
    </div>
  );
}
