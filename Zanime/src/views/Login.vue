<template>
    <div class="tt-smooth-scroll x-hidden">
        <!-- Login -->
        <div id="scroll-container">
            <HeaderMenu />
            <LoginForm />
            <Footer />
        </div>
        <MobileMenu />
    </div>
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import LoginForm from '@/components/LoginForm.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';

export default {
    components: {
        HeaderMenu,
        LoginForm,
        Footer,
        MobileMenu,
    },
    mounted() {
        app()
    }
};
import app from '@/static/js/app';
</script>

<style>
@import "@/static/css/app.css";
@import "@/static/css/swal.css";
</style>