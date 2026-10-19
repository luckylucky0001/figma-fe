<!--
  订单详情 V2：按 Figma「OrderDetail 订单详情」(frame 1:1) 实现，数据为静态 mock。
  组件复用方案见 .figmafe/match/h5.md。
-->
<template>
  <div class="order-detail-v2" data-figma-id="1:1">
    <!-- Figma NavBar(1:2) → 项目顶部导航 SimpleHeader -->
    <s-header :name="'订单详情'"></s-header>

    <van-notice-bar
      class="od-notice"
      data-figma-id="1:4"
      :scrollable="false"
      :text="order.notice"
    />

    <!-- Card / 状态 -->
    <van-cell-group inset :border="false" class="od-card" data-figma-id="1:6">
      <div class="od-status">
        <div class="od-status__row" data-figma-id="1:7">
          <van-tag
            class="od-status__tag"
            data-figma-id="1:8"
            size="medium"
            color="var(--color-warning-light)"
            text-color="var(--color-warning)"
          >{{ order.statusTag }}</van-tag>
          <span class="od-status__title" data-figma-id="1:10">{{ order.statusTitle }}</span>
          <van-count-down
            class="od-status__countdown"
            data-figma-id="1:11"
            :time="order.remainTime"
            format="剩余 HH:mm:ss"
          />
        </div>
        <div class="od-status__desc" data-figma-id="1:13">{{ order.statusDesc }}</div>
      </div>
      <van-divider class="od-divider" data-figma-id="1:14" :hairline="false" />
      <van-cell
        data-figma-id="1:15"
        :title="`${order.address.name} ${order.address.phone}`"
        :label="order.address.detail"
        :border="false"
        is-link
      />
    </van-cell-group>

    <!-- Card / 商品 -->
    <van-cell-group inset :border="false" class="od-card" data-figma-id="1:18">
      <div class="od-card__title" data-figma-id="1:19">商品信息</div>
      <van-card
        v-for="goods in order.goods"
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
            width="100%"
            height="100%"
            radius="8"
            fit="cover"
            :src="goods.thumb"
            :show-loading="false"
            :show-error="false"
          />
        </template>
        <template #price>
          <price data-figma-id="1:24" :value="goods.price" size="medium" tone="default" />
        </template>
      </van-card>
      <van-divider class="od-divider" data-figma-id="1:29" :hairline="false" />
      <van-cell
        v-for="row in order.amounts"
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

    <!-- Card / 订单信息 -->
    <van-cell-group inset :border="false" class="od-card" data-figma-id="1:45">
      <div class="od-card__title" data-figma-id="1:46">订单信息</div>
      <van-cell
        v-for="row in order.infos"
        :key="row.figmaId"
        :data-figma-id="row.figmaId"
        :title="row.title"
        :value="row.value"
        :border="false"
      />
    </van-cell-group>

    <van-action-bar class="od-action-bar" data-figma-id="1:56">
      <van-button round size="small" data-figma-id="1:57" @click="onContact">联系客服</van-button>
      <van-button round size="small" data-figma-id="1:59" @click="onEditAddress">修改地址</van-button>
      <van-button round size="small" type="primary" data-figma-id="1:61" @click="onRemind">提醒发货</van-button>
    </van-action-bar>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { showToast } from 'vant'
import sHeader from '@/components/SimpleHeader.vue'
import Price from '@/components/Price.vue'

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
      title: '轻薄羽绒服 女款 短款连帽 冬季保暖外套',
      spec: '米白色 / M',
      price: 299,
      count: 1,
      thumb: ''
    }
  ],
  amounts: [
    { figmaId: '1:30', title: '商品总额', value: '¥299.00' },
    { figmaId: '1:33', title: '运费', value: '免运费' },
    { figmaId: '1:36', title: '店铺优惠', value: '-¥30.00' }
  ],
  payAmount: 269,
  infos: [
    { figmaId: '1:47', title: '订单编号', value: '2026102012345678' },
    { figmaId: '1:50', title: '下单时间', value: '2026-10-20 10:24' },
    { figmaId: '1:53', title: '支付方式', value: '微信支付' }
  ]
})

