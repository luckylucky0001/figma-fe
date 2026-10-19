<!--
  价格：¥ + 整数 + 小数（对应 Figma 组件 Price，变体 Size=Large|Medium|Small、Tone=Primary|Default）
  用法：<price :value="269" size="large" tone="primary" />
-->
<template>
  <span class="price" :class="[`price--${size}`, `price--${tone}`]">
    <span class="price__symbol">{{ symbol }}</span>
    <span class="price__integer">{{ parts.integer }}</span>
    <span class="price__decimal">.{{ parts.decimal }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: {
    type: [Number, String],
    required: true
  },
  size: {
    type: String,
    default: 'medium',
    validator: (v) => ['large', 'medium', 'small'].includes(v)
  },
  tone: {
    type: String,
    default: 'default',
    validator: (v) => ['primary', 'default'].includes(v)
  },
  symbol: {
    type: String,
    default: '¥'
  }
})

const parts = computed(() => {
  const [integer, decimal] = Number(props.value || 0).toFixed(2).split('.')
  return { integer, decimal }
})
</script>

<style lang="less" scoped>
  .price {
    display: inline-flex;
    align-items: baseline;
    white-space: nowrap;
    &--default { color: var(--color-text-1); }
    &--primary { color: var(--color-accent); }
    &__symbol,
    &__decimal { font-weight: 600; }
    &__integer { font-weight: 700; }
    &--small {
      .price__symbol, .price__decimal { font-size: 10px; line-height: 12px; }
      .price__integer { font-size: 14px; line-height: 17px; }
    }
    &--medium {
      .price__symbol, .price__decimal { font-size: 12px; line-height: 14px; }
      .price__integer { font-size: 18px; line-height: 22px; }
    }
    &--large {
      .price__symbol, .price__decimal { font-size: 14px; line-height: 17px; }
      .price__integer { font-size: 22px; line-height: 26px; }
    }
  }
</style>
