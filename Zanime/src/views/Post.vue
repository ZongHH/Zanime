<template>
    <HeaderMenu />
    <PostDetail />
    <MobileMenu />
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import PostDetail from '@/components/discussion/PostDetail.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';
import app from '@/static/js/app';

export default {
    components: {
        HeaderMenu,
        PostDetail,
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
