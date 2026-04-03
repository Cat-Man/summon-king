<template>
  <section class="map-page">
    <div class="map-hero">
      <p class="badge">参考策略</p>
      <h1>{{ title }}</h1>
      <p>探索城市、占领区域、训练幻兽前的全局视角。</p>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="map-grid">
      <article v-for="city in cities" :key="city.city_id" class="map-card">
        <header>
          <h2>{{ city.name }}</h2>
          <span>{{ city.region }}</span>
        </header>
        <p>坐标 {{ city.loc_x }} / {{ city.loc_y }}</p>
        <button class="ghost-btn" type="button">传送到 {{ city.name }}</button>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getWorldMap, type WorldMap } from "@/api/modules/dungeon"

const world = ref<WorldMap | null>(null)
const errorMessage = ref("")

const title = computed(() => (world.value ? `${world.value.name}世界地图` : "环天世界地图"))
const cities = computed(() => world.value?.cities ?? [])

onMounted(async () => {
  try {
    world.value = await getWorldMap()
  } catch (error) {
    if (error instanceof APIError) {
      errorMessage.value = error.message
      return
    }
    errorMessage.value = "世界地图加载失败，请稍后重试。"
  }
})
</script>

<style scoped>
:global(body) {
  background: radial-gradient(circle at 20% 20%, #1a2a4f, #050715 70%);
}
.map-page {
  min-height: 100vh;
  padding: 32px 24px;
  color: #f4f6ff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.status-text {
  margin: 0;
  color: #ffd3c7;
}
.map-hero {
  max-width: 640px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 20px;
  padding: 24px;
  backdrop-filter: blur(12px);
  box-shadow: 0 25px 45px rgba(10, 10, 40, 0.45);
}
.badge {
  display: inline-flex;
  padding: 6px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.15);
  margin-bottom: 8px;
  font-size: 12px;
}
.map-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}
.map-card {
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.02), rgba(255, 255, 255, 0.08));
  border: 1px solid rgba(255, 255, 255, 0.09);
  backdrop-filter: blur(8px);
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: transform 0.3s ease, border-color 0.3s ease;
}
.map-card:hover {
  transform: translateY(-6px);
  border-color: rgba(255, 255, 255, 0.3);
}
.map-card header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.map-card header span {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}
.ghost-btn {
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 10px;
  background: transparent;
  color: #f4f6ff;
  padding: 10px;
  cursor: pointer;
  transition: background 0.2s ease;
}
.ghost-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}
</style>
