<template>
  <section class="arena-page">
    <div class="arena-hero">
      <p class="arena-label">今日斗法</p>
      <h1>竞技场连胜试炼</h1>
      <p class="arena-subtitle">先打通最小真实闭环：查看今日战绩，并模拟一场胜负切磋。</p>
      <div class="arena-actions">
        <button class="primary" type="button" @click="battle(true)">切磋获胜</button>
        <button class="ghost" type="button" @click="battle(false)">切磋失利</button>
      </div>
    </div>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div class="arena-grid">
      <article>
        <h2>{{ record.current_streak }}</h2>
        <p>当前连胜</p>
      </article>
      <article>
        <h2>{{ record.last_win ? "胜" : "败" }}</h2>
        <p>上次结果</p>
      </article>
    </div>

    <article v-if="rewardLines.length" class="reward-card">
      <header>
        <h3>奖励拆分</h3>
        <span>已写回成长钱包</span>
      </header>
      <ul>
        <li v-for="line in rewardLines" :key="line">{{ line }}</li>
      </ul>
    </article>

    <article v-if="battleSummary !== null" class="battle-card">
      <header>
        <h3>战斗摘要</h3>
        <span>{{ battleSummary.battle_type }}</span>
      </header>
      <dl>
        <div>
          <dt>战斗结果</dt>
          <dd>{{ battleSummary.result }}</dd>
        </div>
        <div>
          <dt>回合数</dt>
          <dd>{{ battleSummary.rounds }}</dd>
        </div>
        <div>
          <dt>我方战力</dt>
          <dd>{{ battleSummary.attacker_power }}</dd>
        </div>
        <div>
          <dt>敌方战力</dt>
          <dd>{{ battleSummary.defender_power }}</dd>
        </div>
      </dl>
    </article>

    <div class="arena-links" v-if="rewardLines.length">
      <RouterLink to="/ranking">查看排行榜</RouterLink>
      <RouterLink to="/home">返回首页</RouterLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import {
  challengeArena,
  getArenaStatus,
  type ArenaBattleSummary,
  type ArenaRecord,
} from "@/api/modules/arena"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const record = ref<ArenaRecord>({
  player_id: 0,
  current_streak: 0,
  last_win: false,
})
const battleSummary = ref<ArenaBattleSummary | null>(null)
const rewardLines = ref<string[]>([])
const errorMessage = ref("")

function summarizeRewards(spiritPower: number, soulPieces: number) {
  const lines: string[] = []
  if (spiritPower) {
    lines.push(`灵力 +${spiritPower}`)
  }
  if (soulPieces) {
    lines.push(`魔魂碎片 +${soulPieces}`)
  }
  rewardLines.value = lines
}

async function loadStatus() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载竞技场。"
    return
  }

  try {
    record.value = await getArenaStatus(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "竞技场状态加载失败。"
  }
}

async function battle(won: boolean) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法发起切磋。"
    return
  }

  try {
    const result = await challengeArena(playerId, won)
    record.value = result.record
    battleSummary.value = result.battle ? { ...result.battle } : null
    summarizeRewards(result.reward_delta.spirit_power, result.reward_delta.soul_pieces)
    resourceSyncStore.touch()
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "竞技场切磋失败。"
  }
}

onMounted(async () => {
  await loadStatus()
})
</script>

<style scoped>
.arena-page {
  min-height: 100vh;
  padding: 2rem;
  background: radial-gradient(circle at top left, rgba(255, 145, 77, 0.18), transparent 30%),
    linear-gradient(180deg, #1b1110, #0d0a12 68%);
  color: #fff7ef;
}

.arena-hero {
  max-width: 760px;
  padding: 2rem;
  border-radius: 24px;
  border: 1px solid rgba(255, 247, 239, 0.1);
  background: rgba(255, 255, 255, 0.04);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.35);
}

.arena-label {
  margin: 0;
  letter-spacing: 0.3em;
  color: #ffb16f;
}

.arena-hero h1 {
  margin: 0.5rem 0;
  font-size: 2.6rem;
}

.arena-subtitle {
  color: rgba(255, 247, 239, 0.72);
  line-height: 1.7;
}

.arena-actions {
  margin-top: 1.25rem;
  display: flex;
  gap: 0.75rem;
}

.primary,
.ghost {
  border-radius: 999px;
  padding: 0.8rem 1.2rem;
  font-weight: 700;
  cursor: pointer;
}

.primary {
  border: none;
  background: linear-gradient(135deg, #ff9f43, #ff6b2c);
  color: #2b1200;
}

.ghost {
  border: 1px solid rgba(255, 247, 239, 0.28);
  background: transparent;
  color: #fff7ef;
}

.status-text {
  margin-top: 1rem;
  color: #ffd2cb;
}

.arena-grid {
  margin-top: 1.5rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.arena-grid article {
  padding: 1.5rem;
  border-radius: 18px;
  border: 1px solid rgba(255, 247, 239, 0.1);
  background: rgba(255, 255, 255, 0.04);
}

.arena-grid h2 {
  margin: 0;
  font-size: 2.2rem;
}

.arena-grid p {
  margin-top: 0.5rem;
  color: rgba(255, 247, 239, 0.72);
}

.reward-card {
  margin-top: 1rem;
  padding: 1.25rem 1.5rem;
  border-radius: 18px;
  border: 1px solid rgba(255, 247, 239, 0.12);
  background: rgba(255, 255, 255, 0.04);
}

.reward-card header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.reward-card h3,
.reward-card span {
  margin: 0;
}

.reward-card span {
  color: rgba(255, 247, 239, 0.65);
}

.reward-card ul {
  margin: 0.8rem 0 0;
  padding-left: 1.2rem;
}

.battle-card {
  margin-top: 1rem;
  padding: 1.25rem 1.5rem;
  border-radius: 18px;
  border: 1px solid rgba(255, 247, 239, 0.12);
  background: rgba(255, 255, 255, 0.04);
}

.battle-card header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.battle-card h3,
.battle-card span {
  margin: 0;
}

.battle-card span {
  color: rgba(255, 247, 239, 0.65);
  text-transform: uppercase;
}

.battle-card dl {
  margin: 0.8rem 0 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
}

.battle-card dt {
  color: rgba(255, 247, 239, 0.65);
  font-size: 0.85rem;
}

.battle-card dd {
  margin: 0.3rem 0 0;
  font-size: 1rem;
  font-weight: 700;
}

.arena-links {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
}

.arena-links a {
  display: inline-flex;
  padding: 0.75rem 1rem;
  border-radius: 999px;
  text-decoration: none;
  background: rgba(255, 177, 111, 0.12);
  color: #ffcfaa;
  border: 1px solid rgba(255, 177, 111, 0.28);
}
</style>
