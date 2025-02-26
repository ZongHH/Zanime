import { createRouter, createWebHistory } from 'vue-router'
import Main from '@/views/Main.vue'
import Users from '@/views/Users.vue'
import Login from '@/views/Login.vue'
import Logs from '../views/Logs.vue'

const routes = [
    {
        path: '/',
        redirect: '/dashboard'
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: Main,
        meta: {
            title: '控制台'
        }
    },
    {
        path: '/users',
        name: 'Users',
        component: Users,
        meta: {
            title: '用户管理'
        }
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {
            title: '请登录'
        }
    },
    {
        path: '/logs',
        name: 'Logs',
        component: Logs,
        meta: {
            title: '详细日志'
        }
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    document.title = `${to.meta.title} | Zanime监控系统`
    next()
})

export default router
