<template>
    <div class="order-page">
        <!-- 页面标题 -->
        <div class="page-header">
            <h1>我的订单</h1>
            <div class="filter-bar">
                <button class="filter-btn active">全部</button>
                <button class="filter-btn">待支付</button>
                <button class="filter-btn">已支付</button>
                <button class="filter-btn">已完成</button>
            </div>
        </div>

        <!-- 订单列表 -->
        <div class="order-container">
            <!-- 加载状态 -->
            <div v-if="isLoading" class="loading-state">
                <div class="loader">
                    <div class="loader-bar"></div>
                    <div class="loader-bar"></div>
                    <div class="loader-bar"></div>
                    <div class="loader-bar"></div>
                </div>
                <p>加载订单中...</p>
            </div>

            <!-- 空状态 -->
            <div v-else-if="orders && !orders.length" class="empty-state">
                <div class="empty-icon">
                    <svg viewBox="0 0 24 24" width="64" height="64">
                        <path fill="currentColor"
                            d="M19,3H5C3.89,3 3,3.89 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5C21,3.89 20.1,3 19,3M19,5V19H5V5H19Z" />
                        <path fill="currentColor"
                            d="M12,15.58L14.88,13.96L15.41,15.12L12,17L8.59,15.12L9.12,13.96L12,15.58Z" />
                        <path fill="currentColor"
                            d="M12,10.41L8.12,12.17L7.59,11L12,8.83L16.41,11L15.88,12.17L12,10.41Z" />
                    </svg>
                </div>
                <p>您暂时没有订单</p>
                <!-- <button class="btn-go-shop">去商城逛逛</button> -->
            </div>

            <!-- 订单列表 -->
            <div v-else-if="orders && orders.length" class="order-list">
                <div v-for="order in orders" :key="order.order_id" class="order-card">
                    <div class="order-header">
                        <div class="order-info">
                            <span class="order-id">订单号：{{ order.order_id }}</span>
                            <span class="order-date">{{ formatDate(order.create_time) }}</span>
                        </div>
                        <div class="order-status" :class="getStatusClass(order.status)">
                            <span class="status-dot"></span>
                            {{ order.status }}
                        </div>
                    </div>

                    <div class="order-content">
                        <div class="product-info">
                            <div class="product-image-container">
                                <img :src="order.image" :alt="order.product_name" class="product-image">
                            </div>
                            <div class="product-details">
                                <h5 class="product-name">{{ order.product_name }}</h5>
                                <div class="product-specs">
                                    <span class="specs-item">{{ order.selected_size }}</span>
                                    <span class="specs-divider"></span>
                                    <span class="specs-item">{{ order.selected_color }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="order-price">
                            <span class="price-currency">¥</span>
                            <span class="price-amount">{{ order.price }}</span>
                        </div>
                    </div>

                    <div class="order-footer">
                        <div class="shipping-info">
                            <div class="shipping-item">
                                <span class="shipping-label">收货人</span>
                                <span class="shipping-value">{{ order.user_name }}</span>
                            </div>
                            <div class="shipping-item">
                                <span class="shipping-label">电话</span>
                                <span class="shipping-value">{{ order.phone }}</span>
                            </div>
                            <div class="shipping-item">
                                <span class="shipping-label">地址</span>
                                <span class="shipping-value address">{{ order.address }}</span>
                            </div>
                        </div>
                        <div class="order-actions">
                            <button v-if="order.status === '待支付'" class="btn-primary action-btn"
                                @click="handlePayment(order)">
                                <span class="btn-icon">
                                    <svg viewBox="0 0 24 24" width="16" height="16">
                                        <path fill="currentColor"
                                            d="M20,8H4V6H20M20,18H4V12H20M20,4H4C2.89,4 2,4.89 2,6V18A2,2 0 0,0 4,20H20A2,2 0 0,0 22,18V6C22,4.89 21.1,4 20,4Z" />
                                    </svg>
                                </span>
                                <span>继续支付</span>
                            </button>
                            <button v-else-if="order.status === '已完成'" class="btn-secondary action-btn"
                                @click="handleReview(order)">
                                <span class="btn-icon">
                                    <svg viewBox="0 0 24 24" width="16" height="16">
                                        <path fill="currentColor"
                                            d="M12,4A4,4 0 0,1 16,8A4,4 0 0,1 12,12A4,4 0 0,1 8,8A4,4 0 0,1 12,4M12,14C16.42,14 20,15.79 20,18V20H4V18C4,15.79 7.58,14 12,14Z" />
                                    </svg>
                                </span>
                                <span>评价</span>
                            </button>
                        </div>
                    </div>

                    <div class="order-time-progress" v-if="order.status === '待支付'">
                        <div class="time-bar">
                            <div class="time-progress"></div>
                        </div>
                        <div class="time-text">支付剩余时间: 30分钟</div>
                    </div>
                </div>

                <!-- 加载更多 -->
                <div v-if="hasMore" class="load-more">
                    <button class="btn-load-more" @click="loadMore" :disabled="isLoadingMore">
                        <span v-if="isLoadingMore">
                            <div class="btn-loader"></div>
                            加载中...
                        </span>
                        <span v-else>查看更多订单</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default {
    name: 'OrderPage',
    data() {
        return {
            orders: [],
            isLoading: true,
            isLoadingMore: false,
            page: 1,
            pageSize: 10,
            hasMore: true
        }
    },
    methods: {
        async fetchOrders(loadMore = false) {
            if (!loadMore) {
                this.isLoading = true;
            } else {
                this.isLoadingMore = true;
            }

            try {
                // 发送json格式数据
                const response = await axios.get('/api/get-orders', {
                    params: {
                        page: this.page,
                        pageSize: this.pageSize
                    }
                });

                if (response.data.code == 200) {
                    if (response.data.orders == null) {
                        this.hasMore = false;
                        return;
                    }

                    const newOrders = response.data.orders;

                    if (loadMore) {
                        this.orders = [...this.orders, ...newOrders];
                    } else {
                        this.orders = newOrders;
                    }

                    this.hasMore = newOrders.length === this.pageSize;
                } else {
                    throw new Error(response.data.message);
                }
            } catch (error) {
                ElMessage.warning('获取订单失败: ' + error.message);
            } finally {
                if (loadMore) {
                    this.isLoadingMore = false;
                } else {
                    this.isLoading = false;
                }
            }
        },

        formatDate(dateString) {
            const date = new Date(dateString);
            return date.toLocaleDateString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        },

        getStatusClass(status) {
            const classMap = {
                '待支付': 'status-pending',
                '已支付': 'status-paid',
                '已完成': 'status-completed',
                '已失效': 'status-expired'
            };
            return classMap[status] || '';
        },

        async loadMore() {
            this.page += 1;
            await this.fetchOrders(true);
        },

        handlePayment(order) {
            // 处理继续支付逻辑
        },

        handleReview(order) {
            // 处理评价逻辑
        }
    },
    mounted() {
        this.fetchOrders();
    },
}
</script>

<style scoped>
.order-page {
    padding: 2rem;
    min-height: 100vh;
    background-color: #121212;
    color: #fff;
    font-family: 'Poppins', -apple-system, BlinkMacSystemFont, sans-serif;
    display: flex;
    flex-direction: column;
    align-items: center;
}

/* 页面标题样式 */
.page-header {
    margin-bottom: 3rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 1.5rem;
    width: 60%;
}

.page-header h1 {
    font-size: 2.2rem;
    font-weight: 600;
    color: #fff;
    margin: 0;
    letter-spacing: -0.5px;
}

.filter-bar {
    display: flex;
    gap: 0.8rem;
    flex-wrap: wrap;
}

.filter-btn {
    background: #000000;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    padding: 0.5rem 1.2rem;
    color: #ffffff;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.25s ease;
    font-weight: 500;
}

.filter-btn:hover {
    border-color: rgba(255, 255, 255, 0.3);
}

.filter-btn.active {
    background: #ffffff;
    color: #000000;
    border-color: #ffffff;
    box-shadow: 0 2px 8px rgba(255, 255, 255, 0.15);
}

/* 订单容器样式 */
.order-container {
    width: 60%;
    margin: 0 auto;
}

/* 订单卡片样式 */
.order-card {
    background: #1e1e1e;
    border-radius: 1rem;
    padding: 1.8rem;
    margin-bottom: 1.8rem;
    border: 1px solid rgba(255, 255, 255, 0.03);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    position: relative;
    overflow: hidden;
}

.order-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.25);
    border-color: rgba(255, 255, 255, 0.06);
}

