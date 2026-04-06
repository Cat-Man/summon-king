<template>
  <section class="spirit-tower-page">
    <div class="spirits-header">
      <p>{{ status.label || "战灵塔试炼" }}</p>
      <h1>灵息共鸣，魂魄觉醒</h1>
      <p class="hint">当前塔层、剩余次数与奖励预览已接入真实接口。</p>
      <div class="spirits-status">
        <span>当前层级：第 {{ status.current_floor }} 层</span>
        <span>可挑战次数：{{ status.remaining_challenges }}/5</span>
      </div>
      <div class="spirits-actions">
        <button class="primary" type="button" @click="startChallenge">开始挑战</button>
        <button class="ghost" type="button" @click="loadStatus">刷新状态</button>
      </div>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="spirits-list">
      <article>
        <h3>奖励预览</h3>
        <p>{{ status.reward_preview }}</p>
        <div class="guard-meta">
          <span>最高层数：{{ status.max_floor }}</span>
          <span>最近奖励：{{ lastReward || "尚未挑战" }}</span>
        </div>
      </article>
    </div>
    <article class="prebattle-card">
      <header class="battle-header">
        <h3>战斗前摘要</h3>
        <span>共享队伍</span>
      </header>
      <div class="prebattle-metrics">
        <span>当前队伍战力 {{ prebattleSummary.teamPower }}</span>
        <span>成长总加成 +{{ prebattleSummary.totalBonus }}</span>
        <span>战骨 +{{ prebattleSummary.boneBonus }}</span>
        <span>战灵 +{{ prebattleSummary.spiritBonus }}</span>
        <span>魔魂 +{{ prebattleSummary.soulBonus }}</span>
      </div>
    </article>
    <article v-if="rewardLines.length" class="spirits-list reward-card">
      <h3>奖励拆分</h3>
      <ul>
        <li v-for="line in rewardLines" :key="line">{{ line }}</li>
      </ul>
    </article>
    <article v-if="battleSummary !== null" class="spirits-list reward-card">
      <header class="battle-header">
        <h3>战斗摘要</h3>
        <span>{{ battleSummary.battle_type }}</span>
      </header>
      <dl class="battle-grid">
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
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"

import { APIError } from "@/api/http"
import { getPetCollection, type PetCollection } from "@/api/modules/pet"
import {
  getTowerStatus,
  startTowerChallenge,
  type TowerBattleSummary,
  type TowerRewardDelta,
  type TowerStatus,
} from "@/api/modules/tower"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"
import { summarizePetGrowth } from "@/utils/petGrowthSummary"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const status = ref<TowerStatus>({
  tower: "spirit",
  label: "战灵塔",
  current_floor: 0,
  max_floor: 12,
  remaining_challenges: 0,
  reward_preview: "",
})
const petCollection = ref<PetCollection | null>(null)
const lastReward = ref("")
const battleSummary = ref<TowerBattleSummary | null>(null)
const rewardLines = ref<string[]>([])
const errorMessage = ref("")
const prebattleSummary = computed(() => summarizePetGrowth(petCollection.value))

function summarizeRewards(delta: TowerRewardDelta) {
  const lines: string[] = []
  if (delta.bone_level) {
    lines.push(`战骨 ${delta.bone_level > 0 ? "+" : ""}${delta.bone_level}`)
  }
  if (delta.spirit_power) {
    lines.push(`灵力 ${delta.spirit_power > 0 ? "+" : ""}${delta.spirit_power}`)
  }
  if (delta.soul_pieces) {
    lines.push(`魔魂碎片 ${delta.soul_pieces > 0 ? "+" : ""}${delta.soul_pieces}`)
  }
  rewardLines.value = lines
}

async function loadStatus() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载战灵塔。"
    return
  }

  try {
    const nextStatus = await getTowerStatus("spirit", playerId)
    status.value = nextStatus
    lastReward.value = nextStatus.last_reward ?? ""
    summarizeRewards(nextStatus.last_reward_delta ?? {})
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战灵塔状态加载失败。"
  }
}

async function loadPetSummary() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    return
  }

  try {
    const nextCollection = await getPetCollection(playerId)
    if (nextCollection) {
      petCollection.value = nextCollection
    }
  } catch {}
}

async function startChallenge() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法挑战战灵塔。"
    return
  }

  try {
    const result = await startTowerChallenge("spirit", playerId)
    lastReward.value = result.reward
    battleSummary.value = result.battle ? { ...result.battle } : null
    summarizeRewards(result.reward_delta)
    await loadStatus()
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战灵塔挑战失败。"
  }
}

onMounted(async () => {
  await Promise.all([loadStatus(), loadPetSummary()])
})

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadPetSummary()
  },
  { flush: "sync" },
)
</script>

<style scoped>
:global(:root) {
  --spirit-gradient: linear-gradient(180deg, #0c0d2c, #26156a);
}

.spirit-tower-page {
  padding: 2rem;
  min-height: 100vh;
  background: radial-gradient(circle at 20% 20%, rgba(193, 198, 255, 0.2), transparent 60%), var(--spirit-gradient);
  color: #f6f5ff;
  font-family: 'Cinzel', 'STSong', serif;
}

.spirits-header {
  max-width: 640px;
  margin-bottom: 2rem;
}

.spirits-header p {
  margin: 0;
}

.spirits-header h1 {
  font-size: 2.8rem;
  margin-bottom: 0.5rem;
}

.hint {
  color: rgba(246, 245, 255, 0.75);
  line-height: 1.5;
}

.spirits-status {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
  font-family: 'Source Sans Pro', sans-serif;
  color: #d3d4ff;
}

.spirits-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
}

.primary,
.ghost {
  border-radius: 999px;
  padding: 0.75rem 1.2rem;
  font-weight: 600;
  cursor: pointer;
}

.primary {
  border: none;
  color: #180b2f;
  background: linear-gradient(135deg, #bfc5ff, #8a7cff);
}

.ghost {
  border: 1px solid rgba(255, 255, 255, 0.35);
  background: transparent;
  color: #fff;
}

.status-text {
  color: #ffd8f4;
}

.spirits-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.spirits-list article {
  background: rgba(255, 255, 255, 0.04);
  padding: 1.25rem;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.spirits-list h3 {
  margin-bottom: 0.4rem;
}

.guard-meta {
  margin-top: 0.8rem;
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
  color: #bbb8ff;
}

.reward-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.15);
  padding: 1.5rem;
  border-radius: 16px;
}

.prebattle-card {
  margin-top: 1rem;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.15);
  padding: 1.5rem;
  border-radius: 16px;
}

.prebattle-metrics {
  margin-top: 0.75rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem 1rem;
  color: #f6f5ff;
}

.reward-card ul {
  margin: 0.75rem 0 0;
  padding-left: 18px;
  color: #ffdce8;
}

.reward-card li {
  margin-bottom: 6px;
}

.battle-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.battle-header span {
  color: #bbb8ff;
  text-transform: uppercase;
}

.battle-grid {
  margin: 0.75rem 0 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
}

.battle-grid dt {
  color: #bbb8ff;
  font-size: 0.85rem;
}

.battle-grid dd {
  margin: 0.3rem 0 0;
  font-weight: 700;
}
</style>
