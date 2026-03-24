<script setup lang="ts">
withDefaults(
  defineProps<{
    items: string[]
    tone?: 'blue' | 'green' | 'orange'
    minWidth?: string
  }>(),
  {
    tone: 'blue',
    minWidth: '120px'
  }
)
</script>

<template>
  <div class="ui-chip-group" :class="`ui-chip-group--${tone}`" :style="{ '--ui-chip-min': minWidth }">
    <span v-for="item in items" :key="item" class="ui-chip-group__item">{{ item }}</span>
  </div>
</template>

<style scoped>
.ui-chip-group {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--ui-chip-min), 1fr));
  gap: 12px;
}

.ui-chip-group__item {
  padding: 12px 14px;
  border-radius: 14px;
  text-align: center;
}

.ui-chip-group--blue .ui-chip-group__item {
  background: #eff6ff;
  color: #1d4ed8;
}

.ui-chip-group--green .ui-chip-group__item {
  background: #dcfce7;
  color: #166534;
}

.ui-chip-group--orange .ui-chip-group__item {
  background: #fff7ed;
  color: #c2410c;
}
</style>
