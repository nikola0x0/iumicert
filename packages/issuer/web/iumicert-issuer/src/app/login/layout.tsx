import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Login - IU-MiCert Issuer",
  description: "Sign in to the IU-MiCert Issuer Dashboard",
};

export default function LoginLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // No DashboardLayout wrapper - login page is standalone
  return <>{children}</>;
}
