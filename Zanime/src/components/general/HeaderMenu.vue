<template>
    <!-- HEADER MENU START -->
    <header class="header">
        <nav class="navigation d-flex align-items-center justify-content-between">
            <a href="" class="d-flex align-items-center">
                <img src="@/static/picture/logo.png" alt="/logo" class="header-logo">
            </a>
            <div class="menu-button-right">
                <div class="main-menu__nav">
                    <ul class="main-menu__list">
                        <li>
                            <router-link to="/" :class="{ active: activeMenuItem === 'home' }">首页</router-link>
                        </li>
                        <li class="dropdown">
                            <router-link to="/animeLibrary" :class="{ active: activeMenuItem === 'anime' }">
                                动漫
                                <i class="fas fa-chevron-down"></i>
                            </router-link>
                            <ul class="dropdown-menu">
                                <li><router-link to="/animeLibrary?type=latest">最新更新</router-link></li>
                                <li><router-link to="/animeLibrary?type=popular">热门动漫</router-link></li>
                                <li class="disabled"><a href="#" @click.prevent>即将上线 <span
                                            class="coming-soon">开发中</span></a></li>
                                <li class="disabled"><a href="#" @click.prevent>完结动漫 <span
                                            class="coming-soon">开发中</span></a></li>
                            </ul>
                        </li>
                        <li class="dropdown disabled">
                            <a href="#" @click.prevent>
                                分类
                                <i class="fas fa-chevron-down"></i>
                                <span class="coming-soon">开发中</span>
                            </a>
                        </li>
                        <li class="dropdown disabled">
                            <a href="#" @click.prevent>
                                时间表
                                <i class="fas fa-chevron-down"></i>
                                <span class="coming-soon">开发中</span>
                            </a>
                        </li>
                        <li class="dropdown">
                            <router-link to="/discussion" :class="{ active: activeMenuItem === 'community' }">
                                社区
                                <i class="fas fa-chevron-down"></i>
                            </router-link>
                            <ul class="dropdown-menu">
                                <li><router-link to="/discussion">讨论区</router-link></li>
                                <li class="disabled"><a href="#" @click.prevent>动漫资讯 <span
                                            class="coming-soon">开发中</span></a></li>
                                <li class="disabled"><a href="#" @click.prevent>排行榜 <span
                                            class="coming-soon">开发中</span></a></li>
                            </ul>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="d-flex align-items-center d-xl-flex d-none">
                <div class="input-box">
                    <!-- 修改搜索表单 -->
                    <form @submit.prevent="handleSearch">
                        <input type="text" v-model="searchQuery" placeholder="搜索..." @input="handleSearchInput"
                            @focus="showDropdown = true" @blur="handleBlur" required>
                        <button class="search" type="submit">
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20"
                                fill="none">
                                <g clip-path="url(#clip0_11668_3150)">
                                    <path
                                        d="M8.11719 0C12.593 0 16.2344 3.64137 16.2344 8.11719C16.2344 10.1445 15.4873 12.0007 14.2539 13.4247L19.8284 18.9998C20.0572 19.2286 20.0572 19.5996 19.8284 19.8284C19.5995 20.0573 19.2286 20.0572 18.9997 19.8284L13.4254 14.2534C12.0012 15.4871 10.1448 16.2344 8.11719 16.2344C3.64137 16.2344 1.90735e-06 12.593 1.90735e-06 8.11719C1.90735e-06 3.64137 3.64137 0 8.11719 0ZM8.11719 15.0625C11.9469 15.0625 15.0625 11.9468 15.0625 8.11719C15.0625 4.28754 11.9468 1.17187 8.11719 1.17187C4.28754 1.17187 1.17188 4.28754 1.17188 8.11719C1.17188 11.9468 4.28754 15.0625 8.11719 15.0625Z"
                                        fill="#FAFAFA"></path>
                                </g>
                                <defs>
                                    <clippath id="clip0_11668_3150">
                                        <rect width="20" height="20" fill="white"></rect>
                                    </clippath>
                                </defs>
                            </svg>
                        </button>
                    </form>
                    <i class="fal fa-times close-icon"></i>
                    <!-- 添加搜索结果下拉框 -->
                    <div class="search-dropdown" v-if="showDropdown">
                        <ul v-if="searchResults.length > 0">
                            <li v-for="result in searchResults" :key="result.id" @mousedown="handleResultClick(result)">
                                <div class="search-result-item">
                                    <img :src="result.cover_image_url" :alt="result.video_name"
                                        class="result-thumbnail">
                                    <div class="result-info">
                                        <div class="result-title">{{ result.video_name }}</div>
                                        <div class="result-type">{{ result.anime_type }}</div>
                                    </div>
                                </div>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="icon">
                    <a @click="toggleNotifications">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20" fill="none">
                            <path
                                d="M16.0547 10.268V8.39844C16.0547 5.67102 14.2417 3.35934 11.7578 2.6043V1.75781C11.7578 0.788555 10.9692 0 9.99998 0C9.03073 0 8.24217 0.788555 8.24217 1.75781V2.6043C5.75819 3.35934 3.9453 5.67098 3.9453 8.39844V10.268C3.9453 12.6638 3.03209 14.9355 1.37393 16.6648C1.21143 16.8342 1.16577 17.0843 1.25788 17.3002C1.34998 17.5161 1.56209 17.6562 1.79686 17.6562H7.12924C7.40143 18.9919 8.58518 20 9.99998 20C11.4148 20 12.5985 18.9919 12.8707 17.6562H18.2031C18.4379 17.6562 18.6499 17.5161 18.7421 17.3002C18.8342 17.0843 18.7885 16.8342 18.626 16.6648C16.9679 14.9355 16.0547 12.6638 16.0547 10.268ZM9.41405 1.75781C9.41405 1.43473 9.6769 1.17188 9.99998 1.17188C10.3231 1.17188 10.5859 1.43473 10.5859 1.75781V2.37219C10.3931 2.35359 10.1976 2.34375 9.99998 2.34375C9.80233 2.34375 9.6069 2.35359 9.41405 2.37219V1.75781ZM9.99998 18.8281C9.23612 18.8281 8.58483 18.3382 8.34295 17.6562H11.657C11.4151 18.3382 10.7639 18.8281 9.99998 18.8281ZM3.05975 16.4844C4.39416 14.6956 5.11717 12.5309 5.11717 10.268V8.39844C5.11717 5.70605 7.3076 3.51562 9.99998 3.51562C12.6924 3.51562 14.8828 5.70605 14.8828 8.39844V10.268C14.8828 12.5309 15.6058 14.6956 16.9403 16.4844H3.05975Z"
                                fill="#FAFAFA"></path>
                        </svg>
                    </a>
                    <Notify :showNotifications="showNotifications" :notifications="notifications" />
                </div>
                <div class="dropdown-container">
                    <details class="dropdown right">
                        <summary class="avatar">
                            <img src="@/static/picture/Ellipse-1.png" alt="logo">
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewbox="0 0 16 16"
                                fill="none">
                                <path d="M0 4.01306L8.00002 11.987L16 4.01306H0Z" fill="#FAFAFA"></path>
                            </svg>
                        </summary>
                        <ul>
                            <!-- Menu links -->
                            <li>
                                <a href="#" @click.prevent="JumpPersonal">
                                    <i class="fa-light fa-user"></i>
                                    我的主页
                                </a>
                            </li>
                            <li>
                                <a href="#" @click.prevent="JumpOrder">
                                    <i class="fa-light fa-file-invoice"></i>
                                    我的订单
                                </a>
                            </li>
                            <!-- Optional divider -->
                            <li>
                                <a href="#">
                                    <i class="fa-solid fa-info"></i>
                                    帮助
                                </a>
                            </li>
                            <li>
                                <a href="#" id="logoutLink">
                                    <i class="fa-solid fa-right-from-bracket"></i>
                                    注销
                                </a>
                            </li>
                        </ul>
                    </details>
                </div>
            </div>
            <a href="#" class="d-xl-none d-flex main-menu__toggler mobile-nav__toggler">
                <i class="fa-light fa-bars"></i>
            </a>
        </nav>
    </header>
    <!-- HEADER MENU END -->
