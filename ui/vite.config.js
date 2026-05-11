import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Vite config — proxy /api/* to the live Ambient Weather backend in dev.
// Override the target via VITE_API_TARGET if your station is hosted elsewhere.
export default defineConfig(({ mode }) => ({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: 'http://localhost:8000',
        changeOrigin: true,
        ws: true,
      },
    },
  },
  build: {
    outDir: "dist",
    sourcemap: true,
  },
}));
