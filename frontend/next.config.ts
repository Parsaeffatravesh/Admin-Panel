import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  allowedDevOrigins: [
    "*.replit.app",
    "*.replit.dev",
    "*.janeway.replit.dev",
    "ff9f88a6-a464-47f3-8adf-2265615ae524-00-1mbe8sgsut9v9.sisko.replit.dev",
    "localhost:5000",
    "127.0.0.1:5000",
  ],
  compress: true,
  poweredByHeader: false,
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ];
  },
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on',
          },
        ],
      },
      {
        source: '/api/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'private, max-age=60',
          },
        ],
      },
    ];
  },
};

export default nextConfig;
