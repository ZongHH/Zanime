<template>
    <transition name="modal">
        <div class="history-modal-overlay" @click.self="$emit('close')" v-if="show">
            <transition name="modal-content">
                <div class="history-modal-content">
                    <div class="modal-header">
                        <h2>观看历史</h2>
                        <button class="close-btn" @click="$emit('close')">
                            <i class="fas fa-times"></i>
                        </button>
                    </div>

                    <div class="modal-body" @scroll="handleScroll" ref="modalBody">
                        <div class="history-list">
                            <!-- 按日期分组显示 -->
                            <div v-for="(group, date) in groupedHistory" :key="date" class="history-group">
                                <div class="date-header">{{ formatDateHeader(date) }}</div>
                                <div class="history-grid">
                                    <div v-for="anime in group" :key="anime.id" class="history-item"
                                        @click="continueWatching(anime)">
                                        <div class="anime-cover">
                                            <img :src="anime.cover_image_url" :alt="anime.title">
                                            <div class="progress-overlay">
                                                <div class="progress-bar" :style="{ width: anime.progress + '%' }">
                                                </div>
                                            </div>
                                        </div>
                                        <div class="anime-info">
                                            <h3>{{ anime.title }}</h3>
                                            <div class="episode-info">
                                                <span class="video-name">{{ anime.video_name }}</span>
                                                <div class="episode-progress">
                                                    <span class="episode">{{ anime.episode }}</span>
                                                    <span class="progress-text">{{ formatProgress(anime.progress)
                                                        }}</span>
                                                </div>
                                            </div>
                                            <p class="watch-time">{{ formatTime(anime.updated_at) }}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- 加载状态 -->
                        <div v-if="isLoading" class="loading-more">
                            <i class="fas fa-spinner fa-spin"></i>
                            加载中...
                        </div>

                        <!-- 没有更多数据提示 -->
                        <div v-if="!hasMore && historyList.length > 0" class="no-more">
                            没有更多记录了
                        </div>
                    </div>
                </div>
            </transition>
        </div>
    </transition>
</template>

<script>
import axios from 'axios';

export default {
    name: 'HistoryModal',
    data() {
        return {
            historyList: [],
            currentPage: 1,
            isLoading: false,
            hasMore: true
        }
    },
    props: {
        userInfo: {
            type: Object,
            required: true
        },
        show: {
            type: Boolean,
            default: true
        }
    },
    computed: {
        groupedHistory() {
            const groups = {};
            this.historyList.forEach(anime => {
                const date = new Date(anime.updated_at).toLocaleDateString();
                if (!groups[date]) {
                    groups[date] = [];
                }
                groups[date].push(anime);
            });
            return groups;
        }
    },
    methods: {
        async fetchHistory() {
            if (this.isLoading || !this.hasMore) return;

            try {
                this.isLoading = true;

                const response = await axios.get('/api/watch-history', {
                    params: {
                        user_id: this.userInfo.user_id,
                        page: this.currentPage,
                    }
                });

                if (response.data.code == 200) {
                    const newHistory = response.data.progress;

                    // 如果没有新数据返回,说明已经加载完所有历史
                    if (!newHistory || newHistory.length === 0) {
                        this.hasMore = false;
                    } else {
                        // 将新数据添加到现有列表
                        this.historyList = [...this.historyList, ...newHistory];
                        this.currentPage++;
                    }
                } else {
                    console.error(`获取观看历史失败: `, response.data.message);
                }
            } catch (error) {
                console.error('获取观看历史失败:', error);
            } finally {
                this.isLoading = false;
            }
        },

        handleScroll(e) {
            const { scrollTop, scrollHeight, clientHeight } = e.target;
            // 当距离底部100px时加载更多
            if (scrollHeight - scrollTop - clientHeight < 100) {
                this.fetchHistory();
            }
        },

        continueWatching(anime) {
            this.$router.push({
                path: '/moviesDetail',
                query: {
                    videoId: anime.video_id,
                    episode: anime.episode
                }
            });
        },

        formatProgress(progress) {
            return `${Math.round(progress)}%`;
        },

        formatDateHeader(date) {
            const today = new Date().toLocaleDateString();
            const yesterday = new Date(Date.now() - 86400000).toLocaleDateString();

            if (date === today) return '今天';
            if (date === yesterday) return '昨天';
            return date;
        },

        formatTime(date) {
            return new Date(date).toLocaleTimeString('zh-CN', {
                hour: '2-digit',
                minute: '2-digit'
            });
        }
    },
    mounted() {
        this.fetchHistory();
    }
}
</script>

