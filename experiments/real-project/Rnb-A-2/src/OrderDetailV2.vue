<template>
  <div class="order-detail-v2">
    <van-nav-bar
      title="订单详情"
      left-arrow
      fixed
      placeholder
      :border="false"
      @click-left="goBack"
    />

    <div class="content">
      <!-- 通知条 -->
      <van-notice-bar class="notice" :scrollable="false" :text="order.notice" />

      <!-- 状态 + 地址 -->
      <section class="card">
        <div class="status-row">
          <van-tag class="status-tag" color="#fff3e8" text-color="#ff8f1f">{{ order.statusTag }}</van-tag>
          <span class="status-title">{{ order.statusTitle }}</span>
          <van-count-down class="countdown" :time="order.remainTime" format="剩余 HH:mm:ss" />
        </div>
        <div class="status-desc">{{ order.statusDesc }}</div>
        <div class="divider" />
        <van-cell
          class="address-cell"
          :title="`${order.address.name} ${order.address.phone}`"
          :label="order.address.detail"
          :border="false"
          is-link
          center
          @click="onAction('修改地址')"
        />
      </section>

      <!-- 商品信息 -->
      <section class="card">
        <h3 class="card-title">商品信息</h3>
        <div v-for="item in order.goods" :key="item.id" class="product">
          <div class="product-thumb">
            <img v-if="item.thumb" :src="item.thumb" alt="" />
          </div>
          <div class="product-info">
            <div class="product-title">{{ item.title }}</div>
            <div class="product-spec">{{ item.spec }}</div>
            <div class="product-bottom">
              <span class="price">
                <span class="price-symbol">¥</span><span class="price-int">{{ splitPrice(item.price)[0] }}</span><span class="price-dec">.{{ splitPrice(item.price)[1] }}</span>
              </span>
              <span class="product-count">x{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="divider" />
        <van-cell
          v-for="row in order.priceRows"
          :key="row.label"
          :title="row.label"
          :value="row.value"
          :border="false"
        />
        <div class="total-row">
          <span class="total-label">实付款</span>
          <span class="price price--large">
            <span class="price-symbol">¥</span><span class="price-int">{{ splitPrice(order.payAmount)[0] }}</span><span class="price-dec">.{{ splitPrice(order.payAmount)[1] }}</span>
          </span>
        </div>
      </section>

      <!-- 订单信息 -->
      <section class="card">
        <h3 class="card-title">订单信息</h3>
        <van-cell
          v-for="row in order.infoRows"
          :key="row.label"
          :title="row.label"
          :value="row.value"
          :border="false"
        />
      </section>
    </div>

    <!-- 吸底操作栏 -->
    <div class="action-bar">
      <van-button round class="action-btn" @click="onAction('联系客服')">联系客服</van-button>
      <van-button round class="action-btn" @click="onAction('修改地址')">修改地址</van-button>
      <van-button round class="action-btn action-btn--primary" @click="onAction('提醒发货')">提醒发货</van-button>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'

const router = useRouter()

// 静态 mock 数据（与设计稿一致）
const order = reactive({
  notice: '商家已接单，预计 10月22日 前发货',
  statusTag: '待发货',
  statusTitle: '等待商家发货',
  statusDesc: '超时未发货将自动赔付 5 元无门槛券',
  remainTime: (23 * 3600 + 59 * 60 + 59) * 1000,
  address: {
    name: '张三',
    phone: '138****8888',
    detail: '上海市浦东新区世纪大道 100 号 8 楼'
  },
  goods: [
    {
      id: 1,
      thumb: '',
      title: '轻薄羽绒服 女款 短款连帽 冬季保暖外套',
      spec: '米白色 / M',
      price: 299,
      count: 1
    }
  ],
  priceRows: [
    { label: '商品总额', value: '¥299.00' },
    { label: '运费', value: '免运费' },
    { label: '店铺优惠', value: '-¥30.00' }
  ],
  payAmount: 269,
  infoRows: [
    { label: '订单编号', value: '2026102012345678' },
    { label: '下单时间', value: '2026-10-20 10:24' },
    { label: '支付方式', value: '微信支付' }
  ]
})

const splitPrice = (price) => Number(price).toFixed(2).split('.')

const goBack = () => router.back()

const onAction = (name) => showToast(name)
</script>

<style lang="less" scoped>
@text-primary: #1f2329;
@text-secondary: #646a73;
@text-tertiary: #8f959e;
@brand: #ff4d1a;
@warning: #ff8f1f;
@warning-light: #fff3e8;
@danger: #f53f3f;
@border: #e5e6eb;
@bg: #f5f6f7;

.order-detail-v2 {
  min-height: 100%;
  background: @bg;
  color: @text-primary;
  font-family: 'PingFang SC', 'Noto Sans SC', 'Helvetica Neue', Arial, sans-serif;

  // NavBar
  --van-nav-bar-height: 44px;
  --van-nav-bar-background: #fff;
  --van-nav-bar-title-font-size: 17px;
  --van-nav-bar-title-text-color: @text-primary;
  --van-nav-bar-icon-color: @text-primary;
  --van-nav-bar-arrow-size: 20px;

  // NoticeBar
  --van-notice-bar-height: 36px;
  --van-notice-bar-padding: 0 12px;
  --van-notice-bar-font-size: 12px;
  --van-notice-bar-line-height: 36px;
  --van-notice-bar-text-color: @warning;
  --van-notice-bar-background: @warning-light;

  // Tag
  --van-tag-padding: 0 6px;
  --van-tag-font-size: 12px;
  --van-tag-line-height: 20px;
  --van-tag-radius: 4px;

  // Cell
  --van-cell-background: transparent;
  --van-cell-vertical-padding: 12px;
  --van-cell-horizontal-padding: 0;
  --van-cell-font-size: 14px;
  --van-cell-line-height: 20px;
  --van-cell-text-color: @text-primary;
  --van-cell-value-color: @text-secondary;
  --van-cell-label-color: @text-tertiary;
  --van-cell-label-font-size: 12px;
  --van-cell-label-line-height: 18px;
  --van-cell-label-margin-top: 2px;
  --van-cell-right-icon-color: @text-tertiary;
  --van-cell-active-color: transparent;

  // CountDown
  --van-count-down-text-color: @danger;
  --van-count-down-font-size: 12px;
  --van-count-down-line-height: 18px;
}

:deep(.van-nav-bar__title) {
  font-weight: 600;
}

.content {
  padding: 12px 12px calc(92px + env(safe-area-inset-bottom));
}

.notice {
  border-radius: 8px;
  margin-bottom: 12px;
}

.card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  &:last-child {
    margin-bottom: 0;
  }
}

