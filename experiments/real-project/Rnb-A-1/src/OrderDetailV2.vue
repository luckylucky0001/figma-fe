<template>
  <div class="order-detail-v2">
    <van-nav-bar
      class="od-nav"
      title="订单详情"
      left-arrow
      fixed
      placeholder
      :border="false"
      @click-left="goBack"
    />

    <div class="od-body">
      <!-- 通知条 -->
      <van-notice-bar class="od-notice" :text="order.notice" :scrollable="false" />

      <!-- 状态 + 地址 -->
      <section class="od-card">
        <div class="status-row">
          <van-tag class="status-tag">{{ order.statusTag }}</van-tag>
          <span class="status-title">{{ order.statusTitle }}</span>
          <van-count-down class="status-countdown" :time="order.remainTime">
            <template #default="t">剩余 {{ pad(t.days * 24 + t.hours) }}:{{ pad(t.minutes) }}:{{ pad(t.seconds) }}</template>
          </van-count-down>
        </div>
        <p class="status-desc">{{ order.statusDesc }}</p>
        <div class="od-divider" />
        <van-cell
          class="od-cell od-cell--address"
          :title="`${order.address.name} ${order.address.phone}`"
          :label="order.address.detail"
          :border="false"
          is-link
          center
          @click="onAction('修改地址')"
        />
      </section>

      <!-- 商品信息 -->
      <section class="od-card">
        <h3 class="od-card__title">商品信息</h3>
        <div v-for="item in order.goods" :key="item.id" class="product">
          <div class="product__thumb">
            <img v-if="item.thumb" :src="item.thumb" alt="" />
          </div>
          <div class="product__info">
            <div class="product__title">{{ item.title }}</div>
            <div class="product__spec">{{ item.spec }}</div>
            <div class="product__bottom">
              <span class="price price--medium">
                <span class="price__symbol">¥</span><span class="price__int">{{ splitPrice(item.price).int }}</span><span class="price__dec">.{{ splitPrice(item.price).dec }}</span>
              </span>
              <span class="product__count">x{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="od-divider" />
        <van-cell class="od-cell" title="商品总额" :value="`¥${formatPrice(order.goodsAmount)}`" :border="false" />
        <van-cell class="od-cell" title="运费" :value="order.freight > 0 ? `¥${formatPrice(order.freight)}` : '免运费'" :border="false" />
        <van-cell class="od-cell" title="店铺优惠" :value="`-¥${formatPrice(order.discount)}`" :border="false" />
        <div class="total-row">
          <span class="total-row__label">实付款</span>
          <span class="price price--large">
            <span class="price__symbol">¥</span><span class="price__int">{{ splitPrice(payAmount).int }}</span><span class="price__dec">.{{ splitPrice(payAmount).dec }}</span>
          </span>
        </div>
      </section>

      <!-- 订单信息 -->
      <section class="od-card">
        <h3 class="od-card__title">订单信息</h3>
        <van-cell class="od-cell" title="订单编号" :value="order.orderNo" :border="false" />
        <van-cell class="od-cell" title="下单时间" :value="order.createTime" :border="false" />
        <van-cell class="od-cell" title="支付方式" :value="order.payType" :border="false" />
      </section>
    </div>

    <!-- 吸底操作栏 -->
    <div class="od-action-bar">
      <van-button class="od-btn" size="small" round @click="onAction('联系客服')">联系客服</van-button>
      <van-button class="od-btn" size="small" round @click="onAction('修改地址')">修改地址</van-button>
      <van-button class="od-btn od-btn--primary" size="small" round type="primary" @click="onAction('已提醒商家发货')">提醒发货</van-button>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'

const router = useRouter()

// 静态 mock，与设计稿一致
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
  goodsAmount: 299,
  freight: 0,
  discount: 30,
  orderNo: '2026102012345678',
  createTime: '2026-10-20 10:24',
  payType: '微信支付'
})

const payAmount = computed(() => order.goodsAmount + order.freight - order.discount)

const formatPrice = (n) => Number(n).toFixed(2)
const splitPrice = (n) => {
  const [int, dec] = formatPrice(n).split('.')
  return { int, dec }
}
const pad = (n) => String(n).padStart(2, '0')

const goBack = () => router.back()
const onAction = (msg) => showToast(msg)
</script>

<style lang="less" scoped>
// 设计稿 token
@text-1: #1f2329;
@text-2: #646a73;
@text-3: #8f959e;
@line: #e5e6eb;
@bg-page: #f5f6f7;
@bg-fill: #f2f3f5;
@brand: #ff4d1a;
@warning: #ff8f1f;
@warning-bg: #fff3e8;
@danger: #f53f3f;

