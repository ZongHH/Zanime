<template>
    <!-- Back To Top -->
    <button class="scrollToTopBtn">
        <i class="fa fa-arrow-up"></i>
    </button>

    <!-- Mobile Menu Start -->
    <div class="mobile-nav__wrapper">
        <div class="mobile-nav__overlay mobile-nav__toggler"></div>
        <div class="mobile-nav__content">
            <span class="mobile-nav__close mobile-nav__toggler"><i class="fa fa-times"></i></span>
            <div class="logo-box">
                <a href="" aria-label="logo image"><img src="@/static/picture/logo.png" alt="logo"></a>
            </div>
            <!-- 添加移动端搜索栏 -->
            <div class="mobile-search-box">
                <form @submit.prevent="handleSearch">
                    <input type="text" v-model="searchQuery" placeholder="搜索..." @input="handleSearchInput"
                        @focus="showDropdown = true" @blur="handleBlur" required>
                    <button class="search" type="submit">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewbox="0 0 20 20" fill="none">
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
                <!-- 搜索结果下拉框 -->
                <div class="search-dropdown" v-if="showDropdown">
                    <ul v-if="searchResults.length > 0">
                        <li v-for="result in searchResults" :key="result.id" @mousedown="handleResultClick(result)">
                            <div class="search-result-item">
                                <img :src="result.cover_image_url" :alt="result.video_name" class="result-thumbnail">
                                <div class="result-info">
                                    <div class="result-title">{{ result.video_name }}</div>
                                    <div class="result-type">{{ result.anime_type }}</div>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="mobile-nav__container"></div>
            <ul class="mobile-nav__contact list-unstyled">
                <li>
                    <i class="fa-thin fa-envelope"></i>
                    <a href="mailto:example@company.com">example@company.com</a>
                </li>
                <li>
                    <i class="fa-light fa-phone-volume"></i>
                    <a href="tel:+12345678">+123 (4567) -890</a>
                </li>
            </ul>
        </div>
    </div>
    <!-- Mobile Menu End -->
</template>

<script>
import axios from 'axios';

export default {
    name: "MobileMenu",
    data() {
        return {
            searchQuery: '',
            searchResults: [],
            showDropdown: false,
            searchTimeout: null,
        }
    },
    methods: {
        handleSearchInput() {
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
                this.searchResults = response.data.data;
            } catch (error) {
                console.error('搜索出错:', error);
                this.searchResults = [];
            }
        },

        handleResultClick(result) {
            this.$router.push(`/moviesDetail?videoId=${result.video_id}`);
        },

        handleSearch(e) {
            e.preventDefault();
            if (this.searchQuery.trim()) {
                this.$router.push(`/animeDetail?params=${encodeURIComponent(this.searchQuery.trim())}`);
            }
        },

        handleBlur() {
            setTimeout(() => {
                this.showDropdown = false;
            }, 200);
        },
    }
}
</script>

<style scoped>
.mobile-search-box {
    padding: 15px;
    position: relative;
}

.mobile-search-box form {
    position: relative;
    display: flex;
    align-items: center;
}

.mobile-search-box input {
    width: 100%;
    height: 40px;
    background: #1A1A1A;
    border: 1px solid #333;
    border-radius: 20px;
    padding: 0 40px 0 15px;
    color: #FAFAFA;
    font-size: 14px;
}

.mobile-search-box .search {
    position: absolute;
    right: 15px;
    background: none;
    border: none;
    cursor: pointer;
}

.search-dropdown {
    position: absolute;
    top: 100%;
    left: 15px;
    right: 15px;
    background: #1A1A1A;
    border: 1px solid #333;
    border-radius: 4px;
    margin-top: 4px;
    max-height: 300px;
    overflow-y: auto;
    z-index: 1000;
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

.search-dropdown ul {
    list-style: none;
    padding: 0;
    margin: 0;
}
</style>