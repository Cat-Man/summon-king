<template>
  <section class="pagoda-page">
    <div class="pagoda-hero">
      <p class="pagoda-label">{{ status.label || "通天塔挑战" }}</p>
      <h1>登顶每一层，换取秘宝</h1>
      <p class="pagoda-subtitle">
        当前塔层与挑战次数已接入真实接口，先完成最小挑战闭环，后续再扩奖励详情与战斗展开。
      </p>
      <div class="pagoda-actions">
        <button class="primary" type="button" @click="startChallenge">开始挑战</button>
        <button class="ghost" type="button" @click="loadStatus">刷新状态</button>
      </div>
    </div>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div class="pagoda-grid">
      <article class="pagoda-card">
        <header>
          <strong>当前层级</strong>
          <span>第 {{ status.current_floor }} 层</span>
        </header>
        <p>最高 {{ status.max_floor }} 层，当前剩余挑战 {{ status.remaining_challenges }}/5。</p>
      </article>
      <article class="pagoda-card">
        <header>
          <strong>奖励预览</strong>
          <span>{{ status.reward_preview }}</span>
        </header>
        <p>最近挑战奖励：{{ lastReward || "尚未挑战" }}</p>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getTowerStatus, startTowerChallenge, type TowerStatus } from "@/api/modules/tower"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const status = ref<TowerStatus>({
  tower: "pagoda",
  label: "通天塔",
  current_floor: 0,
  max_floor: 10,
  remaining_challenges: 0,
  reward_preview: "",
})
const lastReward = ref("")
const errorMessage = ref("")

async function loadStatus() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载通天塔。"
    return
  }

  try {
    status.value = await getTowerStatus("pagoda", playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "通天塔状态加载失败。"
  }
}

async function startChallenge() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法挑战通天塔。"
    return
  }

  try {
    const result = await startTowerChallenge("pagoda", playerId)
    lastReward.value = result.reward
    await loadStatus()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "通天塔挑战失败。"
  }
}

onMounted(async () => {
  await loadStatus()
})
</script>

<style scoped>
:global(:root) {
  --pagoda-bg: radial-gradient(circle at top right, #fdf4f4, #fff5ef 50%, #f1e7d7);
  --pagoda-card: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(255, 230, 210, 0.95));
}

.pagoda-page {
  min-height: 100vh;
  padding: 2rem;
  background: var(--pagoda-bg);
  font-family: 'Playfair Display', 'STKaiti', serif;
  color: #1f1b1b;
}

.pagoda-hero {
  max-width: 720px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 28px;
  padding: 2.5rem;
  box-shadow: 0 40px 80px rgba(31, 27, 27, 0.12);
}

.pagoda-label {
  letter-spacing: 0.3em;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
  color: #c54b00;
  font-size: 0.85rem;
}

.pagoda-hero h1 {
  font-size: 2.6rem;
  margin-bottom: 0.75rem;
}

.pagoda-subtitle {
  font-size: 1rem;
  line-height: 1.8;
  color: #4c4238;
}

.pagoda-actions {
  margin-top: 1.5rem;
  display: flex;
  gap: 0.75rem;
}

.pagoda-actions button {
  border: none;
  border-radius: 999px;
  padding: 0.75rem 1.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.primary {
  background: #e34900;
  color: #fff;
}

.ghost {
  background: transparent;
  border: 1px solid #e34900;
  color: #e34900;
}

.pagoda-actions button:hover {
  transform: translateY(-2px);
}

.status-text {
  margin-top: 1rem;
  color: #b13c00;
}

.pagoda-grid {
  margin-top: 2.5rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.pagoda-card {
  background: var(--pagoda-card);
  padding: 1.4rem;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(31, 27, 27, 0.1);
  border: 1px solid rgba(197, 75, 0, 0.3);
}

.pagoda-card header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 1rem;
}

.pagoda-card p {
  color: #2b231f;
  line-height: 1.6;
  font-family: 'Source Sans Pro', 'PingFang SC', sans-serif;
}
</style>
