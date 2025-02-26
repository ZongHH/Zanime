<template>
    <div>
        <div class="add-comment mb-24">
            <img :src="currentUser.avatar_url" alt="avatar">
            <div class="comment-input">
                <input type="text" class="form-control" placeholder="写下你的评论" v-model="commentContent"
                    @keyup.enter="postComment">
                <button class="post-btn light-gray" @click="postComment" :disabled="!commentContent.trim()">发布</button>
            </div>
        </div>

        <h4 class="white mb-24">热门评论</h4>

        <div id="comment-container">
            <div v-if="!comments.length && !isLoading" class="no-comments">
                暂无评论，来发表第一条评论吧！
            </div>

            <div v-for="comment in comments" :key="comment.comment_id" class="review-card">
                <div>
                    <img :src="comment.userInfo.avatar_url" :alt="comment.userInfo.username">
                </div>
                <div class="textBlock">
                    <div class="d-flex align-items-end justify-content-between mb-16">
                        <div class="d-flex align-items-center">
                            <div>
                                <h6 class="fw-500 light-gray">
                                    {{ comment.userInfo.username }}
                                    <span v-if="comment.userInfo.id === currentUser.id" class="author-tag">我</span>
                                </h6>
                                <span class="subtitle light-gray">{{ formatRelativeTime(comment.created_at) }}</span>
                            </div>
                        </div>
                        <button class="reply-btn fw-500" @click="showReplyInput(comment)"
                            v-if="currentUser.id">回复</button>
                    </div>

                    <p class="mb-12">{{ comment.content }}</p>

                    <!-- 修改主评论的回复框部分 -->
                    <transition name="fade">
                        <div v-if="activeReplyComment === comment.comment_id" class="reply-comment-input">
                            <div class="comment-input">
                                <img :src="currentUser.avatar_url" :alt="currentUser.username">
                                <div class="input-wrapper">
                                    <input type="text" :placeholder="`回复 @${comment.userInfo.username}:`"
                                        class="minimal-input" v-model="replyContent" @keyup.enter="sendReply(comment)"
                                        :ref="el => { if (el) replyInputRef = el }">
                                    <div class="reply-actions">
                                        <button class="cancel-btn" @click="cancelReply">取消</button>
                                        <button class="submit-btn" :disabled="!replyContent.trim()"
                                            @click="sendReply(comment)">发送</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </transition>

                    <!-- 查看更多回复按钮 -->
                    <div class="reply-wrapper mb-24" v-if="comment.reply_num > (comment.replies?.length || 0)">
                        <div class="replies-actions">
                            <button class="get-replies-btn" @click="showReplyDetail(comment)">
                                <span>查看全部 {{ comment.reply_num }} 条回复</span>
                                <i class="fas fa-chevron-right"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="pagination-container" v-if="comments.length">
            <div class="pagination">
                <button class="page-btn" :disabled="currentPage === 1" @click="prevPage">
                    <i class="fas fa-chevron-left"></i>
                </button>
                <span class="current-page">{{ currentPage }}</span>
                <button class="page-btn" :disabled="currentPage === totalPages" @click="nextPage">
                    <i class="fas fa-chevron-right"></i>
                </button>
            </div>
        </div>

        <!-- 修改弹窗部分，添加过渡动画 -->
        <transition name="modal">
            <div v-if="showReplyModal" class="reply-modal" @click.self="closeReplyModal">
                <transition name="slide-down">
                    <div v-if="showReplyModal" class="reply-modal-content">
                        <div class="reply-modal-header">
                            <h3>全部回复 ({{ currentComment?.reply_num || 0 }})</h3>
                            <button class="close-btn" @click="closeReplyModal">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>

                        <div class="reply-modal-body" @scroll="handleModalScroll" ref="modalBody">
                            <!-- 原评论 -->
                            <div class="original-comment">
                                <div class="comment-author">
                                    <img :src="currentComment?.userInfo.avatar_url" alt="avatar">
                                    <div class="comment-meta">
                                        <h4>{{ currentComment?.userInfo.username }}</h4>
                                        <span class="comment-time">{{ formatRelativeTime(currentComment?.created_at)
                                        }}</span>
                                    </div>
                                </div>
                                <div class="comment-content">{{ currentComment?.content }}</div>
                            </div>

                            <!-- 回复列表 -->
                            <div class="modal-replies-list">
                                <div v-for="reply in modalReplies" :key="reply.comment_id" class="reply-item">
                                    <div class="reply-author">
                                        <img :src="reply.userInfo.avatar_url" alt="avatar">
                                        <div class="reply-meta">
                                            <h4>
                                                {{ reply.userInfo.username }}
                                                <span v-if="reply.userInfo.id === currentUser.id"
                                                    class="author-tag">我</span>
                                            </h4>
                                            <span class="reply-time">{{ formatRelativeTime(reply.created_at) }}</span>
                                        </div>
                                    </div>
                                    <div class="reply-content">
                                        <span v-if="reply.replied_name" class="reply-to">
                                            回复 @{{ reply.replied_name }}:
                                        </span>
                                        {{ reply.content }}
                                    </div>
                                    <div class="reply-actions">
                                        <button class="reply-btn" @click="showReplyInput(reply)">
                                            回复
                                        </button>
                                    </div>

                                    <!-- 修改弹窗中的回复框部分 -->
                                    <div v-if="activeReplyComment === reply.comment_id" class="modal-reply-input">
                                        <div class="comment-input">
                                            <img :src="currentUser.avatar_url" :alt="currentUser.username">
                                            <div class="input-wrapper">
                                                <input type="text" :placeholder="`回复 @${reply.userInfo.username}:`"
                                                    class="minimal-input" v-model="replyContent"
                                                    @keyup.enter="sendReply(reply)"
                                                    :ref="el => { if (el) replyInputRef = el }">
                                                <div class="reply-actions">
                                                    <button class="cancel-btn" @click="cancelReply">取消</button>
                                                    <button class="submit-btn" :disabled="!replyContent.trim()"
                                                        @click="sendReply(reply)">发送</button>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <!-- 加载状态和底部提示 -->
                                <div v-if="isLoadingMore" class="loading-more">
                                    <div class="loading-spinner"></div>
                                    <span>加载更多回复...</span>
                                </div>
                                <div v-else-if="!hasMoreModalReplies && modalReplies.length > 0"
                                    class="no-more-replies">
                                    已经到底啦 ~
                                </div>
                                <div v-else-if="!hasMoreModalReplies && modalReplies.length === 0"
                                    class="no-more-replies">
                                    暂无回复，来发表第一条回复吧！
                                </div>
                            </div>
                        </div>
                    </div>
                </transition>
            </div>
        </transition>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'Comment',
    props: {
        videoId: {
            type: Number,
            required: true
        }
    },
    data() {
        return {
            comments: [], // 存储当前视频的所有评论数据
            commentContent: '', // 用户输入的评论内容
            replyContent: '', // 用户输入的回复内容
            activeReplyComment: null, // 当前正在回复的评论对象，用于定位和显示回复框
            isLoading: false, // 是否正在加载评论数据
            currentPage: 1, // 当前评论列表的页码
            totalPages: 1, // 评论列表的总页数
            currentUser: { // 当前登录用户的信息
                id: 1, // 用户唯一标识
                username: '测试用户', // 用户名
                avatar_url: 'https://api.dicebear.com/7.x/avataaars/svg?seed=test', // 用户头像URL
            },
            showReplyModal: false, // 是否显示回复弹窗
            currentComment: null, // 当前查看回复的评论对象
            modalReplies: [], // 弹窗中显示的回复列表
            modalCurrentPage: 1, // 弹窗中回复列表的当前页码
            isLoadingMore: false, // 是否正在加载更多回复
            hasMoreModalReplies: true, // 弹窗中是否还有更多回复可加载
            replyInputRef: null, // 回复输入框的DOM引用，用于自动聚焦
        }
    },
    methods: {
        /**
         * 加载评论列表
         * 从服务器获取指定视频的评论列表数据
         * @returns {Promise<void>}
         */
        async loadComments() {
            // 如果正在加载中则返回
            if (this.isLoading) return;
            this.isLoading = true;

            try {
                // 发送GET请求获取评论列表
                const response = await axios.get('/api/movie/comments', {
                    params: {
                        video_id: this.videoId,
                        page: this.currentPage,
                    }
                });

                // 处理响应数据
                if (response.data.code === 200) {
                    this.totalPages = response.data.total_page;
                    this.comments = response.data.comments;
                } else {
                    throw new Error(response.data.message);
                }

            } catch (error) {
                // 错误处理
                ElMessage.error('加载评论失败: ' + error.message);
            } finally {
                // 无论成功失败都要关闭加载状态
                this.isLoading = false;
            }
        },

        /**
         * 发表评论
         * 向服务器提交新的评论
         * @returns {Promise<void>}
         */
        async postComment() {
            // 检查评论内容是否为空
            if (!this.commentContent.trim()) return;

            try {
                // 发送POST请求提交评论
                const response = await axios.post('/api/movie/submit', {
                    video_id: this.videoId,
                    content: this.commentContent,
                });

                // 处理响应
                if (response.data.code === 200) {
                    ElMessage.success('评论发表成功');
                    this.commentContent = ''; // 清空评论内容
                    await this.loadComments(); // 重新加载评论列表
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('发表评论失败: ' + error.message);
            }
        },

        /**
         * 显示回复输入框
         * @param {Object} comment - 要回复的评论对象
         */
        showReplyInput(comment) {
            this.activeReplyComment = comment.comment_id;
            this.replyContent = '';
            // 等待DOM更新后聚焦输入框
            this.$nextTick(() => {
                if (this.replyInputRef) {
                    this.replyInputRef.focus();
                }
            });
        },

        /**
         * 取消回复
         * 清空回复状态和内容
         */
        cancelReply() {
            this.activeReplyComment = null;
            this.replyContent = '';
        },

        /**
         * 发送回复
         * 向服务器提交回复内容
         * @param {Object} comment - 被回复的评论对象
         * @returns {Promise<void>}
         */
        async sendReply(comment) {
            // 检查回复内容是否为空
            if (!this.replyContent.trim()) return;

            try {
                // 发送POST请求提交回复
                const response = await axios.post('/api/movie/submitReply', {
                    video_id: this.videoId,
                    parent_id: comment.comment_id,
                    root_id: comment.root_id || comment.comment_id,
                    content: this.replyContent,
                    to_user_id: comment.userInfo.id,
                });

                // 处理响应
                if (response.data.code === 200) {
                    ElMessage.success('回复发送成功');
                    this.replyContent = ''; // 清空回复内容
                    this.activeReplyComment = null; // 关闭回复框
                    await this.loadComments(); // 重新加载评论列表
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('发送回复失败: ' + error.message);
            }
        },

        /**
         * 显示回复详情弹窗
         * @param {Object} comment - 要显示回复的评论对象
         */
        showReplyDetail(comment) {
            this.currentComment = comment;
            this.showReplyModal = true;
            this.modalCurrentPage = 1;
            this.hasMoreModalReplies = true;
            this.loadModalReplies();
        },

        /**
         * 关闭回复详情弹窗
         * 重置相关状态
         */
        closeReplyModal() {
            this.showReplyModal = false;
            this.currentComment = null;
            this.modalReplies = [];
            this.modalCurrentPage = 1;
            this.hasMoreModalReplies = true;
        },

        /**
         * 格式化相对时间
         * 将时间戳转换为相对时间描述
         * @param {number|string} timestamp - 时间戳
         * @returns {string} 格式化后的相对时间字符串
         */
        formatRelativeTime(timestamp) {
            const now = new Date();
            const date = new Date(timestamp);
            const diff = now - date;

            const minutes = Math.floor(diff / 60000);
            const hours = Math.floor(minutes / 60);
            const days = Math.floor(hours / 24);

            if (minutes < 60) return `${minutes}分钟前`;
            if (hours < 24) return `${hours}小时前`;
            if (days < 30) return `${days}天前`;
            return date.toLocaleDateString();
        },

        /**
         * 加载弹窗回复列表
         * 从服务器获取评论的回复列表
         * @returns {Promise<void>}
         */
        async loadModalReplies() {
            // 如果正在加载或没有更多回复则返回
            if (this.isLoadingMore || !this.hasMoreModalReplies) return;

            this.isLoadingMore = true;
            try {
                // 发送GET请求获取回复列表
                const response = await axios.get('/api/movie/replies', {
                    params: {
                        root_id: this.currentComment.comment_id,
                        page: this.modalCurrentPage,
                    }
                });

                // 处理响应
                if (response.data.code === 200) {
                    this.modalReplies.push(...response.data.replies);
                    this.hasMoreModalReplies = response.data.replies.length !== 0;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('加载回复失败: ' + error.message);
            } finally {
                this.isLoadingMore = false;
            }
        },

        /**
         * 加载当前用户信息
         * @returns {Promise<void>}
         */
        async loadCurrentUser() {
            try {
                const response = await axios.get('/api/user/current');
                if (response.data.code === 200) {
                    this.currentUser = response.data.user;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('加载当前用户失败: ' + error.message);
            }
        },

        /**
         * 上一页
         * 加载上一页评论并滚动到顶部
         * @returns {Promise<void>}
         */
        async prevPage() {
            this.currentPage--;
            await this.loadComments();
            this.scrollToFirstComment();
        },

        /**
         * 下一页
         * 加载下一页评论并滚动到顶部
         * @returns {Promise<void>}
         */
        async nextPage() {
            this.currentPage++;
            await this.loadComments();
            this.scrollToFirstComment();
        },

        /**
         * 滚动到第一条评论
         * 使用平滑滚动效果
         */
        scrollToFirstComment() {
            this.$nextTick(() => {
                const firstComment = document.querySelector('.review-card');
                if (firstComment) {
                    firstComment.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start'
                    });
                }
            });
        },

        /**
         * 处理弹窗滚动事件
         * 实现无限滚动加载
         * @param {Event} e - 滚动事件对象
         */
        handleModalScroll(e) {
            const { scrollTop, scrollHeight, clientHeight } = e.target;
            // 当滚动到距离底部100px时加载更多
            if (scrollHeight - scrollTop - clientHeight < 100 && !this.isLoadingMore && this.hasMoreModalReplies) {
                this.modalCurrentPage++;
                this.loadModalReplies();
            }
        },
    },
    async created() {
        // 组件创建时加载初始数据
        await this.loadComments();
        await this.loadCurrentUser();
    }
}
</script>

<style scoped>
.reply-comment-input {
    margin: 12px 0;
}

.reply-comment-input .comment-input {
    display: flex;
    gap: 12px;
    align-items: flex-start;
    width: 100%;
    max-width: calc(100% - 48px);
}

.reply-comment-input .input-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
}

