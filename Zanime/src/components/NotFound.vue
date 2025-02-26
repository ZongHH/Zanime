<template>
    <!-- HERO BANNER START -->
    <div class="page">
        <div class="container-fluid">
            <div class="row justify-content-center text-center">
                <div class="col-lg-6 col-md-8 col-10">
                    <img src="@/static/picture/Frame.png" alt="pic">
                    <!-- 添加操作按钮 -->
                    <div class="error-actions mt-4">
                        <button class="btn btn-primary me-3" @click="goHome">
                            <i class="fa fa-home me-2"></i>返回首页
                        </button>
                        <button class="btn btn-outline-primary" @click="contactSupport">
                            <i class="fa fa-envelope me-2"></i>联系客服
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "404",
    methods: {
        goHome() {
            this.$router.push('/');
        },
        contactSupport() {
            window.location.href = 'mailto:support@example.com';
        },
    },
}
</script>

<style scoped>
.error-actions {
    margin-top: 2rem;
}

.btn {
    padding: 10px 20px;
    border-radius: 5px;
    transition: all 0.3s ease;
}

.btn-primary {
    background-color: #AB0511;
    border-color: #AB0511;
    color: white;
}

.btn-outline-primary {
    border-color: #AB0511;
    color: #AB0511;
}

.btn-outline-primary:hover {
    background-color: #AB0511;
    color: white;
}

.me-2 {
    margin-right: 0.5rem;
}

.me-3 {
    margin-right: 1rem;
}

.mt-4 {
    margin-top: 1.5rem;
}
</style>