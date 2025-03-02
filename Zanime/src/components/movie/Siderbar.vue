<template>
    <div :class="['sidebar', { 'sidebar--collapsed': isCollapsed }]">
        <div class="sidebar-header">
            <div class="header-content">
                <h2 v-if="!isCollapsed" class="sidebar-title">
                    <i class="fas fa-gift"></i> 动漫周边
                </h2>
                <button class="collapse-btn" @click="toggleSidebar">
                    <i class="fas" :class="isCollapsed ? 'fa-chevron-right' : 'fa-chevron-left'"></i>
                </button>
            </div>
            <div class="header-gradient"></div>
        </div>
        <div class="sidebar-content" v-show="!isCollapsed" @scroll="handleScroll">
            <div v-if="isLoading && !products.length" class="loading-state">
                <div class="spinner">
                    <div class="bounce1"></div>
                    <div class="bounce2"></div>
                    <div class="bounce3"></div>
                </div>
                <p>加载商品中...</p>
            </div>
            <ul v-else class="product-list">
                <li v-for="item in products" :key="item.product_id" @click="DisplayProduct(formatProduct(item))"
                    class="product-item-container" ref="productItems">
                    <div class="product-item">
                        <div class="product-image-wrapper">
                            <img :src="item.image" :alt="item.product_name" class="product-image" />
                        </div>
                        <div class="product-info">
                            <h4 class="product-name">{{ item.product_name }}</h4>
                            <div class="product-price">
                                <span class="price-label">价格:</span>
                                <span class="price-range">¥{{ getMinPrice(item.prices) }} ~ ¥{{ getMaxPrice(item.prices)
                                }}</span>
                            </div>
                        </div>
                    </div>
                </li>
                <div v-if="isLoadingMore" class="loading-more">
                    <div class="spinner small">
                        <div class="bounce1"></div>
                        <div class="bounce2"></div>
                        <div class="bounce3"></div>
                    </div>
                    <p>加载更多...</p>
                </div>
            </ul>
        </div>
    </div>
    <ProductDetail :isModalVisible="isModalVisible" :product="display || null" @closeVisible="closeVisible" />
</template>

<script>
// 导入所需的依赖
import axios from 'axios';
import ProductDetail from './ProductDetail.vue';

export default {
    // 组件数据
    data() {
        return {
            isCollapsed: true, // 侧边栏是否折叠
            isModalVisible: false, // 商品详情模态框是否可见
            display: null, // 当前显示的商品
            products: [], // 商品列表
            isLoading: true, // 是否正在加载
            isLoadingMore: false, // 是否正在加载更多
            page: 1, // 当前页码
            pageSize: 15, // 每页显示数量
            hasMore: true, // 是否还有更多数据
            loadingThreshold: 10 // 触发加载更多的阈值
        };
    },
    // 组件挂载时执行
    mounted() {
        this.fetchProducts();
    },
    methods: {
        // 获取商品数据
        async fetchProducts(loadMore = false) {
            if (!loadMore) {
                this.isLoading = true;
            } else {
                this.isLoadingMore = true;
            }

            try {
                // 发送API请求获取商品数据
                const response = await axios.get('/api/get-products', {
                    params: {
                        page: this.page,
                        pageSize: this.pageSize
                    }
                });

                // 处理返回的数据
                if (loadMore) {
                    this.products = [...this.products, ...response.data.products];
                } else {
                    this.products = response.data.products;
                }

                // 判断是否还有更多数据
                this.hasMore = response.data.products.length === this.pageSize;

                // 更新加载状态
                if (loadMore) {
                    this.isLoadingMore = false;
                } else {
                    this.isLoading = false;
                }
            } catch (error) {
                console.error('Error fetching products:', error);
                if (loadMore) {
                    this.isLoadingMore = false;
                } else {
                    this.isLoading = false;
                }
            }
        },

        // 处理滚动事件，实现无限加载
        handleScroll(e) {
            if (this.isLoading || this.isLoadingMore || !this.hasMore) return;

            const container = e.target;
            const items = this.$refs.productItems || [];

            if (items.length === 0) return;

            const triggerItem = items[items.length - this.loadingThreshold];
            if (!triggerItem) return;

            const triggerItemPosition = triggerItem.offsetTop;
            const containerScrollPosition = container.scrollTop + container.clientHeight;

            // 当滚动到触发点时加载更多数据
            if (containerScrollPosition > triggerItemPosition) {
                this.page += 1;
                this.fetchProducts(true);
            }
        },

        // 切换侧边栏折叠状态
        toggleSidebar() {
            this.isCollapsed = !this.isCollapsed;
        },

        // 显示商品详情
        DisplayProduct(item) {
            if (!item || !item.stocks || !item.prices) {
                console.warn('Invalid product data:', item);
                return;
            }
            this.isModalVisible = true;
            this.display = item;
        },

        // 关闭商品详情
        closeVisible() {
            this.isModalVisible = false;
            this.display = null;
        },

        // 格式化商品数据
        formatProduct(item) {
            // 验证商品数据的完整性
            if (!item ||
                !Array.isArray(item.sizes) ||
                !Array.isArray(item.colors) ||
                !Array.isArray(item.prices) ||
                !Array.isArray(item.stocks) ||
                item.sizes.length === 0 ||
                item.colors.length === 0 ||
                item.prices.length === 0 ||
                item.stocks.length === 0) {
                console.warn('Invalid product data:', item);
                return null;
            }

            // 获取唯一的尺寸和颜色
            const uniqueSizes = [...new Set(item.sizes)];
            const uniqueColors = [...new Set(item.colors)];

            // 返回格式化后的商品对象
            return {
                productID: item.product_id || 0,
                productName: item.product_name || '未知商品',
                description: item.description || '暂无描述',
                image: item.image || '',
                prices: Array.isArray(item.prices) ? item.prices : [],
                stocks: Array.isArray(item.stocks) ? item.stocks : [],
                sizes: item.sizes,
                colors: item.colors,
                category: item.category || '未分类',
                uniqueSizes: uniqueSizes,
                uniqueColors: uniqueColors,
                // 获取指定尺寸和颜色的价格和库存
                getPriceAndStock(size, color) {
                    if (!size || !color || !this.sizes || !this.colors) {
                        return { price: 0, stock: 0 };
                    }

                    for (let i = 0; i < this.sizes.length; i++) {
                        if (this.sizes[i] === size && this.colors[i] === color) {
                            return {
                                price: this.prices[i] || 0,
                                stock: this.stocks[i] || 0
                            };
                        }
                    }
                    return { price: 0, stock: 0 };
                }
            };
        },

        // 获取最低价格
        getMinPrice(prices) {
            if (!Array.isArray(prices) || prices.length === 0) {
                return '0.00';
            }
            return Math.min(...prices).toFixed(2);
        },

        // 获取最高价格
        getMaxPrice(prices) {
            if (!Array.isArray(prices) || prices.length === 0) {
                return '0.00';
            }
            return Math.max(...prices).toFixed(2);
        }
    },
    // 注册子组件
    components: {
        ProductDetail,
    }
};
</script>

