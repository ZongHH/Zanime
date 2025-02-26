import { createStore } from 'vuex'; // 导入 Vuex 的 createStore 方法，用于创建 Vuex Store
import websocket from '@/static/js/websocket.js'; // 导入自定义的 websocket 模块

export default createStore({
  modules: {
    websocket, // 在 Vuex Store 中注册 websocket 模块
  },
});
