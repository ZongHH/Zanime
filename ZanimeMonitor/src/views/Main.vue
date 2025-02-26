<template>
    <div class="main-dashboard">
        <!-- 统计卡片区域 -->
        <el-row :gutter="20" class="statistics">
            <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6" v-for="(item, index) in statisticsData" :key="index">
                <el-card shadow="hover" class="mb-20" @click="handleCardClick(item)">
                    <div class="statistics-item">
                        <div class="icon" :style="{ backgroundColor: item.color }">
                            <el-icon>
                                <component :is="item.icon" />
                            </el-icon>
                        </div>
                        <div class="content">
                            <div class="value">{{ item.value }}</div>
                            <div class="label">{{ item.label }}</div>
                        </div>
                    </div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 图表区域 -->
        <el-row :gutter="20" class="charts">
            <el-col :xs="24" :sm="24" :md="14" :lg="14" :xl="14" class="mb-20">
                <el-card shadow="hover">
                    <template #header>访问趋势</template>
                    <div ref="lineChartRef" style="height: 300px"></div>
                </el-card>
            </el-col>
            <el-col :xs="24" :sm="24" :md="10" :lg="10" :xl="10">
                <el-card shadow="hover" class="log-card" v-loading="loading">
                    <template #header>
                        <div class="log-header">
                            <span class="log-title">系统日志</span>
                            <div class="log-actions">
                                <div class="log-filters">
                                    <el-select v-model="selectedService" placeholder="选择服务" clearable
                                        class="service-select">
                                        <el-option v-for="service in services" :key="service.value"
                                            :label="service.label" :value="service.value" />
                                    </el-select>
                                    <el-select v-model="selectedLogLevel" placeholder="日志级别" clearable
                                        class="level-select">
                                        <el-option v-for="item in logLevels" :key="item.value" :label="item.label"
                                            :value="item.value" />
                                    </el-select>
                                </div>
                                <el-button type="primary" :icon="Expand" circle class="expand-button"
                                    @click="handleExpandLogs" />
                            </div>
                        </div>
                    </template>
                    <div class="log-content">
                        <el-timeline>
                            <el-timeline-item v-for="log in filteredLogs" :key="log.id" :type="log.levelType"
                                :timestamp="log.timestamp" :hollow="true">
                                <div class="log-item">
                                    <el-tag :type="log.levelType" size="small" effect="dark">
                                        {{ log.level }}
                                    </el-tag>
                                    <p class="log-message">{{ log.message }}</p>
                                    <p class="log-detail">{{ log.detail }}</p>
                                </div>
                            </el-timeline-item>
                        </el-timeline>
                    </div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 最近活动列表 -->
        <el-card shadow="hover" class="activity-list">
            <template #header>
                <div class="header">
                    <span>最近活动</span>
                    <el-button type="primary" link>查看全部</el-button>
                </div>
            </template>
            <el-timeline>
                <el-timeline-item v-for="(activity, index) in activities" :key="index" :timestamp="activity.time"
                    :type="activity.type">
                    {{ activity.content }}
                </el-timeline-item>
            </el-timeline>
        </el-card>
    </div>
</template>

<script>
import * as echarts from 'echarts'
import { Expand } from '@element-plus/icons-vue'

