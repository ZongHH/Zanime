<template>
    <transition name="dropdown-fade">
        <div v-if="showNotifications" class="notification-dropdown">
            <div class="notify-header">
                <h5>通知</h5>
                <span class="notify-badge" v-if="notifications.length">
                    {{ notifications.length }}
                </span>
            </div>

            <div class="notify-body">
                <div v-for="notification in notifications" :key="notification.id" class="notify-item">
                    <div class="notify-icon">
                        <i class="fas fa-bell"></i>
                    </div>
                    <div class="notify-content">
                        <h6>{{ notification.title }}</h6>
                        <p>{{ notification.message }}</p>
                        <div class="notify-meta">
                            <span class="notify-time">
                                <i class="far fa-clock"></i>
                                {{ notification.time }}
                            </span>
                        </div>
                    </div>
                </div>

                <div v-if="!notifications.length" class="notify-empty">
                    <i class="fas fa-bell-slash"></i>
                    <p>暂无新消息</p>
                </div>
            </div>

            <div class="notify-footer" v-if="notifications.length">
                <button class="view-all">
                    查看全部消息
                    <i class="fas fa-chevron-right"></i>
                </button>
            </div>
        </div>
    </transition>
</template>

<script>
export default {
    props: {
        showNotifications: {
            type: Boolean,
            required: true
        },
        notifications: {
            type: Array,
            required: true
        }
    }
};
</script>

<style scoped>
@import "@/static/css/notify.css";
</style>