"use client";

import { Inter } from "next/font/google";
import "./globals.css";
import { Web3Provider } from "@/providers/Web3Provider";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { usePathname } from "next/navigation";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const isLoginPage = pathname === "/login";

  return (
    <html lang="en">
      <body className={`${inter.variable} font-sans antialiased`}>
        <Web3Provider>
          {isLoginPage ? (
            // Login page renders without DashboardLayout
            children
          ) : (
            // All other pages use DashboardLayout
            <DashboardLayout>{children}</DashboardLayout>
          )}
        </Web3Provider>
      </body>
    </html>
  );
}