.order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.8rem;
    padding-bottom: 1.2rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.order-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.order-id {
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.85);
    font-weight: 500;
}

.order-date {
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.5);
}

.order-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    padding: 0.3rem 1rem;
    border-radius: 20px;
    background-color: rgba(255, 255, 255, 0.03);
    color: #e0e0e0;
}

.status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    display: inline-block;
}

.status-pending .status-dot {
    background-color: #ff9800;
    box-shadow: 0 0 4px rgba(255, 152, 0, 0.4);
}

.status-paid .status-dot {
    background-color: #4caf50;
    box-shadow: 0 0 4px rgba(76, 175, 80, 0.4);
}

.status-completed .status-dot {
    background-color: #2196f3;
    box-shadow: 0 0 4px rgba(33, 150, 243, 0.4);
}

.status-expired .status-dot {
    background-color: #f44336;
    box-shadow: 0 0 4px rgba(244, 67, 54, 0.4);
}

/* 订单内容样式 */
.order-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.8rem;
}

.product-info {
    display: flex;
    gap: 1.5rem;
    align-items: center;
    flex: 1;
}

.product-image-container {
    width: 80px;
    height: 80px;
    border-radius: 12px;
    overflow: hidden;
    background: rgba(0, 0, 0, 0.3);
    position: relative;
}

