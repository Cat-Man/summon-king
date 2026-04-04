<template>
  <section class="growth-page growth-page--manor">
    <header>
      <p class="tag">庄园</p>
      <h1>Elder Grove</h1>
      <p>已接入真实地块状态，当前共 {{ plots.length }} 个地块</p>
    </header>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="plots">
      <article v-for="plot in plots" :key="plot.plot_id">
        <strong>地块 {{ plot.plot_id }}</strong>
        <span>{{ plot.state }}</span>
      </article>
    </div>
    <button class="primary" type="button" @click="harvestNow">收获庄园</button>
    <div class="panel">
      <p>庄园扩建和收获动作后续再补，这轮先打通状态读取。</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getManorPlots, harvestManor, type ManorPlot } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const plots = ref<ManorPlot[]>([])
const errorMessage = ref("")

async function loadPlots() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载庄园。"
    return
  }

  try {
    plots.value = await getManorPlots(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "庄园状态加载失败。"
  }
}

async function harvestNow() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法收获庄园。"
    return
  }

  try {
    const result = await harvestManor(playerId)
    plots.value = result.plots
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "庄园收获失败。"
  }
}

onMounted(async () => {
  await loadPlots()
})
</script>

<style scoped>
.growth-page--manor {
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(150deg, #1d4028, #0a1b10);
  color: #f0ffec;
  min-height: 52vh;
}
.plots {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 24px 0;
}
.plots article {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  justify-content: space-between;
}
.status-text {
  margin-top: 1rem;
  color: #ffddb8;
}
.primary {
  padding: 12px 18px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #95d36a, #d4f48d);
  color: #12320f;
  font-weight: 700;
  cursor: pointer;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
}
.tag {
  font-size: 12px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
}
</style>
