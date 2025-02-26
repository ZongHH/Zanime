<template>
    <div class="modal-overlay" @click.self="$emit('close')">
        <div class="modal-content">
            <div class="modal-header">
                <h2>发布新帖子</h2>
                <button class="close-btn" @click="$emit('close')">
                    <i class="fas fa-times"></i>
                </button>
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label>标题</label>
                    <input type="text" v-model="title" placeholder="请输入帖子标题...">
                </div>

                <div class="form-group">
                    <label>选择分类</label>
                    <div class="tags-section">
                        <!-- 热门分类标签 -->
                        <div class="hot-tags">
                            <h4>热门分类</h4>
                            <div class="tags-container">
                                <button v-for="tag in availableTags" :key="tag.id"
                                    :class="['tag-btn', { active: selectedTags.includes(tag.name) }]"
                                    @click="addTag(tag.name)">
                                    {{ tag.name }}
                                </button>
                            </div>
                        </div>

                        <!-- 自定义分类 -->
                        <div class="custom-tags">
                            <h4>已选分类</h4>
                            <div class="selected-tags">
                                <div v-for="(tag, index) in selectedTags" :key="index" class="selected-tag">
                                    <span>{{ tag }}</span>
                                    <button class="remove-tag" @click="removeTag(index)">
                                        <i class="fas fa-times"></i>
                                    </button>
                                </div>
                            </div>

                            <!-- 添加新分类 -->
                            <div class="add-tag-section" v-if="selectedTags.length < 5">
                                <div v-if="showTagInput" class="tag-input-group">
                                    <input ref="tagInput" v-model="newTag" @keyup.enter="confirmNewTag"
                                        @blur="handleTagInputBlur" placeholder="输入分类名称（按回车确认）" maxlength="10">
                                    <span class="char-count">{{ newTag.length }}/10</span>
                                </div>
                                <button v-else class="add-tag-btn" @click="showTagInput = true">
                                    <i class="fas fa-plus"></i> 添加分类
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="form-group">
                    <label>内容</label>
                    <textarea v-model="content" placeholder="请输入帖子内容..."></textarea>
                </div>

                <div class="form-group">
                    <label>添加图片</label>
                    <div class="image-upload">
                        <div class="upload-area" @click="triggerFileInput">
                            <i class="fas fa-cloud-upload-alt"></i>
                            <p>点击或拖拽图片到此处</p>
                            <span>支持 jpg、png、gif 格式，最多 9 张</span>
                        </div>
                        <input type="file" ref="fileInput" multiple accept="image/*" @change="handleFileUpload" hidden>

                        <div class="image-preview" v-if="uploadedImages.length">
                            <div v-for="(image, index) in uploadedImages" :key="index" class="preview-item">
                                <img :src="image.url" :alt="image.name">
                                <button class="remove-btn" @click="removeImage(index)">
                                    <i class="fas fa-times"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="modal-footer">
                <button class="cancel-btn" @click="$emit('close')">取消</button>
                <button class="submit-btn" @click="submitPost" :disabled="!isValid">
                    <i class="fas fa-paper-plane"></i>
                    发布
                </button>
            </div>
        </div>

        <!-- 添加 Toast 提示组件 -->
        <div v-if="showToast" :class="['toast-message', toastType]">
            <i :class="toastIcon"></i>
            <span>{{ toastMessage }}</span>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    // 组件名称
    name: 'CreatePostModal',
    // 组件接收的属性
    props: {
        currentCategoryId: {
            type: Number,
            required: true
        }
    },
    // 组件数据
    data() {
        return {
            title: '', // 帖子标题
            content: '', // 帖子内容
            selectedTags: [], // 已选择的标签
            uploadedImages: [], // 已上传的图片
            newTag: '', // 新标签输入值
            showTagInput: false, // 是否显示标签输入框
            availableTags: [ // 预设的标签列表
                { id: 1, name: '动画' },
                { id: 2, name: '漫画' },
                { id: 3, name: '轻小说' },
                { id: 4, name: '周边' },
                { id: 5, name: '讨论' },
                { id: 6, name: '资讯' },
                { id: 7, name: '安利' },
                { id: 8, name: '杂谈' }
            ],
            showToast: false,
            toastMessage: '',
            toastType: 'success', // 'success' 或 'error'
        }
    },
    // 计算属性
    computed: {
        // 验证表单是否有效
        isValid() {
            return this.title.trim() && this.content.trim() && this.selectedTags.length > 0;
        },
        toastIcon() {
            return {
                'success': 'fas fa-check-circle',
                'error': 'fas fa-times-circle'
            }[this.toastType];
        }
    },
    // 组件方法
    methods: {
        /**
         * 添加预设标签到已选标签列表中
         * @param {string} tagName - 要添加的标签名称
         * @description 如果标签不存在且已选标签数量小于5个,则添加该标签
         */
        addTag(tagName) {
            if (!this.selectedTags.includes(tagName) && this.selectedTags.length < 5) {
                this.selectedTags.push(tagName);
            }
        },

        /**
         * 从已选标签列表中移除指定索引的标签
         * @param {number} index - 要移除的标签索引
         */
        removeTag(index) {
            this.selectedTags.splice(index, 1);
        },

        /**
         * 确认并添加新的自定义标签
         * @description 
         * 1. 对输入的标签进行trim处理
         * 2. 检查标签是否有效且不重复且未超出限制
         * 3. 添加标签后清空输入框并隐藏
         */
        confirmNewTag() {
            const tag = this.newTag.trim();
            if (tag && !this.selectedTags.includes(tag) && this.selectedTags.length < 5) {
                this.selectedTags.push(tag);
                this.newTag = '';
                this.showTagInput = false;
            }
        },

        /**
         * 处理标签输入框失焦事件
         * @description 
         * 1. 如果输入框有内容,则尝试添加标签
         * 2. 如果输入框为空,则隐藏输入框
         */
        handleTagInputBlur() {
            if (this.newTag.trim()) {
                this.confirmNewTag();
            } else {
                this.showTagInput = false;
            }
        },

        /**
         * 触发文件选择器打开
         * @description 通过程序触发隐藏的文件input的点击事件
         */
        triggerFileInput() {
            this.$refs.fileInput.click();
        },

        /**
         * 处理文件上传事件
         * @param {Event} event - 文件上传事件对象
         * @description
         * 1. 将FileList转换为数组进行处理
         * 2. 检查上传图片数量限制(最多9张)
         * 3. 使用FileReader读取图片并添加到预览列表
         */
        handleFileUpload(event) {
            const files = Array.from(event.target.files);
            files.forEach(file => {
                if (this.uploadedImages.length >= 9) return;

                const reader = new FileReader();
                reader.onload = (e) => {
                    this.uploadedImages.push({
                        url: e.target.result,
                        name: file.name,
                        file: file
                    });
                };
                reader.readAsDataURL(file);
            });
        },

        /**
         * 移除已上传的图片
         * @param {number} index - 要移除的图片索引
         */
        removeImage(index) {
            this.uploadedImages.splice(index, 1);
        },

        /**
         * 显示提示信息
         * @param {string} message - 提示信息
         * @param {string} type - 提示类型 ('success' 或 'error')
         */
        showToastMessage(message, type = 'success') {
            this.toastMessage = message;
            this.toastType = type;
            this.showToast = true;

            // 3秒后自动关闭
            setTimeout(() => {
                this.showToast = false;
            }, 3000);
        },

        /**
         * 修改提交方法
         * @description
         * 1. 验证表单数据是否有效
         * 2. 构建请求数据对象,包含用户ID、分类ID、标题、内容、标签和图片
         * 3. 发送POST请求到服务器创建帖子
         * 4. 处理成功/失败响应
         * 5. 关闭模态框
         */
        async submitPost() {
            if (!this.isValid) return;

            try {
                const request = {
                    category_id: this.currentCategoryId,
                    title: this.title,
                    content: this.content,
                    tags: this.selectedTags,
                    images: this.uploadedImages.map(img => img.url)
                };

                const response = await axios.post('/api/create-post', request);

                if (response.data.code === 200) {
                    this.showToastMessage('发布成功！');
                    setTimeout(() => {
                        this.$emit('close');
                    }, 1500);
                } else {
                    this.showToastMessage(response.data.message || '发布失败，请重试', 'error');
                }
            } catch (error) {
                this.showToastMessage(error.message || '发布失败，请重试', 'error');
            }
        }
    }
}
</script>

