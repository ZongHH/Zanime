<template>
    <transition name="dropdown-fade">
        <div v-if="showNotifications" class="notification-dropdown">
            <div class="notify-header">
                <h5>互动消息</h5>
                <div class="filter-dropdown">
                    <span class="selected-filter" @click="toggleFilterDropdown">
                        {{ getFilterText(activeFilter) }} <i class="fas fa-chevron-down"
                            :class="{ 'rotate': showFilterDropdown }"></i>
                    </span>
                    <transition name="slide-down">
                        <div class="filter-options" v-if="showFilterDropdown">
                            <div class="filter-option" @click="setFilter(0)">全部消息</div>
                            <div class="filter-option" @click="setFilter(1)">点赞评论</div>
                            <div class="filter-option" @click="setFilter(4)">点赞帖子</div>
                            <div class="filter-option" @click="setFilter(2)">回复</div>
                            <div class="filter-option" @click="setFilter(3)">收藏</div>
                            <div class="filter-option" @click="setFilter(5)">关注</div>
                        </div>
                    </transition>
                </div>
            </div>

            <div class="notify-body" ref="notifyBody" @scroll="handleScroll">
                <div v-for="notification in filteredNotifications" :key="notification.id" class="notify-item">
                    <div class="user-avatar">
                        <img :src="notification.avatar" alt="用户头像">
                        <span v-if="getNotificationIcon(notification.type)" class="notify-icon"
                            :class="'type-' + notification.type">
                            <i :class="getNotificationIcon(notification.type)"></i>
                        </span>
                    </div>
                    <div class="notify-content">
                        <div class="user-name">{{ notification.username }}</div>
                        <p class="message-text">{{ notification.message }}</p>
                        <div class="notify-meta">
                            <span class="notify-time">
                                {{ getNotificationText(notification.type) }} {{ formatDate(notification.time) }}
                            </span>
                        </div>
                    </div>
                    <div class="action-buttons" v-if="notification.type === 5">
                        <button class="follow-button" v-if="!notification.isFollowing"
                            @click.stop="followUser(notification.id)">
                            关注
                        </button>
                        <button class="unfollow-button" v-else @click.stop="unfollowUser(notification.id)">
                            已关注
                        </button>
                    </div>
                </div>

                <!-- 加载更多指示器 -->
                <div v-if="loading" class="loading-indicator">
                    <div class="spinner"></div>
                    <span>加载中...</span>
                </div>

                <!-- 无更多数据提示 -->
                <div v-if="!loading && !hasMore && notifications.length > 0" class="no-more-data">
                    没有更多消息了
                </div>

                <!-- 空状态提示 -->
                <div v-if="!loading && notifications.length === 0" class="notify-empty">
                    <i class="fas fa-comments"></i>
                    <p>暂无{{ getFilterText(activeFilter) }}</p>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    props: {
        showNotifications: {
            type: Boolean,
            required: true
        }
    },
    data() {
        return {
            notifications: [],
            activeFilter: 0, // 0: 全部, 1-6: 对应通知类型
            showFilterDropdown: false,
            loading: false,
            page: 1,
            hasMore: true,
            scrollDebounceTimer: null
        }
    },
    computed: {
        filteredNotifications() {
            return this.notifications; // 直接返回从后端获取的数据，不需要前端筛选
        }
    },
    watch: {
        showNotifications(newVal) {
            if (newVal && this.notifications.length === 0) {
                // 当通知面板打开且没有数据时，获取通知
                this.fetchNotifications();
            }
        }
    },
    methods: {
        /**
         * 关注用户
         * @param {number} userId - 要关注的用户ID
         * @description 向后端发送关注请求，成功后更新本地通知数据中的关注状态
         */
        followUser(userId) {
            // 发送关注请求到后端
            axios.post('/api/follow', { userId })
                .then(response => {
                    if (response.data.code === 200) {
                        // 请求成功，更新本地通知数据中的关注状态
                        const notification = this.notifications.find(n => n.id === userId);
                        if (notification) {
                            notification.isFollowing = true;
                        }
                    }
                })
                .catch(error => {
                    console.error('关注失败:', error);
                });
        },

        /**
         * 取消关注用户
         * @param {number} userId - 要取消关注的用户ID
         * @description 向后端发送取消关注请求，成功后更新本地通知数据中的关注状态
         */
        unfollowUser(userId) {
            // 发送取消关注请求到后端
            axios.post('/api/unfollow', { userId })
                .then(response => {
                    if (response.data.code === 200) {
                        // 请求成功，更新本地通知数据中的关注状态
                        const notification = this.notifications.find(n => n.id === userId);
                        if (notification) {
                            notification.isFollowing = false;
                        }
                    }
                })
                .catch(error => {
                    console.error('取消关注失败:', error);
                });
        },

        /**
         * 获取通知文本描述
         * @param {number} type - 通知类型(1-6)
         * @returns {string} 对应类型的通知文本描述
         * @description 根据通知类型返回对应的文本描述
         */
        getNotificationText(type) {
            switch (type) {
                case 1: return '点赞了你的评论';
                case 2: return '回复了你的评论';
                case 3: return '收藏了你的帖子';
                case 4: return '点赞了你的帖子';
                case 5: return '关注了你';
                default: return '';
            }
        },

        /**
         * 获取通知图标
         * @param {number} type - 通知类型(1-6)
         * @returns {string} 对应类型的Font Awesome图标类名
         * @description 根据通知类型返回对应的图标类名
         */
        getNotificationIcon(type) {
            switch (type) {
                case 1: return 'fas fa-thumbs-up';  // 点赞评论图标
                case 2: return 'fas fa-comment';    // 回复图标
                case 3: return 'fas fa-bookmark';   // 收藏图标
                case 4: return 'fas fa-heart';      // 点赞帖子图标
                case 5: return 'fas fa-user-plus';  // 关注和回关图标
                default: return '';
            }
        },

        /**
         * 获取筛选器文本
         * @param {number} filterType - 筛选类型(0-6)
         * @returns {string} 对应筛选类型的文本描述
         * @description 根据筛选类型返回对应的文本描述，用于显示在筛选器上
         */
        getFilterText(filterType) {
            switch (filterType) {
                case 0: return '全部消息';
                case 1: return '点赞评论';
                case 2: return '回复';
                case 3: return '收藏';
                case 4: return '点赞帖子';
                case 5: return '关注';
                default: return '全部消息';
            }
        },

        /**
         * 切换筛选下拉菜单的显示状态
         * @description 显示/隐藏筛选下拉菜单
         */
        toggleFilterDropdown() {
            this.showFilterDropdown = !this.showFilterDropdown;
        },

        /**
         * 设置通知筛选类型
         * @param {number} filterType - 筛选类型(0-6)
         * @description 设置当前活动的筛选类型，并重新获取对应类型的通知数据
         */
        setFilter(filterType) {
            if (this.activeFilter === filterType) {
                this.showFilterDropdown = false;
                return; // 如果选择的筛选类型相同，则不重复请求
            }

            this.activeFilter = filterType;
            this.showFilterDropdown = false;
            this.resetAndFetch(); // 重置并获取新数据
        },

        /**
         * 获取通知数据
         * @async
         * @description 从后端API获取通知数据，支持分页和类型筛选
         */
        async fetchNotifications() {
            if (this.loading || !this.hasMore) return;

            this.loading = true;

            try {
                // 构建请求参数
                const params = {
                    page: this.page,
                    type: this.activeFilter,
                };

                // 发送请求获取通知
                const response = await axios.get('/api/user/notifications', { params });

                if (response.data.code === 200) {
                    const newNotifications = response.data.notifications;

                    // 如果是第一页，则替换通知数据，否则追加
                    if (this.page === 1) {
                        this.notifications = newNotifications;
                    } else {
                        this.notifications = [...this.notifications, ...newNotifications];
                    }

                    // 判断是否还有更多数据
                    this.hasMore = newNotifications.length !== 0;

                    // 增加页码
                    this.page += 1;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error(`获取通知失败，请稍后重试: ${error.message}`);
            } finally {
                this.loading = false;
            }
        },

        /**
         * 处理滚动事件
         * @param {Event} event - 滚动事件对象
         * @description 监听滚动事件，实现滚动到底部时自动加载更多通知数据，使用防抖处理
         */
        handleScroll(event) {
            // 防抖处理
            clearTimeout(this.scrollDebounceTimer);
            this.scrollDebounceTimer = setTimeout(() => {
                const element = event.target;
                // 当滚动到底部附近时加载更多
                if (element.scrollHeight - element.scrollTop - element.clientHeight < 50) {
                    this.fetchNotifications();
                }
            }, 200);
        },

        /**
         * 重置并获取通知数据
         * @description 重置页码和通知列表，然后重新获取通知数据，用于筛选条件变更时
         */
        resetAndFetch() {
            this.page = 1;
            this.hasMore = true;
            this.notifications = [];
            this.fetchNotifications();
        },

        /**
         * 格式化日期时间
         * @param {string} dateString - ISO格式的时间字符串
         * @returns {string} 格式化后的日期时间字符串
         * @description 将ISO格式的时间转换为更易读的格式
         */
        formatDate(dateString) {
            if (!dateString) return '';

            const date = new Date(dateString);
            if (isNaN(date.getTime())) return dateString; // 如果解析失败，返回原始字符串

            const now = new Date();
            const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            const yesterday = new Date(today);
            yesterday.setDate(yesterday.getDate() - 1);

            // 格式化时间部分
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            const timeStr = `${hours}:${minutes}`;

            // 判断是今天、昨天还是更早
            if (date >= today) {
                return `今天 ${timeStr}`;
            } else if (date >= yesterday) {
                return `昨天 ${timeStr}`;
            } else if (date.getFullYear() === now.getFullYear()) {
                // 今年的其他日期，显示月-日
                const month = (date.getMonth() + 1).toString().padStart(2, '0');
                const day = date.getDate().toString().padStart(2, '0');
                return `${month}-${day} ${timeStr}`;
            } else {
                // 往年的日期，显示完整年-月-日
                const year = date.getFullYear();
                const month = (date.getMonth() + 1).toString().padStart(2, '0');
                const day = date.getDate().toString().padStart(2, '0');
                return `${year}-${month}-${day} ${timeStr}`;
            }
        },
    },
};
</script>

