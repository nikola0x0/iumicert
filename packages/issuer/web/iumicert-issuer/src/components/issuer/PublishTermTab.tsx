"use client";

import { useState, useEffect } from "react";
import { useAccount, useWalletClient } from "wagmi";
import { Upload, Loader2, CheckCircle2, XCircle, ExternalLink, AlertTriangle, Info } from "lucide-react";
import { apiService } from "@/lib/api";
import {
  publishTermRoot,
  waitForTransactionConfirmation,
} from "@/lib/blockchain";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";

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
      console.error("Publish error:", error);

      // Parse the error message to provide better feedback
      let errorMessage = error.message || "An unknown error occurred";
      let errorType = "error";

      // Check for specific error patterns
      if (errorMessage.includes("Root already published")) {
        errorType = "already_published";
        errorMessage = `This term (${selectedTerm}) has already been published to the blockchain. Each term can only be published once to ensure immutability and prevent tampering.`;
      } else if (errorMessage.includes("User rejected") || errorMessage.includes("User denied")) {
        errorType = "rejected";
        errorMessage = "Transaction was rejected. Please try again and approve the transaction in your wallet.";
      } else if (errorMessage.includes("insufficient funds")) {
        errorType = "insufficient_funds";
        errorMessage = "Insufficient funds to complete the transaction. Please ensure you have enough ETH for gas fees.";
      } else if (errorMessage.includes("network") || errorMessage.includes("connection")) {
        errorType = "network";
        errorMessage = "Network error occurred. Please check your connection and try again.";
      }

      setPublishStatus({
        status: errorType,
        error: errorMessage,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="relative overflow-hidden rounded-3xl shadow-2xl border-2 border-blue-200 hover:border-blue-300 transition-all duration-300 bg-white">
      <CardHeader className="relative">
        <div className="flex items-start gap-4">
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 via-indigo-600 to-purple-600 flex items-center justify-center shadow-xl shadow-blue-500/40">
              <Upload className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <CardTitle className="text-3xl bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              Publish Term Root
            </CardTitle>
            <CardDescription className="text-base text-gray-600 leading-relaxed mt-2">
              Publish academic term roots to the Ethereum Sepolia testnet for permanent verification
            </CardDescription>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Term Selection */}
        <div className="relative">
          <label className="flex items-center gap-2 text-sm font-semibold text-slate-800 mb-3">
            <span className="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-xs">1</span>
            Select Academic Term
          </label>
          <select
            value={selectedTerm}
            onChange={(e) => setSelectedTerm(e.target.value)}
            className="w-full px-5 py-4 border-2 border-gray-200 rounded-2xl focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 hover:border-blue-300 transition-all bg-white shadow-sm text-lg font-semibold text-gray-900"
          >
            <option value="">-- Select a term to publish --</option>
            {terms.map((term) => (
              <option key={term.id} value={term.id}>
                {term.name || term.id.replace(/_/g, ' ')}
              </option>
            ))}
          </select>
          <p className="text-sm text-gray-500 mt-3 ml-8">
            Choose the academic term whose Verkle root you want to publish to the blockchain
          </p>
        </div>

        {/* Publish Button */}
        <div className="pt-4">
          <Button
            onClick={handlePublish}
            disabled={!selectedTerm || isLoading}
            className="w-full text-white rounded-2xl bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 hover:from-blue-700 hover:via-indigo-700 hover:to-purple-700 shadow-2xl shadow-blue-500/50 hover:shadow-blue-600/60 transition-all transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none group relative overflow-hidden"
            size="lg"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            <div className="relative flex items-center justify-center gap-3 py-1">
              {isLoading ? (
                <>
                  <Loader2 className="w-6 h-6 animate-spin" />
                  <span className="text-lg font-bold">Publishing to Blockchain...</span>
                </>
              ) : (
                <>
                  <Upload className="w-6 h-6" />
                  <span className="text-lg font-bold">Publish to Blockchain</span>
                </>
              )}
            </div>
          </Button>
        </div>

        {/* Status Messages */}
        {publishStatus && (
          <Alert className={`rounded-xl flex items-start gap-3 ${
            publishStatus.status === "success"
              ? "bg-gradient-to-r from-green-50 to-emerald-50 border-green-200 text-green-900"
              : publishStatus.status === "already_published"
              ? "bg-gradient-to-r from-yellow-50 to-amber-50 border-yellow-200 text-yellow-900"
              : publishStatus.status === "rejected"
              ? "bg-gradient-to-r from-orange-50 to-orange-50 border-orange-200 text-orange-900"
              : publishStatus.status === "insufficient_funds"
              ? "bg-gradient-to-r from-red-50 to-rose-50 border-red-200 text-red-900"
              : publishStatus.status === "error"
              ? "bg-gradient-to-r from-red-50 to-rose-50 border-red-200 text-red-900"
              : "bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200 text-blue-900"
          }`}>
            <div className={`w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 ${
              publishStatus.status === "success"
                ? "bg-green-100"
                : publishStatus.status === "already_published"
                ? "bg-yellow-100"
                : publishStatus.status === "rejected"
                ? "bg-orange-100"
                : publishStatus.status === "insufficient_funds"
                ? "bg-red-100"
                : publishStatus.status === "error"
                ? "bg-red-100"
                : "bg-blue-100"
            }`}>
              {publishStatus.status === "success" && <CheckCircle2 className="h-5 w-5 text-green-600" />}
              {publishStatus.status === "already_published" && <Info className="h-5 w-5 text-yellow-600" />}
              {publishStatus.status === "rejected" && <AlertTriangle className="h-5 w-5 text-orange-600" />}
              {publishStatus.status === "insufficient_funds" && <XCircle className="h-5 w-5 text-red-600" />}
              {publishStatus.status === "error" && <XCircle className="h-5 w-5 text-red-600" />}
              {publishStatus.status === "publishing" && <Loader2 className="h-5 w-5 text-blue-600 animate-spin" />}
            </div>

            <AlertDescription className="flex-1">
              {publishStatus.status === "success" && (
                <div>
                  <p className="font-semibold mb-2">Successfully Published!</p>
                  <p className="text-sm font-mono bg-white/50 p-2 rounded-lg break-all mb-3">
                    {publishStatus.txHash}
                  </p>
                  <a
                    href={`https://sepolia.etherscan.io/tx/${publishStatus.txHash}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-700 hover:underline"
                  >
                    View on Etherscan
                    <ExternalLink className="w-4 h-4" />
                  </a>
                </div>
              )}
              {publishStatus.status === "already_published" && (
                <div>
                  <p className="font-semibold mb-2">Term Already Published</p>
                  <p className="text-sm">{publishStatus.error}</p>
                </div>
              )}
              {publishStatus.status === "rejected" && (
                <div>
                  <p className="font-semibold mb-1">Transaction Rejected</p>
                  <p className="text-sm">{publishStatus.error}</p>
                </div>
              )}
              {publishStatus.status === "insufficient_funds" && (
                <div>
                  <p className="font-semibold mb-1">Insufficient Funds</p>
                  <p className="text-sm">{publishStatus.error}</p>
                </div>
              )}
              {publishStatus.status === "error" && (
                <div>
                  <p className="font-semibold mb-1">Publication Failed</p>
                  <p className="text-sm">{publishStatus.error}</p>
                </div>
              )}
              {publishStatus.status === "publishing" && (
                <p className="font-semibold">Publishing to blockchain...</p>
              )}
            </AlertDescription>
          </Alert>
        )}
      </CardContent>
    </Card>
  );
}
