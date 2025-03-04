<template>
    <!-- HEADER MENU START -->
    <header class="header">
        <nav class="navigation d-flex align-items-center justify-content-between">
            <router-link to="/" class="d-flex align-items-center">
                <img src="@/static/picture/logo.png" alt="/logo" class="header-logo">
            </router-link>
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
                        <li class="dropdown">
                            <router-link to="/feedback">
                                反馈中心
                                <span class="beta-badge">Beta</span>
                                <i class="fas fa-chevron-down"></i>
                            </router-link>
                            <ul class="dropdown-menu">
                                <li><router-link to="/feedback?type=suggestion">反馈建议</router-link></li>
                                <li><router-link to="/feedback?type=bug">反馈BUG</router-link></li>
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
                        <span class="tooltip-text search-tooltip">搜索</span>
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
                        <span class="tooltip-text notification-tooltip">通知</span>
                    </a>
                    <Notify :showNotifications="showNotifications" />
                </div>
                <div class="dropdown-container">
                    <details class="dropdown right">
                        <summary class="avatar">
                            <img :src="avatarURL" alt="logo">
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
                                <a href="#" @click.prevent="showHelpModal">
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
        <Help :isVisible="showHelp" @close="showHelp = false" />
    </header>
    <!-- HEADER MENU END -->
</template>

<script>
import Notify from '@/components/general/Notify.vue';
import Help from '@/components/general/Help.vue';
import Logout from '@/static/js/general.js';
import axios from 'axios';

export default {
    data() {
        return {
            showNotifications: false,
            showHelp: false,
            searchQuery: '',
            searchResults: [],
            showDropdown: false,
            searchTimeout: null,
            activeMenuItem: '', // 默认不激活
        };
    },
    computed: {
        avatarURL() {
            return localStorage.getItem('avatar_url') || 'src/static/picture/Ellipse-1.png';
        }
    },
    methods: {
        /**
         * 切换通知面板的显示状态
         * 当用户点击通知图标时调用此方法
         * 通过取反当前状态来实现显示/隐藏切换
         */
        toggleNotifications() {
            this.showNotifications = !this.showNotifications;
        },

        /**
         * 显示帮助模态框
         * 当用户点击帮助选项时调用此方法
         * 将showHelp状态设置为true以显示帮助组件
         */
        showHelpModal() {
            this.showHelp = true;
        },

        /**
         * 处理搜索输入事件
         * 实现防抖功能，避免频繁API调用
         * 当用户输入内容时，延迟300ms后才执行搜索
         * 如果搜索框为空，则清空搜索结果
         */
        handleSearchInput() {
            // 防抖处理：清除之前的定时器
            clearTimeout(this.searchTimeout);
            if (this.searchQuery.trim()) {
                // 设置新的定时器，延迟300ms执行搜索
                this.searchTimeout = setTimeout(() => {
                    this.fetchSearchResults();
                }, 300);
            } else {
                // 搜索框为空时清空结果
                this.searchResults = [];
            }
        },

        /**
         * 获取搜索结果
         * 异步方法，通过API获取匹配的动画列表
         * 成功时更新searchResults数组
         * 失败时记录错误并清空结果
         */
        async fetchSearchResults() {
            try {
                // 发送GET请求到搜索API
                const response = await axios.get(`/api/search?query=${this.searchQuery}`);
                if (response.data.code === 200) {
                    // 请求成功，更新搜索结果
                    this.searchResults = response.data.animes;
                } else {
                    // 服务器返回错误码
                    throw new Error(response.data.message)
                }
            } catch (error) {
                // 捕获并记录错误
                console.error('搜索出错:', error);
                this.searchResults = [];
            }
        },

        /**
         * 处理搜索结果点击事件
         * 当用户点击某个搜索结果项时调用
         * 跳转到对应的详情页面
         * @param {Object} result - 包含视频信息的结果对象
         */
        handleResultClick(result) {
            // 使用vue-router导航到详情页，传递视频ID作为参数
            this.$router.push(`/moviesDetail?videoId=${result.video_id}`);
        },

        /**
         * 处理搜索表单提交
         * 当用户提交搜索表单时调用
         * 阻止默认表单提交行为，并跳转到搜索结果页
         * @param {Event} e - 表单提交事件对象
         */
        handleSearch(e) {
            // 阻止表单默认提交行为
            e.preventDefault();
            // 处理搜索表单提交
            if (this.searchQuery.trim()) {
                // 注释掉的代码：this.fetchSearchResults();
                // 跳转到动画详情页，将搜索词作为参数传递
                this.$router.push(`/animeDetail?params=${encodeURIComponent(this.searchQuery.trim())}`);
            }
        },

        /**
         * 处理搜索框失焦事件
         * 当搜索框失去焦点时，延迟关闭下拉框
         * 延迟是为了确保点击事件能够被处理
         */
        handleBlur() {
            // 延迟200ms关闭下拉框，以便能够处理点击事件
            setTimeout(() => {
                this.showDropdown = false;
            }, 200);
        },

        /**
         * 跳转到订单页面
         * 当用户点击"我的订单"选项时调用
         * 使用vue-router导航到订单页面
         */
        JumpOrder() {
            this.$router.push("/orders")
        },

        /**
         * 跳转到个人中心页面
         * 当用户点击个人中心选项时调用
         * 使用vue-router导航到个人中心页面
         */
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
        Help,
    },
};
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

/* 添加提示词样式 */
.input-box form,
.icon a {
    position: relative;
}

.tooltip-text {
    position: absolute;
    bottom: -25px;
    left: 50%;
    transform: translateX(-50%);
    background-color: rgba(26, 26, 26, 0.9);
    color: #fff;
    font-size: 0.75rem;
    padding: 3px 8px;
    border-radius: 4px;
    white-space: nowrap;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    z-index: 1001;
}

.search-tooltip {
    bottom: -50px;
}

.notification-tooltip {
    bottom: -40px;
}

.input-box form:hover .tooltip-text,
.icon a:hover .tooltip-text {
    opacity: 1;
}

.tooltip-text::after {
    content: '';
    position: absolute;
    bottom: 100%;
    left: 50%;
    margin-left: -4px;
    border-width: 4px;
    border-style: solid;
    border-color: transparent transparent rgba(26, 26, 26, 0.9) transparent;
}

@media (max-width: 768px) {
    .tooltip-text {
        display: none;
    }
}

/* 添加Beta标签样式 */
.beta-badge {
    display: inline-block;
    font-size: 0.65rem;
    background: #833232;
    color: white;
    padding: 1px 5px;
    border-radius: 4px;
    font-weight: 500;
    vertical-align: middle;
    position: relative;
    top: 1px;
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
</style>