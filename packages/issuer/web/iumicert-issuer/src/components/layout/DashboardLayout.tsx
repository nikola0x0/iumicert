"use client";

import { useState, useEffect } from "react";
import { ConnectKitButton } from "connectkit";
import { Sidebar } from "./Sidebar";
import { MobileWarning } from "./MobileWarning";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <>
      {/* Mobile Warning - shows on screens < 1024px */}
      <MobileWarning />

      {/* Main Dashboard - hidden on mobile */}
      <div className="hidden lg:flex min-h-screen bg-slate-50">
        {/* Sidebar */}
        <Sidebar />

        {/* Main Content Area */}
        <div className="flex-1 ml-[260px]">
          {/* Header with Wallet Connect */}
          <header className="sticky top-0 z-40 bg-white border-b border-gray-200 shadow-sm">
            <div className="px-8 py-4 flex justify-end items-center">
              <ConnectKitButton.Custom>
                {({ isConnected, show, truncatedAddress }) => (
                  <button
                    onClick={show}
                    className="bg-gradient-to-r from-blue-500 to-blue-600 text-white px-6 py-2.5 rounded-xl hover:shadow-lg hover:shadow-blue-500/30 transition-all duration-200 font-medium text-sm"
                  >
                    {isConnected ? truncatedAddress : "Connect Wallet"}
                  </button>
                )}
              </ConnectKitButton.Custom>
            </div>
          </header>

          {/* Page Content */}
          <main className="p-8">{children}</main>
        </div>
      </div>
    </>
  );
}
