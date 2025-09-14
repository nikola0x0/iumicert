import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Suppress CSS preload warnings in development
  experimental: {
    optimizePackageImports: ['@headlessui/react', '@heroicons/react'],
  },
  // Reduce preloading in development
  ...(process.env.NODE_ENV === 'development' && {
    webpack: (config: any) => {
      config.watchOptions = {
        poll: 1000,
        aggregateTimeout: 300,
      };
      return config;
    },
  }),
};

export default nextConfig;
