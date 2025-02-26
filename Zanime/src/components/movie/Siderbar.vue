<template>
    <div :class="['sidebar', { 'sidebar--collapsed': isCollapsed }]">
        <div class="sidebar-header d-flex align-items-center justify-content-between mb-3">
            <button class="collapse-btn btn btn-dark" @click="toggleSidebar">
                <span v-if="isCollapsed">+</span>
                <span v-else>-</span>
            </button>
            <h2 v-if="!isCollapsed" class="sidebar-title">动漫周边</h2>
        </div>
        <div class="sidebar-content" v-show="!isCollapsed" @scroll="handleScroll">
            <div v-if="isLoading && !products.length" class="loading-state">
                <div class="spinner-border text-light" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
                <p class="text-light mt-3">加载商品中...</p>
            </div>
            <ul v-else class="list-unstyled m-0">
                <li v-for="item in products" :key="item.product_id" @click="DisplayProduct(formatProduct(item))"
                    class="product-item-enter" ref="productItems">
                    <div class="product-item d-flex align-items-center mb-3 p-3 rounded shadow-sm">
                        <img :src="item.image" :alt="item.product_name" class="product-image img-fluid mr-3" />
                        <div class="product-info d-flex flex-column">
                            <span class="product-name">{{ item.product_name }}</span>
                            <span class="product-price">¥{{ getMinPrice(item.prices) }} ~ ¥{{ getMaxPrice(item.prices)
                                }}</span>
                        </div>
                    </div>
                </li>
                <div v-if="isLoadingMore" class="loading-more text-center py-3">
                    <div class="spinner-border spinner-border-sm text-light" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                    <p class="text-light mt-2 mb-0">加载更多...</p>
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
@import '@/static/css/siderbar.css';
</style>
