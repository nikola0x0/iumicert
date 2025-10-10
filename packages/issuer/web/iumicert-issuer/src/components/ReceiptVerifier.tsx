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
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900">Receipt Verifier</h1>
          <p className="text-gray-600 mt-2">
            Upload a student's academic journey receipt to verify credentials
            and view their academic history
          </p>
        </div>

        <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-8">
          <div className="border-2 border-dashed border-gray-300 rounded-xl p-12 text-center hover:border-blue-400 hover:bg-blue-50/30 transition-all duration-200">
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
              <div className="w-16 h-16 bg-blue-100 rounded-2xl flex items-center justify-center mb-4">
                <svg
                  className="w-8 h-8 text-blue-600"
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
              </div>
              <span className="text-xl font-semibold text-slate-900 mb-2">
                Click to upload receipt
              </span>
              <span className="text-sm text-gray-500">
                or drag and drop JSON file here
              </span>
            </label>
          </div>

          {error && (
            <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 flex items-start gap-3">
              <svg className="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
              <span>{error}</span>
            </div>
          )}
        </div>
      </div>
    );
  }

  const overallResult = verificationResults.overall;

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Header */}
      <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-slate-900">
              Academic Journey
            </h1>
            <p className="text-gray-600 mt-1">
              Student ID: <span className="font-mono font-semibold text-slate-900">{receipt.student_id}</span>
            </p>
          </div>
          <button
            onClick={() => {
              setReceipt(null);
              setVerificationResults({});
            }}
            className="px-5 py-2.5 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors font-medium"
          >
            Upload New Receipt
          </button>
        </div>

        {/* Verification Button */}
        <div className="flex items-center gap-4">
          <button
            onClick={verifyJourney}
            disabled={isVerifying}
            className={`px-6 py-3 rounded-xl font-semibold transition-all duration-200 ${
              isVerifying
                ? "bg-gray-400 cursor-not-allowed"
                : "bg-gradient-to-r from-blue-500 to-blue-600 hover:shadow-lg hover:shadow-blue-500/30"
            } text-white`}
          >
            {isVerifying ? (
              <span className="flex items-center gap-2">
                <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Verifying...
              </span>
            ) : (
              "Verify Entire Journey"
            )}
          </button>

          {overallResult && (
            <div
              className={`px-4 py-2 rounded-xl font-semibold border ${
                overallResult.verified
                  ? "bg-green-100 text-green-800 border-green-200"
                  : "bg-red-100 text-red-800 border-red-200"
              }`}
            >
              {overallResult.verified ? "Verified" : "Verification Failed"}
            </div>
          )}
        </div>

        {overallResult?.details && (
          <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-xl">
            <p className="text-sm text-blue-800">{overallResult.details}</p>
          </div>
        )}

        {error && (
          <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-3">
            <svg className="w-5 h-5 mt-0.5 flex-shrink-0 text-red-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}
      </div>

      {/* Journey Summary */}
      <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <h2 className="text-xl font-bold text-slate-900 mb-4">
          Journey Summary
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-blue-50 p-5 rounded-xl border border-blue-100">
            <div className="text-2xl font-bold text-blue-600">
              {Object.keys(receipt.term_receipts).length}
            </div>
            <div className="text-sm text-gray-600 mt-1">Terms Completed</div>
          </div>
          <div className="bg-green-50 p-5 rounded-xl border border-green-100">
            <div className="text-2xl font-bold text-green-600">
              {Object.values(receipt.term_receipts).reduce(
                (sum, term) => sum + term.total_courses,
                0
              )}
            </div>
            <div className="text-sm text-gray-600 mt-1">Total Courses</div>
          </div>
          <div className="bg-purple-50 p-5 rounded-xl border border-purple-100">
            <div className="text-2xl font-bold text-purple-600">
              {receipt.blockchain_ready ? "Ready" : "Pending"}
            </div>
            <div className="text-sm text-gray-600 mt-1">Blockchain Status</div>
          </div>
          <div className="bg-indigo-50 p-5 rounded-xl border border-indigo-100">
            <div className="text-sm font-bold text-indigo-600">
              {receipt.receipt_type.selective_disclosure
                ? "Selective"
                : "Full"}
            </div>
            <div className="text-sm text-gray-600 mt-1">Disclosure Type</div>
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
              className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden"
            >
              {/* Term Header */}
              <div
                className="p-6 cursor-pointer hover:bg-gray-50/50 transition-colors"
                onClick={() => toggleTerm(termId)}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${isExpanded ? 'bg-blue-100' : 'bg-gray-100'}`}>
                      <svg className={`w-5 h-5 ${isExpanded ? 'text-blue-600' : 'text-gray-600'}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        {isExpanded ? (
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z" />
                        ) : (
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        )}
                      </svg>
                    </div>
                    <div>
                      <h3 className="text-xl font-bold text-slate-900">
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
                        className={`px-4 py-1.5 rounded-full text-sm font-semibold border ${
                          termResult.status === "completed" && (termResult.courses_verified || 0) > 0
                            ? "bg-green-100 text-green-800 border-green-200"
                            : "bg-red-100 text-red-800 border-red-200"
                        }`}
                      >
                        {termResult.status === "completed" && (termResult.courses_verified || 0) > 0
                          ? `Verified (${termResult.courses_verified || 0} courses)`
                          : "Failed"}
                      </div>
                    )}
                    <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={isExpanded ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"} />
                    </svg>
                  </div>
                </div>
              </div>

              {/* Term Content */}
              {isExpanded && (
                <div className="border-t border-gray-200 bg-gray-50 p-6">
                  {/* Verification Status */}
                  {termResult && termResult.blockchain_verified && (
                    <div className="mb-4 p-4 bg-green-50 rounded-xl border border-green-200">
                      <div className="flex items-start gap-3">
                        <div className="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                          <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                          </svg>
                        </div>
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
                                  {termResult.blockchain_published_at ? new Date(termResult.blockchain_published_at).toLocaleString() : 'N/A'}
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
                  <div className="mb-4 p-4 bg-white rounded-xl border border-gray-200">
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <span className="text-gray-600 font-medium">Verkle Root:</span>
                        <p className="font-mono text-xs mt-1 break-all text-slate-900 bg-gray-50 p-2 rounded">
                          {termData.verkle_root}
                        </p>
                      </div>
                      <div>
                        <span className="text-gray-600 font-medium">Generated:</span>
                        <p className="text-slate-900 mt-1">
                          {new Date(termData.generated_at).toLocaleString()}
                        </p>
                      </div>
                    </div>
                  </div>

                  {/* Courses */}
                  <div className="space-y-2">
                    <h4 className="font-semibold text-slate-900 mb-3 flex items-center gap-2">
                      <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                      </svg>
                      Courses
                    </h4>
                      {termData.receipt.revealed_courses.map((course, idx) => {
                        const courseKey = `${termId}-${course.course_id || idx}`;
                        const isCourseExpanded = expandedCourses.has(courseKey);

                      return (
                        <div
                          key={courseKey}
                          className="bg-white rounded-xl border border-gray-200 overflow-hidden"
                        >
                          {/* Course Header */}
                          <div
                            className="p-4 cursor-pointer hover:bg-gray-50/50 transition-colors"
                            onClick={() =>
                              toggleCourse(termId, course.course_id || idx.toString())
                            }
                          >
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${isCourseExpanded ? 'bg-blue-100' : 'bg-gray-100'}`}>
                                  <svg className={`w-4 h-4 ${isCourseExpanded ? 'text-blue-600' : 'text-gray-600'}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                                  </svg>
                                </div>
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
                                  className={`px-3 py-1 rounded-full text-sm font-semibold border ${
                                    course.grade === "A" || course.grade === "A+"
                                      ? "bg-green-100 text-green-800 border-green-200"
                                      : course.grade === "B" ||
                                        course.grade === "B+"
                                      ? "bg-blue-100 text-blue-800 border-blue-200"
                                      : course.grade === "C" ||
                                        course.grade === "C+"
                                      ? "bg-amber-100 text-amber-800 border-amber-200"
                                      : "bg-gray-100 text-gray-800 border-gray-200"
                                  }`}
                                >
                                  Grade: {course.grade}
                                </span>
                                <span className="text-sm text-gray-600">
                                  {course.credits} credits
                                </span>
                                <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={isCourseExpanded ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"} />
                                </svg>
                              </div>
                              </div>
                            </div>

                          {/* Course Details (Proof) */}
                          {isCourseExpanded && (
                            <div className="border-t border-gray-200 bg-gray-50 p-4">
                              <h5 className="text-sm font-semibold text-slate-900 mb-2 flex items-center gap-2">
                                <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                                </svg>
                                Cryptographic Proof
                              </h5>
                              <div className="bg-white p-3 rounded-lg border border-gray-200">
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
                              <div className="mt-3 flex flex-wrap items-center gap-2">
                                <span className="text-xs text-gray-600">
                                  Proof Type:
                                </span>
                                <span className="px-2 py-1 bg-indigo-100 text-indigo-800 text-xs rounded-lg font-medium">
                                  {termData.receipt.proof_type}
                                </span>
                                <span className="text-xs text-gray-600">
                                  Verification Path:
                                </span>
                                <span className="px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-lg font-medium">
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
  );
}