<style scoped>
.modal-overlay {
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
}

.modal-content {
    background: linear-gradient(145deg, rgba(30, 30, 30, 0.95), rgba(20, 20, 20, 0.95));
    border-radius: 20px;
    width: 100%;
    max-width: 800px;
    height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.05);
    animation: modalShow 0.3s ease;
}

.modal-header {
    padding: 24px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-shrink: 0;
}

.modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    scrollbar-width: thin;
    scrollbar-color: rgba(15, 15, 15, 0.5) rgba(255, 255, 255, 0.1);
}

.modal-body::-webkit-scrollbar {
    width: 8px;
}

.modal-body::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb {
    background: linear-gradient(45deg, rgba(139, 0, 0, 0.8), rgba(204, 0, 0, 0.8));
    border-radius: 4px;
    border: 2px solid transparent;
    background-clip: padding-box;
}

.modal-body::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(45deg, rgba(139, 0, 0, 0.9), rgba(204, 0, 0, 0.9));
}

.modal-footer {
    padding: 24px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: flex-end;
    gap: 16px;
    flex-shrink: 0;
}

.modal-header h2 {
    color: white;
    font-size: 1.5rem;
    font-weight: 600;
}

.close-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
}

.close-btn:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: rotate(90deg);
}

.form-group {
    margin-bottom: 24px;
}

