import { fileURLToPath, URL } from 'node:url'
import { resolve } from 'path'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
      alias: {
          '@': fileURLToPath(new URL('./src', import.meta.url)),
          '@idx': fileURLToPath(new URL('./src/pages/index', import.meta.url)),
          '@adm': fileURLToPath(new URL('./src/pages/admin', import.meta.url))
      }
  },
  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        admin: resolve(__dirname, 'admin.html'),
      },
    },
    outDir: '../public'
  }
})
