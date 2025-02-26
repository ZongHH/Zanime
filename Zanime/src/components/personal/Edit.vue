<template>
    <transition name="modal">
        <div class="edit-modal-overlay" v-if="show">
            <div class="edit-modal-content">
                <div class="modal-header">
                    <h2>编辑个人资料</h2>
                    <button class="close-btn" @click="$emit('close')">
                        <i class="fas fa-times"></i>
                    </button>
                </div>

                <div class="modal-body">
                    <form @submit.prevent="handleSubmit" class="edit-form">
                        <!-- 头像上传 -->
                        <div class="avatar-upload">
                            <div class="avatar-preview" @click="triggerFileInput">
                                <img :src="formData.avatar_url || userInfo.avatar_url" alt="头像">
                                <div class="upload-overlay">
                                    <i class="fas fa-camera"></i>
                                    <span>更换头像</span>
                                </div>
                            </div>
                            <input type="file" ref="fileInput" @change="handleAvatarChange" accept="image/*"
                                style="display: none">
                        </div>

                        <!-- 基本信息表单 -->
                        <div class="form-group">
                            <label>用户名</label>
                            <div class="text-field editable" @click="startEditUsername" v-if="!isEditingUsername">
                                <span>{{ formData.username }}</span>
                                <i class="fas fa-edit"></i>
                            </div>
                            <div class="text-field editing" v-else>
                                <input ref="usernameInput" type="text" v-model="formData.username"
                                    @blur="finishEditUsername" @keyup.enter="finishEditUsername"
                                    :class="{ 'error': errors.username }">
                            </div>
                            <span class="error-message" v-if="errors.username">{{ errors.username }}</span>
                        </div>

                        <div class="form-group">
                            <label>邮箱</label>
                            <div class="text-field">
                                <span>{{ formData.email }}</span>
                                <i class="fas fa-lock"></i>
                            </div>
                            <span class="field-hint">邮箱地址不可修改</span>
                        </div>

                        <div class="form-group">
                            <label>性别</label>
                            <div class="gender-options">
                                <label class="radio-label" :class="{ 'active': formData.gender === 'Male' }">
                                    <input type="radio" v-model="formData.gender" value="Male">
                                    <div class="radio-content">
                                        <i class="fas fa-mars"></i>
                                        <span>男</span>
                                    </div>
                                </label>
                                <label class="radio-label" :class="{ 'active': formData.gender === 'Female' }">
                                    <input type="radio" v-model="formData.gender" value="Female">
                                    <div class="radio-content">
                                        <i class="fas fa-venus"></i>
                                        <span>女</span>
                                    </div>
                                </label>
                                <label class="radio-label" :class="{ 'active': formData.gender === 'Other' }">
                                    <input type="radio" v-model="formData.gender" value="Other">
                                    <div class="radio-content">
                                        <i class="fas fa-genderless"></i>
                                        <span>其他</span>
                                    </div>
                                </label>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>个性签名</label>
                            <textarea v-model="formData.signature" placeholder="写点什么介绍一下自己吧..."
                                maxlength="200"></textarea>
                            <span class="char-count">{{ formData.signature.length }}/200</span>
                        </div>

                        <div class="form-actions">
                            <button type="button" class="cancel-btn" @click="$emit('close')">取消</button>
                            <button type="submit" class="submit-btn" :disabled="isSubmitting">
                                <i class="fas fa-spinner fa-spin" v-if="isSubmitting"></i>
                                {{ isSubmitting ? '保存中...' : '保存修改' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'EditProfileModal',
    props: {
        show: {
            type: Boolean,
            default: true
        }
    },
    data() {
        return {
            userInfo: {
                username: '',
                email: '',
                gender: '',
                signature: '',
                avatar_url: ''
            },
            formData: {
                username: '',
                email: '',
                gender: '',
                signature: '',
                avatar_url: ''
            },
            errors: {},
            isSubmitting: false,
            avatarFile: null,
            isEditingUsername: false
        }
    },
    watch: {
        userInfo: {
            // 立即执行一次handler函数
            immediate: true,
            /**
             * 监听userInfo变化的处理函数
             * @param {Object} newVal - 新的userInfo值
             * 当userInfo发生变化时,使用新值更新formData表单数据
             * 使用 || '' 确保在值为null或undefined时使用空字符串
             */
            handler(newVal) {
                this.formData = {
                    username: newVal.username || '',
                    email: newVal.email || '',
                    gender: newVal.gender || '',
                    signature: newVal.signature || '',
                    avatar_url: newVal.avatar_url || ''
                }
            }
        }
    },
    methods: {
        /**
         * 触发文件选择框的点击事件
         * 用于打开文件选择对话框以选择头像图片
         */
        triggerFileInput() {
            this.$refs.fileInput.click();
        },

        /**
         * 处理头像文件变更事件
         * @param {Event} event - 文件变更事件对象
         */
        async handleAvatarChange(event) {
            const file = event.target.files[0];
            if (!file) return;

            // 验证文件类型和大小
            if (!file.type.startsWith('image/')) {
                ElMessage.error('请上传图片文件');
                return;
            }

            if (file.size > 5 * 1024 * 1024) {
                ElMessage.error('图片大小不能超过5MB');
                return;
            }

            this.avatarFile = file;
            this.formData.avatar_url = URL.createObjectURL(file);
        },

        /**
         * 验证表单数据
         */
        validateForm() {
            this.errors = {};

            if (!this.formData.username.trim()) {
                this.errors.username = '用户名不能为空';
            }

            return Object.keys(this.errors).length === 0;
        },

        /**
         * 处理表单提交
         */
        async handleSubmit() {
            if (!this.validateForm()) return;

            try {
                this.isSubmitting = true;

                // 如果有新的头像文件，先上传头像
                if (this.avatarFile) {
                    const formData = new FormData();
                    formData.append('avatar', this.avatarFile);

                    const uploadResponse = await axios.post('/api/user/upload-avatar', formData, {
                        headers: {
                            'Content-Type': 'multipart/form-data'
                        }
                    });

                    if (uploadResponse.data.code === 200) {
                        this.formData.avatar_url = uploadResponse.data.url;
                    } else {
                        throw new Error('头像上传失败');
                    }
                }

                // 提交更新后的个人资料
                const response = await axios.post('/api/user/update', {
                    ...this.formData
                });

                if (response.data.code === 200) {
                    ElMessage.success('个人资料更新成功');
                    this.$emit('close');
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('更新失败：' + error.message);
            } finally {
                this.isSubmitting = false;
            }
        },

        /**
         * 开始编辑用户名
         * 激活用户名编辑状态并聚焦输入框
         */
        startEditUsername() {
            this.isEditingUsername = true;
            this.$nextTick(() => {
                this.$refs.usernameInput.focus();
            });
        },

        /**
         * 完成用户名编辑
         * 验证并保存用户名修改
         */
        finishEditUsername() {
            if (!this.formData.username.trim()) {
                this.errors.username = '用户名不能为空';
                return;
            }
            this.isEditingUsername = false;
            this.errors.username = '';
        },

        async loadUserInfo() {
            try {
                const response = await axios.get('/api/user/profile');
                if (response.data.code == 200) {
                    this.userInfo = response.data.profile;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.error('获取用户信息失败:', error);
            }
        }
    },
    mounted() {
        this.loadUserInfo();
    }
}
</script>

<style scoped>
.edit-modal-overlay {
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

.edit-modal-content {
    background: linear-gradient(145deg, rgba(30, 30, 30, 0.95), rgba(20, 20, 20, 0.95));
    border-radius: 20px;
    width: 90%;
    max-width: 500px;
    min-height: 550px;
    max-height: 85vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.05);
    animation: modal-slide-down 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
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
    padding-left: 12px;
    padding-right: 12px;
    padding-top: 8px;
    padding-bottom: 8px;
    border-radius: 50%;
    transition: all 0.3s ease;
}

.close-btn:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.modal-body {
    padding: 24px;
    overflow-y: auto;
    max-height: calc(100% - 70px);
    scrollbar-width: thin;
    scrollbar-color: rgba(252, 250, 250, 0.5) rgba(255, 255, 255, 0.1);
}

.edit-form {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 24px;
    align-items: center;
}

.avatar-upload {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;
}

.avatar-preview {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    overflow: hidden;
    position: relative;
    cursor: pointer;
    border: 3px solid rgba(255, 255, 255, 0.1);
    transition: all 0.3s ease;
}

.avatar-preview:hover .upload-overlay {
    opacity: 1;
}

.avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.upload-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
    gap: 8px;
}

.upload-overlay i {
    color: white;
    font-size: 1.5rem;
}

.upload-overlay span {
    font-size: 0.85rem;
    color: white;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    position: relative;
    width: 100%;
}

.form-group label {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.9rem;
}

.form-group input,
.form-group textarea {
    width: 100%;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 12px;
    color: white;
    font-size: 0.9rem;
    transition: all 0.3s ease;
}

.form-group input:focus,
.form-group textarea:focus {
    border-color: #cc0000;
    outline: none;
    background: rgba(255, 255, 255, 0.08);
}

.form-group input.error {
    border-color: #ff4444;
}

.error-message {
    color: #ff4444;
    font-size: 0.8rem;
}

.form-group textarea {
    min-height: 100px;
    padding-bottom: 32px;
    line-height: 1.5;
    resize: vertical;
}

.char-count {
    position: absolute;
    bottom: 8px;
    right: 12px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.85rem;
    background: transparent;
    padding: 0;
    pointer-events: none;
}

.gender-options {
    width: 100%;
    display: flex;
    justify-content: space-between;
    gap: 12px;
    overflow: hidden;
}

.radio-label {
    flex: 1;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 10px;
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
    min-width: 0;
}

.radio-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    width: 100%;
    min-width: 0;
}

.radio-content i {
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.7);
    transition: all 0.3s ease;
}

