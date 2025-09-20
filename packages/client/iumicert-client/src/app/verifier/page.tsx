"use client";

import { useState } from "react";
import DashboardLayout from "../components/DashboardLayout";
import VerificationUpload from "../components/VerificationUpload";
import VerificationResultsHeader from "../components/VerificationResultsHeader";
import SingleTermView from "../components/SingleTermView";
import AggregatedJourneyView from "../components/AggregatedJourneyView";
import { detectProofType } from "../utils/proofDetection";
import { VerificationResult } from "../types/proofs";

export default function VerifierDashboard() {
  const [credentialData, setCredentialData] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [verificationResult, setVerificationResult] =
    useState<VerificationResult | null>(null);
  const [selectedTerm, setSelectedTerm] = useState<string | null>(null);
  const [viewMode, setViewMode] = useState<"overview" | "single_term">(
    "overview"
  );

  const handleVerify = async () => {
    if (!credentialData.trim()) return;

    setIsLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 2000));

      const parsedData = JSON.parse(credentialData);
      const proofData = detectProofType(parsedData);

      if (!proofData) {
        throw new Error("Invalid proof format");
      }

      setVerificationResult({
        isValid: true,
        message: `${
          proofData.type === "individual_term" ||
          proofData.type === "single_term"
            ? "Individual Term"
            : "Academic Journey"
        } verified successfully on blockchain`,
        proofData,
      });

      if (proofData.type === "aggregated_journey") {
        setViewMode("overview");
        setSelectedTerm(null);
      } else {
        setViewMode("single_term");
      }
    } catch {
      setVerificationResult({
        isValid: false,
        message:
          "Credential verification failed - Invalid JSON format or proof structure",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleTermSelect = (term: string) => {
    setSelectedTerm(term);
    setViewMode("single_term");
  };

  const handleBackToOverview = () => {
    setViewMode("overview");
    setSelectedTerm(null);
  };

  const handleReset = () => {
    setCredentialData("");
    setVerificationResult(null);
    setSelectedTerm(null);
    setViewMode("overview");
  };

  const handleFileChange = async (file: File) => {
    try {
      const text = await file.text();
      setCredentialData(text);
    } catch (error) {
      console.error("Error reading file:", error);
      setVerificationResult({
        isValid: false,
        message: "Failed to read file. Please ensure it's a valid JSON file.",
      });
    }
  };

  const handleTypeError = () => {
    setVerificationResult({
      isValid: false,
      message: "Please upload a JSON file containing the proof package.",
    });
  };

  const handleCopyData = () => {
    if (verificationResult?.proofData) {
      navigator.clipboard.writeText(
        JSON.stringify(verificationResult.proofData, null, 2)
      );
    }
  };

  const renderVerifyTab = () => {
    return (
      <div className="h-full flex flex-col">
        {!verificationResult ? (
          <div className="h-full flex items-center justify-center p-4">
            <VerificationUpload
              credentialData={credentialData}
              isLoading={isLoading}
              onCredentialChange={setCredentialData}
              onVerify={handleVerify}
              onFileChange={handleFileChange}
              onTypeError={handleTypeError}
            />
          </div>
        ) : (
          <div className="h-full flex flex-col">
            <VerificationResultsHeader
              result={verificationResult}
              onReset={handleReset}
              onCopyData={handleCopyData}
            />

            {/* Verification Results */}
            {verificationResult.isValid && verificationResult.proofData && (
              <div className="flex-1 min-h-0">
                {verificationResult.proofData.type === "aggregated_journey" &&
                viewMode === "overview" ? (
                  <AggregatedJourneyView
                    proofData={verificationResult.proofData}
                    onTermSelect={handleTermSelect}
                  />
                ) : (
                  <SingleTermView
                    proofData={verificationResult.proofData}
                    selectedTerm={selectedTerm}
                    onBackToOverview={handleBackToOverview}
                  />
                )}
              </div>
            )}
          </div>
        )}
      </div>
    );
  };

  return (
    <DashboardLayout activeSection="verify">
      <div className="h-full flex flex-col">{renderVerifyTab()}</div>
    </DashboardLayout>
  );
}
