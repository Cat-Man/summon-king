<script setup lang="ts">
withDefaults(
  defineProps<{
    items: Array<{ label: string; value: string; subtext?: string }>
    tone?: 'neutral' | 'warm' | 'mint'
    minWidth?: string
  }>(),
  {
    tone: 'neutral',
    minWidth: '140px'
  }
)
</script>

<template>
  <div class="ui-stat-grid" :class="`ui-stat-grid--${tone}`" :style="{ '--ui-stat-grid-min': minWidth }">
    <div v-for="item in items" :key="item.label" class="ui-stat-grid__item">
      <span class="ui-stat-grid__label">{{ item.label }}</span>
      <strong class="ui-stat-grid__value">{{ item.value }}</strong>
      <small v-if="item.subtext" class="ui-stat-grid__subtext">{{ item.subtext }}</small>
    </div>
  </div>
</template>

<style scoped>
.ui-stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--ui-stat-grid-min), 1fr));
  gap: 14px;
}

.ui-stat-grid__item {
  padding: 12px;
  border-radius: 14px;
}

.ui-stat-grid--neutral .ui-stat-grid__item { background: #f8fafc; }
.ui-stat-grid--warm .ui-stat-grid__item { background: #fff7ed; }
.ui-stat-grid--mint .ui-stat-grid__item { background: #f0fdfa; }

.ui-stat-grid__label {
  display: block;
  margin-bottom: 8px;
  color: #6b7280;
  font-size: 13px;
}

.ui-stat-grid--warm .ui-stat-grid__label { color: #9a3412; }

.ui-stat-grid__value {
  display: block;
  font-size: 16px;
  color: #1f2937;
}

.ui-stat-grid__subtext {
  display: block;
  margin-top: 6px;
  color: #4b5563;
  font-size: 12px;
}
</style>
