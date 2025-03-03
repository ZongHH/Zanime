<template>
    <div class="personal-page">
        <!-- 基本信息区域 -->
        <section class="profile-section">
            <div class="container-fluid">
                <div class="profile-container">
                    <div class="profile-header">
                        <div class="profile-avatar">
                            <img :src="userInfo.avatar_url" :alt="userInfo.username">
                            <div class="vip-badge" v-if="userInfo.is_vip">VIP</div>
                        </div>
                        <div class="profile-info">
                            <h2 class="username">{{ userInfo.username }}</h2>
                            <p class="user-id">ID: {{ userInfo.user_id }}</p>
                            <p class="join-date">注册时间：{{ formatDate(userInfo.register_time) }}</p>
                            <p class="signature">{{ userInfo.signature || '这个人很懒，什么都没写~' }}</p>
                            <div class="profile-actions">
                                <button class="edit-btn" @click="editProfile">
                                    <i class="fas fa-edit"></i> 编辑资料
                                </button>
                                <!-- <button class="settings-btn" @click="openSettings">
                                    <i class="fas fa-cog"></i> 设置
                                </button> -->
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- 数据统计区域 -->
        <section class="stats-section">
            <div class="container-fluid">
                <div class="stats-container">
                    <div class="stat-item">
                        <i class="fas fa-user-friends"></i>
                        <div class="stat-info">
                            <h3>{{ userStats.following_count }}</h3>
                            <p>我的关注</p>
                        </div>
                    </div>
                    <div class="stat-item">
                        <i class="fas fa-edit"></i>
                        <div class="stat-info">
                            <h3>{{ userStats.post_count }}</h3>
                            <p>我的帖子</p>
                        </div>
                    </div>
                    <div class="stat-item">
                        <i class="fas fa-star"></i>
                        <div class="stat-info">
                            <h3>{{ userStats.favorite_post_count }}</h3>
                            <p>收藏帖子</p>
                        </div>
                    </div>
                    <div class="stat-item">
                        <i class="fas fa-comment"></i>
                        <div class="stat-info">
                            <h3>{{ userStats.comment_count }}</h3>
                            <p>我的评论</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- 主要内容区域 -->
        <section class="content-section">
            <div class="container-fluid">
                <div class="content-container">
                    <!-- 左侧内容 -->
                    <div class="main-content">
                        <!-- 观看历史 -->
                        <div class="history-section section-card">
                            <div class="section-header">
                                <h3>观看历史</h3>
                                <button class="more-btn" @click="showHistoryModal"
                                    v-if="watchHistory && watchHistory.length > displayLimit">
                                    更多 <i class="fas fa-chevron-right"></i>
                                </button>
                            </div>
                            <div class="anime-grid history-grid">
                                <template v-if="watchHistory && watchHistory.length > 0">
                                    <div v-for="(anime, index) in limitedWatchHistory" :key="anime.id"
                                        class="anime-card" @click="continueWatching(anime)">
                                        <div class="card-image">
                                            <img :src="anime.cover_image_url" :alt="anime.title">
                                            <div class="progress-overlay">
                                                <div class="progress-bar" :style="{ width: anime.progress + '%' }">
                                                </div>
                                            </div>
                                            <div class="status-badge" :class="anime.status">
                                                {{ getStatusText(anime.status) }}
                                            </div>
                                        </div>
                                        <div class="card-info">
                                            <h4>{{ anime.title }}</h4>
                                            <div class="episode-info">
                                                <span class="video-name">{{ anime.video_name }}</span>
                                                <div class="episode-progress">
                                                    <span class="episode">{{ anime.episode }}</span>
                                                    <span class="progress-text">{{ formatProgress(anime.progress)
                                                        }}</span>
                                                </div>
                                            </div>
                                            <p class="update-time">{{ formatLastWatched(anime.updated_at) }}</p>
                                        </div>
                                    </div>
                                </template>
                                <div v-else class="empty-state">
                                    <p>还没有观看记录~</p>
                                </div>
                            </div>
                        </div>

                        <!-- 收藏动漫 (原我的追番) -->
                        <div class="favorite-anime-section section-card">
                            <div class="section-header">
                                <h3>收藏动漫</h3>
                                <button class="more-btn" @click="showCollectionModal"
                                    v-if="favoriteAnimes && favoriteAnimes.length > displayLimit">
                                    更多 ({{ favoriteTotal }}) <i class="fas fa-chevron-right"></i>
                                </button>
                            </div>
                            <div class="anime-grid">
                                <template v-if="favoriteAnimes && favoriteAnimes.length > 0">
                                    <div v-for="anime in limitedFavoriteAnimes" :key="anime.video_id" class="anime-card"
                                        @click="goToAnime(anime.video_id)">
                                        <div class="card-image">
                                            <img :src="anime.cover_image_url" :alt="anime.title">
                                            <div class="favorite-badge">
                                                <i class="fas fa-heart"></i>
                                            </div>
                                        </div>
                                        <div class="card-info">
                                            <h4>{{ anime.title }}</h4>
                                            <p class="favorite-date">收藏于 {{ formatDate(anime.collected_at) }}</p>
                                        </div>
                                    </div>
                                </template>
                                <div v-else class="empty-state">
                                    <p>还没有收藏动漫~</p>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 右侧边栏 -->
                    <div class="sidebar-content">
                        <!-- 会员状态 -->
                        <div class="vip-status section-card">
                            <div class="vip-header">
                                <div class="vip-crown">
                                    <i class="fas fa-crown"></i>
                                </div>
                                <h4>会员状态</h4>
                            </div>
                            <div class="vip-info">
                                <div class="vip-level-badge">
                                    <span>{{ userInfo.vip_level || '普通用户' }}</span>
                                </div>
                                <div class="vip-details">
                                    <p class="expire-date">{{ userInfo.vip_expire_date ?
                                        `到期时间：${formatDate(userInfo.vip_expire_date)}` : '开通会员享受更多权益' }}</p>
                                    <button class="upgrade-btn" @click="upgradeVip">
                                        <span>{{ userInfo.is_vip ? '续费会员' : '开通会员' }}</span>
                                        <i class="fas fa-arrow-right"></i>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- 我的帖子列表 -->
                        <div class="my-posts section-card">
                            <h4 class="mb-3">最近发布</h4>
                            <div class="posts-list">
                                <template v-if="recentPosts && recentPosts.length > 0">
                                    <div v-for="post in recentPosts" :key="post.id" class="post-item"
                                        @click="goToPost(post.id)">
                                        <div class="post-content">
                                            <h4 class="post-title">{{ post.title }}</h4>
                                            <div class="post-meta">
                                                <div class="meta-stats">
                                                    <span class="meta-item">
                                                        <i class="fas fa-eye"></i>
                                                        {{ formatNumber(post.view_count) }}
                                                    </span>
                                                    <span class="meta-item">
                                                        <i class="fas fa-comment"></i>
                                                        {{ formatNumber(post.comment_count) }}
                                                    </span>
                                                    <span class="meta-item">
                                                        <i class="fas fa-heart"></i>
                                                        {{ formatNumber(post.like_count) }}
                                                    </span>
                                                </div>
                                                <span class="post-time">{{ formatLastPosted(post.created_at) }}</span>
                                            </div>
                                        </div>
                                        <div class="post-arrow">
                                            <i class="fas fa-chevron-right"></i>
                                        </div>
                                    </div>
                                </template>
                                <div v-else class="empty-state">
                                    <p>还没有发布过帖子~</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- 添加历史记录弹窗 -->
        <History v-if="showHistory" :show="showHistory" @close="closeHistory" :userInfo="userInfo" />

        <!-- 添加编辑资料弹窗 -->
        <Edit v-if="showEdit" :show="showEdit" @close="closeEdit" @profile-updated="handleProfileUpdate" />
        <Collection v-if="showCollection" :show="showCollection" @close="closeCollection" />
    </div>
