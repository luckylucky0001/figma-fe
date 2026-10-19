<template>
  <span :class="['price', `price--${size}`, `price--${tone}`]">
    <span class="price__symbol">{{ symbol }}</span>
    <span class="price__integer">{{ parts.integer }}</span>
    <span v-if="decimals > 0" class="price__decimal">.{{ parts.decimal }}</span>
  </span>
</template>

<script setup>
// 价格展示：¥ + 整数 + 小数（对应 Figma 组件 Price，变体 Size × Tone）
import { computed } from 'vue'

const props = defineProps({
  value: {
    type: [Number, String],
    required: true
  },
  // large | medium | small
  size: {
    type: String,
    default: 'medium'
  },
  // primary（品牌色）| default（正文色）
  tone: {
    type: String,
    default: 'default'
  },
  symbol: {
    type: String,
    default: '¥'
  },
  decimals: {
    type: Number,
    default: 2
  }
})

const parts = computed(() => {
  const [integer, decimal = ''] = Number(props.value).toFixed(props.decimals).split('.')
  return { integer, decimal }
})
</script>

<style lang="less" scoped>
.price {
  display: inline-flex;
  align-items: baseline;
  font-weight: 600;
  color: var(--ds-color-text-1);
  white-space: nowrap;
  &--primary {
    color: var(--ds-color-brand);
  }
  &__integer {
    font-weight: 700;
  }
  &--large {
    font-size: 14px;
    line-height: 26px;
    .price__integer {
      font-size: 22px;
    }
  }
  &--medium {
    font-size: 12px;
    line-height: 22px;
    .price__integer {
      font-size: 18px;
    }
  }
  &--small {
    font-size: 10px;
    line-height: 18px;
    .price__integer {
      font-size: 14px;
    }
  }
}
</style>
