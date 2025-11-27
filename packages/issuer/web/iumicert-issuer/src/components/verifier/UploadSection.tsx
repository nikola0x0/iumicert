import React from "react";

interface UploadSectionProps {
  onFileUpload: (event: React.ChangeEvent<HTMLInputElement>) => void;
  error: string;
}

export function UploadSection({ onFileUpload, error }: UploadSectionProps) {
  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="relative overflow-hidden bg-gradient-to-br from-blue-500 via-indigo-600 to-purple-600 rounded-3xl p-10 shadow-2xl shadow-blue-500/20">
        <div className="absolute top-0 right-0 w-72 h-72 bg-white/5 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-12 -left-12 w-64 h-64 bg-white/5 rounded-full blur-3xl"></div>
        <div className="absolute top-1/2 right-1/4 w-32 h-32 border-4 border-white/10 rounded-full"></div>
        <div className="absolute bottom-1/4 left-1/3 w-20 h-20 border-4 border-white/10 rounded-full"></div>

        <div className="relative z-10">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 mb-4">
            <svg
              className="w-4 h-4 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
              />
            </svg>
            <span className="text-xs font-semibold text-white/90">
              CREDENTIAL VERIFICATION
            </span>
          </div>
          <h1 className="text-5xl font-extrabold text-white mb-3 tracking-tight">
            Receipt Verifier
          </h1>
          <p className="text-blue-100 text-lg max-w-2xl leading-relaxed">
            Upload an academic journey receipt to verify credentials, view the
            academic history
          </p>
        </div>
      </div>

      <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-8">
        <div className="border-2 border-dashed border-gray-300 rounded-xl p-12 text-center hover:border-blue-400 hover:bg-blue-50/30 transition-all duration-200">
          <input
            type="file"
            accept=".json"
            onChange={onFileUpload}
            className="hidden"
            id="receipt-upload"
          />
          <label
            htmlFor="receipt-upload"
            className="cursor-pointer flex flex-col items-center"
          >
            <div className="w-16 h-16 bg-blue-100 rounded-2xl flex items-center justify-center mb-4">
              <svg
                className="w-8 h-8 text-blue-600"
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
            <span className="text-xl font-semibold text-slate-900 mb-2">
              Click to upload receipt
            </span>
            <span className="text-sm text-gray-500">
              or drag and drop JSON file here
            </span>
          </label>
        </div>

        {error && (
          <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 flex items-start gap-3">
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