/* 简约输入框样式 */
.minimal-input {
    width: 100%;
    background: transparent;
    border: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    color: rgba(255, 255, 255, 0.9);
    padding: 8px 0;
    font-size: 14px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    min-width: 300px;
}

.minimal-input:focus {
    outline: none;
    border-bottom-color: #AB0511;
    transform: translateY(-1px);
    box-shadow: 0 2px 0 rgba(171, 5, 17, 0.3);
}

.minimal-input::placeholder {
    color: rgba(255, 255, 255, 0.5);
}

.reply-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 4px;
}

.cancel-btn {
    padding: 6px 12px;
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.6);
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.cancel-btn:hover {
    color: rgba(255, 255, 255, 0.9);
}

.submit-btn {
    padding: 6px 16px;
    background: transparent;
    border: none;
    color: #AB0511;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.submit-btn:hover:not(:disabled) {
    color: #8B0000;
}

.submit-btn:disabled {
    color: rgba(171, 5, 17, 0.3);
    cursor: not-allowed;
}

/* 回复框动画 */
.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    max-height: 200px;
    /* 设置最大高度用于动画 */
    overflow: hidden;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    max-height: 0;
    transform: translateY(-10px);
    margin: 0;
}

.fade-enter-to,
.fade-leave-from {
    opacity: 1;
    transform: translateY(0);
}

