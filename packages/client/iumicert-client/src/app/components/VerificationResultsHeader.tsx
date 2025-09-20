import { CheckCircle, AlertTriangle, Download } from "lucide-react";
import { VerificationResult } from "../types/proofs";

interface VerificationResultsHeaderProps {
  result: VerificationResult;
  onReset: () => void;
  onCopyData: () => void;
}

export default function VerificationResultsHeader({
  result,
  onReset,
  onCopyData
}: VerificationResultsHeaderProps) {
  return (
    <div className={`p-4 rounded-xl mb-6 border ${
      result.isValid
        ? "bg-green-500/20 border-green-400/30"
        : "bg-red-500/20 border-red-400/30"
    }`}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${
            result.isValid 
              ? "bg-green-500/20" 
              : "bg-red-500/20"
          }`}>
            {result.isValid ? (
              <CheckCircle className="w-5 h-5 text-green-400" />
            ) : (
              <AlertTriangle className="w-5 h-5 text-red-400" />
            )}
          </div>
          <div>
            <div className={`text-lg font-bold font-space-grotesk ${
              result.isValid ? "text-green-100" : "text-red-100"
            }`}>
              {result.message}
            </div>
            {result.isValid && result.proofData && (
              <div className="text-sm text-purple-200 mt-1 font-inter">
                Type: {result.proofData.type === "individual_term" 
                  ? "Individual Term Proof" 
                  : "Aggregated Academic Journey"}
              </div>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          {result.isValid && (
            <button
              onClick={onCopyData}
              className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white border border-white/20 rounded-lg transition duration-300 flex items-center gap-2 font-inter"
            >
              <Download className="w-4 h-4" />
              Copy Data
            </button>
          )}
          <button
            onClick={onReset}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition duration-300 font-inter"
          >
            New Verification
          </button>
        </div>
      </div>
    </div>
  );
}