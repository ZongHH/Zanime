<template>
    <!-- HERO BANNER START -->
    <div class="anime-detail">
        <div class="trailer">
            <div class="videoplayer no_variable_height">
                <div id="video">
                    <!-- 添加加载动画 -->
                    <div class="video-loading-wrapper" v-show="isLoading">
                        <div class="loading-spinner"></div>
                        <p class="loading-text">视频加载中...</p>
                    </div>
                </div>
                <div class="progress-notification" :class="{ show: showNotification }">
                    上次观看到: {{ formattedProgress }}
                </div>
            </div>
        </div>
    </div>
    <!-- blogs-start -->
    <section class="anime-streaming bg-dark-black py-40">
        <div class="container-fluid">
            <h3 class="white mb-24">{{ videoInfo.video_name }}</h3>

            <!-- 添加集数选择区域 -->
            <div class="episodes-section mb-24">
                <h5 class="white mb-16">选集</h5>
                <div class="episodes-grid">
                    <button v-for="episode in videoInfo.episodes" :key="episode" :class="[
                        'episode-btn',
                        episode === currentEpisode ? 'active' : ''
                    ]" @click="changeEpisode(episode)">
                        {{ episode }}
                    </button>
                </div>
            </div>

            <Comment :videoId="parseInt($route.query.videoId)" />

            <div class="heading">
                <h4 class="white mb-12">相关推荐</h4>
                <a href="#" class="p white">查看全部 <span><i class="fa-regular fa-chevron-right"></i></span></a>
            </div>
            <div class="slider-arrows mt-sm-0 mt-48 d-xxl-flex d-none">
                <a href="javascript:;" class="arrow-btn btn-prev" data-slide="trendingSlider">
                    <i class="fas fa-chevron-left"></i>
                </a>
                <a href="javascript:;" class="arrow-btn btn-next" data-slide="trendingSlider">
                    <i class="fas fa-chevron-right"></i>
                </a>
            </div>
            <div class="trendingSlider row">
                <div class="col-12" v-for="anime in recommendedAnimes" :key="anime.id">
                    <a :href="`/moviesDetail?videoId=${anime.id}`" class="anime-card gap-24">
                        <img :src="anime.coverUrl" :alt="anime.title" class="fixed-size-img">
                        <div class="overlay"></div>
                        <div class="text">
                            <h6 class="white">{{ anime.title }}</h6>
                            <div class="rating">
                                <i class="fas fa-star"></i>
                                <span>{{ anime.rating }}</span>
                            </div>
                        </div>
                    </a>
                </div>
            </div>
        </div>
    </section>
    <!--blogs-->
    <!-- HERO BANNER END -->

    <!--  -->
    <!-- modal-popup area start  -->
    <div class="modal fade" id="videoModal" role="dialog" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="top_bar">
                    <h4 class="modal-title">Demon Slayer Season 4</h4>
                    <button type="button" class="close" id="closeVideoModalButton" data-dismiss="modal"
                        aria-label="Close">
                        <span aria-hidden="true"><i class="fas fa-times"></i> <b>Close</b></span>
                    </button>
                </div>
                <div class="modal-body">
                    <video controls="" title="Video">
                        <a href="index.html" aria-label="logo image"><img src="@/static/picture/logo.png"
                                alt="logo"></a>
                    </video>
                </div>
            </div>
        </div>
    </div>
    <!-- modal-popup area end  -->

    <Siderbar />

</template>

<script>
import axios from 'axios';
import Siderbar from './Siderbar.vue';
import app from '@/static/js/app.js';
import aksVideoPlayer from '@/static/js/aksVideoPlayer.js';
import { VideoPlayer, clearVideoPlayer, clearVideoPlayerOLD } from '@/static/js/app.js';
import Comment from './Comment.vue'
import { ElMessage } from 'element-plus';

