<template>
  <div>
    <!-- 只在非登录页面显示导航栏和侧边栏 -->
    <template v-if="!isLoginPage">
      <div class="min-h-screen bg-gray-50">
        <NavBar @toggle-sidebar="handleSidebarToggle" />
        <div class="flex pt-16">
          <SideBar :is-collapsed="isSidebarCollapsed" :active-index="currentMenu" :menu-items="menuItems"
            @menu-click="handleMenuClick" />
          <main class="flex-1" :class="{ 'ml-64': !isSidebarCollapsed, 'ml-20': isSidebarCollapsed }">
            <router-view></router-view>
          </main>
        </div>
      </div>
    </template>

    <!-- 登录页面全屏显示 -->
    <template v-else>
      <router-view></router-view>
    </template>
  </div>
</template>

<script>
import NavBar from './components/NavBar.vue';
import SideBar from './components/SideBar.vue';

export default {
  components: {
    NavBar,
    SideBar
  },

  data() {
    return {
      isSidebarCollapsed: false,
      currentMenu: 0,
      menuItems: [
        { name: '首页概览', icon: 'fas fa-home' },
        { name: '内容管理', icon: 'fas fa-film' },
        { name: '用户管理', icon: 'fas fa-users' },
        { name: '数据统计', icon: 'fas fa-chart-bar' },
        { name: '系统设置', icon: 'fas fa-cog' }
      ]
    };
  },

  computed: {
    // 判断当前是否是登录页面
    isLoginPage() {
      return this.$route.path === '/login';
    }
  },

  methods: {
    handleSidebarToggle(collapsed) {
      this.isSidebarCollapsed = collapsed;
    },

    handleMenuClick(index) {
      this.currentMenu = index;
      // 根据菜单索引跳转到相应路由
      const routes = ['/', '/content', '/users', '/statistics', '/settings'];
      if (routes[index]) {
        this.$router.push(routes[index]);
      }
    }
  },

  created() {
    // 根据当前路由设置活跃菜单
    const currentPath = this.$router.currentRoute.value.path;
    const routes = ['/', '/content', '/users', '/statistics', '/settings'];
    const index = routes.findIndex(route => route === currentPath);
    if (index !== -1) {
      this.currentMenu = index;
    }
  }
}
</script>

<style>
/* 全局样式 */
body {
  margin: 0;
  font-family: 'Inter', sans-serif;
}

/* 页面过渡动画 */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.3s;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>