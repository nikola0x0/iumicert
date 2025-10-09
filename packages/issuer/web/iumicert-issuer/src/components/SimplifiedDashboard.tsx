"use client";

import { useState, useEffect } from "react";
import { useAccount, useWalletClient } from "wagmi";
import { ConnectKitButton } from "connectkit";
import { apiService } from "@/lib/api";
import {
  publishTermRoot,
  waitForTransactionConfirmation,
} from "@/lib/blockchain";

export function SimplifiedDashboard() {
  const { isConnected } = useAccount();
  const [activeTab, setActiveTab] = useState<"publish" | "verify-course" | "verify-journey">("publish");

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                IU-MiCert Dashboard
              </h1>
              <p className="text-sm text-gray-500">
                Verkle-based academic credential management system
              </p>
            </div>
            <ConnectKitButton.Custom>
              {({ isConnected, show, truncatedAddress }) => (
                <button
                  onClick={show}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  {isConnected ? truncatedAddress : "Connect Wallet"}
                </button>
              )}
            </ConnectKitButton.Custom>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Navigation Tabs */}
        <div className="border-b border-gray-200 mb-8">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: "publish", name: "Publish Terms", icon: "⛓️" },
              { id: "verify-course", name: "Verify Course", icon: "🎓" },
              { id: "verify-journey", name: "Verify Journey", icon: "🔐" },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`${
                  activeTab === tab.id
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
              >
                <span>{tab.icon}</span>
                {tab.name}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        {isConnected ? (
          <>
            {activeTab === "publish" && <PublishTermTab />}
            {activeTab === "verify-course" && <VerifyCourseTab />}
            {activeTab === "verify-journey" && <VerifyJourneyTab />}
          </>
        ) : (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">🔒</div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              Wallet Connection Required
            </h2>
            <p className="text-gray-600 mb-8">
              Please connect your wallet to access the dashboard.
            </p>
            <ConnectKitButton.Custom>
              {({ show }) => (
                <button
                  onClick={show}
                  className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Connect Wallet
                </button>
              )}
            </ConnectKitButton.Custom>
          </div>
        )}
      </main>
    </div>
  );
}

// Tab 1: Publish Terms to Blockchain
function PublishTermTab() {
  const { address } = useAccount();
  const { data: walletClient } = useWalletClient();
  const [terms, setTerms] = useState<any[]>([]);
  const [selectedTerm, setSelectedTerm] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [publishStatus, setPublishStatus] = useState<{
    status: string;
    txHash?: string;
    error?: string;
  } | null>(null);

  useEffect(() => {
    loadTerms();
  }, []);

  const loadTerms = async () => {
    try {
      const data = await apiService.getTerms();
      setTerms(data);
    } catch (error) {
      console.error("Failed to load terms:", error);
    }
  };

  const handlePublish = async () => {
    if (!selectedTerm) return;

    setIsLoading(true);
    setPublishStatus({ status: "publishing" });

    try {
      // Get term root
      const rootData = await apiService.getTermRoot(selectedTerm);

      // Publish to blockchain with wagmi wallet client
      const result = await publishTermRoot({
        term_id: selectedTerm,
        verkle_root: rootData.verkle_root,
        total_students: rootData.total_students,
      }, address, walletClient);

      // Wait for confirmation
      await waitForTransactionConfirmation(result.transactionHash);

      setPublishStatus({
        status: "success",
        txHash: result.transactionHash,
      });
    } catch (error: any) {
      setPublishStatus({
        status: "error",
        error: error.message,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-2xl font-bold mb-6">Publish Term Roots to Blockchain</h2>

      <div className="space-y-6">
        {/* Term Selection */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Select Term
          </label>
          <select
            value={selectedTerm}
            onChange={(e) => setSelectedTerm(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">-- Select a term --</option>
            {terms.map((term) => (
              <option key={term.id} value={term.id}>
                {term.name || term.id}
              </option>
            ))}
          </select>
        </div>

        {/* Publish Button */}
        <button
          onClick={handlePublish}
          disabled={!selectedTerm || isLoading}
          className="w-full bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? "Publishing..." : "Publish to Blockchain"}
        </button>

        {/* Status Display */}
        {publishStatus && (
          <div
            className={`p-4 rounded-lg ${
              publishStatus.status === "success"
                ? "bg-green-50 border border-green-200"
                : publishStatus.status === "error"
                ? "bg-red-50 border border-red-200"
                : "bg-blue-50 border border-blue-200"
            }`}
          >
            {publishStatus.status === "success" && (
              <>
                <p className="text-green-800 font-medium mb-2">✅ Successfully published!</p>
                <p className="text-sm text-green-700">
                  Transaction: {publishStatus.txHash}
                </p>
                <a
                  href={`https://sepolia.etherscan.io/tx/${publishStatus.txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-blue-600 hover:underline"
                >
                  View on Etherscan →
                </a>
              </>
            )}
            {publishStatus.status === "error" && (
              <p className="text-red-800">❌ Error: {publishStatus.error}</p>
            )}
            {publishStatus.status === "publishing" && (
              <p className="text-blue-800">⏳ Publishing to blockchain...</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

// Tab 2: Verify Individual Course from Receipt
function VerifyCourseTab() {
  const [receiptData, setReceiptData] = useState("");
  const [courseId, setCourseId] = useState("");
  const [termId, setTermId] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [verificationResult, setVerificationResult] = useState<any>(null);

  const handleVerify = async () => {
    setIsLoading(true);
    setVerificationResult(null);

    try {
      const receipt = JSON.parse(receiptData);
      const result = await apiService.verifyCourse({
        receipt,
        course_id: courseId,
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

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-2xl font-bold mb-6">Verify Individual Course</h2>

      <div className="space-y-6">
        {/* Receipt Input */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Student Receipt JSON
          </label>
          <textarea
            value={receiptData}
            onChange={(e) => setReceiptData(e.target.value)}
            placeholder='Paste receipt JSON here...'
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
            rows={8}
          />
        </div>

        {/* Course ID */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Course ID
          </label>
          <input
            type="text"
            value={courseId}
            onChange={(e) => setCourseId(e.target.value)}
            placeholder="e.g., IT153IU"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Term ID */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Term ID
          </label>
          <input
            type="text"
            value={termId}
            onChange={(e) => setTermId(e.target.value)}
            placeholder="e.g., Semester_1_2023"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Verify Button */}
        <button
          onClick={handleVerify}
          disabled={!receiptData || !courseId || !termId || isLoading}
          className="w-full bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? "Verifying..." : "Verify Course"}
        </button>

        {/* Verification Result */}
        {verificationResult && (
          <div
            className={`p-4 rounded-lg ${
              verificationResult.verified
                ? "bg-green-50 border border-green-200"
                : "bg-red-50 border border-red-200"
            }`}
          >
            {verificationResult.verified ? (
              <>
                <p className="text-green-800 font-medium mb-4">✅ Course Verified Successfully!</p>

                {verificationResult.course && (
                  <div className="space-y-2 text-sm">
                    <p><strong>Course:</strong> {verificationResult.course.course_name} ({verificationResult.course.course_id})</p>
                    <p><strong>Grade:</strong> {verificationResult.course.grade}</p>
                    <p><strong>Credits:</strong> {verificationResult.course.credits}</p>
                  </div>
                )}

                {verificationResult.verification_details && (
                  <div className="mt-4 pt-4 border-t border-green-200">
                    <p className="font-medium mb-2">Verification Details:</p>
                    <ul className="space-y-1 text-sm">
                      <li>
                        {verificationResult.verification_details.ipa_verified ? "✅" : "❌"} IPA Cryptographic Proof
                      </li>
                      <li>
                        {verificationResult.verification_details.state_diff_verified ? "✅" : "❌"} State Diff Verification
                      </li>
                      <li>
                        {verificationResult.verification_details.blockchain_anchored ? "✅" : "❌"} Blockchain Anchored
                      </li>
                    </ul>
                  </div>
                )}
              </>
            ) : (
              <p className="text-red-800">
                ❌ Verification Failed: {verificationResult.verification_error || "Unknown error"}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

// Tab 3: Verify Full Learning Journey (Accumulated Receipt)
function VerifyJourneyTab() {
  const [studentId, setStudentId] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [verificationResult, setVerificationResult] = useState<any>(null);
  const [accumulatedReceipt, setAccumulatedReceipt] = useState<any>(null);

  const handleLoadAndVerify = async () => {
    setIsLoading(true);
    setVerificationResult(null);
    setAccumulatedReceipt(null);

    try {
      // First, get the accumulated receipt
      const receiptData = await apiService.getStudentAccumulatedReceipt(studentId);
      setAccumulatedReceipt(receiptData.receipt);

      // Then verify it with full IPA
      const result = await apiService.verifyReceiptIPA(receiptData.receipt);
      setVerificationResult(result);
    } catch (error: any) {
      setVerificationResult({
        status: "error",
        error: error.message,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-2xl font-bold mb-6">Verify Full Learning Journey</h2>

      <div className="space-y-6">
        {/* Student ID Input */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Student ID
          </label>
          <input
            type="text"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
            placeholder="e.g., ITITIU00001"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Verify Button */}
        <button
          onClick={handleLoadAndVerify}
          disabled={!studentId || isLoading}
          className="w-full bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? "Verifying Full Journey..." : "Load & Verify Journey"}
        </button>

        {/* Accumulated Receipt Summary */}
        {accumulatedReceipt && (
          <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="font-medium mb-2">📊 Academic Summary</p>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-600">Completed Terms</p>
                <p className="font-bold">{accumulatedReceipt.CompletedTerms}</p>
              </div>
              <div>
                <p className="text-gray-600">Total Courses</p>
                <p className="font-bold">{accumulatedReceipt.TotalCourses}</p>
              </div>
              <div>
                <p className="text-gray-600">Total Credits</p>
                <p className="font-bold">{accumulatedReceipt.TotalCredits}</p>
              </div>
              <div>
                <p className="text-gray-600">GPA</p>
                <p className="font-bold">{accumulatedReceipt.GPA.toFixed(2)}</p>
              </div>
            </div>
          </div>
        )}

        {/* Verification Result */}
        {verificationResult && (
          <div
            className={`p-4 rounded-lg ${
              verificationResult.status === "success"
                ? "bg-green-50 border border-green-200"
                : "bg-red-50 border border-red-200"
            }`}
          >
            {verificationResult.status === "success" ? (
              <>
                <p className="text-green-800 font-medium mb-4">
                  ✅ Full Journey Verified with IPA Cryptographic Proofs!
                </p>

                <div className="space-y-2 text-sm">
                  <p><strong>Total Courses:</strong> {verificationResult.total_courses}</p>
                  <p><strong>Verified Courses:</strong> {verificationResult.verified_courses}</p>
                  <p><strong>Failed Courses:</strong> {verificationResult.failed_courses}</p>
                </div>

                {verificationResult.failed_list && verificationResult.failed_list.length > 0 && (
                  <div className="mt-4 pt-4 border-t border-green-200">
                    <p className="font-medium text-red-700 mb-2">Failed Courses:</p>
                    <ul className="text-sm text-red-600 list-disc list-inside">
                      {verificationResult.failed_list.map((course: string, idx: number) => (
                        <li key={idx}>{course}</li>
                      ))}
                    </ul>
                  </div>
                )}

                <div className="mt-4 pt-4 border-t border-green-200">
                  <p className="text-xs text-gray-600">{verificationResult.computation_note}</p>
                </div>
              </>
            ) : (
              <p className="text-red-800">
                ❌ Verification Failed: {verificationResult.error || "Unknown error"}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
