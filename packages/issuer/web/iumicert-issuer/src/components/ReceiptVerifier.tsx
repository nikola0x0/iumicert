"use client";

import React, { useState } from "react";
import { apiService } from "@/lib/api";

interface Course {
  course_id?: string;
  course_name?: string;
  grade: string;
  credits: number;
}

interface TermReceipt {
  term_id: string;
  student_id: string;
  receipt: {
    student_id: string;
    term_id: string;
    revealed_courses: Course[];
    verkle_root: string;
    course_proofs: Record<string, any>;
    proof_type: string;
    selective_disclosure: boolean;
    verification_path: string;
    timestamp: string;
  };
  verkle_root: string;
  revealed_courses: number;
  total_courses: number;
  generated_at: string;
}

interface JourneyReceipt {
  student_id: string;
  receipt_type: {
    selective_disclosure: boolean;
    specific_courses: boolean;
    specific_terms: boolean;
  };
  generation_timestamp: string;
  terms_included: string[];
  courses_filter: string[];
  term_receipts: Record<string, TermReceipt>;
  blockchain_ready: boolean;
}

interface VerificationResult {
  verified?: boolean;
  status?: string;
  courses_verified?: number;
  courses_failed?: number;
  ipa_verified?: boolean;
  blockchain_verified?: boolean;
  blockchain_anchored?: boolean;
  blockchain_tx_hash?: string;
  blockchain_block?: number;
  blockchain_published_at?: string;
  details?: string;
  error?: string;
}

