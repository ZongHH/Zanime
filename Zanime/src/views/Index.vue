<template>
    <div class="tt-smooth-scroll x-hidden">
        <!-- PRELOADER START -->
        <div id="preloader">
            <div class="loader"></div>
        </div>
        <!-- PRELOADER END -->

        <!-- Main Wrapper Start -->
        <div id="scroll-container">
            <HeaderMenu />
            <IndexPage />
            <Footer />
        </div>
        <!-- Main Wrapper End -->
        <MobileMenu />
    </div>
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import IndexPage from '@/components/IndexPage.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';

export default {
    name: "Index",
    components: {
        HeaderMenu,
        IndexPage,
        Footer,
        MobileMenu,
    },
    mounted() {
        // app() 多个组件只要其中一个的mounted初始化app()，所有组件都能享受到
    }
};
</script>

<style>
@import "@/static/css/app.css";
@import "@/static/css/swal.css";
</style>