</template>

<script>
import axios from 'axios'
import { ElMessage } from 'element-plus'
import History from './History.vue'
import Edit from './Edit.vue'
import Collection from './Collection.vue'

export default {
    components: {
        History,
        Edit,
        Collection
    },
    data() {
        return {
            userInfo: {
                username: '',
                user_id: '',
                avatar_url: '',
                email: '',
                gender: '',
                register_time: new Date(),
                signature: '',
                is_vip: false,
                vip_level: '',
                vip_expire_date: null
            },
            userStats: {
                following_count: 0,      // 我的关注数
                post_count: 0,           // 我的帖子数
                favorite_post_count: 0,   // 收藏帖子数
                comment_count: 0         // 我的评论数
            },
            watchHistory: [],           // 观看历史
            favoriteAnimes: [],
            favoriteTotal: 0,
            recentPosts: [],            // 最近发布的帖子
            displayLimit: 6,            // 默认显示数量
            showHistory: false,
            showEdit: false,
            showCollection: false  // 添加收藏弹窗控制变量
        }
    },
    computed: {
        /**
         * 限制观看历史显示数量
         * @returns {Array} 返回截取后的观看历史数组
         */
        limitedWatchHistory() {
            if (this.watchHistory) {
                return this.watchHistory.slice(0, this.displayLimit);
            }
            return [];
        },
        /**
         * 限制收藏动漫显示数量
         * @returns {Array} 返回截取后的收藏动漫数组
         */
        limitedFavoriteAnimes() {
            if (this.favoriteAnimes) {
                return this.favoriteAnimes.slice(0, this.displayLimit);
            }
            return [];
        }
    },
    methods: {
        /**
         * 格式化日期为中文格式
         * @param {Date} date - 需要格式化的日期
         * @returns {string} 格式化后的日期字符串
         */
        formatDate(date) {
            return new Date(date).toLocaleDateString('zh-CN');
        },
        /**
         * 将观看时长(分钟)转换为天和小时的格式
         * @param {number} minutes - 观看时长(分钟)
         * @returns {string} 格式化后的时长字符串
         */
        formatWatchTime(minutes) {
            const hours = Math.floor(minutes / 60);
            const days = Math.floor(hours / 24);
            if (days > 0) {
                return `${days}天${hours % 24}小时`;
            }
            return `${hours}小时${minutes % 60}分钟`;
        },
        /**
         * 计算并格式化最后观看时间
         * @param {Date} date - 最后观看的日期
         * @returns {string} 格式化后的时间描述
         */
        formatLastWatched(date) {
            // 计算距离现在的时间
            const now = new Date();
            const watched = new Date(date);
            const diff = now - watched;
            const days = Math.floor(diff / (1000 * 60 * 60 * 24));

            if (days === 0) return '今天';
            if (days === 1) return '昨天';
            if (days < 7) return `${days}天前`;
            return this.formatDate(date);
        },
        /**
         * 获取用户统计信息
         * 包括关注数、帖子数、收藏数等
         */
        async fetchUserStats() {
            try {
                const response = await axios.get('/api/user/stats');
                if (response.data.code == 200) {
                    this.userStats = response.data.data;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                console.error('获取用户统计信息失败:', error);
            }
        },
        /**
         * 获取用户观看历史记录
         */
        async fetchWatchHistory() {
            try {
                const response = await axios.get('/api/watch-history', {
                    params: {
                        user_id: this.userInfo.user_id,
                    }
                });
                if (response.data.code == 200) {
                    this.watchHistory = response.data.progress;
                } else {
                    console.error(`获取观看历史失败: `, response.data.message);
                }
            } catch (error) {
                console.error('获取观看历史失败:', error);
            }
        },
        /**
         * 获取用户收藏的动漫列表
         */
        async fetchFavoriteAnimes() {
            try {
                const response = await axios.get('/api/user/collection', {
                    params: {
                        page: 1,
                    }
                });
                if (response.data.code === 200) {
                    this.favoriteAnimes = response.data.anime_collection;
                    this.favoriteTotal = response.data.total;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                console.error('获取收藏动漫失败:', error);
            }
        },
        /**
         * 获取用户最近发布的帖子
         */
        async fetchRecentPosts() {
            try {
                const response = await axios.get('/api/post/recent');
                if (response.data.code == 200) {
                    this.recentPosts = response.data.posts;
                } else {
                    console.error('获取最近帖子失败:', response.data.message);
                }
            } catch (error) {
                console.error('获取最近帖子失败:', error);
            }
        },
        /**
         * 显示编辑个人资料弹窗
         */
        editProfile() {
            this.showEdit = true;
        },
        /**
         * 继续观看指定动漫
         * @param {Object} anime - 动漫信息对象
         */
        continueWatching(anime) {
            this.$router.push({
                path: '/moviesDetail',
                query: {
                    videoId: anime.video_id,
                    episode: anime.episode,
                }
            });
        },
        /**
         * 升级VIP会员
         * TODO: 实现升级会员功能
         */
        upgradeVip() {
            this.handleServiceUnavailable();
        },
        /**
         * 从localStorage加载用户信息
         */
        async loadUserInfo() {
            try {
                const response = await axios.get('/api/user/profile');
                if (response.data.code == 200) {
                    this.userInfo = response.data.profile;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('获取用户信息失败: ' + error.message);
            }
        },
        /**
         * 获取观看状态的中文描述
         * @param {string} status - 观看状态
         * @returns {string} 状态的中文描述
         */
        getStatusText(status) {
            const statusMap = {
                'watching': '观看中',
                'completed': '已看完',
                'paused': '已暂停',
                'dropped': '已放弃',
                'planning': '计划看'
            };
            return statusMap[status] || status;
        },
        /**
         * 显示收藏列表弹窗
         */
        showCollectionModal() {
            this.showCollection = true;
        },
        /**
         * 关闭收藏列表弹窗
         */
        closeCollection() {
            this.showCollection = false;
        },
        /**
         * 格式化视频进度时间
         * @param {number} seconds - 秒数
         * @returns {string} 格式化后的时间字符串 (分:秒)
         */
        formatProgress(seconds) {
            if (!seconds) return '0:00';
            const minutes = Math.floor(seconds / 60);
            const remainingSeconds = Math.floor(seconds % 60);
            // 确保秒数始终显示两位数
            const formattedSeconds = remainingSeconds.toString().padStart(2, '0');
            return `${minutes}:${formattedSeconds}`;
        },
        /**
         * 计算并格式化最后发帖时间
         * @param {Date} date - 发帖日期
         * @returns {string} 格式化后的时间描述
         */
        formatLastPosted(date) {
            const now = new Date();
            const posted = new Date(date);
            const diff = now - posted;
            const minutes = Math.floor(diff / 60000);
            const hours = Math.floor(minutes / 60);
            const days = Math.floor(hours / 24);

            if (minutes < 60) return `${minutes}分钟前`;
            if (hours < 24) return `${hours}小时前`;
            if (days < 30) return `${days}天前`;
            return this.formatDate(date);
        },
        /**
         * 跳转到指定帖子详情页
         * @param {string|number} postId - 帖子ID
         */
        goToPost(postId) {
            this.$router.push(`/post/${postId}`);
        },
        /**
         * 显示观看历史弹窗
         */
        showHistoryModal() {
            this.showHistory = true;
        },
        /**
         * 关闭观看历史弹窗
         */
        closeHistory() {
            this.showHistory = false;
        },
        /**
         * 关闭编辑资料弹窗
         */
        closeEdit() {
            this.showEdit = false;
        },
        /**
         * 处理个人资料更新
         * @param {Object} updatedData - 更新后的个人资料数据
         */
        handleProfileUpdate(updatedData) {
            // 更新本地用户信息
            this.userInfo = { ...this.userInfo, ...updatedData };

            // 更新 localStorage
            Object.entries(updatedData).forEach(([key, value]) => {
                localStorage.setItem(key, value);
            });
        },
        /**
         * 格式化数字为带单位的字符串
         * @param {number} num - 需要格式化的数字
         * @returns {string} 格式化后的字符串
         */
        formatNumber(num) {
            if (num >= 10000) {
                return (num / 10000).toFixed(1) + 'w';
            } else if (num >= 1000) {
                return (num / 1000).toFixed(1) + 'k';
            }
            return num;
        },
        /**
         * 跳转到动漫详情页
         * @param {string|number} videoId - 视频ID
         */
        goToAnime(videoId) {
            this.$router.push({
                path: '/moviesDetail',
                query: {
                    videoId: videoId
                }
            });
        },
        /**
         * 处理服务暂未开放的情况
         */
        handleServiceUnavailable() {
            ElMessage.warning("服务暂未开放");
        }
    },
    async mounted() {
        this.loadUserInfo();
        await Promise.all([
            this.fetchUserStats(),
            this.fetchWatchHistory(),
            this.fetchFavoriteAnimes(),
            this.fetchRecentPosts()
        ]);
    }
};
</script>

<style scoped>
@import "@/static/css/personal.css";

/* 最近发布样式 */
.posts-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.post-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    background: linear-gradient(145deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.02));
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.post-item:hover {
    transform: translateX(4px);
    background: linear-gradient(145deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.03));
    border-color: rgba(255, 255, 255, 0.1);
}

.post-content {
    flex: 1;
    min-width: 0;
    /* 防止文本溢出 */
}

.post-title {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.95rem;
    margin: 0 0 8px 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.post-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
}

.meta-stats {
    display: flex;
    gap: 16px;
}

.meta-item {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.85rem;
    display: flex;
    align-items: center;
    gap: 4px;
}

.meta-item i {
    font-size: 0.9rem;
}

.post-time {
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.85rem;
}

.post-arrow {
    color: rgba(255, 255, 255, 0.3);
    font-size: 0.9rem;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
}

.post-item:hover .post-arrow {
    color: rgba(255, 255, 255, 0.6);
    transform: translateX(4px);
}

/* 特殊图标颜色 */
.meta-item .fa-eye {
    color: #64B5F6;
}

.meta-item .fa-comment {
    color: #81C784;
}

.meta-item .fa-heart {
    color: #E57373;
}

/* 响应式调整 */
@media (max-width: 768px) {
    .post-item {
        padding: 12px;
    }

    .meta-stats {
        gap: 12px;
    }

    .post-title {
        font-size: 0.9rem;
    }
}

@media (max-width: 480px) {
    .meta-stats {
        gap: 8px;
    }

    .meta-item {
        font-size: 0.8rem;
    }

    .post-time {
        font-size: 0.8rem;
    }
}

/* 优化收藏动漫卡片样式 */
.anime-card {
    position: relative;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
}

.anime-card:hover {
    transform: translateY(-4px);
    background: rgba(255, 255, 255, 0.08);
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

.favorite-date {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.85rem;
    margin-top: 4px;
}

.empty-state p {
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.9rem;
    text-align: center;
    margin: 0;
    white-space: nowrap;
}
</style>
