<template>
  <section class="ranking-page">
    <div class="ranking-header">
      <p>排行榜</p>
      <h1>巅峰荣耀榜单</h1>
      <p>当前排行榜展示最近上榜的战力和更新时间，榜内玩家可获取额外奖励。</p>
    </div>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div class="ranking-table">
      <article v-for="entry in leaderboard" :key="entry.player_id">
        <span class="rank">NO. {{ entry.rank }}</span>
        <div>
          <h3>{{ entry.name }}</h3>
          <p>战力 {{ entry.score }} · {{ formatUpdated(entry.updated) }}</p>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getLeaderboard, type LeaderboardEntry } from "@/api/modules/ranking"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const leaderboard = ref<LeaderboardEntry[]>([])
const errorMessage = ref("")

function formatUpdated(updated: number) {
  const diffMinutes = Math.max(0, Math.floor((Date.now() - updated) / 60_000))
  if (diffMinutes <= 0) {
    return "刚刚"
  }
  return `${diffMinutes} 分钟前`
}

onMounted(async () => {
  if (!sessionStore.playerId) {
    errorMessage.value = "当前未登录，无法加载排行榜。"
    return
  }

  try {
    leaderboard.value = await getLeaderboard(sessionStore.playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "排行榜加载失败，请稍后重试。"
  }
})
</script>

<style scoped>
:global(:root) {
  --rank-bg: linear-gradient(180deg, #0d111b, #1b2035);
}

.ranking-page {
  min-height: 100vh;
  padding: 2.5rem;
  background: var(--rank-bg);
  color: #f8fbff;
  font-family: 'Sahitya', 'PingFang HK', serif;
}

.ranking-header h1 {
  font-size: 2.8rem;
  margin: 0.3rem 0;
}

.ranking-header p {
  max-width: 640px;
  color: rgba(248, 251, 255, 0.7);
}

.status-text {
  margin-top: 1rem;
  color: #ffb7b7;
}

.ranking-table {
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.ranking-table article {
  border: 1px solid rgba(248, 251, 255, 0.07);
  border-radius: 18px;
  padding: 1rem 1.5rem;
  background: rgba(248, 251, 255, 0.03);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rank {
  font-size: 1rem;
  letter-spacing: 0.2em;
  color: #8ae3ff;
  min-width: 80px;
}

.ranking-table h3 {
  margin: 0;
  font-size: 1.4rem;
}

.ranking-table p {
  margin: 0.2rem 0 0;
  color: rgba(248, 251, 255, 0.6);
}
</style>
