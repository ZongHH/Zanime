<template>
  <section class="login">
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
              <a href="forgot" class="color-primary">忘记密码？</a>
            </div>

            <!-- Facebook 登录按钮 -->
            <div class="login-cards mb-32">
              <div class="text-center">
                <input type="submit" value="登录"
                  style="background-color: transparent; color: white; border: none; cursor: pointer; padding: 10px 20px;" />
              </div>
            </div>

            <!-- Google 登录按钮 -->
            <a href="#" class="link-btn mb-24">
              <img src="@/static/picture/google.png" alt="logo" />
              或使用谷歌账号登录
            </a>

            <!-- 注册链接 -->
            <a href="signup" class="white d-flex justify-content-center gap-8">
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
export default {
  name: "LoginForm",
  data() {
    return {
      email: "",
      password: "",
      userIp: "", // 添加用户IP字段
    };
  },
  methods: {
    // 获取用户IP的方法
    async getUserIp() {
      try {
        const response = await axios.get('https://api.ipify.org?format=json');
        this.userIp = response.data.ip;
      } catch (error) {
        console.error('Failed to get IP:', error);
        this.userIp = '0.0.0.0'; // 获取失败时的默认值
      }
    },

    async handleSubmit() {
      const formData = new FormData();
      formData.append('email', this.email);
      formData.append('password', this.password);
      formData.append('user_ip', this.userIp);

      try {
        const response = await axios.post('/api/loginInfo', formData);

        if (response.data.code === 200) {
          // 清空表单数据
          this.email = '';
          this.password = '';

          // 存储用户信息到 localStorage
          const userData = response.data.user_info;
          localStorage.setItem('user_name', userData.user_name);
          localStorage.setItem('user_id', userData.user_id);
          localStorage.setItem('email', userData.email);
          localStorage.setItem('gender', userData.gender);
          localStorage.setItem('avatar_url', userData.avatar_url);

          // 获取重定向地址（如果有的话）
          const redirectPath = this.$route.query.redirect || '/';

          // 跳转到重定向地址或首页
          this.$router.push(redirectPath);
        } else {
          showError(response.data.message);
        }
      } catch (error) {
        console.error(error);
        showError(error.response.data.message || '登录时发生错误。');
      }
    },
    // 映射 Vuex 中的 actions 到组件的方法
    ...mapActions('websocket', ['connectWebSocket', 'sendMessage', 'disconnectWebSocket']),
  },
  async mounted() {
    await this.getUserIp(); // 组件加载时获取IP
  },
  computed: {
    // 映射 Vuex 中的 state 到组件的 computed 属性
    ...mapState('websocket', ['socket', 'isConnected', 'user']),
  }
};
import axios from "axios";
import showError from "@/static/js/showError.js";
import { mapState, mapActions } from "vuex";
</script>

<style scoped></style>
