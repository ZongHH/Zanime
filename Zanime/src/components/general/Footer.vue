<template>
    <footer class="py-64">
        <div class="container-fluid">
            <div class="row row-gap-48">
                <div class="col-lg-6">
                    <div class="text">
                        <h6 class="mb-16 color-primary">开启您的娱乐之旅</h6>
                        <h3 class="white mb-48">开启您的娱乐之旅</h3>
                        <div class="newsletter-field">
                            <form @submit.prevent="handleSubscribe">
                                <input type="email" v-model="email" name="email" class="form-input"
                                    placeholder="输入您的邮箱">
                                <button type="submit" class="submit-arrow">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewbox="0 0 24 24"
                                        fill="none">
                                        <g clip-path="url(#clip0_11735_1263)">
                                            <path
                                                d="M23.7939 0.205977C23.595 0.00703961 23.2966 -0.0541321 23.0356 0.0503052L0.441886 9.08769C0.183371 9.1911 0.0102616 9.43729 0.000324102 9.71558C-0.0095665 9.99383 0.145683 10.2517 0.396183 10.3732L9.3056 14.6942L13.6266 23.6036C13.7445 23.8466 13.9906 23.9999 14.2591 23.9999C14.2675 23.9999 14.2758 23.9998 14.2842 23.9994C14.5624 23.9896 14.8087 23.8164 14.9121 23.5579L23.9496 0.964319C24.054 0.703132 23.9928 0.404867 23.7939 0.205977ZM2.44176 9.80235L20.3773 2.62819L9.68848 13.3169L2.44176 9.80235ZM14.1975 21.5581L10.6828 14.3112L21.3717 3.6225L14.1975 21.5581Z"
                                                fill="#AB0511"></path>
                                        </g>
                                        <defs>
                                            <clippath id="clip0_11735_1263">
                                                <rect width="24" height="24" fill="white"></rect>
                                            </clippath>
                                        </defs>
                                    </svg>
                                </button>
                            </form>
                        </div>
                    </div>
                </div>
                <div class="col-lg-6">
                    <img src="@/static/picture/pngwing.com1.png" alt="pic">
                </div>
                <div class="col-lg-6">
                    <a href="#">
                        <img src="@/static/picture/logo.png" alt="logo" class="header-logo pt-12 mb-24">
                    </a>
                    <p>Copyright &copy; 2024 版权所有</p>
                </div>
                <div class="col-lg-6">
                    <div class="row row-gap-4">
                        <div class="col-md-3 col-5">
                            <ul class="list-unstyled">
                                <li class="mb-12">
                                    <a href="#">首页</a>
                                </li>
                                <li class="mb-12">
                                    <a href="#">动漫</a>
                                </li>
                                <li class="mb-12">
                                    <a href="#">电影</a>
                                </li>
                                <li>
                                    <a href="#">博客</a>
                                </li>
                            </ul>
                        </div>
                        <div class="col-md-4 col-7">
                            <ul class="list-unstyled">
                                <li class="mb-12">
                                    <a href="#">服务条款</a>
                                </li>
                                <li class="mb-12">
                                    <a href="#">隐私政策</a>
                                </li>
                                <li>
                                    <a href="#">Cookie政策</a>
                                </li>
                            </ul>
                        </div>
                        <div class="col-md-4 col-12">
                            <ul class="list-unstyled">
                                <li class="mb-12">
                                    <a href="#" class="white">在线客服</a>
                                </li>
                                <li class="mb-12">
                                    <a href="#">support@example.com</a>
                                </li>
                                <li class="white">
                                    <a href="#" class="white">关注我们</a>
                                </li>
                            </ul>
                            <ul class="list-unstyled gap-16 d-flex">
                                <li>
                                    <a href="#">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewbox="0 0 24 24" fill="none">
                                            <path
                                                d="M16.7273 4H14.3471C13.295 4 12.286 4.42143 11.5421 5.17157C10.7981 5.92172 10.3802 6.93913 10.3802 8V10.4H8V13.6H10.3802V20H13.5537V13.6H15.9339L16.7273 10.4H13.5537V8C13.5537 7.78783 13.6373 7.58434 13.7861 7.43431C13.9349 7.28429 14.1367 7.2 14.3471 7.2H16.7273V4Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                        </svg>
                                    </a>
                                </li>

                                <li>
                                    <a href="#">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewbox="0 0 24 24" fill="none">
                                            <path d="M19 5L10 14" stroke="#EAEAEC" stroke-linecap="round"
                                                stroke-linejoin="round"></path>
                                            <path d="M19 5L13.4 21L10.2 13.8L3 10.6L19 5Z" stroke="#EAEAEC"
                                                stroke-linecap="round" stroke-linejoin="round"></path>
                                        </svg>
                                    </a>
                                </li>
                                <li>
                                    <a href="#">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewbox="0 0 24 24" fill="none">
                                            <path fill-rule="evenodd" clip-rule="evenodd"
                                                d="M8.26748 5C6.4629 5 5 6.4629 5 8.26748V15.8024C5 17.607 6.4629 19.0699 8.26748 19.0699H15.8024C17.607 19.0699 19.0699 17.607 19.0699 15.8024V8.26748C19.0699 6.4629 17.607 5 15.8024 5H8.26748ZM4 8.26748C4 5.91062 5.91062 4 8.26748 4H15.8024C18.1593 4 20.0699 5.91062 20.0699 8.26748V15.8024C20.0699 18.1593 18.1593 20.0699 15.8024 20.0699H8.26748C5.91062 20.0699 4 18.1593 4 15.8024V8.26748ZM12.4364 9.51543C11.9134 9.43786 11.3792 9.52721 10.9098 9.77075C10.4404 10.0143 10.0598 10.3996 9.82209 10.872C9.58436 11.3443 9.50161 11.8795 9.58562 12.4016C9.66962 12.9237 9.91611 13.406 10.29 13.7799C10.6639 14.1538 11.1462 14.4003 11.6683 14.4843C12.1903 14.5683 12.7256 14.4855 13.1979 14.2478C13.6703 14.0101 14.0556 13.6294 14.2991 13.1601C14.5427 12.6907 14.632 12.1565 14.5545 11.6335C14.4753 11.0999 14.2267 10.606 13.8453 10.2246C13.4639 9.84316 12.97 9.59454 12.4364 9.51543ZM10.4492 8.88312C11.1053 8.54271 11.852 8.41783 12.5831 8.52624C13.3289 8.63683 14.0193 8.98435 14.5524 9.51746C15.0855 10.0506 15.4331 10.741 15.5436 11.4868C15.6521 12.2179 15.5272 12.9646 15.1868 13.6207C14.8463 14.2767 14.3077 14.8087 13.6475 15.141C12.9873 15.4733 12.2391 15.589 11.5094 15.4716C10.7797 15.3541 10.1055 15.0096 9.58291 14.487C9.06027 13.9643 8.71574 13.2902 8.59832 12.5605C8.48089 11.8308 8.59656 11.0826 8.92886 10.4224C9.26116 9.76216 9.79317 9.22354 10.4492 8.88312ZM16.1792 7.39074C15.903 7.39074 15.6792 7.6146 15.6792 7.89074C15.6792 8.16688 15.903 8.39074 16.1792 8.39074H16.1878C16.464 8.39074 16.6878 8.16688 16.6878 7.89074C16.6878 7.6146 16.464 7.39074 16.1878 7.39074H16.1792Z"
                                                fill="#EAEAEC"></path>
                                        </svg>
                                    </a>
                                </li>
                                <li>
                                    <a href="#">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewbox="0 0 24 24" fill="none">
                                            <path
                                                d="M7 6.5C7 5.83696 7.26339 5.20107 7.73223 4.73223C8.20107 4.26339 8.83696 4 9.5 4H12V9H9.5C8.83696 9 8.20107 8.73661 7.73223 8.26777C7.26339 7.79893 7 7.16304 7 6.5Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                            <path
                                                d="M12 4H14.5C14.8283 4 15.1534 4.06466 15.4567 4.1903C15.76 4.31594 16.0356 4.50009 16.2678 4.73223C16.4999 4.96438 16.6841 5.23998 16.8097 5.54329C16.9353 5.84661 17 6.1717 17 6.5C17 6.8283 16.9353 7.15339 16.8097 7.45671C16.6841 7.76002 16.4999 8.03562 16.2678 8.26777C16.0356 8.49991 15.76 8.68406 15.4567 8.8097C15.1534 8.93534 14.8283 9 14.5 9H12V4Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                            <path
                                                d="M12 11.5C12 11.1717 12.0647 10.8466 12.1903 10.5433C12.3159 10.24 12.5001 9.96438 12.7322 9.73223C12.9644 9.50009 13.24 9.31594 13.5433 9.1903C13.8466 9.06466 14.1717 9 14.5 9C14.8283 9 15.1534 9.06466 15.4567 9.1903C15.76 9.31594 16.0356 9.50009 16.2678 9.73223C16.4999 9.96438 16.6841 10.24 16.8097 10.5433C16.9353 10.8466 17 11.1717 17 11.5C17 11.8283 16.9353 12.1534 16.8097 12.4567C16.6841 12.76 16.4999 13.0356 16.2678 13.2678C16.0356 13.4999 15.76 13.6841 15.4567 13.8097C15.1534 13.9353 14.8283 14 14.5 14C14.1717 14 13.8466 13.9353 13.5433 13.8097C13.24 13.6841 12.9644 13.4999 12.7322 13.2678C12.5001 13.0356 12.3159 12.76 12.1903 12.4567C12.0647 12.1534 12 11.8283 12 11.5V11.5Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                            <path
                                                d="M7 16.5C7 15.837 7.26339 15.2011 7.73223 14.7322C8.20107 14.2634 8.83696 14 9.5 14H12V16.5C12 17.163 11.7366 17.7989 11.2678 18.2678C10.7989 18.7366 10.163 19 9.5 19C8.83696 19 8.20107 18.7366 7.73223 18.2678C7.26339 17.7989 7 17.163 7 16.5Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                            <path
                                                d="M7 11.5C7 10.837 7.26339 10.2011 7.73223 9.73223C8.20107 9.26339 8.83696 9 9.5 9H12V14H9.5C8.83696 14 8.20107 13.7366 7.73223 13.2678C7.26339 12.7989 7 12.163 7 11.5Z"
                                                stroke="#EAEAEC" stroke-linecap="round" stroke-linejoin="round"></path>
                                        </svg>
                                    </a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </footer>
</template>

<script>
export default {
    name: "FooterComponent",
    data() {
        return {
            email: '',
        }
    },
    methods: {
        handleSubscribe(e) {
            e.preventDefault();
            // 处理订阅逻辑
            console.log('订阅邮箱:', this.email);
            this.email = ''; // 清空输入
        },
    }
};
</script>

<style scoped></style>