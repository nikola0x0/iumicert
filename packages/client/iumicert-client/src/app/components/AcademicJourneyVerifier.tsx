"use client";

import React, { useState } from "react";
import { apiService } from "@/lib/api";
import FileUploaderWrapper from "./FileUploaderWrapper";
import JourneyTimelineModal from "./JourneyTimelineModal";
import { Shield, ChevronDown, ChevronRight, Lock, GitBranch } from "lucide-react";

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
  verkle_root?: string;
  details?: string;
  error?: string;
}

export default function AcademicJourneyVerifier() {
  const [receipt, setReceipt] = useState<JourneyReceipt | null>(null);
  const [expandedYears, setExpandedYears] = useState<Set<string>>(new Set());
  const [expandedTerms, setExpandedTerms] = useState<Set<string>>(new Set());
  const [expandedCourses, setExpandedCourses] = useState<Set<string>>(
    new Set()
  );
  const [verificationResults, setVerificationResults] = useState<
    Record<string, VerificationResult>
  >({});
  const [isVerifying, setIsVerifying] = useState(false);
  const [error, setError] = useState<string>("");
  const [isTimelineModalOpen, setIsTimelineModalOpen] = useState(false);

  // Handle file upload
  const handleFileUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const json = JSON.parse(e.target?.result as string);
        setReceipt(json);
        setError("");
        // Start with everything collapsed for journey view
        setExpandedYears(new Set());
        setExpandedTerms(new Set());
        setExpandedCourses(new Set());
      } catch (err) {
        setError("Invalid JSON file. Please upload a valid receipt.");
      }
    };
    reader.readAsText(file);
  };

  // Toggle year expansion
  const toggleYear = (year: string) => {
    setExpandedYears((prev) => {
      const next = new Set(prev);
      if (next.has(year)) {
        next.delete(year);
      } else {
        next.add(year);
      }
      return next;
    });
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

      // Calculate total courses in the entire receipt
      const totalCoursesInReceipt = Object.values(receipt.term_receipts).reduce(
        (sum, term) => sum + term.total_courses,
        0
      );

      // Store overall result
      setVerificationResults({
        overall: {
          verified: result.status === "success",
          ipa_verified: result.verified_courses === result.total_courses,
          blockchain_anchored: true,
          details: `Cryptographically verified: ${result.verified_courses} out of ${totalCoursesInReceipt} courses in this receipt`,
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
      <div className="h-full flex items-center justify-center p-4">
        <div className="glass-effect rounded-xl p-8 max-w-lg w-full">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <Shield className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-white font-space-grotesk mb-2">
              Verify Academic Journey
            </h1>
            <p className="text-purple-200 font-inter text-sm">
              Upload an academic journey receipt to verify credentials
            </p>
          </div>

          <FileUploaderWrapper
            handleChange={handleFileUpload}
            name="receipt"
            types={["JSON"]}
          >
            <div className="border-2 border-dashed border-white/20 rounded-xl p-12 text-center hover:border-blue-400 hover:bg-blue-500/10 transition-all duration-200 cursor-pointer">
              <div className="w-16 h-16 bg-blue-500/20 rounded-2xl flex items-center justify-center mb-4 mx-auto">
                <svg
                  className="w-8 h-8 text-blue-400"
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
              <span className="text-xl font-semibold text-white mb-2 block font-space-grotesk">
                Click or drag to upload receipt
              </span>
              <span className="text-sm text-purple-200 font-inter">
                JSON files only
              </span>
            </div>
          </FileUploaderWrapper>

          {error && (
            <div className="mt-6 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-200 flex items-start gap-3 font-inter">
              <svg
                className="w-5 h-5 mt-0.5 flex-shrink-0"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
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
    <div className="h-full overflow-y-auto p-6 space-y-6">
      {/* Header */}
      <div className="glass-effect rounded-xl p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-white font-space-grotesk">
              Academic Journey
            </h1>
            <p className="text-purple-200 mt-1 font-inter">
              Student ID:{" "}
              <span className="font-mono font-semibold text-white">
                {receipt.student_id}
              </span>
            </p>
          </div>
          <button
            onClick={() => {
              setReceipt(null);
              setVerificationResults({});
            }}
            className="px-5 py-2.5 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors font-medium font-inter"
          >
            Upload New Receipt
          </button>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-wrap items-center gap-4">
          <button
            onClick={verifyJourney}
            disabled={isVerifying}
            className={`px-6 py-3 rounded-xl font-semibold transition-all duration-200 font-inter ${
              isVerifying
                ? "bg-gray-400 cursor-not-allowed"
                : "bg-gradient-to-r from-blue-500 to-blue-600 hover:shadow-lg hover:shadow-blue-500/30"
            } text-white`}
          >
            {isVerifying ? (
              <span className="flex items-center gap-2">
                <svg
                  className="w-5 h-5 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                Verifying...
              </span>
            ) : (
              "Verify Entire Journey"
            )}
          </button>

          {overallResult && (
            <>
              <button
                onClick={() => setIsTimelineModalOpen(true)}
                className="px-6 py-3 rounded-xl font-semibold transition-all duration-200 font-inter bg-gradient-to-r from-purple-500 to-indigo-600 hover:shadow-lg hover:shadow-purple-500/30 text-white flex items-center gap-2"
              >
                <GitBranch className="w-5 h-5" />
                View Journey Timeline
              </button>

              <div
                className={`px-4 py-2 rounded-xl font-semibold border font-inter ${
                  overallResult.verified
                    ? "bg-green-500/20 text-green-300 border-green-500/30"
                    : "bg-red-500/20 text-red-300 border-red-500/30"
                }`}
              >
                {overallResult.verified ? "Verified" : "Verification Failed"}
              </div>
            </>
          )}
        </div>

        {overallResult?.details && (
          <div className="mt-4 p-4 bg-blue-500/20 border border-blue-500/30 rounded-xl">
            <p className="text-sm text-blue-200 font-inter">
              {overallResult.details}
            </p>
          </div>
        )}

        {error && (
          <div className="mt-4 p-4 bg-red-500/20 border border-red-500/30 rounded-xl flex items-start gap-3">
            <svg
              className="w-5 h-5 mt-0.5 flex-shrink-0 text-red-300"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                clipRule="evenodd"
              />
            </svg>
            <p className="text-sm text-red-200 font-inter">{error}</p>
          </div>
        )}
      </div>

      {/* Journey Summary */}
      <div className="glass-effect rounded-xl p-6">
        <h2 className="text-xl font-bold text-white mb-4 font-space-grotesk">
          Journey Summary
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-blue-500/20 p-5 rounded-xl border border-blue-500/30">
            <div className="text-2xl font-bold text-blue-300 font-space-grotesk">
              {Object.keys(receipt.term_receipts).length}
            </div>
            <div className="text-sm text-purple-200 mt-1 font-inter">
              Terms Completed
            </div>
          </div>
          <div className="bg-green-500/20 p-5 rounded-xl border border-green-500/30">
            <div className="text-2xl font-bold text-green-300 font-space-grotesk">
              {Object.values(receipt.term_receipts).reduce(
                (sum, term) => sum + term.total_courses,
                0
              )}
            </div>
            <div className="text-sm text-purple-200 mt-1 font-inter">
              Total Courses
            </div>
          </div>
          <div className="bg-purple-500/20 p-5 rounded-xl border border-purple-500/30">
            <div className="text-2xl font-bold text-purple-300 font-space-grotesk">
              {receipt.blockchain_ready ? "Valid" : "Invalid"}
            </div>
            <div className="text-sm text-purple-200 mt-1 font-inter">
              Receipt Format
            </div>
          </div>
          <div className="bg-indigo-500/20 p-5 rounded-xl border border-indigo-500/30">
            <div className="text-2xl font-bold text-indigo-300 font-space-grotesk">
              {receipt.receipt_type.selective_disclosure ? "Selective" : "Full"}
            </div>
            <div className="text-sm text-purple-200 mt-1 font-inter">
              Disclosure Type
            </div>
          </div>
        </div>
      </div>

      {/* Journey Timeline */}
      <div className="space-y-6">
        <h2 className="text-2xl font-bold text-white font-space-grotesk">
          Academic Journey Timeline
        </h2>

        {(() => {
          // Group terms by year
          const termsByYear: Record<string, Array<[string, TermReceipt]>> = {};

          Object.entries(receipt.term_receipts).forEach(([termId, termData]) => {
            const yearMatch = termId.match(/(\d{4})/);
            const year = yearMatch ? yearMatch[1] : "Unknown";
            if (!termsByYear[year]) {
              termsByYear[year] = [];
            }
            termsByYear[year].push([termId, termData]);
          });

          // Sort years
          const sortedYears = Object.keys(termsByYear).sort();

          return sortedYears.map((year) => {
            const isYearExpanded = expandedYears.has(year);
            const yearTerms = termsByYear[year].sort(([termIdA], [termIdB]) => {
              const order = (id: string) =>
                id.includes("Semester_1") ? 1 : id.includes("Semester_2") ? 2 : 3;
              return order(termIdA) - order(termIdB);
            });

            // Calculate year stats
            const totalTerms = yearTerms.length;
            const totalCourses = yearTerms.reduce(
              (sum, [, term]) => sum + term.total_courses,
              0
            );

            // Check if all terms in this year are blockchain verified
            const verifiedTermsCount = yearTerms.filter(([termId]) => {
              const termResult = verificationResults[termId];
              return termResult?.blockchain_verified === true;
            }).length;
            const hasUnverifiedTerms = verifiedTermsCount < totalTerms;

            return (
              <div key={year} className="relative">
                {/* Year Card */}
                <div className="glass-effect rounded-xl overflow-hidden">
                  {/* Year Header */}
                  <div
                    className="p-6 cursor-pointer hover:bg-white/5 transition-colors"
                    onClick={() => toggleYear(year)}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div>
                          <h3 className="text-2xl font-bold text-white font-space-grotesk">
                            {year}
                          </h3>
                          <p className="text-sm text-purple-200 font-inter">
                            {totalTerms} {totalTerms === 1 ? "term" : "terms"} · {totalCourses} courses
                          </p>
                        </div>
                        {verificationResults.overall && hasUnverifiedTerms && (
                          <div className="px-3 py-1.5 bg-amber-500/20 text-amber-300 rounded-full text-xs font-semibold border border-amber-500/30 font-inter" title={`${verifiedTermsCount} of ${totalTerms} terms are blockchain-published`}>
                            ⚠ {verifiedTermsCount}/{totalTerms} terms on blockchain
                          </div>
                        )}
                      </div>
                      <ChevronDown
                        className={`w-6 h-6 text-white/60 transition-transform duration-200 ${
                          isYearExpanded ? "" : "-rotate-90"
                        }`}
                      />
                    </div>
                  </div>

                  {/* Year Content - Terms */}
                  {isYearExpanded && (
                    <div className="border-t border-white/10 bg-black/20 p-6 space-y-4">
                      {yearTerms.map(([termId, termData]) => {
                        const isTermExpanded = expandedTerms.has(termId);
                        const termResult = verificationResults[termId];
                        const isBlockchainVerified = termResult?.blockchain_verified === true;
                        const hasVerificationResult = termResult && termResult.status === "completed";

                        return (
                          <div
                            key={termId}
                            className="bg-white/5 rounded-xl overflow-hidden border border-white/10"
                          >
                            {/* Term Header */}
                            <div
                              className="p-5 cursor-pointer hover:bg-white/5 transition-colors"
                              onClick={() => toggleTerm(termId)}
                            >
                              <div className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                  <div
                                    className={`w-10 h-10 rounded-lg flex items-center justify-center ${
                                      isTermExpanded
                                        ? "bg-blue-500/30"
                                        : "bg-white/10"
                                    }`}
                                  >
                                    {isTermExpanded ? (
                                      <ChevronDown className="w-5 h-5 text-blue-300" />
                                    ) : (
                                      <ChevronRight className="w-5 h-5 text-white/60" />
                                    )}
                                  </div>
                                  <div>
                                    <h4 className="text-lg font-bold text-white font-space-grotesk">
                                      {termId.replace(/_/g, " ")}
                                    </h4>
                                    <p className="text-sm text-purple-200 font-inter">
                                      {termData.revealed_courses} of{" "}
                                      {termData.total_courses} courses
                                    </p>
                                  </div>
                                </div>
                                <div className="flex items-center gap-3">
                                  {hasVerificationResult && (
                                    <div
                                      className={`px-3 py-1.5 rounded-full text-xs font-semibold border font-inter ${
                                        isBlockchainVerified
                                          ? "bg-green-500/20 text-green-300 border-green-500/30"
                                          : "bg-amber-500/20 text-amber-300 border-amber-500/30"
                                      }`}
                                      title={isBlockchainVerified
                                        ? `${termResult.courses_verified || 0} courses in this term verified against blockchain`
                                        : "This term's root is not published on blockchain"}
                                    >
                                      {isBlockchainVerified
                                        ? `✓ ${termResult.courses_verified || 0}/${termData.total_courses} courses verified`
                                        : "⚠ Not Published"}
                                    </div>
                                  )}
                                </div>
                              </div>
                            </div>

                            {/* Term Content */}
                            {isTermExpanded && (
                  <div className="border-t border-white/10 bg-black/20 p-6">
                    {/* Verification Status */}
                    {termResult && termResult.blockchain_verified && (
                      <div className="mb-4 p-4 bg-green-500/20 rounded-xl border border-green-500/30">
                        <div className="flex items-start gap-3">
                          <div className="w-10 h-10 bg-green-500/30 rounded-xl flex items-center justify-center flex-shrink-0">
                            <Shield className="w-5 h-5 text-green-300" />
                          </div>
                          <div className="flex-1">
                            <h4 className="font-semibold text-green-300 mb-2 font-space-grotesk">
                              Blockchain Verified
                            </h4>
                            <div className="space-y-2 text-sm font-inter">
                              <div>
                                <span className="text-green-200">
                                  Verkle Root:
                                </span>
                                <span className="ml-2 font-mono text-xs text-green-300 break-all">
                                  {termResult.verkle_root}
                                </span>
                              </div>
                              {termResult.blockchain_tx_hash && (
                                <div>
                                  <span className="text-green-200">
                                    Transaction:
                                  </span>
                                  <a
                                    href={`https://sepolia.etherscan.io/tx/${termResult.blockchain_tx_hash}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="ml-2 font-mono text-xs text-blue-300 hover:text-blue-200 underline break-all"
                                  >
                                    {termResult.blockchain_tx_hash.slice(0, 10)}...{termResult.blockchain_tx_hash.slice(-8)}
                                  </a>
                                </div>
                              )}
                              {termResult.blockchain_block && (
                                <div>
                                  <span className="text-green-200">
                                    Block:
                                  </span>
                                  <a
                                    href={`https://sepolia.etherscan.io/block/${termResult.blockchain_block}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="ml-2 text-blue-300 hover:text-blue-200 underline"
                                  >
                                    #{termResult.blockchain_block}
                                  </a>
                                </div>
                              )}
                              {termResult.blockchain_published_at && (
                                <div>
                                  <span className="text-green-200">
                                    Published:
                                  </span>
                                  <span className="ml-2 text-green-300">
                                    {new Date(parseInt(termResult.blockchain_published_at) * 1000).toLocaleString()}
                                  </span>
                                </div>
                              )}
                              <p className="text-green-200/80 text-xs mt-2">
                                ✓ Cryptographically verified using IPA (Inner Product Argument) with 32-byte proofs against Ethereum Sepolia testnet
                              </p>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}

                    {/* Term Metadata */}
                    <div className="mb-4 p-4 bg-white/5 rounded-xl border border-white/10">
                      <div className="grid grid-cols-2 gap-4 text-sm font-inter">
                        <div>
                          <span className="text-purple-200 font-medium">
                            Verkle Root:
                          </span>
                          <p className="font-mono text-xs mt-1 break-all text-white bg-black/30 p-2 rounded">
                            {termData.verkle_root}
                          </p>
                        </div>
                        <div>
                          <span className="text-purple-200 font-medium">
                            Receipt Issued:
                          </span>
                          <p className="text-white mt-1">
                            {new Date(termData.generated_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>

                    {/* Courses */}
                    <div className="space-y-2">
                      <h4 className="font-semibold text-white mb-3 flex items-center gap-2 font-space-grotesk">
                        <svg
                          className="w-5 h-5 text-blue-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                          />
                        </svg>
                        Courses
                      </h4>
                      {termData.receipt.revealed_courses.map((course, idx) => {
                        const courseKey = `${termId}-${
                          course.course_id || idx
                        }`;
                        const isCourseExpanded = expandedCourses.has(courseKey);

                        return (
                          <div
                            key={courseKey}
                            className="bg-white/5 rounded-xl border border-white/10 overflow-hidden"
                          >
                            {/* Course Header */}
                            <div
                              className="p-4 cursor-pointer hover:bg-white/5 transition-colors"
                              onClick={() =>
                                toggleCourse(
                                  termId,
                                  course.course_id || idx.toString()
                                )
                              }
                            >
                              <div className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                  <div
                                    className={`w-8 h-8 rounded-lg flex items-center justify-center ${
                                      isCourseExpanded
                                        ? "bg-blue-500/30"
                                        : "bg-white/10"
                                    }`}
                                  >
                                    {isCourseExpanded ? (
                                      <ChevronDown className="w-4 h-4 text-blue-300" />
                                    ) : (
                                      <ChevronRight className="w-4 h-4 text-white/60" />
                                    )}
                                  </div>
                                  <div>
                                    <p className="font-semibold text-white font-inter">
                                      {course.course_id || `Course ${idx + 1}`}
                                    </p>
                                    <p className="text-sm text-purple-200 font-inter">
                                      {course.course_name || "N/A"}
                                    </p>
                                  </div>
                                </div>
                                <div className="flex items-center gap-3">
                                  <span
                                    className={`px-3 py-1 rounded-full text-sm font-semibold border font-inter ${
                                      course.grade === "A" ||
                                      course.grade === "A+"
                                        ? "bg-green-500/20 text-green-300 border-green-500/30"
                                        : course.grade === "B" ||
                                          course.grade === "B+"
                                        ? "bg-blue-500/20 text-blue-300 border-blue-500/30"
                                        : course.grade === "C" ||
                                          course.grade === "C+"
                                        ? "bg-amber-500/20 text-amber-300 border-amber-500/30"
                                        : "bg-gray-500/20 text-gray-300 border-gray-500/30"
                                    }`}
                                  >
                                    Grade: {course.grade}
                                  </span>
                                  <span className="text-sm text-purple-200 font-inter">
                                    {course.credits} credits
                                  </span>
                                </div>
                              </div>
                            </div>

                            {/* Course Details (Proof) */}
                            {isCourseExpanded && (
                              <div className="border-t border-white/10 bg-black/20 p-4">
                                <h5 className="text-sm font-semibold text-white mb-2 flex items-center gap-2 font-space-grotesk">
                                  <Lock className="w-4 h-4 text-blue-400" />
                                  Cryptographic Proof
                                </h5>
                                <div className="bg-black/30 p-3 rounded-lg border border-white/10">
                                  <pre className="text-xs font-mono text-purple-200 overflow-x-auto">
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
                                  <span className="text-xs text-purple-200 font-inter">
                                    Proof Type:
                                  </span>
                                  <span className="px-2 py-1 bg-indigo-500/20 text-indigo-300 text-xs rounded-lg font-medium border border-indigo-500/30 font-inter">
                                    {termData.receipt.proof_type}
                                  </span>
                                  <span className="text-xs text-purple-200 font-inter">
                                    Verification Path:
                                  </span>
                                  <span className="px-2 py-1 bg-purple-500/20 text-purple-300 text-xs rounded-lg font-medium border border-purple-500/30 font-inter">
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
                  )}
                </div>
              </div>
            );
          });
        })()}
      </div>

      {/* Journey Timeline Modal */}
      <JourneyTimelineModal
        isOpen={isTimelineModalOpen}
        onClose={() => setIsTimelineModalOpen(false)}
        receipt={receipt}
        verificationResults={verificationResults}
      />
    </div>
  );
}