.form-group label {
    display: block;
    color: rgba(255, 255, 255, 0.9);
    margin-bottom: 8px;
    font-weight: 500;
}

.form-group input,
.form-group textarea {
    width: 100%;
    padding: 12px 16px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
    font-size: 1rem;
    transition: all 0.3s ease;
    scrollbar-width: thin;
    scrollbar-color: rgba(139, 0, 0, 0.5) rgba(255, 255, 255, 0.1);
}

.form-group textarea {
    height: 200px;
    resize: vertical;
    line-height: 1.6;
}

.form-group input:focus,
.form-group textarea:focus {
    background: rgba(255, 255, 255, 0.12);
    border-color: rgba(139, 0, 0, 0.5);
    box-shadow: 0 0 0 2px rgba(139, 0, 0, 0.2);
    outline: none;
}

.form-group textarea::-webkit-scrollbar {
    width: 8px;
}

.form-group textarea::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
}

.form-group textarea::-webkit-scrollbar-thumb {
    background: linear-gradient(45deg, rgba(139, 0, 0, 0.8), rgba(204, 0, 0, 0.8));
    border-radius: 4px;
    border: 2px solid transparent;
    background-clip: padding-box;
}

.form-group textarea::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(45deg, rgba(139, 0, 0, 0.9), rgba(204, 0, 0, 0.9));
}

.tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
}

.tag-btn {
    padding: 8px 16px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
    cursor: pointer;
    transition: all 0.3s ease;
}

.tag-btn:hover {
    background: rgba(255, 255, 255, 0.12);
    transform: translateY(-2px);
}

