import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: 'localhost', // 明确指定监听 localhost，确保可以通过 localhost:5173 访问
    port: 5173,
    open: true,
    strictPort: true, // 如果端口被占用，报错而不是尝试其他端口
    hmr: {
      host: 'localhost', // HMR 也使用 localhost
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'),
      },
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        // 注意：Sass 的弃用警告（legacy-js-api 和 @import）是正常的
        // 这些警告不影响功能，只是提示未来版本的变化
        // 我们已经将 @import 改为 @use，但 legacy-js-api 警告来自 Sass 编译器内部
        // 无法通过配置完全消除，但不影响项目运行
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'element-plus': ['element-plus', '@element-plus/icons-vue'],
          'vue-flow': ['@vue-flow/core'],
        },
      },
    },
  },
})

