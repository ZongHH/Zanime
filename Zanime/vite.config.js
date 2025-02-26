import { fileURLToPath, URL } from 'node:url' // Node.js 文件路径处理模块

import { defineConfig } from 'vite' // Vite 核心配置函数
import vue from '@vitejs/plugin-vue' // Vue 3 单文件组件支持插件
import vueDevTools from 'vite-plugin-vue-devtools' // Vue DevTools 调试工具集成插件

/**
 * Vite 配置文件
 * @see https://vite.dev/config/
 */
export default defineConfig({
  /**
   * 插件配置
   * 1. vue(): 处理 Vue 单文件组件
   * 2. vueDevTools(): 开发阶段集成 Vue DevTools
   */
  plugins: [
    vue({
      template: {
        compilerOptions: {
          // SVG 相关元素处理配置
          // 将 <clippath> 标记为原生元素，避免 Vue 的解析警告
          isCustomElement: tag => tag === 'clippath'
        }
      }
    }),
    vueDevTools(), // 开发工具插件（仅在开发环境生效）
  ],

  /**
   * 模块解析配置
   */
  resolve: {
    /**
     * 路径别名配置
     * - '@': 指向项目源码目录 src/
     * 使用示例：import MyComponent from '@/components/MyComponent.vue'
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
    port: 5173,        // 开发服务器端口（默认5173）
    https: false,      // 是否启用 HTTPS（默认false）
    open: false,       // 启动时自动打开浏览器（默认false）

    /**
     * 代理配置
     * 用于解决开发环境跨域问题
     */
    proxy: {
      // REST API 请求代理
      '/api': {
        target: 'http://127.0.0.1:9092', // 后端服务地址
        changeOrigin: true, // 修改请求头中的 Host 为目标地址
        rewrite: path => path // 路径重写函数（此处保留原始路径）
      },

      // WebSocket 代理配置
      '/conn': {
        target: 'ws://127.0.0.1:9092',  // WebSocket 服务地址
        ws: true,  // 启用 WebSocket 协议代理
        changeOrigin: true,  // 修改请求头 Origin 字段
        withCredentials: true,  // 允许跨域携带 cookies
        /**
         * 注意：WebSocket 代理需要同时配置：
         * - 协议头 ws:// 或 wss://
         * - ws: true 显式声明代理 WebSocket
         */
      }
    }
  },
})
