<template>
    <div class="comments-section">
        <h3 class="section-title">
            <i class="fas fa-comments"></i>
            评论 ({{ post.comment_count }})
        </h3>

        <!-- 评论输入框 -->
        <div class="comment-input">
            <img :src="currentUser.avatar_url" :alt="currentUser.username">
            <div class="input-wrapper">
                <div class="textarea-container">
                    <textarea v-model="commentText" placeholder="发表你的评论..." @input="autoResizeReply($event)"
                        ref="commentTextarea" maxlength="500"></textarea>
                    <div class="word-count" :class="{ 'near-limit': commentText.length > 450 }">
                        {{ commentText.length }}/500
                    </div>
                </div>
                <button class="submit-btn" :disabled="!commentText.trim() || commentText.length > 500"
                    @click="submitComment">
                    发表评论
                </button>
            </div>
        </div>

        <!-- 评论列表 -->
        <div class="comments-list">
            <!-- 添加一个用于滚动定位的元素 -->
            <div ref="commentAnchor"></div>

            <div v-for="comment in comments" :key="comment.id" class="comment-item">
                <div class="comment-author">
                    <img :src="comment.author.avatar_url" :alt="comment.author.username">
                    <div class="comment-meta">
                        <h4>{{ comment.author.username }}</h4>
                        <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
                    </div>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
                <div class="comment-actions">
                    <button class="action-btn" @click="replyToComment(comment)">
                        <i class="fas fa-reply"></i>
                        <span>回复</span>
                    </button>
                    <button class="action-btn" :class="{ active: comment.is_liked }"
                        @click="toggleCommentLike(comment)">
                        <i class="fas fa-heart"></i>
                        <span class="like-count">{{ comment.like_count }}</span>
                    </button>
                    <button v-if="false" class="action-btn delete" @click="deleteComment(comment)">
                        <i class="fas fa-trash-alt"></i>
                        <span>删除</span>
                    </button>
                </div>

                <!-- 修改展开回复按钮 -->
                <div v-if="comment.reply_num > 0" class="expand-replies">
                    <button class="expand-btn" @click="showReplyDetail(comment)">
                        <i class="fas fa-chevron-down"></i>
                        查看全部 {{ comment.reply_num }} 条回复
                    </button>
                </div>

                <!-- 添加回复框组件 -->
                <div v-if="comment.showReplyBox" class="reply-input-box">
                    <img :src="currentUser.avatar_url" :alt="currentUser.username">
                    <div class="input-wrapper">
                        <div class="textarea-container">
                            <textarea v-model="comment.replyText"
                                :placeholder="comment.replyTo ? `回复 @${comment.replyTo.username}:` : '发表你的回复...'"
                                @input="autoResizeReply($event)" maxlength="500"></textarea>
                            <div class="word-count" :class="{ 'near-limit': comment.replyText?.length > 450 }">
                                {{ comment.replyText?.length || 0 }}/500
                            </div>
                        </div>
                        <div class="reply-actions-wrapper">
                            <button class="cancel-btn" @click="cancelReply(comment)">取消</button>
                            <button class="submit-btn"
                                :disabled="!comment.replyText?.trim() || comment.replyText?.length > 500"
                                @click="submitReply(comment)">
                                发表回复
                            </button>
                        </div>
                    </div>
                </div>

                <!-- 评论的回复 -->
                <div v-if="comment.replies && comment.replies.length" class="replies-list">
                    <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                        <div class="reply-author">
                            <img :src="reply.author.avatar_url" :alt="reply.author.username">
                            <div class="reply-meta">
                                <h4>{{ reply.author.username }}</h4>
                                <span class="reply-time">{{ formatTime(reply.created_at) }}</span>
                            </div>
                        </div>
                        <div class="reply-content">
                            <span class="reply-to" v-if="reply.reply_to">
                                回复 @{{ reply.reply_to.username }}:
                            </span>
                            {{ reply.content }}
                        </div>
                        <div class="reply-actions">
                            <button class="action-btn" @click="replyToComment(reply, comment)">
                                <i class="fas fa-reply"></i>
                                <span>回复</span>
                            </button>
                            <button class="action-btn" :class="{ active: reply.is_liked }"
                                @click="toggleCommentLike(reply)">
                                <i class="fas fa-heart"></i>
                                <span class="like-count">{{ reply.like_count }}</span>
                            </button>
                            <button v-if="false" class="action-btn delete" @click="deleteComment(reply)">
                                <i class="fas fa-trash-alt"></i>
                                <span>删除</span>
                            </button>
                        </div>

                        <!-- 在每个回复下面添加回复框 -->
                        <div v-if="reply.showReplyBox" class="reply-input-box">
                            <img :src="currentUser.avatar_url" :alt="currentUser.username">
                            <div class="input-wrapper">
                                <div class="textarea-container">
                                    <textarea v-model="reply.replyText" :placeholder="`回复 @${reply.author.username}:`"
                                        @input="autoResizeReply($event)" maxlength="500"></textarea>
                                    <div class="word-count" :class="{ 'near-limit': reply.replyText?.length > 450 }">
                                        {{ reply.replyText?.length || 0 }}/500
                                    </div>
                                </div>
                                <div class="reply-actions-wrapper">
                                    <button class="cancel-btn" @click="cancelReply(reply)">取消</button>
                                    <button class="submit-btn"
                                        :disabled="!reply.replyText?.trim() || reply.replyText?.length > 500"
                                        @click="submitReply(reply, comment)">
                                        发表回复
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 修改分页部分 -->
        <div v-if="totalPages > 1" class="pagination">
            <button class="page-btn prev" :disabled="currentPage === 1" @click="changePage(currentPage - 1)">
                <i class="fas fa-chevron-left"></i>
            </button>

            <div class="page-numbers">
                <button v-for="pageNum in displayedPages" :key="pageNum" class="page-btn" :class="{
                    active: pageNum === currentPage,
                    ellipsis: pageNum === '...'
                }" :disabled="pageNum === '...'" @click="changePage(pageNum)">
                    {{ pageNum }}
                </button>
            </div>

            <button class="page-btn next" :disabled="currentPage === totalPages" @click="changePage(currentPage + 1)">
                <i class="fas fa-chevron-right"></i>
            </button>
        </div>
    </div>

    <!-- 添加回复详情弹窗 -->
    <div v-if="showReplyModal" class="reply-modal">
        <div class="reply-modal-content">
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
                        <img :src="currentComment?.author.avatar_url" :alt="currentComment?.author.username">
                        <div class="comment-meta">
                            <h4>{{ currentComment?.author.username }}</h4>
                            <span class="comment-time">{{ formatTime(currentComment?.created_at) }}</span>
                        </div>
                    </div>
                    <div class="comment-content">{{ currentComment?.content }}</div>
                </div>

                <!-- 回复列表 -->
                <div class="modal-replies-list">
                    <div v-for="reply in modalReplies" :key="reply.id" :id="'reply-' + reply.id" class="reply-item">
                        <div class="reply-author">
                            <img :src="reply.author.avatar_url" :alt="reply.author.username">
                            <div class="reply-meta">
                                <h4>{{ reply.author.username }}</h4>
                                <span class="reply-time">{{ formatTime(reply.created_at) }}</span>
                            </div>
                        </div>
                        <div class="reply-content">
                            <span class="reply-to" v-if="reply.reply_to">
                                回复 @{{ reply.reply_to.username }}:
                            </span>
                            {{ reply.content }}
                        </div>
                        <div class="reply-actions">
                            <button class="action-btn" @click="replyToModalComment(reply)">
                                <i class="fas fa-reply"></i>
                                <span>回复</span>
                            </button>
                            <button class="action-btn" :class="{ active: reply.is_liked }"
                                @click="toggleCommentLike(reply)">
                                <i class="fas fa-heart"></i>
                                <span class="like-count">{{ reply.like_count }}</span>
                            </button>
                            <button v-if="false" class="action-btn delete" @click="deleteComment(reply)">
                                <i class="fas fa-trash-alt"></i>
                                <span>删除</span>
                            </button>
                        </div>

                        <!-- 在每条回复下添加回复框 -->
                        <div v-if="reply.showReplyBox" class="reply-input-box">
                            <img :src="currentUser.avatar_url" :alt="currentUser.username">
                            <div class="input-wrapper">
                                <div class="textarea-container">
                                    <textarea v-model="reply.replyText" :placeholder="`回复 @${reply.author.username}:`"
                                        @input="autoResizeReply($event)" maxlength="500"></textarea>
                                    <div class="word-count" :class="{ 'near-limit': reply.replyText?.length > 450 }">
                                        {{ reply.replyText?.length || 0 }}/500
                                    </div>
                                </div>
                                <div class="reply-actions-wrapper">
                                    <button class="cancel-btn" @click="cancelModalReply(reply)">取消</button>
                                    <button class="submit-btn"
                                        :disabled="!reply.replyText?.trim() || reply.replyText?.length > 500"
                                        @click="submitModalReply(reply, currentComment)">
                                        发表回复
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 加载状态提示 -->
                    <div v-if="isLoadingMore" class="loading-more">
                        <i class="fas fa-spinner fa-spin"></i>
                        加载中...
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Comment',
    emits: [
        'submit-comment',
        'toggle-like',
        'delete-comment',
        'submit-reply',
        'change-page',
        'load-replies'
    ],
    props: {
        // 当前帖子信息
        post: {
            type: Object,
            required: true
        },
        // 评论列表
        comments: {
            type: Array,
            default: () => []
        },
        // 当前登录用户信息
        currentUser: {
            type: Object,
            required: true
        },
        // 是否还有更多评论可以加载
        hasMoreComments: {
            type: Boolean,
            default: false
        },
        // 当前页码
        currentPage: {
            type: Number,
            required: true
        },
        // 总页数
        totalPages: {
            type: Number,
            required: true
        }
    },
    data() {
        return {
            commentText: '', // 评论输入框的文本内容
            showReplyModal: false, // 是否显示回复弹窗
            currentComment: null, // 当前选中的评论
            modalReplies: [], // 弹窗中的回复列表
            modalCurrentPage: 1, // 弹窗中的当前页码
            isLoadingMore: false, // 是否正在加载更多回复
            hasMoreReplies: true, // 是否还有更多回复可以加载
            modalReplyText: '', // 弹窗中回复输入框的文本内容
            replyToUser: null, // 要回复的用户信息
            virtualID: -1, // 虚拟评论ID,用于临时标识新创建的评论
        }
    },
    computed: {
        displayedPages() {
            const pages = [];
            const maxDisplayPages = 5; // 最多显示5个页码

            if (this.totalPages <= maxDisplayPages) {
                // 如果总页数小于等于5，显示所有页码
                for (let i = 1; i <= this.totalPages; i++) {
                    pages.push(i);
                }
            } else {
                // 始终显示第一页
                pages.push(1);

                // 计算中间页码
                let start = Math.max(this.currentPage - 1, 2);
                let end = Math.min(start + 2, this.totalPages - 1);

                // 调整 start，确保显示3个中间页码
                start = Math.max(Math.min(end - 2, start), 2);

                // 添加省略号和中间页码
                if (start > 2) pages.push('...');
                for (let i = start; i <= end; i++) {
                    pages.push(i);
                }
                if (end < this.totalPages - 1) pages.push('...');

                // 始终显示最后一页
                pages.push(this.totalPages);
            }

            return pages;
        }
    },
    methods: {
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
         * 重置文本框
         */
        resetTextarea() {
            this.commentText = ''; // 清空内容
            const textarea = this.$refs.commentTextarea;
            textarea.style.height = '60px'; // 重置为初始高度
        },

        /**
         * 提交新评论
         */
        submitComment() {
            if (this.commentText.length > 500) {
                return;
            }
            if (!this.commentText.trim()) return;
            this.$emit('submit-comment', {
                id: this.virtualID--, // 使用虚拟ID
                content: this.commentText,
                author: this.currentUser,
                created_at: new Date().toISOString()
            });
            // 提交成功后重置文本框
            this.resetTextarea();
        },

        /**
         * 判断评论是否由当前用户发布
         * @param {Object} comment - 评论对象
         * @returns {boolean} 是否为评论作者
         */
        isCommentAuthor(comment) {
            return comment.author.id === this.currentUser.id;
        },

        /**
         * 切换评论点赞状态
         * @param {Object} comment - 评论对象
         */
        toggleCommentLike(comment) {
            this.$emit('toggle-like', comment);
        },

        /**
         * 删除评论
         * @param {Object} comment - 评论对象
         */
        deleteComment(comment) {
            this.$emit('delete-comment', comment);
        },

        /**
         * 回复评论
         * @param {Object} reply - 要回复的评论对象
         * @param {Object} parentComment - 父评论对象(如果有)
         */
        replyToComment(reply, parentComment) {
            // 关闭其他所有回复框
            this.comments.forEach(c => {
                c.showReplyBox = false;
                c.replyText = '';
                if (c.replies) {
                    c.replies.forEach(r => {
                        r.showReplyBox = false;
                        r.replyText = '';
                    });
                }
            });

            // 如果是回复主评论
            if (!parentComment) {
                reply.showReplyBox = true;
                reply.replyText = '';
                reply.replyTo = null;
            }
            // 如果是回复子评论
            else {
                reply.showReplyBox = true;
                reply.replyText = '';
                reply.replyTo = reply;
            }
        },

        /**
         * 切换评论页码
         * @param {number} page - 目标页码
         */
        changePage(page) {
            if (page === this.currentPage || page < 1 || page > this.totalPages) return;
            this.$emit('change-page', page);

            // 等待评论加载完成后滚动
            this.$nextTick(() => {
                this.scrollToComments();
            });
        },

        /**
         * 取消回复
         * @param {Object} comment - 评论对象
         */
        cancelReply(comment) {
            comment.showReplyBox = false;
            comment.replyText = '';
        },

        /**
         * 提交评论回复
         * @param {Object} reply - 要回复的评论对象（可以是主评论或子评论）
         * @param {Object} rootComment - 所属的根评论对象
         */
        submitReply(reply, rootComment) {
            if (!reply.replyText.trim()) return;

            const replyData = {
                id: this.virtualID--, // 使用虚拟ID
                content: reply.replyText,
                parent_id: reply.id,
                root_id: rootComment ? rootComment.id : reply.id,
                to_user_id: reply.author.id,
                created_at: new Date().toISOString(),
                reply_to: reply.author
            };

            this.$emit('submit-reply', replyData);
            this.cancelReply(reply);
        },

        /**
         * 滚动到评论区域
         */
        scrollToComments() {
            const anchor = this.$refs.commentAnchor;
            if (anchor) {
                anchor.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        },

        /**
         * 加载评论回复
         * @param {Object} comment - 评论对象
         */
        async loadReplies(comment) {
            // 发出加载回复的事件
            this.$emit('load-replies', comment);
        },

        /**
         * 显示回复详情弹窗
         * @param {Object} comment - 评论对象
         */
        showReplyDetail(comment) {
            this.currentComment = comment;
            this.showReplyModal = true;
            this.modalCurrentPage = 1;
            this.modalReplies = [];
            this.hasMoreReplies = true;
            this.loadModalReplies();
        },

        /**
         * 关闭回复详情弹窗
         */
        closeReplyModal() {
            this.showReplyModal = false;
            this.currentComment = null;
            this.modalReplies = [];
            this.modalCurrentPage = 1;
            this.hasMoreReplies = true;
            this.modalReplyText = '';
            this.replyToUser = null;
        },

        /**
         * 加载弹窗中的回复列表
         */
        async loadModalReplies() {
            if (this.isLoadingMore || !this.hasMoreReplies) return;

            this.isLoadingMore = true;
            this.$emit('load-replies', {
                comment: this.currentComment,
                page: this.modalCurrentPage
            });
        },

        /**
         * 处理弹窗滚动事件
         * @param {Event} e - 滚动事件对象
         */
        handleModalScroll(e) {
            const { scrollTop, scrollHeight, clientHeight } = e.target;
            // 当滚动到距离底部100px时加载更多
            if (scrollHeight - scrollTop - clientHeight < 100 && !this.isLoadingMore && this.hasMoreReplies) {
                this.modalCurrentPage++;
                this.loadModalReplies();
            }
        },

        /**
         * 处理加载回复的结果
         * @param {Array} replies - 回复列表
         */
        handleLoadRepliesResult(replies) {
            // 如果返回的数据少于10条（假设每页10条），说明没有更多数据了
            const pageSize = 10;
            this.modalReplies = [...this.modalReplies, ...replies];
            this.hasMoreReplies = replies.length === pageSize;
            this.isLoadingMore = false;
        },

        /**
         * 回复弹窗中的评论
         * @param {Object} reply - 回复对象
         */
        replyToModalComment(reply) {
            // 关闭其他所有回复框
            this.modalReplies.forEach(r => {
                if (r.id !== reply.id) {
                    r.showReplyBox = false;
                    r.replyText = '';
                }
            });

            // 切换当前回复的回复框
            reply.showReplyBox = !reply.showReplyBox;
            if (reply.showReplyBox) {
                reply.replyText = '';
                this.$nextTick(() => {
                    // 滚动到回复框
                    const replyBox = document.querySelector(`#reply-${reply.id} .reply-input-box`);
                    if (replyBox) {
                        replyBox.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    }
                });
            }
        },

        /**
         * 取消弹窗中的回复
         * @param {Object} reply - 回复对象
         */
        cancelModalReply(reply) {
            reply.showReplyBox = false;
            reply.replyText = '';
        },

        /**
         * 提交弹窗中的回复
         * @param {Object} reply - 回复对象
         * @param {Object} rootComment - 根评论对象
         */
        submitModalReply(reply, rootComment) {
            if (!reply.replyText?.trim() || reply.replyText?.length > 500) return;

            const replyData = {
                id: this.virtualID--, // 使用虚拟ID
                content: reply.replyText,
                parent_id: reply.id,
                root_id: rootComment ? rootComment.id : reply.id,
                to_user_id: reply.author.id,
                created_at: new Date().toISOString(),
                reply_to: reply.author
            };

            // 构造新回复对象
            const newReply = {
                id: replyData.id, // 使用虚拟ID
                content: reply.replyText,
                author: this.currentUser,
                created_at: new Date().toISOString(),
                reply_to: reply.author,
                like_count: 0,
                is_liked: false,
                showReplyBox: false,
                replyText: ''
            };

            // 找到当前回复在数组中的索引
            const currentIndex = this.modalReplies.findIndex(r => r.id === reply.id);

            // 在当前回复后插入新回复
            if (currentIndex !== -1) {
                this.modalReplies.splice(currentIndex + 1, 0, newReply);
            } else {
                this.modalReplies.push(newReply);
            }

            // 发送事件到父组件
            this.$emit('submit-reply', replyData);
            this.cancelModalReply(reply);

            // 滚动到新回复的位置
            this.$nextTick(() => {
                const newReplyElement = document.querySelector(`#reply-${newReply.id}`);
                if (newReplyElement) {
                    newReplyElement.scrollIntoView({
                        behavior: 'smooth',
                        block: 'center'
                    });
                }
            });
        },

        /**
         * 自动调整回复框高度
         * @param {Event} event - 输入事件
         * @param {Object} reply - 回复对象
         */
        autoResizeReply(event) {
            const textarea = event.target;
            textarea.style.height = 'auto';
            textarea.style.height = textarea.scrollHeight + 2 + 'px';
        }
    },
    watch: {
        // 监听当前页码变化
        currentPage: {
            handler(newPage, oldPage) {
                if (newPage !== oldPage) {
                    this.$nextTick(() => {
                        this.scrollToComments();
                    });
                }
            }
        }
    }
}
</script>