/* 回复框内容动画 */
.reply-comment-input,
.modal-reply-input {
    transform-origin: top;
    animation: slideDown 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideDown {
    from {
        opacity: 0;
        transform: translateY(-10px) scale(0.98);
    }

    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.me-icon {
    width: 20px;
    height: 20px;
    margin-left: 5px;
}

.pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 40px;
    margin-bottom: 40px;
}

.pagination {
    display: flex;
    gap: 8px;
    align-items: center;
}

.page-btn {
    background: rgba(139, 0, 0, 0.9);
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
    min-width: 40px;
    font-size: 0.95rem;
}

.page-btn:hover:not(:disabled) {
    background: rgb(159, 0, 0);
    transform: translateY(-2px);
    box-shadow: 0 4px 15px rgba(139, 0, 0, 0.3);
}

.page-btn:disabled {
    background: rgb(84 70 70 / 30%);
    cursor: not-allowed;
    transform: none;
}

.current-page {
    color: white;
    font-size: 1.1rem;
    font-weight: 500;
    padding: 0 12px;
}

.loading-indicator {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    color: rgba(255, 255, 255, 0.7);
}

.loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(255, 255, 255, 0.1);
    border-top-color: rgba(139, 0, 0, 0.9);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
}

.loading-replies {
    text-align: center;
    padding: 10px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
}

