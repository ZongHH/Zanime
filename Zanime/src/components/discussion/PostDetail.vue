<template>
    <div class="post-detail-page">
        <div class="container-fluid">
            <!-- 返回按钮 -->
            <button class="back-btn" @click="$router.back()">
                <i class="fas fa-arrow-left"></i>
                返回讨论区
            </button>

            <!-- 帖子主体内容 -->
            <div class="post-content">
                <div class="post-header">
                    <div class="author-info">
                        <img :src="post.author.avatar_url" :alt="post.author.username">
                        <div class="author-meta">
                            <h4>{{ post.author.username }}</h4>
                            <div class="post-meta">
                                <span class="post-time">
                                    <i class="far fa-clock"></i>
                                    {{ formatTime(post.created_at) }}
                                </span>
                                <span class="post-views">
                                    <i class="far fa-eye"></i>
                                    {{ post.view_count }} 次浏览
                                </span>
                            </div>
                        </div>
                    </div>
                    <div class="post-actions">
                        <button class="action-btn follow" v-if="!isAuthor" @click="followAuthor">
                            <i class="fas fa-user-plus"></i>
                            关注作者
                        </button>
                        <button class="action-btn edit" v-if="isAuthor" @click="followAuthor">
                            <i class="fas fa-edit"></i>
                            编辑
                        </button>
                        <button class="action-btn delete" v-if="isAuthor" @click="followAuthor">
                            <i class="fas fa-trash-alt"></i>
                            删除
                        </button>
                    </div>
                </div>

                <div class="post-main">
                    <h1 class="post-title">{{ post.title }}</h1>
                    <div class="post-tags">
                        <span v-for="tag in post.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
                    </div>
                    <div class="post-content-text" v-html="formattedContent"></div>
                    <div v-if="post.images && post.images.length" class="post-images">
                        <div v-for="image in post.images" :key="image.id" class="image-item"
                            @click="showImagePreview(image.url)">
                            <img :src="image.url" :alt="post.title">
                        </div>
                    </div>
                </div>

                <div class="post-footer">
                    <div class="interaction-btns">
                        <button class="like-btn" :class="{ active: post.is_liked }" @click="toggleLike">
                            <i class="fas fa-heart"></i>
                            <span>{{ post.like_count }} 点赞</span>
                        </button>
                        <button class="favorite-btn" :class="{ active: post.is_favorited }" @click="toggleFavorite">
                            <i class="fas fa-star"></i>
                            <span>{{ post.favorite_count }} 收藏</span>
                        </button>
                        <button class="share-btn" @click="sharePost">
                            <i class="fas fa-share-alt"></i>
                            <span>分享</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- 评论区域 -->
            <Comment ref="commentComponent" :post="post" :comments="comments" :current-user="currentUser"
                :has-more-comments="hasMoreComments" :current-page="currentPage" :total-pages="totalPages"
                @submit-comment="handleCommentSubmit" @toggle-like="handleCommentLike"
                @delete-comment="handleCommentDelete" @submit-reply="handleCommentReply" @change-page="loadComments"
                @load-replies="handleLoadReplies" />
        </div>

        <!-- 图片预览弹窗 -->
        <div v-if="showPreview" class="image-preview-modal" @click="closeImagePreview">
            <img :src="previewImage" alt="预览图片">
            <button class="close-preview-btn" @click="closeImagePreview">
                <i class="fas fa-times"></i>
            </button>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import Comment from './Comment.vue';
import { ElMessage } from 'element-plus';