export default {
    name: "MoviesDetail",
    data() {
        return {
            videoInfo: {},
            currentEpisode: "",
            showNotification: false,
            formattedProgress: "",
            notificationTimeout: null,
            userProgress: {},
            isLoading: true, // 添加加载状态
            recommendedAnimes: []
        }
    },
    components: {
        Siderbar,
        Comment
    },
    watch: {
        '$route.query': {
            handler: async function (newQuery) {
                await clearVideoPlayer(); // 清理旧的视频播放器
                this.isLoading = true; // 显示加载动画
                await this.fetchUserProgress();
                this.loadSavedProgress();
                await VideoPlayer(); // 重新初始化视频，放在loadSavedProgress后面，因为loadSavedProgress是监听DOM变化
                // clearVideoPlayerOLD() //存在历史遗留问题 当快速点击多集时，会残留多个视频播放器
            },
            deep: true // 深度监听对象的变化
        }
    },
    methods: {
        async changeEpisode(episode) {
            if (episode === this.currentEpisode) return; // 如果是同一集则不处理
            this.$router.push(`/moviesDetail?videoId=${this.videoInfo.video_id}&episode=${episode}`);
        },

        async fetchVideoInfo() {
            try {
                const urlParams = new URLSearchParams(window.location.search);
                const videoId = urlParams.get('videoId');

                const response = await axios.get(`/api/video-info?videoId=${videoId}`);
                this.videoInfo = response.data.video_info;
                this.videoInfo.video_id = videoId;

            } catch (error) {
                console.error('获取视频信息失败:', error);
            }
        },

        async fetchUserProgress() {
            try {
                const urlParams = new URLSearchParams(window.location.search)
                const videoId = urlParams.get('videoId')
                const episode = urlParams.get('episode')

                const response = await axios.get(`/api/load-progress?videoId=${videoId}&episode=${episode}`)
                if (response.data.code == 200) {
                    this.userProgress = response.data.progress
                    this.currentEpisode = response.data.progress.episode
                } else {
                    console.error("获取用户进度失败: ", response.data.message)
                }

            } catch (error) {
                console.error('获取用户进度失败:', error)
            }
        },

        async saveProgress(currentTime) {
            try {
                this.userProgress.episode = this.currentEpisode
                this.userProgress.progress = currentTime

                const response = await axios.post(`/api/save-progress`, this.userProgress)
                if (response.data.code != 200) {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                console.error('更新进度失败', error)
            }
            // localStorage.setItem(`videoProgress_${this.videoInfo.video_id}`, currentTime);
        },

        async saveProgressOnExit() {
            // 保存最新的观看进度
            const video = document.getElementById('aks-video');
            if (video) {
                const currentTime = Math.floor(video.currentTime);
                this.saveProgress(currentTime);
            }
        },

        async loadSavedProgress() {
            // 只能监听到 DOM 还没加载元素的变化，不能监听到已经加载完成的元素
            // 监听 DOM 变化，确保 aks-video 加载完成
            const observer = new MutationObserver(() => {
                const videoElement = document.getElementById('aks-video')
                if (videoElement) {
                    this.isLoading = false;
                    observer.disconnect(); // 找到目标后停止监听
                    // const savedProgress = localStorage.getItem(`videoProgress_${this.videoInfo.video_id}`);
                    const savedProgress = this.userProgress.progress
                    if (savedProgress) {
                        videoElement.currentTime = parseFloat(savedProgress);
                        this.showProgressNotification(savedProgress);
                    }
                }
            });

            // 开始监听 DOM 根节点的变化
            observer.observe(document.body, { childList: true, subtree: true });
        },

        showProgressNotification(savedProgress) {
            this.formattedProgress = this.formatTime(parseFloat(savedProgress));
            this.showNotification = true;

            if (this.notificationTimeout) {
                clearTimeout(this.notificationTimeout);
            }

            this.notificationTimeout = setTimeout(() => {
                this.showNotification = false;
            }, 10000);
        },

        formatTime(seconds) {
            const minutes = Math.floor(seconds / 60);
            const remainingSeconds = Math.floor(seconds % 60);
            return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
        },

        async initializeVideo() {
            await this.fetchUserProgress();     //获取用户进度
            this.loadSavedProgress();           //加载用户进度
            this.fetchVideoInfo();              //获取选集信息
            aksVideoPlayer();                   //初始化aksVideoPlayer
            app();                              //初始化页面事件
        },

        async fetchRecommendedAnimes() {
            try {
                const response = await axios.get('/api/movie/recommend', {
                    params: {
                        video_id: this.$route.query.videoId
                    }
                })
                if (response.data.code == 200) {
                    this.recommendedAnimes = response.data.recommendations
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                ElMessage.error('获取推荐动漫失败:' + error.message)
            }
        }
    },

    async mounted() {
        await this.fetchRecommendedAnimes();
        this.initializeVideo();
        window.addEventListener('beforeunload', this.saveProgressOnExit);
    },

    async beforeUnmount() {
        this.saveProgressOnExit()
        window.removeEventListener('beforeunload', this.saveProgressOnExit);
    },
}
</script>

<style>
@import "@/static/css/app.css";
@import "@/static/css/aksVideoPlayer.css";
@import "@/static/css/moviesDetail.css";
@import "@/static/css/swal.css";

.progress-notification {
    position: absolute;
    bottom: 60px;
    left: 20px;
    background: rgba(0, 0, 0, 0.85);
    color: white;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 14px;
    opacity: 0;
    transition: opacity 0.3s ease-in-out;
    z-index: 1000;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    display: flex;
    align-items: center;
    gap: 8px;
}

.progress-notification::before {
    content: '继续观看：';
    color: rgba(255, 255, 255, 0.8);
}

.progress-notification.show {
    opacity: 1;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .progress-notification {
        bottom: 50px;
        /* 稍微上移一点，避免与移动端控制栏重叠 */
        left: 10px;
        font-size: 12px;
        padding: 8px 12px;
        max-width: 80%;
        /* 防止文字太长溢出屏幕 */
    }

    .progress-notification::before {
        content: '续看：';
        /* 缩短提示文字 */
    }
}

/* 超小屏幕适配 */
@media (max-width: 480px) {
    .progress-notification {
        bottom: 45px;
        padding: 6px 10px;
        font-size: 11px;
    }
}

.episodes-section {
    background: rgba(26, 26, 26, 0.5);
    border-radius: 8px;
    padding: 20px;
}

.episodes-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 15px;
    margin-top: 15px;
}

.episode-btn {
    background: #2A2A2A;
    border: 1px solid #333;
    color: #FAFAFA;
    padding: 12px 8px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 14px;
    font-weight: 500;
    display: flex;
    align-items: center;
    justify-content: center;
}

.episode-btn:hover {
    background: #AB0511;
    border-color: #AB0511;
}

.episode-btn.active {
    background: #AB0511;
    border-color: #AB0511;
}

.episode-btn.watched {
    position: relative;
}

.episode-btn.watched::after {
    content: '';
    position: absolute;
    bottom: 4px;
    left: 50%;
    transform: translateX(-50%);
    width: 6px;
    height: 6px;
    background: #AB0511;
    border-radius: 50%;
}

@media (max-width: 768px) {
    .episodes-grid {
        grid-template-columns: repeat(auto-fill, minmax(60px, 1fr));
        gap: 10px;
    }

    .episode-btn {
        padding: 10px 6px;
        font-size: 13px;
    }
}

/* 响应式加载动画样式 */
.video-loading-wrapper {
    width: 100%;
    height: 56.25vw;
    /* 16:9 宽高比 */
    max-height: 850px;
    min-height: 250px;
    background: #1a1a1a;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    position: relative;
}

.loading-spinner {
    width: min(8vw, 50px);
    height: min(8vw, 50px);
    border: min(0.6vw, 4px) solid #333;
    border-top: min(0.6vw, 4px) solid #AB0511;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: min(2.5vw, 15px);
}

.loading-text {
    color: #FAFAFA;
    font-size: clamp(14px, 2.5vw, 16px);
    text-align: center;
    padding: 0 20px;
}

/* 媒体查询优化 */
@media (max-width: 768px) {
    .video-loading-wrapper {
        height: 70vw;
        /* 移动端调整高度比例 */
    }

    .loading-spinner {
        width: 35px;
        height: 35px;
        border-width: 3px;
        margin-bottom: 10px;
    }

    .loading-text {
        font-size: 14px;
    }
}

@media (max-width: 480px) {
    .video-loading-wrapper {
        height: 80vw;
    }

    .loading-spinner {
        width: 30px;
        height: 30px;
        border-width: 2px;
        margin-bottom: 8px;
    }

    .loading-text {
        font-size: 12px;
    }
}

/* 横屏模式优化 */
@media (orientation: landscape) and (max-height: 480px) {
    .video-loading-wrapper {
        height: 45vw;
        min-height: 180px;
    }

    .loading-spinner {
        width: 25px;
        height: 25px;
        margin-bottom: 5px;
    }

    .loading-text {
        font-size: 12px;
    }
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}

.fixed-size-img {
    width: 400px;
    /* 设置你想要的宽度 */
    height: 250px;
    /* 设置你想要的高度 */
    object-fit: cover;
    /* 保持图片比例，裁剪多余部分 */
}
</style>