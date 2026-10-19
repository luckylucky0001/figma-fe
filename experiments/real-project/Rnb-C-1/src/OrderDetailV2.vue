<template>
  <div class="order-detail-v2" data-figma-id="1:1">
    <!-- NavBar → 复用业务组件 SimpleHeader（顶部返回 + 居中标题） -->
    <div class="od-header" data-figma-id="1:2">
      <s-header :name="'订单详情'"></s-header>
    </div>

    <!-- NoticeBar / Tone=Warning -->
    <van-notice-bar
      class="od-notice"
      data-figma-id="1:4"
      :text="order.notice"
      :scrollable="false"
      color="var(--ds-color-warning)"
      background="var(--ds-color-warning-light)"
    />

    <!-- Card / 状态（Title=Off） -->
    <van-cell-group inset class="od-card" data-figma-id="1:6">
      <div class="od-status" data-figma-id="1:7">
        <van-tag
          class="od-status__tag"
          data-figma-id="1:8"
          color="var(--ds-color-warning-light)"
          text-color="var(--ds-color-warning)"
        >{{ order.statusTag }}</van-tag>
        <span class="od-status__title" data-figma-id="1:10">{{ order.statusTitle }}</span>
        <van-count-down
          class="od-status__countdown"
          data-figma-id="1:11"
          :time="order.remainTime"
          format="剩余 HH:mm:ss"
        />
      </div>
      <p class="od-status__desc" data-figma-id="1:13">{{ order.statusDesc }}</p>
      <van-divider class="od-divider" data-figma-id="1:14" :hairline="false" />
      <van-cell
        data-figma-id="1:15"
        :title="`${order.address.name} ${order.address.phone}`"
        :label="order.address.detail"
        :border="false"
        is-link
        center
        @click="onEditAddress"
      />
    </van-cell-group>

    <!-- Card / 商品（Title=On） -->
    <van-cell-group inset class="od-card" data-figma-id="1:18">
      <h3 class="od-card__title" data-figma-id="1:19">商品信息</h3>
      <van-card
        v-for="goods in order.goodsList"
        :key="goods.id"
        class="od-goods"
        data-figma-id="1:20"
        :title="goods.title"
        :desc="goods.spec"
        :num="goods.count"
      >
        <template #thumb>
          <van-image
            class="od-goods__thumb"
            data-figma-id="1:21"
            :src="goods.thumb"
            width="80"
            height="80"
            radius="8"
            :show-loading="false"
          />
        </template>
        <template #price>
          <price data-figma-id="1:24" :value="goods.price" size="medium" tone="default" />
        </template>
      </van-card>
      <van-divider class="od-divider" data-figma-id="1:29" :hairline="false" />
      <van-cell
        v-for="row in order.amountRows"
        :key="row.figmaId"
        :data-figma-id="row.figmaId"
        :title="row.title"
        :value="row.value"
        :border="false"
      />
      <div class="od-total" data-figma-id="1:39">
        <span class="od-total__label" data-figma-id="1:40">实付款</span>
        <price data-figma-id="1:41" :value="order.payAmount" size="large" tone="primary" />
      </div>
    </van-cell-group>

    <!-- Card / 订单信息（Title=On） -->
    <van-cell-group inset class="od-card" data-figma-id="1:45">
      <h3 class="od-card__title" data-figma-id="1:46">订单信息</h3>
      <van-cell
        v-for="row in order.infoRows"
        :key="row.figmaId"
        :data-figma-id="row.figmaId"
        :title="row.title"
        :value="row.value"
        :border="false"
      />
    </van-cell-group>

    <!-- ActionBar：按钮右对齐 -->
    <van-action-bar class="od-action-bar" data-figma-id="1:56">
      <van-button data-figma-id="1:57" size="small" round @click="onContact">联系客服</van-button>
      <van-button data-figma-id="1:59" size="small" round @click="onEditAddress">修改地址</van-button>
      <van-button data-figma-id="1:61" size="small" round type="primary" @click="onRemind">提醒发货</van-button>
    </van-action-bar>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { showToast } from 'vant'
import sHeader from '@/components/SimpleHeader.vue'
import Price from '@/components/Price.vue'

// 静态 mock，与设计稿一致；后续接入 getOrderDetail 时替换
const order = reactive({
  notice: '商家已接单，预计 10月22日 前发货',
  statusTag: '待发货',
  statusTitle: '等待商家发货',
  statusDesc: '超时未发货将自动赔付 5 元无门槛券',
  remainTime: ((23 * 60 + 59) * 60 + 59) * 1000,
  address: {
    name: '张三',
    phone: '138****8888',
    detail: '上海市浦东新区世纪大道 100 号 8 楼'
  },
  goodsList: [
    {
      id: 1,
      title: '轻薄羽绒服 女款 短款连帽 冬季保暖外套',
      spec: '米白色 / M',
      price: 299,
      count: 1,
      thumb: ''
    }
  ],
  amountRows: [
    { figmaId: '1:30', title: '商品总额', value: '¥299.00' },
    { figmaId: '1:33', title: '运费', value: '免运费' },
    { figmaId: '1:36', title: '店铺优惠', value: '-¥30.00' }
  ],
  payAmount: 269,
  infoRows: [
    { figmaId: '1:47', title: '订单编号', value: '2026102012345678' },
    { figmaId: '1:50', title: '下单时间', value: '2026-10-20 10:24' },
    { figmaId: '1:53', title: '支付方式', value: '微信支付' }
  ]
})

