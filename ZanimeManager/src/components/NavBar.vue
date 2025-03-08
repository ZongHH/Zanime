<template>
    <nav class="fixed top-0 left-0 right-0 bg-white shadow-sm h-16 z-50">
        <div class="max-w-[1440px] mx-auto px-6 h-full flex items-center justify-between">
            <!-- 左侧Logo和标题 -->
            <div class="flex items-center gap-4">
                <i class="fas fa-play-circle text-2xl text-indigo-600" v-if="logo"></i>
                <h1 class="text-xl font-bold">动漫管理系统</h1>
                <!-- 折叠按钮 -->
                <button @click="toggleSidebar"
                    class="p-2 rounded-lg hover:bg-gray-100 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500">
                    <i :class="['fas', isCollapse ? 'fa-expand-alt' : 'fa-compress-alt', 'text-gray-600']"></i>
                </button>
            </div>

            <!-- 中间搜索框 -->
            <div class="relative w-96">
                <input type="text" v-model="searchKey" placeholder="搜索动漫、分类、标签..."
                    class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm">
                <i class="fas fa-search absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
            </div>

            <!-- 右侧功能区 -->
            <div class="flex items-center gap-6">
                <!-- 全屏按钮 -->
                <button @click="toggleFullScreen"
                    class="p-2 rounded-lg hover:bg-gray-100 transition-colors group relative" title="全屏显示">
                    <i class="fas fa-expand text-gray-600"></i>
                    <!-- 工具提示 - 改为向下显示 -->
                    <div
                        class="absolute top-full left-1/2 -translate-x-1/2 mt-2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                        全屏显示
                    </div>
                </button>

                <!-- 主题切换 -->
                <button @click="toggleTheme" class="p-2 rounded-lg hover:bg-gray-100 transition-colors group relative"
                    title="切换主题">
                    <i :class="['fas', isDark ? 'fa-sun' : 'fa-moon', 'text-gray-600']"></i>
                    <!-- 工具提示 - 改为向下显示 -->
                    <div
                        class="absolute top-full left-1/2 -translate-x-1/2 mt-2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                        切换主题
                    </div>
                </button>

                <!-- 通知按钮 -->
                <button class="relative p-2 rounded-lg hover:bg-gray-100 transition-colors">
                    <i class="fas fa-bell text-gray-600"></i>
                    <span
                        class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full text-white text-xs flex items-center justify-center">
                        3
                    </span>
                </button>

                <!-- 用户信息下拉菜单 -->
                <div class="relative" v-click-outside="closeUserMenu">
                    <button @click="toggleUserMenu"
                        class="flex items-center gap-2 p-1 rounded-lg hover:bg-gray-100 transition-colors">
                        <img :src="userAvatar" alt="用户头像" class="w-8 h-8 rounded-full object-cover">
                        <span class="text-sm font-medium">{{ username }}</span>
                        <i class="fas fa-chevron-down text-xs text-gray-500"></i>
                    </button>

                    <!-- 下拉菜单 -->
                    <div v-show="showUserMenu"
                        class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-1 border border-gray-100 z-50">
                        <a href="#" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">
                            <i class="fas fa-user-circle w-5 text-gray-400"></i>
                            个人信息
                        </a>
                        <a href="#" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">
                            <i class="fas fa-key w-5 text-gray-400"></i>
                            修改密码
                        </a>
                        <div class="h-px bg-gray-200 my-1"></div>
                        <button @click="handleLogout"
                            class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-gray-50">
                            <i class="fas fa-sign-out-alt w-5"></i>
                            退出登录
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </nav>
</template>

<script>
export default {
    name: 'NavBar',

    props: {
        logo: {
            type: String,
            default: 'True'
        }
    },

    data() {
        return {
            searchKey: '',
            isCollapse: false,
            isDark: false,
            username: '陈思远',
            userAvatar: 'https://ai-public.mastergo.com/ai/img_res/37b2391a50f6a4893ef76095947a7646.jpg',
            showUserMenu: false
        }
    },

    directives: {
        'click-outside': {
            mounted(el, binding) {
                el.clickOutsideEvent = function (event) {
                    if (!(el === event.target || el.contains(event.target))) {
                        binding.value(event);
                    }
                };
                document.addEventListener('click', el.clickOutsideEvent);
            },
            unmounted(el) {
                document.removeEventListener('click', el.clickOutsideEvent);
            }
        }
    },

    methods: {
        toggleSidebar() {
            this.isCollapse = !this.isCollapse;
            this.$emit('toggle-sidebar', this.isCollapse);
        },

        toggleFullScreen() {
            if (!document.fullscreenElement) {
                document.documentElement.requestFullscreen();
            } else {
                document.exitFullscreen();
            }
        },

        toggleTheme() {
            this.isDark = !this.isDark;
            // 这里可以添加切换主题的具体实现
        },

        toggleUserMenu() {
            this.showUserMenu = !this.showUserMenu;
        },

        closeUserMenu() {
            this.showUserMenu = false;
        },

        handleLogout() {
            this.$message.success('退出成功');
            // 这里可以添加退出登录的具体实现
            this.showUserMenu = false;
        }
    }
}
</script>

<style scoped>
/* 响应式设计 */
@media screen and (max-width: 768px) {
    .relative.w-96 {
        display: none;
    }

    h1 {
        display: none;
    }
}
</style>