.no-comments {
    text-align: center;
    padding: 40px;
    color: rgba(255, 255, 255, 0.7);
    font-style: italic;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.reply-btn {
    opacity: 0.8;
    transition: opacity 0.3s ease;
}

.reply-btn:hover {
    opacity: 1;
}

.send-btn:disabled,
.post-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.reply-input-wrapper {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 12px;
    margin-top: 12px;
}

.reply-input {
    width: 100%;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 4px;
    color: white;
    padding: 10px 12px;
    margin-bottom: 12px;
    transition: all 0.3s ease;
}

.reply-input:focus {
    outline: none;
    border-color: rgba(139, 0, 0, 0.5);
    background: rgba(255, 255, 255, 0.15);
}

.reply-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

.cancel-btn {
    padding: 6px 16px;
    background: transparent;
    color: rgba(255, 255, 255, 0.7);
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.cancel-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: white;
}

.send-btn {
    padding: 6px 16px;
    background: rgba(139, 0, 0, 0.9);
    border: none;
    color: white;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.send-btn:hover:not(:disabled) {
    background: rgb(159, 0, 0);
    transform: translateY(-1px);
}

.get-replies-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    color: rgba(255, 255, 255, 0.7);
    transition: all 0.3s ease;
}

.get-replies-btn i {
    transition: transform 0.3s ease;
}

.get-replies-btn.expanded i {
    transform: rotate(180deg);
}

.reply-container {
    padding: 12px 0;
    margin-left: 40px;
}

.reply-card {
    margin-bottom: 12px;
}

.no-replies {
    text-align: center;
    padding: 20px;
    color: rgba(255, 255, 255, 0.5);
    font-style: italic;
}

.reply-relation {
    display: flex;
    align-items: center;
    gap: 6px;
    margin: 4px 0;
}

.reply-to {
    color: rgba(255, 255, 255, 0.5);
    font-size: 14px;
}

.replied-name {
    color: rgba(255, 255, 255, 0.8);
    font-size: 14px;
    font-weight: 500;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .reply-relation {
        gap: 4px;
    }

    .reply-to,
    .replied-name {
        font-size: 13px;
    }

    .me-icon {
        width: 16px;
        height: 16px;
    }

    .reply-comment-input {
        padding-left: 40px;
    }

    .reply-comment-input .comment-input {
        max-width: calc(100% - 40px);
    }

    .minimal-input {
        min-width: 200px;
    }
}

/* 超小屏幕适配 */
@media (max-width: 480px) {
    .reply-comment-input {
        padding-left: 32px;
    }

    .reply-comment-input .comment-input {
        max-width: calc(100% - 32px);
    }

    .reply-to,
    .replied-name {
        font-size: 12px;
    }
}

.replies-actions {
    display: flex;
    gap: 12px;
    align-items: center;
    margin-top: 8px;
}

.get-replies-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    color: rgba(255, 255, 255, 0.7);
    transition: all 0.3s ease;
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 4px 8px;
}

