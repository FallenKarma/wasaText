import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: true, // <-- allow external access (i.e., from host to container)
    port: 4173, // <-- ensure this matches your Docker EXPOSE and docker-compose
    strictPort: true, // <-- fail if port is taken, for clarity
  },
})
