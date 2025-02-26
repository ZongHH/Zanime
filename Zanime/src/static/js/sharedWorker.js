let connections = []; // 存储所有连接的页面端口
let ws = null; // WebSocket 实例

// 监听页面的连接
self.onconnect = (event) => {
  const port = event.ports[0]; // 获取当前页面的通信端口
  connections.push(port);

  if (!ws || ws.readyState === WebSocket.CLOSED) {
    ws = new WebSocket('/conn/ws')

    ws.onopen = () => {
        console.log('WebSocket connected in SharedWorker');
        broadcast({ type: 'WS_CONNECTED' });
      };

      ws.onmessage = (event) => {
        // 将 WebSocket 消息广播给所有连接的页面
        let jsonData;
        try {
            jsonData = JSON.parse(event.data);
        } catch (error) {
            console.error('解析WebSocket消息为JSON时出错：', error);
        }
        broadcast({ type: 'WS_MESSAGE', data: jsonData });
      };

      ws.onclose = () => {
        console.log('WebSocket closed');
        broadcast({ type: 'WS_DISCONNECTED' });
      };
  }

  // 监听页面发来的消息
  port.onmessage = (msgEvent) => {
    const { type, payload } = msgEvent.data;

    if (type === 'SEND_MESSAGE') {
      // 发送消息到 WebSocket
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(payload);
      }
    }
  };

  // 移除关闭的页面
  port.onclose = () => {
    connections = connections.filter((conn) => conn !== port);

    // 如果没有页面连接到 SharedWorker，关闭 WebSocket
    if (connections.length === 0) {
        console.log('No more connections, closing WebSocket.');
        if (ws) {
            ws.close();
            ws = null; // 确保重新连接时创建新的 WebSocket 实例
        }
    }
  };

  port.start(); // 开启通信
};

// 将消息广播给所有页面
function broadcast(message) {
  connections.forEach((port) => {
    port.postMessage(message);
  });
}