<style scoped>
.sidebar {
    position: fixed;
    top: 50%;
    right: 0;
    transform: translateY(-50%);
    width: 320px;
    height: 80vh;
    background: rgba(20, 20, 20, 0.95);
    border-radius: 12px 0 0 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: -5px 0 25px rgba(0, 0, 0, 0.4);
    overflow: hidden;
    z-index: 1000;
    backdrop-filter: blur(10px);
    border-left: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar--collapsed {
    width: 50px;
    height: 120px;
    border-radius: 8px 0 0 8px;
}

.sidebar-header {
    position: relative;
    padding: 0;
}

.header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    position: relative;
    z-index: 2;
}

.header-gradient {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 60px;
    background: linear-gradient(135deg, #AB0511, #8C0511);
    z-index: 1;
    opacity: 0.85;
}

.sidebar-title {
    margin: 0;
    color: white;
    font-size: 1.2rem;
    font-weight: 600;
    white-space: nowrap;
    position: relative;
    z-index: 2;
}

.sidebar-title i {
    margin-right: 8px;
    color: #f8f8f8;
}

.collapse-btn {
    background: rgba(255, 255, 255, 0);
    color: white;
    border: none;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    position: relative;
    z-index: 2;
    padding: 0;
}

.collapse-btn:hover {
    transform: scale(1.05);
}

.sidebar-content {
    height: calc(80vh - 60px);
    overflow-y: auto;
    padding: 10px;
    scrollbar-width: thin;
    scrollbar-color: #444 #222;
}

.sidebar-content::-webkit-scrollbar {
    width: 6px;
}

.sidebar-content::-webkit-scrollbar-track {
    background: #222;
    border-radius: 3px;
}

.sidebar-content::-webkit-scrollbar-thumb {
    background-color: #444;
    border-radius: 3px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
    background-color: #555;
}

.product-list {
    list-style: none;
    margin: 0;
    padding: 0;
}

.product-item-container {
    margin-bottom: 15px;
    transition: transform 0.3s ease, opacity 0.3s ease;
    animation: fadeIn 0.5s ease-in-out;
}

.product-item-container:hover {
    transform: translateY(-3px);
}

.product-item {
    background: rgba(40, 40, 40, 0.7);
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.product-item:hover {
    background: rgba(50, 50, 50, 0.8);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.product-image-wrapper {
    width: 100%;
    height: 150px;
    overflow: hidden;
    position: relative;
}

.product-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.5s ease;
}

.product-item:hover .product-image {
    transform: scale(1.05);
}

.product-info {
    padding: 12px;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
}

.product-name {
    margin: 0 0 8px 0;
    font-size: 1.25rem;
    font-weight: 500;
    color: white;
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.product-price {
    display: flex;
    align-items: center;
    margin-top: auto;
}

.price-label {
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.7);
    margin-right: 5px;
}

.price-range {
    font-size: 1rem;
    font-weight: 600;
    color: #ffffff;
}

/* 加载动画 */
.loading-state,
.loading-more {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
    color: rgba(255, 255, 255, 0.8);
}

.loading-more {
    padding: 10px;
}

.loading-more p {
    margin: 5px 0 0 0;
    font-size: 0.8rem;
}

.spinner {
    display: flex;
    align-items: center;
    justify-content: center;
}

.spinner>div {
    width: 10px;
    height: 10px;
    margin: 0 3px;
    background-color: #AB0511;
    border-radius: 100%;
    display: inline-block;
    animation: bounce 1.4s infinite ease-in-out both;
}

.spinner.small>div {
    width: 6px;
    height: 6px;
}

.spinner .bounce1 {
    animation-delay: -0.32s;
}

.spinner .bounce2 {
    animation-delay: -0.16s;
}

@keyframes bounce {

    0%,
    80%,
    100% {
        transform: scale(0);
    }

    40% {
        transform: scale(1.0);
    }
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* 响应式调整 */
@media (max-width: 768px) {
    .sidebar {
        width: 280px;
        height: 70vh;
    }

    .sidebar--collapsed {
        width: 40px;
        height: 100px;
    }

    .sidebar-content {
        height: calc(70vh - 60px);
    }

    .product-image-wrapper {
        height: 120px;
    }

    .product-name {
        font-size: 0.85rem;
    }
}
</style>