<style scoped>
.notification-dropdown {
    position: absolute;
    top: 60px;
    right: -110px;
    width: 385px;
    background: #0a0a0a;
    backdrop-filter: blur(10px);
    border-radius: 16px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.03);
    z-index: 1000;
    overflow: hidden;
    transform-origin: top right;
    min-height: 300px;
}

.notify-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 18px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    position: relative;
    background: #080808;
}

.notify-header:after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 1px;
    background: linear-gradient(90deg,
            rgba(100, 100, 100, 0),
            rgba(100, 100, 100, 0.1),
            rgba(100, 100, 100, 0));
}

.notify-header h5 {
    color: #ffffff;
    font-size: 1.2rem;
    font-weight: 600;
    margin: 0;
    letter-spacing: 0.5px;
}

/* 筛选器下拉菜单样式 */
.filter-dropdown {
    position: relative;
    cursor: pointer;
    z-index: 1100;
}

.selected-filter {
    color: #d1c3c3;
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    transition: all 0.3s ease;
}

.selected-filter i {
    font-size: 0.8rem;
    margin-left: 5px;
    transition: transform 0.3s;
}

.selected-filter:hover {
    color: #ffffff;
}

.selected-filter .rotate {
    transform: rotate(180deg);
}

.filter-options {
    position: absolute;
    top: 100%;
    right: 0;
    width: 120px;
    background: #151515;
    border-radius: 8px;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.05);
    margin-top: 10px;
    z-index: 1100;
    overflow: visible;
}