.product-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
}

.product-image-container:hover .product-image {
    transform: scale(1.05);
}

.product-details {
    flex: 1;
}

.product-name {
    margin-bottom: 0.7rem;
    font-size: 1.1rem;
    font-weight: 500;
    color: #ffffff;
}

.product-specs {
    display: flex;
    align-items: center;
    gap: 0.8rem;
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
}

.specs-item {
    background: rgba(255, 255, 255, 0.04);
    padding: 0.2rem 0.8rem;
    border-radius: 12px;
}

.specs-divider {
    width: 3px;
    height: 3px;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.2);
}

.order-price {
    font-weight: 600;
    display: flex;
    align-items: baseline;
}

.price-currency {
    font-size: 1rem;
    color: rgba(255, 255, 255, 0.6);
    margin-right: 0.2rem;
}

.price-amount {
    font-size: 1.4rem;
    color: #ffffff;
}

/* 订单底部样式 */
.order-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 1.2rem;
    padding-top: 1.2rem;
    border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.shipping-info {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.shipping-item {
    display: flex;
    align-items: center;
    gap: 0.8rem;
}

.shipping-label {
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.85rem;
    min-width: 45px;
}

.shipping-value {
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.9rem;
}

.shipping-value.address {
    max-width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.order-actions {
    display: flex;
    gap: 1rem;
}

.action-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1.2rem;
    border-radius: 20px;
    border: none;
    font-weight: 500;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    outline: none;
}

.btn-icon {
    display: flex;
    align-items: center;
    justify-content: center;
}

.btn-primary {
    background: #3a3a3a;
    color: #fff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn-primary:hover {
    background: #454545;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.04);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.08);
}

.btn-secondary:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.12);
}

/* 加载中样式 */
.loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 300px;
    color: rgba(255, 255, 255, 0.6);
}

.loader {
    display: flex;
    align-items: flex-end;
    margin-bottom: 1.5rem;
    height: 40px;
}

.loader-bar {
    width: 4px;
    margin: 0 3px;
    background: #555;
    border-radius: 2px;
    animation: loader 1s ease-in-out infinite;
}

.loader-bar:nth-child(1) {
    animation-delay: 0s;
    height: 20px;
}

.loader-bar:nth-child(2) {
    animation-delay: 0.2s;
    height: 30px;
}

