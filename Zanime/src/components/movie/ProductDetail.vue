<template>
    <div>
        <!-- 背景遮罩 -->
        <div class="modal-backdrop fade" :class="{ show: isModalVisible }" v-if="isModalVisible"></div>

        <!-- 商品详情弹窗 -->
        <div class="modal fade" :class="{ show: isModalVisible }" tabindex="-1" style="display: block"
            v-if="isModalVisible">
            <div class="modal-dialog modal-dialog-centered modal-size">
                <div class="modal-content">
                    <!-- 模态框头部 -->
                    <div class="modal-header border-bottom-0">
                        <h5 class="modal-title fw-bold">{{ product.productName }}</h5>
                        <button v-if="!isLoading" type="button" class="btn-close" @click="closeModal">
                        </button>
                    </div>

                    <!-- 模态框主体 -->
                    <div class="modal-body">
                        <!-- 商品详情视图 -->
                        <div v-if="!isPaymentView && !isPaymentSuccess" class="product-container">
                            <!-- 左侧商品图片 -->
                            <div class="product-image-section">
                                <div class="product-image-container mb-3">
                                    <img :src="product.image" :alt="product.productName" class="img-fluid rounded">
                                </div>

                                <!-- 评价部分 -->
                                <div class="product-rating-container">
                                    <h6 class="text-muted mb-3">商品评价</h6>

                                    <!-- 总体评分 -->
                                    <div class="rating-summary mb-3">
                                        <div class="rating-score">
                                            <span class="score">{{ averageRating }}</span>
                                            <span class="total">/5</span>
                                        </div>
                                        <div class="rating-stars">
                                            <i v-for="n in 5" :key="n" class="fas fa-star"
                                                :class="{ 'active': n <= averageRating }">
                                            </i>
                                        </div>
                                        <div class="rating-count">
                                            {{ product.ratings?.length || 0 }} 条评价
                                        </div>
                                    </div>

                                    <!-- 评价列表 -->
                                    <div class="rating-list">
                                        <div v-for="(rating, index) in displayRatings" :key="index" class="rating-item">
                                            <div class="rating-header">
                                                <div class="rating-user">{{ rating.userName }}</div>
                                                <div class="rating-date">{{ rating.date }}</div>
                                            </div>
                                            <div class="rating-stars mb-2">
                                                <i v-for="n in 5" :key="n" class="fas fa-star"
                                                    :class="{ 'active': n <= rating.score }">
                                                </i>
                                            </div>
                                            <div class="rating-content">{{ rating.comment }}</div>
                                        </div>
                                    </div>

                                    <!-- 查看更多按钮 -->
                                    <button v-if="hasMoreRatings" class="btn btn-link text-light mt-3"
                                        @click="showMoreRatings">
                                        查看更多评价
                                    </button>
                                </div>
                            </div>

                            <!-- 右侧商品详情 -->
                            <div class="product-info-section">
                                <div class="product-details">
                                    <div class="product-info">
                                        <div class="mb-4">
                                            <h6 class="text-muted mb-2">商品描述</h6>
                                            <p class="mb-0">{{ product.description }}</p>
                                        </div>

                                        <div class="mb-4">
                                            <h6 class="text-muted mb-2">价格</h6>
                                            <p class="fs-4 text-primary fw-bold mb-0">
                                                {{ selectedPrice ? `¥${selectedPrice}` : `¥${getMinPrice} ~
                                                ¥${getMaxPrice}` }}
                                            </p>
                                        </div>

                                        <div class="mb-4">
                                            <h6 class="text-muted mb-2">库存状态</h6>
                                            <span class="badge" :class="[
                                                currentStock > 0 ? 'bg-success' : 'bg-danger',
                                                (!selectedSize && !selectedColor) ? 'bg-info' : ''
                                            ]">
                                                {{ getStockStatusText }}
                                            </span>
                                        </div>

                                        <!-- 尺寸选择 -->
                                        <div class="mb-4" v-if="product.uniqueSizes">
                                            <h6 class="text-muted mb-2">选择尺寸</h6>
                                            <div class="btn-group">
                                                <button v-for="size in product.uniqueSizes" :key="size"
                                                    class="btn btn-outline-secondary" :class="{
                                                        'active': selectedSize === size,
                                                        'is-invalid': errors.size
                                                    }" @click="handleSizeSelect(size)">
                                                    {{ size }}
                                                </button>
                                            </div>
                                            <div class="invalid-feedback" v-if="errors.size">{{ errors.size }}</div>
                                        </div>

                                        <!-- 颜色选择 -->
                                        <div class="mb-4" v-if="product.uniqueColors">
                                            <h6 class="text-muted mb-2">选择颜色</h6>
                                            <div class="btn-group">
                                                <button v-for="color in product.uniqueColors" :key="color"
                                                    class="btn btn-outline-secondary" :class="{
                                                        'active': selectedColor === color,
                                                        'is-invalid': errors.color
                                                    }" @click="handleColorSelect(color)">
                                                    {{ color }}
                                                </button>
                                            </div>
                                            <div class="invalid-feedback" v-if="errors.color">{{ errors.color }}</div>
                                        </div>

                                        <!-- 添加分隔线和联系信息表单 -->
                                        <div class="contact-section mt-4 pt-4 border-top border-secondary">
                                            <h6 class="text-muted mb-3">收货信息</h6>

                                            <!-- 联系电话 -->
                                            <div class="custom-form-group mb-3">
                                                <label class="custom-form-label">联系电话</label>
                                                <input type="tel" class="custom-form-input"
                                                    :class="{ 'is-invalid': errors.phone }" v-model="contactInfo.phone"
                                                    placeholder="请输入手机号码" ref="phoneInput">
                                                <div class="invalid-feedback" v-if="errors.phone">{{ errors.phone }}
                                                </div>
                                            </div>

                                            <!-- 收货人 -->
                                            <div class="custom-form-group mb-3">
                                                <label class="custom-form-label">收货人</label>
                                                <input type="text" class="custom-form-input"
                                                    :class="{ 'is-invalid': errors.name }" v-model="contactInfo.name"
                                                    placeholder="请输入收货人姓名" ref="nameInput">
                                                <div class="invalid-feedback" v-if="errors.name">{{ errors.name }}</div>
                                            </div>

                                            <!-- 收货地址 -->
                                            <div class="custom-form-group mb-3">
                                                <label class="custom-form-label">收货地址</label>
                                                <textarea class="custom-form-input"
                                                    :class="{ 'is-invalid': errors.address }"
                                                    v-model="contactInfo.address" rows="2" placeholder="请输入详细收货地"
                                                    ref="addressInput"></textarea>
                                                <div class="invalid-feedback" v-if="errors.address">{{ errors.address }}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- 支付视图 -->
                        <div v-else-if="isPaymentView && !isPaymentSuccess" class="payment-container">
                            <!-- 加载状态 -->
                            <div v-if="isLoading" class="loading-overlay">
                                <div class="spinner-border text-light" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                                <p class="mt-3 text-light">订单创建中...</p>
                            </div>

                            <!-- 支付信息展示 -->
                            <div v-else class="payment-content">
                                <div class="order-info mb-4">
                                    <h6 class="text-muted mb-3">订单信息</h6>
                                    <div class="order-details">
                                        <div class="info-item">
                                            <span class="label">订单编号：</span>
                                            <span class="value">{{ orderInfo.orderID }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">商品名称：</span>
                                            <span class="value">{{ orderInfo.productName }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">尺码：</span>
                                            <span class="value">{{ orderInfo.selectedSize }};{{ orderInfo.selectedColor
                                                }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">收货人：</span>
                                            <span class="value">{{ orderInfo.userName }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">联系电话：</span>
                                            <span class="value">{{ orderInfo.phone }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">收货地址：</span>
                                            <span class="value">{{ orderInfo.address }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">支付金额：</span>
                                            <span class="value price">¥{{ orderInfo.price }}</span>
                                        </div>
                                        <div class="info-item">
                                            <span class="label">剩余支付时间：</span>
                                            <span class="value countdown" :class="{
                                                'text-danger': remainingTime < 300,
                                                'text-muted': remainingTime <= 0
                                            }">
                                                {{ formatRemainingTime }}
                                            </span>
                                        </div>
                                    </div>
                                </div>

                                <div class="qr-code-section text-center">
                                    <h6 class="text-muted mb-3">请扫码支付</h6>
                                    <div class="qr-code-wrapper" :class="{ 'expired': remainingTime <= 0 }"
                                        @click="handleQrCodeClick">
                                        <img src="@/static/image/收款码.jpg" alt="支付二维码" class="qr-code">
                                        <div v-if="remainingTime <= 0" class="expired-overlay">
                                            <i class="fas fa-ban"></i>
                                            <span>二维码已失效</span>
                                        </div>
                                    </div>
                                    <p class="text-muted mt-2">
                                        <small v-if="remainingTime > 0">请使用微信或支付宝扫描二维码完成支付</small>
                                        <small v-else class="text-danger">订单已超时，请重新下单</small>
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- 支付成功视图 -->
                        <div v-else class="success-container">
                            <div class="success-content text-center">
                                <div class="success-icon mb-4">
                                    <i class="fas fa-check-circle"></i>
                                </div>
                                <h4 class="mb-3">支付成功</h4>
                                <p class="text-muted mb-4">
                                    订单号：{{ orderInfo.orderID }}<br>
                                    支付金额：¥{{ orderInfo.price }}
                                </p>
                                <button class="btn btn-primary" @click="handlePaymentSuccess">
                                    确认
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- 模态框底部 -->
                    <div class="modal-footer border-top-0">
                        <button v-if="!isPaymentView && !isPaymentSuccess" type="button" class="btn btn-outline-primary"
                            @click="addToCart" :disabled="currentStock <= 0">
                            加入购物车
                        </button>
                        <button v-if="!isPaymentView && !isPaymentSuccess" type="button" class="btn btn-primary"
                            @click="buyNow" :disabled="currentStock <= 0">
                            立即购买
                        </button>
                        <button v-if="isPaymentView && !isPaymentSuccess && !isLoading" type="button"
                            class="btn btn-secondary" @click="backToProduct">
                            返回商品
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 添加 Toast 提示组件 -->
        <div class="toast-container position-fixed top-0 start-50 translate-middle-x">
            <div class="toast align-items-center text-white bg-danger border-0" role="alert"
                :class="{ 'show': showErrorToast }" aria-live="assertive" aria-atomic="true">
                <div class="d-flex">
                    <div class="toast-body">
                        <i class="fas fa-exclamation-circle me-2"></i>
                        {{ errorMessage }}
                    </div>
                    <button type="button" class="btn-close btn-close-white me-2 m-auto" @click="showErrorToast = false">
                    </button>
                </div>
            </div>
        </div>

        <!-- 修改确认离开的模态框 -->
        <div v-if="showLeaveConfirm" class="modal fade show" style="display: block; background: rgba(0, 0, 0, 0.5);">
            <div class="modal-dialog modal-dialog-centered leave-confirm-dialog">
                <div class="modal-content">
                    <div class="modal-body text-center py-4">
                        <p>订单尚未支付，确定要离开吗？</p>
                        <div class="mt-4">
                            <button type="button" class="btn btn-secondary me-3" @click="cancelLeave">
                                返回支付
                            </button>
                            <button type="button" class="btn btn-danger" @click="confirmLeave">
                                确认离开
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'ProductDetail',
    props: {
        // 控制模态框显示/隐藏
        isModalVisible: {
            type: Boolean,
            default: false
        },
        // 商品详情数据
        product: {
            type: Object,
            required: false,
            default: () => ({
                productName: '',
                description: '',
                image: '',
                prices: [],
                stocks: [],
                sizes: [],
                colors: [],
                ratings: []
            })
        }
    },
    data() {
        return {
            selectedSize: null,        // 选中的尺寸
            selectedColor: null,       // 选中的颜色
            contactInfo: {             // 联系信息
                phone: '',
                name: '',
                address: ''
            },
            displayCount: 3,           // 评价显示数量
            errors: {                  // 表单错误信息
                phone: '',
                name: '',
                address: '',
                size: '',
                color: ''
            },
            isPaymentView: false,      // 是否显示支付视图
            isLoading: false,          // 加载状态
            orderInfo: null,           // 订单信息
            isPaymentSuccess: false,   // 支付是否成功
            selectedPrice: null,       // 选中规格的价格
            currentStock: 0,           // 当前库存
            showErrorToast: false,     // 是否显示错误提示
            errorMessage: '',          // 错误信息
            remainingTime: 0,          // 支付剩余时间
            timer: null,               // 倒计时定时器
            showLeaveConfirm: false,   // 是否显示离开确认框
            pendingAction: null        // 待执行的操作
        }
    },
    watch: {
        // 监听模态框显示状态
        isModalVisible(newVal) {
            if (newVal && this.product) {
                this.resetForm();
                this.$nextTick(() => {
                    if (this.product && Array.isArray(this.product.stocks)) {
                        this.updatePriceAndStock();
                    }
                });
            }
        },
        // 监听商品数据变化
        product: {
            immediate: true,
            handler(newProduct) {
                if (newProduct && Array.isArray(newProduct.stocks)) {
                    this.updatePriceAndStock();
                }
            }
        }
    },
    computed: {
        // 计算平均评分
        averageRating() {
            if (!this.product.ratings?.length) return 0;
            const sum = this.product.ratings.reduce((acc, curr) => acc + curr.score, 0);
            return (sum / this.product.ratings.length).toFixed(1);
        },
        // 获取要显示的评价列表
        displayRatings() {
            return this.product.ratings?.slice(0, this.displayCount) || [];
        },
        // 是否还有更多评价
        hasMoreRatings() {
            return this.product.ratings?.length > this.displayCount;
        },
        // 获取最低价格
        getMinPrice() {
            if (!this.product || !Array.isArray(this.product.prices) || this.product.prices.length === 0) {
                return '0.00';
            }
            return Math.min(...this.product.prices).toFixed(2);
        },
        // 获取最高价格
        getMaxPrice() {
            if (!this.product || !Array.isArray(this.product.prices) || this.product.prices.length === 0) {
                return '0.00';
            }
            return Math.max(...this.product.prices).toFixed(2);
        },
        // 获取库存状态文本
        getStockStatusText() {
            if (!this.selectedSize && !this.selectedColor) {
                return '请选择规格';
            }
            return this.currentStock > 0 ? `库存: ${this.currentStock}` : '暂时缺货';
        },
        // 格式化剩余支付时间
        formatRemainingTime() {
            if (this.remainingTime <= 0) {
                return '订单已超时';
            }
            const minutes = Math.floor(this.remainingTime / 60);
            const seconds = this.remainingTime % 60;
            return `${minutes}分${seconds < 10 ? '0' : ''}${seconds}秒`;
        }
    },
    methods: {
        // 关闭模态框
        closeModal() {
            if (this.isPaymentView && !this.isPaymentSuccess && this.remainingTime > 0) {
                this.showLeaveConfirm = true;
                this.pendingAction = 'close';
            } else {
                this.$emit('closeVisible');
            }
        },
        // 添加到购物车
        addToCart() {
            if (!this.validateForm()) return;
            const productInfo = {
                ...this.product,
                selectedSize: this.selectedSize,
                selectedColor: this.selectedColor,
                phone: this.contactInfo.phone,
                userName: this.contactInfo.name,
                address: this.contactInfo.address
            }
            this.$emit('add-to-cart', productInfo)
        },
        // 立即购买
        async buyNow() {
            if (!this.validateForm()) return;

            const productInfo = {
                productID: this.product.productID,
                productName: this.product.productName,
                selectedSize: this.selectedSize,
                selectedColor: this.selectedColor,
                price: parseFloat(this.selectedPrice) || parseFloat(this.getMinPrice()),
                phone: this.contactInfo.phone,
                userName: this.contactInfo.name,
                address: this.contactInfo.address
            };

            this.isPaymentView = true;
            this.isLoading = true;

            try {
                const response = await axios.post('/api/order', productInfo);
                if (response.data.code === 200) {
                    this.orderInfo = {
                        ...response.data,
                        price: this.selectedPrice || this.getMinPrice
                    };
                    this.startCountdown(response.data.createdTime);
                    this.isLoading = false;
                } else {
                    this.isLoading = false;
                    this.isPaymentView = false;
                    this.showError(response.data.message);
                }
            } catch (error) {
                console.error('Error placing order:', error);
                this.isLoading = false;
                this.isPaymentView = false;
                this.showError(error.response.data.message);
            }
        },
        // 更新价格和库存信息
        updatePriceAndStock() {
            if (!this.product || !Array.isArray(this.product.stocks)) {
                this.currentStock = 0;
                this.selectedPrice = null;
                return;
            }

            // 如果都没选择，显示总库存
            if (!this.selectedSize && !this.selectedColor) {
                this.currentStock = this.product.stocks.reduce((sum, stock) => sum + stock, 0);
                this.selectedPrice = null;
                return;
            }

            // 如果只选择了尺寸
            if (this.selectedSize && !this.selectedColor && Array.isArray(this.product.sizes)) {
                let sizeStock = 0;
                let sizePrices = new Set();
                this.product.sizes.forEach((size, index) => {
                    if (size === this.selectedSize) {
                        sizeStock += this.product.stocks[index] || 0;
                        if (this.product.prices[index]) {
                            sizePrices.add(this.product.prices[index]);
                        }
                    }
                });
                this.currentStock = sizeStock;
                this.selectedPrice = sizePrices.size === 1 ?
                    Array.from(sizePrices)[0].toFixed(2) : null;
                return;
            }

            // 如果只选择了颜色
            if (!this.selectedSize && this.selectedColor && Array.isArray(this.product.colors)) {
                let colorStock = 0;
                let colorPrices = new Set();
                this.product.colors.forEach((color, index) => {
                    if (color === this.selectedColor) {
                        colorStock += this.product.stocks[index] || 0;
                        if (this.product.prices[index]) {
                            colorPrices.add(this.product.prices[index]);
                        }
                    }
                });
                this.currentStock = colorStock;
                this.selectedPrice = colorPrices.size === 1 ?
                    Array.from(colorPrices)[0].toFixed(2) : null;
                return;
            }

            // 如果同时选择了尺寸和颜色
            if (this.product.getPriceAndStock) {
                const result = this.product.getPriceAndStock(this.selectedSize, this.selectedColor);
                if (result) {
                    this.selectedPrice = result.price.toFixed(2);
                    this.currentStock = result.stock;
                } else {
                    this.selectedPrice = null;
                    this.currentStock = 0;
                }
            }
        },
        // 清空表单错误信息
        clearErrors() {
            this.errors = {
                phone: '',
                name: '',
                address: '',
                size: '',
                color: ''
            }
        },
        // 验证表单
        validateForm() {
            this.clearErrors()
            let isValid = true

            // 验证手机号
            if (!this.contactInfo.phone) {
                this.errors.phone = '请输入联系电话'
                this.scrollToElement('phoneInput')
                isValid = false
            } else if (!/^1[3-9]\d{9}$/.test(this.contactInfo.phone)) {
                this.errors.phone = '请输入正确的手机号码'
                this.scrollToElement('phoneInput')
                isValid = false
            }

            // 验证收货人
            if (!this.contactInfo.name) {
                this.errors.name = '请输入收货人姓名'
                if (isValid) this.scrollToElement('nameInput')
                isValid = false
            }

            // 验证地址
            if (!this.contactInfo.address) {
                this.errors.address = '请输入收货地址'
                if (isValid) this.scrollToElement('addressInput')
                isValid = false
            }

            // 证尺寸选择（如果有）
            if (this.product.sizes && !this.selectedSize) {
                this.errors.size = '请选择尺寸'
                if (isValid) this.scrollToElement('sizeGroup')
                isValid = false
            }

            // 验证颜色选择（如果有）
            if (this.product.colors && !this.selectedColor) {
                this.errors.color = '请选择颜色'
                if (isValid) this.scrollToElement('colorGroup')
                isValid = false
            }

            return isValid
        },
        // 滚动到指定元素
        scrollToElement(refName) {
            const element = this.$refs[refName]
            if (element) {
                element.scrollIntoView({
                    behavior: 'smooth',
                    block: 'center'
                })
                // 添加闪烁效果
                element.classList.add('highlight-error')
                setTimeout(() => {
                    element.classList.remove('highlight-error')
                }, 2000)
            }
        },
        // 显示更多评价
        showMoreRatings() {
            this.displayCount += 3;
        },
        // 重置表单
        resetForm() {
            this.selectedSize = null;
            this.selectedColor = null;
            this.contactInfo = {
                phone: '',
                name: '',
                address: ''
            };
            this.clearErrors();
            this.displayCount = 3;
            this.selectedPrice = null;
            this.currentStock = 0;
            this.isPaymentView = false;
            this.isPaymentSuccess = false;
            this.orderInfo = null;
        },
        // 返回商品
        backToProduct() {
            if (this.remainingTime > 0) {
                this.showLeaveConfirm = true;
                this.pendingAction = 'back';
            } else {
                this.isPaymentView = false;
                this.orderInfo = null;
            }
        },
        // 模拟支付
        async SimulatePayment() {
            try {
                const response = await axios.post('/api/call-pay', {
                    orderID: this.orderInfo.orderID
                });

                if (response.data.code === 200) {  // 假设 200 表示支付成功
                    this.isPaymentSuccess = true;
                } else {
                    alert('支付失败，请重试');
                }
            } catch (error) {
                console.error('Error simulating payment:', error);
            }
        },
        // 处理支付成功
        handlePaymentSuccess() {
            this.isPaymentView = false;
            this.isPaymentSuccess = false;
            this.orderInfo = null;
            this.$emit('closeVisible');
        },
        // 处理尺寸选择
        handleSizeSelect(size) {
            this.selectedSize = size;
            this.updatePriceAndStock();
        },
        // 处理颜色选择
        handleColorSelect(color) {
            this.selectedColor = color;
            this.updatePriceAndStock();
        },
        // 显示错误
        showError(message) {
            this.errorMessage = message;
            this.showErrorToast = true;
            setTimeout(() => {
                this.showErrorToast = false;
            }, 3000); // 3秒后自动关闭
        },
        // 处理二维码点击
        handleQrCodeClick() {
            if (this.remainingTime <= 0) {
                this.showError('订单已超时，请重新下单');
                return;
            }
            this.SimulatePayment();
        },

        handleOrderTimeout() {
            this.showError('订单已超时，请重新下单');
        },

        startCountdown(createTime) {
            if (this.timer) {
                clearInterval(this.timer);
            }
            const createDateTime = new Date(createTime);
            const paymentDeadline = new Date(createDateTime.getTime() + (5 * 60 - 1) * 1000);
            this.remainingTime = Math.max(0, Math.floor((paymentDeadline - new Date()) / 1000));

            this.timer = setInterval(() => {
                if (this.remainingTime > 0) {
                    this.remainingTime--;
                } else {
                    clearInterval(this.timer);
                    this.handleOrderTimeout();
                }
            }, 1000);
        },

        // 取消离开
        cancelLeave() {
            this.showLeaveConfirm = false;
            this.pendingAction = null;
        },

        // 确认离开
        confirmLeave() {
            this.showLeaveConfirm = false;
            if (this.timer) {
                clearInterval(this.timer);
            }

            if (this.pendingAction === 'close') {
                this.$emit('closeVisible');
            } else if (this.pendingAction === 'back') {
                this.isPaymentView = false;
                this.orderInfo = null;
            }

            this.pendingAction = null;
        }
    },
    mounted() {
        if (this.product && Array.isArray(this.product.stocks)) {
            this.updatePriceAndStock();
        }
    },
    beforeDestroy() {
        if (this.timer) {
            clearInterval(this.timer);
        }
    }
}
</script>

<style scoped>
@import '@/static/css/productDetail.css';
</style>