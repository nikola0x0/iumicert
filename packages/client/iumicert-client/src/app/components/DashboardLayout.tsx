import React from "react";
import AnimatedBackground from "./AnimatedBackground";
import DashboardSidebar from "./DashboardSidebar";

interface DashboardLayoutProps {
  children: React.ReactNode;
  activeSection?: string;
}

const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  children,
  activeSection,
}) => {
  return (
    <div className="h-full w-full relative overflow-hidden">
      {/* Animated Background */}
      <AnimatedBackground
        gradient="from-slate-900 via-purple-900 to-blue-900"
        className="transition-all duration-1000"
      />

      {/* Dashboard Content */}
      <div className="relative z-10 h-full flex">
        {/* Sidebar */}
        <DashboardSidebar activeSection={activeSection} />

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col min-h-0">
          {/* Content Area */}
          <div className="flex-1 px-6 pb-6 min-h-0">
            <div className="h-full">{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardLayout;