.loader-bar:nth-child(3) {
    animation-delay: 0.4s;
    height: 10px;
}

.loader-bar:nth-child(4) {
    animation-delay: 0.6s;
    height: 25px;
}

@keyframes loader {
    0% {
        height: 10px;
    }

    50% {
        height: 30px;
    }

    100% {
        height: 10px;
    }
}

/* 空状态样式 */
.empty-state {
    text-align: center;
    padding: 5rem 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.empty-icon {
    margin-bottom: 1.8rem;
    color: rgba(255, 255, 255, 0.15);
    animation: float 3s ease-in-out infinite;
}

@keyframes float {
    0% {
        transform: translateY(0px);
    }

    50% {
        transform: translateY(-8px);
    }

    100% {
        transform: translateY(0px);
    }
}

.empty-state p {
    color: rgba(255, 255, 255, 0.5);
    font-size: 1.1rem;
    margin-bottom: 1.8rem;
}

.btn-go-shop {
    background: #3a3a3a;
    color: #fff;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 25px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn-go-shop:hover {
    transform: translateY(-2px);
    background: #454545;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

/* 加载更多 */
.load-more {
    text-align: center;
    padding: 2rem 0 3rem;
}

.btn-load-more {
    background: rgb(0 0 0);
    color: #ffffff;
    border: 1px solid rgba(255, 255, 255, 0.08);
    padding: 0.8rem 2.5rem;
    border-radius: 25px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
    gap: 0.7rem;
}

.btn-load-more:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.12);
}

.btn-load-more:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.btn-loader {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    border-top-color: #ffffff;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

/* 订单倒计时 */
.order-time-progress {
    margin-top: 1rem;
    padding-top: 0.5rem;
}

.time-bar {
    height: 2px;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 1px;
    overflow: hidden;
    margin-bottom: 0.5rem;
}

.time-progress {
    height: 100%;
    width: 40%;
    background: #aaaaaa;
    border-radius: 1px;
    animation: timeProgress 30s linear;
}

@keyframes timeProgress {
    0% {
        width: 100%;
    }

    100% {
        width: 0%;
    }
}

.time-text {
    font-size: 0.8rem;
    color: #aaaaaa;
    text-align: right;
}

/* 响应式样式 */
@media (max-width: 1200px) {

    .page-header,
    .order-container {
        width: 80%;
    }
}

@media (max-width: 992px) {

    .page-header,
    .order-container {
        width: 90%;
    }
}

@media (max-width: 768px) {
    .order-page {
        padding: 1.5rem 1rem;
    }

    .page-header,
    .order-container {
        width: 100%;
    }

    .page-header {
        flex-direction: column;
        align-items: flex-start;
    }

    .order-content {
        flex-direction: column;
        align-items: flex-start;
    }

    .product-info {
        width: 100%;
        margin-bottom: 1rem;
    }

    .order-price {
        align-self: flex-end;
    }

    .order-footer {
        flex-direction: column;
        align-items: flex-start;
        gap: 1.5rem;
    }

    .shipping-info {
        width: 100%;
    }

    .order-actions {
        align-self: flex-end;
    }
}

/* 移动端优化 */
@media (max-width: 576px) {
    .order-page {
        padding: 1rem 0.8rem;
    }

    .page-header h1 {
        font-size: 1.5rem;
    }

    .filter-bar {
        width: 100%;
        overflow-x: auto;
        padding-bottom: 0.5rem;
        scroll-behavior: smooth;
        -webkit-overflow-scrolling: touch;
    }

    .filter-btn {
        flex-shrink: 0;
    }

    .order-card {
        padding: 1.2rem;
        margin-bottom: 1rem;
    }

    .product-image-container {
        width: 70px;
        height: 70px;
    }

    .product-name {
        font-size: 1rem;
    }

    .price-amount {
        font-size: 1.3rem;
    }

    .shipping-value.address {
        max-width: 200px;
    }

    .action-btn {
        padding: 0.4rem 1rem;
        font-size: 0.85rem;
    }
}
</style>
