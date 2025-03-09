<template>
    <div class="p-6">
        <!-- 头部卡片：统计信息 -->
        <div class="grid grid-cols-4 gap-6 mb-6">
            <div v-for="(stat, index) in ['最新动漫', '本周更新', '本月更新', '总动漫数']" :key="index"
                class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-500 text-sm">{{ stat }}</h3>
                    <i :class="['text-2xl',
                        index === 0 ? 'fas fa-clock text-indigo-600' :
                            index === 1 ? 'fas fa-calendar-week text-green-600' :
                                index === 2 ? 'fas fa-calendar-alt text-amber-500' :
                                    'fas fa-film text-purple-600']"></i>
                </div>
                <div class="flex items-end gap-2">
                    <span class="text-2xl font-bold">{{
                        index === 0 ? todayCount :
                            index === 1 ? weekCount :
                                index === 2 ? monthCount : totalCount
                    }}</span>
                    <span class="text-sm text-green-500">{{
                        index === 0 ? '+3.5%' :
                            index === 1 ? '+8.2%' :
                                index === 2 ? '+12.4%' : '+5.7%'
                    }}</span>
                </div>
            </div>
        </div>

        <!-- 主要内容卡片 -->
        <div class="bg-white rounded-lg shadow-sm overflow-hidden">
            <!-- 卡片头部 -->
            <div class="p-6 border-b border-gray-200">
                <div class="flex flex-wrap items-center justify-between gap-4">
                    <div class="flex items-center gap-3">
                        <h2 class="text-lg font-bold">最新上线动漫</h2>
                        <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">总数：{{ totalCount
                        }}</span>
                    </div>
                    <div class="flex items-center gap-3">
                        <!-- 时间范围选择 -->
                        <div class="relative">
                            <select v-model="timeRange"
                                class="pl-10 pr-4 py-2 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm appearance-none">
                                <option value="all">全部时间</option>
                                <option value="today">今日更新</option>
                                <option value="week">本周更新</option>
                                <option value="month">本月更新</option>
                            </select>
                            <i class="fas fa-clock absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <!-- 类型筛选 -->
                        <div class="relative">
                            <select v-model="animeType"
                                class="pl-10 pr-4 py-2 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm appearance-none">
                                <option value="">全部类型</option>
                                <option v-for="type in animeTypes" :key="type" :value="type">{{ type }}</option>
                            </select>
                            <i class="fas fa-tags absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <!-- 搜索框 -->
                        <div class="relative">
                            <input type="text" v-model="searchQuery" placeholder="搜索动漫名称..."
                                class="pl-10 pr-4 py-2 w-64 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm"
                                @keyup.enter="handleSearch">
                            <i class="fas fa-search absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <button @click="handleSearch" class="btn btn-primary">
                            <i class="fas fa-search mr-1"></i> 搜索
                        </button>

                        <!-- 视图切换 -->
                        <div class="flex rounded-lg border border-gray-200 overflow-hidden">
                            <button @click="viewMode = 'grid'"
                                class="p-2 flex items-center justify-center w-10 transition-colors"
                                :class="viewMode === 'grid' ? 'bg-indigo-50 text-indigo-600' : 'bg-white text-gray-600 hover:bg-gray-50'">
                                <i class="fas fa-th-large"></i>
                            </button>
                            <button @click="viewMode = 'list'"
                                class="p-2 flex items-center justify-center w-10 transition-colors"
                                :class="viewMode === 'list' ? 'bg-indigo-50 text-indigo-600' : 'bg-white text-gray-600 hover:bg-gray-50'">
                                <i class="fas fa-list"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 动漫列表 - 网格视图 -->
            <div v-if="viewMode === 'grid'" class="p-6" v-loading="loading">
                <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
                    <div v-for="(anime, index) in displayAnimes" :key="index" class="group">
                        <div class="aspect-[3/4] rounded-lg overflow-hidden">
                            <img :src="anime.image" :alt="anime.title"
                                class="w-full h-full object-cover transition-transform group-hover:scale-105">
                        </div>
                        <div class="mt-3">
                            <h3
                                class="font-medium text-sm mb-1 line-clamp-1 group-hover:text-indigo-600 transition-colors">
                                {{ anime.title }}</h3>
                            <p class="text-gray-500 text-xs">{{ formatTime(anime.updateTime) }} 更新</p>
                        </div>
                    </div>
                </div>

                <!-- 空状态 -->
                <div v-if="displayAnimes.length === 0 && !loading" class="py-16 text-center text-gray-500">
                    <i class="fas fa-film text-5xl mb-4 block"></i>
                    <p>没有找到符合条件的动漫</p>
                </div>
            </div>

            <!-- 动漫列表 - 列表视图 -->
            <div v-if="viewMode === 'list'" v-loading="loading">
                <div class="divide-y divide-gray-200">
                    <div v-for="(anime, index) in displayAnimes" :key="index"
                        class="p-4 hover:bg-gray-50 transition-colors flex items-center gap-4">
                        <img :src="anime.image" :alt="anime.title" class="w-16 h-20 rounded object-cover flex-shrink-0">
                        <div class="flex-1">
                            <h3 class="font-medium mb-1 hover:text-indigo-600 transition-colors">{{ anime.title }}</h3>
                            <p class="text-sm text-gray-500">{{ formatTime(anime.updateTime) }} 更新</p>
                        </div>
                        <div class="flex-shrink-0">
                            <button class="btn-sm btn-primary">
                                <i class="fas fa-play mr-1"></i> 观看
                            </button>
                        </div>
                    </div>
                </div>

                <!-- 空状态 -->
                <div v-if="displayAnimes.length === 0 && !loading" class="py-16 text-center text-gray-500">
                    <i class="fas fa-film text-5xl mb-4 block"></i>
                    <p>没有找到符合条件的动漫</p>
                </div>
            </div>

            <!-- 分页部分 -->
            <div class="p-4 border-t border-gray-200 flex justify-end">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" :total="totalCount"
                    @size-change="handleSizeChange" @current-change="handleCurrentChange" background
                    class="custom-pagination" />
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'NewAnime',

    data() {
        return {
            searchQuery: '',
            timeRange: 'all',
            animeType: '',
            currentPage: 1,
            pageSize: 20,
            totalCount: 0,
            loading: false,
            animes: [],
            viewMode: 'grid', // 'grid' 或 'list'
            animeTypes: ['TV动画', 'OVA', '剧场版', 'WEB动画']
        }
    },

    computed: {
        // 今日更新数量
        todayCount() {
            // 获取今天的日期（yyyy-MM-dd格式）
            const today = new Date().toISOString().split('T')[0];
            return this.animes.filter(anime => anime.updateTime.startsWith(today)).length || 5;
        },

        // 本周更新数量
        weekCount() {
            const now = new Date();
            const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
            return this.animes.filter(anime => new Date(anime.updateTime) >= oneWeekAgo).length || 18;
        },

        // 本月更新数量
        monthCount() {
            const now = new Date();
            const oneMonthAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
            return this.animes.filter(anime => new Date(anime.updateTime) >= oneMonthAgo).length || 42;
        },

        // 当前显示的动漫
        displayAnimes() {
            return this.animes;
        }
    },

    methods: {
        async fetchAnimes() {
            try {
                this.loading = true;
                const response = await axios.get('/api/newAnime', {
                    params: {
                        page: this.currentPage,
                        page_size: this.pageSize,
                        time_range: this.timeRange,
                        type: this.animeType || undefined,
                        search: this.searchQuery || undefined
                    }
                });

                if (response.data.code === 200) {
                    this.animes = response.data.animes;
                    this.totalCount = response.data.total || response.data.animes.length;
                } else {
                    throw new Error(response.data.message || '获取动漫列表失败');
                }
            } catch (error) {
                ElMessage.error('获取动漫列表失败：' + error);
                // 可以设置一些默认数据，避免页面为空
                this.animes = [];
            } finally {
                this.loading = false;
            }
        },

        handleSearch() {
            this.currentPage = 1;
            this.fetchAnimes();
        },

        handleSizeChange(val) {
            this.pageSize = val;
            this.fetchAnimes();
        },

        handleCurrentChange(val) {
            this.currentPage = val;
            this.fetchAnimes();
        },

        // 将ISO时间格式转换为相对时间
        formatTime(isoTime) {
            const now = new Date();
            const animeTime = new Date(isoTime);
            const diff = Math.floor((now.getTime() - animeTime.getTime()) / 1000); // 差异（秒）

            if (diff < 60) {
                return '刚刚';
            } else if (diff < 3600) {
                return Math.floor(diff / 60) + '分钟前';
            } else if (diff < 86400) {
                return Math.floor(diff / 3600) + '小时前';
            } else if (diff < 2592000) {
                return Math.floor(diff / 86400) + '天前';
            } else {
                // 对于较早的更新，返回具体日期
                return animeTime.toLocaleDateString();
            }
        }
    },

    created() {
        this.fetchAnimes();
    }
}
</script>

<style scoped>
/* 自定义分页器样式 */
.custom-pagination :deep(.el-pagination__sizes .el-select .el-input .el-input__inner) {
    @apply border-gray-300 rounded;
}

.custom-pagination :deep(.el-pagination button) {
    @apply bg-white rounded border border-gray-300 mx-1;
}

.custom-pagination :deep(.el-pagination button:hover) {
    @apply bg-gray-50;
}

.custom-pagination :deep(.el-pagination button.is-active) {
    @apply bg-indigo-600 text-white border-indigo-600;
}

/* 统一按钮样式 */
.btn {
    @apply inline-flex items-center justify-center px-4 py-2 rounded-lg text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-sm {
    @apply inline-flex items-center justify-center px-2.5 py-1.5 rounded text-xs font-medium transition-colors focus:outline-none focus:ring-1 focus:ring-offset-1;
}

.btn-primary {
    @apply bg-indigo-600 hover:bg-indigo-700 text-white focus:ring-indigo-500;
}

/* 内容限制行数 */
.line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>
