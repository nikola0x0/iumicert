import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Performance optimizations
  compiler: {
    removeConsole: process.env.NODE_ENV === "production",
  },

  // Optimize images
  images: {
    formats: ["image/avif", "image/webp"],
  },

  // Reduce bundle size
  experimental: {
    optimizePackageImports: ["lucide-react", "gsap"],
  },

  // Production optimizations
  reactStrictMode: true,
};

export default nextConfig;
