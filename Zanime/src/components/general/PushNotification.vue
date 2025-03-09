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
      default: 10000, // 自动关闭时间 (毫秒)
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
  width: 320px;
  background-color: #ffffff;
  border-radius: 10px;
  overflow: hidden;
  z-index: 1000;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.08), 0 2px 5px rgba(0, 0, 0, 0.03);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  color: #333333;
  border-left: 3px solid #666666;
  animation: elegant-fade-in 0.4s cubic-bezier(0.215, 0.61, 0.355, 1);
  backdrop-filter: blur(10px);
  transform-origin: center;
}

/* 根据位置设置位置 */
.push-notification[class*="top-right"] {
  top: 20px;
  right: 20px;
  bottom: auto;
  left: auto;
}

.push-notification[class*="top-left"] {
  top: 20px;
  left: 20px;
  bottom: auto;
  right: auto;
}

.push-notification[class*="bottom-left"] {
  bottom: 20px;
  left: 20px;
  top: auto;
  right: auto;
}

/* 不同类型的样式 */
.push-notification.like {
  border-left-color: #f44336;
  box-shadow: 0 10px 25px rgba(244, 67, 54, 0.08), 0 2px 5px rgba(244, 67, 54, 0.03);
}

.push-notification.comment {
  border-left-color: #2196f3;
  box-shadow: 0 10px 25px rgba(33, 150, 243, 0.08), 0 2px 5px rgba(33, 150, 243, 0.03);
}

.push-notification.reply {
  border-left-color: #03a9f4;
  box-shadow: 0 10px 25px rgba(3, 169, 244, 0.08), 0 2px 5px rgba(3, 169, 244, 0.03);
}

.push-notification.favorite {
  border-left-color: #ff9800;
  box-shadow: 0 10px 25px rgba(255, 152, 0, 0.08), 0 2px 5px rgba(255, 152, 0, 0.03);
}

.push-notification.follow {
  border-left-color: #673ab7;
  box-shadow: 0 10px 25px rgba(103, 58, 183, 0.08), 0 2px 5px rgba(103, 58, 183, 0.03);
}

.push-notification.system {
  border-left-color: #607d8b;
  box-shadow: 0 10px 25px rgba(96, 125, 139, 0.08), 0 2px 5px rgba(96, 125, 139, 0.03);
}

.push-notification.success {
  border-left-color: #4caf50;
  box-shadow: 0 10px 25px rgba(76, 175, 80, 0.08), 0 2px 5px rgba(76, 175, 80, 0.03);
}

.push-notification.error {
  border-left-color: #f44336;
  box-shadow: 0 10px 25px rgba(244, 67, 54, 0.08), 0 2px 5px rgba(244, 67, 54, 0.03);
}

.push-notification.warning {
  border-left-color: #ff9800;
  box-shadow: 0 10px 25px rgba(255, 152, 0, 0.08), 0 2px 5px rgba(255, 152, 0, 0.03);
}

.push-notification.info {
  border-left-color: #2196f3;
  box-shadow: 0 10px 25px rgba(33, 150, 243, 0.08), 0 2px 5px rgba(33, 150, 243, 0.03);
}

.push-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  font-size: 15px;
  background-color: #fafafa;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  position: relative;
}

.push-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, rgba(0, 0, 0, 0.01), rgba(0, 0, 0, 0.05), rgba(0, 0, 0, 0.01));
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
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background-color: rgba(102, 102, 102, 0.08);
  transition: transform 0.3s ease;
}

.push-notification:hover .notification-icon {
  transform: scale(1.05);
}

.notification-icon svg {
  fill: #666666;
  width: 18px;
  height: 18px;
}

.like .notification-icon {
  background-color: rgba(244, 67, 54, 0.08);
}

.like .notification-icon svg {
  fill: #f44336;
}

.comment .notification-icon {
  background-color: rgba(33, 150, 243, 0.08);
}

.comment .notification-icon svg {
  fill: #2196f3;
}

.reply .notification-icon {
  background-color: rgba(3, 169, 244, 0.08);
}

.reply .notification-icon svg {
  fill: #03a9f4;
}

