import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
    state: {
        webSocket: null,
        wsConnected: false,
        lastMessage: null,
        systemLogs: [] // 添加系统日志数组
    },
    mutations: {
        setWebSocket(state, ws) {
            state.webSocket = ws
        },
        setWsConnected(state, status) {
            state.wsConnected = status
        },
        setLastMessage(state, message) {
            state.lastMessage = message
        },
        clearWebSocket(state) {
            if (state.webSocket) {
                state.webSocket.close()
                state.webSocket = null
            }
            state.wsConnected = false
        },
        setInitialLogs(state, logs) {
            state.systemLogs = logs
        },
        addLog(state, log) {
            state.systemLogs = [log, ...state.systemLogs]
        }
    },
    actions: {
        initWebSocket({ commit, state }, url) {
            return new Promise((resolve, reject) => {
                // 如果已经存在连接并且处于连接状态，直接返回
                if (state.webSocket && state.wsConnected) {
                    resolve(state.webSocket)
                    return
                }

                // 检查登录状态
                const isLoggedIn = localStorage.getItem('isLoggedIn')
                if (isLoggedIn !== 'true') {
                    reject(new Error('未登录状态'))
                    return
                }

                // 如果存在旧连接，先关闭
                if (state.webSocket) {
                    state.webSocket.close()
                }

                try {
                    const ws = new WebSocket(url)

                    ws.onopen = () => {
                        commit('setWebSocket', ws)
                        commit('setWsConnected', true)
                        console.log('WebSocket connected')
                        resolve(ws)
                    }

                    ws.onmessage = (event) => {
                        try {
                            const data = JSON.parse(event.data)
                            // 转换成日志格式，确保类型值符合 Element Plus 要求
                            const log = {
                                id: Date.now(),
                                level: data.level || 'Info',
                                levelType: data.levelType, // 使用转换函数
                                service: data.service || 'system',
                                message: data.message || '系统消息',
                                detail: data.detail || '',
                                timestamp: data.timestamp
                            }
                            commit('addLog', log)
                            commit('setLastMessage', data)
                        } catch (error) {
                            console.error('解析消息失败:', error)
                        }
                    }

                    ws.onclose = () => {
                        commit('setWsConnected', false)
                        console.log('WebSocket disconnected')
                    }

                    ws.onerror = (error) => {
                        console.error('WebSocket error:', error)
                        commit('setWsConnected', false)
                        reject(error)
                    }

                } catch (error) {
                    console.error('Failed to create WebSocket:', error)
                    reject(error)
                }
            })
        },
        sendWsMessage({ state }, message) {
            if (state.webSocket && state.wsConnected) {
                state.webSocket.send(JSON.stringify(message))
            }
        },
        closeWebSocket({ commit }) {
            commit('clearWebSocket')
        },
        async fetchHistoryLogs({ commit, state }) {
            // 如果已有日志数据，直接返回
            if (state.systemLogs.length > 0) {
                return state.systemLogs
            }

            try {
                const response = await axios.get('http://127.0.0.1:9999/api/logs/recent', {
                    params: {
                        limit: 50
                    }
                })
                if (response.data && Array.isArray(response.data)) {
                    const sortedLogs = response.data.sort((a, b) =>
                        new Date(b.timestamp) - new Date(a.timestamp)
                    )
                    commit('setInitialLogs', sortedLogs)
                    return sortedLogs
                }
            } catch (error) {
                console.error('获取历史日志失败:', error)
                throw error
            }
        }
    },
    getters: {
        isWsConnected: state => state.wsConnected,
        getLastMessage: state => state.lastMessage,
        getSystemLogs: state => state.systemLogs
    }
})
