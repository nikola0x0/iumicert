"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import Image from "next/image";
import { FileText, ShieldCheck, Database, AlertTriangle } from "lucide-react";
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
    label: "Revocations",
    href: "/revocations",
    icon: AlertTriangle,
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
    <aside className="fixed left-0 top-0 h-screen w-[260px] bg-gradient-to-b from-slate-900 via-blue-900 to-indigo-900 border-r border-blue-800/50 flex flex-col shadow-2xl">
      {/* Logo Area */}
      <div className="py-6 px-6 border-b border-blue-800/50">
        <div className="flex items-center gap-3 mb-2">
          <div className="w-12 h-12 rounded-xl bg-white/10 backdrop-blur-sm flex items-center justify-center shadow-lg shadow-blue-500/30 p-1.5">
            <Image
              src="/logo.png"
              alt="IU-MiCert Logo"
              width={40}
              height={40}
              className="w-full h-full object-contain"
            />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">IU-MiCert</h1>
            <p className="text-xs text-blue-200/80 font-medium">
              Issuer Dashboard
            </p>
          </div>
        </div>
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
                      ? "bg-gradient-to-r from-blue-500 to-indigo-500 text-white shadow-lg shadow-blue-500/50"
                      : "text-blue-100/70 hover:bg-white/10 hover:text-white"
                  )}
                >
                  <Icon
                    className={cn(
                      "w-5 h-5 transition-colors",
                      isActive
                        ? "text-white"
                        : "text-blue-200/70 group-hover:text-white"
                    )}
                  />
                  <span
                    className={cn(
                      "font-semibold text-sm transition-colors",
                      isActive ? "text-white" : ""
                    )}
                  >
                    {item.label}
                  </span>
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Footer */}
      <div className="px-4 py-4 border-t border-blue-800/50">
        <p className="text-[11px] text-blue-300/50">Version 1.0.0</p>
      </div>
    </aside>
  );
}