<style scoped>
@import "@/static/css/postDetail.css";

/* 添加滚动定位元素的样式 */
.comments-list {
    position: relative;
}

.comments-list>div:first-child {
    position: absolute;
    top: -20px;
    /* 调整滚动位置的偏移量 */
}

/* 添加省略号按钮样式 */
.page-btn.ellipsis {
    background: none;
    cursor: default;
    pointer-events: none;
}

.page-btn.ellipsis:hover {
    transform: none;
    background: none;
}

/* 添加展开回复按钮样式 */
.expand-replies {
    padding-left: 52px;
    margin-top: 12px;
}

.expand-btn {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 6px 12px;
    border-radius: 20px;
    transition: all 0.3s ease;
}

.expand-btn:hover {
    background: rgba(255, 255, 255, 0.08);
    color: white;
}

.expand-btn i {
    font-size: 0.8rem;
    transition: transform 0.3s ease;
}

.expand-btn:hover i {
    transform: translateY(2px);
}

/* 响应式调整 */
@media (max-width: 768px) {
    .expand-replies {
        padding-left: 25px;
    }

    .expand-btn {
        font-size: 0.85rem;
    }
}

.input-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.textarea-container {
    position: relative;
    width: 100%;
}

textarea {
    width: 100%;
    height: 60px;
    /* 设置初始高度 */
    min-height: 60px;
    max-height: 300px;
    padding: 12px 40px 12px 12px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    color: white;
    font-size: 14px;
    line-height: 1.5;
    resize: none;
    overflow-y: hidden;
    transition: all 0.3s ease;
}