.get-replies-btn:hover {
    color: white;
}

.get-replies-btn i {
    font-size: 12px;
}

.loading-more-replies {
    text-align: center;
    padding: 10px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
    margin-top: 12px;
}

.loading-more-replies .loading-spinner {
    width: 24px;
    height: 24px;
    border-width: 2px;
    margin-bottom: 8px;
}

/* 添加弹窗相关样式 */
.reply-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.75);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}

.reply-modal-content {
    background: #1a1a1a;
    width: 90%;
    max-width: 600px;
    max-height: 80vh;
    border-radius: 12px;
    overflow: hidden;
    position: relative;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.reply-modal-header {
    padding: 16px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.reply-modal-header h3 {
    color: white;
    margin: 0;
    font-size: 1.2rem;
}

.close-btn {
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    padding: 8px;
    transition: all 0.3s ease;
}

.close-btn:hover {
    color: white;
}

.reply-modal-body {
    padding: 20px;
    overflow-y: auto;
    max-height: calc(80vh - 60px);

    /* 添加滚动条样式 */
    scrollbar-width: thin;
    /* Firefox */
    scrollbar-color: rgba(255, 255, 255, 0.5) rgba(255, 255, 255, 0.1);
    /* Firefox */
}

/* Webkit浏览器的滚动条样式 */
.reply-modal-body::-webkit-scrollbar {
    width: 6px;
}

.reply-modal-body::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
}

.reply-modal-body::-webkit-scrollbar-thumb {
    background: rgba(171, 5, 17, 0.5);
    border-radius: 3px;
    transition: background 0.3s ease;
}

.reply-modal-body::-webkit-scrollbar-thumb:hover {
    background: rgba(171, 5, 17, 0.8);
}

.original-comment {
    padding: 16px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    margin-bottom: 20px;
}

.modal-replies-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 12px;
}

