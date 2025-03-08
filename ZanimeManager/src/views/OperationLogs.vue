<template>
    <div class="p-6">
        <!-- 头部卡片：统计信息 -->
        <div class="grid grid-cols-4 gap-6 mb-6">
            <div v-for="(stat, index) in ['总日志数', '管理员操作', '用户操作', '今日操作']" :key="index"
                class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-500 text-sm">{{ stat }}</h3>
                    <i :class="['text-2xl',
                        index === 0 ? 'fas fa-clipboard-list text-indigo-600' :
                            index === 1 ? 'fas fa-user-shield text-purple-600' :
                                index === 2 ? 'fas fa-user text-blue-600' :
                                    'fas fa-calendar-day text-green-600']"></i>
                </div>
                <div class="flex items-end gap-2">
                    <span class="text-2xl font-bold">{{
                        index === 0 ? total :
                            index === 1 ? adminCount :
                                index === 2 ? userCount : todayCount
                    }}</span>
                    <span :class="['text-sm', 'text-green-500']">{{
                        index === 0 ? '+8.3%' :
                            index === 1 ? '+5.7%' :
                                index === 2 ? '+12.4%' : '+15.2%'
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
                                <option value="admin">管理员操作</option>
                                <option value="user">用户操作</option>
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
            <div class="divide-y divide-gray-200">
                <div v-for="(log, index) in displayLogs" :key="index" class="p-4 hover:bg-gray-50 transition-colors">
                    <div class="flex items-start gap-4">
                        <div class="w-10 h-10 rounded-full flex items-center justify-center"
                            :class="log.type === 'admin' ? 'bg-purple-100' : 'bg-blue-100'">
                            <i
                                :class="[log.type === 'admin' ? 'fas fa-user-shield text-purple-600' : 'fas fa-user text-blue-600']"></i>
                        </div>
                        <div class="flex-1">
                            <div class="flex items-center justify-between">
                                <h3 class="font-medium text-sm">{{ log.user }}</h3>
                                <span class="text-xs text-gray-400">{{ log.time }}</span>
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
                <div v-if="displayLogs.length === 0" class="p-16 text-center text-gray-500">
                    <i class="fas fa-clipboard-list text-4xl mb-4 block"></i>
                    <p>没有找到符合条件的操作日志</p>
                </div>
            </div>

            <!-- 分页部分 -->
            <div class="p-4 border-t border-gray-200 flex justify-end">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next"
                    :total="filteredLogs.length" @size-change="handleSizeChange" @current-change="handleCurrentChange"
                    background class="custom-pagination" />
            </div>
        </div>

        <!-- 日志详情弹窗 -->
        <el-dialog v-model="detailVisible" title="操作日志详情" width="600px" destroy-on-close>
            <div v-if="selectedLog" class="p-4">
                <div class="flex items-center gap-4 mb-6">
                    <div class="w-12 h-12 rounded-full flex items-center justify-center"
                        :class="selectedLog.type === 'admin' ? 'bg-purple-100' : 'bg-blue-100'">
                        <i
                            :class="[selectedLog.type === 'admin' ? 'fas fa-user-shield text-purple-600' : 'fas fa-user text-blue-600']"></i>
                    </div>
                    <div>
                        <h3 class="text-lg font-bold">{{ selectedLog.user }}</h3>
                        <p class="text-sm text-gray-500">{{ selectedLog.type === 'admin' ? '系统管理员' : '普通用户' }}</p>
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-500 mb-1">操作时间</p>
                        <p class="text-sm">{{ selectedLog.time }}</p>
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
                    <p class="text-xs text-gray-500 mb-1">操作IP</p>
                    <div class="bg-gray-50 p-3 rounded-lg">
                        <p class="text-sm">{{ selectedLog.ip || '192.168.1.' + Math.floor(Math.random() * 255) }}</p>
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
import { ElMessage } from 'element-plus'

export default {
    name: 'OperationLogs',

    data() {
        return {
            searchQuery: '',
            filterType: '',
            filterModule: '',
            currentPage: 1,
            pageSize: 20,
            detailVisible: false,
            selectedLog: null,
            loading: false,
            logs: [
                {
                    id: 1,
                    user: '系统管理员',
                    type: 'admin',
                    action: '更新了《咒术回战 第二季》的播放源信息',
                    time: '2024-04-15 10:30:25',
                    module: '内容管理',
                    ip: '192.168.1.100'
                },
                {
                    id: 2,
                    user: '王梓晨',
                    type: 'user',
                    action: '发表了《进击的巨人》的评论',
                    time: '2024-04-15 09:45:18',
                    module: '用户互动',
                    ip: '192.168.1.101'
                },
                {
                    id: 3,
                    user: '系统管理员',
                    type: 'admin',
                    action: '删除了违规用户评论',
                    time: '2024-04-15 09:12:33',
                    module: '内容审核',
                    ip: '192.168.1.100'
                },
                {
                    id: 4,
                    user: '林雨晴',
                    type: 'user',
                    action: '收藏了《间谍过家家》',
                    time: '2024-04-15 08:55:21',
                    module: '用户行为',
                    ip: '192.168.1.102'
                },
                {
                    id: 5,
                    user: '系统管理员',
                    type: 'admin',
                    action: '新增了《海贼王》最新话章节',
                    time: '2024-04-15 08:30:15',
                    module: '内容管理',
                    ip: '192.168.1.100'
                },
                {
                    id: 6,
                    user: '张伟',
                    type: 'user',
                    action: '评分了《鬼灭之刃》',
                    time: '2024-04-14 22:15:42',
                    module: '用户评分',
                    ip: '192.168.1.103'
                },
                {
                    id: 7,
                    user: '系统管理员',
                    type: 'admin',
                    action: '更新了《鬼灭之刃》的剧集信息',
                    time: '2024-04-14 18:20:11',
                    module: '内容管理',
                    ip: '192.168.1.100'
                },
                {
                    id: 8,
                    user: '李明',
                    type: 'user',
                    action: '分享了《间谍过家家》到微信',
                    time: '2024-04-14 17:05:33',
                    module: '社交分享',
                    ip: '192.168.1.104'
                },
                {
                    id: 9,
                    user: '系统管理员',
                    type: 'admin',
                    action: '封禁了违规用户账号',
                    time: '2024-04-14 16:40:27',
                    module: '用户管理',
                    ip: '192.168.1.100'
                },
                {
                    id: 10,
                    user: '赵芳',
                    type: 'user',
                    action: '订阅了《海贼王》的更新提醒',
                    time: '2024-04-14 15:22:38',
                    module: '订阅管理',
                    ip: '192.168.1.105'
                },
                {
                    id: 11,
                    user: '系统管理员',
                    type: 'admin',
                    action: '发布了系统公告',
                    time: '2024-04-14 14:10:55',
                    module: '系统管理',
                    ip: '192.168.1.100'
                },
                {
                    id: 12,
                    user: '陈晓',
                    type: 'user',
                    action: '上传了头像',
                    time: '2024-04-14 13:05:22',
                    module: '用户资料',
                    ip: '192.168.1.106'
                }
            ]
        }
    },

    computed: {
        // 总日志数
        total() {
            return this.logs.length;
        },

        // 管理员操作数量
        adminCount() {
            return this.logs.filter(log => log.type === 'admin').length;
        },

        // 用户操作数量
        userCount() {
            return this.logs.filter(log => log.type === 'user').length;
        },

        // 今日操作数量（假设是时间包含今天日期的）
        todayCount() {
            const today = new Date().toISOString().split('T')[0]; // 获取今天的日期部分
            return this.logs.filter(log => log.time.includes(today)).length || 5; // 如果没有今天的数据，显示5
        },

        // 所有模块列表
        modules() {
            const moduleSet = new Set();
            this.logs.forEach(log => moduleSet.add(log.module));
            return Array.from(moduleSet);
        },

        // 筛选后的日志
        filteredLogs() {
            let result = [...this.logs];

            // 按类型筛选
            if (this.filterType) {
                result = result.filter(log => log.type === this.filterType);
            }

            // 按模块筛选
            if (this.filterModule) {
                result = result.filter(log => log.module === this.filterModule);
            }

            // 搜索查询
            if (this.searchQuery) {
                const query = this.searchQuery.toLowerCase();
                result = result.filter(log =>
                    log.user.toLowerCase().includes(query) ||
                    log.action.toLowerCase().includes(query) ||
                    log.module.toLowerCase().includes(query)
                );
            }

            return result;
        },

        // 当前页显示的日志
        displayLogs() {
            const start = (this.currentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredLogs.slice(start, end);
        }
    },

    methods: {
        handleSearch() {
            this.currentPage = 1;
            ElMessage({
                message: '搜索完成',
                type: 'success'
            });
        },

        handleExport() {
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
        },

        handleCurrentChange(val) {
            this.currentPage = val;
        }
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