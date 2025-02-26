<template>
    <div class="login-container">
        <div class="login-box">
            <div class="login-header">
                <h2>Zanime Monitor</h2>
                <p>监控系统登录</p>
            </div>
            <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
                <el-form-item prop="username">
                    <el-input v-model="loginForm.username" placeholder="用户名">
                        <template #prefix>
                            <el-icon>
                                <User />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input v-model="loginForm.password" type="password" placeholder="密码" show-password>
                        <template #prefix>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <div class="login-options">
                    <el-checkbox v-model="loginForm.remember">记住密码</el-checkbox>
                    <el-link type="primary">忘记密码？</el-link>
                </div>
                <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin">
                    {{ loading ? '登录中...' : '登录' }}
                </el-button>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Login',
    data() {
        return {
            loading: false,
            loginForm: {
                username: '',
                password: '',
                remember: false
            },
            rules: {
                username: [
                    { required: true, message: '请输入用户名', trigger: 'blur' },
                    { min: 3, max: 20, message: '用户名长度应在 3 到 20 个字符之间', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: '请输入密码', trigger: 'blur' },
                    { min: 6, max: 20, message: '密码长度应在 6 到 20 个字符之间', trigger: 'blur' }
                ]
            },
            wsUrl: 'ws://127.0.0.1:9999/ws' // WebSocket 服务器地址
        }
    },
    methods: {
        handleLogin() {
            this.$refs.loginFormRef.validate(valid => {
                if (valid) {
                    this.loading = true
                    // 模拟登录API调用
                    setTimeout(() => {
                        // 登录成功后初始化WebSocket连接
                        this.$store.dispatch('initWebSocket', this.wsUrl)
                            .then(() => {
                                this.loading = false
                                this.$router.push('/')
                                this.$message.success('登录成功')
                            })
                            .catch(error => {
                                this.loading = false
                                this.$message.error('连接失败：' + error.message)
                            })
                    }, 1000)
                }
            })
        }
    },
    beforeDestroy() {
        // 组件销毁前关闭WebSocket连接
        this.$store.dispatch('closeWebSocket')
    }
}
</script>

<style scoped>
.login-container {
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: #f5f7fa;
    position: relative;
}

.login-box {
    width: 100%;
    max-width: 400px;
    padding: 40px;
    background: #ffffff;
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

.login-header {
    text-align: center;
    margin-bottom: 36px;
}

.login-header h2 {
    font-size: 28px;
    color: #2c3e50;
    margin-bottom: 8px;
    font-weight: 500;
}

.login-header p {
    color: #606266;
    font-size: 16px;
}

.login-form {
    margin-top: 24px;
}

.login-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 24px 0;
}

.login-button {
    width: 100%;
    padding: 12px 0;
    font-size: 16px;
    border-radius: 4px;
}

.login-button:hover {
    opacity: 0.9;
}

:deep(.el-input__wrapper) {
    box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__wrapper:hover) {
    box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input__wrapper.is-focus) {
    box-shadow: 0 0 0 1px #409eff inset;
}

/* 移动端适配 */
@media screen and (max-width: 480px) {
    .login-box {
        margin: 20px;
        padding: 30px 20px;
    }

    .login-header h2 {
        font-size: 24px;
    }

    .login-header p {
        font-size: 14px;
    }
}
</style>
