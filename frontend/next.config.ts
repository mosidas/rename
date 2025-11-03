import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // Wails で静的配信できるように静的書き出し
  output: 'export',
  images: { unoptimized: true },
};

export default nextConfig;
