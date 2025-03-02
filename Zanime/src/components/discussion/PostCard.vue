<template>
    <div class="post-card-content" @click="navigateToPost">
        <div class="post-header">
            <div class="post-author">
                <img :src="post.author.avatar_url" :alt="post.author.username">
                <div class="author-info">
                    <h4>{{ post.author.username }}</h4>
                    <span class="post-time">{{ formatTime(post.created_at) }}</span>
                </div>
            </div>
            <div class="post-badges">
                <span v-if="post.is_featured" class="badge featured">
                    <i class="fas fa-star"></i> 精华
                </span>
                <span v-if="post.is_pinned" class="badge pinned">
                    <i class="fas fa-thumbtack"></i> 置顶
                </span>
            </div>
        </div>

        <div class="post-body">
            <h3 class="post-title">{{ post.title }}</h3>
            <div class="post-tags">
                <span v-for="tag in post.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
            </div>
            <p class="post-preview">{{ post.content }}</p>
            <div v-if="post.images && post.images.length" :class="[
                'post-images',
                { 'single-image': post.images.length === 1 },
                { 'double-images': post.images.length === 2 }
            ]">
                <img v-for="image in post.images.slice(0, 3)" :key="image.id" :src="image.url" :alt="post.title"
                    loading="lazy">
            </div>
        </div>

        <div class="post-footer">
            <div class="post-stats">
                <span class="stat-item">
                    <i class="fas fa-eye"></i> {{ post.view_count }}
                </span>
                <span class="stat-item">
                    <i class="fas fa-comment"></i> {{ post.comment_count }}
                </span>
                <span class="stat-item">
                    <i class="fas fa-heart"></i> {{ post.like_count }}
                </span>
            </div>
            <div class="post-actions">
                <button class="action-btn like" :class="{ active: post.is_liked }" @click.stop="toggleLike">
                    <i class="fas fa-heart"></i>
                </button>
                <button class="action-btn favorite" :class="{ active: post.is_favorited }" @click.stop="toggleFavorite">
                    <i class="fas fa-star"></i>
                </button>
                <button class="action-btn share" @click.stop="sharePost">
                    <i class="fas fa-share-alt"></i>
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    props: {
        post: {
            type: Object,
            required: true
        }
    },
    methods: {
        /**
         * 格式化时间为相对时间
         * @param {string|Date} time - 需要格式化的时间
         * @returns {string} - 格式化后的相对时间字符串（如：5分钟前，2小时前）
         */
        formatTime(time) {
            const now = new Date();
            const postTime = new Date(time);
            const diff = now - postTime;
            const minutes = Math.floor(diff / 60000);
            const hours = Math.floor(minutes / 60);
            const days = Math.floor(hours / 24);

            if (minutes < 60) return `${minutes}分钟前`;
            if (hours < 24) return `${hours}小时前`;
            if (days < 30) return `${days}天前`;
            return postTime.toLocaleDateString();
        },

        /**
         * 导航到帖子详情页
         * 当用户点击帖子时触发，使用Vue Router导航到对应的帖子详情页
         */
        navigateToPost() {
            this.$router.push(`/post/${this.post.id}`);
        },

        /**
         * 切换帖子点赞状态
         * 如果用户未登录，会重定向到登录页
         * 发送API请求更新点赞状态，并在本地更新点赞计数和状态
         */
        async toggleLike() {
            // 检查用户是否已登录
            if (!localStorage.getItem('user_id')) {
                this.$router.push('/login');
                return;
            }
            try {
                // 发送点赞/取消点赞请求
                const response = await axios.post(`/api/post/like`, {
                    post_id: this.post.id,
                    status: this.post.is_liked ? 0 : 1  // 0表示取消点赞，1表示点赞
                });
                if (response.data.code === 200) {
                    // 更新本地点赞状态
                    this.post.is_liked = !this.post.is_liked;
                    // 更新点赞数量
                    this.post.like_count += this.post.is_liked ? 1 : -1;
                } else {
                    console.error('点赞失败:', response.data.message);
                }
            } catch (error) {
                console.error('点赞失败:', error);
            }
        },

        /**
         * 切换帖子收藏状态
         * 如果用户未登录，会重定向到登录页
         * 发送API请求更新收藏状态，并在本地更新收藏状态
         */
        async toggleFavorite() {
            // 检查用户是否已登录
            if (!localStorage.getItem('user_id')) {
                this.$router.push('/login');
                return;
            }
            try {
                // 发送收藏/取消收藏请求
                const response = await axios.post(`/api/post/favorite`, {
                    post_id: this.post.id,
                    status: this.post.is_favorited ? 0 : 1  // 0表示取消收藏，1表示收藏
                });
                if (response.data.code === 200) {
                    // 更新本地收藏状态
                    this.post.is_favorited = !this.post.is_favorited;
                } else {
                    console.error('收藏失败:', response.data.message);
                }
            } catch (error) {
                console.error('收藏失败:', error);
            }
        },

        /**
         * 分享帖子
         * 复制帖子链接到剪贴板，并显示成功提示
         * 优先使用现代Clipboard API，如不可用则使用备用方法
         */
        sharePost() {
            // 构建完整的帖子URL
            const postUrl = `${window.location.origin}/post/${this.post.id}`;

            // 检查剪贴板API是否可用
            if (navigator.clipboard && navigator.clipboard.writeText) {
                navigator.clipboard.writeText(postUrl).then(() => {
                    ElMessage.success('复制链接成功，快分享给好友吧！');
                }).catch(err => {
                    console.error('复制链接失败:', err);
                    this.fallbackCopyToClipboard(postUrl);
                });
            } else {
                // 使用备用复制方法
                this.fallbackCopyToClipboard(postUrl);
            }
        },

        /**
         * 备用复制到剪贴板方法
         * 当现代Clipboard API不可用时使用
         * 创建临时文本区域并使用document.execCommand进行复制
         * 
         * @param {string} text - 要复制到剪贴板的文本
         */
        fallbackCopyToClipboard(text) {
            try {
                // 创建临时文本区域元素
                const textArea = document.createElement('textarea');
                textArea.value = text;
                // 设置样式以防止滚动和视觉干扰
                textArea.style.top = '0';
                textArea.style.left = '0';
                textArea.style.position = 'fixed';
                textArea.style.opacity = '0';
                textArea.style.pointerEvents = 'none';
                document.body.appendChild(textArea);
                textArea.focus();
                textArea.select();

                // 执行复制命令
                const successful = document.execCommand('copy');
                // 移除临时元素
                document.body.removeChild(textArea);

                if (successful) {
                    ElMessage.success('复制链接成功，快分享给好友吧！');
                } else {
                    ElMessage.warning('无法自动复制，请手动复制链接分享');
                }
            } catch (err) {
                console.error('备用复制方法失败:', err);
                ElMessage.warning('无法自动复制，请手动复制链接分享');
            }
        }
    }
};
</script>

