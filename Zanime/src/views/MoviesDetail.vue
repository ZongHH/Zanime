<template>
    <div class="x-hidden">
        <!-- Main Wrapper Start -->
        <div id="main-wrapper">
            <HeaderMenu />
            <MoviesPage />
            <Footer />
        </div>
        <MobileMenu />
    </div>
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import MoviesPage from '@/components/movie/MoviesPage.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';

export default {
    name: "Movies",
    components: {
        HeaderMenu,
        MoviesPage,
        Footer,
        MobileMenu,
    },
};
</script>

<style></style>