.order-detail-v2 {
  min-height: 100%;
  background: @bg-page;
  font-family: 'Noto Sans SC', 'PingFang SC', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  color: @text-1;
  // 底部：操作栏 56 + 设计稿留白 36
  padding-bottom: calc(92px + env(safe-area-inset-bottom));
  box-sizing: border-box;

  // Vant 主题变量（仅作用于本页）
  --van-nav-bar-height: 44px;
  --van-nav-bar-title-font-size: 17px;
  --van-nav-bar-title-text-color: @text-1;
  --van-nav-bar-icon-color: @text-1;
  --van-nav-bar-arrow-size: 20px;

  --van-notice-bar-height: 36px;
  --van-notice-bar-padding: 0 12px;
  --van-notice-bar-font-size: 12px;
  --van-notice-bar-line-height: 36px;
  --van-notice-bar-text-color: @warning;
  --van-notice-bar-background: @warning-bg;

  --van-cell-vertical-padding: 12px;
  --van-cell-horizontal-padding: 0;
  --van-cell-font-size: 14px;
  --van-cell-line-height: 20px;
  --van-cell-text-color: @text-1;
  --van-cell-value-color: @text-2;
  --van-cell-background: transparent;
  --van-cell-active-color: transparent;
  --van-cell-label-font-size: 12px;
  --van-cell-label-line-height: 18px;
  --van-cell-label-margin-top: 2px;
  --van-cell-label-color: @text-3;
  --van-cell-right-icon-color: @text-3;
  --van-cell-icon-size: 16px;

  --van-tag-padding: 0 6px;
  --van-tag-font-size: 12px;
  --van-tag-line-height: 20px;
  --van-tag-radius: 4px;

  --van-count-down-text-color: @danger;
  --van-count-down-font-size: 12px;
  --van-count-down-line-height: 18px;

  --van-button-small-height: 32px;
  --van-button-small-padding: 0 16px;
  --van-button-small-font-size: 14px;
  --van-button-default-color: @text-1;
  --van-button-default-border-color: @line;
  --van-button-primary-background: @brand;
  --van-button-primary-border-color: @brand;
}

.od-body {
  padding: 12px 12px 0;
}

.od-notice {
  border-radius: 8px;
}

.od-card {
  margin-top: 12px;
  padding: 16px;
  background: #fff;
  border-radius: 12px;

  &__title {
    margin: 0 0 12px;
    font-size: 15px;
    line-height: 22px;
    font-weight: 600;
    color: @text-1;
  }
}

.od-divider {
  height: 1px;
  background: @line;
}

// 状态
.status-row {
  display: flex;
  align-items: center;
  height: 24px;
  gap: 8px;
}
.status-tag {
  flex-shrink: 0;
  font-weight: 500;
  color: @warning;
  background: @warning-bg;
}
.status-title {
  font-size: 17px;
  line-height: 24px;
  font-weight: 600;
  white-space: nowrap;
}
.status-countdown {
  margin-left: auto;
  white-space: nowrap;
}
.status-desc {
  margin: 4px 0 12px;
  font-size: 12px;
  line-height: 18px;
  color: @text-3;
}

// 单元格
.od-cell {
  :deep(.van-cell__value) {
    flex: none;
  }
  :deep(.van-cell__title) {
    flex: 1;
    min-width: 0;
  }
}
.od-cell--address {
  :deep(.van-cell__right-icon) {
    margin-left: 4px;
  }
}

// 商品
.product {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;

  &__thumb {
    flex-shrink: 0;
    width: 80px;
    height: 80px;
    border-radius: 8px;
    background: @bg-fill;
    overflow: hidden;
    img {
      display: block;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  &__info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
  }
  &__title {
    font-size: 14px;
    line-height: 20px;
    color: @text-1;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
  }
  &__spec {
    font-size: 12px;
    line-height: 18px;
    color: @text-3;
  }
  &__bottom {
    margin-top: auto;
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
  }
  &__count {
    font-size: 12px;
    line-height: 18px;
    color: @text-3;
  }
}

// 价格：¥ + 整数 + 小数
.price {
  display: inline-flex;
  align-items: baseline;
  font-weight: 600;
  white-space: nowrap;
  &__int {
    font-weight: 700;
  }

  &--medium {
    color: @text-1;
    .price__symbol,
    .price__dec {
      font-size: 12px;
      line-height: 14px;
    }
    .price__int {
      font-size: 18px;
      line-height: 22px;
    }
  }
  &--large {
    color: @brand;
    .price__symbol,
    .price__dec {
      font-size: 14px;
      line-height: 17px;
    }
    .price__int {
      font-size: 22px;
      line-height: 26px;
    }
  }
}

.total-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  height: 32px;
  margin-top: 4px;
  &__label {
    font-size: 14px;
    line-height: 20px;
    color: @text-1;
  }
}

// 吸底操作栏
.od-action-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  height: 56px;
  padding: 0 16px;
  padding-bottom: env(safe-area-inset-bottom);
  box-sizing: content-box;
  background: #fff;
  border-top: 1px solid @line;
}
.od-btn {
  font-weight: 500;
  &--primary {
    border: 0;
  }
}
</style>
