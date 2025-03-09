import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

/**
 * Vite 配置文件
 * @see https://vite.dev/config/
 */
export default defineConfig({
  /**
   * 插件配置
   * - vue(): Vue 3 单文件组件支持
   * - vueDevTools(): Vue DevTools 集成
   */
  plugins: [
    vue(), // Vue 3 核心插件
    vueDevTools(), // 开发工具调试插件
  ],

  /**
   * 模块解析配置
   */
  resolve: {
    /**
     * 路径别名配置
     * '@': 指向项目源码目录 src/
     */
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },

  /**
   * 开发服务器配置
   */
  server: {
    host: '127.0.0.1', // 指定监听地址（默认localhost）
    port: 5174,        // 指定开发服务器端口
    open: false,       // 是否自动打开浏览器（默认false）
    https: false,      // 是否启用 HTTPS（默认false）

    /**
     * 代理配置
     * 将 /api 前缀的请求转发到后端服务器
     */
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:9999',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    },

    /**
     * 其他可用配置示例：
     * cors: true       // 启用 CORS（默认true）
     * strictPort: true // 端口被占用时直接退出（默认false）
     */
  }
})
