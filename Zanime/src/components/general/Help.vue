<template>
    <transition name="modal-fade">
        <div class="help-modal" v-if="isVisible" @click.self="closeModal">
            <div class="modal-content">
                <span class="close" @click="closeModal">&times;</span>
                <h2 class="section-title">常见问题</h2>
                <div class="faq-item" v-for="(faq, index) in faqs" :key="index">
                    <h3 class="faq-question">Q：{{ faq.question }}</h3>
                    <p class="faq-answer">A：{{ faq.answer }}</p>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
export default {
    props: {
        isVisible: {
            type: Boolean,
            required: true
        }
    },
    data() {
        return {
            faqs: [
                { question: '视频资源能出来吗？', answer: '可以的，请耐心等待。如果失败，请刷新页面重试。' }
            ]
        };
    },
    methods: {
        closeModal() {
            this.$emit('close');
        }
    }
};
</script>

<style scoped>
@import '@/static/css/personal.css';

.help-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(8px);
}

.modal-content {
    background: linear-gradient(145deg, #222, #333);
    padding: 30px;
    border-radius: 16px;
    width: 90%;
    max-width: 600px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.1);
    position: relative;
    transform: translateY(0);
}

.section-title {
    font-size: 2rem;
    margin-bottom: 25px;
    color: #fff;
    font-weight: 600;
    text-align: center;
    background: linear-gradient(90deg, #cfc4c4, #ff3333);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.faq-item {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    margin-bottom: 15px;
    padding: 20px;
    transition: all 0.3s;
    border-left: 4px solid #cc0000;
}

.faq-item:hover {
    background: rgba(255, 255, 255, 0.1);
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
}

.faq-question {
    color: #fff;
    font-size: 1.2rem;
    margin-bottom: 10px;
    font-weight: 500;
}

.faq-answer {
    color: #ddd;
    line-height: 1.6;
    font-size: 1.1rem;
}

.close {
    color: #999;
    position: absolute;
    top: 15px;
    right: 20px;
    font-size: 28px;
    font-weight: bold;
    transition: all 0.3s;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
}

.close:hover {
    color: #fff;
    background: rgba(255, 255, 255, 0.1);
    cursor: pointer;
    transform: rotate(90deg);
}

/* 弹窗过渡效果 */
.modal-fade-enter-active,
.modal-fade-leave-active {
    transition: opacity 0.3s, transform 0.3s;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
    opacity: 0;
    transform: scale(0.9);
}

.modal-fade-enter-to,
.modal-fade-leave-from {
    opacity: 1;
    transform: scale(1);
}
</style>