textarea:focus {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.2);
    outline: none;
}

.word-count {
    position: absolute;
    bottom: 8px;
    right: 12px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.5);
    pointer-events: none;
    /* 防止干扰文本输入 */
    transition: color 0.3s ease;
}

.word-count.near-limit {
    color: #ff4d4d;
    /* 接近字数限制时显示红色 */
}

.submit-btn {
    align-self: flex-end;
    padding: 8px 20px;
    background: linear-gradient(45deg, #8b0000, #cc0000);
    border: none;
    border-radius: 20px;
    color: white;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.submit-btn:disabled {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.3);
    cursor: not-allowed;
}

.submit-btn:not(:disabled):hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(139, 0, 0, 0.3);
}

/* 移动端适配 */
@media (max-width: 768px) {
    textarea {
        font-size: 16px;
        /* 移动端增大字体，防止自动缩放 */
        padding: 10px 35px 10px 10px;
    }

    .word-count {
        bottom: 6px;
        right: 8px;
    }
}

.reply-input-box {
    margin-left: 40px;
    margin-top: 12px;
    display: flex;
    gap: 12px;
}

.reply-input-box img {
    width: 32px;
    height: 32px;
    border-radius: 50%;
}

.reply-input-box .input-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.textarea-container {
    position: relative;
    width: 100%;
}

.textarea-container textarea {
    width: 100%;
    height: 60px;
    min-height: 60px;
    max-height: 200px;
    padding: 12px 40px 12px 12px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    color: white;
    font-size: 14px;
    line-height: 1.5;
    resize: none;
    overflow-y: hidden;
    transition: all 0.3s ease;
}

.textarea-container textarea:focus {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.2);
    outline: none;
}

.word-count {
    position: absolute;
    bottom: 8px;
    right: 12px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.5);
    pointer-events: none;
}

.word-count.near-limit {
    color: #ff4d4d;
}

.reply-actions-wrapper {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .reply-input-box {
        margin-left: 20px;
    }

    .textarea-container textarea {
        font-size: 16px;
        padding: 10px 35px 10px 10px;
    }
}
</style>
