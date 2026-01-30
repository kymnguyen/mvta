import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api/auth': {
        target: 'http://localhost:50000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/auth/, '/api/v1'),
      },
      '/api/vehicle': {
        target: 'http://localhost:50001',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/vehicle/, '/api/v1'),
      },
      '/api/tracking': {
        target: 'http://localhost:50002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/tracking/, '/api/v1'),
      }
    }
  }
})
