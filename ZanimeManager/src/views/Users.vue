<template>
    <div class="p-6">
        <!-- 头部卡片：统计信息 -->
        <div class="grid grid-cols-4 gap-6 mb-6">
            <div v-for="(stat, index) in ['总用户数', '在线用户', '今日注册', '异常账户']" :key="index"
                class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-500 text-sm">{{ stat }}</h3>
                    <i :class="['text-2xl',
                        index === 0 ? 'fas fa-users text-indigo-600' :
                            index === 1 ? 'fas fa-user-check text-green-600' :
                                index === 2 ? 'fas fa-user-plus text-amber-500' :
                                    'fas fa-user-shield text-red-600']"></i>
                </div>
                <div class="flex items-end gap-2">
                    <span class="text-2xl font-bold">{{
                        index === 0 ? total :
                            index === 1 ? '42' :
                                index === 2 ? '18' : '5'
                    }}</span>
                    <span :class="['text-sm', index === 3 ? 'text-red-500' : 'text-green-500']">{{
                        index === 0 ? '+5.2%' :
                            index === 1 ? '+12.3%' :
                                index === 2 ? '+8.7%' : '-2.1%'
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
                        <h2 class="text-lg font-bold">用户管理</h2>
                        <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">总数：{{ total }}</span>
                    </div>
                    <div class="flex items-center gap-3">
                        <div class="relative">
                            <input type="text" v-model="searchQuery" placeholder="搜索用户名、邮箱..."
                                class="pl-10 pr-4 py-2 w-64 rounded-lg border border-gray-200 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm"
                                @keyup.enter="handleSearch">
                            <i class="fas fa-search absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
                        </div>
                        <button @click="handleSearch" class="btn btn-primary">
                            <i class="fas fa-search mr-1"></i> 搜索
                        </button>
                        <button @click="showAddUserDialog" class="btn btn-success">
                            <i class="fas fa-plus mr-1"></i> 新增用户
                        </button>
                    </div>
                </div>
            </div>

            <!-- 表格部分 -->
            <div class="overflow-x-auto">
                <el-table :data="userList" style="width: 100%" v-loading="loading" class="custom-table">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column prop="avatar" label="头像" width="80">
                        <template #default="{ row }">
                            <img :src="row.avatar" class="w-10 h-10 rounded-full object-cover" :alt="row.username">
                        </template>
                    </el-table-column>
                    <el-table-column prop="username" label="用户名" />
                    <el-table-column prop="email" label="邮箱" />
                    <el-table-column prop="status" label="状态" width="120">
                        <template #default="{ row }">
                            <span class="px-2 py-1 text-xs rounded-full"
                                :class="row.status === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                                {{ row.status === 1 ? '正常' : '已封禁' }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="lastLogin" label="最后登录" width="180" />
                    <el-table-column label="操作" width="280" fixed="right">
                        <template #default="{ row }">
                            <div class="flex items-center gap-2">
                                <button @click="handleToggleBan(row)" class="btn-sm"
                                    :class="row.status === 1 ? 'btn-danger' : 'btn-success'">
                                    <i :class="row.status === 1 ? 'fas fa-ban mr-1' : 'fas fa-unlock mr-1'"></i>
                                    {{ row.status === 1 ? '封禁' : '解封' }}
                                </button>
                                <button @click="handleForceLogout(row)" class="btn-sm btn-warning"
                                    :disabled="row.status === 0">
                                    <i class="fas fa-sign-out-alt mr-1"></i> 强退
                                </button>
                                <button @click="handleViewDetails(row)" class="btn-sm btn-primary">
                                    <i class="fas fa-eye mr-1"></i> 详情
                                </button>
                                <button @click="handleReset(row)" class="btn-sm btn-info">
                                    <i class="fas fa-key mr-1"></i> 重置
                                </button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>

            <!-- 分页部分 -->
            <div class="p-4 border-t border-gray-200 flex justify-end">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" :total="total"
                    @size-change="handleSizeChange" @current-change="handleCurrentChange" background
                    class="custom-pagination" />
            </div>
        </div>

        <!-- 用户详情弹窗 -->
        <el-dialog v-model="detailsVisible" title="用户详情" width="700px" destroy-on-close>
            <div v-if="selectedUser" class="p-4">
                <div class="flex items-start gap-6 mb-6">
                    <img :src="selectedUser.avatar" class="w-20 h-20 rounded-lg object-cover"
                        :alt="selectedUser.username">
                    <div class="flex-1">
                        <h3 class="text-xl font-bold mb-2">{{ selectedUser.username }}</h3>
                        <div class="grid grid-cols-2 gap-4">
                            <div class="flex items-center gap-2 text-gray-600">
                                <i class="fas fa-envelope text-indigo-500"></i>
                                <span>{{ selectedUser.email }}</span>
                            </div>
                            <div class="flex items-center gap-2 text-gray-600">
                                <i class="fas fa-id-card text-indigo-500"></i>
                                <span>ID: {{ selectedUser.id }}</span>
                            </div>
                        </div>
                    </div>
                    <span class="px-3 py-1 text-sm rounded-full"
                        :class="selectedUser.status === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                        {{ selectedUser.status === 1 ? '正常' : '已封禁' }}
                    </span>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                    <div class="bg-gray-50 p-4 rounded-lg">
                        <h4 class="text-sm font-medium text-gray-500 mb-2">用户活动</h4>
                        <div class="text-lg font-semibold">最近登录时间</div>
                        <p class="text-gray-600">{{ selectedUser.lastLogin }}</p>
                    </div>
                    <div class="bg-gray-50 p-4 rounded-lg">
                        <h4 class="text-sm font-medium text-gray-500 mb-2">账户信息</h4>
                        <div class="text-lg font-semibold">注册时间</div>
                        <p class="text-gray-600">{{ selectedUser.registerTime }}</p>
                    </div>
                </div>

                <div class="bg-gray-50 p-4 rounded-lg mb-6">
                    <h4 class="text-sm font-medium text-gray-500 mb-2">安全信息</h4>
                    <div class="text-lg font-semibold">最后登录IP</div>
                    <p class="text-gray-600">{{ selectedUser.lastLoginIp }}</p>
                </div>

                <div class="flex justify-end gap-3">
                    <button @click="handleToggleBan(selectedUser); detailsVisible = false" class="btn"
                        :class="selectedUser.status === 1 ? 'btn-danger' : 'btn-success'">
                        <i :class="selectedUser.status === 1 ? 'fas fa-ban mr-2' : 'fas fa-unlock mr-2'"></i>
                        {{ selectedUser.status === 1 ? '封禁用户' : '解除封禁' }}
                    </button>
                    <button @click="detailsVisible = false" class="btn btn-gray">
                        <i class="fas fa-times mr-2"></i>关闭
                    </button>
                </div>
            </div>
        </el-dialog>

        <!-- 新增用户弹窗 -->
        <el-dialog v-model="addUserVisible" title="新增用户" width="500px" destroy-on-close>
            <div class="p-4">
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
                    <input type="text" v-model="newUser.username"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                    <input type="email" v-model="newUser.email"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">密码</label>
                    <input type="password" v-model="newUser.password"
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                </div>
                <div class="flex justify-end gap-3 mt-6">
                    <button @click="handleAddUser" class="btn btn-primary">
                        <i class="fas fa-save mr-2"></i>添加
                    </button>
                    <button @click="addUserVisible = false" class="btn btn-gray">
                        <i class="fas fa-times mr-2"></i>取消
                    </button>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
    name: 'UsersView',

    components: {
        Search
    },

    data() {
        return {
            searchQuery: '',
            loading: false,
            total: 100,
            currentPage: 1,
            pageSize: 10,
            detailsVisible: false,
            addUserVisible: false,
            selectedUser: null,
            newUser: {
                username: '',
                email: '',
                password: ''
            },
            userList: [
                {
                    id: 1,
                    username: 'john_doe',
                    email: 'john@example.com',
                    avatar: 'https://randomuser.me/api/portraits/men/32.jpg',
                    status: 1,
                    lastLogin: '2024-03-15 10:30:25',
                    registerTime: '2024-01-01 12:00:00',
                    lastLoginIp: '192.168.1.100'
                },
                {
                    id: 2,
                    username: 'jane_smith',
                    email: 'jane@example.com',
                    avatar: 'https://randomuser.me/api/portraits/women/44.jpg',
                    status: 1,
                    lastLogin: '2024-03-14 15:45:18',
                    registerTime: '2024-01-05 09:30:00',
                    lastLoginIp: '192.168.1.101'
                },
                {
                    id: 3,
                    username: 'mike_jackson',
                    email: 'mike@example.com',
                    avatar: 'https://randomuser.me/api/portraits/men/45.jpg',
                    status: 0,
                    lastLogin: '2024-03-10 08:20:13',
                    registerTime: '2024-01-10 14:20:00',
                    lastLoginIp: '192.168.1.102'
                },
                {
                    id: 4,
                    username: 'sarah_connor',
                    email: 'sarah@example.com',
                    avatar: 'https://randomuser.me/api/portraits/women/22.jpg',
                    status: 1,
                    lastLogin: '2024-03-15 12:10:05',
                    registerTime: '2024-01-15 11:45:00',
                    lastLoginIp: '192.168.1.103'
                },
                {
                    id: 5,
                    username: 'alex_williams',
                    email: 'alex@example.com',
                    avatar: 'https://randomuser.me/api/portraits/men/67.jpg',
                    status: 1,
                    lastLogin: '2024-03-14 23:05:42',
                    registerTime: '2024-01-20 16:30:00',
                    lastLoginIp: '192.168.1.104'
                }
            ]
        }
    },

    methods: {
        handleSearch() {
            this.loading = true
            // 模拟搜索请求
            setTimeout(() => {
                this.loading = false
                ElMessage({
                    message: '搜索完成',
                    type: 'success'
                })
            }, 500)
        },

        handleToggleBan(user) {
            ElMessageBox.confirm(
                `确认${user.status === 1 ? '封禁' : '解封'}用户 "${user.username}" ?`,
                '提示',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(() => {
                // 调用封禁/解封API
                ElMessage.success(`${user.status === 1 ? '封禁' : '解封'}成功`)
                user.status = user.status === 1 ? 0 : 1
            }).catch(() => { })
        },

        handleReset(user) {
            ElMessageBox.confirm(
                `确认为用户 "${user.username}" 重置密码?`,
                '提示',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(() => {
                ElMessage.success(`密码重置成功，新密码已发送至用户邮箱`)
            }).catch(() => { })
        },

        handleForceLogout(user) {
            ElMessageBox.confirm(
                `确认强制用户 "${user.username}" 退出登录?`,
                '提示',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(() => {
                // 调用强制退出API
                ElMessage.success('强制退出成功')
            }).catch(() => { })
        },

        handleViewDetails(user) {
            this.selectedUser = user
            this.detailsVisible = true
        },

        handleSizeChange(val) {
            this.pageSize = val
            this.fetchUsers()
        },

        handleCurrentChange(val) {
            this.currentPage = val
            this.fetchUsers()
        },

        fetchUsers() {
            this.loading = true
            // 模拟API请求
            setTimeout(() => {
                this.loading = false
            }, 500)
        },

        showAddUserDialog() {
            this.newUser = {
                username: '',
                email: '',
                password: ''
            }
            this.addUserVisible = true
        },

        handleAddUser() {
            if (!this.newUser.username || !this.newUser.email || !this.newUser.password) {
                ElMessage.error('请填写完整信息')
                return
            }

            // 模拟添加用户
            this.loading = true
            setTimeout(() => {
                // 添加新用户到列表顶部
                const newId = Math.max(...this.userList.map(u => u.id)) + 1
                this.userList.unshift({
                    id: newId,
                    username: this.newUser.username,
                    email: this.newUser.email,
                    avatar: `https://randomuser.me/api/portraits/men/${Math.floor(Math.random() * 100)}.jpg`,
                    status: 1,
                    lastLogin: '刚刚',
                    registerTime: new Date().toLocaleString(),
                    lastLoginIp: '127.0.0.1'
                })

                this.loading = false
                this.addUserVisible = false
                ElMessage.success('用户添加成功')
            }, 800)
        }
    },

    created() {
        this.fetchUsers()
    }
}
</script>

<style scoped>
/* 表格样式 */
.custom-table :deep(.el-table__header) {
    @apply bg-gray-50;
}

.custom-table :deep(.el-table__header th) {
    @apply bg-gray-50 text-gray-700 font-medium py-4;
}

.custom-table :deep(.el-table__row) {
    @apply hover:bg-gray-50 transition-colors;
}

.custom-table :deep(.el-table__row td) {
    @apply py-4;
}

/* 分页样式 */
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

.btn-success {
    @apply bg-green-600 hover:bg-green-700 text-white focus:ring-green-500;
}

.btn-danger {
    @apply bg-red-600 hover:bg-red-700 text-white focus:ring-red-500;
}

.btn-warning {
    @apply bg-amber-500 hover:bg-amber-600 text-white focus:ring-amber-500;
}

.btn-info {
    @apply bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500;
}

.btn-gray {
    @apply bg-gray-200 hover:bg-gray-300 text-gray-800 focus:ring-gray-500;
}

button:disabled {
    @apply opacity-50 cursor-not-allowed;
}
</style>
