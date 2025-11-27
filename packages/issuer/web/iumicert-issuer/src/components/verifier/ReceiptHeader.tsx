import { JourneyReceipt, VerificationResult } from "./types";

interface ReceiptHeaderProps {
  receipt: JourneyReceipt;
  overallResult?: VerificationResult;
  isVerifying: boolean;
  isSelectiveMode: boolean;
  selectedCourseCount: number;
  error: string;
  onUploadNew: () => void;
  onVerify: () => void;
  onInitializeSelective: () => void;
  onDownloadSelective: () => void;
  onCancelSelective: () => void;
}

export function ReceiptHeader({
  receipt,
  overallResult,
  isVerifying,
  isSelectiveMode,
  selectedCourseCount,
  error,
  onUploadNew,
  onVerify,
  onInitializeSelective,
  onDownloadSelective,
  onCancelSelective,
}: ReceiptHeaderProps) {
  return (
    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">
            Academic Journey
          </h1>
          <p className="text-gray-600 mt-1">
            Student ID:{" "}
            <span className="font-mono font-semibold text-slate-900">
              {receipt.student_id}
            </span>
          </p>
        </div>
        <button
          type="button"
          onClick={onUploadNew}
          className="px-5 py-2.5 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors font-medium"
        >
          Upload New Receipt
        </button>
      </div>

      {/* Action Buttons */}
      <div className="flex flex-wrap items-center gap-4">
        <button
          type="button"
          onClick={onVerify}
          disabled={isVerifying}
          className={`px-6 py-3 rounded-xl font-semibold transition-all duration-200 ${
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

        {!isSelectiveMode ? (
          <button
            type="button"
            onClick={onInitializeSelective}
            className="px-6 py-3 rounded-xl font-semibold transition-all duration-200 bg-gradient-to-r from-purple-500 to-indigo-600 hover:shadow-lg hover:shadow-purple-500/30 text-white"
          >
            Create Selective Receipt
          </button>
        ) : (
          <>
            <button
              type="button"
              onClick={onDownloadSelective}
              disabled={selectedCourseCount === 0}
              className={`px-6 py-3 rounded-xl font-semibold transition-all duration-200 ${
                selectedCourseCount === 0
                  ? "bg-gray-400 cursor-not-allowed"
                  : "bg-gradient-to-r from-green-500 to-emerald-600 hover:shadow-lg hover:shadow-green-500/30"
              } text-white`}
            >
              <span className="flex items-center gap-2">
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                  />
                </svg>
                Download Selective Receipt ({selectedCourseCount} courses)
              </span>
            </button>
            <button
              type="button"
              onClick={onCancelSelective}
              className="px-6 py-3 rounded-xl font-semibold transition-all duration-200 bg-gray-200 hover:bg-gray-300 text-gray-700"
            >
              Cancel
            </button>
          </>
        )}

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

      {/* Selective Disclosure Info */}
      {isSelectiveMode && (
        <div className="mt-4 p-4 bg-purple-50 border border-purple-200 rounded-xl">
          <div className="flex items-start gap-3">
            <svg
              className="w-5 h-5 text-purple-600 mt-0.5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <div className="flex-1">
              <p className="text-sm font-semibold text-purple-800 mb-1">
                Selective Disclosure Mode
              </p>
              <p className="text-sm text-purple-700">
                Click on terms and courses below to include/exclude them from
                your selective receipt. Unchecked items will be removed, but all
                remaining courses will still verify cryptographically.
              </p>
            </div>
          </div>
        </div>
      )}

      {overallResult?.details && (
        <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-xl">
          <p className="text-sm text-blue-800">{overallResult.details}</p>
        </div>
      )}

      {error && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-3">
          <svg
            className="w-5 h-5 mt-0.5 flex-shrink-0 text-red-600"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clipRule="evenodd"
            />
          </svg>
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}
    </div>
  );
}
