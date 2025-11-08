import React from "react";
import { Link } from "next-view-transitions";
import {
  Shield,
  X,
  Settings,
  Home,
  Book,
  Lock,
} from "lucide-react";
import Image from "next/image";
import SystemStatus from "./ui/SystemStatus";

interface DashboardSidebarProps {
  activeSection?: string;
}

const DashboardSidebar: React.FC<DashboardSidebarProps> = ({
  activeSection,
}) => {
  const navigation = [
    {
      id: "home",
      name: "Home",
      href: "/",
      icon: Home,
      description: "Landing page",
    },
    {
      id: "verify",
      name: "Verify",
      href: "/verifier",
      icon: Shield,
      description: "Verify credentials",
    },
    {
      id: "selective",
      name: "Selective",
      href: "/selective",
      icon: Lock,
      description: "Create filtered receipt",
    },
    {
      id: "revoke",
      name: "Revoke",
      href: "/revoke",
      icon: X,
      description: "Contact issuer",
    },
  ];

  const secondaryNavigation = [
    {
      id: "documentation",
      name: "Documentation",
      href: "https://nikolaempire.gitbook.io/iu-micert/",
      icon: Book,
      external: true,
      description: "System documentation",
    },
    {
      id: "settings",
      name: "Settings",
      href: "/settings",
      icon: Settings,
      description: "App settings",
    },
  ];

  return (
    <div className="w-64 flex-shrink-0 flex flex-col h-full">
      {/* Logo Section */}
      <div className="p-6">
        <Link href="/" className="flex items-center space-x-3 group">
          <div className="relative">
            <Image
              src="/logo.svg"
              alt="IU-MiCert Logo"
              width={40}
              height={40}
              className="object-contain transition-all duration-300 ease-in-out group-hover:scale-110 group-hover:rotate-3"
            />
          </div>
          <div className="flex flex-col">
            <h1 className="text-xl font-bold text-white font-crimson">
              IU-MiCert
            </h1>
            <p className="text-purple-200 text-xs font-inter">Dashboard</p>
          </div>
        </Link>
      </div>

      {/* Main Navigation */}
      <div className="flex-1 px-4">
        <div className="mb-4">
          <h3 className="text-xs font-semibold text-white/60 uppercase tracking-wider px-3 py-2 font-inter">
            Main
          </h3>
        </div>

        {navigation.map((item) => {
          const isActive = activeSection === item.id;

          return (
            <Link key={item.id} href={item.href} className="block mb-2">
              <div
                className={`flex items-center px-3 py-3 rounded-lg transition-all duration-75 group cursor-pointer border ${
                  isActive
                    ? "bg-white/15 text-white shadow-sm border-white/20"
                    : "text-white/70 hover:bg-white/10 hover:text-white border-transparent"
                }`}
              >
                <item.icon
                  className={`w-5 h-5 mr-3 transition-colors duration-75 ${
                    isActive
                      ? "text-blue-300"
                      : "text-white/60 group-hover:text-blue-300"
                  }`}
                />
                <div className="flex-1 min-w-0">
                  <div
                    className={`font-medium font-inter ${
                      isActive ? "text-white" : ""
                    }`}
                  >
                    {item.name}
                  </div>
                  <div className="text-xs text-white/50 font-inter">
                    {item.description}
                  </div>
                </div>
              </div>
            </Link>
          );
        })}

        {/* Divider */}
        <div className="my-6 mx-3">
          <div className="border-t border-white/10"></div>
        </div>

        {/* Secondary Navigation */}
        <div className="mb-4">
          <h3 className="text-xs font-semibold text-white/60 uppercase tracking-wider px-3 py-2 font-inter">
            Other
          </h3>
        </div>

        {secondaryNavigation.map((item) => {
          const isActive = activeSection === item.id;

          if (item.external) {
            return (
              <a
                key={item.id}
                href={item.href}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center px-3 py-3 rounded-lg transition-all duration-75 group cursor-pointer text-white/70 hover:bg-white/10 hover:text-white mb-2 border border-transparent"
              >
                <item.icon className="w-5 h-5 mr-3 text-white/60 group-hover:text-blue-300 transition-colors duration-75" />
                <div className="flex-1 min-w-0">
                  <div className="font-medium font-inter">{item.name}</div>
                  <div className="text-xs text-white/50 font-inter">
                    {item.description}
                  </div>
                </div>
                <div className="w-4 h-4 text-white/40">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                    />
                  </svg>
                </div>
              </a>
            );
          }

          return (
            <Link key={item.id} href={item.href} className="block mb-2">
              <div
                className={`flex items-center px-3 py-3 rounded-lg transition-all duration-75 group cursor-pointer border ${
                  isActive
                    ? "bg-white/15 text-white shadow-sm border-white/20"
                    : "text-white/70 hover:bg-white/10 hover:text-white border-transparent"
                }`}
              >
                <item.icon
                  className={`w-5 h-5 mr-3 transition-colors duration-75 ${
                    isActive
                      ? "text-blue-300"
                      : "text-white/60 group-hover:text-blue-300"
                  }`}
                />
                <div className="flex-1 min-w-0">
                  <div
                    className={`font-medium font-inter ${
                      isActive ? "text-white" : ""
                    }`}
                  >
                    {item.name}
                  </div>
                  <div className="text-xs text-white/50 font-inter">
                    {item.description}
                  </div>
                </div>
              </div>
            </Link>
          );
        })}
      </div>

      {/* Status Footer */}
      <div className="p-4">
        <SystemStatus status="development" />
      </div>
    </div>
  );
};

export default DashboardSidebar;
