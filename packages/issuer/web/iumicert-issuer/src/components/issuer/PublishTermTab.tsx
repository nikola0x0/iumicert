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
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-2xl font-bold mb-6">
        Publish Term Roots to Blockchain
      </h2>

      <div className="space-y-6">
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

        <button
          onClick={handlePublish}
          disabled={!selectedTerm || isLoading}
          className="w-full bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? "Publishing..." : "Publish to Blockchain"}
        </button>

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
                <p className="text-green-800 font-medium mb-2">
                  ✅ Successfully published!
                </p>
                <p className="text-sm text-green-700 break-all mb-2">
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
