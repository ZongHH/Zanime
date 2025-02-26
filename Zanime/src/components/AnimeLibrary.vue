<template>
    <div class="anime-library">
        <!-- 筛选区域 -->
        <section class="filter-section" v-if="filters">
            <div class="container-fluid">
                <div class="filter-container">
                    <!-- 修改筛选组件 -->
                    <div class="filter-group" v-for="(filterItems, type) in filterGroups" :key="type">
                        <h6 class="filter-title">{{ filterTitles[type] }}</h6>
                        <div class="filter-options" :class="{ 'collapsed': !expanded[type] }">
                            <span v-for="item in displayedItems(type)" :key="item"
                                :class="['filter-tag', { active: selectedFilters[type] === item }]"
                                @click="updateFilter(type, item)">
                                {{ item }}
                            </span>
                        </div>
                        <button v-if="needsExpansion(type)" class="expand-btn" @click="toggleExpand(type)">
                            {{ expanded[type] ? '收起' : '更多' }}
                            <i :class="['fas', expanded[type] ? 'fa-chevron-up' : 'fa-chevron-down']"></i>
                        </button>
                    </div>

                    <!-- 重置按钮 -->
                    <button class="reset-btn" @click="resetFilters">
                        重置筛选
                        <i class="fas fa-redo-alt"></i>
                    </button>
                </div>
            </div>
        </section>

        <!-- 动漫列表 -->
        <section class="anime-grid-section">
            <div class="container-fluid">
                <!-- 动漫网格 -->
                <div class="anime-grid">
                    <div v-for="anime in filteredAnimes" :key="anime.video_id" class="anime-card"
                        @click="navigateToDetail(anime)">
                        <div class="card-image">
                            <img :src="anime.cover_image_url" :alt="anime.video_name">
                            <div class="card-overlay">
                                <div class="card-rating">
                                    <i class="fas fa-star"></i>
                                    <span>{{ anime.rating }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="card-info">
                            <h6>{{ anime.video_name }}</h6>
                            <div class="card-meta">
                                <span>{{ anime.release_date }}</span>
                                <span>{{ anime.genres }}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 替换加载更多为分页组件 -->
                <div class="pagination-container">
                    <div class="pagination">
                        <button class="page-btn" :disabled="currentPage === 1" @click="changePage(currentPage - 1)">
                            <i class="fas fa-chevron-left"></i>
                        </button>

                        <button v-for="pageNum in displayedPages" :key="pageNum"
                            :class="['page-btn', { active: currentPage === pageNum }]" @click="changePage(pageNum)">
                            {{ pageNum }}
                        </button>

                        <button class="page-btn" :disabled="currentPage === totalPages"
                            @click="changePage(currentPage + 1)">
                            <i class="fas fa-chevron-right"></i>
                        </button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<script>
export default {
    data() {
        return {
            filters: {},
            selectedFilters: {
                region: '',
                year: '',
                type: '',
                letter: ''
            },
            animes: [],
            loading: false,
            currentPage: 1,
            totalPages: 1,
            pageSize: 24,
            expanded: {
                region: false,
                year: false,
                type: false,
                letter: false
            },
            filterTitles: {
                region: '地区',
                year: '年份',
                type: '类型',
                letter: '首字母'
            },
            itemsPerRow: 8, // 每行显示的标签数量
            defaultRows: 2, // 默认显示行数
        };
    },
    computed: {
        totalAnime() {
            return this.filteredAnimes.length;
        },
        filteredAnimes() {
            // 实现筛选逻辑
            return this.animes;
        },
        displayedPages() {
            const delta = 2; // 当前页码左右显示的页码数
            let pages = [];
            let left = Math.max(1, this.currentPage - delta);
            let right = Math.min(this.totalPages, this.currentPage + delta);

            // 处理省略号显示
            if (left > 2) {
                pages.push(1);
                pages.push('...');
            }

            for (let i = left; i <= right; i++) {
                pages.push(i);
            }

            if (right < this.totalPages - 1) {
                pages.push('...');
                pages.push(this.totalPages);
            }

            return pages;
        },
        filterGroups() {
            return {
                region: this.filters?.regions || [],
                year: this.filters?.years || [],
                type: this.filters?.types || [],
                letter: this.filters?.letters || []
            }
        }
    },
    methods: {
        updateFilter(type, value) {
            this.selectedFilters[type] = value;
            this.currentPage = 1
            this.fetchAnimes();
        },
        resetFilters() {
            this.currentPage = 1
            this.selectedFilters = {
                region: '',
                year: '',
                type: '',
                letter: ''
            };
            this.fetchAnimes();
        },
        async fetchFilters() {
            try {
                const response = await axios.get('/api/animeFilters');
                if (response.data.code == 200) {
                    this.filters = response.data.video_filters;
                } else {
                    console.error("获取动漫筛选信息失败: ", response.data.data)
                }
            } catch (error) {
                console.error('获取动漫筛选信息失败:', error);
            }
        },
        async fetchAnimes() {
            this.loading = true;
            try {
                const response = await axios.get('/api/animeLibrary', {
                    params: {
                        ...this.selectedFilters,
                        page: this.currentPage,
                        page_size: this.pageSize
                    }
                });
                if (response.data.code == 200) {
                    this.animes = response.data.videos;
                    this.totalPages = Math.ceil(response.data.total / this.pageSize);
                } else {
                    console.error("获取动漫列表失败: response.data.code == ", response.data.code)
                }
            } catch (error) {
                console.error('获取动漫列表失败:', error);
            } finally {
                this.loading = false;
            }
        },
        navigateToDetail(anime) {
            this.$router.push(`/moviesDetail?videoId=${anime.video_id}`)
        },
        changePage(page) {
            if (page === '...') return;
            this.currentPage = page;
            this.fetchAnimes();
            // 滚动到顶部
            window.scrollTo({ top: 0, behavior: 'smooth' });
        },
        displayedItems(type) {
            const items = this.filterGroups[type];
            if (!items) return [];

            if (this.expanded[type]) {
                return items;
            }

            return items.slice(0, this.itemsPerRow * this.defaultRows);
        },

        needsExpansion(type) {
            return this.filterGroups[type]?.length > this.itemsPerRow * this.defaultRows;
        },

        toggleExpand(type) {
            this.expanded[type] = !this.expanded[type];
        },
    },
    mounted() {
        this.fetchFilters()
        this.fetchAnimes();
    },
};
import axios from 'axios';
</script>

<style scoped>
@import "@/static/css/animeLibrary.css";
</style>
