import { defineConfig } from "vite";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 3000,
  },
  preview: {
    allowedHosts: ["frontend-production-f7d0.up.railway.app"],
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@poojanbhardwaj/openapi": path.resolve(
        __dirname,
        "../../packages/openapi/src"
      ),
      "@poojanbhardwaj/zod": path.resolve(__dirname, "../../packages/zod/src"),
    },
  },
});
