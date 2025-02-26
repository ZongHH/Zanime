<template>
    <div class="tt-smooth-scroll x-hidden">
        <!-- NotFound -->
        <div id="scroll-container">
            <HeaderMenu />
            <NotFound />
            <Footer />
        </div>
        <MobileMenu />
    </div>
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import NotFound from '@/components/NotFound.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';
import app from '@/static/js/app';

export default {
    components: {
        HeaderMenu,
        NotFound,
        Footer,
        MobileMenu,
    },
    mounted() {
        app()
    }
};
</script>

<style>
@import "@/static/css/app.css";
@import "@/static/css/swal.css";
</style>