import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  define: {
    'process.env': {
      NODE_ENV: JSON.stringify(process.env.NODE_ENV || 'production')
    }
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    lib: {
      entry: resolve(__dirname, 'src/main.ts'),
      name: 'GrowDeskWidget',
      fileName: (format) => `growdesk-widget.${format}.js`,
    },
    rollupOptions: {
      external: ['vue'],
      output: {
        globals: {
          vue: 'Vue'
        },
        preserveModules: false,
      },
    },
    sourcemap: true,
    minify: false,
    target: 'es2018',
    assetsInlineLimit: 4096
  },
  css: {
    devSourcemap: true,
  },
  server: {
    host: true,
    port: 3000,
    cors: true
  },
}); 