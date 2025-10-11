"use client";

import { useState, useEffect } from "react";
import { useAccount } from "wagmi";
import { ConnectKitButton } from "connectkit";
import { PublishTermTab } from "./issuer/PublishTermTab";

export function IssuerDashboard() {
  const { isConnected } = useAccount();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

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
            <div className="flex items-center gap-4">
              <a
                href="/demo-data"
                className="text-gray-600 hover:text-gray-900 px-4 py-2 text-sm font-medium transition-colors"
              >
                🎲 Demo Data
              </a>
              <a
                href="/verifier"
                className="text-gray-600 hover:text-gray-900 px-4 py-2 text-sm font-medium transition-colors"
              >
                🔍 Verifier
              </a>
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
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Tab Content */}
        {isConnected ? (
          <PublishTermTab />
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