export default {
    name: 'PostDetail',
    components: {
        Comment
    },
    data() {
        return {
            // 帖子详情数据
            post: {
                id: null,                // 帖子ID
                title: '',               // 帖子标题
                content: '',             // 帖子内容
                author: {},              // 作者信息
                created_at: '',          // 创建时间
                view_count: 0,           // 浏览数
                like_count: 0,           // 点赞数
                comment_count: 0,        // 评论数
                favorite_count: 0,       // 收藏数
                is_liked: false,         // 当前用户是否点赞
                is_favorited: false,     // 当前用户是否收藏
                tags: [],                // 帖子标签
                images: []               // 帖子图片
            },
            // 当前登录用户信息
            currentUser: {
                id: null,                // 用户ID
                username: '',            // 用户名
                avatar_url: ''           // 头像URL
            },
            comments: [
                {
                    id: 0,                    // 评论唯一标识
                    content: '',              // 评论文本内容
                    author: {                 // 评论作者信息
                        id: 0,                // 作者ID
                        username: '',         // 作者用户名
                        avatar_url: ''        // 作者头像
                    },
                    created_at: '',           // 创建时间
                    like_count: 0,            // 点赞数
                    is_liked: false,          // 是否已点赞
                    replies: [                // 回复列表
                        {
                            id: 0,            // 回复ID
                            content: '',      // 回复内容
                            author: {         // 回复作者
                                id: 0,
                                username: '',
                                avatar_url: ''
                            },
                            created_at: '',   // 回复时间
                            like_count: 0,    // 点赞数
                            is_liked: false,  // 是否已点赞
                            reply_to: {       // 回复目标用户
                                id: 0,
                                username: '',
                                avatar_url: ''
                            }
                        }
                    ],
                    reply_num: 0              // 回复总数
                }
            ],
            commentText: '',             // 评论输入框内容
            showPreview: false,          // 是否显示图片预览
            previewImage: '',            // 预览图片URL
            currentPage: 1,              // 当前评论页码
            totalPages: 1,               // 总页数（从后台获取）
            isLoadingComments: false,     // 是否正在加载评论
            hasMoreComments: false,      // 是否还有更多评论可以加载
        }
    },
    computed: {
        isAuthor() {
            return this.currentUser.id === this.post.author.id;
        },
        formattedContent() {
            // 可以在这里处理内容的格式化，比如将换行符转换为<br>标签等
            return this.post.content.replace(/\n/g, '<br>');
        }
    },
    methods: {
        /**
         * 获取帖子详情数据
         */
        async fetchPostDetail() {
            try {
                // 从URL中获取帖子ID
                const postId = this.$route.params.id;
                // 获取帖子详情数据
                const response = await axios.get(`/api/post/detail`, {
                    params: {
                        post_id: postId
                    }
                });
                if (response.data.code === 200) {
                    this.post = response.data.post;
                } else {
                    console.error('获取帖子详情失败:', response.data.message);
                }
            } catch (error) {
                console.error('获取帖子详情失败:', error);
            }
        },

        /**
         * 获取评论列表数据
         * @param {number} page - 页码，默认为1
         */
        async fetchComments(page = 1) {
            if (this.isLoadingComments) return;

            try {
                this.isLoadingComments = true;

                const response = await axios.get(`/api/post/comments`, {
                    params: {
                        page: page,
                        id: this.post.id
                    }
                });

                const data = response.data;

                // 更新评论列表和总页数
                this.comments = data.comments;
                this.totalPages = data.total_page;
                this.currentPage = page;
                this.hasMoreComments = page < this.totalPages;

            } catch (error) {
                console.error('获取评论失败:', error);
            } finally {
                this.isLoadingComments = false;
            }
        },

        /**
         * 格式化时间显示
         * @param {string} time - ISO格式的时间字符串
         * @returns {string} 格式化后的时间字符串
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
         * 显示图片预览
         * @param {string} url - 图片URL
         */
        showImagePreview(url) {
            this.previewImage = url;
            this.showPreview = true;
        },

        /**
         * 关闭图片预览
         */
        closeImagePreview() {
            this.showPreview = false;
            this.previewImage = '';
        },

        /**
         * 切换点赞状态
         */
        async toggleLike() {
            try {
                this.post.is_liked = !this.post.is_liked;
                this.post.like_count += this.post.is_liked ? 1 : -1;
                const response = await axios.post('/api/post/like', {
                    post_id: this.post.id,
                    status: this.post.is_liked ? 1 : 0
                });
                if (response.data.code === 200) {
                    // 点赞成功
                } else {
                    // 点赞失败
                    this.post.is_liked = !this.post.is_liked;
                    this.post.like_count += this.post.is_liked ? 1 : -1;
                }
            } catch (error) {
                console.error('点赞失败:', error);
            }
        },

        /**
         * 切换收藏状态
         */
        async toggleFavorite() {
            try {
                this.post.is_favorited = !this.post.is_favorited;
                this.post.favorite_count += this.post.is_favorited ? 1 : -1;
                const response = await axios.post('/api/post/favorite', {
                    post_id: this.post.id,
                    status: this.post.is_favorited ? 1 : 0
                });
                if (response.data.code === 200) {
                    // 收藏成功
                } else {
                    // 收藏失败
                    this.post.is_favorited = !this.post.is_favorited;
                    this.post.favorite_count += this.post.is_favorited ? 1 : -1;
                }
            } catch (error) {
                console.error('收藏失败:', error);
            }
        },

        /**
         * 关注作者
         */
        async followAuthor() {
            ElMessage.warning("服务暂未开放")
        },

        /**
         * 处理评论提交
         * @param {Object} comment - 评论数据对象
         */
        async handleCommentSubmit(comment) {
            const newComment = {
                id: comment.id,
                content: comment.content,
                author: this.currentUser,
                created_at: new Date().toISOString(),
                like_count: 0,
                is_liked: false,
                reply_num: 0,
                replies: []
            };

            const response = await axios.post('/api/comment/submit', {
                post_id: this.post.id,
                comment_id: newComment.id,
                content: comment.content,
                created_at: newComment.created_at,
            });

            if (response.data.code === 200) {
                // 如果在第一页，直接添加到列表开头
                if (this.currentPage === 1) {
                    this.comments.unshift(newComment);
                    this.post.comment_count++;
                } else {
                    // 如果不在第一页，跳转到第一页
                    this.loadComments(1);
                }
            }
        },

        /**
         * 处理评论点赞
         * @param {Object} comment - 评论对象
         */
        async handleCommentLike(comment) {
            comment.is_liked = !comment.is_liked;
            comment.like_count += comment.is_liked ? 1 : -1;

            try {
                const response = await axios.post('/api/comment/like', {
                    comment_id: comment.id,
                    status: comment.is_liked
                });
                if (response.data.code === 200) {
                    // 点赞成功
                } else {
                    // 点赞失败
                    comment.is_liked = !comment.is_liked;
                    comment.like_count += comment.is_liked ? 1 : -1;
                }
            } catch (error) {
                // 如果请求失败,回滚状态
                comment.is_liked = !comment.is_liked;
                comment.like_count += comment.is_liked ? 1 : -1;
            }
        },

        /**
         * 处理评论删除
         * @param {Object} comment - 评论对象
         */
        handleCommentDelete(comment) {
            this.comments = this.comments.filter(c => c.id !== comment.id);
            this.post.comment_count--;
        },

        /**
         * 处理评论回复
         * @param {Object} replyData - 回复数据对象
         */
        async handleCommentReply(replyData) {
            try {
                // 创建新回复对象
                const newReply = {
                    post_id: this.post.id,
                    id: replyData.id,
                    content: replyData.content,
                    author: this.currentUser,
                    created_at: replyData.created_at,
                    like_count: 0,
                    is_liked: false,
                    reply_to: replyData.reply_to
                };

                const response = await axios.post('/api/comment/reply', {
                    post_id: this.post.id,
                    content: replyData.content,
                    parent_id: replyData.parent_id,
                    root_id: replyData.root_id,
                    to_user_id: replyData.to_user_id,
                    created_at: replyData.created_at,
                    comment_id: replyData.id
                });

                if (response.data.code === 200) {
                    // 更新根评论的回复数和回复列表
                    const targetComment = this.comments.find(c => c.id === replyData.root_id);
                    if (targetComment) {
                        targetComment.reply_num++;

                        // 确保replies数组存在
                        if (!targetComment.replies) {
                            targetComment.replies = [];
                        }

                        // 添加到回复列表
                        targetComment.replies.push(newReply);
                    }
                }
            } catch (error) {
                console.error('回复评论失败:', error);
            }
        },

        /**
         * 加载指定页码的评论
         * @param {number} page - 目标页码
         */
        loadComments(page) {
            if (page === this.currentPage) return;
            this.fetchComments(page);
        },

        /**
         * 处理加载评论回复
         * @param {Object} params - 参数对象
         * @param {Object} params.comment - 评论对象
         * @param {number} params.page - 页码，默认为1
         */
        async handleLoadReplies({ comment, page = 1 }) {
            try {
                const response = await axios.get(`/api/comment/replies`, {
                    params: {
                        root_id: comment.id,
                        page: page
                    }
                });

                // 直接传递回复数据数组
                this.$refs.commentComponent.handleLoadRepliesResult(response.data);

            } catch (error) {
                console.error('加载回复失败:', error);
                // 出错时也要重置加载状态
                this.$refs.commentComponent.isLoadingMore = false;
            }
        },

        /**
         * 获取当前用户信息
         */
        async getCurrentUser() {
            try {
                const response = await axios.get('/api/user/current');
                if (response.data.code === 200) {
                    this.currentUser = response.data.user;
                }
            } catch (error) {
                console.error('获取用户信息失败:', error);
            }
        },
        /**
         * 分享帖子
         * 生成当前帖子的URL并复制到剪贴板，优先使用现代Clipboard API，
         * 如不可用则回退到传统的document.execCommand方法
         */
        sharePost() {
            // 构建完整的帖子URL，包含域名和帖子ID
            const postUrl = `${window.location.origin}/post/${this.post.id}`;

            // 检查现代剪贴板API是否可用
            if (navigator.clipboard && navigator.clipboard.writeText) {
                // 使用现代Clipboard API复制文本
                navigator.clipboard.writeText(postUrl).then(() => {
                    // 复制成功时显示成功提示
                    ElMessage.success('复制链接成功，快分享给好友吧！');
                }).catch(err => {
                    // 复制失败时记录错误并尝试使用备用方法
                    console.error('复制链接失败:', err);
                    this.fallbackCopyToClipboard(postUrl);
                });
            } else {
                // 浏览器不支持现代Clipboard API时使用备用复制方法
                this.fallbackCopyToClipboard(postUrl);
            }
        },

        /**
         * 备用复制到剪贴板方法
         * 使用传统的document.execCommand('copy')方法实现复制功能
         * 适用于不支持现代Clipboard API的浏览器
         * 
         * @param {string} text - 需要复制到剪贴板的文本
         */
        fallbackCopyToClipboard(text) {
            try {
                // 创建一个临时textarea元素用于选择和复制文本
                const textArea = document.createElement('textarea');
                textArea.value = text;

                // 设置样式使textarea不可见且不影响页面布局
                textArea.style.top = '0';
                textArea.style.left = '0';
                textArea.style.position = 'fixed';
                textArea.style.opacity = '0';
                textArea.style.pointerEvents = 'none';

                // 将textarea添加到DOM中
                document.body.appendChild(textArea);

                // 选中textarea中的文本
                textArea.focus();
                textArea.select();

                // 执行复制命令
                const successful = document.execCommand('copy');

                // 操作完成后从DOM中移除textarea
                document.body.removeChild(textArea);

                // 根据复制结果显示相应提示
                if (successful) {
                    ElMessage.success('复制链接成功，快分享给好友吧！');
                } else {
                    ElMessage.warning('无法自动复制，请手动复制链接分享');
                }
            } catch (err) {
                // 捕获并记录可能出现的错误
                console.error('备用复制方法失败:', err);
                ElMessage.warning('无法自动复制，请手动复制链接分享');
            }
        }
    },
    async mounted() {
        await this.fetchPostDetail();
        this.fetchComments();
        this.getCurrentUser();
    }
}
</script>

<style scoped>
@import "@/static/css/postDetail.css";
</style>