<template>
    <div class="order-page">
        <!-- 页面标题 -->
        <div class="page-header">
            <h1>我的订单</h1>
        </div>

        <!-- 订单列表 -->
        <div class="order-container">
            <!-- 加载状态 -->
            <div v-if="isLoading" class="loading-state">
                <div class="spinner-border text-light" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
                <p class="text-light mt-3">加载订单中...</p>
            </div>

            <!-- 空状态 -->
            <div v-else-if="!orders.length" class="empty-state">
                <i class="fas fa-box-open mb-3"></i>
                <p>暂无订单</p>
            </div>

            <!-- 订单列表 -->
            <div v-else class="order-list">
                <div v-for="order in orders" :key="order.order_id" class="order-card">
                    <div class="order-header">
                        <div class="order-info">
                            <span class="order-id">订单号：{{ order.order_id }}</span>
                            <span class="order-date">{{ formatDate(order.create_time) }}</span>
                        </div>
                        <div class="order-status" :class="getStatusClass(order.status)">
                            {{ order.status }}
                        </div>
                    </div>

                    <div class="order-content">
                        <div class="product-info">
                            <img :src="order.image" :alt="order.product_name" class="product-image">
                            <div class="product-details">
                                <h5 class="product-name">{{ order.product_name }}</h5>
                                <p class="product-specs">
                                    规格：{{ order.selected_size }} / {{ order.selected_color }}
                                </p>
                            </div>
                        </div>
                        <div class="order-price">
                            <span class="price">¥{{ order.price }}</span>
                        </div>
                    </div>

                    <div class="order-footer">
                        <div class="shipping-info">
                            <p>收货人：{{ order.user_name }}</p>
                            <p>电话：{{ order.phone }}</p>
                            <p>地址：{{ order.address }}</p>
                        </div>
                        <div class="order-actions">
                            <button v-if="order.status === '待支付'" class="btn btn-primary" @click="handlePayment(order)">
                                继续支付
                            </button>
                            <button v-else-if="order.status === '已完成'" class="btn btn-outline-primary"
                                @click="handleReview(order)">
                                评价
                            </button>
                        </div>
                    </div>
                </div>

                <!-- 加载更多 -->
                <div v-if="hasMore" class="load-more text-center py-4">
                    <button class="btn btn-outline-light" @click="loadMore" :disabled="isLoadingMore">
                        <span v-if="isLoadingMore">
                            <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                            加载中...
                        </span>
                        <span v-else>加载更多</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

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
                    const newOrders = response.data.orders;

                    if (loadMore) {
                        this.orders = [...this.orders, ...newOrders];
                    } else {
                        this.orders = newOrders;
                    }

                    this.hasMore = newOrders.length === this.pageSize;
                } else {
                    console.error('Error fetching orders: response.data.code == ', response.data.code);
                }
            } catch (error) {
                console.error('Error fetching orders:', error);
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
    background: linear-gradient(145deg, #0e0e0e 0%, #272727 100%);
    color: #fff;
}

.page-header {
    margin-bottom: 2rem;
}

.page-header h1 {
    font-size: 2rem;
    font-weight: 600;
    color: #fff;
}

.order-card {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 1rem;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.order-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.order-info {
    display: flex;
    gap: 1rem;
    color: rgba(255, 255, 255, 0.7);
}

.order-status {
    padding: 0.25rem 1rem;
    border-radius: 2rem;
    font-size: 0.875rem;
    font-weight: 500;
}

.status-pending {
    background-color: rgba(255, 215, 0, 0.1);
    background-color: transparent;
    /* 背景是无色透明的 */
    color: #d4dfd5;
}

.status-paid {
    background-color: rgba(76, 175, 80, 0.1);
    background-color: transparent;
    /* 背景是无色透明的 */
    color: #d4dfd5;
}

.status-completed {
    background-color: rgba(33, 150, 243, 0.1);
    background-color: transparent;
    /* 背景是无色透明的 */
    color: #d4dfd5;
}

.status-expired {
    background-color: transparent;
    /* 背景是无色透明的 */
    color: #d4dfd5;
}

.order-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.product-info {
    display: flex;
    gap: 1rem;
    align-items: center;
}

.product-image {
    width: 80px;
    height: 80px;
    border-radius: 0.5rem;
    object-fit: cover;
}

.product-details {
    flex: 1;
}

.product-name {
    margin-bottom: 0.5rem;
    font-size: 1.1rem;
}

.product-specs {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
    margin: 0;
}

.order-price {
    font-size: 1.25rem;
    font-weight: 600;
    color: #dad2d2;
}

.order-footer {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.shipping-info {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
}

.shipping-info p {
    margin: 0.25rem 0;
}

.order-actions {
    display: flex;
    gap: 1rem;
}

.empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: rgba(255, 255, 255, 0.6);
}

.empty-state i {
    font-size: 3rem;
    margin-bottom: 1rem;
}

.loading-state {
    text-align: center;
    padding: 4rem 2rem;
}

.btn-outline-light {
    background-color: transparent !important;
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-outline-light:hover,
.btn-outline-light.active {
    background-color: rgba(255, 255, 255, 0.1) !important;
    border-color: rgba(255, 255, 255, 0.3);
}

/* 修改按钮样式 */
.btn-dark {
    background-color: transparent !important;
    border: 1px solid rgba(0, 0, 0, 0.2);
    color: #000;
}

.btn-dark:hover {
    background-color: rgba(0, 0, 0, 0.05) !important;
    border-color: rgba(0, 0, 0, 0.3);
}

.btn-outline-secondary {
    background-color: transparent !important;
    border: 1px solid rgba(0, 0, 0, 0.2);
    color: #666;
}

.btn-outline-secondary:hover {
    background-color: rgba(0, 0, 0, 0.05) !important;
    border-color: rgba(0, 0, 0, 0.3);
    color: #333;
}

/* 加载更多按钮 */
.load-more .btn-outline-light {
    background-color: transparent !important;
    border: 1px solid rgba(255, 255, 255, 0.2);
    padding: 0.5rem 2rem;
}

.load-more .btn-outline-light:hover {
    background-color: rgba(255, 255, 255, 0.1) !important;
    border-color: rgba(255, 255, 255, 0.3);
}

/* 响应式调整 */
@media (max-width: 768px) {
    .order-page {
        padding: 1rem;
    }

    .order-content {
        flex-direction: column;
        align-items: flex-start;
        gap: 1rem;
    }

    .order-price {
        align-self: flex-end;
    }

    .order-footer {
        flex-direction: column;
        gap: 1rem;
    }

    .order-actions {
        align-self: flex-end;
    }
}

/* 移动端样式优化 */
@media (max-width: 767px) {
    .order-page {
        padding: 1.5rem;
    }

    .page-header h1 {
        font-size: 1.5rem;
    }

    .order-card {
        padding: 1rem;
    }

    .order-header {
        flex-direction: column;
        gap: 0.5rem;
        align-items: flex-end;
        /* 改为右对齐 */
    }

    .order-info {
        width: 100%;
        text-align: right;
        /* 添加右对齐 */
    }

    .order-info .order-id,
    .order-info .order-date {
        text-align: right;
        width: 100%;
        display: block;
    }

    .order-status {
        align-self: flex-end;
        /* 添加右对齐 */
    }

    .order-content {
        flex-direction: column;
        gap: 1rem;
    }

    .product-info {
        width: 100%;
    }

    .product-image {
        width: 70px;
        height: 70px;
    }

    .order-price {
        align-self: flex-end;
    }

    .order-footer {
        flex-direction: column;
        gap: 1rem;
    }

    .shipping-info {
        width: 100%;
    }

    .shipping-info p {
        display: flex;
        justify-content: space-between;
        margin: 0.5rem 0;
        text-align: right;
        /* 添加右对齐 */
    }

    .shipping-info p::before {
        content: attr(data-label);
        color: rgba(255, 255, 255, 0.5);
    }

    .order-actions {
        align-self: flex-end;
    }
}

/* 响应式设计优化 */
/* 大屏幕设备 (1200px 及以上) */
@media (min-width: 1200px) {
    .order-page {
        padding: 2.5rem;
        max-width: 1400px;
        margin: 0 auto;
    }

    .order-card {
        padding: 2rem;
    }

    .product-image {
        width: 100px;
        height: 100px;
    }
}

/* 中等屏幕设备 (768px 到 1199px) */
@media (max-width: 1199px) {
    .order-page {
        padding: 2rem;
    }
}

/* 平板设备 (576px 到 767px) */
@media (max-width: 767px) {
    .order-page {
        padding: 1.5rem;
    }

    .page-header h1 {
        font-size: 1.5rem;
    }

    .order-card {
        padding: 1rem;
    }

    .order-header {
        flex-direction: column;
        gap: 0.5rem;
        align-items: flex-start;
    }

    .order-content {
        flex-direction: column;
        gap: 1rem;
    }

    .product-info {
        width: 100%;
    }

    .product-image {
        width: 70px;
        height: 70px;
    }

    .order-price {
        align-self: flex-end;
    }

    .order-footer {
        flex-direction: column;
        gap: 1rem;
    }

    .shipping-info {
        width: 100%;
    }

    .order-actions {
        width: 100%;
        justify-content: flex-end;
    }
}

/* 手机设备 (575px 及以下) */
@media (max-width: 575px) {
    .order-page {
        padding: 1rem;
    }

    .page-header h1 {
        font-size: 1.25rem;
    }

    .order-card {
        padding: 1rem;
        margin-bottom: 1rem;
    }

    .order-info {
        flex-direction: column;
        gap: 0.25rem;
    }

    .product-info {
        flex-direction: column;
        text-align: center;
    }

    .product-image {
        width: 100%;
        height: auto;
        max-width: 200px;
        margin: 0 auto;
    }

    .product-details {
        text-align: center;
    }

    .order-actions {
        flex-direction: column;
        gap: 0.5rem;
    }

    .btn {
        width: 100%;
        margin-bottom: 0.5rem;
    }

    .shipping-info p {
        font-size: 0.85rem;
    }
}

/* 超小屏幕设备优化 */
@media (max-width: 360px) {
    .order-page {
        padding: 0.5rem;
    }

    .order-card {
        padding: 0.75rem;
        border-radius: 0.5rem;
    }

    .product-name {
        font-size: 1rem;
    }

    .product-specs {
        font-size: 0.8rem;
    }

    .order-price {
        font-size: 1.1rem;
    }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
    .order-card {
        background: rgba(255, 255, 255, 0.03);
    }

    .btn-outline-light {
        border-color: rgba(255, 255, 255, 0.1);
    }
}
</style>
