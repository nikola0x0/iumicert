"use client";

import { Monitor } from "lucide-react";

export function MobileWarning() {
  return (
    <div className="fixed inset-0 bg-white z-50 flex items-center justify-center p-6 lg:hidden">
      <div className="max-w-md text-center">
        <div className="w-20 h-20 bg-blue-100 rounded-2xl flex items-center justify-center mx-auto mb-6">
          <Monitor className="w-10 h-10 text-blue-600" />
        </div>
        <h1 className="text-2xl font-bold text-slate-900 mb-4">
          Desktop Access Required
        </h1>
        <p className="text-gray-600 leading-relaxed">
          Please access this dashboard from a desktop computer for the best
          experience. The issuer dashboard is optimized for larger screens to
          provide you with better functionality and usability.
        </p>
      </div>
    </div>
  );
}
