import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    react(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        manualChunks(id: string) {
          // Monaco editor — largest chunk, always isolated
          if (id.includes('monaco-editor') || id.includes('@monaco-editor')) {
            return 'vendor-monaco';
          }
          // Recharts / react-is
          if (id.includes('recharts') || id.includes('react-is')) {
            return 'vendor-charts';
          }
          // Framer Motion
          if (id.includes('framer-motion')) {
            return 'vendor-motion';
          }
          // Forms / validation
          if (id.includes('react-hook-form') || id.includes('@hookform') || id.includes('zod')) {
            return 'vendor-forms';
          }
          // Networking + toasts
          if (id.includes('axios') || id.includes('sonner')) {
            return 'vendor-net';
          }
          // React core
          if (id.includes('react-dom') || id.includes('react-router')) {
            return 'vendor-react';
          }
          // All other node_modules
          if (id.includes('node_modules')) {
            return 'vendor-misc';
          }
        },
      },
    },
  },
})