.loading-more {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 16px;
    color: rgba(255, 255, 255, 0.7);
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    margin-top: 16px;
}

.loading-more .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-top-color: rgba(171, 5, 17, 0.7);
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

/* 底部提示样式 */
.no-more-replies {
    text-align: center;
    padding: 16px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 14px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    margin-top: 16px;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .reply-modal-content {
        width: 95%;
        max-height: 90vh;
    }

    .reply-modal-header {
        padding: 12px 16px;
    }

    .reply-modal-body {
        padding: 16px;
    }

    .reply-modal-body::-webkit-scrollbar {
        width: 4px;
        /* 移动端滚动条更细 */
    }

    .no-more-replies {
        font-size: 13px;
        padding: 12px;
    }

    .loading-more {
        padding: 12px;
    }

    .loading-more .loading-spinner {
        width: 18px;
        height: 18px;
    }
}

/* 修复评论头像被挤压的问题 */
.review-card {
    display: flex;
    gap: 12px;
}

.review-card>div:first-child {
    flex: 0 0 36px;
    /* 固定宽度，不允许缩放 */
}

.review-card>.textBlock {
    flex: 1;
    min-width: 0;
    /* 允许文本内容自动换行 */
}

/* 保持原有的头像样式 */
.add-comment img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
}

.review-card img {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    object-fit: cover;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .review-card>div:first-child {
        flex: 0 0 32px;
    }

    .review-card img {
        width: 32px;
        height: 32px;
    }
}

