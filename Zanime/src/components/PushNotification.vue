<template>
  <transition name="slide-fade">
    <div v-if="visible" class="push-notification">
      <div class="push-header">
        <span class="push-title">{{ title }}</span>
        <button class="close-btn" @click="closeNotification">×</button>
      </div>
      <div class="push-body">
        <slot>{{ message }}</slot>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: 'Notification',
    },
    message: {
      type: String,
      default: 'This is a push notification.',
    },
    duration: {
      type: Number,
      default: 5000, // 自动关闭时间 (毫秒)
    },
  },
  data() {
    return {
      visible: false,
    };
  },
  methods: {
    showNotification() {
      this.visible = true;
      if (this.duration > 0) {
        setTimeout(() => {
          this.closeNotification();
        }, this.duration);
      }
    },
    closeNotification() {
      this.visible = false;
      this.$emit('closed');
    },
  },
  mounted() {
    this.showNotification();
  },
};
</script>
<style scoped>
/* 推送窗口样式 */
.push-notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 320px;
  background-color: #1e1e2f;
  /* 暗色背景 */
  border: 1px solid #2d2d3d;
  /* 边框颜色 */
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  /* 深色阴影 */
  border-radius: 10px;
  overflow: hidden;
  z-index: 1000;
  animation: slide-in 0.3s ease-out;
  font-family: 'Poppins', Arial, sans-serif;
  /* 动漫风格字体 */
  color: #f4f4f9;
  /* 浅色字体 */
}

.push-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, #b71c1c, #ff5722);
  /* 深红到浅红渐变 */
  color: #ffffff;
  padding: 12px 16px;
  font-size: 16px;
}

.push-title {
  font-weight: bold;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  color: #f4f4f9;
  font-size: 20px;
  cursor: pointer;
  transition: color 0.3s;
}

.close-btn:hover {
  color: #ff6f91;
  /* 按钮悬停时的亮色 */
}

.push-body {
  padding: 15px;
  color: #d1c4e9;
  /* 较浅的文字颜色 */
  font-size: 14px;
  line-height: 1.6;
  background-color: #2d2d3d;
  /* 身体背景 */
  border-top: 1px solid #3a3a4a;
  /* 分隔线 */
}

/* 过渡动画 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.4s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

@keyframes slide-in {
  from {
    opacity: 0;
    transform: translateY(50px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>