"use client";

import { useState, useEffect } from "react";
import { useAccount, useWalletClient } from "wagmi";
import { apiService } from "@/lib/api";
import {
  publishTermRoot,
  waitForTransactionConfirmation,
} from "@/lib/blockchain";

export function PublishTermTab() {
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
      const rootData = await apiService.getTermRoot(selectedTerm);

      const result = await publishTermRoot(
        {
          term_id: selectedTerm,
          verkle_root: rootData.verkle_root,
          total_students: rootData.total_students,
        },
        address,
        walletClient
      );

      await waitForTransactionConfirmation(result.transactionHash);

      // Update database with blockchain info
      try {
        await fetch(
          `http://localhost:8080/api/terms/${selectedTerm}/blockchain`,
          {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              tx_hash: result.transactionHash,
              block_number: result.blockNumber ? Number(result.blockNumber) : 0,
              publisher_address: address,
            }),
          }
        );
      } catch (error) {
        console.error("Failed to update blockchain status in database:", error);
      }

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
    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-8">
      <div className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Select Term
          </label>
          <select
            value={selectedTerm}
            onChange={(e) => setSelectedTerm(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          >
            <option value="">-- Select a term --</option>
            {terms.map((term) => (
              <option key={term.id} value={term.id}>
                {term.name || term.id}
              </option>
            ))}
          </select>
        </div>

        <button
          onClick={handlePublish}
          disabled={!selectedTerm || isLoading}
          className="w-full bg-gradient-to-r from-blue-500 to-blue-600 text-white px-6 py-3.5 rounded-xl hover:shadow-lg hover:shadow-blue-500/30 disabled:from-gray-300 disabled:to-gray-300 disabled:cursor-not-allowed disabled:shadow-none transition-all duration-200 font-medium"
        >
          {isLoading ? "Publishing..." : "Publish to Blockchain"}
        </button>

        {publishStatus && (
          <div
            className={`p-5 rounded-xl border ${
              publishStatus.status === "success"
                ? "bg-green-50 border-green-200"
                : publishStatus.status === "error"
                ? "bg-red-50 border-red-200"
                : "bg-blue-50 border-blue-200"
            }`}
          >
            {publishStatus.status === "success" && (
              <>
                <p className="text-green-800 font-semibold mb-2 flex items-center gap-2">
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  Successfully published!
                </p>
                <p className="text-sm text-green-700 break-all mb-3 font-mono bg-white/50 p-2 rounded-lg">
                  {publishStatus.txHash}
                </p>
                <a
                  href={`https://sepolia.etherscan.io/tx/${publishStatus.txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium inline-flex items-center gap-1"
                >
                  View on Etherscan
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </>
            )}
            {publishStatus.status === "error" && (
              <p className="text-red-800 flex items-center gap-2">
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
                Error: {publishStatus.error}
              </p>
            )}
            {publishStatus.status === "publishing" && (
              <p className="text-blue-800 flex items-center gap-2">
                <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Publishing to blockchain...
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
