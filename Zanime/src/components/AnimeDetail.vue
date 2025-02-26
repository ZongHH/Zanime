<template>
    <div id="main-wrapper">
        <section class="anime-detail-section py-40">
            <div class="container">
                <!-- 添加无搜索结果提示 -->
                <div v-if="animeInfo.length === 0 && !isLoading" class="no-results">
                    <div class="no-results-content">
                        <i class="fas fa-search"></i>
                        <h3>未找到相关动漫</h3>
                        <p>试试其他关键词，或者返回首页查看更多动漫</p>
                        <button class="back-btn" @click="$router.push('/')">
                            返回首页
                        </button>
                    </div>
                </div>

                <div class="anime-detail-card mb-40" v-else v-for="info in animeInfo" :key="info.video_name">
                    <div class="info-section mb-32">
                        <div class="row">
                            <div class="col-md-2">
                                <div class="anime-cover">
                                    <img :src="info.cover_image_url" :alt="info.video_name">
                                </div>
                            </div>
                            <div class="col-md-10">
                                <div class="anime-info">
                                    <h2 class="white mb-16">{{ info.video_name }}</h2>
                                    <div class="info-grid mb-24">
                                        <div class="info-item">
                                            <span class="label">发布时间：</span>
                                            <span class="value">{{ info.release_date }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">地区：</span>
                                            <span class="value">{{ info.area }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">类型：</span>
                                            <span class="value">{{ info.genres }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">状态：</span>
                                            <span class="value"
                                                :class="info.status === '连载中' ? 'status-ongoing' : 'status-completed'">
                                                {{ info.status }}
                                            </span>
                                        </div>
                                    </div>
                                    <div class="description mb-24">
                                        <h5 class="white mb-12">简介</h5>
                                        <p>{{ info.description }}</p>
                                    </div>
                                    <div class="action-buttons">
                                        <button class="play-btn" @click="playLatestEpisode(info)">
                                            <i class="fas fa-play"></i>
                                            立即观看
                                        </button>
                                        <button class="collect-btn" :class="{ 'collected': info.is_collected }"
                                            @click="toggleCollect(info)">
                                            <i :class="info.is_collected ? 'fas fa-heart' : 'far fa-heart'"></i>
                                            {{ info.is_collected ? '已收藏' : '收藏' }}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="divider mb-32"></div>

                    <div class="episodes-section">
                        <div class="section-header mb-24">
                            <h4 class="white">剧集列表</h4>
                            <div class="episode-filters">
                                <button v-for="filter in episodeFilters" :key="filter.value"
                                    :class="['filter-btn', currentFilter === filter.value ? 'active' : '']"
                                    @click="currentFilter = filter.value">
                                    {{ filter.label }}
                                </button>
                            </div>
                        </div>
                        <div class="episodes-grid">
                            <button v-for="episode in paginatedEpisodes(info)" :key="info.video_name + episode" :class="[
                                'episode-btn',
                                // info.current_episode === episode ? 'active' : '',
                                // episode.watched ? 'watched' : '' //标记哪些是已观看状态
                            ]" @click="playEpisode(info.video_id, episode)">
                                {{ episode }}
                            </button>
                        </div>
                        <div class="pagination-controls mt-24">
                            <button class="page-btn" :disabled="info.currentPage === 1" @click="info.currentPage--">
                                上一页
                            </button>
                            <span class="page-info">{{ info.currentPage }}/{{ totalPages(info) }}</span>
                            <button class="page-btn" :disabled="info.currentPage === totalPages(info)"
                                @click="info.currentPage++">
                                下一页
                            </button>
                        </div>
                    </div>
                </div>

                <div ref="trigger" class="load-more-trigger"></div>

                <!-- <div class="recommendations-section">
                    <h4 class="white mb-24">相关推荐</h4>
                    <div class="recommendations-grid">
                        <div v-for="anime in recommendations" :key="anime.id" class="recommendation-card"
                            @click="goToAnime(anime.id)">
                            <div class="card-image">
                                <img :src="anime.cover_image_url" :alt="anime.video_name">
                            </div>
                            <div class="card-info">
                                <h6 class="white">{{ anime.video_name }}</h6>
                                <p>{{ anime.brief }}</p>
                            </div>
                        </div>
                    </div>
                </div> -->
            </div>
        </section>
    </div>
</template>

<script>
import axios from 'axios'
import { ElMessage } from 'element-plus'
import 'element-plus/es/components/message/style/css'  // 按需引入 Message 组件样式

export default {
    name: 'AnimeDetail',
    data() {
        return {
            animeInfo: [],
            Detailpage: 1, // 每次请求5条动漫数据
            pageSize: 30, // 动漫最多展示30集
            currentFilter: 'all',
            episodeFilters: [
                { label: '全部', value: 'all' },
                { label: '已观看', value: 'watched' },
                { label: '未观看', value: 'unwatched' }
            ],
            recommendations: [
                {
                    id: 1,
                    video_name: '进击的巨人 最终季',
                    cover_image_url: 'https://pic.imgdb.cn/item/65bde74f871b83018ac3c8a0.jpg',
                    brief: '人类与巨人的最终决战，艾伦发动地鸣，马莱与艾尔迪亚的命运走向终结。'
                },
            ]
        }
    },
    methods: {
        async fetchAnimeInfo() {
            if (this.isLoading) return;
            try {
                this.isLoading = true;
                const animeparams = this.$route.query.params
                const response = await axios.get(`/api/searchDetail?params=${animeparams}&page=${this.Detailpage}`)
                if (response.data.code === 200) {
                    const newdata = response.data.animes
                    newdata.forEach(anime => {
                        anime.currentPage = 1
                        anime.current_episode = anime.episodes[0]
                        anime.status = '连载中'
                    })
                    this.animeInfo = [...this.animeInfo, ...newdata]
                    this.Detailpage++
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                console.error('获取动漫信息失败:', error)
            } finally {
                this.isLoading = false; // 请求完成后重置标志位
            }
        },
        async fetchRecommendations() {
            try {
                const animeId = this.$route.query.id
                const response = await axios.get(`/api/anime/${animeId}/recommendations`)
                this.recommendations = response.data.data
            } catch (error) {
                console.error('获取推荐失败:', error)
            }
        },
        playEpisode(videoId, episode) {
            this.$router.push(`/moviesDetail?videoId=${videoId}&episode=${episode}`)
        },
        playLatestEpisode(info) {
            // 播放最新集或上次观看的集数
            // const lastWatched = this.episodes.findLast(ep => ep.watched)
            // const episodeToPlay = lastWatched ? lastWatched.episode_number + 1 : 1
            this.playEpisode(info.video_id, null)
        },
        async toggleCollect(info) {
            try {
                const response = await axios.post(`/api/user/update-collection`, {
                    video_id: info.video_id,
                    status: !info.is_collected
                })

                if (response.data.code === 200) {
                    // 更新当前动漫的收藏状态
                    this.animeInfo = this.animeInfo.map(anime => {
                        if (anime.video_id === info.video_id) {
                            return {
                                ...anime,
                                is_collected: !anime.is_collected
                            }
                        }
                        return anime
                    })

                    // 显示成功提示
                    ElMessage({
                        message: info.is_collected ? '已取消收藏' : '收藏成功',
                        type: 'success',
                        duration: 2000,
                        customClass: 'custom-message'
                    })
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                console.error('收藏操作失败:', error)
                // 显示错误提示
                ElMessage({
                    message: '操作失败：' + error.message,
                    type: 'error',
                    duration: 3000,
                    customClass: 'custom-message'
                })
            }
        },
        goToAnime(animeId) {
            this.$router.push(`/anime?id=${animeId}`)
        },
        totalPages(info) {
            return Math.ceil(info.episodes.length / this.pageSize)
        },
        filteredEpisodes(info) {
            if (this.currentFilter === 'all') return info.episodes
            return info.episodes.filter(episode =>
                this.currentFilter === 'watched' ? episode.watched : !episode.watched
            )
        },
        paginatedEpisodes(info) {
            const start = (info.currentPage - 1) * this.pageSize
            const end = start + this.pageSize
            return this.filteredEpisodes(info).slice(start, end)
        },
        createObserver() {
            const options = {
                root: null, // 整个视口
                threshold: 1.0, // 完全进入视口时触发
            };

            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        // 当第3个动漫进入视口，加载下一页数据
                        this.fetchAnimeInfo();
                    }
                });
            }, options);

            // 观察触发器元素
            this.$nextTick(() => {
                const triggerElement = this.$refs.trigger;
                if (triggerElement) observer.observe(triggerElement);
            });
        },
    },
    mounted() {
        this.createObserver()
        // this.fetchAnimeInfo()
        // this.fetchRecommendations()
    },
    watch: {
        '$route.query.params'() {
            this.animeInfo = [] // 清空现有数据
            this.Detailpage = 1 // 重置页码
            this.fetchAnimeInfo()
            // this.fetchRecommendations()
        }
    }
}
</script>

