"use client";

import { useState } from "react";

type OperationStatus = {
  status: "idle" | "running" | "success" | "error";
  message?: string;
  output?: string;
};

export default function DemoDataPage() {
  const [resetStatus, setResetStatus] = useState<OperationStatus>({ status: "idle" });
  const [generateStatus, setGenerateStatus] = useState<OperationStatus>({ status: "idle" });
  const [addTermStatus, setAddTermStatus] = useState<OperationStatus>({ status: "idle" });

  // Add New Term state
  const [termId, setTermId] = useState("Semester_1_2026");

  // Full Data Generation state
  const [numStudentsForFull, setNumStudentsForFull] = useState(5);
  const [availableTerms, setAvailableTerms] = useState<string[]>([
    "Semester_1_2023",
    "Semester_2_2023",
    "Summer_2023",
    "Semester_1_2024",
    "Semester_2_2024",
    "Summer_2024",
    "Semester_1_2025",
  ]);
  const [selectedTerms, setSelectedTerms] = useState<string[]>([
    "Semester_1_2023",
    "Semester_2_2023",
    "Summer_2023",
    "Semester_1_2024",
    "Semester_2_2024",
    "Summer_2024",
    "Semester_1_2025",
  ]);

  const handleReset = async () => {
    if (!confirm("⚠️ WARNING: This will delete ALL generated data. Are you sure?")) {
      return;
    }

    setResetStatus({ status: "running", message: "Executing reset.sh..." });

    try {
      const response = await fetch("http://localhost:8080/api/demo/reset", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      const data = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Reset failed");
      }

      setResetStatus({
        status: "success",
        message: "System reset completed successfully",
        output: data.data?.output,
      });
    } catch (error: any) {
      setResetStatus({
        status: "error",
        message: error.message || "Failed to execute reset",
      });
    }
  };

  const handleGenerateFull = async () => {
    if (selectedTerms.length === 0) {
      alert("Please select at least one term to generate");
      return;
    }

    setGenerateStatus({
      status: "running",
      message: `Generating ${numStudentsForFull} students across ${selectedTerms.length} terms...`,
    });

    try {
      const response = await fetch("http://localhost:8080/api/demo/generate-full", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          num_students: numStudentsForFull,
          terms: selectedTerms,
        }),
      });

      const data = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Generation failed");
      }

      setGenerateStatus({
        status: "success",
        message: data.data?.message || "Full dataset generated successfully",
        output: data.data?.output,
      });
    } catch (error: any) {
      setGenerateStatus({
        status: "error",
        message: error.message || "Failed to execute generation",
      });
    }
  };

  const toggleTerm = (term: string) => {
    setSelectedTerms((prev) =>
      prev.includes(term) ? prev.filter((t) => t !== term) : [...prev, term]
    );
  };

  const toggleAllTerms = () => {
    setSelectedTerms((prev) =>
      prev.length === availableTerms.length ? [] : [...availableTerms]
    );
  };

  const handleAddTerm = () => {
    if (!termId.trim()) {
      alert("Please enter a term ID");
      return;
    }

    if (availableTerms.includes(termId)) {
      alert("This term already exists in the list");
      return;
    }

    // Add term to available terms and automatically select it
    setAvailableTerms((prev) => [...prev, termId]);
    setSelectedTerms((prev) => [...prev, termId]);

    // Clear the input and show success
    setTermId("");
    setAddTermStatus({
      status: "success",
      message: `Added "${termId}" to term list and selected it`,
    });

    // Clear success message after 3 seconds
    setTimeout(() => {
      setAddTermStatus({ status: "idle" });
    }, 3000);
  };

  const StatusBadge = ({ status, message, output }: OperationStatus) => {
    if (status === "idle") return null;

    const colors = {
      running: "bg-blue-50 border-blue-200 text-blue-800",
      success: "bg-green-50 border-green-200 text-green-800",
      error: "bg-red-50 border-red-200 text-red-800",
    };

    const icons = {
      running: "⏳",
      success: "✅",
      error: "❌",
    };

    return (
      <div className={`p-4 rounded-lg border ${colors[status]} mt-4`}>
        <p className="font-medium">
          {icons[status]} {message}
        </p>
        {output && (
          <details className="mt-2">
            <summary className="cursor-pointer text-sm font-medium">View Output</summary>
            <pre className="mt-2 p-2 bg-white rounded text-xs overflow-x-auto border">
              {output}
            </pre>
          </details>
        )}
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header with Navigation */}
        <div className="mb-8">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Demo Data Management</h1>
              <p className="text-gray-600 mt-2">
                Generate and manage demo data for testing the IU-MiCert system
              </p>
            </div>
            <div className="flex gap-3">
              <a
                href="/"
                className="bg-gray-600 text-white px-4 py-2 rounded-lg hover:bg-gray-700 transition-colors text-sm font-medium"
              >
                ← Back to Dashboard
              </a>
              <a
                href="/verifier"
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
              >
                🔍 Verifier
              </a>
            </div>
          </div>
        </div>

        {/* System Reset Section */}
        <div className="bg-white rounded-lg shadow p-6 mb-6 border-l-4 border-red-500">
          <h2 className="text-2xl font-bold mb-4 text-red-700">🧹 System Reset</h2>
          <p className="text-gray-600 mb-6">
            <strong>Warning:</strong> This will delete ALL generated data, including terms,
            receipts, and Verkle trees. Use this to start fresh.
          </p>

          <button
            onClick={handleReset}
            disabled={resetStatus.status === "running"}
            className="bg-red-600 text-white px-6 py-3 rounded-lg hover:bg-red-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {resetStatus.status === "running" ? "Resetting..." : "Execute Reset (./reset.sh)"}
          </button>

          <StatusBadge {...resetStatus} />
        </div>

        {/* Full Data Generation Section */}
        <div className="bg-white rounded-lg shadow p-6 mb-6 border-l-4 border-blue-500">
          <h2 className="text-2xl font-bold mb-4 text-blue-700">🚀 Full Data Generation</h2>
          <p className="text-gray-600 mb-6">
            Generate a complete dataset with custom parameters. This creates student journeys,
            converts them to Verkle format, builds trees, and generates receipts for all selected terms.
          </p>

          <div className="space-y-6 mb-6">
            {/* Number of Students Input */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Number of Students <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                value={numStudentsForFull}
                onChange={(e) => setNumStudentsForFull(parseInt(e.target.value) || 1)}
                min="1"
                max="100"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <p className="text-xs text-gray-500 mt-1">
                Number of students to generate (1-100). Each student gets 3-6 courses per term.
              </p>
            </div>

            {/* Term Selection */}
            <div>
              <div className="flex items-center justify-between mb-3">
                <label className="block text-sm font-medium text-gray-700">
                  Select Terms <span className="text-red-500">*</span>
                </label>
                <button
                  onClick={toggleAllTerms}
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                >
                  {selectedTerms.length === availableTerms.length ? "Deselect All" : "Select All"}
                </button>
              </div>
              <div className="grid grid-cols-2 gap-3 p-4 bg-gray-50 rounded-lg border border-gray-200">
                {availableTerms.map((term) => (
                  <label
                    key={term}
                    className="flex items-center space-x-3 cursor-pointer hover:bg-gray-100 p-2 rounded transition-colors"
                  >
                    <input
                      type="checkbox"
                      checked={selectedTerms.includes(term)}
                      onChange={() => toggleTerm(term)}
                      className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                    />
                    <span className="text-sm font-medium text-gray-700">{term}</span>
                  </label>
                ))}
              </div>
              <p className="text-xs text-gray-500 mt-2">
                Selected: {selectedTerms.length} / {availableTerms.length} terms
              </p>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-4 mb-4">
            <button
              onClick={handleGenerateFull}
              disabled={generateStatus.status === "running" || selectedTerms.length === 0}
              className="flex-1 bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
            >
              {generateStatus.status === "running"
                ? "Generating..."
                : `Generate ${numStudentsForFull} Students × ${selectedTerms.length} Terms`}
            </button>
          </div>

          <StatusBadge {...generateStatus} />

          {/* Add New Term to List */}
          <div className="mt-6 pt-6 border-t border-gray-200">
            <h3 className="text-lg font-bold mb-3 text-gray-900">➕ Add Custom Term</h3>
            <p className="text-gray-600 mb-4 text-sm">
              Add a new term to the selection list above (it will be automatically selected).
            </p>

            <div className="flex gap-4 items-end">
              {/* Term ID Input */}
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  New Term ID <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={termId}
                  onChange={(e) => setTermId(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      handleAddTerm();
                    }
                  }}
                  placeholder="e.g., Semester_1_2026"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Format: Semester_1_YYYY, Semester_2_YYYY, or Summer_YYYY
                </p>
              </div>

              {/* Add Button */}
              <button
                onClick={handleAddTerm}
                disabled={!termId.trim()}
                className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium whitespace-nowrap"
              >
                ➕ Add Term
              </button>
            </div>

            <StatusBadge {...addTermStatus} />
          </div>
        </div>

        {/* Info Section */}
        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <p className="text-sm text-blue-800">
            <strong>💡 Tip:</strong> After generating data, use the main issuer dashboard to
            process terms, build Verkle trees, and publish to blockchain.
          </p>
        </div>
      </div>
    </div>
  );
}
