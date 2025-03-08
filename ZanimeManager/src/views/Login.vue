<template>
    <div
        class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
        <!-- 背景装饰元素 -->
        <div class="absolute top-0 left-0 w-full h-full overflow-hidden z-0">
            <div class="absolute top-10 left-10 w-20 h-20 rounded-full bg-indigo-100 opacity-40"></div>
            <div class="absolute top-20 right-20 w-32 h-32 rounded-full bg-indigo-200 opacity-30"></div>
            <div class="absolute bottom-10 left-1/4 w-40 h-40 rounded-full bg-indigo-100 opacity-20"></div>
            <div class="absolute -bottom-10 right-1/3 w-64 h-64 rounded-full bg-indigo-50 opacity-40"></div>

            <!-- 右上角动漫图标装饰 -->
            <div class="absolute top-24 right-10 text-indigo-200 opacity-30">
                <i class="fas fa-film text-8xl"></i>
            </div>

            <!-- 左下角动漫图标装饰 -->
            <div class="absolute bottom-20 left-16 text-indigo-200 opacity-30">
                <i class="fas fa-tv text-8xl"></i>
            </div>
        </div>

        <div class="max-w-sm w-full space-y-8 bg-white rounded-lg shadow-md p-8 relative z-10">
            <!-- Logo和标题 -->
            <div class="text-center">
                <div class="flex justify-center">
                    <i class="fas fa-play-circle text-4xl text-indigo-600"></i>
                </div>
                <div class="mt-1 bg-indigo-50 rounded-full mx-auto max-w-fit px-3 py-1">
                    <p class="text-xs text-indigo-700 font-medium">动漫内容管理系统</p>
                </div>
                <h2 class="mt-4 text-xl font-medium text-gray-700">欢迎回来</h2>
                <p class="mt-1 text-sm text-gray-500">请登录您的账户以继续访问</p>
            </div>

            <!-- 登录表单 -->
            <div class="space-y-4">
                <!-- 用户名输入 -->
                <div class="mb-4">
                    <div class="text-left mb-1 text-sm text-gray-600">用户名</div>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                            <i class="fas fa-user text-indigo-600"></i>
                        </div>
                        <input type="text" v-model="loginForm.username" placeholder="请输入用户名"
                            class="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 text-sm" />
                    </div>
                    <div v-if="usernameError" class="text-left text-red-500 text-xs mt-1">{{ usernameError }}</div>
                </div>

                <!-- 密码输入 -->
                <div class="mb-4">
                    <div class="text-left mb-1 text-sm text-gray-600">密码</div>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                            <i class="fas fa-lock text-indigo-600"></i>
                        </div>
                        <input :type="showPassword ? 'text' : 'password'" v-model="loginForm.password"
                            placeholder="请输入密码"
                            class="w-full pl-10 pr-10 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 text-sm" />
                        <div class="absolute inset-y-0 right-0 flex items-center pr-3">
                            <i :class="['fas cursor-pointer text-gray-500 hover:text-indigo-500 transition-colors', showPassword ? 'fa-eye-slash' : 'fa-eye']"
                                @click="showPassword = !showPassword">
                            </i>
                        </div>
                    </div>
                    <div v-if="passwordError" class="text-left text-red-500 text-xs mt-1">{{ passwordError }}</div>
                </div>

                <!-- 记住我和忘记密码 -->
                <div class="flex items-center justify-between pt-2">
                    <div class="flex items-center">
                        <input id="remember_me" type="checkbox"
                            class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                            v-model="loginForm.remember" />
                        <label for="remember_me" class="ml-2 block text-sm text-gray-700">
                            记住我
                        </label>
                    </div>

                    <div class="text-sm">
                        <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">
                            忘记密码?
                        </a>
                    </div>
                </div>

                <!-- 登录按钮 -->
                <div class="pt-4">
                    <button type="button" @click="handleLogin"
                        class="w-full py-2 px-4 rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors flex justify-center items-center"
                        :disabled="loading">
                        <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                            xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4">
                            </circle>
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                            </path>
                        </svg>
                        {{ loading ? '登录中...' : '登录' }}
                    </button>
                </div>
            </div>

            <!-- 页脚 -->
            <div class="mt-6 text-center text-sm text-gray-500">
                <p>&copy; {{ new Date().getFullYear() }} ZanimeManager. 版权所有。</p>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Login',
    data() {
        return {
            loading: false,
            showPassword: false,
            usernameError: '',
            passwordError: '',
            loginForm: {
                username: '',
                password: '',
                remember: false
            },
            wsUrl: 'ws://127.0.0.1:9999/ws' // WebSocket 服务器地址
        }
    },
    methods: {
        handleLogin() {
            // 验证表单
            this.usernameError = '';
            this.passwordError = '';

            if (!this.loginForm.username) {
                this.usernameError = '请输入用户名';
                return;
            }

            if (this.loginForm.username.length < 3 || this.loginForm.username.length > 20) {
                this.usernameError = '用户名长度应在 3 到 20 个字符之间';
                return;
            }

            if (!this.loginForm.password) {
                this.passwordError = '请输入密码';
                return;
            }

            if (this.loginForm.password.length < 6 || this.loginForm.password.length > 20) {
                this.passwordError = '密码长度应在 6 到 20 个字符之间';
                return;
            }

            this.loading = true;
            // 模拟登录API调用
            setTimeout(() => {
                // 保存登录状态
                localStorage.setItem('isLoggedIn', 'true');

                // 登录成功后初始化WebSocket连接
                this.$store.dispatch('initWebSocket', this.wsUrl)
                    .then(() => {
                        this.loading = false;
                        this.$router.push('/');
                        this.$message({
                            message: '登录成功',
                            type: 'success'
                        });
                    })
                    .catch(error => {
                        this.loading = false;
                        this.$message({
                            message: '连接失败：' + error.message,
                            type: 'error'
                        });
                    });
            }, 1500);
        }
    },
    unmounted() {
        // 组件销毁前关闭WebSocket连接
        this.$store.dispatch('closeWebSocket');
    }
}
</script>

<style scoped>
/* 确保所有输入框有一致的样式 */
input {
    height: 40px;
    box-sizing: border-box;
}

/* 确保图标垂直居中对齐 */
.absolute.inset-y-0 {
    display: flex;
    align-items: center;
}

/* 背景动画效果 */
@keyframes float {
    0% {
        transform: translateY(0px);
    }

    50% {
        transform: translateY(-10px);
    }

    100% {
        transform: translateY(0px);
    }
}

.absolute.top-10 {
    animation: float 7s ease-in-out infinite;
}

.absolute.top-20 {
    animation: float 9s ease-in-out infinite;
}

.absolute.bottom-10 {
    animation: float 8s ease-in-out infinite;
}

.absolute.-bottom-10 {
    animation: float 10s ease-in-out infinite;
}
</style>