<style scoped>
@import "@/static/css/animeDetail.css";

.collect-btn {
    transition: all 0.3s ease;
}

.collect-btn.collected {
    background: rgba(204, 0, 0, 0.1);
    color: #cc0000;
}

.collect-btn.collected i {
    color: #cc0000;
}

.collect-btn:hover {
    transform: translateY(-1px);
}

.collect-btn.collected:hover {
    background: rgba(204, 0, 0, 0.15);
}

/* 自定义消息提示样式 */
.custom-message {
    min-width: 240px;
    padding: 12px 20px;
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 500;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 成功提示样式 */
.el-message.el-message--success {
    background: rgba(40, 167, 69, 0.95);
    border: none;
}

.el-message.el-message--success .el-message__content {
    color: white;
}

/* 错误提示样式 */
.el-message.el-message--error {
    background: rgba(220, 53, 69, 0.95);
    border: none;
}

.el-message.el-message--error .el-message__content {
    color: white;
}

/* 图标样式 */
.el-message .el-message__icon {
    margin-right: 10px;
    font-size: 1.1rem;
}

.el-message--success .el-message__icon {
    color: white;
}

.el-message--error .el-message__icon {
    color: white;
}

/* 无搜索结果样式 */
.no-results {
    min-height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 40px;
}

.no-results-content {
    max-width: 400px;
}

.no-results i {
    font-size: 48px;
    color: rgba(255, 255, 255, 0.2);
    margin-bottom: 20px;
}

.no-results h3 {
    color: #FFFFFF;
    font-size: 24px;
    margin-bottom: 12px;
}

.no-results p {
    color: rgba(255, 255, 255, 0.6);
    font-size: 14px;
    margin-bottom: 24px;
}

.back-btn {
    padding: 10px 24px;
    background: #AB0511;
    color: white;
    border: none;
    border-radius: 20px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.back-btn:hover {
    background: #8B0000;
    transform: translateY(-1px);
}

/* 移动端适配 */
@media (max-width: 768px) {
    .no-results {
        min-height: 300px;
        padding: 30px;
    }

    .no-results i {
        font-size: 40px;
        margin-bottom: 16px;
    }

    .no-results h3 {
        font-size: 20px;
    }

    .no-results p {
        font-size: 13px;
        margin-bottom: 20px;
    }

    .back-btn {
        padding: 8px 20px;
        font-size: 13px;
    }
}
</style>
