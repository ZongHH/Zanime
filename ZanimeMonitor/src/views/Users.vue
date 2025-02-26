<template>
    <div class="users-container">
        <el-card>
            <template #header>
                <div class="header">
                    <div class="left">
                        <span class="title">用户管理</span>
                        <el-tag class="count" type="info">总数：{{ total }}</el-tag>
                    </div>
                    <div class="right">
                        <el-input v-model="searchQuery" placeholder="搜索用户..." class="search-input" clearable
                            @clear="handleSearch" @keyup.enter="handleSearch">
                            <template #prefix>
                                <el-icon>
                                    <Search />
                                </el-icon>
                            </template>
                        </el-input>
                        <el-button type="primary" @click="handleSearch">搜索</el-button>
                    </div>
                </div>
            </template>

            <el-table :data="userList" style="width: 100%" v-loading="loading">
                <el-table-column prop="id" label="ID" width="80" />
                <el-table-column prop="avatar" label="头像" width="80">
                    <template #default="{ row }">
                        <el-avatar :size="40" :src="row.avatar" />
                    </template>
                </el-table-column>
                <el-table-column prop="username" label="用户名" />
                <el-table-column prop="email" label="邮箱" />
                <el-table-column prop="status" label="状态" width="100">
                    <template #default="{ row }">
                        <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                            {{ row.status === 1 ? '正常' : '已封禁' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="lastLogin" label="最后登录" width="180" />
                <el-table-column label="操作" width="250" fixed="right">
                    <template #default="{ row }">
                        <el-button :type="row.status === 1 ? 'danger' : 'success'" size="small"
                            @click="handleToggleBan(row)">
                            {{ row.status === 1 ? '封禁' : '解封' }}
                        </el-button>
                        <el-button type="warning" size="small" @click="handleForceLogout(row)"
                            :disabled="row.status === 0">
                            强制退出
                        </el-button>
                        <el-button type="primary" size="small" @click="handleViewDetails(row)">
                            详情
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>

            <div class="pagination">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" :total="total"
                    @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>

        <!-- 用户详情弹窗 -->
        <el-dialog v-model="detailsVisible" title="用户详情" width="600px">
            <div v-if="selectedUser" class="user-details">
                <div class="user-header">
                    <el-avatar :size="80" :src="selectedUser.avatar" />
                    <div class="user-info">
                        <h3>{{ selectedUser.username }}</h3>
                        <p>{{ selectedUser.email }}</p>
                    </div>
                </div>
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="用户ID">{{ selectedUser.id }}</el-descriptions-item>
                    <el-descriptions-item label="注册时间">{{ selectedUser.registerTime }}</el-descriptions-item>
                    <el-descriptions-item label="最后登录">{{ selectedUser.lastLogin }}</el-descriptions-item>
                    <el-descriptions-item label="状态">
                        <el-tag :type="selectedUser.status === 1 ? 'success' : 'danger'">
                            {{ selectedUser.status === 1 ? '正常' : '已封禁' }}
                        </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="登录IP" :span="2">{{ selectedUser.lastLoginIp }}</el-descriptions-item>
                </el-descriptions>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { Search } from '@element-plus/icons-vue'  // 添加图标引入

export default {
    name: 'UsersView',

    components: {
        Search  // 注册Search图标组件
    },

    data() {
        return {
            searchQuery: '',
            loading: false,
            total: 100,
            currentPage: 1,
            pageSize: 10,
            detailsVisible: false,
            selectedUser: null,
            userList: [
                {
                    id: 1,
                    username: 'john_doe',
                    email: 'john@example.com',
                    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
                    status: 1,
                    lastLogin: '2024-03-15 10:30:25',
                    registerTime: '2024-01-01 12:00:00',
                    lastLoginIp: '192.168.1.100'
                },
                {
                    id: 2,
                    username: 'jim_doe',
                    email: 'john@example.com',
                    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
                    status: 1,
                    lastLogin: '2024-03-15 10:30:25',
                    registerTime: '2024-01-01 12:00:00',
                    lastLoginIp: '192.168.1.100'
                },
                {
                    id: 3,
                    username: 'jack_d',
                    email: 'john@example.com',
                    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
                    status: 1,
                    lastLogin: '2024-03-15 10:30:25',
                    registerTime: '2024-01-01 12:00:00',
                    lastLoginIp: '192.168.1.100'
                },
            ]
        }
    },

    methods: {
        handleSearch() {
            this.loading = true
            // 模拟搜索请求
            setTimeout(() => {
                this.loading = false
            }, 500)
        },

        handleToggleBan(user) {
            this.$confirm(
                `确认${user.status === 1 ? '封禁' : '解封'}用户 "${user.username}" ?`,
                '提示',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(() => {
                // 调用封禁/解封API
                this.$message.success(`${user.status === 1 ? '封禁' : '解封'}成功`)
                user.status = user.status === 1 ? 0 : 1
            }).catch(() => { })
        },

        handleForceLogout(user) {
            this.$confirm(
                `确认强制用户 "${user.username}" 退出登录?`,
                '提示',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(() => {
                // 调用强制退出API
                this.$message.success('强制退出成功')
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
        }
    },

    created() {
        this.fetchUsers()
    }
}
</script>

<style scoped>
.users-container {
    padding: 20px;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.left {
    display: flex;
    align-items: center;
    gap: 12px;
}

.title {
    font-size: 18px;
    font-weight: 600;
}

.right {
    display: flex;
    gap: 12px;
}

.search-input {
    width: 300px;
}

.pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}

.user-details {
    padding: 20px;
}

.user-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 20px;
}

.user-info h3 {
    margin: 0 0 8px 0;
}

.user-info p {
    margin: 0;
    color: #909399;
}

@media screen and (max-width: 768px) {
    .header {
        flex-direction: column;
        align-items: stretch;
        gap: 12px;
    }

    .search-input {
        width: 100%;
    }

    .right {
        flex-direction: column;
    }
}
</style>
