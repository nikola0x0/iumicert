"use client";

import { useState } from "react";
import DashboardLayout from "../components/DashboardLayout";
import { X, AlertTriangle, Shield, Lock, FileText, Settings } from "lucide-react";

export default function RevokePage() {
  const [isLoading, setIsLoading] = useState(false);
  const [certificateId, setCertificateId] = useState("");
  const [reason, setReason] = useState("");
  const [revocationResult, setRevocationResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);

  const handleRevoke = async () => {
    if (!certificateId.trim() || !reason.trim()) return;
    
    setIsLoading(true);
    
    // Simulate API call
    setTimeout(() => {
      setRevocationResult({
        success: false,
        message: "Revocation functionality not yet implemented. This feature requires institutional administrator privileges."
      });
      setIsLoading(false);
    }, 1500);
  };

  const handleReset = () => {
    setCertificateId("");
    setReason("");
    setRevocationResult(null);
  };

  return (
    <DashboardLayout activeSection="revoke">
      <div className="h-full flex items-center justify-center p-4">
        {!revocationResult ? (
          <div className="glass-effect rounded-xl p-8 max-w-md w-full">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 bg-gradient-to-br from-red-500 to-orange-600 rounded-full flex items-center justify-center mx-auto mb-4">
                <X className="w-8 h-8 text-white" />
              </div>
              <h1 className="text-2xl font-bold text-white font-space-grotesk mb-2">
                Certificate Revocation
              </h1>
              <p className="text-purple-200 font-inter text-sm">
                Revoke issued certificates and manage certificate status
              </p>
            </div>

            {/* Warning Alert */}
            <div className="bg-red-500/20 border border-red-500/30 rounded-lg p-4 mb-6">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="w-4 h-4 text-red-300 flex-shrink-0" />
                <span className="text-red-300 font-semibold text-sm font-space-grotesk">
                  Administrative Access Required
                </span>
              </div>
              <p className="text-red-200 text-xs font-inter">
                This feature requires institutional administrator privileges and proper authentication.
              </p>
            </div>

            {/* Form */}
            <div className="space-y-4">
              <div>
                <label className="block text-white font-medium text-sm font-space-grotesk mb-2">
                  Certificate ID
                </label>
                <div className="relative">
                  <FileText className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-white/50" />
                  <input
                    type="text"
                    placeholder="Enter certificate ID or transaction hash"
                    className="w-full pl-10 pr-4 py-3 bg-white/10 border border-white/20 rounded-xl 
                             text-white placeholder-white/50 backdrop-blur-sm
                             focus:outline-none focus:ring-2 focus:ring-red-400/50 focus:border-transparent 
                             transition duration-300"
                    value={certificateId}
                    onChange={(e) => setCertificateId(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-white font-medium text-sm font-space-grotesk mb-2">
                  Revocation Reason
                </label>
                <textarea
                  rows={3}
                  placeholder="Provide reason for revocation..."
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl 
                           text-white placeholder-white/50 backdrop-blur-sm
                           focus:outline-none focus:ring-2 focus:ring-red-400/50 focus:border-transparent 
                           transition duration-300 resize-none"
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                />
              </div>

              <button
                onClick={handleRevoke}
                disabled={!certificateId.trim() || !reason.trim() || isLoading}
                className="w-full bg-gradient-to-r from-red-500 to-orange-600 hover:from-red-600 hover:to-orange-700
                         disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed
                         text-white font-bold py-3 px-6 rounded-xl transition-all duration-300
                         hover:scale-105 hover:shadow-xl flex items-center justify-center gap-2 font-space-grotesk"
              >
                {isLoading ? (
                  <>
                    <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Processing...
                  </>
                ) : (
                  <>
                    <Lock className="w-5 h-5" />
                    Revoke Certificate
                  </>
                )}
              </button>
            </div>

            {/* Info */}
            <div className="mt-6 pt-4 border-t border-white/20">
              <div className="flex items-center gap-2 text-white/60 text-xs font-inter">
                <Settings className="w-3 h-3" />
                <span>Revocation requires administrator authentication</span>
              </div>
            </div>
          </div>
        ) : (
          /* Results */
          <div className="glass-effect rounded-xl p-8 max-w-md w-full text-center">
            <div className={`w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 ${
              revocationResult.success 
                ? 'bg-gradient-to-br from-green-500 to-emerald-600'
                : 'bg-gradient-to-br from-red-500 to-orange-600'
            }`}>
              {revocationResult.success ? (
                <Shield className="w-8 h-8 text-white" />
              ) : (
                <AlertTriangle className="w-8 h-8 text-white" />
              )}
            </div>
            
            <h2 className="text-2xl font-bold text-white font-space-grotesk mb-4">
              {revocationResult.success ? 'Certificate Revoked' : 'Revocation Failed'}
            </h2>
            
            <p className="text-purple-200 font-inter mb-6 text-sm leading-relaxed">
              {revocationResult.message}
            </p>

            <button
              onClick={handleReset}
              className="bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700
                       text-white font-bold py-3 px-6 rounded-xl transition-all duration-300
                       hover:scale-105 hover:shadow-xl font-space-grotesk"
            >
              Try Another Certificate
            </button>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}