export default {
    name: 'MainDashboard',

    data() {
        return {
            statisticsData: [
                { label: '总用户数', value: '1,234', icon: 'User', color: '#409EFF' },
                { label: '视频数量', value: '856', icon: 'VideoPlay', color: '#67C23A' },
                { label: '文章数量', value: '432', icon: 'Document', color: '#E6A23C' },
                { label: '收藏数', value: '2,345', icon: 'Star', color: '#F56C6C' }
            ],
            activities: [
                { content: '用户张三注册了账号', time: '2023-05-20 12:00', type: 'primary' },
                { content: '新增视频《动画片名》', time: '2023-05-19 15:30', type: 'success' },
                { content: '用户李四发表了评论', time: '2023-05-19 10:00', type: 'info' }
            ],
            lineChartRef: null,
            logLevels: [
                { value: 'fatal', label: 'Fatal' },
                { value: 'error', label: 'Error' },
                { value: 'warning', label: 'Warning' },
                { value: 'info', label: 'Info' }
            ],
            selectedLogLevel: '',
            services: [
                { value: 'user', label: '用户服务' },
                { value: 'video', label: '视频服务' },
                { value: 'auth', label: '认证服务' },
                { value: 'upload', label: '上传服务' },
                { value: 'system', label: '系统服务' }
            ],
            selectedService: '',
            logs: [], // 移除硬编码的日志数据
            Expand, // 添加图标
            loading: false
        }
    },

    computed: {
        filteredLogs() {
            // 确保计算属性能够正确监听状态变化
            const logs = this.$store.getters.getSystemLogs
            let filtered = [...logs]

            if (this.selectedService) {
                filtered = filtered.filter(log =>
                    log.service === this.selectedService
                )
            }

            if (this.selectedLogLevel) {
                filtered = filtered.filter(log =>
                    log.level.toLowerCase() === this.selectedLogLevel.toLowerCase()
                )
            }

            return filtered
        }
    },

    methods: {
        initLineChart() {
            const lineChart = echarts.init(this.$refs.lineChartRef)
            this.lineChart = lineChart

            lineChart.setOption({
                tooltip: {
                    trigger: 'axis'
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
                },
                yAxis: {
                    type: 'value'
                },
                series: [{
                    data: [150, 230, 224, 218, 135, 147, 260],
                    type: 'line',
                    areaStyle: {},
                    smooth: true
                }]
            })
        },

        handleResize() {
            if (this.lineChart) {
                this.lineChart.resize()
            }
        },

        handleCardClick(item) {
            if (item.label === '总用户数') {
                this.$router.push('/users')
            }
        },

        handleExpandLogs() {
            this.$router.push('/logs')
        },

        async initData() {
            const currentLogs = this.$store.getters.getSystemLogs
            if (currentLogs.length === 0) {
                try {
                    this.loading = true
                    await this.$store.dispatch('fetchHistoryLogs')
                } catch (error) {
                    this.$message.error('加载历史日志失败：' + error.message)
                } finally {
                    this.loading = false
                }
            }
        }
    },

    async created() {
        await this.initData()
    },

    mounted() {
        this.$nextTick(() => {
            this.initLineChart()
            window.addEventListener('resize', this.handleResize)
        })
    },

    beforeDestroy() {
        window.removeEventListener('resize', this.handleResize)
        if (this.lineChart) {
            this.lineChart.dispose()
        }
    }
}
</script>

<style scoped>
.main-dashboard {
    padding: 10px;
}

@media screen and (min-width: 768px) {
    .main-dashboard {
        padding: 20px;
    }
}

.statistics {
    margin-bottom: 20px;
}

.mb-20 {
    margin-bottom: 20px;
}

.statistics-item {
    display: flex;
    align-items: center;
    cursor: pointer;
}

.statistics-item .icon {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 8px;
}

@media screen and (min-width: 768px) {
    .statistics-item .icon {
        width: 48px;
        height: 48px;
        margin-right: 12px;
    }
}

.statistics-item .icon :deep(svg) {
    font-size: 24px;
    color: #fff;
}

.statistics-item .content .value {
    font-size: 20px;
    font-weight: bold;
    margin-bottom: 4px;
}

.statistics-item .content .label {
    font-size: 14px;
    color: #909399;
}

.charts {
    margin-bottom: 20px;
}

.activity-list .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

/* 移动端适配 */
@media screen and (max-width: 767px) {
    .statistics-item .content .value {
        font-size: 16px;
    }

    .statistics-item .content .label {
        font-size: 12px;
    }

    .statistics-item .icon :deep(svg) {
        font-size: 20px;
    }
}

.log-card {
    height: 100%;
}

.log-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
}

.log-filters {
    display: flex;
    gap: 8px;
}

.log-title {
    flex-shrink: 0;
    font-size: 16px;
    font-weight: 500;
}

.service-select {
    width: 120px;
    flex-shrink: 0;
}

.level-select {
    width: 110px;
    flex-shrink: 0;
}

.log-actions {
    display: flex;
    align-items: center;
    gap: 12px;
}

.expand-button {
    flex-shrink: 0;
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
    .log-header {
        flex-direction: column;
        align-items: flex-start;
    }

    .log-filters {
        width: 100%;
        justify-content: space-between;
    }

    .service-select,
    .level-select {
        width: 48%;
    }

    .log-actions {
        width: 100%;
        justify-content: space-between;
    }
}

.log-content {
    height: 300px;
    overflow-y: auto;
}

.log-item {
    padding: 8px;
    border-radius: 4px;
}

.log-message {
    margin: 8px 0;
    font-weight: 500;
}

.log-detail {
    font-size: 12px;
    color: #909399;
    margin: 4px 0;
}

/* 自定义滚动条样式 */
.log-content::-webkit-scrollbar {
    width: 6px;
}

.log-content::-webkit-scrollbar-thumb {
    background-color: #dcdfe6;
    border-radius: 3px;
}

.log-content::-webkit-scrollbar-track {
    background-color: #f5f7fa;
}
</style>
