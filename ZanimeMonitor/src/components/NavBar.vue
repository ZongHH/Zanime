<template>
    <div class="navbar">
        <div class="left-part">
            <!-- Logo -->
            <div class="logo">
                <img src="@/assets/logo.png" alt="Logo" v-if="logo">
                <span class="title">Zanime Monitor</span>
            </div>
            <!-- 折叠按钮 -->
            <el-icon class="fold-btn" @click="toggleSidebar">
                <component :is="isCollapse ? 'Expand' : 'Fold'" />
            </el-icon>
        </div>

        <div class="right-part">
            <!-- 搜索框 -->
            <el-input v-model="searchKey" placeholder="搜索..." class="search-input" :prefix-icon="Search" />

            <!-- 全屏按钮 -->
            <el-tooltip content="全屏显示" placement="bottom">
                <el-icon class="right-menu-item" @click="toggleFullScreen">
                    <FullScreen />
                </el-icon>
            </el-tooltip>

            <!-- 主题切换 -->
            <el-tooltip content="主题切换" placement="bottom">
                <el-icon class="right-menu-item" @click="toggleTheme">
                    <component :is="isDark ? 'Sunny' : 'Moon'" />
                </el-icon>
            </el-tooltip>

            <!-- 用户信息 -->
            <el-dropdown trigger="click">
                <div class="user-info">
                    <el-avatar :size="32" :src="userAvatar" />
                    <span class="username">{{ username }}</span>
                </div>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>个人信息</el-dropdown-item>
                        <el-dropdown-item>修改密码</el-dropdown-item>
                        <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
</template>

<script>
import { Search, FullScreen, Sunny, Moon, Expand, Fold } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

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
            username: 'Admin',
            userAvatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
        }
    },

    components: {
        Search,
        FullScreen,
        Sunny,
        Moon,
        Expand,
        Fold
    },

    methods: {
        toggleSidebar() {
            this.isCollapse = !this.isCollapse
            this.$emit('toggle-sidebar', this.isCollapse)
        },

        toggleFullScreen() {
            if (!document.fullscreenElement) {
                document.documentElement.requestFullscreen()
            } else {
                document.exitFullscreen()
            }
        },

        toggleTheme() {
            this.isDark = !this.isDark
            // 这里可以添加切换主题的具体实现
        },

        handleLogout() {
            ElMessage.success('退出成功')
            // 这里可以添加退出登录的具体实现
        }
    }
}
</script>

<style scoped>
.navbar {
    height: 60px;
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 20px;
    background-color: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.left-part {
    display: flex;
    align-items: center;
}

.logo {
    display: flex;
    align-items: center;
    margin-right: 20px;
}

.logo img {
    height: 32px;
    margin-right: 12px;
}

.title {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
}

.fold-btn {
    font-size: 20px;
    cursor: pointer;
    padding: 8px;
    border-radius: 4px;
    color: #303133;
    transition: background-color 0.3s;
}

.fold-btn:hover {
    background-color: #f5f7fa;
}

.right-part {
    display: flex;
    align-items: center;
    gap: 20px;
}

.search-input {
    width: 200px;
}

.right-menu-item {
    font-size: 20px;
    cursor: pointer;
    padding: 8px;
    border-radius: 4px;
    color: #303133;
    transition: background-color 0.3s;
}

.right-menu-item:hover {
    background-color: #f5f7fa;
}

.user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding-right: 20px;
    max-width: 150px;
    /* 添加最大宽度限制 */
}

.username {
    margin-left: 8px;
    font-size: 14px;
    color: #303133;
    white-space: nowrap;
    /* 文本不换行 */
    overflow: hidden;
    /* 溢出隐藏 */
    text-overflow: ellipsis;
    /* 显示省略号 */
    max-width: 80px;
    /* 用户名最大宽度 */
}

@media screen and (max-width: 768px) {
    .search-input {
        display: none;
    }

    .title {
        display: none;
    }

    .username {
        display: none;
    }
}
</style>
