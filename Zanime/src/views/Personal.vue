<template>
    <HeaderMenu />
    <PersonalPage />
    <Footer />
    <MobileMenu />
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import PersonalPage from '@/components/personal/PersonalPage.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';
import app from '@/static/js/app';

export default {
    components: {
        HeaderMenu,
        PersonalPage,
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
