<template>
    <div class="feedback-center">
        <div class="feedback-container">
            <div class="feedback-header">
                <div class="header-content">
                    <h2 class="feedback-title">
                        <i class="fas fa-comments"></i> 反馈中心
                    </h2>
                    <div class="feedback-tabs">
                        <button :class="['tab-btn', { active: activeTab === 'suggestion' }]"
                            @click="activeTab = 'suggestion'">
                            <i class="fas fa-lightbulb"></i> 功能建议
                        </button>
                        <button :class="['tab-btn', { active: activeTab === 'bug' }]" @click="activeTab = 'bug'">
                            <i class="fas fa-bug"></i> 问题反馈
                        </button>
                    </div>
                </div>
                <div class="header-gradient"></div>
            </div>

            <div class="feedback-content">
                <!-- 功能建议表单 -->
                <div v-if="activeTab === 'suggestion'" class="feedback-form">
                    <div class="form-intro">
                        <p>您的建议对我们至关重要！请分享您对Zanime的功能改进或新功能想法。</p>
                    </div>

                    <div class="form-group">
                        <label for="suggestion-title">建议标题</label>
                        <input type="text" id="suggestion-title" v-model="suggestionForm.title"
                            :class="['form-input', { 'is-invalid': errors.suggestion.title }]" placeholder="简短描述您的建议">
                        <div class="error-message" v-if="errors.suggestion.title">
                            {{ errors.suggestion.title }}
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="suggestion-type">建议类型</label>
                        <select id="suggestion-type" v-model="suggestionForm.type" class="form-input">
                            <option value="feature">新功能请求</option>
                            <option value="improvement">功能改进</option>
                            <option value="design">界面设计</option>
                            <option value="other">其他建议</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="suggestion-description">详细描述</label>
                        <textarea id="suggestion-description" v-model="suggestionForm.description"
                            :class="['form-input textarea', { 'is-invalid': errors.suggestion.description }]"
                            placeholder="请详细描述您的建议，包括为什么这个建议对您很重要" rows="5"></textarea>
                        <div class="error-message" v-if="errors.suggestion.description">
                            {{ errors.suggestion.description }}
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="suggestion-contact">联系方式 (选填)</label>
                        <input type="text" id="suggestion-contact" v-model="suggestionForm.contact" class="form-input"
                            placeholder="留下您的邮箱或其他联系方式，以便我们进一步沟通">
                    </div>

                    <div class="form-actions">
                        <button class="action-btn submit" @click="submitSuggestion" :disabled="isSubmitting">
                            <i class="fas fa-paper-plane"></i>
                            {{ isSubmitting ? '提交中...' : '提交建议' }}
                        </button>
                    </div>
                </div>

                <!-- 问题反馈表单 -->
                <div v-if="activeTab === 'bug'" class="feedback-form">
                    <div class="form-intro">
                        <p>发现了问题或错误？请告诉我们详细情况，帮助我们改进Zanime。</p>
                    </div>

                    <div class="form-group">
                        <label for="bug-title">问题标题</label>
                        <input type="text" id="bug-title" v-model="bugForm.title"
                            :class="['form-input', { 'is-invalid': errors.bug.title }]" placeholder="简短描述您遇到的问题">
                        <div class="error-message" v-if="errors.bug.title">
                            {{ errors.bug.title }}
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="bug-severity">问题严重程度</label>
                        <select id="bug-severity" v-model="bugForm.severity" class="form-input">
                            <option value="critical">严重 - 无法使用核心功能</option>
                            <option value="major">主要 - 功能受到显著影响</option>
                            <option value="minor">次要 - 有问题但可以正常使用</option>
                            <option value="cosmetic">界面 - 视觉或排版问题</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="bug-description">问题描述</label>
                        <textarea id="bug-description" v-model="bugForm.description"
                            :class="['form-input textarea', { 'is-invalid': errors.bug.description }]"
                            placeholder="请详细描述问题：1. 您做了什么操作 2. 期望看到什么 3. 实际发生了什么" rows="5"></textarea>
                        <div class="error-message" v-if="errors.bug.description">
                            {{ errors.bug.description }}
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="bug-reproduction">重现步骤</label>
                        <textarea id="bug-reproduction" v-model="bugForm.reproduction"
                            :class="['form-input textarea', { 'is-invalid': errors.bug.reproduction }]"
                            placeholder="请列出详细的步骤，帮助我们重现这个问题" rows="3"></textarea>
                        <div class="error-message" v-if="errors.bug.reproduction">
                            {{ errors.bug.reproduction }}
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="bug-browser">浏览器/设备信息</label>
                        <input type="text" id="bug-browser" v-model="bugForm.environment" class="form-input"
                            placeholder="例如：Chrome 96.0.4664.110, Windows 10">
                    </div>

                    <div class="form-group">
                        <label for="bug-contact">联系方式 (选填)</label>
                        <input type="text" id="bug-contact" v-model="bugForm.contact" class="form-input"
                            placeholder="留下您的邮箱或其他联系方式，以便我们进一步沟通">
                    </div>

                    <div class="form-actions">
                        <button class="action-btn submit" @click="submitBug" :disabled="isSubmitting">
                            <i class="fas fa-bug"></i>
                            {{ isSubmitting ? '提交中...' : '提交问题' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 成功提示 -->
        <div class="feedback-toast" v-if="showSuccessToast">
            <div class="toast-content">
                <i class="fas fa-check-circle"></i>
                <span>{{ successMessage }}</span>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'FeedbackCenter',
    data() {
        return {
            activeTab: 'suggestion', // 当前激活的标签页
            isSubmitting: false, // 是否正在提交
            showSuccessToast: false, // 是否显示成功提示
            successMessage: '', // 成功提示信息

            // 功能建议表单
            suggestionForm: {
                title: '',
                type: 'feature',
                description: '',
                contact: ''
            },

            // 问题反馈表单
            bugForm: {
                title: '',
                severity: 'major',
                description: '',
                reproduction: '',
                environment: this.getBrowserInfo(),
                contact: ''
            },

            // 表单错误信息
            errors: {
                suggestion: {
                    title: '',
                    description: ''
                },
                bug: {
                    title: '',
                    description: '',
                    reproduction: ''
                }
            }
        };
    },
    methods: {
        // 功能暂未开放
        handleServiceUnavailable() {
            ElMessage.warning('功能暂未开放');
        },

        // 获取浏览器信息
        getBrowserInfo() {
            const userAgent = navigator.userAgent;
            let browserName = "未知浏览器";

            if (userAgent.match(/chrome|chromium|crios/i)) {
                browserName = "Chrome";
            } else if (userAgent.match(/firefox|fxios/i)) {
                browserName = "Firefox";
            } else if (userAgent.match(/safari/i)) {
                browserName = "Safari";
            } else if (userAgent.match(/opr\//i)) {
                browserName = "Opera";
            } else if (userAgent.match(/edg/i)) {
                browserName = "Edge";
            } else if (userAgent.match(/msie|trident/i)) {
                browserName = "Internet Explorer";
            }

            return `${browserName}, ${navigator.platform}`;
        },

        // 验证功能建议表单
        validateSuggestionForm() {
            let isValid = true;
            this.errors.suggestion.title = '';
            this.errors.suggestion.description = '';

            if (!this.suggestionForm.title.trim()) {
                this.errors.suggestion.title = '请输入建议标题';
                isValid = false;
            }

            if (!this.suggestionForm.description.trim()) {
                this.errors.suggestion.description = '请输入详细描述';
                isValid = false;
            } else if (this.suggestionForm.description.trim().length < 10) {
                this.errors.suggestion.description = '描述内容过短，请详细说明您的建议，10字以上';
                isValid = false;
            }

            return isValid;
        },

        // 验证问题反馈表单
        validateBugForm() {
            let isValid = true;
            this.errors.bug.title = '';
            this.errors.bug.description = '';
            this.errors.bug.reproduction = '';

            if (!this.bugForm.title.trim()) {
                this.errors.bug.title = '请输入问题标题';
                isValid = false;
            }

            if (!this.bugForm.description.trim()) {
                this.errors.bug.description = '请输入问题描述';
                isValid = false;
            } else if (this.bugForm.description.trim().length < 10) {
                this.errors.bug.description = '描述内容过短，请详细说明问题，10字以上';
                isValid = false;
            }

            if (!this.bugForm.reproduction.trim()) {
                this.errors.bug.reproduction = '请输入重现步骤';
                isValid = false;
            }

            return isValid;
        },

        // 提交功能建议
        async submitSuggestion() {
            if (!this.validateSuggestionForm()) return;

            this.isSubmitting = true;

            try {
                // 这里是模拟API请求，实际项目中应替换为真实的API请求
                // const response = await axios.post('/api/feedback/suggestion', this.suggestionForm);

                // 模拟API延迟
                await new Promise(resolve => setTimeout(resolve, 1000));

                // 显示成功提示
                this.successMessage = '感谢您的建议！我们会认真考虑您的反馈。';
                this.showSuccessToast = true;
                setTimeout(() => {
                    this.showSuccessToast = false;
                }, 3000);

                // 清空表单
                this.clearSuggestionForm();
            } catch (error) {
                console.error('提交建议失败:', error);
                // 这里可以添加错误处理逻辑
            } finally {
                this.isSubmitting = false;
            }
        },

        // 提交问题反馈
        async submitBug() {
            if (!this.validateBugForm()) return;

            this.isSubmitting = true;

            try {
                // 这里是模拟API请求，实际项目中应替换为真实的API请求
                // const response = await axios.post('/api/feedback/bug', this.bugForm);

                // 模拟API延迟
                await new Promise(resolve => setTimeout(resolve, 1000));

                // 显示成功提示
                this.successMessage = '感谢您的反馈！我们会尽快修复这个问题。';
                this.showSuccessToast = true;
                setTimeout(() => {
                    this.showSuccessToast = false;
                }, 3000);

                // 清空表单
                this.clearBugForm();
            } catch (error) {
                console.error('提交问题反馈失败:', error);
                // 这里可以添加错误处理逻辑
            } finally {
                this.isSubmitting = false;
            }
        },

        // 清空功能建议表单
        clearSuggestionForm() {
            this.suggestionForm = {
                title: '',
                type: 'feature',
                description: '',
                contact: ''
            };
            this.errors.suggestion = {
                title: '',
                description: ''
            };
        },

        // 清空问题反馈表单
        clearBugForm() {
            this.bugForm = {
                title: '',
                severity: 'major',
                description: '',
                reproduction: '',
                environment: this.getBrowserInfo(),
                contact: ''
            };
            this.errors.bug = {
                title: '',
                description: '',
                reproduction: ''
            };
        }
    },

    watch: {
        '$route.query': {
            handler(newVal) {
                if (newVal.type === 'bug') {
                    this.activeTab = 'bug';
                } else {
                    this.activeTab = 'suggestion';
                }
            },
            immediate: true
        }
    },

    mounted() {
        this.handleServiceUnavailable();
    }
};
</script>

<style scoped>
.feedback-center {
    width: 100%;
    max-width: 100%;
    margin: 0;
    padding: 20px;
    color: #f0f0f0;
    background-color: #101010;
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: flex-start;
}

.feedback-container {
    width: 100%;
    max-width: 900px;
    background: rgba(25, 25, 28, 0.95);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
}

.feedback-header {
    position: relative;
    margin-bottom: 20px;
}

.header-content {
    display: flex;
    flex-direction: column;
    padding: 20px 20px 10px;
    position: relative;
    z-index: 2;
}

.header-gradient {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 100%;
    background: linear-gradient(135deg, #7e040e, #590511);
    z-index: 1;
    opacity: 0.75;
}

.feedback-title {
    margin: 0 0 20px 0;
    font-size: 1.8rem;
    font-weight: 600;
    color: white;
}

.feedback-title i {
    margin-right: 10px;
}

.feedback-tabs {
    display: flex;
    gap: 15px;
}

.tab-btn {
    background: rgba(255, 255, 255, 0.08);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 8px 8px 0 0;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s ease;
}

.tab-btn i {
    margin-right: 8px;
}

.tab-btn.active {
    background: rgba(35, 35, 40, 0.9);
    box-shadow: 0 -4px 10px rgba(0, 0, 0, 0.1);
}

.tab-btn:hover:not(.active) {
    background: rgba(255, 255, 255, 0.15);
}

.feedback-content {
    padding: 30px;
}

.form-intro {
    margin-bottom: 20px;
    color: #ccc;
    font-size: 0.95rem;
    line-height: 1.5;
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #f0f0f0;
}

.form-input {
    width: 100%;
    padding: 12px 15px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    background: rgba(30, 30, 30, 0.7);
    color: #f0f0f0;
    font-size: 0.95rem;
    transition: all 0.2s ease;
}

.form-input::placeholder {
    color: #999;
}

.form-input:focus {
    outline: none;
    border-color: #7e040e;
    box-shadow: 0 0 0 2px rgba(126, 4, 14, 0.2);
}

.form-input.is-invalid {
    border-color: #ee5555;
}

.textarea {
    resize: vertical;
    min-height: 100px;
}

.error-message {
    margin-top: 5px;
    color: #ee5555;
    font-size: 0.85rem;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 15px;
    margin-top: 30px;
}

.action-btn {
    padding: 10px 20px;
    border: none;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 8px;
}

.action-btn.clear {
    background: rgba(45, 45, 50, 0.7);
    color: #f0f0f0;
}

.action-btn.submit {
    background: #7e040e;
    color: white;
}

.action-btn.clear:hover {
    background: rgba(60, 60, 65, 0.8);
}

.action-btn.submit:hover {
    background: #960512;
}

.action-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.feedback-toast {
    position: fixed;
    bottom: 30px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(40, 150, 80, 0.9);
    color: white;
    padding: 15px 25px;
    border-radius: 8px;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
    z-index: 1000;
    animation: toast-in 0.3s ease, toast-out 0.3s ease 2.7s forwards;
}

.toast-content {
    display: flex;
    align-items: center;
    gap: 12px;
}

.toast-content i {
    font-size: 1.3rem;
}

@keyframes toast-in {
    from {
        transform: translate(-50%, 20px);
        opacity: 0;
    }

    to {
        transform: translate(-50%, 0);
        opacity: 1;
    }
}

@keyframes toast-out {
    from {
        transform: translate(-50%, 0);
        opacity: 1;
    }

    to {
        transform: translate(-50%, -20px);
        opacity: 0;
    }
}

/* 响应式设计 */
@media (max-width: 768px) {
    .feedback-center {
        padding: 10px;
    }

    .feedback-container {
        border-radius: 8px;
    }

    .feedback-header {
        margin-bottom: 0;
    }

    .header-content {
        padding: 15px 15px 5px;
    }

    .feedback-title {
        font-size: 1.5rem;
        margin-bottom: 15px;
    }

    .feedback-tabs {
        gap: 8px;
    }

    .tab-btn {
        padding: 8px 12px;
        font-size: 0.9rem;
    }

    .feedback-content {
        padding: 20px 15px;
    }

    .form-actions {
        flex-direction: column;
        gap: 10px;
    }

    .action-btn {
        width: 100%;
        justify-content: center;
    }
}
</style>