</template>

<script>
export default {
    data() {
        return {
            showNotifications: false,
            notifications: [
                { id: 1, title: "新番更新！", message: "《鬼灭之刃》最新一集已上线，快来观看！", time: "2分钟前" },
                { id: 2, title: "漫画更新", message: "《海贼王》新章节已发布，立即查看！", time: "15分钟前" },
                { id: 3, title: "活动提醒", message: "下周是漫展最后一天，别错过哦！", time: "1小时前" },
            ],
            searchQuery: '',
            searchResults: [],
            showDropdown: false,
            searchTimeout: null,
            activeMenuItem: '', // 默认不激活
        };
    },
    methods: {
        toggleNotifications() {
            this.showNotifications = !this.showNotifications;
        },
        handleSearchInput() {
            // 防抖处理
            clearTimeout(this.searchTimeout);
            if (this.searchQuery.trim()) {
                this.searchTimeout = setTimeout(() => {
                    this.fetchSearchResults();
                }, 300);
            } else {
                this.searchResults = [];
            }
        },

        async fetchSearchResults() {
            try {
                const response = await axios.get(`/api/search?query=${this.searchQuery}`);
                if (response.data.code === 200) {
                    this.searchResults = response.data.animes;
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                console.error('搜索出错:', error);
                this.searchResults = [];
            }
        },

        handleResultClick(result) {
            // 跳转到详情页
            this.$router.push(`/moviesDetail?videoId=${result.video_id}`);
        },

        handleSearch(e) {
            e.preventDefault();
            // 处理搜索表单提交
            if (this.searchQuery.trim()) {
                // this.fetchSearchResults();
                this.$router.push(`/animeDetail?params=${encodeURIComponent(this.searchQuery.trim())}`);
            }
        },

        handleBlur() {
            // 延迟关闭下拉框，以便能够处理点击事件
            setTimeout(() => {
                this.showDropdown = false;
            }, 200);
        },

        JumpOrder() {
            this.$router.push("/orders")
        },

        JumpPersonal() {
            this.$router.push("/personal")
        }
    },
    mounted() {
        // 根据当前路由路径设置激活项
        const path = window.location.pathname;
        if (path === '/') this.activeMenuItem = 'home';
        else if (path.includes('animeLibrary')) this.activeMenuItem = 'anime';
        else if (path.includes('discussion')) this.activeMenuItem = 'community';
        Logout()
    },
    components: {
        Notify,
    },
};

import Notify from '@/components/Notify.vue';
import axios from 'axios';
import Logout from '@/static/js/general';
</script>

<style scoped>
.input-box {
    position: relative;
}

.search-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: #1A1A1A;
    border: 1px solid #333;
    border-radius: 4px;
    margin-top: 4px;
    max-height: 300px;
    overflow-y: auto;
    z-index: 1000;
}

