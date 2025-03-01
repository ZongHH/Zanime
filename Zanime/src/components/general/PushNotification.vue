<template>
  <transition name="slide-fade">
    <div v-if="visible" class="push-notification" :class="type">
      <div class="push-header">
        <div class="icon-title">
          <div class="notification-icon">
            <slot name="icon">
              <svg v-if="type === 'like'" viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z" />
              </svg>
              <svg v-else-if="type === 'comment'" viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M20,2H4C2.9,2,2,2.9,2,4v18l4-4h14c1.1,0,2-0.9,2-2V4C22,2.9,21.1,2,20,2z M13,14h-2v-2h2V14z M13,10h-2V6h2V10z" />
              </svg>
              <svg v-else-if="type === 'reply'" viewBox="0 0 24 24" width="24" height="24">
                <path d="M10,9V5L3,12L10,19V14.9C15,14.9 18.5,16.5 21,20C20,15 17,10 10,9Z" />
              </svg>
              <svg v-else-if="type === 'favorite'" viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M17,3H7C4.791,3,3,4.791,3,7v10c0,2.209,1.791,4,4,4h10c2.209,0,4-1.791,4-4V7C21,4.791,19.209,3,17,3z M16,13h-3v3c0,0.552-0.448,1-1,1s-1-0.448-1-1v-3H8c-0.552,0-1-0.448-1-1s0.448-1,1-1h3V8c0-0.552,0.448-1,1-1s1,0.448,1,1v3h3c0.552,0,1,0.448,1,1S16.552,13,16,13z" />
              </svg>
              <svg v-else-if="type === 'follow'" viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M16,9c-1.66,0-3-1.34-3-3s1.34-3,3-3s3,1.34,3,3S17.66,9,16,9z M16,5c-0.55,0-1,0.45-1,1s0.45,1,1,1s1-0.45,1-1S16.55,5,16,5z M8,9C6.34,9,5,7.66,5,6s1.34-3,3-3s3,1.34,3,3S9.66,9,8,9z M8,5C7.45,5,7,5.45,7,6s0.45,1,1,1s1-0.45,1-1S8.55,5,8,5z M19,13c0-0.55-0.45-1-1-1h-4c-0.55,0-1,0.45-1,1s0.45,1,1,1h4C18.55,14,19,13.55,19,13z M8,17c-2.01,0-3.74,1.23-4.46,3h8.92c-0.72-1.77-2.45-3-4.46-3z M12,17c0.56,0,1.08,0.15,1.53,0.42c0.44-0.79,1.11-1.45,1.9-1.91C14.7,15.2,13.43,15,12,15c-1.42,0-2.7,0.2-3.43,0.5c0.8,0.46,1.46,1.13,1.9,1.91C10.92,17.15,11.44,17,12,17z" />
              </svg>
              <svg v-else-if="type === 'system'" viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
              </svg>
              <svg v-else viewBox="0 0 24 24" width="24" height="24">
                <path
                  d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
              </svg>
            </slot>
          </div>
          <span class="push-title">{{ title }}</span>
        </div>
        <button class="close-btn" @click="closeNotification">×</button>
      </div>
      <div class="push-body">
        <slot>{{ message }}</slot>
      </div>
      <div v-if="showProgressBar && duration > 0" class="progress-bar">
        <div class="progress" :style="{ animationDuration: duration + 'ms' }"></div>
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
    showProgressBar: {
      type: Boolean,
      default: true,
    },
    type: {
      type: String,
      default: 'default',
      validator: (value) => ['default', 'like', 'comment', 'reply', 'favorite', 'follow', 'system'].includes(value),
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
  width: 340px;
  background-color: #1e1e2f;
  border-radius: 16px;
  overflow: hidden;
  z-index: 1000;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.5), 0 0 1px rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  font-family: 'Poppins', Arial, sans-serif;
  color: #f4f4f9;
  border-left: 4px solid #7c4dff;
  animation: pop-in 0.4s cubic-bezier(0.23, 1, 0.32, 1);
  transform-origin: bottom right;
}

/* 根据位置设置位置 */
.push-notification[class*="top-right"] {
  top: 20px;
  right: 20px;
  bottom: auto;
  left: auto;
  transform-origin: top right;
}

.push-notification[class*="top-left"] {
  top: 20px;
  left: 20px;
  bottom: auto;
  right: auto;
  transform-origin: top left;
}

.push-notification[class*="bottom-left"] {
  bottom: 20px;
  left: 20px;
  top: auto;
  right: auto;
  transform-origin: bottom left;
}

/* 不同类型的样式 */
.push-notification.like {
  border-left-color: #ff5252;
  /* 红色系 - 点赞 */
}