.card-title {
  margin: 0 0 12px;
  font-size: 15px;
  line-height: 22px;
  font-weight: 600;
  color: @text-primary;
}

.divider {
  height: 1px;
  background: @border;
  transform: scaleY(0.99);
}

// 状态区
.status-row {
  display: flex;
  align-items: center;
  height: 24px;
}
.status-tag {
  flex-shrink: 0;
  font-weight: 500;
}
.status-title {
  margin-left: 8px;
  font-size: 17px;
  line-height: 24px;
  font-weight: 600;
  color: @text-primary;
}
.countdown {
  margin-left: auto;
  flex-shrink: 0;
}
.status-desc {
  margin: 4px 0 12px;
  font-size: 12px;
  line-height: 18px;
  color: @text-tertiary;
}
.address-cell {
  padding-bottom: 0;
  :deep(.van-cell__right-icon) {
    font-size: 16px;
  }
}

// 商品卡
.product {
  display: flex;
  margin-bottom: 12px;
}
.product-thumb {
  flex-shrink: 0;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  background: #f2f3f5;
  overflow: hidden;
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
.product-info {
  flex: 1;
  min-width: 0;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
}
.product-title {
  font-size: 14px;
  line-height: 20px;
  color: @text-primary;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}
.product-spec {
  font-size: 12px;
  line-height: 18px;
  color: @text-tertiary;
}
.product-bottom {
  margin-top: auto;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}
.product-count {
  font-size: 12px;
  line-height: 18px;
  color: @text-tertiary;
}

// 价格：¥ + 整数 + 小数
.price {
  display: inline-flex;
  align-items: baseline;
  color: @text-primary;
  .price-symbol,
  .price-dec {
    font-size: 12px;
    line-height: 14px;
    font-weight: 600;
  }
  .price-int {
    font-size: 18px;
    line-height: 22px;
    font-weight: 700;
  }
  &--large {
    color: @brand;
    .price-symbol,
    .price-dec {
      font-size: 14px;
      line-height: 17px;
    }
    .price-int {
      font-size: 22px;
      line-height: 26px;
    }
  }
}

.total-row {
  margin-top: 4px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
.total-label {
  margin-right: 8px;
  font-size: 14px;
  line-height: 20px;
  color: @text-primary;
}

// 吸底操作栏
.action-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  height: 56px;
  padding: 0 16px;
  padding-bottom: env(safe-area-inset-bottom);
  box-sizing: content-box;
  background: #fff;
  border-top: 1px solid @border;
}
.action-btn {
  height: 32px;
  padding: 0 16px;
  margin-left: 8px;
  font-size: 14px;
  font-weight: 500;
  color: @text-primary;
  background: #fff;
  border: 1px solid @border;
  &--primary {
    color: #fff;
    background: @brand;
    border: none;
  }
}
</style>