.search-dropdown ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.search-result-item {
    display: flex;
    align-items: center;
    padding: 10px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.search-result-item:hover {
    background-color: #2A2A2A;
}

.result-thumbnail {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 4px;
    margin-right: 10px;
}

.result-info {
    flex: 1;
}

.result-title {
    color: #FFFFFF;
    font-size: 14px;
    margin-bottom: 4px;
}

.result-type {
    color: #999;
    font-size: 12px;
}

/* 自定义滚动条样式 */
.search-dropdown::-webkit-scrollbar {
    width: 6px;
}

.search-dropdown::-webkit-scrollbar-track {
    background: #1A1A1A;
}

.search-dropdown::-webkit-scrollbar-thumb {
    background: #444;
    border-radius: 3px;
}

.search-dropdown::-webkit-scrollbar-thumb:hover {
    background: #555;
}

/* 添加下拉菜单样式 */
.dropdown {
    position: relative;
}

.dropdown-menu {
    position: absolute;
    top: 100%;
    left: 0;
    background: #1A1A1A;
    border: 1px solid #333;
    border-radius: 8px;
    min-width: 180px;
    opacity: 0;
    visibility: hidden;
    transform: translateY(10px);
    transition: all 0.3s ease;
    z-index: 1000;
    padding: 8px 0;
}

.dropdown:hover .dropdown-menu {
    opacity: 1;
    visibility: visible;
    transform: translateY(0);
}

.dropdown-menu li {
    padding: 0;
}

.dropdown-menu a {
    padding: 10px 20px;
    display: block;
    color: #FAFAFA;
    font-size: 0.9rem;
    transition: all 0.2s ease;
}

.dropdown-menu a:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #cc0000;
}

/* 修改下拉菜单箭头的动画样式 */
.main-menu__nav .dropdown i {
    font-size: 0.8rem;
    margin-left: 4px;
    transition: transform 0.2s ease;
}

.main-menu__nav .dropdown:hover i {
    transform: rotate(-180deg);
}

/* 用户下拉菜单样式 */
.dropdown-container .dropdown ul li a i {
    transform: none !important;
    /* 防止用户菜单图标旋转 */
    transition: none !important;
    margin-right: 8px;
    width: 16px;
    text-align: center;
}

/* 移除之前影响所有图标的样式 */
.dropdown i {
    font-size: 0.8rem;
    margin-left: 4px;
}

.dropdown:hover i {
    transform: none;
    /* 移除通用的hover效果 */
}

/* 响应式调整 */
@media (max-width: 1200px) {
    .main-menu__list {
        gap: 20px;
    }
}

@media (max-width: 992px) {
    .dropdown-menu {
        position: static;
        background: transparent;
        border: none;
        padding-left: 20px;
    }

    .coming-soon {
        display: inline-block;
        margin-top: 4px;
    }
}

/* 添加禁用状态样式 */
.disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.disabled a {
    cursor: not-allowed;
}

.disabled:hover .dropdown-menu {
    display: none;
}

.coming-soon {
    font-size: 0.7rem;
    background: rgba(204, 0, 0, 0.2);
    color: #ff4444;
    padding: 2px 6px;
    border-radius: 4px;
    margin-left: 8px;
    font-weight: normal;
}

/* 修改下拉菜单样式 */
.dropdown-menu li.disabled a {
    color: rgba(255, 255, 255, 0.5);
}

.dropdown-menu li.disabled a:hover {
    background: transparent;
    color: rgba(255, 255, 255, 0.5);
    cursor: not-allowed;
}
</style>