<!-- 代码已包含 CSS：使用 TailwindCSS , 安装 TailwindCSS 后方可看到布局样式效果 -->
<template>
    <div class="p-6">
        <div class="grid grid-cols-4 gap-6 mb-8">
            <div v-for="(stat, index) in statistics" :key="index" class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-gray-500 text-sm">{{ stat.title }}</h3>
                    <i :class="stat.icon" class="text-2xl" :style="{ color: stat.color }"></i>
                </div>
                <div class="flex items-end gap-2">
                    <span class="text-2xl font-bold">{{ stat.value }}</span>
                    <span class="text-sm text-green-500">{{ stat.growth }}</span>
                </div>
            </div>
        </div>
        <div class="bg-white rounded-lg p-6 mb-8 shadow-sm">
            <div class="flex items-center justify-between mb-6">
                <h2 class="text-lg font-bold">最新上线</h2>
                <button
                    class="text-sm text-indigo-600 hover:text-indigo-700 !rounded-button whitespace-nowrap">查看全部</button>
            </div>
            <div class="grid grid-cols-4 gap-6">
                <div v-for="(anime, index) in latestAnimes" :key="index" class="group relative">
                    <div class="aspect-[3/4] rounded-lg overflow-hidden">
                        <img :src="anime.image" :alt="anime.title"
                            class="w-full h-full object-cover transition-transform group-hover:scale-105">
                    </div>
                    <div class="mt-3">
                        <h3 class="font-medium text-sm mb-1">{{ anime.title }}</h3>
                        <p class="text-gray-500 text-xs">{{ anime.updateTime }} 更新</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="grid grid-cols-2 gap-6">
            <div class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-lg font-bold">热门推荐</h2>
                    <div class="flex gap-2">
                        <button v-for="(tab, index) in rankingTabs" :key="index"
                            class="px-3 py-1 text-sm !rounded-button whitespace-nowrap"
                            :class="currentTab === index ? 'bg-indigo-50 text-indigo-600' : 'text-gray-600 hover:bg-gray-50'"
                            @click="currentTab = index">
                            {{ tab }}
                        </button>
                    </div>
                </div>
                <div class="space-y-4">
                    <div v-for="(rank, index) in rankings" :key="index"
                        class="flex items-center gap-4 p-4 rounded-lg hover:bg-gray-50">
                        <span class="w-6 text-center font-bold"
                            :class="index < 3 ? 'text-indigo-600' : 'text-gray-400'">{{ index + 1 }}</span>
                        <img :src="rank.image" :alt="rank.title" class="w-16 h-16 rounded object-cover">
                        <div class="flex-1">
                            <h3 class="font-medium mb-1">{{ rank.title }}</h3>
                            <p class="text-sm text-gray-500">{{ rank.category }}</p>
                        </div>
                        <div class="text-right">
                            <div class="text-lg font-bold mb-1">{{ rank.score }}</div>
                            <div class="text-sm text-green-500">
                                <i class="fas fa-arrow-up mr-1"></i>
                                {{ rank.increase }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-lg p-6 shadow-sm">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-lg font-bold">操作日志</h2>
                    <button @click="$router.push('/operation-logs')"
                        class="text-sm text-indigo-600 hover:text-indigo-700 !rounded-button whitespace-nowrap">查看全部</button>
                </div>
                <div class="space-y-4">
                    <div v-for="(log, index) in operationLogs" :key="index"
                        class="flex items-start gap-4 p-4 rounded-lg hover:bg-gray-50">
                        <div class="w-8 h-8 rounded-full flex items-center justify-center"
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
                            <div class="mt-2 text-xs text-gray-400">
                                <i class="fas fa-tag mr-1"></i>
                                {{ log.module }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    data() {
        return {
            currentTab: 0,
            statistics: [
                {
                    title: '总动漫数量',
                    icon: 'fas fa-film',
                    value: '2,461',
                    growth: '+12.5%',
                    color: '#4F46E5'
                },
                {
                    title: '活跃用户数',
                    icon: 'fas fa-users',
                    value: '35,687',
                    growth: '+8.2%',
                    color: '#10B981'
                },
                {
                    title: '今日播放量',
                    icon: 'fas fa-play',
                    value: '89,254',
                    growth: '+15.3%',
                    color: '#F59E0B'
                },
                {
                    title: '收藏总数',
                    icon: 'fas fa-heart',
                    value: '126,519',
                    growth: '+9.1%',
                    color: '#EF4444'
                }
            ],
            latestAnimes: [
                {
                    title: '进击的巨人 最终季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/da0433310dda9b88a6cd293dca7ba515.jpg',
                    updateTime: '2023-12-01'
                },
                {
                    title: '咒术回战 第二季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/8b5fe3117f6aea8af4518e7246b6bedc.jpg',
                    updateTime: '2023-12-02'
                },
                {
                    title: '间谍过家家 第二季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/27ce123b6864a9d5621b64a2443ce08a.jpg',
                    updateTime: '2023-12-03'
                },
                {
                    title: '海贼王 第1089话',
                    image: 'https://ai-public.mastergo.com/ai/img_res/c5d3173638b3a3387849ff05aadb1b00.jpg',
                    updateTime: '2023-12-04'
                }
            ],
            rankingTabs: ['日榜', '周榜', '月榜'],
            operationLogs: [
                {
                    user: '系统管理员',
                    type: 'admin',
                    action: '更新了《咒术回战 第二季》的播放源信息',
                    time: '10分钟前',
                    module: '内容管理'
                },
                {
                    user: '王梓晨',
                    type: 'user',
                    action: '发表了《进击的巨人》的评论',
                    time: '25分钟前',
                    module: '用户互动'
                },
                {
                    user: '系统管理员',
                    type: 'admin',
                    action: '删除了违规用户评论',
                    time: '42分钟前',
                    module: '内容审核'
                },
                {
                    user: '林雨晴',
                    type: 'user',
                    action: '收藏了《间谍过家家》',
                    time: '1小时前',
                    module: '用户行为'
                },
                {
                    user: '系统管理员',
                    type: 'admin',
                    action: '新增了《海贼王》最新话章节',
                    time: '2小时前',
                    module: '内容管理'
                }
            ],
            rankings: [
                {
                    title: '咒术回战 第二季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/ee633a086e14ce80489204f88b38ee92.jpg',
                    category: '奇幻 / 动作',
                    score: '9.8',
                    increase: '15%'
                },
                {
                    title: '进击的巨人 最终季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/55857f0472b426b671fe7e6e2ab4f03a.jpg',
                    category: '剧情 / 动作',
                    score: '9.7',
                    increase: '12%'
                },
                {
                    title: '间谍过家家 第二季',
                    image: 'https://ai-public.mastergo.com/ai/img_res/ea480a57fa75d0f63fa83eff4158935a.jpg',
                    category: '喜剧 / 动作',
                    score: '9.5',
                    increase: '10%'
                },
                {
                    title: '海贼王',
                    image: 'https://ai-public.mastergo.com/ai/img_res/b986c79860122b6d5961d745a8aab7da.jpg',
                    category: '冒险 / 动作',
                    score: '9.4',
                    increase: '8%'
                },
                {
                    title: '鬼灭之刃',
                    image: 'https://ai-public.mastergo.com/ai/img_res/98102a62be3d56331287bdeb76835beb.jpg',
                    category: '奇幻 / 动作',
                    score: '9.3',
                    increase: '6%'
                }
            ]
        };
    }
});
</script>
<style scoped>
.custom-input::-webkit-inner-spin-button,
.custom-input::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
</style>