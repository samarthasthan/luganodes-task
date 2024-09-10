/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/deposits",
        destination: "http://3.7.73.40:8000/deposits?page=1&limit=10",
      },
    ];
  },
};

export default nextConfig;
