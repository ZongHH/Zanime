<template>
  <router-view />
  <PushNotification v-if="showPush" :title="title" :message="message" :duration="5000" @closed="showPush = false" />
</template>

<script>
import PushNotification from './components/PushNotification.vue';

/**
 * WebSocket 连接管理类
 * 提供WebSocket连接的建立、重连、心跳维护和关闭等功能
 */
class WebSocketManager {
  /**
   * 构造函数
   * @param {string} url - WebSocket服务器地址
   * @param {Function} messageCallback - 消息到达时的回调函数
   */
  constructor(url, messageCallback) {
    this.url = url;                    // WebSocket服务端地址
    this.messageCallback = messageCallback; // 消息回调处理函数
    this.ws = null;                    // WebSocket实例
    this.reconnectAttempts = 0;        // 当前重连尝试次数
    this.maxReconnectAttempts = 5;     // 最大重连尝试次数
    this.reconnectDelay = 3000;        // 重连延迟时间（毫秒）
    this.heartbeatInterval = 30000;    // 心跳包发送间隔（毫秒）
    this.isManualClose = false;        // 是否手动关闭连接标志
  }

  /**
   * 建立WebSocket连接
   * 初始化WebSocket实例并绑定事件处理器
   */
  connect() {
    this.ws = new WebSocket(this.url);

    // 连接成功回调
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;      // 重置重连计数器
      this.startHeartbeat();           // 启动心跳检测
    };

    // 消息接收回调
    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.messageCallback(data);    // 调用外部传入的消息处理回调
      } catch (error) {
        console.error('WebSocket message parse error:', error);
      }
    };

    // 连接关闭回调
    this.ws.onclose = (event) => {
      console.log(`WebSocket closed (code ${event.code})`);
      this.stopHeartbeat();            // 停止心跳检测
      // 非主动关闭且未达最大重连次数时尝试重连
      if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
        setTimeout(() => this.reconnect(), this.reconnectDelay);
      }
    };

    // 错误处理回调
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  /**
   * 执行重连操作
   * 增加重连计数器并重新建立连接
   */
  reconnect() {
    this.reconnectAttempts++;
    console.log(`Reconnecting attempt ${this.reconnectAttempts}`);
    this.connect();
  }

  /**
   * 启动心跳检测机制
   * 定时发送心跳包保持连接活跃
   */
  startHeartbeat() {
    this.heartbeatTimer = setInterval(() => {
      if (this.ws.readyState === WebSocket.OPEN) {
        // 发送心跳包（需服务端支持心跳协议）
        this.ws.send(JSON.stringify({ type: 'ping' }));
      }
    }, this.heartbeatInterval);
  }

  /**
   * 停止心跳检测
   * 清除心跳定时器
   */
  stopHeartbeat() {
    clearInterval(this.heartbeatTimer);
  }

  /**
   * 主动关闭连接
   * 设置手动关闭标志并关闭WebSocket
   */
  close() {
    this.isManualClose = true;  // 标记为手动关闭
    if (this.ws) {
      this.ws.close();
    }
  }
}

export default {
  name: "App",
  data() {
    return {
      showPush: false,
      title: "",
      message: "",
      wsManager: null
    };
  },
  components: {
    PushNotification,
  },
  methods: {
    handleWebSocketMessage(data) {
      if (data.type === 'WS_MESSAGE') {
        this.title = `${data.send_username} 回复你: `;
        this.message = data.content;
        this.showPush = true;
      } else if (data.type === 'WS_COMMENT_LIKE') {
        this.title = `${data.user_name} 点赞了你的评论`;
        this.message = `"${data.post_title}"`;
        this.showPush = true;
      }
    }
  },
  mounted() {
    // 初始化 WebSocket 连接
    this.wsManager = new WebSocketManager(
      `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/conn/ws`,
      this.handleWebSocketMessage.bind(this)
    );
    this.wsManager.connect();

    // 监听网络恢复事件
    window.addEventListener('online', () => {
      if (!this.wsManager.ws || this.wsManager.ws.readyState !== WebSocket.OPEN) {
        this.wsManager.connect();
      }
    });
  },
  beforeDestroy() {
    if (this.wsManager) {
      this.wsManager.close();
    }
  }
};
</script>

<style></style>
