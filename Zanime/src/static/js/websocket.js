const state = {
  socket: null, // 存储 WebSocket 实例
  isConnected: false, // WebSocket 是否连接
  user: null, // 存储当前用户信息
};

const mutations = {
  SET_SOCKET(state, socket) {
    state.socket = socket;
  },
  SET_CONNECTED(state, isConnected) {
    state.isConnected = isConnected;
  },
  SET_USER(state, user) {
    state.user = user;
  },
};

const actions = {
  connectWebSocket({ commit }, pushUrl) {
    const socket = new WebSocket('/ws');
    commit('SET_SOCKET', socket);

    socket.onopen = () => {
      commit('SET_CONNECTED', true);
      console.log('WebSocket 连接已打开');
    };

    socket.onmessage = (event) => {
      console.log('收到消息:', event.data);
    };

    socket.onclose = () => {
      commit('SET_CONNECTED', false);
      console.log('WebSocket 连接已关闭');
    };
  },

  sendMessage({ state }, message) {
    if (state.socket && state.isConnected) {
      state.socket.send(JSON.stringify(message));
    } else {
      console.log('WebSocket 未连接');
    }
  },

  disconnectWebSocket({ commit }) {
    if (state.socket) {
      state.socket.close();
      commit('SET_CONNECTED', false);
      console.log('WebSocket 已关闭');
    }
  },
};

const getters = {
  socket: (state) => state.socket,
  isConnected: (state) => state.isConnected,
  user: (state) => state.user,
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
};