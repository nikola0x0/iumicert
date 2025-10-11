"use client";

import { useState, useEffect } from "react";
import Image from "next/image";
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
      <div className="hidden lg:flex min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/30">
        {/* Sidebar */}
        <Sidebar />

        {/* Main Content Area */}
        <div className="flex-1 ml-[260px]">
          {/* Header with Wallet Connect */}
          <header className="sticky top-0 z-40 bg-gradient-to-r from-slate-50 via-white to-blue-50/50 border-b border-blue-200/50 shadow-lg backdrop-blur-sm">
            <div className="px-8 py-5 flex justify-between items-center">
              <div className="flex items-center gap-5">
                <div>
                  <h2 className="text-xl font-extrabold bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 bg-clip-text text-transparent tracking-tight">
                    Issuer Dashboard
                  </h2>
                  <p className="text-xs text-slate-600 font-semibold tracking-wide">
                    ACADEMIC CREDENTIAL MANAGEMENT
                  </p>
                </div>
              </div>
              <ConnectKitButton.Custom>
                {({ isConnected, show, truncatedAddress }) => (
                  <button
                    onClick={show}
                    className="relative group bg-gradient-to-r from-blue-600 via-indigo-600 to-blue-600 text-white px-7 py-3.5 rounded-xl hover:shadow-2xl hover:shadow-blue-500/50 transition-all duration-300 font-bold text-sm shadow-xl shadow-blue-500/40 hover:scale-105"
                  >
                    <span className="relative z-10">
                      {isConnected ? truncatedAddress : "Connect Wallet"}
                    </span>
                    <div className="absolute inset-0 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                  </button>
                )}
              </ConnectKitButton.Custom>
            </div>
          </header>

          {/* Page Content */}
          <main className="p-8 min-h-screen">{children}</main>
        </div>
      </div>
    </>
  );
}
