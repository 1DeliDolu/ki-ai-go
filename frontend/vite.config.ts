import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: "localhost", // Force localhost instead of 0.0.0.0
    port: 5173,
    strictPort: true, // Fail if port is already in use
    open: true, // Automatically open browser
    cors: true,
    proxy: {
      "/api": {
        target: "http://localhost:8082",
        changeOrigin: true,
        secure: false,
        configure: (proxy) => {
          proxy.on("error", (err) => {
            console.error("❌ Proxy error:", err.message);
          });
          proxy.on("proxyReq", (_, req) => {
            console.log("🔄 Proxying:", req.method, req.url);
          });
          proxy.on("proxyRes", (proxyRes, req) => {
            console.log("✅ Proxy response:", proxyRes.statusCode, req.url);
          });
        },
      },
    },
  },
  build: {
    outDir: "dist",
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
});
