<template>
    <div class="collection-modal-overlay" v-if="show" @click.self="$emit('close')">
        <div class="collection-modal-content">
            <div class="modal-header">
                <h2>我的收藏</h2>
                <button class="close-btn" @click="$emit('close')">
                    <i class="fas fa-times"></i>
                </button>
            </div>

            <div class="modal-body" @scroll="handleScroll" ref="modalBody">
                <div v-for="(group, date) in groupedCollection" :key="date" class="collection-group">
                    <div class="date-header">{{ formatDateHeader(date) }}</div>
                    <div class="collection-row">
                        <div v-for="anime in group" :key="anime.video_id" class="collection-item"
                            @click="goToAnime(anime.video_id)">
                            <div class="anime-cover">
                                <img :src="anime.cover_image_url" :alt="anime.title">
                                <div class="favorite-badge">
                                    <i class="fas fa-heart"></i>
                                </div>
                            </div>
                            <div class="anime-info">
                                <h3>{{ anime.title }}</h3>
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
                <div v-if="!hasMore && collectionList.length > 0" class="no-more">
                    没有更多收藏了
                </div>

                <!-- 空状态 -->
                <div v-if="collectionList.length === 0 && !isLoading" class="empty-state">
                    <i class="fas fa-heart-broken"></i>
                    <p>还没有收藏任何动漫哦</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'CollectionModal',
    props: {
        show: {
            type: Boolean,
            default: false
        }
    },
    data() {
        return {
            collectionList: [],
            currentPage: 1,
            isLoading: false,
            hasMore: true,
            total: 0
        }
    },
    computed: {
        groupedCollection() {
            const groups = {};
            this.collectionList.forEach(anime => {
                const date = new Date(anime.collected_at).toLocaleDateString();
                if (!groups[date]) {
                    groups[date] = [];
                }
                groups[date].push(anime);
            });
            return groups;
        }
    },
    methods: {
        async fetchCollection() {
            if (this.isLoading || !this.hasMore) return;

            try {
                this.isLoading = true;
                const response = await axios.get('/api/user/collection', {
                    params: {
                        page: this.currentPage
                    }
                });

                if (response.data.code === 200) {
                    const newCollection = response.data.anime_collection;
                    this.total = response.data.total;

                    if (!newCollection || newCollection.length === 0) {
                        this.hasMore = false;
                    } else {
                        this.collectionList = [...this.collectionList, ...newCollection];
                        this.currentPage++;
                    }
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                console.error('获取收藏列表失败:', error);
            } finally {
                this.isLoading = false;
            }
        },

        handleScroll(e) {
            const { scrollTop, scrollHeight, clientHeight } = e.target;
            // 当距离底部100px时加载更多
            if (scrollHeight - scrollTop - clientHeight < 100) {
                this.fetchCollection();
            }
        },

        goToAnime(videoId) {
            this.$router.push({
                path: '/moviesDetail',
                query: { videoId }
            });
            this.$emit('close');
        },

        formatDateHeader(date) {
            const collected = new Date(date);
            const now = new Date();
            const diff = now - collected;
            const days = Math.floor(diff / (1000 * 60 * 60 * 24));

            if (days === 0) return '今天';
            if (days === 1) return '昨天';
            if (days < 7) return `${days}天前`;
            return date;
        }
    },
    mounted() {
        this.fetchCollection();
    },
    watch: {
        show(newVal) {
            if (newVal) {
                this.currentPage = 1;
                this.collectionList = [];
                this.hasMore = true;
                this.fetchCollection();
            }
        }
    }
}
</script>

<style scoped>
.collection-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
    animation: fade-in 0.3s;
}

.collection-modal-content {
    width: 90%;
    max-width: 800px;
    /* 调整最大宽度 */
    height: 85vh;
    background: #1a1a1a;
    border-radius: 16px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    animation: slide-in 0.3s;
}

.modal-header {
    padding: 16px 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    background: linear-gradient(to bottom, rgba(255, 255, 255, 0.05), transparent);
}

.modal-header h2 {
    color: white;
    margin: 0;
    font-size: 1.2rem;
    font-weight: 500;
}

.close-btn {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.6);
    cursor: pointer;
    font-size: 1.1rem;
    padding: 8px;
    transition: all 0.3s ease;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.close-btn:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
}

.collection-group {
    margin-bottom: 32px;
}

.date-header {
    color: rgba(255, 255, 255, 0.6);
    font-size: 1.25rem;
    margin-bottom: 12px;
    padding-left: 16px;
    position: relative;
    font-weight: 500;
}

.date-header::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 16px;
    background: linear-gradient(to bottom, #FFD700, #FFA500);
    border-radius: 2px;
}

.collection-row {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 16px;
    margin-bottom: 8px;
    padding: 0 4px;
}

.collection-item {
    width: 100%;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.collection-item:hover {
    transform: translateY(-4px);
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.1);
}

.anime-cover {
    position: relative;
    width: 100%;
    padding-top: 130%;
    /* 调整图片比例 */
}

.anime-cover img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.favorite-badge {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 32px;
    height: 32px;
    background: rgba(204, 0, 0, 0.8);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.favorite-badge i {
    color: white;
    font-size: 0.9rem;
}

.anime-info {
    padding: 12px;
}

.anime-info h3 {
    margin: 0;
    color: white;
    font-size: 0.95rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
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

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: rgba(255, 255, 255, 0.4);
}

.empty-state i {
    font-size: 3rem;
    margin-bottom: 16px;
}

.empty-state p {
    font-size: 1rem;
    margin: 0;
}

/* 滚动条样式 */
.modal-body::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

.modal-body::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
    transition: all 0.3s ease;
}

.modal-body::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
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

/* 响应式调整 */
@media (max-width: 1200px) {
    .collection-row {
        grid-template-columns: repeat(4, 1fr);
    }
}

@media (max-width: 992px) {
    .collection-row {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (max-width: 768px) {
    .collection-row {
        grid-template-columns: repeat(3, 1fr);
        gap: 12px;
    }

    .date-header {
        font-size: 0.95rem;
        padding-left: 12px;
    }

    .date-header::before {
        width: 2px;
        height: 14px;
    }
}

@media (max-width: 480px) {
    .collection-row {
        grid-template-columns: repeat(2, 1fr);
        gap: 10px;
    }

    .collection-group {
        margin-bottom: 20px;
    }

    .date-header {
        font-size: 0.9rem;
        padding-left: 10px;
    }

    .date-header::before {
        height: 12px;
    }
}
</style>
