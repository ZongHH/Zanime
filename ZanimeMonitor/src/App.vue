<template>
  <NavBar />
  <router-view></router-view>
</template>

<script>
import NavBar from './components/NavBar.vue';

export default {
  components: {
    NavBar,
  },
  created() {
    // 检查登录状态并初始化WebSocket
    const isLoggedIn = localStorage.getItem('isLoggedIn')
    if (isLoggedIn === 'true') {
      const wsUrl = 'ws://127.0.0.1:9999/ws'
      this.$store.dispatch('initWebSocket', wsUrl)
    }
  },
}
</script>

<style></style>