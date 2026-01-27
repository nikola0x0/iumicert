import type { NextConfig } from "next";

// Fix for Node.js 25+ which has experimental localStorage that throws
// "Cannot initialize local storage without a `--localstorage-file` path"
// We need to replace it with a mock implementation for SSR
if (typeof globalThis !== 'undefined') {
  const mockStorage = {
    getItem: () => null,
    setItem: () => {},
    removeItem: () => {},
    clear: () => {},
    key: () => null,
    length: 0,
  };

  try {
    // Test if localStorage throws
    if (typeof globalThis.localStorage !== 'undefined') {
      globalThis.localStorage.getItem('test');
    }
  } catch {
    // If it throws, replace with mock
    // @ts-ignore - replacing broken localStorage in Node.js 25+
    globalThis.localStorage = mockStorage;
    // @ts-ignore - also mock sessionStorage
    globalThis.sessionStorage = mockStorage;
  }
}

const nextConfig: NextConfig = {
  // Performance optimizations
  compiler: {
    removeConsole: process.env.NODE_ENV === "production",
    styledComponents: true,
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