.favorite .notification-icon {
  background-color: rgba(255, 152, 0, 0.08);
}

.favorite .notification-icon svg {
  fill: #ff9800;
}

.follow .notification-icon {
  background-color: rgba(103, 58, 183, 0.08);
}

.follow .notification-icon svg {
  fill: #673ab7;
}

.system .notification-icon {
  background-color: rgba(96, 125, 139, 0.08);
}

.system .notification-icon svg {
  fill: #607d8b;
}

.success .notification-icon {
  background-color: rgba(76, 175, 80, 0.08);
}

.success .notification-icon svg {
  fill: #4caf50;
}

.error .notification-icon {
  background-color: rgba(244, 67, 54, 0.08);
}

.error .notification-icon svg {
  fill: #f44336;
}

.warning .notification-icon {
  background-color: rgba(255, 152, 0, 0.08);
}

.warning .notification-icon svg {
  fill: #ff9800;
}

.info .notification-icon {
  background-color: rgba(33, 150, 243, 0.08);
}

.info .notification-icon svg {
  fill: #2196f3;
}

.push-title {
  font-weight: 500;
  font-size: 15px;
  letter-spacing: 0.2px;
  color: #222222;
}

.close-btn {
  background: none;
  border: none;
  color: #888888;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  opacity: 0.6;
}

.close-btn:hover {
  color: #333333;
  background-color: rgba(0, 0, 0, 0.05);
  opacity: 1;
}

.push-body {
  padding: 16px 18px;
  font-size: 14px;
  line-height: 1.5;
  background-color: #ffffff;
  color: #555555;
  letter-spacing: 0.1px;
}

/* 进度条 */
.progress-bar {
  height: 2px;
  width: 100%;
  background-color: rgba(0, 0, 0, 0.03);
  position: relative;
  overflow: hidden;
}

.progress {
  height: 100%;
  width: 100%;
  background: linear-gradient(90deg, rgba(102, 102, 102, 0.7), rgba(102, 102, 102, 0.9));
  transform-origin: left;
  animation: progress-shrink linear forwards;
  border-radius: 2px;
}

.like .progress {
  background: linear-gradient(90deg, rgba(244, 67, 54, 0.7), rgba(244, 67, 54, 0.9));
}

.comment .progress {
  background: linear-gradient(90deg, rgba(33, 150, 243, 0.7), rgba(33, 150, 243, 0.9));
}

.reply .progress {
  background: linear-gradient(90deg, rgba(3, 169, 244, 0.7), rgba(3, 169, 244, 0.9));
}

.favorite .progress {
  background: linear-gradient(90deg, rgba(255, 152, 0, 0.7), rgba(255, 152, 0, 0.9));
}

.follow .progress {
  background: linear-gradient(90deg, rgba(103, 58, 183, 0.7), rgba(103, 58, 183, 0.9));
}

.system .progress {
  background: linear-gradient(90deg, rgba(96, 125, 139, 0.7), rgba(96, 125, 139, 0.9));
}

.success .progress {
  background: linear-gradient(90deg, rgba(76, 175, 80, 0.7), rgba(76, 175, 80, 0.9));
}

.error .progress {
  background: linear-gradient(90deg, rgba(244, 67, 54, 0.7), rgba(244, 67, 54, 0.9));
}

.warning .progress {
  background: linear-gradient(90deg, rgba(255, 152, 0, 0.7), rgba(255, 152, 0, 0.9));
}

.info .progress {
  background: linear-gradient(90deg, rgba(33, 150, 243, 0.7), rgba(33, 150, 243, 0.9));
}

/* 优雅的动画 */
@keyframes progress-shrink {
  from {
    transform: scaleX(1);
  }

  to {
    transform: scaleX(0);
  }
}

@keyframes elegant-fade-in {
  0% {
    opacity: 0;
    transform: translateY(8px) scale(0.98);
  }

  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 过渡动画 */
.slide-fade-enter-active {
  transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.19, 1, 0.22, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.98);
}

/* 响应式调整 */
@media (max-width: 480px) {
  .push-notification {
    width: calc(100% - 40px);
    max-width: 320px;
  }
}
</style>