<style scoped>
.history-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 20px;
    animation: fade-in 0.3s;
}

.history-modal-content {
    background: linear-gradient(145deg, rgba(30, 30, 30, 0.95), rgba(20, 20, 20, 0.95));
    border-radius: 20px;
    width: 100%;
    max-width: 800px;
    height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.05);
    animation: slide-in 0.3s;
}

@keyframes fade-in {
    from {
        opacity: 0;
    }

    to {
        opacity: 1;
    }
}

@keyframes slide-in {
    from {
        opacity: 0;
        transform: translateY(-30px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.modal-header {
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.modal-header h2 {
    color: white;
    margin: 0;
    font-size: 1.5rem;
}

.close-btn {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.6);
    font-size: 1.2rem;
    cursor: pointer;
    padding-left: 10px;
    padding-right: 10px;
    padding-top: 6px;
    padding-bottom: 6px;
    border-radius: 50%;
    transition: all 0.3s ease;
}

.close-btn:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    scrollbar-width: thin;
    /* Firefox */
    scrollbar-color: rgb(255 255 255 / 50%) rgba(255, 255, 255, 0.1);
    /* Firefox */
}

/* Webkit (Chrome, Safari, Edge) 滚动条样式 */
.modal-body::-webkit-scrollbar {
    width: 6px;
}

.modal-body::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb {
    background: linear-gradient(to bottom,
            rgba(204, 0, 0, 0.5),
            rgba(139, 0, 0, 0.5));
    border-radius: 3px;
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-body::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(to bottom,
            rgba(204, 0, 0, 0.7),
            rgba(139, 0, 0, 0.7));
}

/* 当滚动条处于活动状态时的样式 */
.modal-body::-webkit-scrollbar-thumb:active {
    background: linear-gradient(to bottom,
            rgba(204, 0, 0, 0.8),
            rgba(139, 0, 0, 0.8));
}

/* 确保在暗色背景下滚动条可见 */
.modal-body::-webkit-scrollbar-corner {
    background: transparent;
}

.history-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.history-group {
    margin-bottom: 24px;
}

.date-header {
    color: white;
    font-size: 1.1rem;
    font-weight: 500;
    margin-bottom: 16px;
    padding-left: 8px;
    border-left: 3px solid #cc0000;
}

.history-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 16px;
}

.history-item {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
}

.history-item:hover {
    background: rgba(255, 255, 255, 0.08);
    transform: translateY(-2px);
}

.anime-cover {
    width: 100%;
    height: 0;
    padding-bottom: 140%;
    /* 调整为 5:7 的比例 */
    position: relative;
}

.anime-cover img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.progress-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: rgba(0, 0, 0, 0.5);
}

.progress-bar {
    height: 100%;
    background: #cc0000;
    transition: width 0.3s ease;
}

.anime-info {
    padding: 8px;
}

.anime-info h3 {
    margin: 0;
    color: white;
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.episode-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.video-name {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.8rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.episode-progress {
    display: flex;
    gap: 6px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.75rem;
}

.watch-time {
    margin-top: 4px;
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.75rem;
}

.loading-more {
    text-align: center;
    padding: 20px 0;
    color: rgba(255, 255, 255, 0.6);
}

.loading-more i {
    margin-right: 8px;
}

.no-more {
    text-align: center;
    padding: 20px 0;
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.9rem;
}

/* 响应式布局调整 */
@media (max-width: 1400px) {
    .history-grid {
        grid-template-columns: repeat(4, 1fr);
    }
}

@media (max-width: 1100px) {
    .history-grid {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (max-width: 768px) {
    .history-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 480px) {
    .history-grid {
        grid-template-columns: repeat(2, 1fr);
        /* 保持手机端也显示两列 */
        gap: 12px;
        /* 手机端稍微减小间距 */
    }
}

/* 删除状态标签相关样式 */
.status-badge,
.status-badge.completed,
.status-badge.planning {
    display: none;
}

/* 修改过渡动画样式 */
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.4s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

.modal-content-enter-active,
.modal-content-leave-active {
    transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-content-enter-from,
.modal-content-leave-to {
    transform: translateY(-50px);
    opacity: 0;
}

.modal-content-enter-to,
.modal-content-leave-from {
    transform: translateY(0);
    opacity: 1;
}
</style>
