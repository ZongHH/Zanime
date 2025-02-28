<template>
    <transition name="fade">
        <div class="help-modal" v-if="isVisible">
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
}

.modal-content {
    background: #2A2A2A;
    padding: 30px;
    border-radius: 12px;
    width: 90%;
    max-width: 700px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    transition: transform 0.3s ease;
}

.section-title {
    font-size: 1.8rem;
    margin-bottom: 20px;
    color: #fff;
}

.faq-item {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 10px;
    margin-bottom: 15px;
    padding: 20px;
    transition: background 0.3s;
}

.faq-item:hover {
    background: rgba(255, 255, 255, 0.1);
}

.faq-question {
    color: #fff;
    font-size: 1.2rem;
    margin: 0;
}

.faq-answer {
    margin-top: 10px;
    color: #ddd;
    transition: opacity 0.3s ease;
}

.close {
    color: #ccc;
    position: absolute;
    top: 15px;
    right: 20px;
    font-size: 28px;
    font-weight: bold;
}

.close:hover,
.close:focus {
    color: #fff;
    text-decoration: none;
    cursor: pointer;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.5s;
}

.fade-enter,
.fade-leave-to

/* .fade-leave-active in <2.1.8 */
    {
    opacity: 0;
}
</style>
