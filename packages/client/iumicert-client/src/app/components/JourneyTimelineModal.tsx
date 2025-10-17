"use client";

import React from "react";
import { X, Calendar, Award, BookOpen, Shield, CheckCircle, AlertTriangle } from "lucide-react";

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

interface JourneyTimelineModalProps {
  isOpen: boolean;
  onClose: () => void;
  receipt: JourneyReceipt;
  verificationResults?: Record<string, any>;
}

export default function JourneyTimelineModal({
  isOpen,
  onClose,
  receipt,
  verificationResults = {},
}: JourneyTimelineModalProps) {
  if (!isOpen) return null;

  // Sort terms chronologically
  const sortedTerms = Object.entries(receipt.term_receipts).sort(
    ([termIdA], [termIdB]) => {
      const extractOrder = (id: string) => {
        const yearMatch = id.match(/(\d{4})/);
        const year = yearMatch ? parseInt(yearMatch[1]) : 0;
        const semester = id.includes("Semester_1")
          ? 1
          : id.includes("Semester_2")
          ? 2
          : 3;
        return year * 10 + semester;
      };
      return extractOrder(termIdA) - extractOrder(termIdB);
    }
  );

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
      <div className="relative w-full h-full bg-gradient-to-br from-slate-900 via-purple-900 to-blue-900 overflow-hidden">
        {/* Animated Background */}
        <div className="absolute inset-0 opacity-30">
          <div className="absolute top-20 left-20 w-96 h-96 bg-blue-500 rounded-full blur-3xl animate-pulse"></div>
          <div className="absolute bottom-20 right-20 w-96 h-96 bg-purple-500 rounded-full blur-3xl animate-pulse delay-1000"></div>
        </div>

        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-6 right-6 z-50 w-12 h-12 bg-white/10 hover:bg-white/20 rounded-full flex items-center justify-center transition-colors backdrop-blur-sm border border-white/20"
        >
          <X className="w-6 h-6 text-white" />
        </button>

        {/* Header */}
        <div className="relative z-10 p-8 border-b border-white/10">
          <div className="max-w-7xl mx-auto">
            <div className="flex items-center gap-4 mb-2">
              <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
                <Calendar className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-white font-space-grotesk">
                  Academic Journey Timeline
                </h1>
                <p className="text-purple-200 font-inter">
                  Student ID: <span className="font-mono">{receipt.student_id}</span>
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Timeline Content */}
        <div className="relative z-10 h-[calc(100%-120px)] overflow-y-auto">
          <div className="max-w-5xl mx-auto py-12 px-8">
            <div className="relative">
              {/* Vertical Timeline Line */}
              <div className="absolute left-[2.5rem] top-0 bottom-0 w-1 bg-gradient-to-b from-blue-500 via-purple-500 to-blue-500"></div>

              {/* Timeline Items */}
              <div className="space-y-12">
                {sortedTerms.map(([termId, termData], index) => {
                  const termResult = verificationResults[termId];
                  const isBlockchainVerified = termResult?.blockchain_verified === true;
                  const totalCoursesInReceipt = Object.values(receipt.term_receipts).reduce(
                    (sum, term) => sum + term.total_courses,
                    0
                  );

                  return (
                    <div key={termId} className="relative pl-24">
                      {/* Timeline Dot */}
                      <div className="absolute left-0 top-0 w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center shadow-2xl shadow-blue-500/50 border-4 border-slate-900">
                        <span className="text-2xl font-bold text-white font-space-grotesk">
                          {index + 1}
                        </span>
                      </div>

                      {/* Term Card */}
                      <div className="glass-effect rounded-2xl p-6 ml-4 hover:scale-[1.02] transition-transform duration-300">
                        {/* Term Header */}
                        <div className="flex items-center justify-between mb-4">
                          <div className="flex items-center gap-3">
                            <BookOpen className="w-6 h-6 text-blue-400" />
                            <h3 className="text-2xl font-bold text-white font-space-grotesk">
                              {termId.replace(/_/g, " ")}
                            </h3>
                          </div>
                          {termResult && (
                            <>
                              {isBlockchainVerified ? (
                                <div className="px-4 py-2 bg-green-500/20 text-green-300 rounded-full text-sm font-semibold border border-green-500/30 flex items-center gap-2 font-inter">
                                  <CheckCircle className="w-4 h-4" />
                                  Verified
                                </div>
                              ) : (
                                <div className="px-4 py-2 bg-amber-500/20 text-amber-300 rounded-full text-sm font-semibold border border-amber-500/30 flex items-center gap-2 font-inter">
                                  <AlertTriangle className="w-4 h-4" />
                                  Not Published
                                </div>
                              )}
                            </>
                          )}
                        </div>

                        {/* Term Stats */}
                        <div className="grid grid-cols-3 gap-4 mb-6">
                          <div className="bg-blue-500/20 rounded-xl p-4 border border-blue-500/30">
                            <div className="text-3xl font-bold text-blue-300 font-space-grotesk">
                              {termData.total_courses}
                            </div>
                            <div className="text-sm text-purple-200 font-inter">
                              Courses
                            </div>
                          </div>
                          <div className="bg-purple-500/20 rounded-xl p-4 border border-purple-500/30">
                            <div className="text-3xl font-bold text-purple-300 font-space-grotesk">
                              {termData.receipt.revealed_courses.reduce(
                                (sum, c) => sum + c.credits,
                                0
                              )}
                            </div>
                            <div className="text-sm text-purple-200 font-inter">
                              Credits
                            </div>
                          </div>
                          <div className="bg-green-500/20 rounded-xl p-4 border border-green-500/30">
                            <div className="text-3xl font-bold text-green-300 font-space-grotesk">
                              {(
                                termData.receipt.revealed_courses.reduce(
                                  (sum, c) => {
                                    const gradePoints: Record<string, number> = {
                                      "A+": 4.0,
                                      A: 4.0,
                                      "B+": 3.5,
                                      B: 3.0,
                                      "C+": 2.5,
                                      C: 2.0,
                                      "D+": 1.5,
                                      D: 1.0,
                                      F: 0.0,
                                    };
                                    return sum + (gradePoints[c.grade] || 0);
                                  },
                                  0
                                ) / termData.receipt.revealed_courses.length
                              ).toFixed(2)}
                            </div>
                            <div className="text-sm text-purple-200 font-inter">
                              GPA
                            </div>
                          </div>
                        </div>

                        {/* Courses Grid */}
                        <div className="space-y-3">
                          <h4 className="text-sm font-semibold text-white/60 uppercase tracking-wider font-inter">
                            Courses
                          </h4>
                          <div className="grid grid-cols-1 gap-2">
                            {termData.receipt.revealed_courses.map(
                              (course, idx) => (
                                <div
                                  key={idx}
                                  className="bg-white/5 rounded-lg p-3 border border-white/10 hover:bg-white/10 transition-colors"
                                >
                                  <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-3">
                                      <Award className="w-4 h-4 text-blue-400" />
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
                                        {course.grade}
                                      </span>
                                      <span className="text-sm text-purple-200 font-inter">
                                        {course.credits} cr
                                      </span>
                                    </div>
                                  </div>
                                </div>
                              )
                            )}
                          </div>
                        </div>

                        {/* Blockchain Verification Details */}
                        {isBlockchainVerified && termResult && (
                          <div className="mt-6 space-y-3">
                            <h4 className="text-sm font-semibold text-white/60 uppercase tracking-wider font-inter flex items-center gap-2">
                              <Shield className="w-4 h-4" />
                              Blockchain Verification
                            </h4>
                            <div className="bg-white/5 rounded-lg p-4 border border-white/10 space-y-3">
                              {/* Verkle Root */}
                              <div>
                                <p className="text-xs text-white/50 font-inter mb-1">
                                  Verkle Root (Cryptographic Commitment)
                                </p>
                                <p className="text-xs text-purple-200 font-mono break-all bg-black/30 p-2 rounded">
                                  {termResult.verkle_root}
                                </p>
                              </div>

                              {/* Blockchain Transaction */}
                              {termResult.blockchain_tx_hash && (
                                <div>
                                  <p className="text-xs text-white/50 font-inter mb-1">
                                    Transaction Hash
                                  </p>
                                  <a
                                    href={`https://sepolia.etherscan.io/tx/${termResult.blockchain_tx_hash}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-xs text-blue-300 font-mono break-all hover:text-blue-200 underline bg-black/30 p-2 rounded block"
                                  >
                                    {termResult.blockchain_tx_hash}
                                  </a>
                                </div>
                              )}

                              {/* Block Number */}
                              {termResult.blockchain_block && (
                                <div>
                                  <p className="text-xs text-white/50 font-inter mb-1">
                                    Block Number
                                  </p>
                                  <a
                                    href={`https://sepolia.etherscan.io/block/${termResult.blockchain_block}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-xs text-blue-300 font-mono hover:text-blue-200 underline"
                                  >
                                    #{termResult.blockchain_block}
                                  </a>
                                </div>
                              )}

                              {/* Published Timestamp */}
                              {termResult.blockchain_published_at && (
                                <div>
                                  <p className="text-xs text-white/50 font-inter mb-1">
                                    Published On-Chain
                                  </p>
                                  <p className="text-sm text-purple-200 font-inter">
                                    {new Date(
                                      parseInt(termResult.blockchain_published_at) * 1000
                                    ).toLocaleString()}
                                  </p>
                                </div>
                              )}

                              {/* Verification Status */}
                              <div className="flex items-start gap-2 pt-2 border-t border-white/10">
                                <CheckCircle className="w-4 h-4 text-green-400 mt-0.5 flex-shrink-0" />
                                <div className="flex-1">
                                  <p className="text-xs text-green-300 font-inter font-semibold">
                                    Cryptographically Verified
                                  </p>
                                  <p className="text-xs text-green-200/70 font-inter mt-1">
                                    {termResult.courses_verified || 0} course{termResult.courses_verified !== 1 ? 's' : ''} verified using IPA (Inner Product Argument) with 32-byte proofs against Ethereum Sepolia testnet
                                  </p>
                                </div>
                              </div>
                            </div>
                          </div>
                        )}

                        {/* Term Metadata */}
                        <div className="mt-4 pt-4 border-t border-white/10">
                          <p className="text-xs text-purple-200 font-inter">
                            Generated:{" "}
                            {new Date(termData.generated_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* End Marker */}
              <div className="relative pl-24 mt-12">
                <div className="absolute left-0 top-0 w-20 h-20 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full flex items-center justify-center shadow-2xl shadow-green-500/50 border-4 border-slate-900">
                  <CheckCircle className="w-10 h-10 text-white" />
                </div>
                <div className="glass-effect rounded-2xl p-6 ml-4">
                  <h3 className="text-2xl font-bold text-white mb-2 font-space-grotesk">
                    Journey Complete
                  </h3>
                  <p className="text-purple-200 font-inter">
                    {sortedTerms.length} terms · Total{" "}
                    {Object.values(receipt.term_receipts).reduce(
                      (sum, term) => sum + term.total_courses,
                      0
                    )}{" "}
                    courses completed
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
