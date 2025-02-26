<template>
    <div class="logs-page">
        <el-card shadow="never" class="logs-card">
            <template #header>
                <div class="logs-header">
                    <div class="title-section">
                        <el-button :icon="Back" circle @click="$router.back()" />
                        <h2>系统日志</h2>
                    </div>
                    <div class="filter-section">
                        <el-select v-model="selectedService" placeholder="选择服务" clearable class="filter-item">
                            <el-option v-for="service in services" :key="service.value" :label="service.label"
                                :value="service.value" />
                        </el-select>
                        <el-select v-model="selectedLogLevel" placeholder="日志级别" clearable class="filter-item">
                            <el-option v-for="item in logLevels" :key="item.value" :label="item.label"
                                :value="item.value" />
                        </el-select>
                        <el-date-picker v-model="dateRange" type="daterange" range-separator="至"
                            start-placeholder="开始日期" end-placeholder="结束日期" class="filter-item" />
                    </div>
                </div>
            </template>

            <div class="logs-content" ref="logsContent" @scroll="handleScroll">
                <el-timeline>
                    <el-timeline-item v-for="log in displayedLogs" :key="log.id" :type="log.levelType"
                        :timestamp="log.timestamp" :hollow="true">
                        <div class="log-item">
                            <div class="log-header">
                                <el-tag :type="log.levelType" size="small" effect="dark">{{ log.level }}</el-tag>
                                <span class="service-tag">{{ log.service }}</span>
                            </div>
                            <p class="log-message">{{ log.message }}</p>
                            <p class="log-detail">{{ log.detail }}</p>
                        </div>
                    </el-timeline-item>
                </el-timeline>

                <div v-if="loading" class="loading-more">
                    <el-icon class="loading-icon">
                        <Loading />
                    </el-icon>
                    加载更多...
                </div>
            </div>
        </el-card>
    </div>
</template>

<script>
import { Back, Loading } from '@element-plus/icons-vue'

export default {
    name: 'LogsPage',
    data() {
        return {
            Back,
            selectedService: '',
            selectedLogLevel: '',
            dateRange: null,
            services: [
                { value: 'GateService', label: 'GateService' },
                { value: 'PushService', label: 'PushService' },
                { value: 'Crawler', label: 'Crawler' },
                { value: 'Recommend', label: 'Recommend' },
            ],
            logLevels: [
                { value: 'fatal', label: 'Fatal' },
                { value: 'error', label: 'Error' },
                { value: 'warning', label: 'Warning' },
                { value: 'info', label: 'Info' }
            ],
            currentPage: 1,
            pageSize: 20,
            currentDisplayCount: 20,
            loading: false,
            Loading
        }
    },
    computed: {
        filteredLogs() {
            let logs = this.$store.getters.getSystemLogs

            if (this.selectedService) {
                logs = logs.filter(log => log.service === this.selectedService)
            }

            if (this.selectedLogLevel) {
                logs = logs.filter(log => log.level.toLowerCase() === this.selectedLogLevel.toLowerCase())
            }

            if (this.dateRange) {
                const [start, end] = this.dateRange
                logs = logs.filter(log => {
                    const logDate = new Date(log.timestamp)
                    return logDate >= start && logDate <= end
                })
            }

            return logs
        },
        displayedLogs() {
            return this.filteredLogs.slice(0, this.currentDisplayCount)
        },
        hasMore() {
            return this.currentDisplayCount < this.filteredLogs.length
        }
    },
    methods: {
        handleScroll(e) {
            const { scrollHeight, scrollTop, clientHeight } = e.target
            if (scrollHeight - scrollTop - clientHeight < 50 && !this.loading && this.hasMore) {
                this.loadMore()
            }
        },
        async loadMore() {
            if (this.loading || !this.hasMore) return

            this.loading = true
            await new Promise(resolve => setTimeout(resolve, 500))
            this.currentDisplayCount += this.pageSize
            this.loading = false
        },
        resetDisplay() {
            this.currentDisplayCount = this.pageSize
            this.loading = false
            if (this.$refs.logsContent) {
                this.$refs.logsContent.scrollTop = 0
            }
        }
    },
    watch: {
        selectedService() {
            this.resetDisplay()
        },
        selectedLogLevel() {
            this.resetDisplay()
        },
        dateRange() {
            this.resetDisplay()
        }
    }
}
</script>

<style scoped>
.logs-page {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
}

.logs-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 20px;
}

.title-section {
    display: flex;
    align-items: center;
    gap: 12px;
}

.title-section h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 500;
}

.filter-section {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}

.filter-item {
    min-width: 160px;
}

.logs-card {
    min-height: calc(100vh - 40px);
    display: flex;
    flex-direction: column;
}

.logs-content {
    flex: 1;
    padding: 20px;
    overflow-y: auto;
    height: calc(100vh - 180px);
}

.log-item {
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;
    margin-bottom: 8px;
}

.log-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
}

.service-tag {
    color: #606266;
    font-size: 14px;
    background: #ebeef5;
    padding: 2px 8px;
    border-radius: 4px;
}

.log-message {
    font-size: 14px;
    color: #303133;
    margin: 8px 0;
    line-height: 1.5;
}

.log-detail {
    font-size: 13px;
    color: #909399;
    margin: 4px 0;
    font-family: monospace;
    white-space: pre-wrap;
}

.loading-more {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    color: #909399;
    font-size: 14px;
}

.loading-icon {
    margin-right: 8px;
    animation: rotating 2s linear infinite;
}

@keyframes rotating {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

@media screen and (max-width: 768px) {
    .logs-page {
        padding: 10px;
    }

    .logs-header {
        flex-direction: column;
        gap: 12px;
    }

    .filter-section {
        width: 100%;
    }

    .filter-item {
        min-width: 0;
        flex: 1;
    }

    .log-item {
        padding: 12px;
    }
}

/* 自定义滚动条样式 */
.logs-content::-webkit-scrollbar {
    width: 6px;
}

.logs-content::-webkit-scrollbar-thumb {
    background-color: #dcdfe6;
    border-radius: 3px;
}

.logs-content::-webkit-scrollbar-track {
    background-color: #f5f7fa;
}
</style>
