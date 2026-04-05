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

    <article class="pagoda-card prebattle-card">
      <header>
        <strong>战斗前摘要</strong>
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

    <article v-if="rewardLines.length" class="pagoda-card reward-card">
      <header>
        <strong>奖励拆分</strong>
        <span>实时写回成长</span>
      </header>
      <ul>
        <li v-for="line in rewardLines" :key="line">{{ line }}</li>
      </ul>
    </article>

    <article v-if="battleSummary !== null" class="pagoda-card battle-card">
      <header>
        <strong>战斗摘要</strong>
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
  tower: "pagoda",
  label: "通天塔",
  current_floor: 0,
  max_floor: 10,
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
    errorMessage.value = "当前未登录，无法挑战通天塔。"
    return
  }

  try {
    const result = await startTowerChallenge("pagoda", playerId)
    lastReward.value = result.reward
    battleSummary.value = result.battle ? { ...result.battle } : null
    summarizeRewards(result.reward_delta)
    await loadStatus()
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "通天塔挑战失败。"
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

.prebattle-card {
  margin-top: 1.5rem;
}

.prebattle-metrics {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem 1rem;
  color: #2b231f;
  font-family: 'Source Sans Pro', 'PingFang SC', sans-serif;
}

.reward-card {
  margin-top: 1.5rem;
  border-color: rgba(255, 255, 255, 0.25);
}

.reward-card ul {
  margin: 0;
  padding-left: 20px;
  color: #1d1d1d;
}

.reward-card li {
  margin-bottom: 6px;
  font-weight: 600;
}

.battle-card dl {
  margin: 0.8rem 0 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
}

.battle-card dt {
  color: #7b614c;
  font-size: 0.85rem;
}

.battle-card dd {
  margin: 0.3rem 0 0;
  font-weight: 700;
}
</style>