const onContact = () => showToast('联系客服')
const onEditAddress = () => showToast('修改地址')
const onRemind = () => showToast('已提醒商家发货')
</script>

<style lang="less" scoped>
  .order-detail-v2 {
    /* 把 Vant 组件变量映射到设计稿 token（只作用于本页） */
    --van-cell-group-inset-padding: 12px 12px 0;
    --van-cell-group-inset-radius: 12px;
    --van-cell-text-color: var(--color-text-1);
    --van-cell-value-color: var(--color-text-2);
    --van-cell-label-color: var(--color-text-3);
    --van-cell-font-size: 14px;
    --van-cell-line-height: 20px;
    --van-cell-vertical-padding: 12px;
    --van-cell-label-margin-top: 2px;
    --van-cell-right-icon-color: var(--color-text-3);
    --van-divider-border-color: var(--color-border-light);
    --van-divider-margin: 12px 16px 0;
    --van-notice-bar-height: 36px;
    --van-notice-bar-padding: 0 12px;
    --van-notice-bar-font-size: 12px;
    --van-notice-bar-text-color: var(--color-warning);
    --van-notice-bar-background: var(--color-warning-light);
    --van-tag-radius: 4px;
    --van-tag-font-size: 12px;
    --van-count-down-text-color: var(--color-danger);
    --van-count-down-font-size: 12px;
    --van-count-down-line-height: 18px;
    --van-card-padding: 0 16px;
    --van-card-background: transparent;
    --van-card-thumb-size: 80px;
    --van-card-thumb-radius: 8px;
    --van-card-desc-color: var(--color-text-3);
    --van-card-desc-line-height: 18px;
    --van-card-num-color: var(--color-text-3);
    --van-action-bar-height: 55px;
    --van-button-small-height: 32px;
    --van-button-small-padding: 0 16px;
    --van-button-small-font-size: 14px;
    --van-button-default-color: var(--color-text-1);
    --van-button-default-border-color: var(--color-border-light);
    --van-button-primary-background: var(--color-accent);
    --van-button-primary-border-color: var(--color-accent);

    min-height: 100vh;
    /* 36px 页尾留白 + 56px 吸底操作栏 + 安全区 */
    padding-bottom: calc(92px + env(safe-area-inset-bottom));
    background: var(--color-bg-page);
  }

  .od-notice {
    margin: 12px 12px 0;
    border-radius: 8px;
  }

  .od-card {
    padding: 16px 0;
    &__title {
      padding: 0 16px 12px;
      color: var(--color-text-1);
      font-size: 15px;
      font-weight: 600;
      line-height: 22px;
    }
  }

  .od-status {
    padding: 0 16px;
    &__row {
      display: flex;
      align-items: center;
      gap: 8px;
      height: 24px;
    }
    &__tag {
      font-weight: 500;
      line-height: 16px;
    }
    &__title {
      color: var(--color-text-1);
      font-size: 17px;
      font-weight: 600;
      line-height: 24px;
    }
    &__countdown {
      margin-left: auto;
    }
    &__desc {
      margin-top: 4px;
      color: var(--color-text-3);
      font-size: 12px;
      line-height: 18px;
    }
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
      color: var(--color-text-1);
    }
    :deep(.van-card__desc),
    :deep(.van-card__num) {
      font-size: 12px;
    }
    &__thumb {
      background: var(--color-fill-1);
    }
  }

  .od-total {
    display: flex;
    align-items: baseline;
    justify-content: flex-end;
    gap: 8px;
    height: 32px;
    margin-top: 4px;
    padding: 0 16px;
    &__label {
      color: var(--color-text-1);
      font-size: 14px;
      line-height: 20px;
    }
  }

  .od-action-bar {
    justify-content: flex-end;
    gap: 8px;
    padding: 0 16px;
    border-top: 1px solid var(--color-border-light);
    :deep(.van-button) {
      font-weight: 500;
    }
  }
</style>