export default function ReceiptVerifier() {
  const [receipt, setReceipt] = useState<JourneyReceipt | null>(null);
  const [expandedTerms, setExpandedTerms] = useState<Set<string>>(new Set());
  const [expandedCourses, setExpandedCourses] = useState<Set<string>>(
    new Set()
  );
  const [verificationResults, setVerificationResults] = useState<
    Record<string, VerificationResult>
  >({});
  const [isVerifying, setIsVerifying] = useState(false);
  const [error, setError] = useState<string>("");

  // Handle file upload
  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const json = JSON.parse(e.target?.result as string);
        setReceipt(json);
        setError("");
        // Auto-expand all terms for better UX
        setExpandedTerms(new Set(Object.keys(json.term_receipts || {})));
      } catch (err) {
        setError("Invalid JSON file. Please upload a valid receipt.");
      }
    };
    reader.readAsText(file);
  };

  // Toggle term expansion
  const toggleTerm = (termId: string) => {
    setExpandedTerms((prev) => {
      const next = new Set(prev);
      if (next.has(termId)) {
        next.delete(termId);
      } else {
        next.add(termId);
      }
      return next;
    });
  };

  // Toggle course expansion
  const toggleCourse = (termId: string, courseId: string) => {
    const key = `${termId}-${courseId}`;
    setExpandedCourses((prev) => {
      const next = new Set(prev);
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }
      return next;
    });
  };

  // Verify entire journey
  const verifyJourney = async () => {
    if (!receipt) return;

    setIsVerifying(true);
    setError("");

    try {
      // Call IPA verification endpoint
      const result = await apiService.verifyReceiptIPA(receipt);

      // Store overall result
      setVerificationResults({
        overall: {
          verified: result.status === "success",
          ipa_verified: result.verified_courses === result.total_courses,
          blockchain_anchored: true,
          details: `${result.verified_courses}/${result.total_courses} courses verified`,
        },
        ...result.term_results,
      });
    } catch (err: any) {
      setError(err.message || "Verification failed");
      setVerificationResults({
        overall: {
          verified: false,
          error: err.message,
        },
      });
    } finally {
      setIsVerifying(false);
    }
  };

  if (!receipt) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-8">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h1 className="text-3xl font-bold text-gray-800 mb-6">
              🔍 Academic Receipt Verifier
            </h1>
            <p className="text-gray-600 mb-6">
              Upload a student's academic journey receipt to verify credentials
              and view their academic history.
            </p>

            <div className="border-2 border-dashed border-gray-300 rounded-lg p-12 text-center hover:border-indigo-500 transition-colors">
              <input
                type="file"
                accept=".json"
                onChange={handleFileUpload}
                className="hidden"
                id="receipt-upload"
              />
              <label
                htmlFor="receipt-upload"
                className="cursor-pointer flex flex-col items-center"
              >
                <svg
                  className="w-16 h-16 text-gray-400 mb-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
                <span className="text-xl font-semibold text-gray-700 mb-2">
                  Click to upload receipt
                </span>
                <span className="text-sm text-gray-500">
                  or drag and drop JSON file here
                </span>
              </label>
            </div>

            {error && (
              <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                {error}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  const overallResult = verificationResults.overall;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold text-gray-800">
                🎓 Academic Journey
              </h1>
              <p className="text-gray-600 mt-1">
                Student ID: <span className="font-mono">{receipt.student_id}</span>
              </p>
            </div>
            <button
              onClick={() => {
                setReceipt(null);
                setVerificationResults({});
              }}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
            >
              Upload New Receipt
            </button>
          </div>

          {/* Verification Button */}
          <div className="flex items-center gap-4">
            <button
              onClick={verifyJourney}
              disabled={isVerifying}
              className={`px-6 py-3 rounded-lg font-semibold ${
                isVerifying
                  ? "bg-gray-400 cursor-not-allowed"
                  : "bg-indigo-600 hover:bg-indigo-700"
              } text-white`}
            >
              {isVerifying ? "🔄 Verifying..." : "🔐 Verify Entire Journey"}
            </button>

            {overallResult && (
              <div
                className={`px-4 py-2 rounded-lg font-semibold ${
                  overallResult.verified
                    ? "bg-green-100 text-green-800"
                    : "bg-red-100 text-red-800"
                }`}
              >
                {overallResult.verified ? "✅ Verified" : "❌ Verification Failed"}
              </div>
            )}
          </div>

          {overallResult?.details && (
            <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <p className="text-sm text-blue-800">{overallResult.details}</p>
            </div>
          )}

          {error && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-800">{error}</p>
            </div>
          )}
        </div>

        {/* Journey Summary */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">
            📊 Journey Summary
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-blue-600">
                {Object.keys(receipt.term_receipts).length}
              </div>
              <div className="text-sm text-gray-600">Terms Completed</div>
            </div>
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-green-600">
                {Object.values(receipt.term_receipts).reduce(
                  (sum, term) => sum + term.total_courses,
                  0
                )}
              </div>
              <div className="text-sm text-gray-600">Total Courses</div>
            </div>
            <div className="bg-purple-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-purple-600">
                {receipt.blockchain_ready ? "✅" : "⏳"}
              </div>
              <div className="text-sm text-gray-600">Blockchain Ready</div>
            </div>
            <div className="bg-indigo-50 p-4 rounded-lg">
              <div className="text-2xl font-bold text-indigo-600">
                {receipt.receipt_type.selective_disclosure ? "🔒" : "📖"}
              </div>
              <div className="text-sm text-gray-600">
                {receipt.receipt_type.selective_disclosure
                  ? "Selective"
                  : "Full Disclosure"}
              </div>
            </div>
          </div>
        </div>

        {/* Terms List */}
        <div className="space-y-4">
          {Object.entries(receipt.term_receipts)
            .sort(([termIdA], [termIdB]) => {
              // Sort chronologically: Semester_1_YYYY -> Semester_2_YYYY -> Summer_YYYY
              const extractOrder = (id: string) => {
                const yearMatch = id.match(/(\d{4})/);
                const year = yearMatch ? parseInt(yearMatch[1]) : 0;
                const semester = id.includes("Semester_1") ? 1 : id.includes("Semester_2") ? 2 : 3;
                return year * 10 + semester;
              };
              return extractOrder(termIdA) - extractOrder(termIdB);
            })
            .map(([termId, termData]) => {
              const isExpanded = expandedTerms.has(termId);
              const termResult = verificationResults[termId];

            return (
              <div
                key={termId}
                className="bg-white rounded-lg shadow-lg overflow-hidden"
              >
                {/* Term Header */}
                <div
                  className="p-6 cursor-pointer hover:bg-gray-50 transition-colors"
                  onClick={() => toggleTerm(termId)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="text-2xl">
                        {isExpanded ? "📂" : "📁"}
                      </span>
                      <div>
                        <h3 className="text-xl font-bold text-gray-800">
                          {termId.replace(/_/g, " ")}
                        </h3>
                        <p className="text-sm text-gray-600">
                          {termData.revealed_courses} of {termData.total_courses}{" "}
                          courses
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      {termResult && (
                        <div
                          className={`px-3 py-1 rounded-full text-sm font-semibold ${
                            termResult.status === "completed" && termResult.courses_verified > 0
                              ? "bg-green-100 text-green-800"
                              : "bg-red-100 text-red-800"
                          }`}
                        >
                          {termResult.status === "completed" && termResult.courses_verified > 0
                            ? `✅ Verified (${termResult.courses_verified} courses)`
                            : "❌ Failed"}
                        </div>
                      )}
                      <span className="text-gray-400">
                        {isExpanded ? "▼" : "▶"}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Term Content */}
                {isExpanded && (
                  <div className="border-t border-gray-200 bg-gray-50 p-6">
                    {/* Verification Status */}
                    {termResult && termResult.blockchain_verified && (
                      <div className="mb-4 p-4 bg-green-50 rounded-lg border border-green-200">
                        <div className="flex items-start gap-3">
                          <span className="text-2xl">⛓️</span>
                          <div className="flex-1">
                            <h4 className="font-semibold text-green-800 mb-2">
                              Blockchain Verified
                            </h4>
                            <div className="space-y-2 text-sm">
                              <div>
                                <span className="text-green-700">Transaction:</span>
                                <a
                                  href={`https://sepolia.etherscan.io/tx/${termResult.blockchain_tx_hash}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="ml-2 font-mono text-xs text-blue-600 hover:text-blue-800 underline break-all"
                                >
                                  {termResult.blockchain_tx_hash}
                                </a>
                              </div>
                              <div>
                                <span className="text-green-700">Published:</span>
                                <span className="ml-2 text-green-800">
                                  {new Date(termResult.blockchain_published_at).toLocaleString()}
                                </span>
                              </div>
                              <div>
                                <span className="text-green-700">Block:</span>
                                <span className="ml-2 font-mono text-green-800">
                                  {termResult.blockchain_block || "N/A"}
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}

                    {/* Term Metadata */}
                    <div className="mb-4 p-4 bg-white rounded-lg border border-gray-200">
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <span className="text-gray-600">Verkle Root:</span>
                          <p className="font-mono text-xs mt-1 break-all text-gray-800">
                            {termData.verkle_root}
                          </p>
                        </div>
                        <div>
                          <span className="text-gray-600">Generated:</span>
                          <p className="text-gray-800 mt-1">
                            {new Date(termData.generated_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>

                    {/* Courses */}
                    <div className="space-y-2">
                      <h4 className="font-semibold text-gray-700 mb-3">
                        📚 Courses:
                      </h4>
                      {termData.receipt.revealed_courses.map((course, idx) => {
                        const courseKey = `${termId}-${course.course_id || idx}`;
                        const isCourseExpanded = expandedCourses.has(courseKey);

                        return (
                          <div
                            key={courseKey}
                            className="bg-white rounded-lg border border-gray-200 overflow-hidden"
                          >
                            {/* Course Header */}
                            <div
                              className="p-4 cursor-pointer hover:bg-gray-50 transition-colors"
                              onClick={() =>
                                toggleCourse(termId, course.course_id || idx.toString())
                              }
                            >
                              <div className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                  <span className="text-lg">
                                    {isCourseExpanded ? "📖" : "📕"}
                                  </span>
                                  <div>
                                    <p className="font-semibold text-gray-800">
                                      {course.course_id || `Course ${idx + 1}`}
                                    </p>
                                    <p className="text-sm text-gray-600">
                                      {course.course_name || "N/A"}
                                    </p>
                                  </div>
                                </div>
                                <div className="flex items-center gap-3">
                                  <span
                                    className={`px-3 py-1 rounded-full text-sm font-semibold ${
                                      course.grade === "A" || course.grade === "A+"
                                        ? "bg-green-100 text-green-800"
                                        : course.grade === "B" ||
                                          course.grade === "B+"
                                        ? "bg-blue-100 text-blue-800"
                                        : course.grade === "C" ||
                                          course.grade === "C+"
                                        ? "bg-yellow-100 text-yellow-800"
                                        : "bg-gray-100 text-gray-800"
                                    }`}
                                  >
                                    Grade: {course.grade}
                                  </span>
                                  <span className="text-sm text-gray-600">
                                    {course.credits} credits
                                  </span>
                                  <span className="text-gray-400">
                                    {isCourseExpanded ? "▼" : "▶"}
                                  </span>
                                </div>
                              </div>
                            </div>

                            {/* Course Details (Proof) */}
                            {isCourseExpanded && (
                              <div className="border-t border-gray-200 bg-gray-50 p-4">
                                <h5 className="text-sm font-semibold text-gray-700 mb-2">
                                  🔐 Cryptographic Proof:
                                </h5>
                                <div className="bg-white p-3 rounded border border-gray-200">
                                  <pre className="text-xs font-mono text-gray-700 overflow-x-auto">
                                    {JSON.stringify(
                                      termData.receipt.course_proofs[
                                        course.course_id || idx.toString()
                                      ],
                                      null,
                                      2
                                    )}
                                  </pre>
                                </div>
                                <div className="mt-3 flex items-center gap-2">
                                  <span className="text-xs text-gray-600">
                                    Proof Type:
                                  </span>
                                  <span className="px-2 py-1 bg-indigo-100 text-indigo-800 text-xs rounded">
                                    {termData.receipt.proof_type}
                                  </span>
                                  <span className="text-xs text-gray-600">
                                    Verification Path:
                                  </span>
                                  <span className="px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded">
                                    {termData.receipt.verification_path}
                                  </span>
                                </div>
                              </div>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
