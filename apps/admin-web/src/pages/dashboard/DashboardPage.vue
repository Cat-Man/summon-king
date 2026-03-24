<script setup lang="ts">
import { computed } from 'vue'
import * as echarts from 'echarts/core'

const chartEngine = echarts
const summaryCards = computed(() => [
  { label: '今日活跃', value: '0' },
  { label: '付费订单', value: '0' },
  { label: '待审核配置', value: '0' }
])
</script>

<template>
  <section class="dashboard-page">
    <header>
      <h1>运营后台</h1>
      <p>管理配置、监控数据与版本投放。</p>
    </header>

    <div class="dashboard-page__cards">
      <article v-for="card in summaryCards" :key="card.label" class="dashboard-page__card">
        <span>{{ card.label }}</span>
        <strong>{{ card.value }}</strong>
      </article>
    </div>

    <p class="dashboard-page__chart-note">ECharts 引擎已接入：{{ typeof chartEngine.use === 'function' ? 'ready' : 'pending' }}</p>
  </section>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dashboard-page__cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.dashboard-page__card {
  padding: 16px;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 4px 12px rgb(15 23 42 / 8%);
}

.dashboard-page__card strong {
  display: block;
  margin-top: 8px;
  font-size: 24px;
}

.dashboard-page__chart-note {
  color: #606266;
}
</style>