const onEditAddress = () => showToast('修改地址')
const onContact = () => showToast('联系客服')
const onRemind = () => showToast('已提醒商家尽快发货')
</script>

<style lang="less" scoped>
.order-detail-v2 {
  // 设计稿 token 映射到 Vant 主题变量（仅作用于本页）
  --van-text-color: var(--ds-color-text-1);
  --van-text-color-2: var(--ds-color-text-3);
  --van-border-color: var(--ds-color-border);
  --van-cell-line-height: 20px;
  --van-cell-vertical-padding: 12px;
  --van-cell-value-color: var(--ds-color-text-2);
  --van-cell-label-margin-top: 2px;
  --van-cell-label-line-height: 18px;
  --van-cell-group-inset-padding: 0 12px;
  --van-cell-group-inset-radius: 12px;
  --van-divider-border-color: var(--ds-color-border);
  --van-divider-margin: 12px 16px 0;
  --van-card-padding: 0 16px;
  --van-card-background: transparent;
  --van-card-thumb-size: 80px;
  --van-card-num-color: var(--ds-color-text-3);

  min-height: 100vh;
  padding-bottom: calc(56px + 12px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  background: var(--ds-color-bg-page);
  color: var(--ds-color-text-1);
}

// SimpleHeader：标题按设计稿 17px / 600，右侧「更多」图标设计稿中没有
.od-header {
  :deep(.simple-header) {
    color: var(--ds-color-text-1);
  }
  :deep(.simple-header-name) {
    font-size: 17px;
    font-weight: 600;
  }
  :deep(.nbmore) {
    visibility: hidden;
  }
}

.od-notice {
  --van-notice-bar-height: 36px;
  --van-notice-bar-padding: 0 12px;
  --van-notice-bar-font-size: 12px;
  margin: 12px 12px 0;
  border-radius: 8px;
}

.od-card {
  margin-top: 12px;
  padding: 16px 0;
  background: var(--ds-color-bg-card);
  &__title {
    margin: 0 0 12px;
    padding: 0 16px;
    font-size: 15px;
    font-weight: 600;
    line-height: 22px;
    color: var(--ds-color-text-1);
  }
}

.od-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 16px;
  &__tag {
    --van-tag-padding: 0 6px;
    --van-tag-radius: 4px;
    --van-tag-line-height: 20px;
    --van-tag-font-size: 12px;
    font-weight: 500;
  }
  &__title {
    flex: 1;
    font-size: 17px;
    font-weight: 600;
    line-height: 24px;
    color: var(--ds-color-text-1);
  }
  &__countdown {
    --van-count-down-text-color: var(--ds-color-danger);
    --van-count-down-font-size: 12px;
    --van-count-down-line-height: 18px;
  }
  &__desc {
    margin: 4px 0 0;
    padding: 0 16px;
    font-size: 12px;
    line-height: 18px;
    color: var(--ds-color-text-3);
  }
}

.od-divider {
  border-width: 1px 0 0;
}

.od-goods {
  :deep(.van-card__thumb) {
    margin-right: 12px;
  }
  :deep(.van-card__title) {
    max-height: 40px;
    font-size: 14px;
    font-weight: 400;
    line-height: 20px;
  }
  :deep(.van-card__desc) {
    font-size: 12px;
    line-height: 18px;
  }
  :deep(.van-card__num) {
    font-size: 12px;
    line-height: 22px;
  }
  &__thumb {
    background: var(--ds-color-fill);
    overflow: hidden;
  }
}

.od-total {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  height: 32px;
  margin-top: 4px;
  padding: 0 16px;
  &__label {
    font-size: 14px;
    line-height: 20px;
    color: var(--ds-color-text-1);
  }
}

.od-action-bar {
  --van-action-bar-height: 56px;
  --van-button-small-height: 32px;
  --van-button-small-padding: 0 16px;
  --van-button-small-font-size: 14px;
  --van-button-default-border-color: var(--ds-color-border);
  --van-button-default-color: var(--ds-color-text-1);
  --van-button-primary-background: var(--ds-color-brand);
  --van-button-primary-border-color: var(--ds-color-brand);
  justify-content: flex-end;
  gap: 8px;
  padding: 0 16px;
  border-top: 1px solid var(--ds-color-border);
  box-sizing: border-box;
  .van-button {
    font-weight: 500;
  }
}
</style>
