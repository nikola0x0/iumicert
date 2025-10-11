"use client";

import { useState, useEffect } from "react";
import { apiService } from "@/lib/api";

export function VerifyJourneyTab() {
  const [students, setStudents] = useState<any[]>([]);
  const [selectedStudent, setSelectedStudent] = useState("");
  const [journeyData, setJourneyData] = useState<any>(null);
  const [termReceipts, setTermReceipts] = useState<any[]>([]);
  const [selectedTerm, setSelectedTerm] = useState<string | null>(null);
  const [selectedCourse, setSelectedCourse] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [verificationResult, setVerificationResult] = useState<any>(null);

  useEffect(() => {
    loadStudents();
  }, []);

  const loadStudents = async () => {
    try {
      const data = await apiService.getStudents();
      setStudents(data);
    } catch (error) {
      console.error("Failed to load students:", error);
    }
  };

  const loadJourney = async (studentId: string) => {
    setIsLoading(true);
    setJourneyData(null);
    setTermReceipts([]);
    setSelectedTerm(null);
    setSelectedCourse(null);
    setVerificationResult(null);

    try {
      // Load accumulated receipt for summary
      const accData = await apiService.getStudentAccumulatedReceipt(studentId);
      setJourneyData(accData.receipt);

      // Load term receipts for detailed course data
      const termData = await apiService.getStudentLatestReceipts(studentId);
      setTermReceipts(termData.receipts);
    } catch (error: any) {
      console.error("Failed to load journey:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleStudentSelect = (studentId: string) => {
    setSelectedStudent(studentId);
    if (studentId) {
      loadJourney(studentId);
    }
  };

  const verifyFullJourney = async () => {
    if (!termReceipts || termReceipts.length === 0) return;

    setIsLoading(true);
    setVerificationResult(null);

    try {
      // Build the complete receipt structure with all terms
      const termReceiptsMap: Record<string, any> = {};

      for (const termReceipt of termReceipts) {
        const courses =
          typeof termReceipt.RevealedCourses === "string"
            ? JSON.parse(termReceipt.RevealedCourses)
            : termReceipt.RevealedCourses;

        termReceiptsMap[termReceipt.TermID] = {
          verkle_root: termReceipt.VerkleRootHex,
          receipt: {
            course_proofs: termReceipt.VerkleProof,
            revealed_courses: courses,
          },
        };
      }

      const fullReceipt = {
        student_id: selectedStudent,
        term_receipts: termReceiptsMap,
      };

      // Call the IPA verification endpoint once for all terms
      const result = await apiService.verifyReceiptIPA(fullReceipt);

      setVerificationResult({
        full_journey: true,
        ...result,
      });
    } catch (error: any) {
      setVerificationResult({
        status: "error",
        error: error.message,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const verifyCourse = async (course: any, termId: string) => {
    setIsLoading(true);
    setVerificationResult(null);

    try {
      const termReceipt = termReceipts.find((r: any) => r.TermID === termId);

      if (!termReceipt) {
        throw new Error("Term receipt not found");
      }

      // Build the receipt structure expected by the backend
      const receiptForVerification = {
        student_id: selectedStudent,
        term_receipts: {
          [termId]: {
            verkle_root: termReceipt.VerkleRootHex,
            receipt: {
              course_proofs: termReceipt.VerkleProof,
              revealed_courses: termReceipt.RevealedCourses,
            },
          },
        },
      };

      const result = await apiService.verifyCourse({
        receipt: receiptForVerification,
        course_id: course.course_id,
        term_id: termId,
      });
      setVerificationResult(result);
    } catch (error: any) {
      setVerificationResult({
        verified: false,
        verification_error: error.message,
      });
    } finally {
      setIsLoading(false);
    }
  };

  // Get courses grouped by term from term receipts
  const getTermCourses = () => {
    const grouped: Record<string, any[]> = {};

    termReceipts.forEach((receipt: any) => {
      try {
        const courses =
          typeof receipt.RevealedCourses === "string"
            ? JSON.parse(receipt.RevealedCourses)
            : receipt.RevealedCourses;
        grouped[receipt.TermID] = courses;
      } catch (error) {
        console.error(`Failed to parse courses for ${receipt.TermID}:`, error);
        grouped[receipt.TermID] = [];
      }
    });

    return grouped;
  };

  const termCourses = getTermCourses();

  return (
    <div className="space-y-6">
      {/* Student Selection */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-2xl font-bold mb-4">Select Student</h2>
        <select
          value={selectedStudent}
          onChange={(e) => handleStudentSelect(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        >
          <option value="">-- Select a student --</option>
          {students.map((student) => (
            <option key={student.student_id} value={student.student_id}>
              {student.name} ({student.student_id})
            </option>
          ))}
        </select>
      </div>

      {/* Loading State */}
      {isLoading && !journeyData && (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading academic journey...</p>
        </div>
      )}

      {/* Journey Overview */}
      {journeyData && (
        <>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-start mb-6">
              <div>
                <h2 className="text-2xl font-bold mb-2">
                  Academic Journey Overview
                </h2>
                <p className="text-gray-600">Student: {selectedStudent}</p>
              </div>
              <div className="flex gap-3">
                <button
                  onClick={async () => {
                    try {
                      await apiService.downloadJourneyReceipt(selectedStudent);
                    } catch (error: any) {
                      alert(`Download failed: ${error.message}`);
                    }
                  }}
                  className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  📥 Download Receipt
                </button>
                <button
                  onClick={verifyFullJourney}
                  disabled={isLoading}
                  className="bg-green-600 text-white px-6 py-2 rounded-lg hover:bg-green-700 disabled:bg-gray-300 transition-colors"
                >
                  {isLoading ? "Verifying..." : "Verify Full Journey"}
                </button>
              </div>
            </div>

            {/* Summary Stats */}
            <div className="grid grid-cols-4 gap-4 mb-6">
              <div className="bg-blue-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Completed Terms</p>
                <p className="text-2xl font-bold text-blue-600">
                  {journeyData.CompletedTerms}
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Total Courses</p>
                <p className="text-2xl font-bold text-purple-600">
                  {journeyData.TotalCourses}
                </p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Total Credits</p>
                <p className="text-2xl font-bold text-orange-600">
                  {journeyData.TotalCredits}
                </p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">GPA</p>
                <p className="text-2xl font-bold text-green-600">
                  {journeyData.GPA.toFixed(2)}
                </p>
              </div>
            </div>

            {/* Terms List */}
            <div>
              <h3 className="text-lg font-semibold mb-4">Terms</h3>
              <div className="grid grid-cols-1 gap-3">
                {Object.keys(termCourses).map((termId) => {
                  const termReceipt = termReceipts.find(
                    (r: any) => r.TermID === termId
                  );
                  const isPublished = termReceipt?.BlockchainVerified || false;

                  return (
                    <button
                      key={termId}
                      onClick={() =>
                        setSelectedTerm(selectedTerm === termId ? null : termId)
                      }
                      className={`text-left p-4 rounded-lg border-2 transition-colors ${
                        selectedTerm === termId
                          ? "border-blue-500 bg-blue-50"
                          : "border-gray-200 hover:border-gray-300 bg-white"
                      }`}
                    >
                      <div className="flex justify-between items-center">
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <p className="font-semibold text-gray-900">
                              {termId}
                            </p>
                            {isPublished ? (
                              <span className="text-xs bg-green-100 text-green-700 px-2 py-1 rounded">
                                ✅ Published
                              </span>
                            ) : (
                              <span className="text-xs bg-yellow-100 text-yellow-700 px-2 py-1 rounded">
                                ⚠️ Unpublished
                              </span>
                            )}
                          </div>
                          <p className="text-sm text-gray-600">
                            {termCourses[termId].length} courses
                          </p>
                        </div>
                        <span className="text-gray-400">
                          {selectedTerm === termId ? "▼" : "▶"}
                        </span>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          </div>

          {/* Selected Term Courses */}
          {selectedTerm && termCourses[selectedTerm] && (
            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="text-xl font-bold mb-4">
                {selectedTerm} - Courses
              </h3>
              <div className="space-y-3">
                {termCourses[selectedTerm].map((course: any, idx: number) => (
                  <div
                    key={idx}
                    className={`p-4 rounded-lg border-2 transition-colors ${
                      selectedCourse?.course_id === course.course_id
                        ? "border-purple-500 bg-purple-50"
                        : "border-gray-200 hover:border-gray-300"
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <p className="font-semibold text-gray-900">
                          {course.course_name || course.course_id}
                        </p>
                        <p className="text-sm text-gray-600">
                          {course.course_id} • Grade: {course.grade} • Credits:{" "}
                          {course.credits}
                        </p>
                      </div>
                      <button
                        onClick={() => {
                          setSelectedCourse(course);
                          verifyCourse(course, selectedTerm);
                        }}
                        disabled={isLoading}
                        className="ml-4 bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 disabled:bg-gray-300 text-sm transition-colors"
                      >
                        {isLoading &&
                        selectedCourse?.course_id === course.course_id
                          ? "Verifying..."
                          : "Verify Course"}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Verification Results */}
          {verificationResult && (
            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="text-xl font-bold mb-4">Verification Results</h3>

              {/* Full Journey Verification (IPA) */}
              {verificationResult.full_journey && (
                <div className="space-y-4">
                  {/* Overall Status */}
                  <div
                    className={`p-4 rounded-lg ${
                      verificationResult.status === "success"
                        ? "bg-green-50 border border-green-200"
                        : verificationResult.status === "partial_failure"
                        ? "bg-yellow-50 border border-yellow-200"
                        : "bg-red-50 border border-red-200"
                    }`}
                  >
                    <p className="font-medium mb-2">
                      {verificationResult.status === "success" ? "✅" : verificationResult.status === "partial_failure" ? "⚠️" : "❌"}{" "}
                      Full Journey IPA Verification: {verificationResult.status}
                    </p>
                    <div className="grid grid-cols-3 gap-4 text-sm mb-3">
                      <div>
                        <p className="text-gray-600">Total Courses</p>
                        <p className="font-bold">{verificationResult.total_courses}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Verified</p>
                        <p className="font-bold text-green-600">{verificationResult.verified_courses}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Failed</p>
                        <p className="font-bold text-red-600">{verificationResult.failed_courses}</p>
                      </div>
                    </div>
                    {verificationResult.computation_note && (
                      <p className="text-xs text-gray-500 italic">{verificationResult.computation_note}</p>
                    )}
                  </div>

                  {/* Per-Term Results */}
                  {verificationResult.term_results && (
                    <div className="space-y-3">
                      <h4 className="font-semibold">Results by Term:</h4>
                      {Object.entries(verificationResult.term_results).map(([termId, termData]: [string, any]) => (
                        <div key={termId} className="border border-gray-200 rounded-lg p-4">
                          <div className="flex justify-between items-start mb-2">
                            <div>
                              <h5 className="font-medium">{termId}</h5>
                              {termData.blockchain_verified ? (
                                <div className="text-xs text-green-600 mt-1">
                                  ⛓️ Blockchain verified
                                </div>
                              ) : (
                                <div className="text-xs text-red-600 mt-1">
                                  ⚠️ Not on blockchain
                                </div>
                              )}
                            </div>
                            <span className={`text-xs px-2 py-1 rounded ${termData.status === "completed" ? "bg-blue-100 text-blue-800" : termData.status === "error" ? "bg-red-100 text-red-800" : "bg-gray-100 text-gray-800"}`}>
                              {termData.status}
                            </span>
                          </div>
                          {termData.error && (
                            <div className="text-xs text-red-600 mb-2 p-2 bg-red-50 rounded">
                              {termData.error}
                            </div>
                          )}
                          <div className="text-sm mb-2">
                            <span className="text-green-600">✓ {termData.courses_verified || 0}</span>
                            {" / "}
                            <span className="text-red-600">✗ {termData.courses_failed || 0}</span>
                          </div>
                          {termData.course_results && (
                            <details className="text-sm">
                              <summary className="cursor-pointer text-blue-600 hover:text-blue-800">View course details</summary>
                              <div className="mt-2 space-y-1 pl-4">
                                {Object.entries(termData.course_results).map(([courseId, status]: [string, any]) => (
                                  <div key={courseId} className="flex justify-between">
                                    <span className="font-mono text-xs">{courseId}</span>
                                    <span className={status === "verified" ? "text-green-600" : "text-red-600"}>
                                      {status === "verified" ? "✓ verified" : `✗ ${status}`}
                                    </span>
                                  </div>
                                ))}
                              </div>
                            </details>
                          )}
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Failed Courses List */}
                  {verificationResult.failed_list && verificationResult.failed_list.length > 0 && (
                    <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                      <p className="font-medium text-red-800 mb-2">Failed Courses:</p>
                      <ul className="text-sm text-red-700 space-y-1">
                        {verificationResult.failed_list.map((failed: string, idx: number) => (
                          <li key={idx} className="font-mono">{failed}</li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              )}

              {/* Individual Course Verification */}
              {verificationResult.verified !== undefined &&
                !verificationResult.full_journey && (
                  <div
                    className={`p-4 rounded-lg ${
                      verificationResult.verified
                        ? "bg-green-50 border border-green-200"
                        : "bg-red-50 border border-red-200"
                    }`}
                  >
                    {verificationResult.verified ? (
                      <>
                        <p className="text-green-800 font-medium mb-4">
                          ✅ Course Verified Successfully!
                        </p>
                        {verificationResult.course && (
                          <div className="space-y-2 text-sm mb-4">
                            <p>
                              <strong>Course:</strong>{" "}
                              {verificationResult.course.course_name} (
                              {verificationResult.course.course_id})
                            </p>
                            <p>
                              <strong>Grade:</strong>{" "}
                              {verificationResult.course.grade}
                            </p>
                            <p>
                              <strong>Credits:</strong>{" "}
                              {verificationResult.course.credits}
                            </p>
                            <p>
                              <strong>Term:</strong>{" "}
                              {verificationResult.term_id}
                            </p>
                          </div>
                        )}
                        {verificationResult.verification_details && (
                          <div className="pt-4 border-t border-green-200">
                            <p className="font-medium mb-2">
                              Verification Details:
                            </p>
                            <ul className="space-y-1 text-sm">
                              <li>
                                {verificationResult.verification_details
                                  .ipa_verified
                                  ? "✅"
                                  : "❌"}{" "}
                                IPA Cryptographic Proof
                              </li>
                              <li>
                                {verificationResult.verification_details
                                  .state_diff_verified
                                  ? "✅"
                                  : "❌"}{" "}
                                State Diff Verification
                              </li>
                              <li>
                                {verificationResult.verification_details
                                  .blockchain_anchored
                                  ? "✅"
                                  : "❌"}{" "}
                                Blockchain Anchored
                              </li>
                            </ul>
                          </div>
                        )}
                        {verificationResult.blockchain_info &&
                          verificationResult.blockchain_info.tx_hash && (
                            <div className="pt-4 border-t border-green-200">
                              <p className="font-medium mb-2">
                                Blockchain Anchor:
                              </p>
                              <div className="space-y-2 text-sm">
                                <div>
                                  <a
                                    href={`https://sepolia.etherscan.io/tx/${verificationResult.blockchain_info.tx_hash}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-blue-600 hover:text-blue-800"
                                  >
                                    🔗 View Transaction on Etherscan
                                  </a>
                                </div>
                                {verificationResult.blockchain_info
                                  .publisher_address && (
                                  <div className="text-xs text-gray-600">
                                    <span className="font-medium">
                                      Publisher:
                                    </span>{" "}
                                    <a
                                      href={`https://sepolia.etherscan.io/address/${verificationResult.blockchain_info.publisher_address}`}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      className="text-blue-600 hover:text-blue-800 font-mono"
                                    >
                                      {
                                        verificationResult.blockchain_info
                                          .publisher_address
                                      }
                                    </a>
                                  </div>
                                )}
                                <div className="text-xs text-gray-600">
                                  Published:{" "}
                                  {new Date(
                                    verificationResult.blockchain_info
                                      .published_at * 1000
                                  ).toLocaleString()}
                                </div>
                              </div>
                            </div>
                          )}
                      </>
                    ) : (
                      <div>
                        <p className="text-red-800 font-medium mb-2">
                          ❌ Verification Failed
                        </p>
                        <p className="text-sm text-red-700">
                          {verificationResult.verification_error?.includes(
                            "does not exist on blockchain"
                          )
                            ? "⚠️ This term has not been published to the blockchain yet. Please go to the 'Publish Terms' tab to publish it first."
                            : verificationResult.verification_error ||
                              "Unknown error"}
                        </p>
                      </div>
                    )}
                  </div>
                )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