<style scoped>
.post-card-content {
    padding: 24px;
    cursor: pointer;
    transition: all 0.3s ease;
    background: rgba(30, 30, 30, 0.8);
    backdrop-filter: blur(20px);
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.05);
    overflow: hidden;
}

.post-card-content:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4);
    background: rgba(30, 30, 30, 0.9);
}

/* 帖子头部样式 */
.post-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.post-author {
    display: flex;
    align-items: center;
    gap: 14px;
}

.post-author img {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    border: 2px solid rgba(139, 0, 0, 0.5);
    object-fit: cover;
}

.author-info h4 {
    color: #fff;
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 4px;
}

.post-time {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.85rem;
}

.post-badges {
    display: flex;
    gap: 8px;
}

.badge {
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 0.85rem;
    display: flex;
    align-items: center;
    gap: 4px;
}

.badge.featured {
    background: linear-gradient(45deg, #ffd700, #ffa500);
    color: #000;
    font-weight: 600;
}

.badge.pinned {
    background: linear-gradient(45deg, #8b0000, #cc0000);
    color: white;
    font-weight: 600;
}

/* 帖子内容样式 */
.post-body {
    margin-bottom: 24px;
}

.post-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #fff;
    margin-bottom: 16px;
    line-height: 1.4;
}

.post-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
}

.tag {
    background: rgba(139, 0, 0, 0.15);
    color: rgba(255, 255, 255, 0.9);
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.85rem;
    border: 1px solid rgba(139, 0, 0, 0.3);
    transition: all 0.3s ease;
}

.tag:hover {
    background: rgba(139, 0, 0, 0.25);
    transform: translateY(-1px);
}

.post-preview {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.95rem;
    line-height: 1.6;
    margin-bottom: 20px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

/* 图片样式 */
.post-images {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-top: 16px;
    max-width: 100%;
}

.post-images.single-image {
    grid-template-columns: minmax(0, 600px);
}

.post-images.double-images {
    grid-template-columns: repeat(2, 1fr);
    max-width: 800px;
}

.post-images img {
    width: 100%;
    aspect-ratio: 16/9;
    object-fit: cover;
    border-radius: 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.2);
}

.post-images img:hover {
    transform: scale(1.02);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
    border-color: rgba(139, 0, 0, 0.3);
}

/* 帖子底部样式 */
.post-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.post-stats {
    display: flex;
    gap: 24px;
}

.stat-item {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    gap: 6px;
    transition: color 0.3s ease;
}

.stat-item:hover {
    color: rgba(255, 255, 255, 0.8);
}

.post-actions {
    display: flex;
    gap: 12px;
}

.action-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
}

.action-btn:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
}

.action-btn.like.active {
    background: linear-gradient(45deg, #8b0000, #cc0000);
    color: white;
}

.action-btn.favorite.active {
    background: linear-gradient(45deg, #ffd700, #ffa500);
    color: #1a1a1a;
}

/* 响应式调整 */
@media (max-width: 1200px) {
    .post-images {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 768px) {
    .post-card-content {
        padding: 20px;
    }

    .post-header {
        flex-direction: column;
        gap: 12px;
    }

    .post-badges {
        margin-top: 8px;
    }

    .post-title {
        font-size: 1.1rem;
    }

    .post-images {
        grid-template-columns: 1fr;
        gap: 8px;
    }

    .post-images.double-images {
        grid-template-columns: 1fr;
    }

    .post-images img {
        max-height: 300px;
    }

    .post-footer {
        flex-direction: column;
        gap: 16px;
    }

    .post-stats {
        width: 100%;
        justify-content: space-around;
    }

    .post-actions {
        width: 100%;
        justify-content: space-around;
    }
}
</style>