.push-notification.comment {
  border-left-color: #2196f3;
  /* 蓝色系 - 评论 */
}

.push-notification.reply {
  border-left-color: #03a9f4;
  /* 浅蓝色系 - 回复 */
}

.push-notification.favorite {
  border-left-color: #ff9800;
  /* 橙色系 - 收藏 */
}

.push-notification.follow {
  border-left-color: #7c4dff;
  /* 紫色系 - 关注 */
}

.push-notification.system {
  border-left-color: #f44336;
  /* 红色系 - 系统消息 */
}

/* 下面是原有的类型样式，可以考虑移除或保留作为备用 */
.push-notification.success {
  border-left-color: #4caf50;
}

.push-notification.error {
  border-left-color: #f44336;
}

.push-notification.warning {
  border-left-color: #ff9800;
}

.push-notification.info {
  border-left-color: #2196f3;
}

.push-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(30, 30, 47, 0.95);
  padding: 16px;
  font-size: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.icon-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.notification-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-icon svg {
  fill: #7c4dff;
  width: 24px;
  height: 24px;
}

.like .notification-icon svg {
  fill: #ff5252;
  /* 红色系 - 点赞 */
}

.comment .notification-icon svg {
  fill: #2196f3;
  /* 蓝色系 - 评论 */
}

.reply .notification-icon svg {
  fill: #03a9f4;
  /* 浅蓝色系 - 回复 */
}

.favorite .notification-icon svg {
  fill: #ff9800;
  /* 橙色系 - 收藏 */
}

.follow .notification-icon svg {
  fill: #7c4dff;
  /* 紫色系 - 关注 */
}

.system .notification-icon svg {
  fill: #f44336;
  /* 红色系 - 系统消息 */
}

/* 原有的类型样式 */
.success .notification-icon svg {
  fill: #4caf50;
}

.error .notification-icon svg {
  fill: #f44336;
}

.warning .notification-icon svg {
  fill: #ff9800;
}

.info .notification-icon svg {
  fill: #2196f3;
}

.push-title {
  font-weight: 600;
  font-size: 16px;
  letter-spacing: 0.5px;
}

.close-btn {
  background: none;
  border: none;
  color: #f4f4f9;
  font-size: 22px;
  cursor: pointer;
  transition: all 0.3s;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.close-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #ffffff;
  transform: rotate(90deg);
}

.push-body {
  padding: 16px;
  color: #d1c4e9;
  font-size: 14px;
  line-height: 1.6;
  background-color: #2d2d3d;
  border-top: 1px solid rgba(255, 255, 255, 0.02);
}

/* 进度条 */
.progress-bar {
  height: 4px;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.05);
  position: relative;
  overflow: hidden;
}

.progress {
  height: 100%;
  width: 100%;
  background: linear-gradient(90deg, #7c4dff, #ff6e7f);
  transform-origin: left;
  animation: progress-shrink 5000ms linear forwards;
}

.like .progress {
  background: linear-gradient(90deg, #ff5252, #ff8a8a);
  /* 红色渐变 - 点赞 */
}

.comment .progress {
  background: linear-gradient(90deg, #2196f3, #64b5f6);
  /* 蓝色渐变 - 评论 */
}

.reply .progress {
  background: linear-gradient(90deg, #03a9f4, #4fc3f7);
  /* 浅蓝色渐变 - 回复 */
}

.favorite .progress {
  background: linear-gradient(90deg, #ff9800, #ffb74d);
  /* 橙色渐变 - 收藏 */
}

.follow .progress {
  background: linear-gradient(90deg, #7c4dff, #b388ff);
  /* 紫色渐变 - 关注 */
}

.system .progress {
  background: linear-gradient(90deg, #f44336, #ff8a80);
  /* 红色渐变 - 系统消息 */
}

/* 原有样式 */
.success .progress {
  background: linear-gradient(90deg, #43a047, #7cb342);
}

.error .progress {
  background: linear-gradient(90deg, #e53935, #ff5252);
}

.warning .progress {
  background: linear-gradient(90deg, #ff9800, #ffb74d);
}

.info .progress {
  background: linear-gradient(90deg, #2196f3, #64b5f6);
}

/* 改进的动画 */
@keyframes progress-shrink {
  from {
    transform: scaleX(1);
  }

  to {
    transform: scaleX(0);
  }
}

@keyframes pop-in {
  0% {
    opacity: 0;
    transform: scale(0.8) translateY(20px);
  }

  70% {
    transform: scale(1.05);
  }

  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* 过渡动画 */
.slide-fade-enter-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.55, 0, 1, 0.45);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(30px) scale(0.9);
}
</style>