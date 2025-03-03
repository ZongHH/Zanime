import { createRouter, createWebHistory } from 'vue-router';
import axios from 'axios';

// 导入需要路由的组件
import Index from '@/views/Index.vue';
import Login from '@/views/Login.vue';
import MoviesDetail from '@/views/MoviesDetail.vue';
import Order from '@/views/Order.vue';
import AnimeDetail from '@/views/AnimeDetail.vue';
import AnimeLibrary from '@/views/AnimeLibrary.vue';
import NotFound from '@/views/NotFound.vue'
import Personal from '@/views/Personal.vue';
import Discussion from '@/views/Discussion.vue';
import Post from '@/views/Post.vue';
import Feedback from '@/views/Feedback.vue';

const routes = [
  {
    path: '/', // 根路径
    name: 'Index',
    component: Index,
    meta: { requiresAuth: true } // 需要登录才能访问
  },
  {
    path: '/login', // 登录页面路径
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false } // 不需要登录就能访问
  },
  {
    path: '/moviesDetail',
    name: 'MoviesDetail',
    component: MoviesDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: Order,
    meta: { requiresAuth: true }
  },
  {
    path: '/animeDetail',
    name: 'AnimeDetail',
    component: AnimeDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/animeLibrary',
    name: 'AnimeLibrary',
    component: AnimeLibrary,
    meta: { requiresAuth: true }
  },
  {
    path: '/personal',
    name: 'Personal',
    component: Personal,
    meta: { requiresAuth: true }
  },
  {
    path: '/discussion',
    name: 'Discussion',
    component: Discussion,
    meta: { requiresAuth: true }
  },
  {
    path: '/post/:id',
    name: 'Post',
    component: Post,
    meta: { requiresAuth: true }
  },
  {
    path: '/feedback',
    name: 'Feedback',
    component: Feedback,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: { requiresAuth: false }
  },
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(), // 使用HTML5历史模式
  routes, // 路由配置
  // 控制页面滚动行为
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      // 如果存在保存的位置信息,则返回到之前的位置
      return savedPosition
    } else {
      // 否则滚动到页面顶部
      return { top: 0 }
    }
  }
});

// 全局前置守卫 - 在路由跳转前进行权限验证
router.beforeEach(async (to, from, next) => {
  // 从localStorage获取用户ID
  const user_id = localStorage.getItem('user_id');

  // 检查目标路由是否需要认证
  if (to.meta.requiresAuth) {
    // 如果需要认证但没有用户ID,重定向到登录页
    if (!user_id) {
      next({
        path: '/login',
        query: { redirect: to.fullPath } // 保存原目标路径用于登录后跳转
      });
      return;
    }

    try {
      // 向后端发送请求验证token有效性
      const response = await axios.get('/api/verify-token?user_id=' + user_id);
      // 如果验证失败
      if (response.data.code != 200) {
        console.error('Token verification failed:', error);
        localStorage.clear(); // 清除本地存储
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
      } else {
        // 验证成功,允许访问
        next();
      }
    } catch (error) {
      // 请求出错,清除用户信息并重定向到登录页
      console.error('Token verification failed:', error);
      localStorage.clear();
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      });
    }
  } else {
    // 不需要认证的路由直接放行
    next();
  }
});

export default router;