@media (max-width: 480px) {
    .review-card>div:first-child {
        flex: 0 0 28px;
    }

    .review-card img {
        width: 28px;
        height: 28px;
    }
}

/* 评论作者信息布局优化 */
.comment-author,
.reply-author {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 8px;
}

.comment-meta,
.reply-meta {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.comment-meta h4,
.reply-meta h4 {
    margin: 0;
    font-size: 0.95rem;
    color: rgba(255, 255, 255, 0.9);
}

.comment-time,
.reply-time {
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.6);
}

/* 修改弹窗动画 */
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

/* 弹窗内容滑动动画 */
.slide-down-enter-active,
.slide-down-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(-20vh);
    opacity: 0;
}

/* 添加作者标签样式 */
.author-tag {
    display: inline-block;
    font-size: 12px;
    padding: 2px 6px;
    margin-left: 8px;
    background: rgba(171, 5, 17, 0.15);
    color: #AB0511;
    border: 1px solid rgba(171, 5, 17, 0.3);
    border-radius: 4px;
    font-weight: normal;
    vertical-align: middle;
    line-height: 1.2;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .author-tag {
        font-size: 11px;
        padding: 1px 4px;
        margin-left: 6px;
    }
}

/* 超小屏幕适配 */
@media (max-width: 480px) {
    .author-tag {
        font-size: 10px;
        padding: 1px 3px;
        margin-left: 4px;
    }
}

/* 修改弹窗回复框样式 */
.modal-reply-input {
    margin: 12px 0;
    padding-left: 48px;
}

.modal-reply-input .comment-input {
    display: flex;
    gap: 12px;
    align-items: flex-start;
}

.modal-reply-input img {
    width: 32px;
    height: 32px;
    min-width: 32px;
    /* 添加最小宽度，防止压缩 */
    border-radius: 50%;
    object-fit: cover;
}

.modal-reply-input .input-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .modal-reply-input {
        padding-left: 40px;
    }

    .modal-reply-input img {
        width: 28px;
        height: 28px;
        min-width: 28px;
    }
}

@media (max-width: 480px) {
    .modal-reply-input {
        padding-left: 32px;
    }

    .modal-reply-input img {
        width: 24px;
        height: 24px;
        min-width: 24px;
    }
}

/* 修改弹窗中的头像样式 */
.reply-modal .comment-author,
.reply-modal .reply-author {
    display: flex;
    align-items: flex-start;
    gap: 12px;
}

.reply-modal .comment-author>img,
.reply-modal .reply-author>img {
    width: 36px;
    height: 36px;
    min-width: 36px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
}

.reply-modal .comment-meta,
.reply-modal .reply-meta {
    flex: 1;
    min-width: 0;
}

/* 移动端适配 */
@media (max-width: 768px) {

    .reply-modal .comment-author>img,
    .reply-modal .reply-author>img {
        width: 32px;
        height: 32px;
        min-width: 32px;
    }
}

@media (max-width: 480px) {

    .reply-modal .comment-author>img,
    .reply-modal .reply-author>img {
        width: 28px;
        height: 28px;
        min-width: 28px;
    }
}
</style>