/* 下拉列表滑动动画 */
.slide-down-enter-active,
.slide-down-leave-active {
    transition: all 0.3s ease;
    max-height: 300px;
}

.slide-down-enter-from,
.slide-down-leave-to {
    opacity: 0;
    max-height: 0;
    transform: translateY(-10px);
}

.filter-option {
    padding: 10px 15px;
    color: #d1c3c3;
    font-size: 0.9rem;
    transition: all 0.2s ease;
}

.filter-option:hover {
    background: #222;
    color: #ffffff;
}

.notify-body {
    max-height: 600px;
    min-height: 300px;
    overflow-y: auto;
    padding: 5px;
    background: #0a0a0a;
}

.notify-body::-webkit-scrollbar {
    width: 4px;
}

.notify-body::-webkit-scrollbar-track {
    background: #0a0a0a;
    border-radius: 2px;
}

.notify-body::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.6);
    border-radius: 2px;
}

.notify-body::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.8);
}

.notify-item {
    display: flex;
    padding: 15px;
    margin: 8px;
    border-radius: 12px;
    background: #121212;
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    position: relative;
    overflow: hidden;
}

.notify-item:hover {
    background: #1c1c1c;
}

.user-avatar {
    position: relative;
    width: 48px;
    height: 48px;
    margin-right: 15px;
    flex-shrink: 0;
}

