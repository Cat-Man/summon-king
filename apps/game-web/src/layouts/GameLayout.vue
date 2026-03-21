<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="game-layout">
    <header class="game-layout__header">召唤之王</header>
    <main class="game-layout__content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.game-layout {
  min-height: 100vh;
  background: #f7f8fa;
}

.game-layout__header {
  padding: 16px;
  font-size: 20px;
  font-weight: 700;
  color: #222;
}

.game-layout__content {
  padding: 16px;
}
</style>
