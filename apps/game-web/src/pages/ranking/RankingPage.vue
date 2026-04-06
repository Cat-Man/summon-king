<template>
  <section class="ranking-page">
    <div class="ranking-header">
      <p>排行榜</p>
      <h1>巅峰荣耀榜单</h1>
      <p>当前排行榜展示最近上榜的战力和更新时间，榜内玩家可获取额外奖励。</p>
    </div>

    <article v-if="selfEntry" class="self-summary">
      <p>我的排名</p>
      <h2>NO. {{ selfEntry.rank }}</h2>
      <span>战力 {{ selfEntry.score }} · 当前连胜 {{ selfEntry.arena_streak }} 场</span>
    </article>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div class="ranking-table">
      <article
        v-for="entry in leaderboard"
        :key="entry.player_id"
        :class="{ 'self-entry': entry.is_self }"
      >
        <span class="rank">NO. {{ entry.rank }}</span>
        <div>
          <h3>{{ entry.name }}</h3>
          <p>战力 {{ entry.score }} · 连胜 {{ entry.arena_streak }} · {{ formatUpdated(entry.updated) }}</p>
        </div>
      </article>
    </div>

    <div class="ranking-links">
      <RouterLink to="/home">返回首页</RouterLink>
      <RouterLink to="/arena">继续斗法</RouterLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import { getLeaderboard, type LeaderboardEntry } from "@/api/modules/ranking"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const leaderboard = ref<LeaderboardEntry[]>([])
const errorMessage = ref("")
const selfEntry = computed(() => leaderboard.value.find((entry) => entry.is_self) ?? null)

function formatUpdated(updated: number) {
  const diffMinutes = Math.max(0, Math.floor((Date.now() - updated) / 60_000))
  if (diffMinutes <= 0) {
    return "刚刚"
  }
  return `${diffMinutes} 分钟前`
}

async function loadLeaderboard() {
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
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadLeaderboard()
  },
  { flush: "sync" },
)

onMounted(async () => {
  await loadLeaderboard()
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

.self-summary {
  margin-top: 1.5rem;
  padding: 1.25rem 1.5rem;
  border-radius: 18px;
  border: 1px solid rgba(138, 227, 255, 0.25);
  background: rgba(138, 227, 255, 0.06);
}

.self-summary p,
.self-summary h2,
.self-summary span {
  margin: 0;
}

.self-summary p {
  color: #8ae3ff;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  font-size: 12px;
}

.self-summary h2 {
  margin-top: 0.4rem;
  font-size: 2rem;
}

.self-summary span {
  display: block;
  margin-top: 0.35rem;
  color: rgba(248, 251, 255, 0.72);
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

.ranking-table article.self-entry {
  border-color: rgba(138, 227, 255, 0.35);
  background: rgba(138, 227, 255, 0.08);
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

.ranking-links {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
}

.ranking-links a {
  display: inline-flex;
  padding: 0.75rem 1rem;
  border-radius: 999px;
  text-decoration: none;
  background: rgba(138, 227, 255, 0.08);
  color: #b9f0ff;
  border: 1px solid rgba(138, 227, 255, 0.25);
}
</style>
