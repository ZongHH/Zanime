<template>
  <section class="login">
    <!-- 体验账号提示框 -->
    <div v-if="showTestAccount" class="test-account-tip">
      <div class="tip-content">
        <div class="tip-header">
          <h5>体验账号(1小时)</h5>
          <button class="close-btn" @click="showTestAccount = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="tip-body">
          <div class="info-item">
            <p>账号：{{ testAccount.email }}</p>
          </div>
          <div class="info-item">
            <p>密码：{{ testAccount.password }}</p>
          </div>
          <div class="tip-footer" @click="fillTestAccount">
            <button class="fill-button">
              <i class="fas fa-sign-in-alt"></i>
              <span>一键填充账号密码</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="row justify-content-center">
      <!-- 图片部分 -->
      <div class="col-xxl-6 pe-xxl-0">
        <div class="pic">
          <img src="@/static/picture/banner1.jpg" alt="pic" />
        </div>
      </div>

      <!-- 登录表单部分 -->
      <div class="col-xxl-4 ps-xxl-0">
        <div class="bg-ter login-form">
          <h4 class="white mb-48">登录您的账户！</h4>
          <div class="d-flex gap-12 mb-48">
            <img src="@/static/picture/logo1.png" alt="pic" />
            <h6 class="white">Logo</h6>
          </div>
          <h6 class="white mb-24">很高兴再次见到您</h6>

          <form @submit.prevent="handleSubmit">
            <!-- Email/Phone Number 输入框 -->
            <div class="input-box-1">
              <label class="white mb-08">账号</label>
              <input v-model="email" type="text" name="email" id="eMail" class="form-input mb-16" placeholder="邮箱或手机号"
                required />
            </div>

            <!-- Password 输入框 -->
            <div class="input-block input-box-1 mb-20">
              <label class="white mb-08">密码</label>
              <input v-model="password" type="password" class="form-control password-input" id="passWord"
                name="password" placeholder="请输入密码" required />
              <i class="fas fa-eye-slash" id="eye"></i>
            </div>

            <!-- 忘记密码链接 -->
            <div class="text-end mb-32">
              <a href="#" class="color-primary" @click.prevent="checkServiceStatus">忘记密码？</a>
            </div>

            <!-- 登录按钮 -->
            <div class="login-cards mb-32">
              <div class="text-center">
                <input type="submit" value="登录"
                  style="background-color: transparent; color: white; border: none; cursor: pointer; padding: 10px 20px;" />
              </div>
            </div>

            <!-- 体验账号按钮 -->
            <button type="button" class="link-btn mb-24" @click.prevent="getTestAccount">
              <i class="fas fa-user-circle"></i>
              获取体验账号
            </button>

            <!-- 注册链接 -->
            <a href="#" class="white d-flex justify-content-center gap-8" @click.prevent="checkServiceStatus">
              还没有账号？
              <span class="color-primary">立即注册</span>
            </a>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import axios from "axios";
import { ElMessage } from "element-plus";

export default {
  name: "LoginForm",
  data() {
    return {
      email: "",
      password: "",
      showTestAccount: false,
      testAccount: {
        email: '',
        password: ''
      }
    };
  },
  methods: {
    // 获取体验账号信息
    async getTestAccount() {
      try {
        const response = await axios.get('/api/user/test-account');
        if (response.data.code === 200) {
          this.testAccount = response.data.account;
          this.showTestAccount = true;
        } else {
          throw new Error(response.data.message);
        }
      } catch (error) {
        ElMessage.warning('获取体验账号信息失败:' + error.message);
      }
    },

    async handleSubmit() {
      try {
        const response = await axios.post('/api/loginInfo', {
          email: this.email,
          password: this.password
        });

        if (response.data.code === 200) {
          // 清空表单数据
          this.email = '';
          this.password = '';

          // 存储用户信息到 localStorage
          const userData = response.data.user_info;
          localStorage.setItem('user_name', userData.username);
          localStorage.setItem('user_id', userData.user_id);
          localStorage.setItem('email', userData.email);
          localStorage.setItem('gender', userData.gender);
          localStorage.setItem('avatar_url', userData.avatar_url);

          // 获取重定向地址（如果有的话）
          const redirectPath = this.$route.query.redirect || '/';

          // 跳转到重定向地址或首页
          this.$router.push(redirectPath);
        } else {
          throw new Error(response.data.message);
        }
      } catch (error) {
        ElMessage.error('登录失败:' + error.message);
      }
    },

    // 提示服务暂未开放
    async checkServiceStatus() {
      ElMessage.warning('服务暂未开放');
    },

    // 一键填充测试账号
    fillTestAccount() {
      this.email = this.testAccount.email;
      this.password = this.testAccount.password;
      ElMessage.success('已自动填充账号密码');
    },
  },
};
</script>

<style scoped>
.test-account-tip {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  min-width: 320px;
  background: rgba(28, 28, 35, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
  animation: slideIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.tip-content {
  padding: 20px;
  color: white;
  position: relative;
}

.tip-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.2), transparent);
}

.tip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.tip-header h5 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  gap: 8px;
  letter-spacing: 0.5px;
}

.close-btn {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  padding: 8px;
  transition: all 0.2s ease;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  margin: -4px;
}

.close-btn:hover {
  color: white;
  background: rgba(255, 255, 255, 0.1);
  transform: rotate(90deg);
}

.tip-body {
  padding-top: 4px;
}

.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.05);
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  margin: 12px 0;
  transition: all 0.2s ease;
}

.info-item i {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
}

.info-item:hover i {
  color: rgba(255, 255, 255, 0.9);
}

.tip-footer {
  text-align: center;
  margin-top: 16px;
  padding: 0 8px;
}

.fill-button {
  width: 100%;
  background: linear-gradient(135deg, #ffffff 0%, #3f4a55 100%);
  border: none;
  border-radius: 12px;
  padding: 12px 20px;
  color: white;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  box-shadow: 0 4px 15px rgba(74, 144, 226, 0.2);
}

.fill-button i {
  font-size: 1rem;
  transition: transform 0.3s ease;
}

.fill-button:hover {
  background: linear-gradient(135deg, #ffffff 0%, #3f4a55 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(74, 144, 226, 0.3);
}

.fill-button:hover i {
  transform: translateX(3px);
}

.fill-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(74, 144, 226, 0.2);
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }

  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .test-account-tip {
    top: auto;
    bottom: 20px;
    left: 20px;
    right: 20px;
    min-width: auto;
  }

  .info-item {
    padding: 12px 14px;
  }

  .info-item p {
    font-size: 0.9rem;
  }
}
</style>
