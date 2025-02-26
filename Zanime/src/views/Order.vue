<template>
    <div class="order-view">
        <div class="">
            <HeaderMenu />
            <OrderPage />
        </div>
        <Footer />
        <MobileMenu />
    </div>
</template>

<script>
import HeaderMenu from '@/components/general/HeaderMenu.vue';
import OrderPage from '@/components/OrderPage.vue';
import Footer from '@/components/general/Footer.vue';
import MobileMenu from '@/components/general/MobileMenu.vue';
import app from '@/static/js/app';

export default {
    name: "Order",
    components: {
        HeaderMenu,
        OrderPage,
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

.order-view {
    position: relative;
    min-height: 100vh;
    background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
}
</style>