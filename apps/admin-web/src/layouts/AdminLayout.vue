<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="admin-layout">
    <aside class="admin-layout__sidebar">运营后台</aside>
    <main class="admin-layout__main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 220px 1fr;
  background: #f5f7fa;
}

.admin-layout__sidebar {
  padding: 24px;
  background: #1f2d3d;
  color: #fff;
  font-weight: 700;
}

.admin-layout__main {
  padding: 24px;
}
</style>