.user-avatar img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #1c1c1c;
}

.notify-icon {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid #0a0a0a;
}

.notify-icon i {
    color: white;
    font-size: 8px;
}

/* 不同类型通知的图标颜色 */
.type-1 {
    background: #4e6ef2;
    /* 点赞评论 - 蓝色 */
}

.type-2 {
    background: #0072ff;
    /* 回复评论 - 深蓝色 */
}

.type-3 {
    background: #f2b94e;
    /* 收藏帖子 - 金色 */
}

.type-4 {
    background: #f24e4e;
    /* 点赞帖子 - 红色 */
}

.type-5,
.type-6 {
    background: #4ef28c;
    /* 关注/回关 - 绿色 */
}

.notify-content {
    flex: 1;
    overflow: hidden;
}

.user-name {
    color: #ffffff;
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 5px;
}

.message-text {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.95rem;
    margin-bottom: 8px;
    line-height: 1.5;
    word-break: break-word;
}

.notify-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.notify-time {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.8rem;
}

.action-buttons {
    margin-left: 10px;
    display: flex;
    align-items: center;
}

.follow-button,
.unfollow-button {
    background: #cc0000;
    color: white;
    border: none;
    border-radius: 20px;
    padding: 6px 16px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
}

.unfollow-button {
    background: transparent;
    color: #cc0000;
    border: 1px solid #cc0000;
}

.follow-button:hover {
    background: #a30000;
    transform: translateY(-2px);
}

.unfollow-button:hover {
    background: rgba(204, 0, 0, 0.1);
    transform: translateY(-2px);
}

/* 加载状态和空状态样式 */
.loading-indicator,
.notify-empty,
.no-more-data,
.error-message {
    text-align: center;
    padding: 20px 0;
    color: rgba(255, 255, 255, 0.5);
}

.loading-indicator {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.spinner {
    width: 30px;
    height: 30px;
    border: 3px solid rgba(204, 0, 0, 0.3);
    border-radius: 50%;
    border-top-color: #cc0000;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.no-more-data {
    font-size: 0.9rem;
    padding: 15px 0;
    opacity: 0.7;
}

.notify-empty i {
    font-size: 3rem;
    margin-bottom: 15px;
    opacity: 0.5;
    color: #ffffff;
}

.notify-empty p {
    font-size: 1rem;
    color: #ffffff;
}

.error-message {
    display: flex;
    flex-direction: column;
    align-items: center;
    color: #f24e4e;
    padding: 20px;
}

.error-message i {
    font-size: 2rem;
    margin-bottom: 10px;
}

.retry-button {
    margin-top: 15px;
    background: rgba(204, 0, 0, 0.8);
    color: white;
    border: none;
    border-radius: 20px;
    padding: 6px 16px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.3s ease;
}

.retry-button:hover {
    background: #cc0000;
    transform: translateY(-2px);
}

/* 过渡动画 */
.dropdown-fade-enter-active,
.dropdown-fade-leave-active {
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
}

/* 响应式调整 */
@media (max-width: 480px) {
    .notification-dropdown {
        width: 320px;
        right: -100px;
    }
}
</style>