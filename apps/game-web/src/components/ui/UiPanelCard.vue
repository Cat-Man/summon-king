<script setup lang="ts">
defineProps<{
  title?: string
  wide?: boolean
}>()
</script>

<template>
  <article class="ui-panel-card" :class="{ 'ui-panel-card--wide': wide }">
    <header v-if="title || $slots.actions" class="ui-panel-card__header">
      <h2 v-if="title" class="ui-panel-card__title">{{ title }}</h2>
      <div v-if="$slots.actions" class="ui-panel-card__actions">
        <slot name="actions" />
      </div>
    </header>
    <slot />
  </article>
</template>

<style scoped>
.ui-panel-card {
  padding: 20px;
  border-radius: 20px;
  background: #fff;
  box-shadow: 0 10px 30px rgb(15 23 42 / 8%);
}

.ui-panel-card--wide {
  grid-column: 1 / -1;
}

.ui-panel-card__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin-bottom: 10px;
}

.ui-panel-card__title {
  margin: 0;
  font-size: 18px;
  color: #1f2937;
}

.ui-panel-card__actions {
  display: flex;
  align-items: center;
}
</style>
