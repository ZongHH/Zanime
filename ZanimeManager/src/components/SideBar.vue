<template>
    <aside class="fixed left-0 top-16 bottom-0 w-64 bg-white border-r border-gray-200 overflow-y-auto"
        :class="{ 'w-20': isCollapsed }">
        <div class="p-4">
            <button v-for="(item, index) in menuItems" :key="index"
                class="flex items-center w-full px-4 py-2 rounded-lg text-left text-sm font-medium transition-colors mb-2 !rounded-button"
                :class="[activeIndex === index ? 'bg-indigo-50 text-indigo-600' : 'text-gray-600 hover:bg-gray-50']"
                @click="handleMenuClick(index)">
                <i :class="[item.icon, 'mr-3']"></i>
                <span v-if="!isCollapsed">{{ item.name }}</span>
            </button>
        </div>
    </aside>
</template>

<script>
export default {
    name: 'SideBar',

    props: {
        // 是否折叠
        isCollapsed: {
            type: Boolean,
            default: false
        },
        // 当前激活的菜单项索引
        activeIndex: {
            type: Number,
            default: 0
        },
        // 菜单项数据
        menuItems: {
            type: Array,
            default: () => [
                { name: '首页概览', icon: 'fas fa-home' },
                { name: '内容管理', icon: 'fas fa-film' },
                { name: '用户管理', icon: 'fas fa-users' },
                { name: '数据统计', icon: 'fas fa-chart-bar' },
                { name: '系统设置', icon: 'fas fa-cog' }
            ]
        }
    },

    methods: {
        handleMenuClick(index) {
            this.$emit('menu-click', index);
        }
    }
}
</script>

<style scoped>
/* 侧边栏动画过渡效果 */
aside {
    transition: width 0.3s;
}
</style>