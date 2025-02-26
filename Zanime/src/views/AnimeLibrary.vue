<template>
    <HeaderMenu />
    <AnimeLibrary />
    <Footer />
    <MobileMenu />
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import AnimeLibrary from '@/components/AnimeLibrary.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';
import app from '@/static/js/app';

export default {
    components: {
        HeaderMenu,
        AnimeLibrary,
        Footer,
        MobileMenu,
    },
    mounted() {
        app()
    }
}
</script>

<style>
@import "@/static/css/app.css";
@import "@/static/css/swal.css";
</style>