"use client";

import { usePathname } from "next/navigation";
import Header from "./Header";
import Footer from "./Footer";

interface ConditionalLayoutProps {
  children: React.ReactNode;
}

const ConditionalLayout: React.FC<ConditionalLayoutProps> = ({ children }) => {
  const pathname = usePathname();

  // Routes that should use dashboard layout (no header/footer)
  const dashboardRoutes = [
    "/revoke",
    "/search",
    "/analytics",
    "/settings",
    "/verifier",
    "/selective",
  ];
  const isDashboardRoute = dashboardRoutes.includes(pathname);

  return (
    <>
      {/* Full-screen content */}
      <main className="absolute inset-0">{children}</main>

      {/* Conditionally render Header and Footer */}
      {!isDashboardRoute && (
        <>
          <Header />
          <Footer />
        </>
      )}
    </>
  );
};

export default ConditionalLayout;