.radio-label input[type="radio"] {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
    margin: 0;
    padding: 0;
    pointer-events: none;
}

.radio-label span {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
    transition: all 0.3s ease;
}

.radio-label:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.2);
}

.radio-label:hover .radio-content i,
.radio-label:hover span {
    color: rgba(255, 255, 255, 0.9);
}

.radio-label input[type="radio"]:checked~.radio-content {
    color: #fff;
}

.radio-label input[type="radio"]:checked~.radio-content i,
.radio-label input[type="radio"]:checked~.radio-content span {
    color: #fff;
}

.radio-label input[type="radio"]:checked~.radio-bg {
    position: absolute;
    content: '';
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(145deg, #cc0000, #dd0000);
    opacity: 0.15;
    z-index: -1;
}

.form-actions {
    width: 100%;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 16px;
}

.cancel-btn,
.submit-btn {
    min-width: 100px;
}

.cancel-btn {
    padding: 10px 24px;
    background: rgba(255, 255, 255, 0.1);
    border: none;
    border-radius: 8px;
    color: white;
    font-size: 0.95rem;
    cursor: pointer;
    transition: all 0.3s ease;
}

.cancel-btn:hover {
    background: rgba(255, 255, 255, 0.15);
}

.submit-btn {
    min-width: 120px;
    padding: 10px 24px;
    background: linear-gradient(145deg, #cc0000, #dd0000);
    border: none;
    border-radius: 8px;
    color: white;
    font-size: 0.95rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.submit-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(204, 0, 0, 0.3);
}

.submit-btn:disabled {
    background: #666;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
}

/* 文本字段样式 */
.text-field {
    width: 100%;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 12px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.95rem;
    min-height: 42px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.text-field span {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

/* 邮箱字段特殊样式 */
.text-field i {
    color: rgba(255, 255, 255, 0.3);
    font-size: 0.9rem;
    margin-left: 12px;
}

.field-hint {
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.8rem;
    margin-top: 4px;
    text-align: right;
}

/* 性别选择优化 */
.radio-label.active {
    background: rgba(204, 0, 0, 0.1);
    border-color: rgba(204, 0, 0, 0.3);
}

.radio-label.active .radio-content i,
.radio-label.active .radio-content span {
    color: #cc0000;
}

/* 动画优化 */
.edit-modal-content {
    animation: modal-slide-down 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes modal-slide-down {
    from {
        transform: translateY(-30px);
        opacity: 0;
    }

    to {
        transform: translateY(0);
        opacity: 1;
    }
}

/* 响应式调整 */
@media (max-width: 768px) {
    .edit-modal-content {
        width: 95%;
        max-width: 450px;
    }
}

@media (max-width: 480px) {
    .edit-modal-content {
        width: 100%;
        height: 100vh;
        max-height: none;
        border-radius: 0;
    }

    .modal-body {
        padding: 20px;
    }

    .gender-options {
        gap: 8px;
    }
}

/* Webkit 滚动条样式 */
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

.text-field.editable {
    cursor: pointer;
    transition: all 0.3s ease;
}

.text-field.editable:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.2);
}

.text-field.editable i {
    color: rgba(255, 255, 255, 0.3);
    font-size: 0.9rem;
    margin-left: 12px;
    opacity: 0;
    transition: all 0.3s ease;
}

.text-field.editable:hover i {
    opacity: 1;
}

.text-field.editing {
    padding: 0;
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.2);
}

.text-field.editing input {
    width: 100%;
    background: transparent;
    border: none;
    padding: 12px;
    color: white;
    font-size: 0.95rem;
    outline: none;
}

.text-field.editing input.error {
    color: #ff4444;
}

.error-message {
    color: #ff4444;
    font-size: 0.8rem;
    margin-top: 4px;
}
</style>
