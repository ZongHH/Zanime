<template>
    <div class="discussion-page">
        <!-- 固定在顶部的分区导航 -->
        <section class="category-section">
            <div class="container-fluid">
                <div class="category-container">
                    <div class="category-tabs">
                        <button v-for="category in categories" :key="category.id"
                            :class="['category-tab', { active: currentCategory === category.id }]"
                            @click="changeCategory(category.id)">
                            <i :class="category.icon"></i>
                            {{ category.name }}
                            <span class="post-count">({{ category.post_count }})</span>
                        </button>
                    </div>

                    <!-- 搜索框 -->
                    <div class="search-box disabled">
                        <input type="text" disabled placeholder="搜索功能暂未开放">
                        <i class="fas fa-search"></i>
                    </div>

                    <!-- 发帖按钮 -->
                    <button class="create-post-btn" @click="showCreatePostModal">
                        <i class="fas fa-edit"></i> 发布帖子
                    </button>
                </div>
            </div>
        </section>

        <!-- 主要内容区域 - 移除内部滚动，改为整页滚动 -->
        <section class="content-section">
            <div class="container-fluid">
                <div class="content-container">
                    <!-- 左侧热门话题 -->
                    <div class="left-sidebar">
                        <div class="hot-topics sticky">
                            <h3>热门话题</h3>
                            <div class="topic-list">
                                <div v-for="topic in hotTopics" :key="topic.id" class="topic-item"
                                    @click.prevent="handleServiceUnavailable">
                                    <span class="topic-rank">#{{ topic.rank }}</span>
                                    <div class="topic-info">
                                        <h4>{{ topic.title }}</h4>
                                        <p>{{ topic.viewCount }}次浏览</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 中间帖子列表 - 移除滚动监听 -->
                    <div class="main-content">
                        <!-- 置顶帖子 -->
                        <div class="pinned-posts" v-if="pinnedPosts.length > 0">
                            <div v-for="post in pinnedPosts" :key="post.id" class="post-card pinned">
                                <PostCard :post="post" />
                            </div>
                        </div>

                        <!-- 普通帖子列表 -->
                        <div class="posts-grid">
                            <div v-for="post in posts" :key="post.id" class="post-card">
                                <PostCard :post="post" />
                            </div>
                        </div>

                        <!-- 加载状态 -->
                        <div v-if="isLoading" class="loading-more">
                            <i class="fas fa-spinner fa-spin"></i>
                            加载中...
                        </div>

                        <!-- 没有更多数据提示 -->
                        <div v-if="!hasMore && posts.length > 0" class="no-more">
                            没有更多内容了
                        </div>
                    </div>

                    <!-- 右侧活动公告 -->
                    <div class="right-sidebar">
                        <div class="announcements sticky">
                            <h3>活动公告</h3>
                            <div class="announcement-list">
                                <div v-for="announcement in announcements" :key="announcement.id"
                                    class="announcement-item" @click.prevent="handleServiceUnavailable">
                                    <div class="announcement-badge" :class="announcement.type">
                                        {{ announcement.type === 'activity' ? '活动' : '公告' }}
                                    </div>
                                    <div class="announcement-content">
                                        <h4>{{ announcement.title }}</h4>
                                        <p>{{ announcement.date }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- 发帖模态框 -->
        <CreatePostModal v-if="showCreatePost" :currentCategoryId="currentCategory" @close="showCreatePost = false" />
    </div>
</template>

<script>
import { debounce } from 'lodash';
import { ElMessage } from 'element-plus';
import axios from 'axios';
import PostCard from '@/components/discussion/PostCard.vue';
import CreatePostModal from '@/components/discussion/CreatePostModal.vue';

export default {
    components: {
        PostCard,
        CreatePostModal
    },
    data() {
        return {
            currentCategory: 1,
            categories: [],
            searchQuery: '',
            currentPage: 1,
            posts: [],
            pinnedPosts: [],
            showCreatePost: false,
            hotTopics: [],
            announcements: [],
            isLoading: false,
            hasMore: true,
        }
    },
    methods: {
        async changeCategory(categoryId) {
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
            this.currentCategory = categoryId;
            this.currentPage = 1;
            this.posts = [];
            this.hasMore = true;
            this.isLoading = false;
            await this.fetchPosts();
        },
        showCreatePostModal() {
            this.showCreatePost = true;
        },
        async fetchPosts() {
            if (this.isLoading || !this.hasMore) return;

            try {
                this.isLoading = true;
                const response = await axios.get('/api/post/list', {
                    params: {
                        category_id: this.currentCategory,
                        page: this.currentPage,
                        page_size: 10
                    }
                });

                if (response.data.code === 200) {
                    const newPosts = response.data.posts || [];

                    // 如果是第一页
                    if (this.currentPage === 1) {
                        this.posts = newPosts;
                        this.pinnedPosts = response.data.pinned_posts || [];
                    } else {
                        this.posts = [...this.posts, ...newPosts];
                    }

                    // 如果没有获取到新数据，说明已经没有更多内容了
                    if (newPosts.length === 0) {
                        this.hasMore = false;
                    } else {
                        this.currentPage += 1;
                    }
                } else {
                    console.error('获取帖子列表失败:', response.data.message);
                }
            } catch (error) {
                console.error('获取帖子列表失败:', error);
                this.hasMore = false; // 发生错误时也设置为没有更多数据
            } finally {
                this.isLoading = false;
            }
        },
        async fetchHotTopics() {
            try {
                const response = await axios.get('/api/hot-topics');
                this.hotTopics = response.data;
            } catch (error) {
                console.error('获取热门话题失败:', error);
            }
        },
        async fetchAnnouncements() {
            try {
                const response = await axios.get('/api/announcements');
                this.announcements = response.data;
            } catch (error) {
                console.error('获取公告失败:', error);
            }
        },
        async fetchCategories() {
            try {
                const response = await axios.get('/api/post/categories');
                if (response.data.code === 200) {
                    this.categories = response.data.categories;
                } else {
                    console.error('获取分类失败:', response.data.message);
                }
            } catch (error) {
                console.error('获取分类失败:', error);
            }
        },
        handleScroll() {
            const mainContent = document.querySelector('.main-content');
            if (!mainContent) return;

            const rect = mainContent.getBoundingClientRect();
            const bottomDistance = mainContent.scrollHeight - (window.innerHeight - rect.top);

            // 当main-content距离底部100px时加载更多
            if (bottomDistance <= 100) {
                this.fetchPosts();
            }
        },
        handleServiceUnavailable() {
            ElMessage.warning("服务暂未开放");
        },
    },
    async mounted() {
        await this.fetchCategories();
        if (this.categories.length > 0) {
            this.changeCategory(this.categories[0].id);
        }

        this.hotTopics = [
            {
                id: 1,
                rank: 1,
                title: "《鬼灭之刃》锻刀村篇开播",
                viewCount: "12.5万"
            },
            {
                id: 2,
                rank: 2,
                title: "《间谍过家家》第二季完结",
                viewCount: "8.9万"
            },
            {
                id: 3,
                rank: 3,
                title: "《咒术回战》第二季最终话",
                viewCount: "7.2万"
            },
            {
                id: 4,
                rank: 4,
                title: "《海贼王》漫画最新话情报",
                viewCount: "6.8万"
            },
            {
                id: 5,
                rank: 5,
                title: "《进击的巨人》完结纪念展",
                viewCount: "5.4万"
            }
        ];

        this.announcements = [
            {
                id: 1,
                title: "网站春节活动预告",
                content: "新年将至，我们准备了丰富的活动和奖励...",
                type: "activity",
                date: "2024-01-20"
            },
            {
                id: 2,
                title: "系统维护通知",
                content: "本周日凌晨将进行系统升级维护...",
                type: "system",
                date: "2024-01-21"
            },
            {
                id: 3,
                title: "新功能上线公告",
                content: "评论区增加了表情包功能...",
                type: "feature",
                date: "2024-01-19"
            }
        ];

        // 添加窗口滚动监听
        window.addEventListener('scroll', this.handleScroll);
    },
    beforeUnmount() {
        // 组件销毁前移除滚动监听
        window.removeEventListener('scroll', this.handleScroll);
    },
    watch: {
        searchQuery: debounce(function () {
            this.currentPage = 1;
            this.fetchPosts();
        }, 300),
    }
};
</script>

<style scoped>
@import "@/static/css/discussion.css";

/* 修改滚动相关样式 */
.discussion-page {
    min-height: 100vh;
}

/* 固定顶部导航 */
.category-section {
    position: sticky;
    top: 0;
    z-index: 1;
    background: #1a1a1a;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 移除main-content的滚动样式 */
.main-content {
    padding: 20px;
}

/* 侧边栏固定效果 */
.sticky {
    position: sticky;
    top: 80px;
    /* 导航栏高度 + 一些间距 */
}

/* 内容布局调整 */
.content-container {
    display: grid;
    grid-template-columns: 250px 1fr 250px;
    gap: 20px;
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
}

/* 加载状态样式 */
.loading-more {
    text-align: center;
    padding: 20px 0;
    color: #666;
}

.loading-more i {
    margin-right: 8px;
}

.no-more {
    text-align: center;
    padding: 20px 0;
    color: #666;
    font-size: 14px;
}

/* 响应式布局 */
@media (max-width: 1200px) {
    .content-container {
        grid-template-columns: 200px 1fr 200px;
    }
}

@media (max-width: 992px) {
    .content-container {
        grid-template-columns: 1fr;
    }

    .left-sidebar,
    .right-sidebar {
        display: none;
    }
}

/* 搜索框禁用状态样式 */
.search-box.disabled {
    position: relative;
    opacity: 0.7;
    cursor: not-allowed;
}

.search-box.disabled input {
    background: rgba(255, 255, 255, 0.05);
    cursor: not-allowed;
    color: rgba(255, 255, 255, 0.5);
}

.search-box.disabled i {
    color: rgba(255, 255, 255, 0.3);
}

/* 移动端优化样式 */
@media (max-width: 768px) {
    .category-container {
        flex-direction: column;
        gap: 12px;
        padding: 12px;
    }

    .category-tabs {
        width: 100%;
        overflow-x: auto;
        padding-bottom: 8px;
        scrollbar-width: none;
        /* Firefox */
        -ms-overflow-style: none;
        /* IE & Edge */
        white-space: nowrap;
        -webkit-overflow-scrolling: touch;
    }

    .category-tabs::-webkit-scrollbar {
        display: none;
        /* Chrome, Safari, Opera */
    }

    .category-tab {
        padding: 8px 14px;
        font-size: 13px;
        margin-right: 8px;
        border-radius: 20px;
        display: inline-flex;
        align-items: center;
        gap: 6px;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.08);
        white-space: nowrap;
        transition: all 0.3s ease;
    }

    .category-tab i {
        font-size: 14px;
    }

    .category-tab.active {
        background: linear-gradient(135deg, rgba(171, 5, 17, 0.9), rgba(171, 5, 17, 0.7));
        border-color: transparent;
        box-shadow: 0 3px 12px rgba(171, 5, 17, 0.3);
        transform: translateY(-1px);
    }

    .post-count {
        font-size: 12px;
        opacity: 0.8;
    }

    .search-box.disabled {
        width: 100%;
        margin: 0;
    }

    .search-box.disabled input {
        font-size: 13px;
        height: 42px;
        border-radius: 8px;
    }

    .create-post-btn {
        width: 100%;
        margin: 0;
        height: 42px;
        border-radius: 8px;
        font-size: 14px;
        background: linear-gradient(135deg, rgba(171, 5, 17, 0.9), rgba(171, 5, 17, 0.7));
        border: none;
        box-shadow: 0 3px 12px rgba(171, 5, 17, 0.2);
        transition: all 0.3s ease;
    }

    .create-post-btn:hover {
        transform: translateY(-2px);
        box-shadow: 0 5px 15px rgba(171, 5, 17, 0.3);
    }

    .create-post-btn:active {
        transform: translateY(0);
    }

    .main-content {
        padding: 15px 12px;
    }

    .posts-grid {
        display: flex;
        flex-direction: column;
        gap: 15px;
    }

    .post-card {
        margin-bottom: 0;
        animation: fadeIn 0.5s ease-out;
    }

    .pinned-posts {
        margin-bottom: 15px;
    }

    .pinned-posts .post-card {
        margin-bottom: 12px;
    }

    .loading-more {
        padding: 15px 0;
        font-size: 14px;
        opacity: 0.8;
    }

    .loading-more i {
        animation: spin 1.2s linear infinite;
    }

    .no-more {
        padding: 12px 0;
        font-size: 12px;
        opacity: 0.6;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }

        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }

        100% {
            transform: rotate(360deg);
        }
    }
}

/* 超小屏幕优化 */
@media (max-width: 480px) {
    .content-container {
        padding: 12px 8px;
    }

    .category-tab {
        padding: 7px 12px;
        font-size: 12px;
    }

    .main-content {
        padding: 12px 8px;
    }

    .search-box.disabled input {
        height: 38px;
        font-size: 12px;
    }

    .create-post-btn {
        height: 38px;
        font-size: 13px;
    }
}
</style>
