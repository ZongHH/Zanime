<template>
    <div class="p-6">
        <!-- 头部卡片：统计信息 -->
        <div class="grid grid-cols-4 gap-6 mb-6">
            <div v-for="(stat, index) in ['总日志数', '管理员操作', '用户操作', '体验用户']" :key="index"
                class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-500 text-sm">{{ stat }}</h3>
                    <i :class="['text-2xl',
                        index === 0 ? 'fas fa-clipboard-list text-indigo-600' :
                            index === 1 ? 'fas fa-user-shield text-purple-600' :
                                index === 2 ? 'fas fa-user text-blue-600' :
                                    'fas fa-user-tag text-orange-500']"></i>
                </div>
                <div class="flex items-end gap-2">
                    <span class="text-2xl font-bold">{{
                        index === 0 ? total :
                            index === 1 ? adminCount :
                                index === 2 ? regularCount : trialCount
                    }}</span>
                    <span :class="['text-sm', 'text-green-500']">{{
                        index === 0 ? '+8.3%' :
                            index === 1 ? '+5.7%' :
                                index === 2 ? '+12.4%' : '+9.2%'
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
                        <h2 class="text-lg font-bold">操作日志</h2>
                        <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">总数：{{ total }}</span>
                    </div>
                    <div class="flex items-center gap-3">
                        <!-- 筛选项 -->
                        <div class="relative">
                            <select v-model="filterType"
                                class="pl-10 pr-4 py-2 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm appearance-none">
                                <option value="">全部类型</option>
                                <option value="admin">管理员</option>
                                <option value="regular">普通用户</option>
                                <option value="trial">体验用户</option>
                            </select>
                            <i class="fas fa-filter absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <div class="relative">
                            <select v-model="filterModule"
                                class="pl-10 pr-4 py-2 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm appearance-none">
                                <option value="">全部模块</option>
                                <option v-for="module in modules" :key="module" :value="module">{{ module }}</option>
                            </select>
                            <i class="fas fa-th-large absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <div class="relative">
                            <input type="text" v-model="searchQuery" placeholder="搜索日志内容..."
                                class="pl-10 pr-4 py-2 w-64 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm"
                                @keyup.enter="handleSearch">
                            <i class="fas fa-search absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>

                        <button @click="handleSearch" class="btn btn-primary">
                            <i class="fas fa-search mr-1"></i> 搜索
                        </button>

                        <button @click="handleExport" class="btn btn-success">
                            <i class="fas fa-file-export mr-1"></i> 导出
                        </button>
                    </div>
                </div>
            </div>

            <!-- 日志列表部分 -->
            <div class="divide-y divide-gray-200" v-loading="loading">
                <div v-for="(log, index) in displayLogs" :key="index" class="p-4 hover:bg-gray-50 transition-colors">
                    <div class="flex items-start gap-4">
                        <div class="w-10 h-10 rounded-full flex items-center justify-center" :class="{
                            'bg-purple-100': log.user_type === 'admin',
                            'bg-blue-100': log.user_type === 'regular',
                            'bg-orange-100': log.user_type === 'trial'
                        }">
                            <i :class="{
                                'fas fa-user-shield text-purple-600': log.user_type === 'admin',
                                'fas fa-user text-blue-600': log.user_type === 'regular',
                                'fas fa-user-tag text-orange-500': log.user_type === 'trial'
                            }"></i>
                        </div>
                        <div class="flex-1">
                            <div class="flex items-center justify-between">
                                <h3 class="font-medium text-sm">{{ log.user_name }}</h3>
                                <span class="text-xs text-gray-400">{{ formatTime(log.time) }}</span>
                            </div>
                            <p class="text-sm text-gray-600 mt-1">{{ log.action }}</p>
                            <div class="mt-2 flex items-center justify-between">
                                <div class="text-xs text-gray-400">
                                    <i class="fas fa-tag mr-1"></i>
                                    {{ log.module }}
                                </div>
                                <button @click="handleLogDetail(log)"
                                    class="text-xs text-indigo-600 hover:text-indigo-800">
                                    <i class="fas fa-info-circle mr-1"></i> 详细信息
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 空状态 -->
                <div v-if="displayLogs.length === 0 && !loading" class="p-16 text-center text-gray-500">
                    <i class="fas fa-clipboard-list text-4xl mb-4 block"></i>
                    <p>没有找到符合条件的操作日志</p>
                </div>
            </div>

            <!-- 分页部分 -->
            <div class="p-4 border-t border-gray-200 flex justify-end">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" :total="totalRecords"
                    @size-change="handleSizeChange" @current-change="handleCurrentChange" background
                    class="custom-pagination" />
            </div>
        </div>

        <!-- 日志详情弹窗 -->
        <el-dialog v-model="detailVisible" title="操作日志详情" width="600px" destroy-on-close>
            <div v-if="selectedLog" class="p-4">
                <div class="flex items-center gap-4 mb-6">
                    <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="{
                        'bg-purple-100': selectedLog.user_type === 'admin',
                        'bg-blue-100': selectedLog.user_type === 'regular',
                        'bg-orange-100': selectedLog.user_type === 'trial'
                    }">
                        <i :class="{
                            'fas fa-user-shield text-purple-600': selectedLog.user_type === 'admin',
                            'fas fa-user text-blue-600': selectedLog.user_type === 'regular',
                            'fas fa-user-tag text-orange-500': selectedLog.user_type === 'trial'
                        }"></i>
                    </div>
                    <div>
                        <h3 class="text-lg font-bold">{{ selectedLog.user_name }}</h3>
                        <p class="text-sm text-gray-500">{{
                            selectedLog.user_type === 'admin' ? '系统管理员' :
                                selectedLog.user_type === 'regular' ? '普通用户' :
                                    '体验用户'
                        }}</p>
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-500 mb-1">操作时间</p>
                        <p class="text-sm">{{ formatFullTime(selectedLog.time) }}</p>
                    </div>
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-500 mb-1">操作模块</p>
                        <p class="text-sm">{{ selectedLog.module }}</p>
                    </div>
                </div>

                <div class="mb-6">
                    <p class="text-xs text-gray-500 mb-1">操作内容</p>
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-sm">{{ selectedLog.action }}</p>
                    </div>
                </div>

                <div class="mb-6">
                    <p class="text-xs text-gray-500 mb-1">操作ID</p>
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-sm">{{ selectedLog.id }}</p>
                    </div>
                </div>

                <div class="flex justify-end">
                    <button @click="detailVisible = false" class="btn btn-gray">
                        <i class="fas fa-times mr-2"></i>关闭
                    </button>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'OperationLogs',

    data() {
        return {
            searchQuery: '',
            filterType: '',
            filterModule: '',
            currentPage: 1,
            pageSize: 10,
            totalRecords: 0,
            detailVisible: false,
            selectedLog: null,
            loading: false,
            logs: []
        }
    },

    computed: {
        // 获取所有日志的数量
        total() {
            return this.totalRecords;
        },

        // 管理员操作数量
        adminCount() {
            return this.logs.filter(log => log.user_type === 'admin').length;
        },

        // 普通用户操作数量
        regularCount() {
            return this.logs.filter(log => log.user_type === 'regular').length;
        },

        // 体验用户操作数量
        trialCount() {
            return this.logs.filter(log => log.user_type === 'trial').length;
        },

        // 所有模块列表
        modules() {
            const moduleSet = new Set();
            this.logs.forEach(log => moduleSet.add(log.module));
            return Array.from(moduleSet);
        },

        // 当前显示的日志
        displayLogs() {
            return this.logs;
        }
    },

    methods: {
        async fetchLogs() {
            try {
                this.loading = true;
                const response = await axios.get('/api/userActionLogs', {
                    params: {
                        page: this.currentPage,
                        page_size: this.pageSize,
                        user_type: this.filterType || undefined,
                        module: this.filterModule || undefined,
                        search: this.searchQuery || undefined
                    }
                });

                if (response.data.code === 200) {
                    this.logs = response.data.userActionLogs;
                    this.totalRecords = response.data.total || response.data.userActionLogs.length;
                } else {
                    throw new Error(response.data.message || '获取操作日志失败');
                }
            } catch (error) {
                ElMessage.error('获取操作日志失败：' + error);
            } finally {
                this.loading = false;
            }
        },

        handleSearch() {
            this.currentPage = 1;
            this.fetchLogs();
        },

        handleExport() {
            // 生成导出文件名
            const timestamp = new Date().toISOString().replace(/:/g, '-').substring(0, 19);
            const fileName = `操作日志_${timestamp}.csv`;

            // 构建CSV内容
            let csvContent = "ID,用户名,用户类型,操作内容,操作时间,模块\n";

            this.logs.forEach(log => {
                const row = [
                    log.id,
                    log.user_name,
                    log.user_type === 'admin' ? '管理员' :
                        log.user_type === 'regular' ? '普通用户' : '体验用户',
                    log.action.replace(/,/g, '，'), // 替换英文逗号避免CSV格式问题
                    this.formatFullTime(log.time),
                    log.module
                ];
                csvContent += row.join(',') + "\n";
            });

            // 创建下载链接
            const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.setAttribute('href', url);
            link.setAttribute('download', fileName);
            link.style.visibility = 'hidden';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            ElMessage({
                message: '日志导出成功',
                type: 'success'
            });
        },

        handleLogDetail(log) {
            this.selectedLog = log;
            this.detailVisible = true;
        },

        handleSizeChange(val) {
            this.pageSize = val;
            this.fetchLogs();
        },

        handleCurrentChange(val) {
            this.currentPage = val;
            this.fetchLogs();
        },

        // 将ISO时间格式转换为相对时间
        formatTime(isoTime) {
            const now = new Date();
            const logTime = new Date(isoTime);
            const diff = Math.floor((now.getTime() - logTime.getTime()) / 1000); // 差异（秒）

            if (diff < 60) {
                return '刚刚';
            } else if (diff < 3600) {
                return Math.floor(diff / 60) + '分钟前';
            } else if (diff < 86400) {
                return Math.floor(diff / 3600) + '小时前';
            } else if (diff < 2592000) {
                return Math.floor(diff / 86400) + '天前';
            } else {
                return logTime.toLocaleDateString();
            }
        },

        // 将ISO时间格式转换为完整时间
        formatFullTime(isoTime) {
            const date = new Date(isoTime);
            return date.toLocaleString();
        }
    },

    created() {
        this.fetchLogs();
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

.btn-primary {
    @apply bg-indigo-600 hover:bg-indigo-700 text-white focus:ring-indigo-500;
}

.btn-success {
    @apply bg-green-600 hover:bg-green-700 text-white focus:ring-green-500;
}

.btn-gray {
    @apply bg-gray-200 hover:bg-gray-300 text-gray-800 focus:ring-gray-500;
}

/* 下拉选择框样式 */
select {
    min-width: 160px;
}
</style>