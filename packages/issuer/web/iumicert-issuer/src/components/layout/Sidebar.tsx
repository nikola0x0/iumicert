"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { FileText, ShieldCheck, Database } from "lucide-react";
import { cn } from "@/lib/utils";

interface NavItem {
  label: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
}

const navItems: NavItem[] = [
  {
    label: "Publish Terms",
    href: "/",
    icon: FileText,
  },
  {
    label: "Verifier",
    href: "/verifier",
    icon: ShieldCheck,
  },
  {
    label: "Demo Data",
    href: "/demo-data",
    icon: Database,
  },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-screen w-[260px] bg-gradient-to-b from-white to-gray-50 border-r border-gray-200 flex flex-col shadow-lg">
      {/* Logo Area with gradient */}
      <div className="py-6 px-6 border-b border-gray-200 bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
        <div className="flex items-center gap-2 mb-1">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-600 to-indigo-600 flex items-center justify-center shadow-md">
            <span className="text-white font-bold text-sm">IU</span>
          </div>
          <h1 className="text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">MiCert</h1>
        </div>
        <p className="text-xs text-gray-600 font-medium ml-10">Issuer Dashboard</p>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 py-6">
        <ul className="space-y-2">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            const Icon = item.icon;

            return (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group",
                    isActive
                      ? "bg-gradient-to-r from-blue-500 via-indigo-600 to-blue-500 text-white shadow-lg shadow-blue-500/30"
                      : "text-gray-600 hover:bg-gradient-to-r hover:from-blue-50 hover:to-indigo-50 hover:shadow-md"
                  )}
                >
                  <div className={cn(
                    "w-8 h-8 rounded-lg flex items-center justify-center transition-all",
                    isActive ? "bg-white/20" : "bg-gray-100 group-hover:bg-blue-100"
                  )}>
                    <Icon
                      className={cn(
                        "w-4 h-4 transition-colors",
                        isActive ? "text-white" : "text-gray-600 group-hover:text-blue-600"
                      )}
                    />
                  </div>
                  <span className={cn(
                    "font-medium text-sm transition-colors",
                    isActive ? "text-white" : "group-hover:text-blue-700"
                  )}>{item.label}</span>
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Footer */}
      <div className="px-4 py-4 border-t border-gray-200">
        <p className="text-[11px] text-gray-400">Version 1.0.0</p>
      </div>
    </aside>
  );
}