.tag-btn.active {
    background: linear-gradient(45deg, #8b0000, #cc0000);
    border-color: transparent;
    box-shadow: 0 4px 15px rgba(139, 0, 0, 0.3);
}

.image-upload {
    margin-top: 12px;
}

.upload-area {
    border: 2px dashed rgba(255, 255, 255, 0.2);
    border-radius: 16px;
    padding: 40px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s ease;
}

.upload-area:hover {
    border-color: rgba(139, 0, 0, 0.5);
    background: rgba(255, 255, 255, 0.02);
}

.upload-area i {
    font-size: 2.5rem;
    color: rgba(255, 255, 255, 0.6);
    margin-bottom: 16px;
}

.upload-area p {
    color: white;
    font-size: 1.1rem;
    margin-bottom: 8px;
}

.upload-area span {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.9rem;
}

.image-preview {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 16px;
    margin-top: 20px;
}

.preview-item {
    position: relative;
    border-radius: 12px;
    overflow: hidden;
    aspect-ratio: 1;
}

.preview-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.remove-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
}

.remove-btn:hover {
    background: rgba(139, 0, 0, 0.8);
    transform: scale(1.1);
}

.cancel-btn,
.submit-btn {
    padding: 12px 24px;
    border-radius: 12px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
}

.cancel-btn {
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.cancel-btn:hover {
    background: rgba(255, 255, 255, 0.15);
}

.submit-btn {
    background: linear-gradient(45deg, #8b0000, #cc0000);
    color: white;
    border: none;
    display: flex;
    align-items: center;
    gap: 8px;
    box-shadow: 0 4px 15px rgba(139, 0, 0, 0.3);
}

.submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(139, 0, 0, 0.4);
}

.submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

@keyframes modalShow {
    from {
        opacity: 0;
        transform: scale(0.9);
    }

    to {
        opacity: 1;
        transform: scale(1);
    }
}

/* 响应式调整 */
@media (max-width: 768px) {
    .modal-content {
        border-radius: 16px;
        height: 100%;
    }

    .modal-header,
    .modal-body,
    .modal-footer {
        padding: 20px;
    }

    .modal-body::-webkit-scrollbar {
        width: 6px;
    }
}

.tags-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.hot-tags,
.custom-tags {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    padding: 16px;
}

.hot-tags h4,
.custom-tags h4 {
    color: rgba(255, 255, 255, 0.7);
    margin-bottom: 12px;
    font-size: 0.9rem;
}

.selected-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 12px;
}

.selected-tag {
    display: flex;
    align-items: center;
    background: rgba(139, 0, 0, 0.2);
    border: 1px solid rgba(139, 0, 0, 0.3);
    padding: 4px 12px;
    border-radius: 20px;
    gap: 8px;
}

.remove-tag {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.6);
    cursor: pointer;
    padding: 2px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.remove-tag:hover {
    color: white;
}

.tag-input-group {
    position: relative;
    display: flex;
    align-items: center;
}

.tag-input-group input {
    width: 100%;
    padding: 8px 40px 8px 12px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    color: white;
}

.char-count {
    position: absolute;
    right: 12px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.8rem;
}

.add-tag-btn {
    width: 100%;
    padding: 8px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px dashed rgba(255, 255, 255, 0.2);
    border-radius: 20px;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    transition: all 0.3s ease;
}

.add-tag-btn:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
    color: white;
}

/* Toast 提示样式 */
.toast-message {
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    padding: 12px 24px;
    border-radius: 8px;
    color: white;
    font-size: 0.95rem;
    display: flex;
    align-items: center;
    gap: 8px;
    z-index: 2000;
    animation: slideIn 0.3s ease, slideOut 0.3s ease 2.7s;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.toast-message.success {
    background: linear-gradient(45deg, #28a745, #34ce57);
}

.toast-message.error {
    background: linear-gradient(45deg, #dc3545, #ff4d5a);
}

.toast-message i {
    font-size: 1.1rem;
}

@keyframes slideIn {
    from {
        transform: translate(-50%, -100%);
        opacity: 0;
    }

    to {
        transform: translate(-50%, 0);
        opacity: 1;
    }
}

@keyframes slideOut {
    from {
        transform: translate(-50%, 0);
        opacity: 1;
    }

    to {
        transform: translate(-50%, -100%);
        opacity: 0;
    }
}

/* 移动端适配 */
@media (max-width: 768px) {
    .toast-message {
        width: 90%;
        text-align: center;
        justify-content: center;
    }
}
</style>