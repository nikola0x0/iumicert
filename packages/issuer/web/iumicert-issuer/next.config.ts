import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Suppress CSS preload warnings in development
  experimental: {
    optimizePackageImports: ['@headlessui/react', '@heroicons/react'],
  },
  webpack: (config: any) => {
    // Ignore optional pino dependencies
    config.resolve.alias = {
      ...config.resolve.alias,
      'pino-pretty': false,
    };

    // Development watch options
    if (process.env.NODE_ENV === 'development') {
      config.watchOptions = {
        poll: 1000,
        aggregateTimeout: 300,
      };
    }

    return config;
  },
};

export default nextConfig;
